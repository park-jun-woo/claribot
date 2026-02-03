# TASK-EXT-011: SectionCard 공통 컴포넌트

## 목표
Project 탭의 각 섹션을 감싸는 공통 카드 컴포넌트 생성

## 작업 내용

### SectionCard.tsx 생성
```typescript
// vscode-extension/webview-ui/src/components/SectionCard.tsx
interface SectionCardProps {
  title: string;
  children: React.ReactNode;
  onEdit?: () => void;  // Edit 버튼 표시 여부
  className?: string;
}

export function SectionCard({ title, children, onEdit, className }: SectionCardProps) {
  return (
    <div className={`border border-vscode-border rounded ${className}`}>
      <div className="flex items-center justify-between px-3 py-2 border-b border-vscode-border bg-vscode-sideBar-background">
        <span className="font-medium">{title}</span>
        {onEdit && (
          <button
            onClick={onEdit}
            className="text-sm px-2 py-0.5 hover:bg-vscode-list-hover rounded"
          >
            Edit
          </button>
        )}
      </div>
      <div className="p-3">
        {children}
      </div>
    </div>
  );
}
```

## 파일
- `vscode-extension/webview-ui/src/components/SectionCard.tsx` (신규)

## 의존성
- 없음

## 완료 조건
- SectionCard 컴포넌트 정상 렌더링
- Edit 버튼 옵션 동작
