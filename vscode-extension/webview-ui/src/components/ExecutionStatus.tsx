import { useMemo } from 'react';
import { useStore } from '../store';
import { SectionCard } from './SectionCard';
import { ProgressBar } from './ProgressBar';
import { RecentTaskList } from './RecentTaskList';

type ExecutionState = 'running' | 'idle' | 'completed' | 'has_failures';

export function ExecutionStatus() {
  const { tasks } = useStore();

  const stats = useMemo(() => {
    const total = tasks.length;
    const pending = tasks.filter(t => t.status === 'pending').length;
    const doing = tasks.filter(t => t.status === 'doing').length;
    const done = tasks.filter(t => t.status === 'done').length;
    const failed = tasks.filter(t => t.status === 'failed').length;
    const progress = total > 0 ? (done / total) * 100 : 0;

    return { total, pending, doing, done, failed, progress };
  }, [tasks]);

  const executionState = useMemo((): ExecutionState => {
    if (stats.doing > 0) return 'running';
    if (stats.failed > 0) return 'has_failures';
    if (stats.total > 0 && stats.done === stats.total) return 'completed';
    return 'idle';
  }, [stats]);

  if (stats.total === 0) {
    return (
      <SectionCard title="Execution Status">
        <div className="text-sm opacity-50">No tasks created yet</div>
      </SectionCard>
    );
  }

  return (
    <SectionCard title="Execution Status">
      <div className="space-y-4">
        {/* Progress */}
        <ProgressBar
          progress={stats.progress}
          total={stats.total}
          done={stats.done}
        />

        {/* Status */}
        <div className="flex items-center gap-2">
          <span className="text-sm opacity-70">Status:</span>
          <ExecutionStateIndicator state={executionState} />
        </div>

        {/* Stats summary */}
        <div className="flex flex-wrap gap-2 text-xs">
          <StatBadge label="Pending" count={stats.pending} color="gray" />
          <StatBadge label="Running" count={stats.doing} color="blue" />
          <StatBadge label="Done" count={stats.done} color="green" />
          <StatBadge label="Failed" count={stats.failed} color="red" />
        </div>

        {/* Recent tasks */}
        <RecentTaskList tasks={tasks} limit={5} />
      </div>
    </SectionCard>
  );
}

function ExecutionStateIndicator({ state }: { state: ExecutionState }) {
  const config: Record<ExecutionState, { color: string; label: string; animate: boolean }> = {
    running: { color: 'bg-blue-500', label: 'Running', animate: true },
    idle: { color: 'bg-gray-500', label: 'Idle', animate: false },
    completed: { color: 'bg-green-500', label: 'Completed', animate: false },
    has_failures: { color: 'bg-red-500', label: 'Has Failures', animate: false },
  };

  const { color, label, animate } = config[state];

  return (
    <span className="flex items-center gap-2">
      <span
        className={`w-2 h-2 rounded-full ${color} ${animate ? 'animate-pulse' : ''}`}
      />
      <span className="text-sm">{label}</span>
    </span>
  );
}

function StatBadge({
  label,
  count,
  color,
}: {
  label: string;
  count: number;
  color: 'gray' | 'blue' | 'green' | 'red';
}) {
  if (count === 0) return null;

  const colorClass: Record<string, string> = {
    gray: 'bg-gray-700',
    blue: 'bg-blue-900',
    green: 'bg-green-900',
    red: 'bg-red-900',
  };

  return (
    <span className={`px-2 py-0.5 rounded ${colorClass[color]}`}>
      {label}: {count}
    </span>
  );
}
