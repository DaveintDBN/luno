# Roadmap

_Last updated: 2025-05-07 14:02:46_

## Phase 1 – Kick-off & Environment (1–2 days) ✅ Completed
1. Clone the Luno Go SDK and run the `_examples/readonly` demo.
2. Install Go (1.20+), configure `GOPATH`/`GOROOT`, export Luno API keys as env vars.
3. Verify basic API calls: tickers, order-book, balances.

## Phase 2 – Core Bot Architecture (2–3 days) ✅ Interfaces & config loader implemented
1. Define modules: **Client** (wrap `luno-go`), **Strategy** (signal gen), **Executor** (order placement), **State** (config store).
2. Sketch Go interfaces/types (e.g. `type Strategy interface { Next(...) Signal }`).

## Phase 3 – Config & Persistence (1–2 days) ✅ Completed
1. Choose storage (SQLite/Postgres/JSON).
2. Define config schema (pairs, thresholds, stake, cooldowns).
3. Build CLI or JSON loader to run headless with a config file.

## Phase 4 – Web Interface MVP (3–5 days) ✅ Completed
1. Pick Go web framework (Gin/Echo/Fiber).
2. Expose REST endpoints:
   - `GET /config`
   - `PUT /config`
   - `GET /status`
   - `GET /logs`
3. Scaffold front-end (vanilla JS or React/Vue): form for parameters, live feed via WebSocket/polling.

## Phase 5 – Strategy & Backtesting (2–4 days) ✅ Completed
1. Implement SMA crossover on XBTZAR.
2. Build backtester feeding historical bars to `Strategy`.
3. Validate metrics: drawdown, win rate.

## Phase 6 – Live Execution & Safety (2–3 days) ✅ Completed
1. Integrate `Executor` to place real orders via LunoClient.
2. Add risk controls: enforce `position_limit`, `max_drawdown`, and breakers.
3. Structured logging for orders, fills, and errors.

## Phase 7 – API & Observability (1–2 days) ✅ Completed
1. Setup Gin REST server and endpoints.
2. Implement `/healthz` health check.
3. Expose `/metrics` for Prometheus.
4. Add `/simulate` and `/execute` endpoints with UI buttons.
5. Register Prometheus counters and gauges.

## Phase 8 – Backtesting & Simulation Enhancements (1–2 days) ✅ Completed
1. Integrate fee calculations into backtester.
2. Extend backtester to use candle data.
3. Report PnL metrics in UI and logs.

## Phase 9 – Deployment & Monitoring (1–2 days) [in progress]
+- Dockerize Go service + front-end.
+- Deploy to cloud VM or Kubernetes.
+- Set up alerts (Slack/Email) on failures or P&L.
+- Test & implement Prometheus & Grafana monitoring and alerting setup. ✅

## Phase 10 – Iteration & Extension (ongoing)
 - Add new strategies (RSI, MACD, ML).
 - Support multiple pairs.
 - ✅ Implement Market Scanner UI: multi-select pairs, filters, start/stop controls, real-time results table & logs.
 - Enhance UI: real-time charts, order-book depth.
 - Build sandbox vs live test harness.
 - Continuous scanning & auto-execution: background market-wide scanning for profitable signals with optional live trade triggers.

## Phase 11 – Go-live & monitoring (1 day): start live trading on selected pairs, watch risk controls.

**Next steps:**
- Complete Phase 9 containerization & deployment (1 day).
- Continue Phase 10 UI enhancements & sandbox testing (2 days).
- Start implementing continuous scanning & auto-execution features (2 days).
- Phase 11 – Go-live & monitoring (1 day): start live trading on selected pairs, watch risk controls.
