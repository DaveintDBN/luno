import React from 'react';
import { Card, CardHeader, CardContent, CardTitle } from '@/components/ui/card';

interface BalanceCardProps {
  balances: { asset: string; balance: number | string }[];
  configPair?: string;
  percentChange: number | null;
}

export default function BalanceCard({ balances, configPair = '', percentChange }: BalanceCardProps) {
  const asset = configPair.slice(-3) || '';
  const entry = balances.find(b => b.asset === asset);
  const balance = entry?.balance ?? '0';
  const changeText = percentChange !== null ? (percentChange >= 0 ? '+' : '') + percentChange.toFixed(2) + '%' : '-';
  const changeClass = percentChange === null
    ? ''
    : percentChange >= 0
      ? 'text-green-600 dark:text-green-400'
      : 'text-red-600 dark:text-red-400';

  return (
    <Card className="bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700">
      <CardHeader className="pb-2">
        <CardTitle className="text-sm font-medium text-gray-500 dark:text-gray-400">Balance</CardTitle>
      </CardHeader>
      <CardContent>
        <p className="text-2xl font-mono font-semibold">{balance} {asset}</p>
        <p className={'text-xs ' + changeClass}>{changeText}</p>
      </CardContent>
    </Card>
  );
}
