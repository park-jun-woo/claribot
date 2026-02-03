# TASK-EXT-010: Project 탭 기본 구조

## 목표
App.tsx에 Project 탭을 추가하고 ProjectPanel.tsx 컴포넌트 생성

## 작업 내용

### 1. App.tsx 수정
- 탭 상태에 `'project'` 추가: `useState<'project' | 'features' | 'tasks'>('project')`
- Project 탭 버튼 추가
- `view === 'project'` 일 때 `<ProjectPanel />` 렌더링

### 2. ProjectPanel.tsx 생성
```typescript
// vscode-extension/webview-ui/src/components/ProjectPanel.tsx
export function ProjectPanel() {
  return (
    <div className="p-4 space-y-4 overflow-auto h-full">
      <ProjectInfo />
      <ContextSection />
      <TechSection />
      <DesignSection />
      <ExecutionStatus />
    </div>
  );
}
```

## 파일
- `vscode-extension/webview-ui/src/App.tsx` (수정)
- `vscode-extension/webview-ui/src/components/ProjectPanel.tsx` (신규)

## 의존성
- 없음 (첫 번째 Task)

## 완료 조건
- Project 탭 클릭 시 ProjectPanel 렌더링
- 기존 Features, Tasks 탭 정상 동작
