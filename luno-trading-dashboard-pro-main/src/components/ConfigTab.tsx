import React, { useState, useEffect } from 'react';
import { Card, CardHeader, CardContent, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Slider } from '@/components/ui/slider';
import { toast } from 'sonner';

interface SavedConfig {
  name: string;
  description?: string;
  cfg: {
    pair: string;
    entry_threshold: number;
    exit_threshold: number;
    stake_size: number;
    cooldown: string;
  };
}

interface ConfigTabProps {
  pairs: string[];
  configObj: any;
  onConfigSave: (cfg: any) => void;
}

export default function ConfigTab({ pairs, configObj, onConfigSave }: ConfigTabProps) {
  const [cfgPair, setCfgPair] = useState<string>(configObj?.pair || '');
  const [cfgEntryThreshold, setCfgEntryThreshold] = useState<number>(configObj?.entry_threshold || 0);
  const [cfgExitThreshold, setCfgExitThreshold] = useState<number>(configObj?.exit_threshold || 0);
  const [cfgStakeSize, setCfgStakeSize] = useState<number>(configObj?.stake_size || 0);
  const [cfgCooldown, setCfgCooldown] = useState<string>(configObj?.cooldown || '');

  const [configError, setConfigError] = useState<string>('');
  const [configFilter, setConfigFilter] = useState<string>('');
  const [savedConfigs, setSavedConfigs] = useState<SavedConfig[]>([]);
  const [selectedPreset, setSelectedPreset] = useState<string>('');

  // Load saved presets
  useEffect(() => {
    const raw = localStorage.getItem('savedConfigs');
    if (raw) {
      try { setSavedConfigs(JSON.parse(raw)); } catch {}
    }
  }, []);

  // Update local state if configObj changes
  useEffect(() => {
    if (configObj) {
      setCfgPair(configObj.pair);
      setCfgEntryThreshold(configObj.entry_threshold);
      setCfgExitThreshold(configObj.exit_threshold);
      setCfgStakeSize(configObj.stake_size);
      setCfgCooldown(configObj.cooldown);
    }
  }, [configObj]);

  const handleSaveConfig = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!cfgPair) { setConfigError('Pair is required'); return; }
    setConfigError('');
    try {
      const newCfg = { pair: cfgPair, entry_threshold: cfgEntryThreshold, exit_threshold: cfgExitThreshold, stake_size: cfgStakeSize, cooldown: cfgCooldown };
      const res = await fetch('/config', { method: 'PUT', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(newCfg) });
      if (!res.ok) throw new Error(`Status ${res.status}`);
      const updated = await res.json();
      onConfigSave(updated);
      toast.success('Config saved');
    } catch (err: any) {
      setConfigError(err.message || 'Failed to save config');
      toast.error('Config save failed');
    }
  };

  const handleSavePreset = () => {
    const name = prompt('Preset name:'); if (!name) return;
    const description = prompt('Preset description:') || '';
    const preset: SavedConfig = { name, description, cfg: { pair: cfgPair, entry_threshold: cfgEntryThreshold, exit_threshold: cfgExitThreshold, stake_size: cfgStakeSize, cooldown: cfgCooldown } };
    const updated = [...savedConfigs.filter(p => p.name !== name), preset];
    setSavedConfigs(updated);
    localStorage.setItem('savedConfigs', JSON.stringify(updated));
    setSelectedPreset(name);
    toast.success('Preset saved');
  };

  const handleSelectPreset = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const name = e.target.value;
    setSelectedPreset(name);
    const p = savedConfigs.find(x => x.name === name);
    if (p) {
      setCfgPair(p.cfg.pair);
      setCfgEntryThreshold(p.cfg.entry_threshold);
      setCfgExitThreshold(p.cfg.exit_threshold);
      setCfgStakeSize(p.cfg.stake_size);
      setCfgCooldown(p.cfg.cooldown);
    }
  };

  const handleDeletePreset = () => {
    if (!selectedPreset) return;
    if (!window.confirm(`Delete preset ${selectedPreset}?`)) return;
    const updated = savedConfigs.filter(p => p.name !== selectedPreset);
    setSavedConfigs(updated);
    localStorage.setItem('savedConfigs', JSON.stringify(updated));
    setSelectedPreset('');
    toast.success('Preset deleted');
  };

  const selectedConfig = savedConfigs.find(p => p.name === selectedPreset);

  return (
    <Card>
      <CardHeader><CardTitle>Configuration</CardTitle></CardHeader>
      <CardContent>
        {configError && <p className="text-red-600">{configError}</p>}
        <div className="mb-4 space-y-2">
          <div>
            <Label>Load Preset</Label>
            <select value={selectedPreset} onChange={handleSelectPreset} className="mt-1 w-full border rounded-md bg-white text-black dark:bg-gray-800 dark:text-white">
              <option value="">-- None --</option>
              {savedConfigs.map(p => <option key={p.name} value={p.name}>{p.name}</option>)}
            </select>
            {selectedConfig?.description && <p className="text-sm text-gray-500 mt-1">{selectedConfig.description}</p>}
          </div>
          <div className="flex space-x-2">
            <Button onClick={handleSavePreset}>Save Preset</Button>
            <Button variant="outline" onClick={handleDeletePreset} disabled={!selectedPreset}>Delete</Button>
          </div>
        </div>
        <form onSubmit={handleSaveConfig} className="space-y-6">
          <div>
            <Label htmlFor="pair">Pair</Label>
            <Input placeholder="Filter pairs" value={configFilter} onChange={e => setConfigFilter(e.target.value)} className="mt-1 w-full mb-1" />
            <select id="pair" value={cfgPair} onChange={e => setCfgPair(e.target.value)} className="mt-2 w-full border rounded-md bg-white text-black dark:bg-gray-800 dark:text-white">
              {pairs.filter(p => p.includes(configFilter)).map(p => <option key={p} value={p}>{p}</option>)}
            </select>
          </div>
          <div>
            <Label htmlFor="entryThreshold">Entry Threshold</Label>
            <Slider id="entryThreshold" min={0} max={1} step={0.01} value={[cfgEntryThreshold]} onValueChange={v => setCfgEntryThreshold(v[0])} />
            <div>Value: {cfgEntryThreshold}</div>
          </div>
          <div>
            <Label htmlFor="exitThreshold">Exit Threshold</Label>
            <Slider id="exitThreshold" min={0} max={1} step={0.01} value={[cfgExitThreshold]} onValueChange={v => setCfgExitThreshold(v[0])} />
            <div>Value: {cfgExitThreshold}</div>
          </div>
          <div><Label htmlFor="stakeSize">Stake Size</Label><Input id="stakeSize" type="number" value={cfgStakeSize} onChange={e => setCfgStakeSize(Number(e.target.value))} /></div>
          <div><Label htmlFor="cooldown">Cooldown</Label><Input id="cooldown" value={cfgCooldown} onChange={e => setCfgCooldown(e.target.value)} /></div>
          <Button type="submit" disabled={!cfgPair || !!configError}>Save Config</Button>
        </form>
      </CardContent>
    </Card>
  );
}
