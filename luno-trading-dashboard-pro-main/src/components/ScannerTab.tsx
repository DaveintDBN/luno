import React, { useState, useEffect, useRef } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Separator } from '@/components/ui/separator';
import { Switch } from '@/components/ui/switch';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip';
import { RefreshCw, Loader2, Info } from 'lucide-react';
import { toast } from 'sonner';

interface ScannerPreset {
  name: string;
  cfg: {
    selectedPairs: string[];
    scanMinVolume: string;
    scanEntryThreshold: string;
    scanExitThreshold: string;
    scanIntervalSec: string;
    scanAutoExecute: boolean;
  };
}

export default function ScannerTab({ pairs }: { pairs: string[] }) {
  const [scannerFilter, setScannerFilter] = useState('');
  const [selectedPairs, setSelectedPairs] = useState<string[]>([]);
  const [scanMinVolume, setScanMinVolume] = useState('0');
  const [scanEntryThreshold, setScanEntryThreshold] = useState('0');
  const [scanExitThreshold, setScanExitThreshold] = useState('0');
  const [scanIntervalSec, setScanIntervalSec] = useState('10');
  const [scanAutoExecute, setScanAutoExecute] = useState(false);

  const [scannerPresets, setScannerPresets] = useState<ScannerPreset[]>([]);
  const [selectedScannerPreset, setSelectedScannerPreset] = useState('');

  const [scannerSaved, setScannerSaved] = useState(false);
  const [scanError, setScanError] = useState('');
  const [scanLoading, setScanLoading] = useState(false);
  const [autoScanning, setAutoScanning] = useState(false);
  const [autoScanLoading, setAutoScanLoading] = useState(false);

  const [scanResults, setScanResults] = useState<any[]>([]);
  const [topPicks, setTopPicks] = useState<any[]>([]);
  const [topLimit, setTopLimit] = useState(5);
  const [autoRefreshPicks, setAutoRefreshPicks] = useState(false);
  const [picksInterval, setPicksInterval] = useState(10);
  const [streamPicks, setStreamPicks] = useState(false);
  const picksIntervalRef = useRef<NodeJS.Timeout | null>(null);

  // CSV export helper
  const exportToCsv = (filename: string, rows: any[]) => {
    if (!rows || !rows.length) return;
    const headers = Object.keys(rows[0]);
    const csvContent = [headers.join(','), ...rows.map(row => headers.map(field => JSON.stringify(row[field])).join(','))].join('\r\n');
    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a'); link.href = url; link.setAttribute('download', filename);
    document.body.appendChild(link); link.click(); document.body.removeChild(link);
  };

  // Load settings & presets
  useEffect(() => {
    const saved = localStorage.getItem('scannerSettings');
    if (saved) {
      try {
        const cfg = JSON.parse(saved);
        setSelectedPairs(cfg.selectedPairs || []);
        setScanMinVolume(String(cfg.scanMinVolume || '0'));
        setScanEntryThreshold(String(cfg.scanEntryThreshold || '0'));
        setScanExitThreshold(String(cfg.scanExitThreshold || '0'));
      } catch {}
    }
    const raw = localStorage.getItem('scannerPresets');
    if (raw) {
      try { setScannerPresets(JSON.parse(raw)); } catch {}
    }
  }, []);

  const handleScan = async () => {
    setScanError(''); setScanLoading(true);
    try {
      const body = { pairs: selectedPairs, min_volume: parseFloat(scanMinVolume)||0, entry_threshold: parseFloat(scanEntryThreshold)||0, exit_threshold: parseFloat(scanExitThreshold)||0 };
      const res = await fetch('/scan',{ method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body) });
      if (!res.ok) throw new Error(`Status ${res.status}`);
      const data = await res.json(); setScanResults(data); toast.success('Scan complete');
    } catch(err:any) {
      setScanError(err.message||'Scan failed'); toast.error('Scan failed');
    } finally { setScanLoading(false); }
  };

  const handleSaveScannerSettings = () => { const cfg = { selectedPairs, scanMinVolume, scanEntryThreshold, scanExitThreshold }; localStorage.setItem('scannerSettings', JSON.stringify(cfg)); setScannerSaved(true); setTimeout(() => setScannerSaved(false), 2000); };

  const handleStartAutoScan = () => {
    setAutoScanLoading(true);
    try {
      setAutoScanning(true);
      const interval = Number(scanIntervalSec)*1000;
      const doScan = async () => { await handleScan(); };
      doScan(); picksIntervalRef.current = setInterval(doScan, interval);
      toast.success('Auto-scan started');
    } catch { setAutoScanning(false); toast.error('Failed to start auto-scan'); }
    finally { setAutoScanLoading(false); }
  };

  const handleStopAutoScan = () => {
    setAutoScanLoading(true);
    try {
      if (picksIntervalRef.current) clearInterval(picksIntervalRef.current);
      setAutoScanning(false); toast.success('Auto-scan stopped');
    } catch { toast.error('Failed to stop auto-scan'); }
    finally { setAutoScanLoading(false); }
  };

  const handleSaveScannerPreset = () => {
    const name = prompt('Scanner preset name:'); if (!name) return;
    const preset = { name, cfg: { selectedPairs, scanMinVolume, scanEntryThreshold, scanExitThreshold, scanIntervalSec, scanAutoExecute } };
    const updated = [...scannerPresets.filter(p=>p.name!==name), preset];
    setScannerPresets(updated); localStorage.setItem('scannerPresets', JSON.stringify(updated)); setSelectedScannerPreset(name);
    toast.success('Scanner preset saved');
  };

  const handleSelectScannerPreset = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const name = e.target.value; setSelectedScannerPreset(name);
    const p = scannerPresets.find(x=>x.name===name);
    if (p) {
      const c = p.cfg;
      setSelectedPairs(c.selectedPairs); setScanMinVolume(c.scanMinVolume); setScanEntryThreshold(c.scanEntryThreshold); setScanExitThreshold(c.scanExitThreshold); setScanIntervalSec(c.scanIntervalSec); setScanAutoExecute(c.scanAutoExecute);
    }
  };

  const handleDeleteScannerPreset = () => { if (!selectedScannerPreset) return; if (!confirm(`Delete scanner preset ${selectedScannerPreset}?`)) return; const updated = scannerPresets.filter(p=>p.name!==selectedScannerPreset); setScannerPresets(updated); localStorage.setItem('scannerPresets', JSON.stringify(updated)); setSelectedScannerPreset(''); toast.success('Scanner preset deleted'); };

  const handleFetchTopPicks = async () => { const body = { pairs: selectedPairs, min_volume: parseFloat(scanMinVolume)||0, limit: topLimit }; const res = await fetch('/opportunities',{ method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body) }); if (!res.ok) return; const data = await res.json(); setTopPicks(data); };

  useEffect(() => {
    let id: NodeJS.Timeout;
    if (autoRefreshPicks) { handleFetchTopPicks(); id = setInterval(handleFetchTopPicks, picksInterval*1000); }
    return () => clearInterval(id);
  }, [autoRefreshPicks, picksInterval, selectedPairs]);

  return (
    <Card>
      <CardHeader><CardTitle>Market Scanner</CardTitle></CardHeader>
      <CardContent>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div className="md:col-span-1 space-y-4">
            <div>
              <Label>Pairs</Label>
              <Input placeholder="Filter pairs" value={scannerFilter} onChange={e=>setScannerFilter(e.target.value)} className="w-full mt-2 mb-1" disabled={autoScanning||scanLoading}/>
              <select multiple value={selectedPairs} onChange={e=>setSelectedPairs(Array.from(e.target.selectedOptions).map(o=>o.value))} className="w-full h-32 mt-1 border rounded-md bg-white text-black dark:bg-gray-800 dark:text-white" disabled={autoScanning||scanLoading}>
                {pairs.filter(p=>p.includes(scannerFilter)).map(p=> <option key={p} value={p}>{p}</option>)}
              </select>
            </div>
            ... (UI for thresholds, buttons, presets, above) ...
          </div>
          ... (Top picks and scanResult UI) ...
        </div>
      </CardContent>
    </Card>
  );
}
