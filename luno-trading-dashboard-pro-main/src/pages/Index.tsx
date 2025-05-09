import { useState, useEffect, useRef } from "react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Switch } from "@/components/ui/switch";
import { Slider } from "@/components/ui/slider";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import { Badge } from "@/components/ui/badge";
import { ScrollArea } from "@/components/ui/scroll-area";
import { ChevronLeft, ChevronRight, BarChart, RefreshCw, Info, Loader2, BrainCircuit, Activity, AlertTriangle } from "lucide-react";
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";
import * as RechartsPrimitive from "recharts";
import { ChartContainer } from "@/components/ui/chart";
import { Toaster, toast } from 'sonner';
import BalanceCard from '@/components/BalanceCard';
import SimulationTab from '@/components/SimulationTab';
import ConfigTab from '@/components/ConfigTab';
import OrderBookTab from '@/components/OrderBookTab';
import ScannerTab from '@/components/ScannerTab';
import BacktestTab from '@/components/BacktestTab';
import ThresholdsTab from '@/components/ThresholdsTab';
import AIInsightsTab from '@/components/AIInsightsTab';
import FeedbackSystem, { FeedbackEvent } from '@/components/FeedbackSystem';

const Index = () => {
  const [darkMode, setDarkMode] = useState(true);
  const toggleTheme = () => setDarkMode(prev => !prev);

  const [pairs, setPairs] = useState<string[]>([])
  const [scanResults, setScanResults] = useState<any[]>([])
  const [topPicks, setTopPicks] = useState<any[]>([])
  const [topLimit, setTopLimit] = useState<number>(5)
  const [autoRefreshPicks, setAutoRefreshPicks] = useState<boolean>(false)
  const [picksInterval, setPicksInterval] = useState<number>(10)
  const [autoScanning, setAutoScanning] = useState(false)
  const [simResult, setSimResult] = useState<any>(null)
  const [execResult, setExecResult] = useState<any>(null)
  const [orderbookData, setOrderbookData] = useState<any>(null)
  const [orderbookAuto, setOrderbookAuto] = useState(false)
  const [orderbookLoading, setOrderbookLoading] = useState<boolean>(false)
  const [refreshIntervalSeconds, setRefreshIntervalSeconds] = useState<number>(5)
  const [apiConnected, setApiConnected] = useState<boolean>(false)
  const [streamPicks, setStreamPicks] = useState<boolean>(false)
  const [orderbookPair, setOrderbookPair] = useState<string>("")
  const [simLoading, setSimLoading] = useState<boolean>(false)
  const [depthData, setDepthData] = useState<{price:number; bids:number; asks:number}[]>([])

  const [configObj, setConfigObj] = useState<any>(null)
  const [cfgPair, setCfgPair] = useState<string>("")
  const [cfgEntryThreshold, setCfgEntryThreshold] = useState<number>(0)
  const [cfgExitThreshold, setCfgExitThreshold] = useState<number>(0)
  const [cfgStakeSize, setCfgStakeSize] = useState<number>(0)
  const [cfgCooldown, setCfgCooldown] = useState<string>("")
  const [selectedPairs, setSelectedPairs] = useState<string[]>([])
  const [scanMinVolume, setScanMinVolume] = useState<string>("")
  const [scanEntryThreshold, setScanEntryThreshold] = useState<string>("")
  const [scanExitThreshold, setScanExitThreshold] = useState<string>("")
  const [scanIntervalSec, setScanIntervalSec] = useState<string>("10")
  const [scanAutoExecute, setScanAutoExecute] = useState<boolean>(false)
  const [btPair, setBtPair] = useState<string>("")
  const [btSinceMinutes, setBtSinceMinutes] = useState<number>(60)
  const [btShortWindow, setBtShortWindow] = useState<number>(10)
  const [btLongWindow, setBtLongWindow] = useState<number>(50)
  const [btFeeRate, setBtFeeRate] = useState<number>(0.001)
  const [btMetrics, setBtMetrics] = useState<any>(null)
  const [btPnlHistory, setBtPnlHistory] = useState<any[]>([])
  const [btDrawdownHistory, setBtDrawdownHistory] = useState<any[]>([])
  const [logs, setLogs] = useState<string[]>([])
  const [logFilter, setLogFilter] = useState<string>("")
  const [logFrom, setLogFrom] = useState<string>("")
  const [logTo, setLogTo] = useState<string>("")
  const [balances, setBalances] = useState<any[]>([])
  const [percentChange, setPercentChange] = useState<number | null>(null)
  const [scannerSaved, setScannerSaved] = useState(false)

  // Preset management
  const [savedConfigs, setSavedConfigs] = useState<{name:string;description:string;cfg:any}[]>([])
  const [selectedPreset, setSelectedPreset] = useState<string>("")

  // Scanner presets
  const [scannerPresets, setScannerPresets] = useState<any[]>([])
  const [selectedScannerPreset, setSelectedScannerPreset] = useState<string>("")
  
  // AI insights state
  const [aiInsights, setAiInsights] = useState<any[]>([])
  const [aiParameters, setAiParameters] = useState<any>(null)
  const [aiEnabled, setAiEnabled] = useState<boolean>(false)
  const [aiAutoExecute, setAiAutoExecute] = useState<boolean>(false)
  const [aiRefreshInterval, setAiRefreshInterval] = useState<number>(5)

  // Backtest presets
  const [backtestPresets, setBacktestPresets] = useState<any[]>([])
  const [selectedBacktestPreset, setSelectedBacktestPreset] = useState<string>("")

  // Threshold presets
  const [thresholdPresets, setThresholdPresets] = useState<any[]>([])
  const [selectedThresholdPreset, setSelectedThresholdPreset] = useState<string>("")

  // Loading and operation states
  const [scanLoading, setScanLoading] = useState(false)
  const [autoScanLoading, setAutoScanLoading] = useState(false)
  const [btLoading, setBtLoading] = useState(false)
  const [thLoading, setThLoading] = useState(false)
  const [aiLoading, setAiLoading] = useState(false)

  // Threshold optimization state
  const [thSelectedPairs, setThSelectedPairs] = useState<string[]>([]);
  const [thSinceMinutes, setThSinceMinutes] = useState<number>(60);
  const [thFeeRate, setThFeeRate] = useState<number>(0.001);
  const [gridStart, setGridStart] = useState<number>(0.01);
  const [gridEnd, setGridEnd] = useState<number>(0.05);
  const [gridStep, setGridStep] = useState<number>(0.005);
  const [thResults, setThResults] = useState<any[]>([]);

  // Error states
  const [configError, setConfigError] = useState<string>("")
  const [scanError, setScanError] = useState<string>("")
  const [btError, setBtError] = useState<string>("")
  const [thError, setThError] = useState<string>("")
  const [simError, setSimError] = useState<string>("")
  const [orderbookError, setOrderbookError] = useState<string>("")
  const [aiError, setAiError] = useState<string>("")
  
  // Feedback system state
  const [feedbackEvents, setFeedbackEvents] = useState<FeedbackEvent[]>([])
  const [showFeedbackPanel, setShowFeedbackPanel] = useState<boolean>(true)

  // CSV export helper
  const exportToCsv = (filename: string, rows: any[]) => {
    if (!rows || !rows.length) return;
    const headers = Object.keys(rows[0]);
    const csvContent = [headers.join(','), ...rows.map(row => headers.map(field => JSON.stringify(row[field])).join(','))].join('\r\n');
    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.setAttribute('download', filename);
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  };
  
  // Feedback system event handlers
  const addFeedbackEvent = (event: Omit<FeedbackEvent, 'id' | 'timestamp'>) => {
    const newEvent = {
      ...event,
      id: `event-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
      timestamp: new Date(),
      acknowledged: false
    };
    
    setFeedbackEvents(prev => [newEvent, ...prev]);
    
    // For error and warning events, show a toast notification
    if (event.category === 'error') {
      toast.error(event.title, {
        description: event.message,
      });
    } else if (event.category === 'warning') {
      toast.warning(event.title, {
        description: event.message,
      });
    } else if (event.category === 'success') {
      toast.success(event.title, {
        description: event.message,
      });
    }

    return newEvent.id;
  };
  
  const clearFeedbackEvents = () => {
    setFeedbackEvents([]);
    toast.info('Feedback events cleared');
  };
  
  const acknowledgeFeedbackEvent = (id: string) => {
    setFeedbackEvents(prev => 
      prev.map(event => 
        event.id === id ? { ...event, acknowledged: true } : event
      )
    );
  };
  
  const exportFeedbackEvents = () => {
    const eventsForExport = feedbackEvents.map(event => ({
      timestamp: event.timestamp.toISOString(),
      category: event.category,
      title: event.title,
      message: event.message,
      source: event.source || '',
      duration: event.duration || '',
    }));
    
    exportToCsv('trading-bot-events.csv', eventsForExport);
    toast.success('Events exported to CSV');
  };

  const handleDownloadLogs = () => {
    const blob = new Blob([logs.join('\n')], { type: 'text/plain;charset=utf-8;' });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url; link.setAttribute('download', 'logs.txt'); document.body.appendChild(link);
    link.click(); document.body.removeChild(link);
  };

  const filteredLogs = logs.filter(line => {
    if (logFilter && !line.includes(logFilter)) return false;
    const ts = Date.parse(line.split(' ')[0]);
    if (logFrom && !isNaN(ts) && ts < Date.parse(logFrom)) return false;
    if (logTo && !isNaN(ts) && ts > Date.parse(logTo)) return false;
    return true;
  });

  useEffect(() => {
    // Set body class based on theme
    document.documentElement.classList.toggle('dark', darkMode);
    
    // This would normally be where you'd set up websocket connections
    // or API polling to get real-time data from your trading bot
  }, [darkMode]);

  useEffect(() => {
    fetch('/pairs').then(r=>r.json()).then(p=>{
      setPairs(p)
      setSelectedPairs(p)
      if(p.length) {
        setBtPair(p[0])
        setOrderbookPair(p[0])
      }
    })
    fetch('/config').then(r=>r.json()).then(cfg=>{ setConfigObj(cfg); setCfgPair(cfg.pair); setCfgEntryThreshold(cfg.entry_threshold); setCfgExitThreshold(cfg.exit_threshold); setCfgStakeSize(cfg.stake_size); setCfgCooldown(cfg.cooldown) })
    const ivLogs = setInterval(async ()=>{ const data = await fetch('/logs').then(r=>r.json()); setLogs(data) },2000)
    return ()=>clearInterval(ivLogs)
  }, [])

  useEffect(() => {
    fetch('/healthz')
      .then(res => setApiConnected(res.ok))
      .catch(() => setApiConnected(false))
  }, [])

  useEffect(() => {
    if (!configObj) return;
    fetch('/balances')
      .then(res => res.json())
      .then(data => setBalances(data));
  }, [configObj])

  useEffect(() => {
    fetch('/percent-change')
      .then(res => res.json())
      .then(data => setPercentChange(data.percent_change))
      .catch(() => setPercentChange(null))
  }, [])

  useEffect(() => {
    // Load saved scanner settings
    const saved = localStorage.getItem('scannerSettings');
    if (saved) {
      try {
        const { selectedPairs: sp, scanMinVolume: mv, scanEntryThreshold: et, scanExitThreshold: xt } = JSON.parse(saved);
        setSelectedPairs(sp);
        setScanMinVolume(mv.toString());
        setScanEntryThreshold(et.toString());
        setScanExitThreshold(xt.toString());
      } catch {}
    }
  }, []);

  useEffect(() => {
    const raw = localStorage.getItem('savedConfigs')
    if (raw) {
      try { setSavedConfigs(JSON.parse(raw)) } catch {}
    }
  }, [])

  useEffect(() => {
    const sp = localStorage.getItem('scannerPresets'); if (sp) setScannerPresets(JSON.parse(sp))
    const bp = localStorage.getItem('backtestPresets'); if (bp) setBacktestPresets(JSON.parse(bp))
    const tp = localStorage.getItem('thresholdPresets'); if (tp) setThresholdPresets(JSON.parse(tp))
  }, [])

  useEffect(() => {
    let iv: ReturnType<typeof setInterval>
    if (autoRefreshPicks) {
      handleFetchTopPicks()
      iv = setInterval(handleFetchTopPicks, picksInterval * 1000)
    }
    return () => clearInterval(iv)
  }, [autoRefreshPicks, picksInterval])

  useEffect(() => {
    let es: EventSource | undefined;
    if (streamPicks) {
      const params = new URLSearchParams();
      params.append('pairs', selectedPairs.join(','));
      params.append('min_volume', (parseFloat(scanMinVolume) || 0).toString());
      params.append('limit', topLimit.toString());
      params.append('interval', picksInterval.toString());
      es = new EventSource(`/stream/opportunities?${params.toString()}`);
      es.addEventListener('opportunity', e => {
        const data = JSON.parse(e.data);
        setTopPicks(data);
        if (data.length > 0) {
          const top = data[0];
          toast.success(`${top.pair}: ${top.potential.toFixed(2)}% (score ${top.score.toFixed(2)})`);
        }
      });
    }
    return () => { if (es) es.close(); };
  }, [streamPicks, selectedPairs, scanMinVolume, topLimit, picksInterval]);

  const autoScanTimerRef = useRef<number | null>(null)

  const handleSaveConfig = async e => {
    if (!cfgPair) { setConfigError('Pair is required'); return }
    e.preventDefault()
    setConfigError("")
    try {
      const newCfg = { pair: cfgPair, entry_threshold: cfgEntryThreshold, exit_threshold: cfgExitThreshold, stake_size: cfgStakeSize, cooldown: cfgCooldown }
      const res = await fetch('/config',{method:'PUT', headers:{'Content-Type':'application/json'}, body:JSON.stringify(newCfg)})
      if (!res.ok) throw new Error(`Status ${res.status}`)
      const updated = await res.json()
      setConfigObj(updated)
      toast.success('Config saved')
    } catch (err: any) {
      setConfigError(err.message || 'Failed to save config')
      toast.error('Config save failed')
    }
  }

  const handleFetchOrderbook = async () => {
    if (!orderbookPair) { setOrderbookError('Select a pair'); return }
    setOrderbookError("")
    setOrderbookLoading(true)
    try {
      const res = await fetch(`/orderbook?pair=${orderbookPair}`)
      if (!res.ok) throw new Error(`Status ${res.status}`)
      const data = await res.json()
      setOrderbookData(data)
      // build depth chart data
      const bidsArr = data.bids.map(o=>({price: Number(o.price), volume: Number(o.volume)})).sort((a,b)=>a.price - b.price)
      const asksArr = data.asks.map(o=>({price: Number(o.price), volume: Number(o.volume)})).sort((a,b)=>a.price - b.price)
      let cum = 0
      const bidDepth = bidsArr.map(item => { cum += item.volume; return {price: item.price, bids: cum, asks: 0} })
      cum = 0
      const askDepth = asksArr.map(item => { cum += item.volume; return {price: item.price, bids: 0, asks: cum} })
      const merged = [...bidDepth, ...askDepth].sort((a,b)=>a.price - b.price)
      setDepthData(merged)
    } catch (err: any) {
      setOrderbookError(err.message || 'Failed to fetch orderbook')
      toast.error('Failed to fetch orderbook')
    } finally {
      setOrderbookLoading(false)
    }
  }

  const handleSimulate = async (e) => {
    e?.preventDefault?.()
    setSimError("")
    setSimLoading(true)
    try {
      const res = await fetch('/simulate',{method:'POST'})
      if (!res.ok) throw new Error(`Status ${res.status}`)
      const data = await res.json()
      setSimResult(data)
      toast.success('Simulation complete')
    } catch (err: any) {
      setSimError(err.message || 'Simulation failed')
      toast.error('Simulation failed')
    } finally {
      setSimLoading(false)
    }
  }

  const handleScan = async () => {
    setScanError("")
    setScanLoading(true)
    
    // Add feedback event for scan start
    const eventId = addFeedbackEvent({
      category: 'trading',
      title: 'Market Scan Started',
      message: `Scanning ${selectedPairs.length} pairs with min volume: ${scanMinVolume}, entry: ${scanEntryThreshold}, exit: ${scanExitThreshold}`,
      source: 'Scanner',
      progress: 0
    });
    
    try {
      // Update progress
      setFeedbackEvents(prev => prev.map(e => 
        e.id === eventId ? { ...e, progress: 30 } : e
      ));
      
      const body = { 
        pairs: selectedPairs, 
        min_volume: parseFloat(scanMinVolume)||0, 
        entry_threshold: parseFloat(scanEntryThreshold)||0, 
        exit_threshold: parseFloat(scanExitThreshold)||0 
      }
      
      // Update progress
      setFeedbackEvents(prev => prev.map(e => 
        e.id === eventId ? { ...e, progress: 60 } : e
      ));
      
      const startTime = Date.now();
      const res = await fetch('/scan',{
        method:'POST',
        headers:{'Content-Type':'application/json'},
        body:JSON.stringify(body)
      })
      
      if (!res.ok) throw new Error(`Status ${res.status}`)
      
      const data = await res.json()
      setScanResults(data)
      const duration = Date.now() - startTime;
      
      // Update success event
      setFeedbackEvents(prev => prev.map(e => 
        e.id === eventId ? { 
          ...e, 
          category: 'success',
          title: 'Market Scan Completed',
          message: `Found ${data.length} opportunities. Scan took ${duration}ms.`,
          progress: 100,
          duration: duration
        } : e
      ));
      
      toast.success('Scan complete')
    } catch (err: any) {
      setScanError(err.message || 'Scan failed')
      
      // Update error event
      setFeedbackEvents(prev => prev.map(e => 
        e.id === eventId ? { 
          ...e, 
          category: 'error',
          title: 'Market Scan Failed',
          message: err.message || 'Scan failed',
          progress: 100
        } : e
      ));
      
      toast.error('Scan failed')
    } finally {
      setScanLoading(false)
    }
  }

  const handleStartAutoScan = async () => {
    setAutoScanLoading(true)
    try {
      const body = { pairs: selectedPairs, min_volume: parseFloat(scanMinVolume)||0, entry_threshold: parseFloat(scanEntryThreshold)||0, exit_threshold: parseFloat(scanExitThreshold)||0, interval_seconds: parseInt(scanIntervalSec)||0, auto_execute: scanAutoExecute }
      await fetch('/autoscan/start',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)})
      setAutoScanning(true)
      toast.success('Auto-scan started')
      autoScanTimerRef.current = window.setInterval(handleScan, (parseInt(scanIntervalSec)||0)*1000)
    } catch {
      toast.error('Failed to start auto-scan')
    } finally {
      setAutoScanLoading(false)
    }
  }

  const handleStopAutoScan = async () => {
    setAutoScanLoading(true)
    try {
      await fetch('/autoscan/stop',{method:'POST'})
      if (autoScanTimerRef.current) {
        clearInterval(autoScanTimerRef.current)
        autoScanTimerRef.current = null
      }
      setAutoScanning(false)
      toast.success('Auto-scan stopped')
    } catch {
      toast.error('Failed to stop auto-scan')
    } finally {
      setAutoScanLoading(false)
    }
  }

  const handleRunBacktest = async e => {
    if (!btPair) { setBtError('Pair is required'); return }
    e.preventDefault()
    setBtError("")
    setBtLoading(true)
    try {
      const body = { pair: btPair, since_minutes: btSinceMinutes, short: btShortWindow, long: btLongWindow, fee_rate: btFeeRate }
      const res = await fetch('/backtest',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)})
      if (!res.ok) throw new Error(`Status ${res.status}`)
      const data = await res.json()
      setBtMetrics(data); setBtPnlHistory(data.pnl_history); setBtDrawdownHistory(data.drawdown_history)
      toast.success('Backtest complete')
    } catch (err: any) {
      setBtError(err.message || 'Backtest failed')
      toast.error('Backtest failed')
    } finally {
      setBtLoading(false)
    }
  }

  const handleRunThresholds = async e => {
    if (thSelectedPairs.length === 0) { setThError('Select at least one pair'); return }
    e.preventDefault()
    setThError("")
    setThLoading(true)
    try {
      const body = { pairs: thSelectedPairs, since_minutes: thSinceMinutes, fee_rate: thFeeRate, grid_start: gridStart, grid_end: gridEnd, grid_step: gridStep }
      const res = await fetch('/thresholds',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)})
      if (!res.ok) throw new Error(`Status ${res.status}`)
      const data = await res.json()
      setThResults(data)
      toast.success('Threshold optimization complete')
    } catch (err: any) {
      setThError(err.message || 'Optimization failed')
      toast.error('Optimization failed')
    } finally {
      setThLoading(false)
    }
  }

  const handleFetchAIInsights = async () => {
    setAiError("")
    setAiLoading(true)
    try {
      const res = await fetch('/ai/insights')
      if (!res.ok) throw new Error(`Status ${res.status}`)
      const data = await res.json()
      setAiInsights(data)
      toast.success('AI insights updated')
    } catch (err: any) {
      setAiError(err.message || 'Failed to fetch AI insights')
      toast.error('Failed to fetch AI insights')
    } finally {
      setAiLoading(false)
    }
  }
  
  const handleToggleAI = async () => {
    try {
      const enabled = !aiEnabled
      await fetch(`/ai/${enabled ? 'enable' : 'disable'}`, {method: 'POST'})
      setAiEnabled(enabled)
      toast.success(`AI engine ${enabled ? 'enabled' : 'disabled'}`)
      if (enabled) {
        handleFetchAIInsights()
      }
    } catch (err: any) {
      toast.error(`Failed to ${aiEnabled ? 'disable' : 'enable'} AI engine`)
    }
  }
  
  const handleUpdateAISettings = async () => {
    try {
      const body = {
        auto_execute: aiAutoExecute,
        refresh_interval: aiRefreshInterval
      }
      await fetch('/ai/settings', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify(body)
      })
      toast.success('AI settings updated')
    } catch (err: any) {
      toast.error('Failed to update AI settings')
    }
  }

  const handleFetchTopPicks = async () => {
    const body = { pairs: selectedPairs, min_volume: parseFloat(scanMinVolume) || 0, limit: topLimit }
    const data = await fetch('/opportunities', {
      method: 'POST', headers: {'Content-Type':'application/json'}, body: JSON.stringify(body)
    }).then(r=>r.json())
    setTopPicks(data)
  }

  const handleSaveScannerSettings = () => {
    const cfg = {
      selectedPairs,
      scanMinVolume: parseFloat(scanMinVolume) || 0,
      scanEntryThreshold: parseFloat(scanEntryThreshold) || 0,
      scanExitThreshold: parseFloat(scanExitThreshold) || 0
    }
    localStorage.setItem('scannerSettings', JSON.stringify(cfg))
    setScannerSaved(true)
    setTimeout(() => setScannerSaved(false), 2000)
  }

  const handleSavePreset = () => {
    const name = prompt("Enter preset name:")
    if (!name) return
    const description = prompt("Enter description:") || ""
    const preset = { name, description, cfg: { pair: cfgPair, entryThreshold: cfgEntryThreshold, exitThreshold: cfgExitThreshold, stakeSize: cfgStakeSize, cooldown: cfgCooldown } }
    const updated = [...savedConfigs.filter(p=>p.name!==name), preset]
    setSavedConfigs(updated)
    localStorage.setItem('savedConfigs', JSON.stringify(updated))
    setSelectedPreset(name)
    toast.success('Preset saved!')
  }

  const handleSelectPreset = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const name = e.target.value
    setSelectedPreset(name)
    const p = savedConfigs.find(x=>x.name===name)
    if (p) {
      setCfgPair(p.cfg.pair)
      setCfgEntryThreshold(p.cfg.entryThreshold)
      setCfgExitThreshold(p.cfg.exitThreshold)
      setCfgStakeSize(p.cfg.stakeSize)
      setCfgCooldown(p.cfg.cooldown)
    }
  }

  const handleDeletePreset = () => {
    if (!selectedPreset) return
    if (!confirm(`Delete preset ${selectedPreset}?`)) return
    const updated = savedConfigs.filter(p=>p.name!==selectedPreset)
    setSavedConfigs(updated)
    localStorage.setItem('savedConfigs', JSON.stringify(updated))
    setSelectedPreset("")
    toast.success('Preset deleted')
  }

  const handleSaveScannerPreset = () => {
    const name = prompt('Scanner preset name:'); if (!name) return
    const preset = { name, cfg: { selectedPairs, scanMinVolume, scanEntryThreshold, scanExitThreshold, scanIntervalSec, scanAutoExecute } }
    const updated = [...scannerPresets.filter(p=>p.name!==name), preset]
    setScannerPresets(updated)
    localStorage.setItem('scannerPresets', JSON.stringify(updated))
    setSelectedScannerPreset(name)
    toast.success('Scanner preset saved')
  }

  const handleSelectScannerPreset = (e:any) => {
    const name = e.target.value; setSelectedScannerPreset(name)
    const p = scannerPresets.find(x=>x.name===name); if (p) {
      setSelectedPairs(p.cfg.selectedPairs)
      setScanMinVolume(p.cfg.scanMinVolume)
      setScanEntryThreshold(p.cfg.scanEntryThreshold)
      setScanExitThreshold(p.cfg.scanExitThreshold)
      setScanIntervalSec(p.cfg.scanIntervalSec)
      setScanAutoExecute(p.cfg.scanAutoExecute)
    }
  }

  const handleDeleteScannerPreset = () => {
    if (!selectedScannerPreset) return
    if (!confirm(`Delete scanner preset ${selectedScannerPreset}?`)) return
    const updated = scannerPresets.filter(p=>p.name!==selectedScannerPreset)
    setScannerPresets(updated)
    localStorage.setItem('scannerPresets', JSON.stringify(updated))
    setSelectedScannerPreset("")
    toast.success('Scanner preset deleted')
  }

  const handleSaveBacktestPreset = () => {
    const name = prompt('Backtest preset name:'); if (!name) return
    const preset = { name, cfg: { btPair, btSinceMinutes, btShortWindow, btLongWindow, btFeeRate } }
    const updated = [...backtestPresets.filter(p=>p.name!==name), preset]
    setBacktestPresets(updated)
    localStorage.setItem('backtestPresets', JSON.stringify(updated))
    setSelectedBacktestPreset(name)
    toast.success('Backtest preset saved')
  }

  const handleSelectBacktestPreset = (e:any) => {
    const name = e.target.value; setSelectedBacktestPreset(name)
    const p = backtestPresets.find(x=>x.name===name); if (p) {
      setBtPair(p.cfg.btPair)
      setBtSinceMinutes(p.cfg.btSinceMinutes)
      setBtShortWindow(p.cfg.btShortWindow)
      setBtLongWindow(p.cfg.btLongWindow)
      setBtFeeRate(p.cfg.btFeeRate)
    }
  }

  const handleDeleteBacktestPreset = () => {
    if (!selectedBacktestPreset) return
    if (!confirm(`Delete backtest preset ${selectedBacktestPreset}?`)) return
    const updated = backtestPresets.filter(p=>p.name!==selectedBacktestPreset)
    setBacktestPresets(updated)
    localStorage.setItem('backtestPresets', JSON.stringify(updated))
    setSelectedBacktestPreset("")
    toast.success('Backtest preset deleted')
  }

  const handleSaveThresholdPreset = () => {
    const name = prompt('Threshold preset name:'); if (!name) return
    const preset = { name, cfg: { thSelectedPairs, thSinceMinutes, thFeeRate, gridStart, gridEnd, gridStep } }
    const updated = [...thresholdPresets.filter(p=>p.name!==name), preset]
    setThresholdPresets(updated)
    localStorage.setItem('thresholdPresets', JSON.stringify(updated))
    setSelectedThresholdPreset(name)
    toast.success('Threshold preset saved')
  }

  const handleSelectThresholdPreset = (e:any) => {
    const name = e.target.value; setSelectedThresholdPreset(name)
    const p = thresholdPresets.find(x=>x.name===name); if (p) {
      setThSelectedPairs(p.cfg.thSelectedPairs)
      setThSinceMinutes(p.cfg.thSinceMinutes)
      setThFeeRate(p.cfg.thFeeRate)
      setGridStart(p.cfg.gridStart)
      setGridEnd(p.cfg.gridEnd)
      setGridStep(p.cfg.gridStep)
    }
  }

  const handleDeleteThresholdPreset = () => {
    if (!selectedThresholdPreset) return
    if (!confirm(`Delete threshold preset ${selectedThresholdPreset}?`)) return
    const updated = thresholdPresets.filter(p=>p.name!==selectedThresholdPreset)
    setThresholdPresets(updated)
    localStorage.setItem('thresholdPresets', JSON.stringify(updated))
    setSelectedThresholdPreset("")
    toast.success('Threshold preset deleted')
  }

  useEffect(() => {
    if (orderbookAuto) {
      const iv = setInterval(handleFetchOrderbook, refreshIntervalSeconds * 1000)
      return () => clearInterval(iv)
    }
  }, [orderbookAuto, refreshIntervalSeconds])

  useEffect(() => {
    return () => {
      if (autoScanTimerRef.current) {
        clearInterval(autoScanTimerRef.current)
      }
    }
  }, [])

  const [configFilter, setConfigFilter] = useState<string>("")
  const [orderbookFilter, setOrderbookFilter] = useState<string>("")
  const [scannerFilter, setScannerFilter] = useState<string>("")
  const [thresholdFilter, setThresholdFilter] = useState<string>("")

  return (
    <div className="min-h-screen bg-gradient-to-b from-gray-50 to-gray-100 dark:from-gray-900 dark:to-gray-800 p-4 md:p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        {/* Feedback System toggle button */}
        <div className="fixed bottom-4 right-4 z-50">
          <Button
            className="rounded-full size-10 shadow-lg flex items-center justify-center bg-indigo-600 hover:bg-indigo-700"
            onClick={() => setShowFeedbackPanel(prev => !prev)}
            title={`${showFeedbackPanel ? 'Hide' : 'Show'} feedback panel`}
          >
            <Activity className="h-5 w-5" />
          </Button>
          {feedbackEvents.filter(e => e.category === 'error' && !e.acknowledged).length > 0 && (
            <Badge className="absolute -top-2 -right-2 bg-red-500">
              {feedbackEvents.filter(e => e.category === 'error' && !e.acknowledged).length}
            </Badge>
          )}
        </div>
        
        {/* Feedback System panel */}
        {showFeedbackPanel && (
          <div className="fixed bottom-16 right-4 z-40 w-full max-w-2xl">
            <FeedbackSystem
              feedbackEvents={feedbackEvents}
              onClearEvents={clearFeedbackEvents}
              onAcknowledgeEvent={acknowledgeFeedbackEvent}
              onExportEvents={exportFeedbackEvents}
              maxEventsDisplayed={100}
            />
          </div>
        )}
        
        {/* Header section with KPI cards */}
        <div className="flex flex-col space-y-4">
          <div className="flex justify-between items-center">
            <h1 className="text-3xl font-bold text-gray-900 dark:text-gray-50">Luno Trading Bot Dashboard</h1>
            <div className="flex items-center space-x-2">
              <Badge className={apiConnected ? 'bg-green-100 dark:bg-green-800 text-green-800 dark:text-green-100' : 'bg-red-100 dark:bg-red-800 text-red-800 dark:text-red-100'}>
                {apiConnected ? 'API Connected' : 'API Down'}
              </Badge>
              <Button
                variant="outline"
                size="icon"
                onClick={toggleTheme}
                aria-label="Toggle theme"
                className="rounded-full"
              >
                {darkMode ? "☀️" : "🌙"}
              </Button>
            </div>
          </div>

          {/* Simulation KPI Cards */}
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-6 gap-4">
            <Card className="bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700">
              <CardHeader className="pb-2">
                <CardTitle className="text-sm font-medium text-gray-500 dark:text-gray-400">Status</CardTitle>
              </CardHeader>
              <CardContent>
                <Badge className={apiConnected ? 'bg-green-100 dark:bg-green-800 text-green-800 dark:text-green-100' : 'bg-red-100 dark:bg-red-800 text-red-800 dark:text-red-100'}>
                  {apiConnected ? 'Connected' : 'Down'}
                </Badge>
              </CardContent>
            </Card>

            <Card className="bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700">
              <CardHeader className="pb-2">
                <CardTitle className="text-sm font-medium text-gray-500 dark:text-gray-400">Signal</CardTitle>
              </CardHeader>
              <CardContent>
                <Badge className="font-bold">{simResult?.signal || "-"}</Badge>
              </CardContent>
            </Card>

            <Card className="bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700">
              <CardHeader className="pb-2">
                <CardTitle className="text-sm font-medium text-gray-500 dark:text-gray-400">Position</CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-2xl font-mono font-semibold">{simResult?.position ?? "-"}</p>
              </CardContent>
            </Card>

            <Card className="bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700">
              <CardHeader className="pb-2">
                <CardTitle className="text-sm font-medium text-gray-500 dark:text-gray-400">Total PnL</CardTitle>
              </CardHeader>
              <CardContent>
                <p className={`text-2xl font-mono font-semibold ${simResult && simResult.total_pnl >= 0 ? 'text-green-600 dark:text-green-400' : 'text-red-600 dark:text-red-400'}`}>
                  {simResult?.total_pnl?.toFixed(2) ?? "-"}
                </p>
              </CardContent>
            </Card>

            <Card className="bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700">
              <CardHeader className="pb-2">
                <CardTitle className="text-sm font-medium text-gray-500 dark:text-gray-400">Drawdown Exceeded</CardTitle>
              </CardHeader>
              <CardContent>
                <p className={`text-2xl font-mono font-semibold ${simResult?.max_drawdown_exceeded ? 'text-red-600 dark:text-red-400' : 'text-green-600 dark:text-green-400'}`}>
                  {simResult?.max_drawdown_exceeded ? 'Yes' : 'No'}
                </p>
              </CardContent>
            </Card>

            <BalanceCard balances={balances} configPair={configObj?.pair} percentChange={percentChange} />
          </div>
        </div>

        {/* Main content tabs */}
        <Tabs defaultValue="config" className="w-full">
          <TabsList className="grid grid-cols-2 md:grid-cols-7 gap-2 bg-gray-100 dark:bg-gray-800 p-1 rounded-lg mb-6">
            <TabsTrigger value="config">Config</TabsTrigger>
            <TabsTrigger value="orderbook">Order Book</TabsTrigger>
            <TabsTrigger value="simulation">Simulation</TabsTrigger>
            <TabsTrigger value="scanner">Scanner</TabsTrigger>
            <TabsTrigger value="backtest">Backtest</TabsTrigger>
            <TabsTrigger value="thresholds">Thresholds</TabsTrigger>
            <TabsTrigger value="ai"><div className="flex items-center gap-1"><BrainCircuit className="w-4 h-4" /> AI Insights</div></TabsTrigger>
            <TabsTrigger value="logs">Logs</TabsTrigger>
          </TabsList>

          {/* Configuration Tab */}
          <TabsContent value="config" className="space-y-4">
            <ConfigTab pairs={pairs} configObj={configObj} onConfigSave={setConfigObj} />
          </TabsContent>

          {/* Order Book Tab */}
          <TabsContent value="orderbook">
            <OrderBookTab pairs={pairs} />
          </TabsContent>

          {/* Simulation Tab */}
          <TabsContent value="simulation" className="space-y-4">
            <SimulationTab />
          </TabsContent>

          {/* Scanner Tab */}
          <TabsContent value="scanner">
            <ScannerTab pairs={pairs} />
          </TabsContent>

          {/* Thresholds Tab */}
          <TabsContent value="thresholds" className="space-y-4">
            <ThresholdsTab pairs={pairs} />
          </TabsContent>

          {/* AI Insights Tab */}
          <TabsContent value="ai" className="space-y-4">
            <AIInsightsTab
              aiInsights={aiInsights}
              aiEnabled={aiEnabled}
              aiAutoExecute={aiAutoExecute}
              aiRefreshInterval={aiRefreshInterval}
              aiLoading={aiLoading}
              aiError={aiError}
              onRefresh={handleFetchAIInsights}
              onToggleAI={handleToggleAI}
              onToggleAutoExecute={() => setAiAutoExecute(!aiAutoExecute)}
              onChangeInterval={(value) => setAiRefreshInterval(value)}
              onSaveSettings={handleUpdateAISettings}
            />
          </TabsContent>

          {/* Logs Tab */}
          <TabsContent value="logs">
            <div className="flex flex-wrap items-center space-x-2 mb-2">
              <Input placeholder="Search logs" value={logFilter} onChange={e => setLogFilter(e.target.value)} className="w-48" />
              <Input type="datetime-local" value={logFrom} onChange={e => setLogFrom(e.target.value)} aria-label="From" />
              <Input type="datetime-local" value={logTo} onChange={e => setLogTo(e.target.value)} aria-label="To" />
              <Button size="sm" onClick={() => { setLogFilter(''); setLogFrom(''); setLogTo(''); }}>Clear</Button>
              <Button size="sm" onClick={handleDownloadLogs}>Download Logs</Button>
            </div>
            <Card>
              <CardHeader><CardTitle>Logs</CardTitle></CardHeader>
              <CardContent>
                <ScrollArea className="h-80"><pre className="whitespace-pre-wrap">{filteredLogs.join('\n')}</pre></ScrollArea>
              </CardContent>
            </Card>
          </TabsContent>

          {/* Backtest Tab */}
          <TabsContent value="backtest">
            <BacktestTab pairs={pairs} />
          </TabsContent>
        </Tabs>
      </div>
      
      {/* Toast container for notifications */}
      <div id="toast-container" className="fixed top-4 right-4 space-y-2 z-50"></div>
      
      {/* Script injection for original JS functionality */}
      <script dangerouslySetInnerHTML={{ 
        __html: `
          // Initialize the chart.js dependency that was originally loaded
          document.addEventListener('DOMContentLoaded', () => {
            // Here we would connect the original JS functionality
            console.log('Trading bot dashboard ready');
          });
        `
      }} />
      <Toaster position="top-right" />
    </div>
  );
};

export default Index;
