# TASK-EXT-017: ProgressBar 컴포넌트

## 목표
Task 진행률을 시각적으로 표시하는 Progress Bar 컴포넌트 생성

## 작업 내용

### ProgressBar.tsx 생성
```typescript
// vscode-extension/webview-ui/src/components/ProgressBar.tsx
interface ProgressBarProps {
  progress: number;  // 0-100
  total: number;
  done: number;
  className?: string;
}

export function ProgressBar({ progress, total, done, className }: ProgressBarProps) {
  return (
    <div className={`space-y-1 ${className}`}>
      <div className="flex items-center gap-2">
        <div className="flex-1 h-2 bg-vscode-input-background rounded overflow-hidden">
          <div
            className="h-full bg-vscode-button-bg transition-all duration-300"
            style={{ width: `${progress}%` }}
          />
        </div>
        <span className="text-sm tabular-nums">
          {progress.toFixed(0)}%
        </span>
      </div>
      <div className="text-xs opacity-70">
        {done} / {total} tasks completed
      </div>
    </div>
  );
}
```

## 파일
- `vscode-extension/webview-ui/src/components/ProgressBar.tsx` (신규)

## 의존성
- 없음

## 완료 조건
- 진행률 바 시각화
- 퍼센트 및 완료/전체 수 표시
- 애니메이션 전환 효과
