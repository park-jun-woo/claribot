interface ProgressBarProps {
  progress: number;
  total: number;
  done: number;
  className?: string;
}

export function ProgressBar({ progress, total, done, className }: ProgressBarProps) {
  return (
    <div className={`space-y-1 ${className || ''}`}>
      <div className="flex items-center gap-2">
        <div className="flex-1 h-2 bg-vscode-input-background rounded overflow-hidden">
          <div
            className="h-full bg-vscode-button-bg transition-all duration-300"
            style={{ width: `${Math.min(100, Math.max(0, progress))}%` }}
          />
        </div>
        <span className="text-sm tabular-nums w-12 text-right">
          {progress.toFixed(0)}%
        </span>
      </div>
      <div className="text-xs opacity-70">
        {done} / {total} tasks completed
      </div>
    </div>
  );
}
