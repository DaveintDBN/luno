import React, { useState, useEffect } from 'react';
import { Card, CardHeader, CardContent, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { ChartContainer } from '@/components/ui/chart';
import * as RechartsPrimitive from 'recharts';
import { Loader2 } from 'lucide-react';
import { toast } from 'sonner';

interface BacktestPreset {
  name: string;
  cfg: {
    btPair: string;
    btSinceMinutes: number;
    btShortWindow: number;
    btLongWindow: number;
    btFeeRate: number;
  };
}

export default function BacktestTab({ pairs }: { pairs: string[] }) {
  const [btPair, setBtPair] = useState<string>('');
  const [btSinceMinutes, setBtSinceMinutes] = useState<number>(60);
  const [btShortWindow, setBtShortWindow] = useState<number>(10);
  const [btLongWindow, setBtLongWindow] = useState<number>(50);
  const [btFeeRate, setBtFeeRate] = useState<number>(0.001);

  const [btMetrics, setBtMetrics] = useState<any>(null);
  const [btPnlHistory, setBtPnlHistory] = useState<any[]>([]);
  const [btDrawdownHistory, setBtDrawdownHistory] = useState<any[]>([]);

  const [btError, setBtError] = useState<string>('');
  const [btLoading, setBtLoading] = useState<boolean>(false);

  const [backtestPresets, setBacktestPresets] = useState<BacktestPreset[]>([]);
  const [selectedBacktestPreset, setSelectedBacktestPreset] = useState<string>('');

  // Load presets
  useEffect(() => {
    const raw = localStorage.getItem('backtestPresets');
    if (raw) {
      try { setBacktestPresets(JSON.parse(raw)); } catch {}
    }
  }, []);

  const handleRunBacktest = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!btPair) { setBtError('Pair is required'); return; }
    setBtError(''); setBtLoading(true);
    try {
      const body = { pair: btPair, since_minutes: btSinceMinutes, short: btShortWindow, long: btLongWindow, fee_rate: btFeeRate };
      const res = await fetch('/backtest', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(body) });
      if (!res.ok) throw new Error(`Status ${res.status}`);
      const data = await res.json();
      setBtMetrics(data);
      setBtPnlHistory(data.pnl_history);
      setBtDrawdownHistory(data.drawdown_history);
      toast.success('Backtest complete');
    } catch (err: any) {
      setBtError(err.message || 'Backtest failed');
      toast.error('Backtest failed');
    } finally {
      setBtLoading(false);
    }
  };

  const handleSaveBacktestPreset = () => {
    const name = prompt('Backtest preset name:'); if (!name) return;
    const preset: BacktestPreset = { name, cfg: { btPair, btSinceMinutes, btShortWindow, btLongWindow, btFeeRate } };
    const updated = [...backtestPresets.filter(p => p.name !== name), preset];
    setBacktestPresets(updated);
    localStorage.setItem('backtestPresets', JSON.stringify(updated));
    setSelectedBacktestPreset(name);
    toast.success('Backtest preset saved');
  };

  const handleSelectBacktestPreset = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const name = e.target.value; setSelectedBacktestPreset(name);
    const p = backtestPresets.find(x => x.name === name);
    if (p) {
      const c = p.cfg;
      setBtPair(c.btPair);
      setBtSinceMinutes(c.btSinceMinutes);
      setBtShortWindow(c.btShortWindow);
      setBtLongWindow(c.btLongWindow);
      setBtFeeRate(c.btFeeRate);
    }
  };

  const handleDeleteBacktestPreset = () => {
    if (!selectedBacktestPreset) return;
    if (!confirm(`Delete backtest preset ${selectedBacktestPreset}?`)) return;
    const updated = backtestPresets.filter(p => p.name !== selectedBacktestPreset);
    setBacktestPresets(updated);
    localStorage.setItem('backtestPresets', JSON.stringify(updated));
    setSelectedBacktestPreset('');
    toast.success('Backtest preset deleted');
  };

  // CSV export
  const exportToCsv = (filename: string, rows: any[]) => {
    if (!rows || !rows.length) return;
    const headers = Object.keys(rows[0]);
    const csvContent = [headers.join(','), ...rows.map(row => headers.map(field => JSON.stringify(row[field])).join(','))].join('\r\n');
    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a'); link.href = url;
    link.setAttribute('download', filename); document.body.appendChild(link);
    link.click(); document.body.removeChild(link);
  };

  return (
    <Card>
      <CardHeader><CardTitle>Backtest</CardTitle></CardHeader>
      <CardContent>
        {btError && <p className="text-red-600">{btError}</p>}
        <form onSubmit={handleRunBacktest} className="space-y-4">
          <div>
            <Label>Pair</Label>
            <select value={btPair} onChange={e => setBtPair(e.target.value)} className="mt-1 w-full border rounded-md">
              {pairs.map(p => <option key={p} value={p}>{p}</option>)}
            </select>
          </div>
          <div><Label>Since (minutes)</Label><Input type="number" value={btSinceMinutes} onChange={e => setBtSinceMinutes(Number(e.target.value))} /></div>
          <div><Label>Short Window</Label><Input type="number" value={btShortWindow} onChange={e => setBtShortWindow(Number(e.target.value))} /></div>
          <div><Label>Long Window</Label><Input type="number" value={btLongWindow} onChange={e => setBtLongWindow(Number(e.target.value))} /></div>
          <div><Label>Fee Rate</Label><Input type="number" step="0.0001" value={btFeeRate} onChange={e => setBtFeeRate(Number(e.target.value))} /></div>
          <Button type="submit" disabled={!btPair || btLoading}>
            {btLoading ? <Loader2 className="animate-spin h-4 w-4 inline" /> : 'Run Backtest'}
          </Button>
        </form>
        <div className="flex items-center space-x-2 my-4">
          <select value={selectedBacktestPreset} onChange={handleSelectBacktestPreset} className="border rounded px-2 py-1">
            <option value="">Load preset</option>
            {backtestPresets.map(p => <option key={p.name} value={p.name}>{p.name}</option>)}
          </select>
          <Button size="sm" onClick={handleSaveBacktestPreset}>Save Preset</Button>
          <Button size="sm" onClick={handleDeleteBacktestPreset} disabled={!selectedBacktestPreset}>Delete</Button>
        </div>
        {btMetrics && (
          <>
            <div className="grid grid-cols-3 gap-4 mb-4">
              <Card><CardHeader><CardTitle>Trades</CardTitle></CardHeader><CardContent>{btMetrics.trades}</CardContent></Card>
              <Card><CardHeader><CardTitle>Win %</CardTitle></CardHeader><CardContent>{btMetrics.win_rate.toFixed(1)}%</CardContent></Card>
              <Card><CardHeader><CardTitle>Total PnL</CardTitle></CardHeader><CardContent>{btMetrics.total_pnl.toFixed(2)}</CardContent></Card>
            </div>
            <h4 className="text-sm font-medium">PnL over time</h4>
            <div className="h-40 overflow-auto mb-4">
              <ChartContainer config={{ pnl: { color: 'blue' } }}>
                <RechartsPrimitive.LineChart data={btPnlHistory}>
                  <RechartsPrimitive.CartesianGrid strokeDasharray="3 3" />
                  <RechartsPrimitive.XAxis type="number" dataKey="time" tickFormatter={t => new Date(t).toLocaleTimeString()} />
                  <RechartsPrimitive.YAxis />
                  <RechartsPrimitive.Tooltip />
                  <RechartsPrimitive.Line dataKey="pnl" stroke="var(--color-pnl)" dot={false} />
                </RechartsPrimitive.LineChart>
              </ChartContainer>
            </div>
            <h4 className="text-sm font-medium">Drawdown over time</h4>
            <div className="h-40 overflow-auto mb-4">
              <ChartContainer config={{ drawdown: { color: 'red' } }}>
                <RechartsPrimitive.LineChart data={btDrawdownHistory}>
                  <RechartsPrimitive.CartesianGrid strokeDasharray="3 3" />
                  <RechartsPrimitive.XAxis type="number" dataKey="time" tickFormatter={t => new Date(t).toLocaleTimeString()} />
                  <RechartsPrimitive.YAxis />
                  <RechartsPrimitive.Tooltip />
                  <RechartsPrimitive.Line dataKey="drawdown" stroke="var(--color-drawdown)" dot={false} />
                </RechartsPrimitive.LineChart>
              </ChartContainer>
            </div>
            <div className="flex space-x-2">
              <Button size="sm" onClick={() => exportToCsv('backtest_pnl_history.csv', btPnlHistory)}>Export PnL CSV</Button>
              <Button size="sm" onClick={() => exportToCsv('backtest_drawdown_history.csv', btDrawdownHistory)}>Export Drawdown CSV</Button>
            </div>
          </>
        )}
      </CardContent>
    </Card>
  );
}
