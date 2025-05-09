import React, { useState, useEffect } from 'react';
import { Card, CardHeader, CardContent, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Loader2 } from 'lucide-react';
import { toast } from 'sonner';

interface ThresholdPreset {
  name: string;
  cfg: {
    selectedPairs: string[];
    sinceMinutes: number;
    feeRate: number;
    gridStart: number;
    gridEnd: number;
    gridStep: number;
  };
}

interface ThresholdsTabProps {
  pairs: string[];
}

export default function ThresholdsTab({ pairs }: ThresholdsTabProps) {
  const [filter, setFilter] = useState<string>('');
  const [selectedPairs, setSelectedPairs] = useState<string[]>([]);
  const [sinceMinutes, setSinceMinutes] = useState<number>(60);
  const [feeRate, setFeeRate] = useState<number>(0.001);
  const [gridStart, setGridStart] = useState<number>(0.01);
  const [gridEnd, setGridEnd] = useState<number>(0.05);
  const [gridStep, setGridStep] = useState<number>(0.005);

  const [presets, setPresets] = useState<ThresholdPreset[]>([]);
  const [selectedPreset, setSelectedPreset] = useState<string>('');

  const [results, setResults] = useState<any[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>('');

  useEffect(() => {
    const raw = localStorage.getItem('thresholdPresets');
    if (raw) {
      try { setPresets(JSON.parse(raw)); } catch {};
    }
  }, []);

  const exportToCsv = (filename: string, rows: any[]) => {
    if (!rows || !rows.length) return;
    const headers = Object.keys(rows[0]);
    const csvContent = [headers.join(','), ...rows.map(row => headers.map(f => JSON.stringify(row[f])).join(','))].join('\r\n');
    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a'); link.href = url;
    link.setAttribute('download', filename); document.body.appendChild(link);
    link.click(); document.body.removeChild(link);
  };

  const handleRun = async (e: React.FormEvent) => {
    e.preventDefault();
    if (selectedPairs.length === 0) { setError('Select at least one pair'); return; }
    setError(''); setLoading(true);
    try {
      const body = { pairs: selectedPairs, since_minutes: sinceMinutes, fee_rate: feeRate, grid_start: gridStart, grid_end: gridEnd, grid_step: gridStep };
      const res = await fetch('/thresholds', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(body) });
      if (!res.ok) throw new Error(`Status ${res.status}`);
      const data = await res.json();
      setResults(data);
      toast.success('Threshold optimization complete');
    } catch (err: any) {
      setError(err.message || 'Optimization failed');
      toast.error('Optimization failed');
    } finally {
      setLoading(false);
    }
  };

  const handleSavePreset = () => {
    const name = prompt('Threshold preset name:'); if (!name) return;
    const preset: ThresholdPreset = { name, cfg: { selectedPairs, sinceMinutes, feeRate, gridStart, gridEnd, gridStep } };
    const updated = [...presets.filter(p => p.name !== name), preset];
    setPresets(updated);
    localStorage.setItem('thresholdPresets', JSON.stringify(updated));
    setSelectedPreset(name);
    toast.success('Threshold preset saved');
  };

  const handleSelectPreset = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const name = e.target.value; setSelectedPreset(name);
    const p = presets.find(x => x.name === name);
    if (p) {
      const c = p.cfg;
      setSelectedPairs(c.selectedPairs);
      setSinceMinutes(c.sinceMinutes);
      setFeeRate(c.feeRate);
      setGridStart(c.gridStart);
      setGridEnd(c.gridEnd);
      setGridStep(c.gridStep);
    }
  };

  const handleDeletePreset = () => {
    if (!selectedPreset) return;
    if (!window.confirm(`Delete threshold preset ${selectedPreset}?`)) return;
    const updated = presets.filter(p => p.name !== selectedPreset);
    setPresets(updated);
    localStorage.setItem('thresholdPresets', JSON.stringify(updated));
    setSelectedPreset('');
    toast.success('Threshold preset deleted');
  };

  return (
    <Card>
      <CardHeader><CardTitle>Threshold Optimization</CardTitle></CardHeader>
      <CardContent>
        {error && <p className="text-red-600 mb-2">{error}</p>}
        <div className="mb-4 space-y-2">
          <Label>Load Preset</Label>
          <div className="flex space-x-2">
            <select value={selectedPreset} onChange={handleSelectPreset} className="border rounded px-2 py-1 bg-white text-black dark:bg-gray-800 dark:text-white">
              <option value="">-- None --</option>
              {presets.map(p => <option key={p.name} value={p.name}>{p.name}</option>)}
            </select>
            <Button size="sm" onClick={handleSavePreset}>Save Preset</Button>
            <Button size="sm" variant="outline" onClick={handleDeletePreset} disabled={!selectedPreset}>Delete</Button>
          </div>
        </div>
        <form onSubmit={handleRun} className="space-y-4">
          <div>
            <Label>Pairs</Label>
            <Input placeholder="Filter pairs" value={filter} onChange={e => setFilter(e.target.value)} className="mt-1 w-full mb-1" disabled={loading} />
            <select multiple value={selectedPairs} onChange={e => setSelectedPairs(Array.from(e.target.selectedOptions).map(o => o.value))} className="w-full h-32 border rounded p-1 bg-white text-black dark:bg-gray-800 dark:text-white" disabled={loading}>
              {pairs.filter(p => p.includes(filter)).map(p => <option key={p} value={p}>{p}</option>)}
            </select>
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <Label>Since (minutes)</Label>
              <Input type="number" value={sinceMinutes} onChange={e => setSinceMinutes(Number(e.target.value))} disabled={loading} />
            </div>
            <div>
              <Label>Fee Rate</Label>
              <Input type="number" step="0.0001" value={feeRate} onChange={e => setFeeRate(Number(e.target.value))} disabled={loading} />
            </div>
            <div>
              <Label>Grid Start</Label>
              <Input type="number" step="0.0001" value={gridStart} onChange={e => setGridStart(Number(e.target.value))} disabled={loading} />
            </div>
            <div>
              <Label>Grid End</Label>
              <Input type="number" step="0.0001" value={gridEnd} onChange={e => setGridEnd(Number(e.target.value))} disabled={loading} />
            </div>
            <div className="col-span-2">
              <Label>Grid Step</Label>
              <Input type="number" step="0.0001" value={gridStep} onChange={e => setGridStep(Number(e.target.value))} disabled={loading} />
            </div>
          </div>
          <Button type="submit" disabled={loading} className="w-full flex items-center justify-center gap-2">
            {loading && <Loader2 className="animate-spin h-4 w-4" />}<span>{loading ? 'Running...' : 'Run Optimization'}</span>
          </Button>
        </form>
        {results.length > 0 && (
          <>
            <div className="my-4">
              <Button size="sm" onClick={() => exportToCsv('threshold_results.csv', results)}>Export CSV</Button>
            </div>
            <ScrollArea className="h-64 border rounded p-2">
              <pre className="whitespace-pre-wrap">{JSON.stringify(results, null, 2)}</pre>
            </ScrollArea>
          </>
        )}
      </CardContent>
    </Card>
  );
}
