// Main JavaScript for config, orderbook, simulate, execute, and market scanner UI

// Panel collapse toggle
document.addEventListener('DOMContentLoaded', () => {
  document.querySelectorAll('.toggle-btn').forEach(btn => {
    btn.classList.add('bg-gray-200','hover:bg-gray-300','text-gray-600','px-2','py-1','rounded','ml-2');
    const panelContent = btn.closest('h2').nextElementSibling;
    btn.addEventListener('click', () => {
      panelContent.classList.toggle('hidden');
      btn.textContent = panelContent.classList.contains('hidden') ? '+' : '-';
    });
  });
});

// Tab navigation logic (deep-linking and shortcuts)
document.addEventListener('DOMContentLoaded', () => {
  const tabs = document.querySelectorAll('.tab-link');
  const panels = document.querySelectorAll('[role="tabpanel"]');
  function activateTab(tab) {
    tabs.forEach(t => {
      t.classList.remove('font-semibold','text-blue-600');
      t.setAttribute('aria-selected','false');
    });
    panels.forEach(p => {
      p.classList.add('hidden');
      p.setAttribute('aria-hidden','true');
    });
    tab.classList.add('font-semibold','text-blue-600');
    tab.setAttribute('aria-selected','true');
    const target = tab.dataset.target;
    const panel = document.getElementById(target);
    panel.classList.remove('hidden');
    panel.setAttribute('aria-hidden','false');
    window.history.pushState(null,'', `#${target}`);
  }
  tabs.forEach((tab, i) => {
    tab.addEventListener('click', e => { e.preventDefault(); activateTab(tab); });
  });
  const hash = window.location.hash.slice(1);
  const initial = [...tabs].find(t => t.dataset.target === hash) || tabs[0];
  activateTab(initial);
  document.addEventListener('keydown', e => {
    if (e.ctrlKey && !e.shiftKey && !e.altKey && !e.metaKey) {
      const idx = parseInt(e.key,10)-1;
      if (idx>=0 && idx<tabs.length) activateTab(tabs[idx]);
    }
  });
});

// Config section
const form = document.getElementById('configForm');
const entryErrEl = document.getElementById('entryThresholdError');
const exitErrEl = document.getElementById('exitThresholdError');
const stakeErrEl = document.getElementById('stakeSizeError');
const cooldownErrEl = document.getElementById('cooldownError');
const rsiPeriodError = document.getElementById('rsiPeriodError');
const rsiOverboughtError = document.getElementById('rsiOverboughtError');
const rsiOversoldError = document.getElementById('rsiOversoldError');
const initialEquityErrorEl = document.getElementById('initialEquityError');
const positionSizerTypeErrorEl = document.getElementById('positionSizerTypeError');
const kellyWinProbErrorEl = document.getElementById('kellyWinProbError');
const kellyWinLossRatioErrorEl = document.getElementById('kellyWinLossRatioError');
const twapSlicesErrorEl = document.getElementById('twapSlicesError');
const twapIntervalSecondsErrorEl = document.getElementById('twapIntervalSecondsError');
async function loadConfig() {
  const res = await fetch('/config');
  const cfg = await res.json();
  document.getElementById('entryThreshold').value = cfg.entry_threshold;
  document.getElementById('exitThreshold').value = cfg.exit_threshold;
  document.getElementById('stakeSize').value = cfg.stake_size;
  document.getElementById('cooldown').value = cfg.cooldown;
  document.getElementById('rsiPeriod').value = cfg.rsi_period || 14;
  document.getElementById('rsiOverbought').value = cfg.rsi_overbought || 70;
  document.getElementById('rsiOversold').value = cfg.rsi_oversold || 30;
  document.getElementById('initialEquity').value = cfg.initial_equity;
  document.getElementById('positionSizerType').value = cfg.position_sizer_type;
  document.getElementById('kellyWinProb').value = cfg.kelly_win_prob;
  document.getElementById('kellyWinLossRatio').value = cfg.kelly_win_loss_ratio;
  document.getElementById('twapSlices').value = cfg.twap_slices;
  document.getElementById('twapIntervalSeconds').value = cfg.twap_interval_seconds;
}
form.addEventListener('submit', async e => {
  e.preventDefault();
  const entry = parseFloat(document.getElementById('entryThreshold').value);
  const exit = parseFloat(document.getElementById('exitThreshold').value);
  const stake = parseFloat(document.getElementById('stakeSize').value);
  const initialEquity = parseFloat(document.getElementById('initialEquity').value);
  const positionSizerType = document.getElementById('positionSizerType').value;
  const kellyWinProb = parseFloat(document.getElementById('kellyWinProb').value);
  const kellyWinLossRatio = parseFloat(document.getElementById('kellyWinLossRatio').value);
  const twapSlices = parseInt(document.getElementById('twapSlices').value, 10);
  const twapIntervalSeconds = parseInt(document.getElementById('twapIntervalSeconds').value, 10);
  clearErrors();
  if (isNaN(entry) || entry <= 0) { entryErrEl.textContent = 'Invalid entry threshold'; appendLog('Invalid entry threshold'); return; }
  if (isNaN(exit) || exit <= 0) { exitErrEl.textContent = 'Invalid exit threshold'; appendLog('Invalid exit threshold'); return; }
  if (exit <= entry) { exitErrEl.textContent = 'Exit threshold must be > entry'; appendLog('Exit threshold must be > entry'); return; }
  if (isNaN(stake) || stake <= 0) { stakeErrEl.textContent = 'Invalid stake size'; appendLog('Invalid stake size'); return; }
  if (!document.getElementById('cooldown').value) { cooldownErrEl.textContent = 'Invalid cooldown'; appendLog('Invalid cooldown'); return; }
  const rsiPeriod = parseInt(document.getElementById('rsiPeriod').value, 10);
  const rsiOverbought = parseFloat(document.getElementById('rsiOverbought').value);
  const rsiOversold = parseFloat(document.getElementById('rsiOversold').value);
  if (isNaN(rsiPeriod) || rsiPeriod <= 0) { rsiPeriodError.textContent = 'Invalid RSI period'; appendLog('Invalid RSI period'); return; }
  if (isNaN(rsiOverbought) || rsiOverbought <= 0) { rsiOverboughtError.textContent = 'Invalid RSI overbought'; appendLog('Invalid RSI overbought'); return; }
  if (isNaN(rsiOversold) || rsiOversold <= 0) { rsiOversoldError.textContent = 'Invalid RSI oversold'; appendLog('Invalid RSI oversold'); return; }
  if (isNaN(initialEquity) || initialEquity <= 0) { initialEquityErrorEl.textContent = 'Invalid initial equity'; appendLog('Invalid initial equity'); return; }
  if (!['fixed','kelly'].includes(positionSizerType)) { positionSizerTypeErrorEl.textContent = 'Invalid position sizer type'; appendLog('Invalid position sizer type'); return; }
  if (positionSizerType === 'kelly') {
    if (isNaN(kellyWinProb) || kellyWinProb <= 0 || kellyWinProb >= 1) { kellyWinProbErrorEl.textContent = 'Invalid Kelly win probability'; appendLog('Invalid Kelly win probability'); return; }
    if (isNaN(kellyWinLossRatio) || kellyWinLossRatio <= 0) { kellyWinLossRatioErrorEl.textContent = 'Invalid Kelly win/loss ratio'; appendLog('Invalid Kelly win/loss ratio'); return; }
  }
  if (isNaN(twapSlices) || twapSlices < 1) { twapSlicesErrorEl.textContent = 'Invalid TWAP slices'; appendLog('Invalid TWAP slices'); return; }
  if (isNaN(twapIntervalSeconds) || twapIntervalSeconds < 1) { twapIntervalSecondsErrorEl.textContent = 'Invalid TWAP interval'; appendLog('Invalid TWAP interval'); return; }
  const newCfg = { entry_threshold: entry, exit_threshold: exit, stake_size: stake, cooldown: document.getElementById('cooldown').value,
                  initial_equity: initialEquity, position_sizer_type: positionSizerType,
                  kelly_win_prob: kellyWinProb, kelly_win_loss_ratio: kellyWinLossRatio,
                  twap_slices: twapSlices, twap_interval_seconds: twapIntervalSeconds,
                  rsi_period: rsiPeriod, rsi_overbought: rsiOverbought, rsi_oversold: rsiOversold };
  try {
    const res = await fetch('/config', { method: 'PUT', headers: {'Content-Type': 'application/json'}, body: JSON.stringify(newCfg) });
    if (!res.ok) { const errText = await res.text(); appendLog('Config update failed: ' + errText); showToast('Config update failed: ' + errText, 'error'); return; }
    await loadConfig();
    fetchOrderBook();
    appendLog('Configuration saved');
    showToast('Configuration saved', 'success');
  } catch (err) {
    appendLog('Config save error: ' + err.message);
  }
});

// clearErrors
function clearErrors() {
  entryErrEl.textContent = '';
  exitErrEl.textContent = '';
  stakeErrEl.textContent = '';
  cooldownErrEl.textContent = '';
  rsiPeriodError.textContent = '';
  rsiOverboughtError.textContent = '';
  rsiOversoldError.textContent = '';
  initialEquityErrorEl.textContent = '';
  positionSizerTypeErrorEl.textContent = '';
  kellyWinProbErrorEl.textContent = '';
  kellyWinLossRatioErrorEl.textContent = '';
  twapSlicesErrorEl.textContent = '';
  twapIntervalSecondsErrorEl.textContent = '';
}

// Order book section
const refreshBtn = document.getElementById('refresh');
const obPre = document.getElementById('orderbook');
async function fetchOrderBook() {
  const res = await fetch('/orderbook');
  const data = await res.json();
  obPre.textContent = JSON.stringify(data, null, 2);
  document.getElementById('orderbookTimestamp').textContent = new Date().toLocaleTimeString();
  // Update price chart
  const bids = data.bids;
  const asks = data.asks;
  if (bids.length && asks.length) {
    const mid = (parseFloat(bids[0].price) + parseFloat(asks[0].price)) / 2;
    priceData.labels.push(new Date().toLocaleTimeString());
    priceData.datasets[0].data.push(mid);
    priceChart.update();
  }
  // Update depth chart
  const depthChartBids = bids.map(bid => ({ x: parseFloat(bid.price), y: parseFloat(bid.volume) }));
  const depthChartAsks = asks.map(ask => ({ x: parseFloat(ask.price), y: parseFloat(ask.volume) }));
  depthChart.data.datasets[0].data = depthChartBids;
  depthChart.data.datasets[1].data = depthChartAsks;
  depthChart.update();
}
refreshBtn.addEventListener('click', fetchOrderBook);

// Auto-refresh controls
let orderbookIntervalId = null;
const autoToggle = document.getElementById('autoRefreshToggle');
const autoIntervalInput = document.getElementById('autoRefreshInterval');
autoToggle.addEventListener('change', () => {
  if (autoToggle.checked) {
    fetchOrderBook();
    orderbookIntervalId = setInterval(fetchOrderBook, autoIntervalInput.value * 1000);
  } else {
    clearInterval(orderbookIntervalId);
  }
});
autoIntervalInput.addEventListener('change', () => {
  if (orderbookIntervalId) {
    clearInterval(orderbookIntervalId);
    orderbookIntervalId = setInterval(fetchOrderBook, autoIntervalInput.value * 1000);
  }
});

// Simulate and execute
const simulateBtn = document.getElementById('simulate');
const executeBtn = document.getElementById('execute');
const simPre = document.getElementById('simulateResult');
const execPre = document.getElementById('executeResult');

// Sandbox mode toggle logic
const sandboxCheckbox = document.getElementById('sandboxMode');
sandboxCheckbox.addEventListener('change', updateSandboxMode);
function updateSandboxMode() {
  if (sandboxCheckbox.checked) {
    simulateBtn.disabled = false;
    executeBtn.disabled = true;
  } else {
    simulateBtn.disabled = true;
    executeBtn.disabled = false;
  }
}
updateSandboxMode();

simulateBtn.addEventListener('click', async () => {
  const res = await fetch('/simulate', { method: 'POST' });
  const data = await res.json();
  simPre.textContent = JSON.stringify(data, null, 2);
  showToast(`Simulated: Position ${data.position}, PnL ${data.total_pnl}`, 'success');
  const now = new Date().toLocaleTimeString();
  const row = document.createElement('tr');
  row.innerHTML = `<td class="px-2 py-1">${now}</td><td class="px-2 py-1">${data.position}</td><td class="px-2 py-1">${data.total_pnl}</td>`;
  document.getElementById('simHistoryBody').appendChild(row);
  if (lastPosition !== null) {
    const marker = { x: new Date(), y: data.total_pnl };
    if (data.position > lastPosition) pnlChart.data.datasets[1].data.push(marker);
    else if (data.position < lastPosition) pnlChart.data.datasets[2].data.push(marker);
  }
  lastPosition = data.position;
  pnlChart.update();
  // Update PnL chart
  pnlData.labels.push(new Date().toLocaleTimeString());
  pnlData.datasets[0].data.push(data.total_pnl);
  pnlChart.update();
  // Update KPI cards
  document.getElementById('currentPosition').textContent = data.position;
  document.getElementById('currentPnL').textContent = data.total_pnl;
  document.getElementById('drawdownExceeded').textContent = data.max_drawdown_exceeded;
});
executeBtn.addEventListener('click', async () => {
  const res = await fetch('/execute', { method: 'POST' });
  const data = await res.json();
  execPre.textContent = JSON.stringify(data, null, 2);
  showToast('Live execution complete', 'success');
});

// Market Scanner section
const pairsSelect = document.getElementById('pairs');
const backtestForm = document.getElementById('backtestForm');
const backtestPairSelect = document.getElementById('backtestPair');
const backtestSinceInput = document.getElementById('backtestSince');
const backtestShortInput = document.getElementById('backtestShort');
const backtestLongInput = document.getElementById('backtestLong');
const backtestFeeRateInput = document.getElementById('backtestFeeRate');
const backtestPairErrorEl = document.getElementById('backtestPairError');
const backtestSinceErrorEl = document.getElementById('backtestSinceError');
const backtestShortErrorEl = document.getElementById('backtestShortError');
const backtestLongErrorEl = document.getElementById('backtestLongError');
const backtestFeeRateErrorEl = document.getElementById('backtestFeeRateError');
const backtestMetrics = document.getElementById('backtestMetrics');
const minVolEl = document.getElementById('minVolume');
const entryEl = document.getElementById('scanEntryThreshold');
const exitEl = document.getElementById('scanExitThreshold');
const startBtn = document.getElementById('startScan');
const stopBtn = document.getElementById('stopScan');
const resultsTbody = document.getElementById('scanResults').querySelector('tbody');
const scanLogs = document.getElementById('scanLogs');
const liveLogs = document.getElementById('liveLogs');
const minVolErrEl = document.getElementById('minVolumeError');
const scanEntryErrEl = document.getElementById('scanEntryThresholdError');
const scanExitErrEl = document.getElementById('scanExitThresholdError');
function appendLog(msg) {
  const ts = new Date().toLocaleTimeString();
  liveLogs.textContent += `[${ts}] ${msg}\n`;
  liveLogs.scrollTop = liveLogs.scrollHeight;
}
function clearErrors() {
  entryErrEl.textContent = '';
  exitErrEl.textContent = '';
  stakeErrEl.textContent = '';
  cooldownErrEl.textContent = '';
  minVolErrEl.textContent = '';
  scanEntryErrEl.textContent = '';
  scanExitErrEl.textContent = '';
  rsiPeriodError.textContent = '';
  rsiOverboughtError.textContent = '';
  rsiOversoldError.textContent = '';
  initialEquityErrorEl.textContent = '';
  positionSizerTypeErrorEl.textContent = '';
  kellyWinProbErrorEl.textContent = '';
  kellyWinLossRatioErrorEl.textContent = '';
  twapSlicesErrorEl.textContent = '';
  twapIntervalSecondsErrorEl.textContent = '';
  backtestPairErrorEl.textContent = '';
  backtestSinceErrorEl.textContent = '';
  backtestShortErrorEl.textContent = '';
  backtestLongErrorEl.textContent = '';
  backtestFeeRateErrorEl.textContent = '';
}
let scanIntervalId = null;
const prevBtn = document.getElementById('prevPage');
const nextBtn = document.getElementById('nextPage');
const pageInfo = document.getElementById('pageInfo');
const pageSizeInput = document.getElementById('pageSize');
let scanResultsData = [];
let currentScanPage = 1;

function renderScanPage() {
  const pageSize = parseInt(pageSizeInput.value) || 10;
  const total = scanResultsData.length;
  const totalPages = Math.max(1, Math.ceil(total / pageSize));
  currentScanPage = Math.min(Math.max(1, currentScanPage), totalPages);
  const start = (currentScanPage - 1) * pageSize;
  const end = start + pageSize;
  resultsTbody.innerHTML = '';
  scanResultsData.slice(start, end).forEach(r => {
    const tr = document.createElement('tr');
    tr.classList.add('odd:bg-white','even:bg-gray-50','dark:odd:bg-gray-800','dark:even:bg-gray-700','hover:bg-gray-100','dark:hover:bg-gray-600');
    tr.style.backgroundColor = r.signal === 'buy' ? '#d4edda' : r.signal === 'sell' ? '#f8d7da' : '';
    ['pair','bid','ask','signal'].forEach(col => {
      const td = document.createElement('td');
      td.textContent = r[col];
      tr.appendChild(td);
    });
    resultsTbody.appendChild(tr);
  });
  pageInfo.textContent = `Page ${currentScanPage} / ${totalPages}`;
  prevBtn.disabled = currentScanPage === 1;
  nextBtn.disabled = currentScanPage === totalPages;
}

prevBtn.addEventListener('click', () => { currentScanPage--; renderScanPage(); });
nextBtn.addEventListener('click', () => { currentScanPage++; renderScanPage(); });
pageSizeInput.addEventListener('change', () => { currentScanPage = 1; renderScanPage(); });

async function loadPairs() {
  try {
    const res = await fetch('/pairs');
    const pairs = await res.json();
    pairs.forEach(pair => {
      const opt = document.createElement('option');
      opt.value = pair;
      opt.textContent = pair;
      pairsSelect.appendChild(opt);
      backtestPairSelect.appendChild(opt.cloneNode(true));
    });
  } catch (err) {
    scanLogs.textContent += `[${new Date().toLocaleTimeString()}] Error loading pairs: ${err.message}\n`;
    appendLog(`Error loading pairs: ${err.message}`);
  }
}

async function runScan() {
  clearErrors();
  if (!pairsSelect.selectedOptions.length) { appendLog('No pairs selected'); return; }
  const minVol = parseFloat(minVolEl.value);
  if (isNaN(minVol) || minVol < 0) { minVolErrEl.textContent = 'Invalid min volume'; appendLog('Invalid min volume'); return; }
  const entryThreshold = parseFloat(entryEl.value);
  if (isNaN(entryThreshold) || entryThreshold <= 0) { scanEntryErrEl.textContent = 'Invalid scan entry threshold'; appendLog('Invalid scan entry threshold'); return; }
  const exitThreshold = parseFloat(exitEl.value);
  if (isNaN(exitThreshold) || exitThreshold <= 0) { scanExitErrEl.textContent = 'Invalid scan exit threshold'; appendLog('Invalid scan exit threshold'); return; }
  if (exitThreshold <= entryThreshold) { scanExitErrEl.textContent = 'Scan exit threshold must be > entry'; appendLog('Scan exit threshold must be > entry'); return; }
  const selected = Array.from(pairsSelect.selectedOptions).map(o => o.value);
  const req = {
    pairs: selected,
    min_volume: minVol,
    entry_threshold: entryThreshold,
    exit_threshold: exitThreshold
  };
  try {
    const res = await fetch('/scan', {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify(req)
    });
    const data = await res.json();
    scanResultsData = data;
    currentScanPage = 1;
    renderScanPage();
    scanLogs.textContent += `[${new Date().toLocaleTimeString()}] Scan completed: ${data.length} results\n`;
    appendLog(`Scan completed: ${data.length} results`);
  } catch (err) {
    scanLogs.textContent += `[${new Date().toLocaleTimeString()}] Scan error: ${err.message}\n`;
    appendLog(`Scan error: ${err.message}`);
  }
}

startBtn.addEventListener('click', () => {
  startBtn.disabled = true;
  stopBtn.disabled = false;
  scanLogs.textContent += `[${new Date().toLocaleTimeString()}] Starting scan...\n`;
  appendLog('Starting scan...');
  runScan();
  scanIntervalId = setInterval(runScan, 5000);
});
stopBtn.addEventListener('click', () => {
  if (scanIntervalId) clearInterval(scanIntervalId);
  startBtn.disabled = false;
  stopBtn.disabled = true;
  scanLogs.textContent += `[${new Date().toLocaleTimeString()}] Scan stopped.\n`;
  appendLog('Scan stopped.');
});

const scanIntervalInput = document.getElementById('scanInterval');
const autoExecuteCheckbox = document.getElementById('autoExecuteOption');
const startAutoScanBtn = document.getElementById('startAutoScan');
const stopAutoScanBtn = document.getElementById('stopAutoScan');

startAutoScanBtn.addEventListener('click', async () => {
  clearErrors();
  const selected = Array.from(pairsSelect.selectedOptions).map(o => o.value);
  const body = {
    pairs: selected,
    min_volume: parseFloat(minVolEl.value) || 0,
    entry_threshold: parseFloat(entryEl.value) || 0,
    exit_threshold: parseFloat(exitEl.value) || 0,
    interval_seconds: parseInt(scanIntervalInput.value) || 10,
    auto_execute: autoExecuteCheckbox.checked
  };
  try {
    const res = await fetch('/autoscan/start', {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify(body)
    });
    if (!res.ok) {
      const err = await res.json();
      appendLog('Auto-scan start failed: ' + err.error);
      return;
    }
    appendLog('Auto-scan started');
    startAutoScanBtn.disabled = true;
    stopAutoScanBtn.disabled = false;
  } catch (err) {
    appendLog('Auto-scan error: ' + err.message);
  }
});

stopAutoScanBtn.addEventListener('click', async () => {
  try {
    const res = await fetch('/autoscan/stop', {method: 'POST'});
    if (!res.ok) {
      const err = await res.json();
      appendLog('Auto-scan stop failed: ' + err.error);
      return;
    }
    appendLog('Auto-scan stopped');
    startAutoScanBtn.disabled = false;
    stopAutoScanBtn.disabled = true;
  } catch (err) {
    appendLog('Auto-scan error: ' + err.message);
  }
});

// Backtest form submission
backtestForm.addEventListener('submit', async e => {
  e.preventDefault();
  backtestPairErrorEl.textContent = '';
  backtestSinceErrorEl.textContent = '';
  backtestShortErrorEl.textContent = '';
  backtestLongErrorEl.textContent = '';
  backtestFeeRateErrorEl.textContent = '';
  const pair = backtestPairSelect.value;
  const since = parseInt(backtestSinceInput.value, 10);
  const shortWin = parseInt(backtestShortInput.value, 10);
  const longWin = parseInt(backtestLongInput.value, 10);
  const feeRate = parseFloat(backtestFeeRateInput.value);
  if (!pair) { backtestPairErrorEl.textContent = 'Select a pair'; return; }
  if (isNaN(since) || since < 1) { backtestSinceErrorEl.textContent = 'Invalid minutes'; return; }
  if (isNaN(shortWin) || shortWin < 1) { backtestShortErrorEl.textContent = 'Invalid short window'; return; }
  if (isNaN(longWin) || longWin < 1 || longWin <= shortWin) { backtestLongErrorEl.textContent = 'Invalid long window'; return; }
  if (isNaN(feeRate) || feeRate < 0) { backtestFeeRateErrorEl.textContent = 'Invalid fee rate'; return; }
  try {
    const res = await fetch('/backtest', {
      method: 'POST', headers: {'Content-Type': 'application/json'}, body: JSON.stringify({ pair, since_minutes: since, short: shortWin, long: longWin, fee_rate: feeRate })
    });
    const data = await res.json();
    if (!res.ok) { appendLog('Backtest failed: ' + (data.error || res.statusText)); showToast('Backtest failed', 'error'); return; }
    backtestMetrics.innerHTML = '';
    const metricsArr = [['Trades', data.trades], ['Wins', data.wins], ['Losses', data.losses], ['Win Rate %', data.win_rate.toFixed(2)], ['Total PnL', data.total_pnl.toFixed(2)], ['Avg PnL', data.avg_pnl.toFixed(2)], ['Sharpe', data.sharpe.toFixed(2)], ['Max Drawdown', data.max_drawdown.toFixed(2)]];
    metricsArr.forEach(([label, value]) => {
      const card = document.createElement('div');
      card.className = 'bg-gray-50 border border-gray-200 text-gray-800 p-4 rounded-lg text-center';
      card.innerHTML = `<div class="text-sm">${label}</div><div class="font-bold">${value}</div>`;
      backtestMetrics.appendChild(card);
    });
    const ctx = document.getElementById('backtestPnlChart').getContext('2d');
    if (window.backtestPnlChart) window.backtestPnlChart.destroy();
    window.backtestPnlChart = new Chart(ctx, { type: 'line', data: { labels: data.pnl_history.map(h => new Date(h.time)), datasets: [{ label: 'Cumulative PnL', data: data.pnl_history.map(h => h.pnl), borderColor: '#3e95cd', fill: false }] }, options: { responsive: true, maintainAspectRatio: false, scales: { x: { type: 'time', time: { tooltipFormat: 'HH:mm:ss', unit: 'minute' } } }, plugins: { tooltip: { mode: 'index', intersect: false } } } });
    // Drawdown chart
    const ddCtx = document.getElementById('backtestDrawdownChart');
    if (window.backtestDrawChart) window.backtestDrawChart.destroy();
    window.backtestDrawChart = new Chart(ddCtx.getContext('2d'), { type: 'line', data: { labels: data.drawdown_history.map(h => new Date(h.time)), datasets: [{ label: 'Drawdown', data: data.drawdown_history.map(h => h.drawdown), borderColor: '#db3236', fill: false }] }, options: { responsive: true, maintainAspectRatio: false, scales: { x: { type: 'time', time: { tooltipFormat: 'HH:mm:ss', unit: 'minute' } } }, plugins: { tooltip: { mode: 'index', intersect: false } } } });
    showToast('Backtest complete', 'success');
  } catch (err) { appendLog('Backtest error: ' + err.message); showToast('Backtest error', 'error'); }
});

// Toast notifications helper
function showToast(msg, type = 'info') {
  const container = document.getElementById('toast-container');
  const toast = document.createElement('div');
  toast.className = `${type === 'success' ? 'bg-green-500' : type === 'error' ? 'bg-red-500' : 'bg-gray-500'} text-white px-4 py-2 rounded shadow`;
  toast.textContent = msg;
  container.appendChild(toast);
  setTimeout(() => toast.remove(), 3000);
}

// Track last simulation position for chart markers
let lastPosition = null;

// Live status update
async function fetchStatus() {
  try {
    const res = await fetch('/status');
    const { status } = await res.json();
    const badge = document.getElementById('statusBadge');
    badge.textContent = status;
    badge.classList.toggle('running', status === 'running');
    badge.classList.toggle('stopped', status !== 'running');
  } catch (e) {
    console.error('Status fetch error', e);
  }
}

// Initialize charts, sliders, and status on DOM load
document.addEventListener('DOMContentLoaded', () => {
  // Theme toggle initial state and persistence
  const htmlEl = document.documentElement;
  const themeToggle = document.getElementById('themeToggle');
  const savedTheme = localStorage.getItem('theme');
  if (savedTheme === 'dark') { htmlEl.classList.add('dark'); themeToggle.textContent = ''; } else { htmlEl.classList.remove('dark'); themeToggle.textContent = ''; }
  themeToggle.addEventListener('click', () => {
    if (htmlEl.classList.toggle('dark')) { localStorage.setItem('theme','dark'); themeToggle.textContent = ''; }
    else { localStorage.setItem('theme','light'); themeToggle.textContent = ''; }
    updateChartTheme();
  });

  fetchStatus();
  setInterval(fetchStatus, 5000);

  // Price chart init
  const priceCtx = document.getElementById('priceChart').getContext('2d');
  window.priceData = { labels: [], datasets: [
    { label: 'Mid Price', data: [], borderColor: '#3b82f6', backgroundColor: '#93c5fd', fill: false },
    { label: 'Markers', data: [], type: 'scatter', backgroundColor: [], pointRadius: 6 }
  ]};
  window.priceChart = new Chart(priceCtx, { type: 'line', data: priceData, options: { scales: { x: { type: 'time', time: { unit: 'minute' } } } } });

  // PnL chart init
  const pnlCtx = document.getElementById('pnlChart').getContext('2d');
  window.pnlData = { labels: [], datasets: [
    { label: 'Total PnL', data: [], borderColor: '#10b981', backgroundColor: '#6ee7b7', fill: false },
    { label: 'Entries', data: [], type: 'scatter', backgroundColor: 'green', pointRadius: 6 },
    { label: 'Exits', data: [], type: 'scatter', backgroundColor: 'red', pointRadius: 6 }
  ]};
  window.pnlChart = new Chart(pnlCtx, { type: 'line', data: pnlData, options: { scales: { x: { type: 'time', time: { unit: 'minute' } } } } });

  // Slider bindings
  const entrySlider = document.getElementById('entryThresholdSlider');
  const entryValEl = document.getElementById('entryThresholdValue');
  entrySlider.addEventListener('input', e => { entryValEl.textContent = e.target.value; document.getElementById('entryThreshold').value = e.target.value; });
  const exitSlider = document.getElementById('exitThresholdSlider');
  const exitValEl = document.getElementById('exitThresholdValue');
  exitSlider.addEventListener('input', e => { exitValEl.textContent = e.target.value; document.getElementById('exitThreshold').value = e.target.value; });

  // Load initial config and sync sliders
  loadConfig().then(() => {
    entrySlider.value = document.getElementById('entryThreshold').value;
    entryValEl.textContent = entrySlider.value;
    exitSlider.value = document.getElementById('exitThreshold').value;
    exitValEl.textContent = exitSlider.value;
  });

  // Chart.js theme adapter for dark/light mode
  function updateChartTheme() {
    const isDark = document.documentElement.classList.contains('dark');
    Chart.defaults.color = isDark ? '#ECEFF4' : '#1F2937';
    Chart.defaults.plugins.legend.labels.color = isDark ? '#ECEFF4' : '#1F2937';
    Chart.defaults.plugins.tooltip.titleColor = isDark ? '#ECEFF4' : '#1F2937';
    Chart.defaults.plugins.tooltip.bodyColor = isDark ? '#ECEFF4' : '#1F2937';
    [window.priceChart, window.pnlChart].forEach(chart => {
      chart.options.scales.x.ticks.color = isDark ? '#ECEFF4' : '#1F2937';
      chart.options.scales.y.ticks.color = isDark ? '#ECEFF4' : '#1F2937';
      if (chart.options.scales.x.grid) chart.options.scales.x.grid.color = isDark ? '#4C566A' : '#E5E7EB';
      if (chart.options.scales.y.grid) chart.options.scales.y.grid.color = isDark ? '#4C566A' : '#E5E7EB';
      chart.update();
    });
  }

  updateChartTheme();
});

// Chart.js setup for price and PnL
const priceCtx = document.getElementById('priceChart').getContext('2d');
const priceData = { labels: [], datasets: [{ label: 'Mid Price', data: [], borderColor: '#3e95cd', fill: false }] };
const priceChart = new Chart(priceCtx, {
  type: 'line',
  data: priceData,
  options: {
    responsive: true,
    maintainAspectRatio: false,
    scales: {
      x: {
        type: 'time',
        time: { unit: 'minute', tooltipFormat: 'HH:mm:ss' }
      }
    },
    plugins: {
      tooltip: { mode: 'index', intersect: false },
      zoom: {
        pan: { enabled: true, mode: 'x' },
        zoom: { wheel: { enabled: true }, pinch: { enabled: true }, mode: 'x' }
      }
    }
  }
});

const pnlCtx = document.getElementById('pnlChart').getContext('2d');
const pnlData = { labels: [], datasets: [{ label: 'Total PnL', data: [], borderColor: '#8e5ea2', fill: false }] };
const pnlChart = new Chart(pnlCtx, {
  type: 'line',
  data: pnlData,
  options: {
    responsive: true,
    maintainAspectRatio: false,
    scales: {
      x: {
        type: 'time',
        time: { unit: 'minute', tooltipFormat: 'HH:mm:ss' }
      }
    },
    plugins: {
      tooltip: { mode: 'index', intersect: false },
      zoom: {
        pan: { enabled: true, mode: 'x' },
        zoom: { wheel: { enabled: true }, pinch: { enabled: true }, mode: 'x' }
      }
    }
  }
});

// Depth chart setup
const depthCtx = document.getElementById('depthChart').getContext('2d');
const depthChart = new Chart(depthCtx, {
  type: 'scatter',
  data: {
    datasets: [
      { label: 'Bids', data: [], backgroundColor: '#3cba9f', showLine: true },
      { label: 'Asks', data: [], backgroundColor: '#db3236', showLine: true }
    ]
  },
  options: {
    scales: {
      x: { type: 'linear', title: { display: true, text: 'Price' } },
      y: { title: { display: true, text: 'Volume' } }
    },
    plugins: { legend: { position: 'top' }, tooltip: { mode: 'nearest', intersect: false } }
  }
});

// Fetch and display health status
async function fetchStatus() {
  const res = await fetch('/status');
  const data = await res.json();
  const badge = document.getElementById('statusBadge');
  badge.textContent = data.status;
  badge.classList.toggle('running', data.status === 'running');
  badge.classList.toggle('stopped', data.status !== 'running');
}
fetchStatus();

// Insert CSS file dynamically
const link = document.createElement('link');
link.rel = 'stylesheet';
link.href = '/assets/style.css';
document.head.appendChild(link);

// Initial load
loadConfig();
fetchOrderBook();
loadPairs();
