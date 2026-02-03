import { useStore } from '../store';
import { SectionCard } from './SectionCard';

export function ProjectInfo() {
  const { project } = useStore();

  if (!project) {
    return (
      <SectionCard title="Project Info">
        <div className="text-sm opacity-50">No project loaded</div>
      </SectionCard>
    );
  }

  return (
    <SectionCard title="Project Info">
      <div className="space-y-2 text-sm">
        <div className="flex">
          <span className="w-24 opacity-70">ID:</span>
          <span className="font-mono">{project.id}</span>
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
  const colors: Record<string, string> = {
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
  if (!dateStr) return '-';
  try {
    return new Date(dateStr).toLocaleDateString();
  } catch {
    return dateStr;
  }
}
