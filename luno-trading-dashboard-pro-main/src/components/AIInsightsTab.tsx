import React, { useState, useEffect, useRef } from 'react';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Separator } from '@/components/ui/separator';
import { Switch } from '@/components/ui/switch';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Badge } from '@/components/ui/badge';
import { Progress } from '@/components/ui/progress';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip';
import { RefreshCw, Loader2, Cpu, TrendingUp, BarChart2, BrainCircuit, Info, ArrowUp, ArrowDown, ArrowRight } from 'lucide-react';
import { toast } from 'sonner';

interface AIInsight {
  pair: string;
  timeframe: string;
  score: number;
  signal: string;
  confidence: number;
  predicted_move: number;
  recommended_size: number;
  pattern_signals?: string[];
  sentiment_score?: number;
  top_features?: Record<string, number>;
  last_updated: string;
}

interface AIInsightsTabProps {
  aiInsights: any[];
  aiEnabled: boolean;
  aiAutoExecute: boolean;
  aiRefreshInterval: number;
  aiLoading: boolean;
  aiError: string;
  onRefresh: () => Promise<void>;
  onToggleAI: () => Promise<void>;
  onToggleAutoExecute: () => void;
  onChangeInterval: (value: number) => void;
  onSaveSettings: () => Promise<void>;
}

interface AIModel {
  last_updated: string;
  last_optimized: string;
  feature_weights: Record<string, number>;
  running_scores: Record<string, number>;
  model_parameters: any;
}

export default function AIInsightsTab({
  aiInsights,
  aiEnabled,
  aiAutoExecute,
  aiRefreshInterval,
  aiLoading,
  aiError,
  onRefresh,
  onToggleAI,
  onToggleAutoExecute,
  onChangeInterval,
  onSaveSettings
}: AIInsightsTabProps) {
  const [opportunities, setOpportunities] = useState<AIInsight[]>(aiInsights || []);
  const [modelInfo, setModelInfo] = useState<AIModel | null>(null);
  const [selectedInsight, setSelectedInsight] = useState<AIInsight | null>(null);
  const [modelLoading, setModelLoading] = useState(false);
  const [optimizeMethod, setOptimizeMethod] = useState('bayesian');
  const [optimizeIterations, setOptimizeIterations] = useState('100');
  const [optimizing, setOptimizing] = useState(false);
  const [minScore, setMinScore] = useState('0.6');

  // Update opportunities when aiInsights prop changes
  useEffect(() => {
    if (aiInsights && aiInsights.length > 0) {
      setOpportunities(aiInsights);
      
      // Select the first insight by default if none selected
      if (!selectedInsight && aiInsights.length > 0) {
        setSelectedInsight(aiInsights[0]);
      }
    }
  }, [aiInsights]);

  // Fetch AI model info
  const fetchModelInfo = async () => {
    setModelLoading(true);
    try {
      const res = await fetch('/api/ai/model');
      if (!res.ok) throw new Error(`Status ${res.status}`);
      const data = await res.json();
      setModelInfo(data);
    } catch (err: any) {
      toast.error(`Failed to fetch AI model info: ${err.message}`);
    } finally {
      setModelLoading(false);
    }
  };

  // Trigger optimization
  const startOptimization = async () => {
    setOptimizing(true);
    try {
      const res = await fetch('/api/ai/optimize', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          method: optimizeMethod,
          iterations: parseInt(optimizeIterations)
        })
      });
      
      if (!res.ok) throw new Error(`Status ${res.status}`);
      
      const data = await res.json();
      toast.success(`Optimization started: ${data.job_id}`);
    } catch (err: any) {
      toast.error(`Failed to start optimization: ${err.message}`);
    } finally {
      setOptimizing(false);
    }
  };
  
  // Initial data fetch for model info
  useEffect(() => {
    fetchModelInfo();
  }, []);
  
  // Handle insight selection
  const handleSelectInsight = (insight: AIInsight) => {
    setSelectedInsight(insight);
  };
  
  // Format date
  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleString();
  };
  
  // Signal indicators
  const SignalIndicator = ({ signal }: { signal: string }) => {
    if (signal === 'buy') {
      return <Badge className="bg-green-500"><ArrowUp size={12} className="mr-1" /> Buy</Badge>
    } else if (signal === 'sell') {
      return <Badge className="bg-red-500"><ArrowDown size={12} className="mr-1" /> Sell</Badge>
    } else {
      return <Badge className="bg-gray-500"><ArrowRight size={12} className="mr-1" /> Hold</Badge>
    }
  };

  // Update opportunities from prop
  useEffect(() => {
    if (aiInsights && aiInsights.length > 0) {
      setOpportunities(aiInsights);
      
      // Select the first insight by default if none selected
      if (!selectedInsight) {
        setSelectedInsight(aiInsights[0]);
      }
    }
  }, [aiInsights]);

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="flex items-center">
              <BrainCircuit className="mr-2" /> AI Market Insights
            </CardTitle>
            <CardDescription>Machine learning powered trading opportunities</CardDescription>
          </div>
          <div className="flex items-center space-x-2">
            <TooltipProvider>
              <Tooltip>
                <TooltipTrigger asChild>
                  <div className="flex items-center space-x-2">
                    <Switch checked={aiEnabled} onCheckedChange={onToggleAI} id="ai-enabled" />
                    <Label htmlFor="ai-enabled">{aiEnabled ? 'AI Enabled' : 'AI Disabled'}</Label>
                  </div>
                </TooltipTrigger>
                <TooltipContent>
                  <p>Enable or disable the AI engine</p>
                </TooltipContent>
              </Tooltip>
            </TooltipProvider>
            
            <Input 
              type="number" 
              value={aiRefreshInterval} 
              onChange={e => onChangeInterval(Number(e.target.value))} 
              min="5"
              max="300"
              className="w-16" 
            />
            <span className="text-sm">sec</span>
            
            <div className="flex justify-end mb-4">
              <Button
                size="sm"
                variant="outline"
                className="mr-2"
                onClick={onRefresh}
                disabled={aiLoading}
              >
                {aiLoading ? <Loader2 className="h-4 w-4 animate-spin mr-2" /> : <RefreshCw className="h-4 w-4 mr-2" />}
                Refresh Insights
              </Button>
            </div>
          </div>
        </div>
      </CardHeader>
      
      <CardContent>
        <Tabs defaultValue="opportunities" className="w-full">
          <TabsList className="mb-4">
            <TabsTrigger value="opportunities">Opportunities</TabsTrigger>
            <TabsTrigger value="model">AI Model</TabsTrigger>
            <TabsTrigger value="optimization">Optimization</TabsTrigger>
          </TabsList>
          
          <TabsContent value="opportunities" className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <div className="md:col-span-1">
                <Card>
                  <CardHeader>
                    <CardTitle className="text-sm font-medium">Top Opportunities</CardTitle>
                    <div className="flex items-center">
                      <Label htmlFor="min-score" className="mr-2">Min Score:</Label>
                      <Input 
                        id="min-score" 
                        type="number" 
                        value={minScore} 
                        onChange={e => setMinScore(e.target.value)} 
                        className="w-20"
                        step="0.05"
                        min="0"
                        max="1"
                      />
                      <Button 
                        variant="ghost" 
                        size="sm" 
                        className="ml-2" 
                        onClick={fetchOpportunities}
                        disabled={loading}
                      >
                        Apply
                      </Button>
                    </div>
                  </CardHeader>
                  <CardContent>
                    <ScrollArea className="h-[500px]">
                      <div className="space-y-2">
                        {opportunities.length === 0 ? (
                          <div className="text-center py-8 text-gray-500">
                            {loading ? (
                              <div className="flex flex-col items-center">
                                <Loader2 className="h-8 w-8 animate-spin mb-2" />
                                <p>Loading insights...</p>
                              </div>
                            ) : (
                              <p>No opportunities found above score threshold</p>
                            )}
                          </div>
                        ) : (
                          opportunities.map((opp) => (
                            <div 
                              key={`${opp.pair}-${opp.timeframe}`}
                              className={`p-3 border rounded-md cursor-pointer hover:bg-gray-100 dark:hover:bg-gray-800 ${
                                selectedInsight && selectedInsight.pair === opp.pair && selectedInsight.timeframe === opp.timeframe
                                  ? 'border-blue-500 bg-blue-50 dark:bg-blue-900/20'
                                  : 'border-gray-200 dark:border-gray-800'
                              }`}
                              onClick={() => handleSelectInsight(opp)}
                            >
                              <div className="flex justify-between items-center mb-2">
                                <div className="font-medium">{opp.pair}</div>
                                <SignalIndicator signal={opp.signal} />
                              </div>
                              <div className="flex justify-between text-sm text-gray-500 dark:text-gray-400 mb-1">
                                <div>Timeframe: {opp.timeframe}</div>
                                <div>Score: {opp.score.toFixed(2)}</div>
                              </div>
                              <Progress value={opp.score * 100} className="h-2" />
                            </div>
                          ))
                        )}
                      </div>
                    </ScrollArea>
                  </CardContent>
                </Card>
              </div>
              
              <div className="md:col-span-2">
                {selectedInsight ? (
                  <Card>
                    <CardHeader>
                      <div className="flex justify-between items-center">
                        <CardTitle>
                          {selectedInsight.pair} ({selectedInsight.timeframe})
                        </CardTitle>
                        <SignalIndicator signal={selectedInsight.signal} />
                      </div>
                      <CardDescription>
                        Last updated: {formatDate(selectedInsight.last_updated)}
                      </CardDescription>
                    </CardHeader>
                    <CardContent>
                      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                        <div className="space-y-4">
                          <div>
                            <h3 className="text-sm font-medium mb-2">Analysis Score</h3>
                            <div className="flex items-center justify-between mb-1">
                              <span>ML Confidence</span>
                              <span className="font-medium">{(selectedInsight.confidence * 100).toFixed(1)}%</span>
                            </div>
                            <Progress value={selectedInsight.confidence * 100} className="h-2 mb-3" />
                            
                            <div className="flex items-center justify-between mb-1">
                              <span>Opportunity Score</span>
                              <span className="font-medium">{(selectedInsight.score * 100).toFixed(1)}%</span>
                            </div>
                            <Progress value={selectedInsight.score * 100} className="h-2" />
                          </div>
                          
                          <div>
                            <h3 className="text-sm font-medium mb-2">Prediction</h3>
                            <div className="p-3 rounded-md bg-gray-100 dark:bg-gray-800">
                              <div className="flex justify-between items-center mb-2">
                                <span>Expected Move</span>
                                <span className={`font-medium ${
                                  selectedInsight.predicted_move > 0
                                    ? 'text-green-500'
                                    : selectedInsight.predicted_move < 0
                                    ? 'text-red-500'
                                    : ''
                                }`}>
                                  {selectedInsight.predicted_move > 0 ? '+' : ''}
                                  {selectedInsight.predicted_move.toFixed(2)}%
                                </span>
                              </div>
                              <div className="flex justify-between items-center">
                                <span>Suggested Size</span>
                                <span className="font-medium">
                                  {(selectedInsight.recommended_size * 100).toFixed(1)}%
                                </span>
                              </div>
                            </div>
                          </div>
                          
                          {selectedInsight.sentiment_score !== undefined && (
                            <div>
                              <h3 className="text-sm font-medium mb-2">Market Sentiment</h3>
                              <div className="p-3 rounded-md bg-gray-100 dark:bg-gray-800">
                                <div className="flex justify-between items-center">
                                  <span>Sentiment Score</span>
                                  <span className={`font-medium ${
                                    selectedInsight.sentiment_score > 0.2
                                      ? 'text-green-500'
                                      : selectedInsight.sentiment_score < -0.2
                                      ? 'text-red-500'
                                      : ''
                                  }`}>
                                    {selectedInsight.sentiment_score.toFixed(2)}
                                  </span>
                                </div>
                              </div>
                            </div>
                          )}
                        </div>
                        
                        <div className="space-y-4">
                          {selectedInsight.pattern_signals && selectedInsight.pattern_signals.length > 0 && (
                            <div>
                              <h3 className="text-sm font-medium mb-2">Detected Patterns</h3>
                              <div className="flex flex-wrap gap-2">
                                {selectedInsight.pattern_signals.map((pattern, idx) => (
                                  <Badge key={idx} variant="outline" className="capitalize">
                                    {pattern.replace(/_/g, ' ')}
                                  </Badge>
                                ))}
                              </div>
                            </div>
                          )}
                          
                          {selectedInsight.top_features && Object.keys(selectedInsight.top_features).length > 0 && (
                            <div>
                              <h3 className="text-sm font-medium mb-2">Key Indicators</h3>
                              <div className="space-y-3">
                                {Object.entries(selectedInsight.top_features).map(([feature, value]) => (
                                  <div key={feature}>
                                    <div className="flex justify-between items-center mb-1">
                                      <span className="capitalize">{feature.replace(/([A-Z])/g, ' $1').toLowerCase()}</span>
                                      <span className="font-medium">{value.toFixed(2)}</span>
                                    </div>
                                    <Progress value={value * 100} className="h-1.5" />
                                  </div>
                                ))}
                              </div>
                            </div>
                          )}
                          
                          <div className="mt-4">
                            <Button variant="outline" className="w-full" disabled>
                              Execute {selectedInsight.signal.toUpperCase()} Trade
                            </Button>
                          </div>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                ) : (
                  <div className="flex items-center justify-center h-full">
                    <div className="text-center text-gray-500">
                      <Cpu className="h-12 w-12 mx-auto mb-4" />
                      <p>Select an opportunity to view detailed analysis</p>
                    </div>
                  </div>
                )}
              </div>
            </div>
          </TabsContent>
          
          <TabsContent value="model" className="space-y-4">
            {modelLoading ? (
              <div className="flex justify-center py-12">
                <Loader2 className="h-8 w-8 animate-spin" />
              </div>
            ) : modelInfo ? (
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <Card>
                  <CardHeader>
                    <CardTitle className="text-sm font-medium">
                      Model Information
                    </CardTitle>
                    <CardDescription>
                      Last Updated: {formatDate(modelInfo.last_updated)}<br />
                      Last Optimized: {formatDate(modelInfo.last_optimized)}
                    </CardDescription>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-6">
                      <div>
                        <h3 className="text-sm font-medium mb-3">Feature Importance</h3>
                        <div className="space-y-3">
                          {Object.entries(modelInfo.feature_weights || {})
                            .sort((a, b) => b[1] - a[1]) // Sort by weight descending
                            .slice(0, 8) // Take top 8
                            .map(([feature, weight]) => (
                              <div key={feature}>
                                <div className="flex justify-between items-center mb-1">
                                  <span className="capitalize">{feature.replace(/([A-Z])/g, ' $1').toLowerCase()}</span>
                                  <span className="font-medium">{(weight as number).toFixed(2)}</span>
                                </div>
                                <Progress value={(weight as number) * 100} className="h-1.5" />
                              </div>
                            ))}
                        </div>
                      </div>
                    </div>
                  </CardContent>
                </Card>
                
                <Card>
                  <CardHeader>
                    <CardTitle className="text-sm font-medium">
                      Running Market Scores
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    <ScrollArea className="h-[400px]">
                      <div className="space-y-3">
                        {Object.entries(modelInfo.running_scores || {})
                          .sort((a, b) => b[1] - a[1]) // Sort by score descending
                          .map(([pair, score]) => (
                            <div key={pair}>
                              <div className="flex justify-between items-center mb-1">
                                <span>{pair}</span>
                                <span className={`font-medium ${
                                  score > 0.7 ? 'text-green-500' : 
                                  score < 0.3 ? 'text-red-500' : 
                                  'text-yellow-500'
                                }`}>
                                  {(score as number).toFixed(2)}
                                </span>
                              </div>
                              <Progress 
                                value={(score as number) * 100} 
                                className={`h-2 ${
                                  score > 0.7 ? 'bg-green-500' : 
                                  score < 0.3 ? 'bg-red-500' : 
                                  'bg-yellow-500'
                                }`} 
                              />
                            </div>
                          ))}
                      </div>
                    </ScrollArea>
                  </CardContent>
                </Card>
              </div>
            ) : (
              <div className="text-center py-12 text-gray-500">
                <p>Failed to load AI model information</p>
              </div>
            )}
          </TabsContent>
          
          <TabsContent value="optimization" className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle>Parameter Optimization</CardTitle>
                <CardDescription>
                  Run AI optimization to improve trading parameters
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div className="space-y-4">
                    <div>
                      <Label htmlFor="optimize-method">Optimization Method</Label>
                      <select 
                        id="optimize-method" 
                        value={optimizeMethod}
                        onChange={e => setOptimizeMethod(e.target.value)}
                        className="w-full p-2 border rounded-md mt-1 bg-white dark:bg-gray-800"
                      >
                        <option value="random">Random Search</option>
                        <option value="bayesian">Bayesian Optimization</option>
                        <option value="walkforward">Walk-Forward Optimization</option>
                      </select>
                    </div>
                    
                    <div>
                      <Label htmlFor="optimize-iterations">Iterations</Label>
                      <Input 
                        id="optimize-iterations" 
                        type="number" 
                        value={optimizeIterations} 
                        onChange={e => setOptimizeIterations(e.target.value)} 
                        className="mt-1"
                        min="10"
                        max="500"
                      />
                    </div>
                    
                    <div className="pt-4">
                      <Button 
                        onClick={startOptimization} 
                        disabled={optimizing}
                        className="w-full"
                      >
                        {optimizing ? 
                          <><Loader2 className="mr-2 h-4 w-4 animate-spin" /> Optimizing...</> : 
                          'Start Optimization'
                        }
                      </Button>
                    </div>
                  </div>
                  
                  <div className="border rounded-md p-4">
                    <h3 className="font-medium mb-2">Optimization Methods</h3>
                    <div className="space-y-3 text-sm">
                      <div>
                        <h4 className="font-medium">Random Search</h4>
                        <p className="text-gray-500 dark:text-gray-400">Fast method that randomly samples parameter space. Good for initial exploration.</p>
                      </div>
                      <div>
                        <h4 className="font-medium">Bayesian Optimization</h4>
                        <p className="text-gray-500 dark:text-gray-400">Intelligent method that learns from previous trials. Efficient for finding optimal parameters.</p>
                      </div>
                      <div>
                        <h4 className="font-medium">Walk-Forward Optimization</h4>
                        <p className="text-gray-500 dark:text-gray-400">Tests parameters on rolling historical windows to ensure robustness across different market regimes.</p>
                      </div>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </CardContent>
    </Card>
  );
}
