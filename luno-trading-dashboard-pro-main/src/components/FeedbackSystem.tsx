import React, { useEffect, useState, useRef } from 'react';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Button } from '@/components/ui/button';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Alert, AlertTitle, AlertDescription } from '@/components/ui/alert';
import { Progress } from '@/components/ui/progress';
import { InfoIcon, AlertTriangleIcon, AlertCircleIcon, CheckCircleIcon, XCircleIcon, RefreshCwIcon, ClockIcon } from 'lucide-react';
import { toast } from 'sonner';

export interface FeedbackEvent {
  id: string;
  timestamp: Date;
  category: 'system' | 'trading' | 'ai' | 'error' | 'success' | 'info' | 'warning';
  title: string;
  message: string;
  duration?: number; // in ms, for operations that take time
  progress?: number; // 0-100 for operations with progress
  source?: string; // component or module that generated the event
  data?: any; // additional data for the event
  acknowledged?: boolean;
}

export interface FeedbackSystemProps {
  feedbackEvents: FeedbackEvent[];
  onClearEvents: () => void;
  onAcknowledgeEvent: (id: string) => void;
  onExportEvents: () => void;
  maxEventsDisplayed?: number;
}

const FeedbackSystem: React.FC<FeedbackSystemProps> = ({
  feedbackEvents,
  onClearEvents,
  onAcknowledgeEvent,
  onExportEvents,
  maxEventsDisplayed = 100
}) => {
  const [activeTab, setActiveTab] = useState('all');
  const [filteredEvents, setFilteredEvents] = useState<FeedbackEvent[]>([]);
  const scrollAreaRef = useRef<HTMLDivElement>(null);

  // Filter events based on active tab
  useEffect(() => {
    let filtered = [...feedbackEvents];
    
    if (activeTab !== 'all') {
      filtered = feedbackEvents.filter(event => 
        activeTab === 'errors' 
          ? event.category === 'error' 
          : activeTab === 'warnings'
            ? event.category === 'warning'
            : activeTab === 'success'
              ? event.category === 'success'
              : activeTab === 'trading'
                ? event.category === 'trading'
                : activeTab === 'ai'
                  ? event.category === 'ai'
                  : true
      );
    }
    
    // Limit to max events
    filtered = filtered.slice(0, maxEventsDisplayed);
    
    setFilteredEvents(filtered);
  }, [feedbackEvents, activeTab, maxEventsDisplayed]);

  // Scroll to bottom when new events arrive
  useEffect(() => {
    if (scrollAreaRef.current) {
      const scrollArea = scrollAreaRef.current;
      scrollArea.scrollTop = scrollArea.scrollHeight;
    }
  }, [filteredEvents]);

  // Format timestamp for display
  const formatTimestamp = (date: Date) => {
    return new Intl.DateTimeFormat('default', {
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      hour12: false
    }).format(date);
  };

  // Get color based on event category
  const getCategoryColor = (category: string) => {
    switch (category) {
      case 'error':
        return 'bg-red-500 text-white';
      case 'warning':
        return 'bg-yellow-500 text-black';
      case 'success':
        return 'bg-green-500 text-white';
      case 'info':
        return 'bg-blue-500 text-white';
      case 'system':
        return 'bg-purple-500 text-white';
      case 'trading':
        return 'bg-cyan-500 text-white';
      case 'ai':
        return 'bg-indigo-500 text-white';
      default:
        return 'bg-gray-500 text-white';
    }
  };

  // Get icon based on event category
  const getCategoryIcon = (category: string) => {
    switch (category) {
      case 'error':
        return <XCircleIcon className="h-4 w-4 text-red-500" />;
      case 'warning':
        return <AlertTriangleIcon className="h-4 w-4 text-yellow-500" />;
      case 'success':
        return <CheckCircleIcon className="h-4 w-4 text-green-500" />;
      case 'info':
        return <InfoIcon className="h-4 w-4 text-blue-500" />;
      case 'system':
        return <RefreshCwIcon className="h-4 w-4 text-purple-500" />;
      case 'trading':
        return <ClockIcon className="h-4 w-4 text-cyan-500" />;
      case 'ai':
        return <InfoIcon className="h-4 w-4 text-indigo-500" />;
      default:
        return <InfoIcon className="h-4 w-4 text-gray-500" />;
    }
  };

  // Format event duration for display
  const formatDuration = (duration: number) => {
    if (duration < 1000) {
      return `${duration}ms`;
    } else if (duration < 60000) {
      return `${(duration / 1000).toFixed(2)}s`;
    } else {
      return `${(duration / 60000).toFixed(2)}m`;
    }
  };

  // Count events by category
  const getEventCounts = () => {
    const counts = {
      all: feedbackEvents.length,
      errors: feedbackEvents.filter(e => e.category === 'error').length,
      warnings: feedbackEvents.filter(e => e.category === 'warning').length,
      success: feedbackEvents.filter(e => e.category === 'success').length,
      trading: feedbackEvents.filter(e => e.category === 'trading').length,
      ai: feedbackEvents.filter(e => e.category === 'ai').length,
      system: feedbackEvents.filter(e => e.category === 'system').length,
    };
    return counts;
  };

  const counts = getEventCounts();

  // Render a single event
  const renderEvent = (event: FeedbackEvent) => (
    <div 
      key={event.id} 
      className={`mb-2 p-2 rounded-md border ${event.acknowledged ? 'opacity-60' : 'opacity-100'}`}
    >
      <div className="flex items-start space-x-2">
        <div className="mt-0.5">
          {getCategoryIcon(event.category)}
        </div>
        <div className="flex-1">
          <div className="flex justify-between items-center">
            <div className="font-medium">{event.title}</div>
            <div className="text-xs text-gray-500 dark:text-gray-400">
              {formatTimestamp(event.timestamp)}
            </div>
          </div>
          <div className="text-sm">{event.message}</div>
          
          {event.duration && (
            <div className="text-xs text-gray-500 dark:text-gray-400 mt-1">
              Duration: {formatDuration(event.duration)}
            </div>
          )}
          
          {event.progress !== undefined && (
            <div className="mt-2">
              <Progress value={event.progress} max={100} className="h-1" />
              <div className="text-xs text-right mt-1">{event.progress}%</div>
            </div>
          )}

          {event.source && (
            <div className="text-xs text-gray-500 dark:text-gray-400 mt-1">
              Source: {event.source}
            </div>
          )}

          {!event.acknowledged && (
            <Button 
              variant="ghost" 
              size="sm" 
              className="mt-1 text-xs h-6 px-2" 
              onClick={() => onAcknowledgeEvent(event.id)}
            >
              Acknowledge
            </Button>
          )}
        </div>
      </div>
    </div>
  );

  // Render stats summary
  const renderStats = () => (
    <div className="flex flex-wrap gap-2 mb-4">
      <Badge className="bg-gray-500">Total: {counts.all}</Badge>
      <Badge className="bg-red-500">Errors: {counts.errors}</Badge>
      <Badge className="bg-yellow-500 text-black">Warnings: {counts.warnings}</Badge>
      <Badge className="bg-green-500">Success: {counts.success}</Badge>
      <Badge className="bg-cyan-500">Trading: {counts.trading}</Badge>
      <Badge className="bg-indigo-500">AI: {counts.ai}</Badge>
      <Badge className="bg-purple-500">System: {counts.system}</Badge>
    </div>
  );

  return (
    <Card className="w-full">
      <CardHeader>
        <div className="flex justify-between items-center">
          <div>
            <CardTitle>System Feedback</CardTitle>
            <CardDescription>Real-time updates and notifications</CardDescription>
          </div>
          <div className="flex space-x-2">
            <Button 
              variant="outline" 
              size="sm" 
              onClick={onExportEvents}
            >
              Export
            </Button>
            <Button 
              variant="outline" 
              size="sm" 
              onClick={onClearEvents}
            >
              Clear
            </Button>
          </div>
        </div>
      </CardHeader>
      
      <CardContent>
        {renderStats()}
        
        <Tabs defaultValue="all" value={activeTab} onValueChange={setActiveTab}>
          <TabsList className="grid grid-cols-7 mb-4">
            <TabsTrigger value="all">All</TabsTrigger>
            <TabsTrigger value="errors">Errors</TabsTrigger>
            <TabsTrigger value="warnings">Warnings</TabsTrigger>
            <TabsTrigger value="success">Success</TabsTrigger>
            <TabsTrigger value="trading">Trading</TabsTrigger>
            <TabsTrigger value="ai">AI</TabsTrigger>
            <TabsTrigger value="system">System</TabsTrigger>
          </TabsList>
          
          <TabsContent value={activeTab}>
            {filteredEvents.length === 0 ? (
              <Alert>
                <InfoIcon className="h-4 w-4" />
                <AlertTitle>No events</AlertTitle>
                <AlertDescription>
                  No {activeTab === 'all' ? '' : activeTab} events to display
                </AlertDescription>
              </Alert>
            ) : (
              <ScrollArea 
                className="h-[400px] pr-4" 
                ref={scrollAreaRef}
              >
                {filteredEvents.map(renderEvent)}
              </ScrollArea>
            )}
          </TabsContent>
        </Tabs>
      </CardContent>
    </Card>
  );
};

export default FeedbackSystem;
