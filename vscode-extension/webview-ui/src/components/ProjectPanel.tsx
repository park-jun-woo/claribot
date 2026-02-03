import { ProjectInfo } from './ProjectInfo';
import { ContextSection } from './ContextSection';
import { TechSection } from './TechSection';
import { DesignSection } from './DesignSection';
import { ExecutionStatus } from './ExecutionStatus';

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
