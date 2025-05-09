import React, { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Loader2 } from 'lucide-react';
import { toast } from 'sonner';

export default function SimulationTab() {
  const [simLoading, setSimLoading] = useState(false);
  const [simResult, setSimResult] = useState<any>(null);
  const [simError, setSimError] = useState<string>('');

  const handleSimulate = async () => {
    setSimError('');
    setSimLoading(true);
    try {
      const res = await fetch('/simulate', { method: 'POST' });
      if (!res.ok) throw new Error(`Status ${res.status}`);
      const data = await res.json();
      setSimResult(data);
      toast.success('Simulation complete');
    } catch (err: any) {
      setSimError(err.message || 'Simulation failed');
      toast.error('Simulation failed');
    } finally {
      setSimLoading(false);
    }
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Simulation</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          <Button onClick={handleSimulate} disabled={simLoading} className="w-full">
            {simLoading ? <Loader2 className="animate-spin h-4 w-4 inline" /> : 'Run Simulation'}
          </Button>
          {simError && <p className="text-red-600">{simError}</p>}
          {simResult && (
            <div className="space-y-2">
              <p>Signal: <strong>{simResult.signal}</strong></p>
              <p>Position: <strong>{simResult.position}</strong></p>
              <p>Total PnL: <strong>{simResult.total_pnl?.toFixed(2)}</strong></p>
              <p>Max Drawdown Exceeded: <strong>{simResult.max_drawdown_exceeded ? 'Yes' : 'No'}</strong></p>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
