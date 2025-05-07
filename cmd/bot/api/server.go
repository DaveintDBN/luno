package api

import (
	"context"
	"net/http"
	"time"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gin-gonic/gin"
	luno "github.com/luno/luno-go"
	"github.com/luno/luno-bot/bot"
	"github.com/luno/luno-bot/config"
)

// Metrics
var (
	simulateCounter = prometheus.NewCounter(prometheus.CounterOpts{Name: "simulation_requests_total", Help: "Total simulation requests"})
	simulationPnLGauge = prometheus.NewGauge(prometheus.GaugeOpts{Name: "simulation_total_pnl", Help: "Latest simulation total PnL"})
	liveExecCounter = prometheus.NewCounter(prometheus.CounterOpts{Name: "live_execute_requests_total", Help: "Total live execution requests"})
)

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

	// Health check endpoint
	r.GET("/healthz", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	// Metrics endpoint for Prometheus
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Serve front-end UI
	r.GET("/", func(c *gin.Context) {
		c.File("web/index.html")
	})
	r.Static("/assets", "web")

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

	// Fetch order book
	r.GET("/orderbook", func(c *gin.Context) {
		cfg, err := store.LoadConfig()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ob, err := client.GetOrderBook(context.Background(), &luno.GetOrderBookRequest{Pair: cfg.Pair})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, ob)
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
			"signal":           sig,
			"position":         simExec.(*bot.SimulatedExecutor).Position,
			"total_pnl":        simExec.(*bot.SimulatedExecutor).TotalPnL,
			"max_drawdown_exceeded": simExec.(*bot.SimulatedExecutor).MaxDrawdownExceeded,
			"error":            nil,
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
			Pair:           cfgRaw.Pair,
			EntryThreshold: cfgRaw.EntryThreshold,
			ExitThreshold:  cfgRaw.ExitThreshold,
			StakeSize:      cfgRaw.StakeSize,
			Cooldown:       cfgRaw.Cooldown,
			PositionLimit:  cfgRaw.PositionLimit,
			MaxDrawdown:    cfgRaw.MaxDrawdown,
			ShortWindow:    cfgRaw.ShortWindow,
			LongWindow:     cfgRaw.LongWindow,
			BaseAccountId:  cfgRaw.BaseAccountId,
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

	return r
}
