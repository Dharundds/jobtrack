import { useMemo, useState } from "react";
import { Link } from "react-router-dom";
import {
  DndContext,
  DragOverlay,
  PointerSensor,
  useDraggable,
  useDroppable,
  useSensor,
  useSensors,
  type DragEndEvent,
  type DragStartEvent,
} from "@dnd-kit/core";
import { CSS } from "@dnd-kit/utilities";
import type { Application, ApplicationStatus } from "@/api/types";
import { useApplications, usePatchApplication } from "@/api/hooks";
import { EmptyState } from "@/components/ui/EmptyState";
import { Skeleton } from "@/components/ui/Skeleton";
import { Badge } from "@/components/ui/Badge";
import { useToast } from "@/components/ui/Toast";
import { useAppOutletContext } from "@/lib/layout-context";
import { formatRelative } from "@/lib/dates";
import { BOARD_STATUSES, STATUS_COLORS, STATUS_LABELS } from "@/lib/status";
import { cn } from "@/lib/utils";

function matchesSearch(app: Application, q: string) {
  if (!q.trim()) return true;
  const needle = q.toLowerCase();
  return (
    app.company.toLowerCase().includes(needle) ||
    app.role.toLowerCase().includes(needle) ||
    app.notes.toLowerCase().includes(needle)
  );
}

function KanbanCard({ app, dragging }: { app: Application; dragging?: boolean }) {
  const { attributes, listeners, setNodeRef, transform, isDragging } = useDraggable({
    id: app.id,
    data: { app },
  });

  const style = transform
    ? { transform: CSS.Translate.toString(transform) }
    : undefined;

  return (
    <div
      ref={setNodeRef}
      style={style}
      className={cn(
        "rounded-lg border border-border bg-surface-elevated p-3 shadow-sm",
        (isDragging || dragging) && "opacity-60 ring-2 ring-accent/40",
      )}
      {...listeners}
      {...attributes}
    >
      <Link
        to={`/applications/${app.id}`}
        className="block font-medium hover:text-accent"
        onClick={(e) => e.stopPropagation()}
      >
        {app.company || "Unknown company"}
      </Link>
      <p className="mt-0.5 line-clamp-2 text-xs text-foreground-muted">{app.role || "—"}</p>
      <div className="mt-2 flex flex-wrap items-center gap-1.5 text-[10px] text-foreground-muted">
        <span>{formatRelative(app.last_inbound_at ?? app.updated_at)}</span>
        {app.awaiting_my_reply ? (
          <Badge className="border-amber-500/30 bg-amber-500/10 text-amber-700 dark:text-amber-300">
            Reply?
          </Badge>
        ) : null}
      </div>
    </div>
  );
}

function KanbanColumn({
  status,
  apps,
  activeId,
}: {
  status: ApplicationStatus;
  apps: Application[];
  activeId: number | null;
}) {
  const { setNodeRef, isOver } = useDroppable({ id: status });

  return (
    <div className="flex w-72 shrink-0 flex-col rounded-xl border border-border bg-surface-muted/40">
      <div className="flex items-center justify-between border-b border-border px-3 py-2">
        <span className="text-sm font-medium">{STATUS_LABELS[status]}</span>
        <Badge className={cn("text-[10px]", STATUS_COLORS[status])}>{apps.length}</Badge>
      </div>
      <div
        ref={setNodeRef}
        className={cn(
          "flex min-h-[12rem] flex-1 flex-col gap-2 p-2 transition-colors",
          isOver && "bg-accent/5",
        )}
      >
        {apps.map((app) => (
          <KanbanCard key={app.id} app={app} dragging={activeId === app.id} />
        ))}
      </div>
    </div>
  );
}

export function BoardPage() {
  const { searchQuery } = useAppOutletContext();
  const { data, isLoading, isError, refetch } = useApplications({});
  const patch = usePatchApplication();
  const { toast } = useToast();
  const [activeId, setActiveId] = useState<number | null>(null);

  const sensors = useSensors(useSensor(PointerSensor, { activationConstraint: { distance: 6 } }));

  const byColumn = useMemo(() => {
    const apps = (data?.applications ?? []).filter((a) => matchesSearch(a, searchQuery));
    const map = Object.fromEntries(BOARD_STATUSES.map((s) => [s, [] as Application[]])) as Record<
      ApplicationStatus,
      Application[]
    >;
    for (const app of apps) {
      const col = BOARD_STATUSES.includes(app.status) ? app.status : "applied";
      map[col].push(app);
    }
    return map;
  }, [data?.applications, searchQuery]);

  const activeApp = useMemo(() => {
    if (activeId == null) return null;
    return data?.applications?.find((a) => a.id === activeId) ?? null;
  }, [activeId, data?.applications]);

  const onDragStart = (e: DragStartEvent) => {
    setActiveId(Number(e.active.id));
  };

  const onDragEnd = (e: DragEndEvent) => {
    setActiveId(null);
    const appId = Number(e.active.id);
    const overId = e.over?.id;
    if (overId == null) return;

    let newStatus: ApplicationStatus | null = null;
    if (typeof overId === "string" && BOARD_STATUSES.includes(overId as ApplicationStatus)) {
      newStatus = overId as ApplicationStatus;
    } else if (typeof overId === "number") {
      const overApp = data?.applications?.find((a) => a.id === overId);
      if (overApp && BOARD_STATUSES.includes(overApp.status)) newStatus = overApp.status;
    }
    if (!newStatus) return;
    const app = data?.applications?.find((a) => a.id === appId);
    if (!app || app.status === newStatus) return;

    patch.mutate(
      { id: appId, body: { status: newStatus } },
      {
        onError: (err) => toast(err.message, "error"),
        onSuccess: () => toast(`Moved to ${STATUS_LABELS[newStatus]}`),
      },
    );
  };

  if (isLoading) return <Skeleton className="h-96 w-full" />;

  if (isError) {
    return (
      <EmptyState
        title="Could not load board"
        description="Check that the backend is running."
        actionLabel="Retry"
        onAction={() => refetch()}
      />
    );
  }

  const total = Object.values(byColumn).reduce((n, col) => n + col.length, 0);

  return (
    <div className="space-y-4">
      <div>
        <h1 className="text-2xl font-semibold tracking-tight">Board</h1>
        <p className="text-sm text-foreground-muted">
          Drag cards between columns to update status. {total} application{total === 1 ? "" : "s"}.
        </p>
      </div>

      {total === 0 ? (
        <EmptyState
          title="Nothing on the board"
          description="Run a Gmail sync or clear your search filter."
          actionLabel="View all applications"
          onAction={() => (window.location.href = "/applications")}
        />
      ) : (
        <DndContext sensors={sensors} onDragStart={onDragStart} onDragEnd={onDragEnd}>
          <div className="flex gap-3 overflow-x-auto pb-4">
            {BOARD_STATUSES.map((status) => (
              <KanbanColumn
                key={status}
                status={status}
                apps={byColumn[status]}
                activeId={activeId}
              />
            ))}
          </div>
          <DragOverlay>
            {activeApp ? (
              <div className="w-72 rounded-lg border border-accent/40 bg-surface-elevated p-3 shadow-lg">
                <p className="font-medium">{activeApp.company}</p>
                <p className="text-xs text-foreground-muted">{activeApp.role}</p>
              </div>
            ) : null}
          </DragOverlay>
        </DndContext>
      )}
    </div>
  );
}
