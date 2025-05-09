# Luno Trading Bot Pro

An advanced trading bot for Luno cryptocurrency exchange with AI-enhanced features, multi-strategy support, and comprehensive monitoring.

## Quickstart
1. Ensure Docker & Docker Compose are installed.
2. Export your Luno API credentials:
   ```bash
   export API_KEY_ID=your_api_key_id
   export API_KEY_SECRET=your_api_key_secret
   ```
3. Start all services:
   ```bash
   docker-compose up -d --build
   ```
4. Access the UI and monitoring:
   - Dashboard: http://localhost:8081
   - Prometheus: http://localhost:9091
   - Grafana: http://localhost:3002 (user: `admin` / pass: `admin`)

## Server Deployment
Deploy on any Linux server (local or Afrihost) with Docker & Docker Compose:
1. SSH into your server.
2. Install Docker & Docker Compose:
   ```bash
   # Ubuntu/Debian
   sudo apt update && sudo apt install -y docker.io docker-compose
   # CentOS/Amazon Linux
   sudo yum install -y docker git
   sudo service docker start
   sudo usermod -aG docker $USER
   curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" \
       -o /usr/local/bin/docker-compose && sudo chmod +x /usr/local/bin/docker-compose
   ```
3. Clone the repo and enter directory:
   ```bash
   git clone https://github.com/yourusername/luno-bot.git
   cd luno-bot
   ```
4. Export your Luno API credentials:
   ```bash
   export API_KEY_ID=your_api_key_id
   export API_KEY_SECRET=your_api_key_secret
   ```
5. Start all services:
   ```bash
   docker-compose up -d --build
   ```
6. Access services:
   - Bot UI: http://<server-ip>:8081
   - Prometheus: http://<server-ip>:9091
   - Grafana: http://<server-ip>:3002 (admin/admin)

## Features

### Core Trading Features
- Multiple technical indicators: RSI, MACD, Bollinger Bands
- Multi-timeframe analysis for more accurate signals
- Position sizing options: Fixed size or Kelly Criterion
- Time-Weighted Average Price (TWAP) execution
- Comprehensive backtesting with performance analytics

### AI-Enhanced Trading
- AI-driven signal reinforcement for better entries/exits
- Sentiment analysis integration for market mood assessment
- Pattern recognition for chart formations
- Self-optimizing parameters based on performance
- Automatic hyperparameter tuning

### System Reliability
- Advanced error recovery system
- Watchdog service for system health monitoring
- Automatic crash recovery
- Resource usage monitoring to prevent system failures
- Comprehensive logging and diagnostics

### User Interface
- Modern, responsive dashboard
- Real-time feedback system with event tracking
- AI insights tab for trading recommendations
- Configuration management via UI
- Performance metrics visualization

## Architecture

The system consists of several components:

1. **Trading Bot Core**: Go-based engine handling strategies and execution
2. **AI Engine**: Advanced ML components for market analysis
3. **React Dashboard**: Modern UI for configuration and monitoring
4. **Monitoring Stack**: Prometheus and Grafana for metrics
5. **Recovery System**: Ensures trading continuity during failures
