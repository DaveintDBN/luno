package api

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"sort"
	"sync"
	"time"
	"math"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/luno/luno-bot/bot"
	"github.com/luno/luno-bot/config"
	luno "github.com/luno/luno-go"
)

// Metrics
var (
	simulateCounter    = prometheus.NewCounter(prometheus.CounterOpts{Name: "simulation_requests_total", Help: "Total simulation requests"})
	simulationPnLGauge = prometheus.NewGauge(prometheus.GaugeOpts{Name: "simulation_total_pnl", Help: "Latest simulation total PnL"})
	liveExecCounter    = prometheus.NewCounter(prometheus.CounterOpts{Name: "live_execute_requests_total", Help: "Total live execution requests"})
)

// SweepRequest defines parameters for market scanning.
type SweepRequest struct {
	Pairs          []string `json:"pairs"`
	MinVolume      float64  `json:"min_volume"`
	EntryThreshold float64  `json:"entry_threshold"`
	ExitThreshold  float64  `json:"exit_threshold"`
}

// SweepResult represents scan result for a market.
type SweepResult struct {
	Pair   string  `json:"pair"`
	Bid    float64 `json:"bid"`
	Ask    float64 `json:"ask"`
	Signal string  `json:"signal"`
}

// AutoScanRequest defines parameters for continuous market scanning.
type AutoScanRequest struct {
	Pairs           []string `json:"pairs"`
	MinVolume       float64  `json:"min_volume"`
	EntryThreshold  float64  `json:"entry_threshold"`
	ExitThreshold   float64  `json:"exit_threshold"`
	IntervalSeconds int      `json:"interval_seconds"`
	AutoExecute     bool     `json:"auto_execute"`
}

// OpportunityResult represents a top market opportunity.
type OpportunityResult struct {
	Pair              string  `json:"pair"`
	Bid               float64 `json:"bid"`
	Ask               float64 `json:"ask"`
	Potential         float64 `json:"potential"`
	Score             float64 `json:"score"`
	RecommendedStake float64 `json:"recommended_stake"`
}

// TopRequest defines parameters for top opportunities.
type TopRequest struct {
	Pairs     []string `json:"pairs"`
	MinVolume float64  `json:"min_volume"`
	Limit     int      `json:"limit"`
}

// ThresholdRequest defines params for grid-search backtest threshold optimization
type ThresholdRequest struct {
	Pairs        []string  `json:"pairs"`
	SinceMinutes int       `json:"since_minutes"`
	FeeRate      float64   `json:"fee_rate"`
	GridStart    float64   `json:"grid_start"`
	GridEnd      float64   `json:"grid_end"`
	GridStep     float64   `json:"grid_step"`
}

// ThresholdResult holds the best entry/exit thresholds per pair
type ThresholdResult struct {
	Pair           string  `json:"pair"`
	EntryThreshold float64 `json:"entry_threshold"`
	ExitThreshold  float64 `json:"exit_threshold"`
	TotalPnl       float64 `json:"total_pnl"`
	WinRate        float64 `json:"win_rate"`
}

// autoScanCancel manages the background auto-scan routine.
var autoScanCancel context.CancelFunc

var logsMu sync.Mutex
var logsBuffer []string

// Tracks consecutive entry threshold hits per pair for scan confirmation
var scanCountMu sync.Mutex
var scanConsecCount = make(map[string]int)

// computeRSI calculates the relative strength index over a period
func computeRSI(prices []float64, period int) float64 {
  if len(prices) < period+1 {
    return 50
  }
  gains, losses := 0.0, 0.0
  for i := 1; i < len(prices); i++ {
    delta := prices[i] - prices[i-1]
    if delta > 0 {
      gains += delta
    } else {
      losses -= delta
    }
  }
  if losses == 0 {
    return 100
  }
  rs := gains / losses
  return 100 - (100 / (1 + rs))
}

// computeStdDev calculates the standard deviation of price series
func computeStdDev(vals []float64) float64 {
  n := float64(len(vals))
  if n == 0 {
    return 0
  }
  sum := 0.0
  for _, v := range vals {
    sum += v
  }
  mean := sum / n
  var sdSum float64
  for _, v := range vals {
    diff := v - mean
    sdSum += diff * diff
  }
  return math.Sqrt(sdSum / n)
}

// computeEMA calculates exponential moving average series
func computeEMA(prices []float64, period int) []float64 {
  k := 2.0 / float64(period+1)
  ema := make([]float64, len(prices))
  ema[0] = prices[0]
  for i := 1; i < len(prices); i++ {
    ema[i] = prices[i]*k + ema[i-1]*(1-k)
  }
  return ema
}

// computeMACD returns MACD line, signal line, and histogram
func computeMACD(prices []float64, fastPeriod, slowPeriod, signalPeriod int) (macd, signal, hist []float64) {
  emaFast := computeEMA(prices, fastPeriod)
  emaSlow := computeEMA(prices, slowPeriod)
  macd = make([]float64, len(prices))
  for i := range prices {
    macd[i] = emaFast[i] - emaSlow[i]
  }
  signal = computeEMA(macd, signalPeriod)
  hist = make([]float64, len(prices))
  for i := range prices {
    hist[i] = macd[i] - signal[i]
  }
  return
}

// SetupRouter initializes REST endpoints for bot management.
func SetupRouter(store config.StateStore, client bot.Client, strat bot.Strategy, simExec, liveExec bot.Executor) *gin.Engine {
	// Register metrics safely (ignore already registered)
	for _, c := range []prometheus.Collector{simulateCounter, simulationPnLGauge, liveExecCounter} {
		if err := prometheus.Register(c); err != nil {
			if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
				panic(err)
			}
		}
	}
	r := gin.Default()

	// Log capture middleware
	r.Use(func(c *gin.Context) {
		c.Next()
		logsMu.Lock()
		logsBuffer = append(logsBuffer, fmt.Sprintf("%s %s %s", time.Now().Format(time.RFC3339), c.Request.Method, c.Request.URL.Path))
		if len(logsBuffer) > 500 {
			// keep last 500 entries
			logsBuffer = logsBuffer[len(logsBuffer)-500:]
		}
		logsMu.Unlock()
	})

	// Health check endpoint
	r.GET("/healthz", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	// Metrics endpoint for Prometheus
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Serve new React dashboard
	r.GET("/", func(c *gin.Context) {
		c.File("luno-trading-dashboard-pro-main/dist/index.html")
	})
	// Static assets
	r.Static("/assets", "luno-trading-dashboard-pro-main/dist/assets")
	r.StaticFile("/favicon.ico", "luno-trading-dashboard-pro-main/dist/favicon.ico")
	r.StaticFile("/robots.txt", "luno-trading-dashboard-pro-main/dist/robots.txt")
	// Fallback to index.html for client-side routing
	r.NoRoute(func(c *gin.Context) {
		c.File("luno-trading-dashboard-pro-main/dist/index.html")
	})

	// Get current config
	r.GET("/config", func(c *gin.Context) {
		cfg, err := store.LoadConfig()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, cfg)
	})

	// Update config
	r.PUT("/config", func(c *gin.Context) {
		var newCfg config.Config
		if err := c.ShouldBindJSON(&newCfg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := store.SaveConfig(&newCfg); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, newCfg)
	})

	// Bot status
	r.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "running"})
	})

	// Recent API logs
	r.GET("/logs", func(c *gin.Context) {
		logsMu.Lock()
		defer logsMu.Unlock()
		c.JSON(http.StatusOK, logsBuffer)
	})

	// Get available markets for scanner
	r.GET("/pairs", func(c *gin.Context) {
		resp, err := client.GetTickers(context.Background(), &luno.GetTickersRequest{Pair: []string{}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		pairs := make([]string, len(resp.Tickers))
		for i, t := range resp.Tickers {
			pairs[i] = t.Pair
		}
		c.JSON(http.StatusOK, pairs)
	})

	// Get account balances
	r.GET("/balances", func(c *gin.Context) {
		resp, err := client.GetBalances(context.Background(), &luno.GetBalancesRequest{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp.Balance)
	})

	// Percent change over the last hour for configured pair
	r.GET("/percent-change", func(c *gin.Context) {
		cfg, err := store.LoadConfig()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		since := time.Now().Add(-1 * time.Hour)
		resp, err := client.GetCandles(context.Background(), &luno.GetCandlesRequest{
			Duration: 3600,
			Pair:     cfg.Pair,
			Since:    luno.Time(since),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		candles := resp.Candles
		if len(candles) == 0 {
			c.JSON(http.StatusOK, gin.H{"percent_change": 0})
			return
		}
		first := candles[0]
		last := candles[len(candles)-1]
		change := (last.Close.Float64() - first.Open.Float64()) / first.Open.Float64() * 100
		c.JSON(http.StatusOK, gin.H{"percent_change": change})
	})

	// Scan selected markets with filters
	r.POST("/scan", func(c *gin.Context) {
		var req SweepRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Load bot config for RSI thresholds
		cfg, errLoad := store.LoadConfig()
		if errLoad != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errLoad.Error()})
			return
		}
		// Fetch current tickers
		resp, err := client.GetTickers(context.Background(), &luno.GetTickersRequest{Pair: req.Pairs})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var results []SweepResult
		scanCountMu.Lock()
		defer scanCountMu.Unlock()
		for _, t := range resp.Tickers {
			vol := t.Rolling24HourVolume.Float64()
			if req.MinVolume > 0 && vol < req.MinVolume {
				continue
			}
			// Volatility filter (Bollinger stddev)
			if cfg.BBPeriod > 1 && cfg.BBMultiplier > 0 {
				sinceBB := time.Now().Add(-time.Duration(cfg.BBPeriod+1) * time.Minute)
				bbRes, errBB := client.GetCandles(context.Background(), &luno.GetCandlesRequest{Pair: t.Pair, Duration: 60, Since: luno.Time(sinceBB)})
				if errBB == nil && len(bbRes.Candles) >= cfg.BBPeriod+1 {
					closes := make([]float64, cfg.BBPeriod+1)
					start := len(bbRes.Candles) - (cfg.BBPeriod + 1)
					if start < 0 { start = 0 }
					for i := start; i < len(bbRes.Candles); i++ {
						closes[i-start] = bbRes.Candles[i].Close.Float64()
					}
					sd := computeStdDev(closes)
					if t.Ask.Float64()-t.Bid.Float64() < sd*cfg.BBMultiplier {
						continue
					}
				}
			}
			// Depth-buffer check: sum top N bid volumes and compare to stake size
			ob, errOb := client.GetOrderBook(context.Background(), &luno.GetOrderBookRequest{Pair: t.Pair})
			if errOb == nil {
				levels := cfg.VWAPOrderbookDepthLevels
				if levels <= 0 { levels = 5 }
				totalDepth := 0.0
				for i, lvl := range ob.Bids {
					if i >= levels { break }
					totalDepth += lvl.Volume.Float64()
				}
				if cfg.StakeSize > 0 && totalDepth < cfg.StakeSize {
					continue
				}
			}
			bid := t.Bid.Float64()
			ask := t.Ask.Float64()
			// Update consecutive entry threshold hits
			hit := req.EntryThreshold > 0 && ask > bid*(1+req.EntryThreshold)
			if hit {
				scanConsecCount[t.Pair]++
			} else {
				scanConsecCount[t.Pair] = 0
			}
			sig := "hold"
			if scanConsecCount[t.Pair] >= 2 {
				// RSI confirmation
				rsiOK := true
				if cfg.RSIPeriod > 0 {
					since := time.Now().Add(-time.Duration(cfg.RSIPeriod+1) * time.Minute)
					candlesRes, errC := client.GetCandles(context.Background(), &luno.GetCandlesRequest{
						Pair:     t.Pair,
						Duration: 60,
						Since:    luno.Time(since),
					})
					if errC == nil && len(candlesRes.Candles) > cfg.RSIPeriod {
						// Prepare last period+1 closes
						closes := make([]float64, cfg.RSIPeriod+1)
						start := len(candlesRes.Candles) - (cfg.RSIPeriod + 1)
						if start < 0 {
							start = 0
						}
						for i := start; i < len(candlesRes.Candles); i++ {
							closes[i-start] = candlesRes.Candles[i].Close.Float64()
						}
						rsi := computeRSI(closes, cfg.RSIPeriod)
						if rsi > cfg.RSIOverBought {
							rsiOK = false
						}
					}
				}
				if rsiOK {
					sig = "buy"
				}
				// Multi-indicator: moving average crossover
				if sig == "buy" && cfg.ShortWindow > 0 && cfg.LongWindow > 0 {
					sinceMA := time.Now().Add(-time.Duration(cfg.LongWindow+1) * time.Minute)
					maRes, errMA := client.GetCandles(context.Background(), &luno.GetCandlesRequest{Pair: t.Pair, Duration: 60, Since: luno.Time(sinceMA)})
					if errMA == nil && len(maRes.Candles) >= cfg.LongWindow {
						closesMA := make([]float64, len(maRes.Candles))
						for i, cnd := range maRes.Candles {
							closesMA[i] = cnd.Close.Float64()
						}
						totalShort, totalLong := 0.0, 0.0
						startShort := len(closesMA) - cfg.ShortWindow
						if startShort < 0 { startShort = 0 }
						startLong := len(closesMA) - cfg.LongWindow
						for i := startShort; i < len(closesMA); i++ {
							totalShort += closesMA[i]
						}
						for i := startLong; i < len(closesMA); i++ {
							totalLong += closesMA[i]
						}
						smaShort := totalShort / float64(cfg.ShortWindow)
						smaLong := totalLong / float64(cfg.LongWindow)
						if smaShort <= smaLong {
							sig = "hold"
						}
					}
				}
				// MACD momentum filter
				if sig == "buy" && cfg.MACDFastPeriod > 0 && cfg.MACDSlowPeriod > 0 && cfg.MACDSignalPeriod > 0 {
					// fetch sufficient candles
					levels := cfg.MACDSlowPeriod
					if cfg.MACDSignalPeriod > levels { levels = cfg.MACDSignalPeriod }
					sinceMACD := time.Now().Add(-time.Duration(levels+1) * time.Minute)
					macdRes, errM := client.GetCandles(context.Background(), &luno.GetCandlesRequest{Pair: t.Pair, Duration: 60, Since: luno.Time(sinceMACD)})
					if errM == nil && len(macdRes.Candles) >= levels+1 {
						closesM := make([]float64, levels+1)
						start := len(macdRes.Candles) - (levels + 1)
						if start < 0 { start = 0 }
						for i := start; i < len(macdRes.Candles); i++ {
							closesM[i-start] = macdRes.Candles[i].Close.Float64()
						}
						_, _, hist := computeMACD(closesM, cfg.MACDFastPeriod, cfg.MACDSlowPeriod, cfg.MACDSignalPeriod)
						if hist[len(hist)-1] <= 0 {
							sig = "hold"
						}
					}
				}
			}
			results = append(results, SweepResult{Pair: t.Pair, Bid: bid, Ask: ask, Signal: sig})
		}
		c.JSON(http.StatusOK, results)
	})

	// Get top market opportunities ranked by potential % difference (ask–bid)/bid
	r.POST("/opportunities", func(c *gin.Context) {
		var req TopRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Load config for position sizing
		cfgRaw, errCfg := store.LoadConfig()
		if errCfg != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errCfg.Error()})
			return
		}
		cfg := cfgRaw
		if req.Limit <= 0 {
			req.Limit = 5
		}
		resp, err := client.GetTickers(context.Background(), &luno.GetTickersRequest{Pair: req.Pairs})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var ops []OpportunityResult
		for _, t := range resp.Tickers {
			vol := t.Rolling24HourVolume.Float64()
			if req.MinVolume > 0 && vol < req.MinVolume {
				continue
			}
			bid := t.Bid.Float64()
			ask := t.Ask.Float64()
			potential := (ask - bid) / bid * 100
			// Liquidity-weighted score
			ob, errOb := client.GetOrderBook(context.Background(), &luno.GetOrderBookRequest{Pair: t.Pair})
			topBidVol, topAskVol := 0.0, 0.0
			if errOb == nil {
				if len(ob.Bids) > 0 {
					topBidVol = ob.Bids[0].Volume.Float64()
				}
				if len(ob.Asks) > 0 {
					topAskVol = ob.Asks[0].Volume.Float64()
				}
			}
			liquidity := topBidVol + topAskVol
			weight := 1.0
			if liquidity > 0 {
				weight = math.Log(liquidity)
			}
			score := potential * weight
			// Determine recommended stake
			var recStake float64
			if cfg.PositionSizerType == "kelly" {
				k := cfg.KellyWinProb - (1-cfg.KellyWinProb)/cfg.KellyWinLossRatio
				recStake = cfg.InitialEquity * k
			} else {
				recStake = cfg.StakeSize
			}
			ops = append(ops, OpportunityResult{Pair: t.Pair, Bid: bid, Ask: ask, Potential: potential, Score: score, RecommendedStake: recStake})
		}
		// sort by descending score and limit
		sort.Slice(ops, func(i, j int) bool { return ops[i].Score > ops[j].Score })
		if len(ops) > req.Limit {
			ops = ops[:req.Limit]
		}
		c.JSON(http.StatusOK, ops)
	})

	// Stream top market opportunities as Server-Sent Events
	r.GET("/stream/opportunities", func(c *gin.Context) {
		// Load config for position sizing
		cfgRaw, errCfg := store.LoadConfig()
		if errCfg != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errCfg.Error()})
			return
		}
		cfg := cfgRaw
		pairsStr := c.Query("pairs")
		pairs := strings.Split(pairsStr, ",")
		minVol, _ := strconv.ParseFloat(c.Query("min_volume"), 64)
		limit, err := strconv.Atoi(c.Query("limit"))
		if err != nil || limit <= 0 {
			limit = 5
		}
		intervalSec, err := strconv.Atoi(c.Query("interval"))
		if err != nil || intervalSec <= 0 {
			intervalSec = 10
		}
		ticker := time.NewTicker(time.Duration(intervalSec) * time.Second)
		defer ticker.Stop()
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		for {
			<-ticker.C
			resp, err := client.GetTickers(context.Background(), &luno.GetTickersRequest{Pair: pairs})
			if err != nil {
				continue
			}
			var ops []OpportunityResult
			for _, t := range resp.Tickers {
				vol := t.Rolling24HourVolume.Float64()
				if minVol > 0 && vol < minVol {
					continue
				}
				bid := t.Bid.Float64()
				ask := t.Ask.Float64()
				potential := (ask - bid) / bid * 100
				ob, errOb := client.GetOrderBook(context.Background(), &luno.GetOrderBookRequest{Pair: t.Pair})
				topBidVol, topAskVol := 0.0, 0.0
				if errOb == nil {
					if len(ob.Bids) > 0 {
						topBidVol = ob.Bids[0].Volume.Float64()
					}
					if len(ob.Asks) > 0 {
						topAskVol = ob.Asks[0].Volume.Float64()
					}
				}
				liquidity := topBidVol + topAskVol
				weight := 1.0
				if liquidity > 0 {
					weight = math.Log(liquidity)
				}
				score := potential * weight
				// Determine recommended stake
				var recStake float64
				if cfg.PositionSizerType == "kelly" {
					k := cfg.KellyWinProb - (1-cfg.KellyWinProb)/cfg.KellyWinLossRatio
					recStake = cfg.InitialEquity * k
				} else {
					recStake = cfg.StakeSize
				}
				ops = append(ops, OpportunityResult{Pair: t.Pair, Bid: bid, Ask: ask, Potential: potential, Score: score, RecommendedStake: recStake})
			}
			c.SSEvent("opportunity", ops)
			c.Writer.Flush()
		}
	})

	// Start continuous auto-scan
	r.POST("/autoscan/start", func(c *gin.Context) {
		var req AutoScanRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if autoScanCancel != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "auto-scan already running"})
			return
		}
		ctx2, cancel := context.WithCancel(context.Background())
		autoScanCancel = cancel
		go func() {
			ticker := time.NewTicker(time.Duration(req.IntervalSeconds) * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ctx2.Done():
					return
				case <-ticker.C:
					resp, err := client.GetTickers(context.Background(), &luno.GetTickersRequest{Pair: req.Pairs})
					if err != nil {
						continue
					}
					for _, t := range resp.Tickers {
						vol := t.Rolling24HourVolume.Float64()
						if req.MinVolume > 0 && vol < req.MinVolume {
							continue
						}
						bid, ask := t.Bid.Float64(), t.Ask.Float64()
						signal := "hold"
						if req.EntryThreshold > 0 && ask > bid*(1+req.EntryThreshold) {
							signal = "buy"
						}
						if signal == "hold" && req.ExitThreshold > 0 && bid < ask*(1-req.ExitThreshold) {
							signal = "sell"
						}
						if req.AutoExecute && signal != "hold" {
							// load config and execute trade
							cfgRaw, err := store.LoadConfig()
							if err == nil {
								// build bot.Config from store Config
								botCfg := bot.Config{
									Pair:             t.Pair,
									EntryThreshold:   cfgRaw.EntryThreshold,
									ExitThreshold:    cfgRaw.ExitThreshold,
									StakeSize:        cfgRaw.StakeSize,
									Cooldown:         cfgRaw.Cooldown,
									PositionLimit:    cfgRaw.PositionLimit,
									MaxDrawdown:      cfgRaw.MaxDrawdown,
									ShortWindow:      cfgRaw.ShortWindow,
									LongWindow:       cfgRaw.LongWindow,
									BaseAccountId:    cfgRaw.BaseAccountId,
									CounterAccountId: cfgRaw.CounterAccountId,
								}
								// map string signal to bot.Signal
								var sigConst bot.Signal
								switch signal {
								case "buy":
									sigConst = bot.SignalBuy
								case "sell":
									sigConst = bot.SignalSell
								default:
									sigConst = bot.SignalNone
								}
								_ = liveExec.Execute(context.Background(), sigConst, bot.MarketData{Bid: bid, Ask: ask, Timestamp: time.Now()}, botCfg)
								liveExecCounter.Inc()
							}
						}
					}
				}
			}
		}()
		c.JSON(http.StatusOK, gin.H{"status": "auto-scan started"})
	})

	// Stop continuous auto-scan
	r.POST("/autoscan/stop", func(c *gin.Context) {
		if autoScanCancel == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no auto-scan running"})
			return
		}
		autoScanCancel()
		autoScanCancel = nil
		c.JSON(http.StatusOK, gin.H{"status": "auto-scan stopped"})
	})

	// Fetch order book
	r.GET("/orderbook", func(c *gin.Context) {
		// allow selecting pair via query param
		pair := c.Query("pair")
		if pair == "" {
			cfg, err := store.LoadConfig()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			pair = cfg.Pair
		}
		ob, err := client.GetOrderBook(context.Background(), &luno.GetOrderBookRequest{Pair: pair})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, ob)
	})

	// Backtest historical candles
	r.POST("/backtest", func(c *gin.Context) {
		var req struct {
			Pair         string  `json:"pair"`
			SinceMinutes int     `json:"since_minutes"`
			Short        int     `json:"short"`
			Long         int     `json:"long"`
			FeeRate      float64 `json:"fee_rate"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		since := time.Now().Add(-time.Duration(req.SinceMinutes) * time.Minute)
		candlesRes, err := client.GetCandles(context.Background(), &luno.GetCandlesRequest{
			Pair:     req.Pair,
			Duration: 60,
			Since:    luno.Time(since),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		n := len(candlesRes.Candles)
		closes := make([]float64, n)
		times := make([]time.Time, n)
		for i, cnd := range candlesRes.Candles {
			closes[i] = cnd.Close.Float64()
			times[i] = time.Time(cnd.Timestamp)
		}
		strat := bot.NewSMAStrategy(req.Short, req.Long)
		var cfg bot.Config
		cfg.EntryThreshold = 0
		cfg.ExitThreshold = 0
		cfg.StakeSize = 1
		inPos := false
		entry := 0.0
		trades, wins, losses := 0, 0, 0
		pnlTotal := 0.0
		pnlHistory := make([]gin.H, 0, n)
		drawdownHistory := make([]gin.H, 0, n)
		var profits []float64
		var peak, maxDD float64
		for i := 0; i < n; i++ {
			price := closes[i]
			md := bot.MarketData{Bid: price, Ask: price, Timestamp: times[i]}
			sig := strat.Next(md, cfg)
			if sig == bot.SignalBuy && !inPos {
				entry = price
				inPos = true
			} else if sig == bot.SignalSell && inPos {
				profitGross := (price - entry) * cfg.StakeSize
				feeCost := (entry + price) * cfg.StakeSize * req.FeeRate
				profit := profitGross - feeCost
				pnlTotal += profit
				profits = append(profits, profit)
				trades++
				if profit > 0 {
					wins++
				} else {
					losses++
				}
				inPos = false
			}
			// track drawdown
			if pnlTotal > peak {
				peak = pnlTotal
			}
			dd := peak - pnlTotal
			if dd > maxDD {
				maxDD = dd
			}
			drawdownHistory = append(drawdownHistory, gin.H{"time": times[i], "drawdown": dd})
			pnlHistory = append(pnlHistory, gin.H{"time": times[i], "pnl": pnlTotal})
		}
		winRate := 0.0
		if trades > 0 {
			winRate = float64(wins) / float64(trades) * 100
		}
		avgPnl := 0.0
		if trades > 0 {
			avgPnl = pnlTotal / float64(trades)
		}
		// compute Sharpe ratio on trade profits
		var sharpe float64
		if len(profits) > 1 {
			mean := 0.0
			for _, p := range profits {
				mean += p
			}
			mean /= float64(len(profits))
			sumsq := 0.0
			for _, p := range profits {
				sumsq += (p - mean) * (p - mean)
			}
			std := math.Sqrt(sumsq / float64(len(profits)-1))
			if std > 0 {
				sharpe = mean / std * math.Sqrt(float64(len(profits)))
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"trades":      trades,
			"wins":        wins,
			"losses":      losses,
			"win_rate":    winRate,
			"total_pnl":   pnlTotal,
			"avg_pnl":     avgPnl,
			"sharpe":      sharpe,
			"max_drawdown": maxDD,
			"pnl_history":  pnlHistory,
			"drawdown_history": drawdownHistory,
		})
	})

	// Simulate one step
	r.POST("/simulate", func(c *gin.Context) {
		simulateCounter.Inc()
		cfgRaw, err := store.LoadConfig()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// convert config to bot.Config
		cfg := bot.Config{
			Pair:           cfgRaw.Pair,
			EntryThreshold: cfgRaw.EntryThreshold,
			ExitThreshold:  cfgRaw.ExitThreshold,
			StakeSize:      cfgRaw.StakeSize,
			Cooldown:       cfgRaw.Cooldown,
			PositionLimit:  cfgRaw.PositionLimit,
			MaxDrawdown:    cfgRaw.MaxDrawdown,
			ShortWindow:    cfgRaw.ShortWindow,
			LongWindow:     cfgRaw.LongWindow,
		}
		// fetch market data
		ob, err := client.GetOrderBook(context.Background(), &luno.GetOrderBookRequest{Pair: cfg.Pair})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// compute mid-price
		bid := ob.Bids[0].Price.Float64()
		ask := ob.Asks[0].Price.Float64()
		md := bot.MarketData{Bid: bid, Ask: ask, Timestamp: time.Now()}
		// strategy signal and execution
		sig := strat.Next(md, cfg)
		execErr := simExec.Execute(context.Background(), sig, md, cfg)
		simulationPnLGauge.Set(simExec.(*bot.SimulatedExecutor).TotalPnL)
		// build response
		resp := gin.H{
			"signal":                sig,
			"position":              simExec.(*bot.SimulatedExecutor).Position,
			"total_pnl":             simExec.(*bot.SimulatedExecutor).TotalPnL,
			"max_drawdown_exceeded": simExec.(*bot.SimulatedExecutor).MaxDrawdownExceeded,
			"error":                 nil,
		}
		if execErr != nil {
			resp["error"] = execErr.Error()
		}
		c.JSON(http.StatusOK, resp)
	})

	// Execute one step live
	r.POST("/execute", func(c *gin.Context) {
		liveExecCounter.Inc()
		cfgRaw, err := store.LoadConfig()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		cfg := bot.Config{
			Pair:             cfgRaw.Pair,
			EntryThreshold:   cfgRaw.EntryThreshold,
			ExitThreshold:    cfgRaw.ExitThreshold,
			StakeSize:        cfgRaw.StakeSize,
			Cooldown:         cfgRaw.Cooldown,
			PositionLimit:    cfgRaw.PositionLimit,
			MaxDrawdown:      cfgRaw.MaxDrawdown,
			ShortWindow:      cfgRaw.ShortWindow,
			LongWindow:       cfgRaw.LongWindow,
			BaseAccountId:    cfgRaw.BaseAccountId,
			CounterAccountId: cfgRaw.CounterAccountId,
		}
		ob, err := client.GetOrderBook(context.Background(), &luno.GetOrderBookRequest{Pair: cfg.Pair})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		bid := ob.Bids[0].Price.Float64()
		ask := ob.Asks[0].Price.Float64()
		md := bot.MarketData{Bid: bid, Ask: ask, Timestamp: time.Now()}
		sig := strat.Next(md, cfg)
		execErr := liveExec.Execute(context.Background(), sig, md, cfg)
		resp := gin.H{"signal": sig, "error": nil}
		if execErr != nil {
			resp["error"] = execErr.Error()
		}
		c.JSON(http.StatusOK, resp)
	})

	// Grid-based threshold optimization endpoint
	r.POST("/thresholds", func(c *gin.Context) {
		var req ThresholdRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		since := time.Now().Add(-time.Duration(req.SinceMinutes) * time.Minute)
		var results []ThresholdResult
		for _, pair := range req.Pairs {
			candlesRes, err := client.GetCandles(context.Background(), &luno.GetCandlesRequest{Pair: pair, Duration: 60, Since: luno.Time(since)})
			if err != nil || len(candlesRes.Candles) == 0 {
				continue
			}
			n := len(candlesRes.Candles)
			closes := make([]float64, n)
			for i, cndl := range candlesRes.Candles {
				closes[i] = cndl.Close.Float64()
			}
			bestPnl := -math.MaxFloat64
			var be, bx, bwr float64
			for e := req.GridStart; e <= req.GridEnd; e += req.GridStep {
				for x := req.GridStart; x <= req.GridEnd; x += req.GridStep {
					inPos := false
					entry := 0.0
					pnlTotal := 0.0
					wins, trades := 0, 0
					for _, price := range closes {
						// simple threshold signals on price changes
						if !inPos && price > closes[0]*(1+e) {
							inPos = true; entry = price
						}
						if inPos && price < entry*(1-x) {
							profit := (price-entry) - (entry+price)*req.FeeRate
							pnlTotal += profit
							trades++
							if profit > 0 { wins++ }
							inPos = false
						}
					}
					wr := 0.0
					if trades > 0 { wr = float64(wins)/float64(trades)*100 }
					if pnlTotal > bestPnl { bestPnl, be, bx, bwr = pnlTotal, e, x, wr }
				}
			}
			results = append(results, ThresholdResult{Pair: pair, EntryThreshold: be, ExitThreshold: bx, TotalPnl: bestPnl, WinRate: bwr})
		}
		c.JSON(http.StatusOK, results)
	})

	return r
}
