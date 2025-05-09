package recovery

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

// ErrorSeverity represents how critical an error is
type ErrorSeverity int

const (
	SeverityLow ErrorSeverity = iota
	SeverityMedium
	SeverityHigh
	SeverityCritical
)

// ErrorContext holds information about a trading error
type ErrorContext struct {
	Timestamp   time.Time
	Message     string
	OrderID     string
	Pair        string
	Amount      float64
	Price       float64
	ErrorType   string
	Severity    ErrorSeverity
	Recoverable bool
	Retries     int
	MaxRetries  int
}

// RecoveryManager handles trading errors and recovery strategies
type RecoveryManager struct {
	activeErrors      map[string]*ErrorContext
	historicalErrors  []*ErrorContext
	retryStrategies   map[string]RetryStrategy
	recoveryListeners []RecoveryListener
	mutex             sync.RWMutex
	maxErrorsStored   int
}

// RetryStrategy defines how to handle retries for different error types
type RetryStrategy struct {
	MaxRetries       int
	BackoffMultiplier float64
	InitialWaitMs    int
}

// RecoveryListener gets notified of recovery events
type RecoveryListener interface {
	OnErrorDetected(ctx *ErrorContext)
	OnRecoveryAttempt(ctx *ErrorContext, attempt int)
	OnRecoverySuccess(ctx *ErrorContext)
	OnRecoveryFailed(ctx *ErrorContext)
}

// NewRecoveryManager creates a new recovery manager
func NewRecoveryManager() *RecoveryManager {
	rm := &RecoveryManager{
		activeErrors:     make(map[string]*ErrorContext),
		historicalErrors: make([]*ErrorContext, 0),
		retryStrategies:  make(map[string]RetryStrategy),
		maxErrorsStored:  1000,
	}

	// Set up default retry strategies
	rm.retryStrategies["api_timeout"] = RetryStrategy{MaxRetries: 5, BackoffMultiplier: 1.5, InitialWaitMs: 1000}
	rm.retryStrategies["insufficient_balance"] = RetryStrategy{MaxRetries: 3, BackoffMultiplier: 2.0, InitialWaitMs: 5000}
	rm.retryStrategies["rate_limit"] = RetryStrategy{MaxRetries: 8, BackoffMultiplier: 2.0, InitialWaitMs: 2000}
	rm.retryStrategies["market_closed"] = RetryStrategy{MaxRetries: 2, BackoffMultiplier: 5.0, InitialWaitMs: 10000}
	rm.retryStrategies["default"] = RetryStrategy{MaxRetries: 3, BackoffMultiplier: 2.0, InitialWaitMs: 3000}

	return rm
}

// RegisterListener adds a recovery event listener
func (rm *RecoveryManager) RegisterListener(listener RecoveryListener) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()
	rm.recoveryListeners = append(rm.recoveryListeners, listener)
}

// HandleError processes a new trading error
func (rm *RecoveryManager) HandleError(errType string, message string, orderID string, pair string, amount float64, price float64) (*ErrorContext, error) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	// Create error context
	ctx := &ErrorContext{
		Timestamp:   time.Now(),
		Message:     message,
		OrderID:     orderID,
		Pair:        pair,
		Amount:      amount,
		Price:       price,
		ErrorType:   errType,
		Recoverable: true,
		Retries:     0,
	}

	// Set severity based on error type
	switch errType {
	case "api_timeout", "rate_limit":
		ctx.Severity = SeverityLow
	case "insufficient_balance", "price_changed":
		ctx.Severity = SeverityMedium
	case "market_closed", "invalid_order":
		ctx.Severity = SeverityHigh
	case "exchange_error", "system_error":
		ctx.Severity = SeverityCritical
		// Critical errors might not be automatically recoverable
		if errType == "system_error" {
			ctx.Recoverable = false
		}
	default:
		ctx.Severity = SeverityMedium
	}

	// Assign retry strategy
	strategy, exists := rm.retryStrategies[errType]
	if !exists {
		strategy = rm.retryStrategies["default"]
	}
	ctx.MaxRetries = strategy.MaxRetries

	// Store error context
	errorKey := fmt.Sprintf("%s-%s", orderID, errType)
	rm.activeErrors[errorKey] = ctx

	// Notify listeners
	for _, listener := range rm.recoveryListeners {
		go listener.OnErrorDetected(ctx)
	}

	// Log the error
	log.Printf("[ERROR] Trading error: %s - %s (Order: %s, Pair: %s)", errType, message, orderID, pair)

	// If error is recoverable, initiate recovery
	if ctx.Recoverable {
		go rm.attemptRecovery(errorKey, strategy)
		return ctx, nil
	}

	return ctx, errors.New("non-recoverable error: " + message)
}

// attemptRecovery tries to recover from an error using the specified strategy
func (rm *RecoveryManager) attemptRecovery(errorKey string, strategy RetryStrategy) {
	rm.mutex.RLock()
	ctx, exists := rm.activeErrors[errorKey]
	rm.mutex.RUnlock()

	if !exists {
		return
	}

	for ctx.Retries < ctx.MaxRetries {
		// Calculate backoff time
		waitTime := float64(strategy.InitialWaitMs) * 
			pow(strategy.BackoffMultiplier, float64(ctx.Retries))
		time.Sleep(time.Duration(waitTime) * time.Millisecond)

		rm.mutex.Lock()
		ctx.Retries++
		rm.mutex.Unlock()

		// Notify listeners about retry attempt
		for _, listener := range rm.recoveryListeners {
			listener.OnRecoveryAttempt(ctx, ctx.Retries)
		}

		// Simulate recovery logic - in production this would call the actual trading API
		recoverySuccess := rm.executeRecoveryStrategy(ctx)

		if recoverySuccess {
			rm.mutex.Lock()
			delete(rm.activeErrors, errorKey)
			rm.historicalErrors = append(rm.historicalErrors, ctx)
			// Trim historical errors if needed
			if len(rm.historicalErrors) > rm.maxErrorsStored {
				rm.historicalErrors = rm.historicalErrors[1:]
			}
			rm.mutex.Unlock()

			// Notify listeners about success
			for _, listener := range rm.recoveryListeners {
				listener.OnRecoverySuccess(ctx)
			}

			log.Printf("[RECOVERY] Successfully recovered from %s error (Order: %s, Pair: %s, Attempts: %d)",
				ctx.ErrorType, ctx.OrderID, ctx.Pair, ctx.Retries)
			return
		}
	}

	// If we get here, all retries failed
	rm.mutex.Lock()
	delete(rm.activeErrors, errorKey)
	rm.historicalErrors = append(rm.historicalErrors, ctx)
	rm.mutex.Unlock()

	// Notify listeners about failure
	for _, listener := range rm.recoveryListeners {
		listener.OnRecoveryFailed(ctx)
	}

	log.Printf("[RECOVERY] Failed to recover from %s error after %d attempts (Order: %s, Pair: %s)",
		ctx.ErrorType, ctx.Retries, ctx.OrderID, ctx.Pair)
}

// executeRecoveryStrategy implements recovery logic for different error types
func (rm *RecoveryManager) executeRecoveryStrategy(ctx *ErrorContext) bool {
	switch ctx.ErrorType {
	case "api_timeout", "rate_limit":
		// Simply retry the same request
		return simulateAPIRetry(ctx.Retries)
	
	case "insufficient_balance":
		// Adjust order amount to available balance
		return simulateBalanceAdjustment(ctx)
	
	case "price_changed":
		// Update price to current market price
		return simulateUpdatePrice(ctx)
	
	case "market_closed":
		// Wait for market to open
		return simulateMarketStatusCheck()
	
	case "invalid_order":
		// Validate and fix order parameters
		return simulateOrderValidation(ctx)
	
	case "exchange_error":
		// Attempt alternate API endpoint or exchange
		return simulateAlternateEndpoint(ctx.Retries)
	
	default:
		// Generic recovery approach
		return simulateDefaultRecovery(ctx.Retries)
	}
}

// GetActiveErrors returns all currently active errors
func (rm *RecoveryManager) GetActiveErrors() []*ErrorContext {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()
	
	active := make([]*ErrorContext, 0, len(rm.activeErrors))
	for _, err := range rm.activeErrors {
		active = append(active, err)
	}
	return active
}

// GetErrorHistory returns historical error records
func (rm *RecoveryManager) GetErrorHistory() []*ErrorContext {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()
	
	// Return a copy to prevent modification
	history := make([]*ErrorContext, len(rm.historicalErrors))
	copy(history, rm.historicalErrors)
	return history
}

// ClearErrorHistory clears the error history
func (rm *RecoveryManager) ClearErrorHistory() {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()
	rm.historicalErrors = make([]*ErrorContext, 0)
}

// Helper for exponential calculations
func pow(a, b float64) float64 {
	result := 1.0
	for i := 0; i < int(b); i++ {
		result *= a
	}
	return result
}

// Simulation functions - these would be replaced with actual API calls in production

func simulateAPIRetry(attempt int) bool {
	// Higher attempt number has better chance of success
	return (attempt > 2) || (time.Now().UnixNano() % 2 == 0)
}

func simulateBalanceAdjustment(ctx *ErrorContext) bool {
	// Simulate adjusting the order to 80% of original amount
	ctx.Amount = ctx.Amount * 0.8
	return time.Now().UnixNano() % 4 != 0 // 75% success rate
}

func simulateUpdatePrice(ctx *ErrorContext) bool {
	// Simulate adjusting price by 1% in the direction of the error
	// In a real system, we'd fetch the current market price
	if ctx.Message == "Price too high" {
		ctx.Price = ctx.Price * 0.99
	} else {
		ctx.Price = ctx.Price * 1.01
	}
	return time.Now().UnixNano() % 5 != 0 // 80% success rate
}

func simulateMarketStatusCheck() bool {
	// Simulate checking if market is now open
	return time.Now().UnixNano() % 3 != 0 // 67% success rate
}

func simulateOrderValidation(ctx *ErrorContext) bool {
	// Simulate validating and fixing order parameters
	return time.Now().UnixNano() % 3 != 0 // 67% success rate
}

func simulateAlternateEndpoint(attempt int) bool {
	// Simulate trying alternate API endpoint
	return attempt > 1 && (time.Now().UnixNano() % 4 != 0) // Higher success rate with more attempts
}

func simulateDefaultRecovery(attempt int) bool {
	// Generic recovery with 50% success rate
	return time.Now().UnixNano() % 2 == 0
}
