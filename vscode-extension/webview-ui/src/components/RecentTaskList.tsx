import type { Task } from '../types';

interface RecentTaskListProps {
  tasks: Task[];
  limit?: number;
}

export function RecentTaskList({ tasks, limit = 5 }: RecentTaskListProps) {
  const recentTasks = getRecentTasks(tasks, limit);

  if (recentTasks.length === 0) {
    return <div className="text-sm opacity-50">No tasks yet</div>;
  }

  return (
    <div className="space-y-1">
      <div className="text-xs opacity-70 mb-2">Recent Tasks:</div>
      {recentTasks.map(task => (
        <TaskLogItem key={task.id} task={task} />
      ))}
    </div>
  );
}

function TaskLogItem({ task }: { task: Task }) {
  const icon = getStatusIcon(task.status);
  const colorClass = getStatusColor(task.status);
  const summary = getTaskSummary(task);

  return (
    <div className={`flex items-center gap-2 text-sm ${colorClass}`}>
      <span className={`w-4 ${task.status === 'doing' ? 'animate-pulse' : ''}`}>{icon}</span>
      <span className="opacity-50">#{task.id}</span>
      <span className="flex-1 truncate">{task.title}</span>
      {summary && (
        <span className="text-xs opacity-70 truncate max-w-[150px]">
          {summary}
        </span>
      )}
    </div>
  );
}

function getStatusIcon(status: Task['status']): string {
  switch (status) {
    case 'done': return '✓';
    case 'doing': return '●';
    case 'failed': return '✗';
    case 'pending': return '○';
  }
}

function getStatusColor(status: Task['status']): string {
  switch (status) {
    case 'done': return 'text-green-400';
    case 'doing': return 'text-blue-400';
    case 'failed': return 'text-red-400';
    case 'pending': return 'text-gray-400';
  }
}

function getTaskSummary(task: Task): string {
  if (task.status === 'done' && task.result) {
    return `"${task.result.slice(0, 30)}${task.result.length > 30 ? '...' : ''}"`;
  }
  if (task.status === 'failed' && task.error) {
    return `"${task.error.slice(0, 30)}${task.error.length > 30 ? '...' : ''}"`;
  }
  if (task.status === 'pending') {
    return '(pending)';
  }
  return '';
}

function getRecentTasks(tasks: Task[], limit: number): Task[] {
  const doing = tasks.filter(t => t.status === 'doing');
  const completed = tasks
    .filter(t => t.status === 'done' || t.status === 'failed')
    .sort((a, b) => {
      const aTime = a.completed_at || a.failed_at || '';
      const bTime = b.completed_at || b.failed_at || '';
      return bTime.localeCompare(aTime);
    });
  const pending = tasks.filter(t => t.status === 'pending');

  return [...doing, ...completed, ...pending].slice(0, limit);
}
