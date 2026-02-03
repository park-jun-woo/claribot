# TASK-EXT-018: RecentTaskList 컴포넌트

## 목표
최근 Task 로그를 표시하는 컴포넌트 생성 (완료/실행중/대기중 상태 표시)

## 작업 내용

### types.ts 확장
```typescript
// 기존 Task 인터페이스 활용 (result 필드 포함)
```

### RecentTaskList.tsx 생성
```typescript
// vscode-extension/webview-ui/src/components/RecentTaskList.tsx
import { Task } from '../types';

interface RecentTaskListProps {
  tasks: Task[];
  limit?: number;  // 표시할 최대 개수 (기본 5)
}

export function RecentTaskList({ tasks, limit = 5 }: RecentTaskListProps) {
  // 최근 완료/실행중 Task 우선 + pending 일부
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
      <span className="w-4">{icon}</span>
      <span className="opacity-50">#{task.id}</span>
      <span className="flex-1 truncate">{task.title}</span>
      {summary && (
        <span className="text-xs opacity-70 truncate max-w-[150px]">
          "{summary}"
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
    return task.result.slice(0, 30);
  }
  if (task.status === 'failed' && task.error) {
    return task.error.slice(0, 30);
  }
  if (task.status === 'pending') {
    return '(pending)';
  }
  return '';
}

function getRecentTasks(tasks: Task[], limit: number): Task[] {
  // doing 먼저, 그 다음 최근 완료/실패, 마지막으로 pending
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
```

## 파일
- `vscode-extension/webview-ui/src/components/RecentTaskList.tsx` (신규)

## 의존성
- 없음

## 완료 조건
- 최근 Task 목록 표시 (최대 5개)
- 상태별 아이콘/색상 구분
- 완료 시 result 한줄평 표시
- 실패 시 error 표시
