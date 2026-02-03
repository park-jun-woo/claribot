interface SectionCardProps {
  title: string;
  children: React.ReactNode;
  onEdit?: () => void;
  className?: string;
}

export function SectionCard({ title, children, onEdit, className }: SectionCardProps) {
  return (
    <div className={`border border-vscode-border rounded ${className || ''}`}>
      <div className="flex items-center justify-between px-3 py-2 border-b border-vscode-border bg-vscode-sideBar-background">
        <span className="font-medium text-sm">{title}</span>
        {onEdit && (
          <button
            onClick={onEdit}
            className="text-xs px-2 py-0.5 hover:bg-vscode-list-hover rounded"
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
