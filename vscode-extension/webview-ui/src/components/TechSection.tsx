import { useState } from 'react';
import { useStore } from '../store';
import { SectionCard } from './SectionCard';
import { EditableField, AddFieldButton } from './EditableField';
import { vscode } from '../vscode';

const REQUIRED_FIELDS = ['backend', 'frontend', 'database'];

export function TechSection() {
  const { tech } = useStore();
  const [editing, setEditing] = useState(false);
  const [draft, setDraft] = useState<Record<string, any>>({});

  const handleEdit = () => {
    setDraft({ ...tech });
    setEditing(true);
  };

  const handleSave = () => {
    vscode.postMessage({ type: 'saveTech', data: draft });
    setEditing(false);
  };

  const handleCancel = () => {
    setEditing(false);
    setDraft({});
  };

  const data = editing ? draft : (tech || {});
  const allKeys = Object.keys(data);
  const customKeys = allKeys.filter(k => !REQUIRED_FIELDS.includes(k));

  return (
    <SectionCard
      title="Tech Stack"
      onEdit={editing ? undefined : handleEdit}
    >
      <div className="space-y-2">
        {REQUIRED_FIELDS.map(key => (
          <EditableField
            key={key}
            label={key}
            value={data[key] || ''}
            onChange={(v) => setDraft({ ...draft, [key]: v })}
            required
            editing={editing}
          />
        ))}
        {customKeys.map(key => (
          <EditableField
            key={key}
            label={key}
            value={data[key] || ''}
            onChange={(v) => setDraft({ ...draft, [key]: v })}
            editing={editing}
            onDelete={editing ? () => {
              const newDraft = { ...draft };
              delete newDraft[key];
              setDraft(newDraft);
            } : undefined}
          />
        ))}
        {editing && (
          <>
            <AddFieldButton onAdd={(key) => setDraft({ ...draft, [key]: '' })} />
            <div className="flex gap-2 pt-3 border-t border-vscode-border mt-3">
              <button
                onClick={handleSave}
                className="px-3 py-1 text-sm bg-vscode-button-bg text-vscode-button-fg rounded"
              >
                Save
              </button>
              <button
                onClick={handleCancel}
                className="px-3 py-1 text-sm hover:bg-vscode-list-hover rounded"
              >
                Cancel
              </button>
            </div>
          </>
        )}
      </div>
    </SectionCard>
  );
}
