# TASK-EXT-012: ProjectInfo 컴포넌트

## 목표
프로젝트 기본 정보를 읽기 전용으로 표시하는 컴포넌트 생성

## 작업 내용

### ProjectInfo.tsx 생성
```typescript
// vscode-extension/webview-ui/src/components/ProjectInfo.tsx
import { useStore } from '../store';
import { SectionCard } from './SectionCard';

export function ProjectInfo() {
  const { project } = useStore();

  if (!project) return null;

  return (
    <SectionCard title="Project Info">
      <div className="space-y-2 text-sm">
        <div className="flex">
          <span className="w-24 opacity-70">ID:</span>
          <span>{project.id}</span>
        </div>
        <div className="flex">
          <span className="w-24 opacity-70">Name:</span>
          <span>{project.name}</span>
        </div>
        <div className="flex items-center">
          <span className="w-24 opacity-70">Status:</span>
          <StatusIndicator status={project.status} />
        </div>
        <div className="flex">
          <span className="w-24 opacity-70">Created:</span>
          <span>{formatDate(project.created_at)}</span>
        </div>
      </div>
    </SectionCard>
  );
}

function StatusIndicator({ status }: { status: string }) {
  const colors = {
    active: 'bg-green-500',
    paused: 'bg-yellow-500',
    completed: 'bg-blue-500',
  };
  return (
    <span className="flex items-center gap-2">
      <span className={`w-2 h-2 rounded-full ${colors[status] || 'bg-gray-500'}`} />
      {status}
    </span>
  );
}

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString();
}
```

## 파일
- `vscode-extension/webview-ui/src/components/ProjectInfo.tsx` (신규)

## 의존성
- TASK-EXT-011 (SectionCard)

## 완료 조건
- 프로젝트 ID, Name, Status, Created 표시
- Status별 색상 점 표시
