import React, { useState, useEffect } from 'react';
import { Card, CardHeader, CardContent, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Button } from '@/components/ui/button';
import { Switch } from '@/components/ui/switch';
import { RefreshCw, Loader2 } from 'lucide-react';
import { ChartContainer } from '@/components/ui/chart';
import * as RechartsPrimitive from 'recharts';
import { toast } from 'sonner';

interface OrderBookTabProps {
  pairs: string[];
}

export default function OrderBookTab({ pairs }: OrderBookTabProps) {
  const [orderbookFilter, setOrderbookFilter] = useState<string>('');
  const [orderbookPair, setOrderbookPair] = useState<string>('');
  const [orderbookData, setOrderbookData] = useState<any>(null);
  const [depthData, setDepthData] = useState<{ price: number; bids: number; asks: number }[]>([]);
  const [orderbookLoading, setOrderbookLoading] = useState<boolean>(false);
  const [orderbookError, setOrderbookError] = useState<string>('');
  const [orderbookAuto, setOrderbookAuto] = useState<boolean>(false);
  const [refreshIntervalSeconds, setRefreshIntervalSeconds] = useState<number>(5);

  useEffect(() => {
    let iv: NodeJS.Timeout;
    if (orderbookAuto && orderbookPair) {
      fetchOrderbook();
      iv = setInterval(fetchOrderbook, refreshIntervalSeconds * 1000);
    }
    return () => iv && clearInterval(iv);
  }, [orderbookAuto, refreshIntervalSeconds, orderbookPair]);

  const fetchOrderbook = async () => {
    if (!orderbookPair) { setOrderbookError('Select a pair'); return; }
    setOrderbookError('');
    setOrderbookLoading(true);
    try {
      const res = await fetch(`/orderbook?pair=${orderbookPair}`);
      if (!res.ok) throw new Error(`Status ${res.status}`);
      const data = await res.json();
      setOrderbookData(data);
      const bidsArr = data.bids.map((o: any) => ({ price: Number(o.price), volume: Number(o.volume) })).sort((a, b) => a.price - b.price);
      const asksArr = data.asks.map((o: any) => ({ price: Number(o.price), volume: Number(o.volume) })).sort((a, b) => a.price - b.price);
      let cum = 0;
      const bidDepth = bidsArr.map(item => { cum += item.volume; return { price: item.price, bids: cum, asks: 0 }; });
      cum = 0;
      const askDepth = asksArr.map(item => { cum += item.volume; return { price: item.price, bids: 0, asks: cum }; });
      const merged = [...bidDepth, ...askDepth].sort((a, b) => a.price - b.price);
      setDepthData(merged);
    } catch (err: any) {
      setOrderbookError(err.message || 'Failed to fetch orderbook');
      toast.error('Failed to fetch orderbook');
    } finally {
      setOrderbookLoading(false);
    }
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center justify-between">
          <span>Order Book</span>
          <span className="text-sm text-gray-500">Every {refreshIntervalSeconds}s</span>
        </CardTitle>
      </CardHeader>
      <CardContent>
        {orderbookError && <p className="text-red-600 mb-2">{orderbookError}</p>}
        <div className="flex items-center space-x-2 mb-4">
          <Input placeholder="Filter pairs" value={orderbookFilter} onChange={e => setOrderbookFilter(e.target.value)} className="w-32" />
          <Label htmlFor="orderbookPair">Pair:</Label>
          <select id="orderbookPair" value={orderbookPair} onChange={e => setOrderbookPair(e.target.value)} className="mt-1 border rounded-md bg-white text-black dark:bg-gray-800 dark:text-white">
            {pairs.filter(p => p.includes(orderbookFilter)).map(p => <option key={p} value={p}>{p}</option>)}
          </select>
          <Button size="sm" variant="outline" onClick={fetchOrderbook} disabled={!orderbookPair || orderbookLoading} className="flex items-center gap-2">
            {orderbookLoading ? <Loader2 className="animate-spin h-4 w-4" /> : <RefreshCw className="h-4 w-4" />}
            <span className="sr-only">Refresh</span>
          </Button>
          <Switch checked={orderbookAuto} onCheckedChange={setOrderbookAuto} />
          <Input type="number" value={refreshIntervalSeconds} min={1} onChange={e => setRefreshIntervalSeconds(Number(e.target.value))} className="w-16" />
        </div>
        <div className="border rounded-md p-4 bg-gray-50 dark:bg-gray-900 mb-4">
          <pre className="text-sm overflow-auto max-h-80">{orderbookData ? JSON.stringify(orderbookData, null, 2) : ''}</pre>
        </div>
        <div className="h-64 border rounded-md p-2 bg-white dark:bg-gray-800">
          {depthData.length > 0 ? (
            <ChartContainer config={{ bids: { color: 'green' }, asks: { color: 'red' } }}>
              <RechartsPrimitive.LineChart data={depthData}>
                <RechartsPrimitive.CartesianGrid strokeDasharray="3 3" />
                <RechartsPrimitive.XAxis dataKey="price" />
                <RechartsPrimitive.YAxis />
                <RechartsPrimitive.Tooltip />
                <RechartsPrimitive.Line dataKey="bids" stroke="var(--color-bids)" dot={false} />
                <RechartsPrimitive.Line dataKey="asks" stroke="var(--color-asks)" dot={false} />
              </RechartsPrimitive.LineChart>
            </ChartContainer>
          ) : (
            <p className="text-center text-sm text-gray-500">Select a pair and refresh</p>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
