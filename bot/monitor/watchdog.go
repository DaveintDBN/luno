package monitor

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// HealthStatus represents the health of a monitored service
type HealthStatus string

const (
	// StatusHealthy means the service is functioning normally
	StatusHealthy HealthStatus = "healthy"
	// StatusDegraded means the service is operational but has issues
	StatusDegraded HealthStatus = "degraded"
	// StatusUnhealthy means the service is not functioning correctly
	StatusUnhealthy HealthStatus = "unhealthy"
	// StatusCrashed means the service has completely crashed
	StatusCrashed HealthStatus = "crashed"
)

// ServiceHealth contains health information for a monitored service
type ServiceHealth struct {
	ServiceName      string
	Status           HealthStatus
	LastChecked      time.Time
	LastError        string
	RestartCount     int
	CPUUsage         float64
	MemoryUsageMB    float64
	ResponseTimeMs   float64
	IsAutoRecovering bool
}

// WatchdogListener is notified on watchdog events
type WatchdogListener interface {
	OnServiceStatusChange(service string, oldStatus, newStatus HealthStatus)
	OnServiceRestart(service string, restartCount int, reason string)
	OnResourceThresholdExceeded(service string, resourceType string, value float64, threshold float64)
}

// Watchdog monitors system health and auto-recovers from failures
type Watchdog struct {
	serviceHealthMap map[string]*ServiceHealth
	restartCommands  map[string]string
	healthCheckFuncs map[string]func() HealthStatus
	cpuThresholds    map[string]float64
	memoryThresholds map[string]float64
	checkInterval    time.Duration
	listeners        []WatchdogListener
	mutex            sync.RWMutex
	stopping         bool
	mainPID          int
	botBinaryPath    string
	lastPings        map[string]time.Time
	pingTimeouts     map[string]time.Duration
}

// NewWatchdog creates a watchdog service for monitoring
func NewWatchdog(interval time.Duration, mainPID int, botBinaryPath string) *Watchdog {
	w := &Watchdog{
		serviceHealthMap: make(map[string]*ServiceHealth),
		restartCommands:  make(map[string]string),
		healthCheckFuncs: make(map[string]func() HealthStatus),
		cpuThresholds:    make(map[string]float64),
		memoryThresholds: make(map[string]float64),
		checkInterval:    interval,
		listeners:        make([]WatchdogListener, 0),
		stopping:         false,
		mainPID:          mainPID,
		botBinaryPath:    botBinaryPath,
		lastPings:        make(map[string]time.Time),
		pingTimeouts:     make(map[string]time.Duration),
	}
	return w
}

// RegisterService adds a service to be monitored
func (w *Watchdog) RegisterService(name string, healthCheckFunc func() HealthStatus, 
	restartCmd string, cpuThreshold, memoryThreshold float64) {
	
	w.mutex.Lock()
	defer w.mutex.Unlock()
	
	w.serviceHealthMap[name] = &ServiceHealth{
		ServiceName:      name,
		Status:           StatusHealthy,
		LastChecked:      time.Now(),
		LastError:        "",
		RestartCount:     0,
		CPUUsage:         0,
		MemoryUsageMB:    0,
		ResponseTimeMs:   0,
		IsAutoRecovering: false,
	}
	
	w.healthCheckFuncs[name] = healthCheckFunc
	w.restartCommands[name] = restartCmd
	w.cpuThresholds[name] = cpuThreshold
	w.memoryThresholds[name] = memoryThreshold
}

// RegisterServiceWithPing registers a service that should ping the watchdog
func (w *Watchdog) RegisterServiceWithPing(name string, timeoutDuration time.Duration, 
	restartCmd string, cpuThreshold, memoryThreshold float64) {
	
	w.mutex.Lock()
	defer w.mutex.Unlock()
	
	w.serviceHealthMap[name] = &ServiceHealth{
		ServiceName:      name,
		Status:           StatusHealthy,
		LastChecked:      time.Now(),
		LastError:        "",
		RestartCount:     0,
		CPUUsage:         0,
		MemoryUsageMB:    0,
		ResponseTimeMs:   0,
		IsAutoRecovering: false,
	}
	
	w.lastPings[name] = time.Now()
	w.pingTimeouts[name] = timeoutDuration
	w.restartCommands[name] = restartCmd
	w.cpuThresholds[name] = cpuThreshold
	w.memoryThresholds[name] = memoryThreshold
	
	// Create a health check based on ping timeout
	w.healthCheckFuncs[name] = func() HealthStatus {
		w.mutex.RLock()
		lastPing, exists := w.lastPings[name]
		timeout := w.pingTimeouts[name]
		w.mutex.RUnlock()
		
		if !exists {
			return StatusUnhealthy
		}
		
		if time.Since(lastPing) > timeout {
			return StatusCrashed
		}
		
		return StatusHealthy
	}
}

// Ping signals that a service is still alive
func (w *Watchdog) Ping(serviceName string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	
	w.lastPings[serviceName] = time.Now()
}

// AddListener registers a new event listener
func (w *Watchdog) AddListener(listener WatchdogListener) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	
	w.listeners = append(w.listeners, listener)
}

// Start begins the monitoring process
func (w *Watchdog) Start() {
	go w.monitorLoop()
	log.Println("Watchdog service started")
}

// Stop ends the monitoring process
func (w *Watchdog) Stop() {
	w.mutex.Lock()
	w.stopping = true
	w.mutex.Unlock()
	log.Println("Watchdog service stopped")
}

// monitorLoop continuously checks service health
func (w *Watchdog) monitorLoop() {
	ticker := time.NewTicker(w.checkInterval)
	defer ticker.Stop()
	
	for {
		w.mutex.RLock()
		stopping := w.stopping
		w.mutex.RUnlock()
		
		if stopping {
			return
		}
		
		w.checkAllServices()
		<-ticker.C
	}
}

// checkAllServices performs health checks on all registered services
func (w *Watchdog) checkAllServices() {
	w.mutex.RLock()
	services := make([]string, 0, len(w.serviceHealthMap))
	for service := range w.serviceHealthMap {
		services = append(services, service)
	}
	w.mutex.RUnlock()
	
	for _, service := range services {
		w.checkServiceHealth(service)
	}
}

// checkServiceHealth checks if a service is healthy and takes action if needed
func (w *Watchdog) checkServiceHealth(serviceName string) {
	w.mutex.RLock()
	healthCheck, hasHealthCheck := w.healthCheckFuncs[serviceName]
	health, hasService := w.serviceHealthMap[serviceName]
	w.mutex.RUnlock()
	
	if !hasService || !hasHealthCheck {
		return
	}
	
	// Check process resource usage
	w.checkResourceUsage(serviceName)
	
	// Perform health check
	newStatus := healthCheck()
	
	w.mutex.Lock()
	oldStatus := health.Status
	health.Status = newStatus
	health.LastChecked = time.Now()
	w.mutex.Unlock()
	
	// Notify listeners if status changed
	if oldStatus != newStatus {
		for _, listener := range w.listeners {
			listener.OnServiceStatusChange(serviceName, oldStatus, newStatus)
		}
		
		log.Printf("[WATCHDOG] Service %s status changed: %s -> %s", 
			serviceName, oldStatus, newStatus)
	}
	
	// Take action based on health status
	if newStatus == StatusCrashed {
		w.restartService(serviceName, "Service crashed")
	} else if newStatus == StatusUnhealthy {
		w.restartService(serviceName, "Service unhealthy")
	}
}

// checkResourceUsage monitors CPU and memory usage
func (w *Watchdog) checkResourceUsage(serviceName string) {
	w.mutex.RLock()
	health, hasService := w.serviceHealthMap[serviceName]
	cpuThreshold := w.cpuThresholds[serviceName]
	memThreshold := w.memoryThresholds[serviceName]
	w.mutex.RUnlock()
	
	if !hasService {
		return
	}
	
	// In a real implementation, this would get actual resource usage
	// For this example, we'll simulate resource checks
	cpuUsage, memUsage := w.getServiceResourceUsage(serviceName)
	
	w.mutex.Lock()
	health.CPUUsage = cpuUsage
	health.MemoryUsageMB = memUsage
	w.mutex.Unlock()
	
	// Check if thresholds are exceeded
	if cpuUsage > cpuThreshold {
		for _, listener := range w.listeners {
			listener.OnResourceThresholdExceeded(serviceName, "CPU", cpuUsage, cpuThreshold)
		}
		
		log.Printf("[WATCHDOG] Service %s CPU usage (%.2f%%) exceeds threshold (%.2f%%)",
			serviceName, cpuUsage, cpuThreshold)
	}
	
	if memUsage > memThreshold {
		for _, listener := range w.listeners {
			listener.OnResourceThresholdExceeded(serviceName, "Memory", memUsage, memThreshold)
		}
		
		log.Printf("[WATCHDOG] Service %s memory usage (%.2f MB) exceeds threshold (%.2f MB)",
			serviceName, memUsage, memThreshold)
	}
}

// restartService attempts to restart a failed service
func (w *Watchdog) restartService(serviceName string, reason string) {
	w.mutex.Lock()
	
	health, hasService := w.serviceHealthMap[serviceName]
	restartCmd, hasRestartCmd := w.restartCommands[serviceName]
	
	if !hasService || !hasRestartCmd || health.IsAutoRecovering {
		w.mutex.Unlock()
		return
	}
	
	health.IsAutoRecovering = true
	health.RestartCount++
	restartCount := health.RestartCount
	health.LastError = reason
	
	w.mutex.Unlock()
	
	log.Printf("[WATCHDOG] Attempting to restart service %s: %s (attempt #%d)",
		serviceName, reason, restartCount)
	
	// Notify listeners
	for _, listener := range w.listeners {
		listener.OnServiceRestart(serviceName, restartCount, reason)
	}
	
	// Execute restart command
	go func() {
		success := w.executeRestartCommand(serviceName, restartCmd)
		
		w.mutex.Lock()
		health.IsAutoRecovering = false
		if success {
			health.Status = StatusHealthy
		}
		w.mutex.Unlock()
		
		if success {
			log.Printf("[WATCHDOG] Successfully restarted service %s", serviceName)
		} else {
			log.Printf("[WATCHDOG] Failed to restart service %s", serviceName)
		}
	}()
}

// executeRestartCommand runs the command to restart a service
func (w *Watchdog) executeRestartCommand(serviceName, command string) bool {
	if serviceName == "main" && w.mainPID > 0 {
		// Special handling for main process
		return w.restartMainProcess()
	}
	
	// For other processes, execute the restart command
	var cmd *exec.Cmd
	
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell", "-Command", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	return cmd.Run() == nil
}

// restartMainProcess restarts the main bot process
func (w *Watchdog) restartMainProcess() bool {
	// Kill the current process
	if w.mainPID > 0 {
		proc, err := os.FindProcess(w.mainPID)
		if err == nil {
			proc.Kill()
		}
	}
	
	// Start a new instance of the bot
	dir := filepath.Dir(w.botBinaryPath)
	binary := filepath.Base(w.botBinaryPath)
	
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", "start", "/B", binary)
	} else {
		cmd = exec.Command("nohup", "./"+binary, "&")
	}
	
	cmd.Dir = dir
	return cmd.Start() == nil
}

// GetAllServicesHealth returns health info for all services
func (w *Watchdog) GetAllServicesHealth() map[string]ServiceHealth {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	
	result := make(map[string]ServiceHealth)
	for name, health := range w.serviceHealthMap {
		result[name] = *health // Return a copy
	}
	
	return result
}

// getServiceResourceUsage gets CPU and memory usage (simulated)
func (w *Watchdog) getServiceResourceUsage(serviceName string) (float64, float64) {
	// In production, this would use OS-specific ways to get process resource usage
	// For this simulation, we'll just return random-ish values
	
	// Simulate some resource usage between 0-100% CPU and 10-1000MB memory
	cpuUsage := 10.0 + float64(time.Now().Unix()%90)
	memUsage := 100.0 + float64(time.Now().Unix()%900)
	
	return cpuUsage, memUsage
}

// AIWatchdogListener is an AI-powered listener for the watchdog
type AIWatchdogListener struct {
	// This would be connected to the AI engine
	aiEngine interface{}
}

// NewAIWatchdogListener creates a new AI-powered watchdog listener
func NewAIWatchdogListener() *AIWatchdogListener {
	return &AIWatchdogListener{}
}

// OnServiceStatusChange is called when a service status changes
func (a *AIWatchdogListener) OnServiceStatusChange(service string, oldStatus, newStatus HealthStatus) {
	// In a real implementation, this would use the AI to analyze and respond to the status change
	log.Printf("[AI WATCHDOG] Service %s status changed from %s to %s - analyzing root cause...", 
		service, oldStatus, newStatus)
	
	// Simulate AI analysis (in reality this would use machine learning)
	if newStatus == StatusUnhealthy || newStatus == StatusCrashed {
		log.Printf("[AI WATCHDOG] Analysis complete: Potential problems detected with service %s", service)
		log.Printf("[AI WATCHDOG] Suggested actions: Check network connectivity, verify API rate limits, inspect logs")
	}
}

// OnServiceRestart is called when a service is restarted
func (a *AIWatchdogListener) OnServiceRestart(service string, restartCount int, reason string) {
	log.Printf("[AI WATCHDOG] Service %s restarted (attempt #%d) due to: %s", service, restartCount, reason)
	
	// Simulate AI-driven escalation logic
	if restartCount > 3 {
		log.Printf("[AI WATCHDOG] Multiple restart attempts detected - escalating to high priority")
		// In a real system, this could send alerts, notifications, etc.
	}
}

// OnResourceThresholdExceeded is called when a service exceeds resource limits
func (a *AIWatchdogListener) OnResourceThresholdExceeded(service string, resourceType string, value float64, threshold float64) {
	log.Printf("[AI WATCHDOG] Resource warning: Service %s is using %.2f%% %s (threshold: %.2f%%)",
		service, value, resourceType, threshold)
	
	// Simulate AI recommendations
	if resourceType == "CPU" {
		log.Printf("[AI WATCHDOG] Analysis: High CPU may indicate excessive trading frequency or infinite loop")
		log.Printf("[AI WATCHDOG] Recommendation: Verify market conditions, check for high-frequency loops")
	} else if resourceType == "Memory" {
		log.Printf("[AI WATCHDOG] Analysis: Memory growth may indicate resource leak")
		log.Printf("[AI WATCHDOG] Recommendation: Check for unclosed connections or growing data structures")
	}
}
