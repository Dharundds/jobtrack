import { Link } from "react-router-dom";
import { useStats } from "@/api/hooks";
import { EmptyState } from "@/components/ui/EmptyState";
import { Skeleton } from "@/components/ui/Skeleton";
import { BOARD_STATUSES, STATUS_COLORS, STATUS_LABELS } from "@/lib/status";
import { cn } from "@/lib/utils";

export function DashboardPage() {
  const { data, isLoading, isError, refetch } = useStats();

  if (isLoading) {
    return (
      <div className="space-y-4">
        <Skeleton className="h-10 w-48" />
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          {[1, 2, 3, 4].map((i) => (
            <Skeleton key={i} className="h-24" />
          ))}
        </div>
      </div>
    );
  }

  if (isError || !data) {
    return (
      <EmptyState
        title="Could not load stats"
        description="Check that the backend is running."
        actionLabel="Retry"
        onAction={() => refetch()}
      />
    );
  }

  const pct = Math.round(data.response_rate * 100);

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-semibold tracking-tight">Dashboard</h1>
        <p className="text-sm text-foreground-muted">Overview of your job search pipeline.</p>
      </div>

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <StatCard label="Total applications" value={String(data.total_applications)} />
        <StatCard label="Response rate" value={`${pct}%`} hint="Applications that reached replied+" />
        <StatCard
          label="Avg days to first reply"
          value={
            data.avg_days_to_first_reply > 0
              ? data.avg_days_to_first_reply.toFixed(1)
              : "—"
          }
        />
        <Link
          to="/board"
          className="rounded-xl border border-border bg-surface-elevated p-4 transition hover:border-accent/40"
        >
          <p className="text-xs text-foreground-muted">Quick link</p>
          <p className="mt-1 text-lg font-semibold">Open board →</p>
        </Link>
      </div>

      <section className="rounded-xl border border-border p-4">
        <h2 className="text-sm font-medium">By status</h2>
        <div className="mt-4 grid gap-3 sm:grid-cols-2 lg:grid-cols-5">
          {BOARD_STATUSES.map((s) => {
            const count = data.by_status[s] ?? 0;
            return (
              <Link
                key={s}
                to={`/applications?status=${s}`}
                className={cn(
                  "rounded-lg border px-3 py-3 transition hover:opacity-90",
                  STATUS_COLORS[s],
                )}
              >
                <p className="text-xs font-medium">{STATUS_LABELS[s]}</p>
                <p className="text-2xl font-semibold">{count}</p>
              </Link>
            );
          })}
        </div>
      </section>
    </div>
  );
}

function StatCard({
  label,
  value,
  hint,
}: {
  label: string;
  value: string;
  hint?: string;
}) {
  return (
    <div className="rounded-xl border border-border bg-surface-elevated p-4">
      <p className="text-xs text-foreground-muted">{label}</p>
      <p className="mt-1 text-2xl font-semibold tabular-nums">{value}</p>
      {hint ? <p className="mt-1 text-[10px] text-foreground-muted">{hint}</p> : null}
    </div>
  );
}
