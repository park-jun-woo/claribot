# TASK-EXT-013: EditableField 공통 컴포넌트

## 목표
Key-Value 형식의 필드를 편집할 수 있는 공통 컴포넌트 생성

## 작업 내용

### EditableField.tsx 생성
```typescript
// vscode-extension/webview-ui/src/components/EditableField.tsx
interface EditableFieldProps {
  label: string;
  value: string;
  onChange: (value: string) => void;
  required?: boolean;
  editing?: boolean;
  placeholder?: string;
  onDelete?: () => void;  // 삭제 버튼 (사용자 정의 필드용)
}

export function EditableField({
  label,
  value,
  onChange,
  required,
  editing,
  placeholder,
  onDelete,
}: EditableFieldProps) {
  return (
    <div className="flex items-start gap-2 text-sm">
      <span className="w-32 opacity-70 flex-shrink-0">
        {label}
        {required && <span className="text-red-400 ml-0.5">*</span>}:
      </span>
      {editing ? (
        <div className="flex-1 flex items-center gap-1">
          <input
            type="text"
            value={value}
            onChange={(e) => onChange(e.target.value)}
            placeholder={placeholder}
            className="flex-1 px-2 py-1 bg-vscode-input-background border border-vscode-input-border rounded"
          />
          {onDelete && (
            <button
              onClick={onDelete}
              className="p-1 hover:bg-vscode-list-hover rounded text-red-400"
              title="Delete field"
            >
              🗑
            </button>
          )}
        </div>
      ) : (
        <span className={!value ? 'opacity-50' : ''}>{value || '-'}</span>
      )}
    </div>
  );
}
```

### AddFieldButton 컴포넌트
```typescript
interface AddFieldButtonProps {
  onAdd: (key: string) => void;
}

export function AddFieldButton({ onAdd }: AddFieldButtonProps) {
  const [adding, setAdding] = useState(false);
  const [newKey, setNewKey] = useState('');

  const handleAdd = () => {
    if (newKey.trim()) {
      onAdd(newKey.trim());
      setNewKey('');
      setAdding(false);
    }
  };

  if (!adding) {
    return (
      <button
        onClick={() => setAdding(true)}
        className="text-sm opacity-70 hover:opacity-100"
      >
        + Add field
      </button>
    );
  }

  return (
    <div className="flex items-center gap-2">
      <input
        type="text"
        value={newKey}
        onChange={(e) => setNewKey(e.target.value)}
        placeholder="field name"
        className="px-2 py-1 text-sm bg-vscode-input-background border border-vscode-input-border rounded"
        autoFocus
      />
      <button onClick={handleAdd} className="text-sm">Add</button>
      <button onClick={() => setAdding(false)} className="text-sm opacity-70">Cancel</button>
    </div>
  );
}
```

## 파일
- `vscode-extension/webview-ui/src/components/EditableField.tsx` (신규)

## 의존성
- 없음

## 완료 조건
- 읽기/편집 모드 전환 동작
- 필수 필드 표시 (*)
- 사용자 정의 필드 추가/삭제
