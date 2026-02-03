import { useState } from 'react';

interface EditableFieldProps {
  label: string;
  value: string;
  onChange: (value: string) => void;
  required?: boolean;
  editing?: boolean;
  placeholder?: string;
  onDelete?: () => void;
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
            className="flex-1 px-2 py-1 bg-vscode-input-background border border-vscode-input-border rounded text-sm"
          />
          {onDelete && (
            <button
              onClick={onDelete}
              className="p-1 hover:bg-vscode-list-hover rounded text-red-400"
              title="Delete field"
            >
              ✕
            </button>
          )}
        </div>
      ) : (
        <span className={!value ? 'opacity-50' : ''}>{value || '-'}</span>
      )}
    </div>
  );
}

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

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      handleAdd();
    } else if (e.key === 'Escape') {
      setAdding(false);
      setNewKey('');
    }
  };

  if (!adding) {
    return (
      <button
        onClick={() => setAdding(true)}
        className="text-xs opacity-70 hover:opacity-100 mt-2"
      >
        + Add field
      </button>
    );
  }

  return (
    <div className="flex items-center gap-2 mt-2">
      <input
        type="text"
        value={newKey}
        onChange={(e) => setNewKey(e.target.value)}
        onKeyDown={handleKeyDown}
        placeholder="field name"
        className="px-2 py-1 text-xs bg-vscode-input-background border border-vscode-input-border rounded"
        autoFocus
      />
      <button onClick={handleAdd} className="text-xs px-2 py-1 bg-vscode-button-bg text-vscode-button-fg rounded">
        Add
      </button>
      <button onClick={() => { setAdding(false); setNewKey(''); }} className="text-xs opacity-70">
        Cancel
      </button>
    </div>
  );
}
