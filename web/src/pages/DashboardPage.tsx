import { Link } from "react-router-dom";
import { Columns3, List } from "lucide-react";
import { useStats } from "@/api/hooks";
import { EmptyState } from "@/components/ui/EmptyState";
import { Skeleton } from "@/components/ui/Skeleton";
import { BOARD_STATUSES, STATUS_DOT, STATUS_LABELS } from "@/lib/status";
import { cn } from "@/lib/utils";
import type { ApplicationStatus } from "@/api/types";

export function DashboardPage() {
  const { data, isLoading, isError, refetch } = useStats();

  if (isLoading) {
    return (
      <div className="space-y-4">
        <Skeleton className="h-14 w-full max-w-xl" />
        <Skeleton className="h-72 w-full" />
      </div>
    );
  }

  if (isError || !data) {
    return (
      <EmptyState
        title="Could not load overview"
        description="Check that the backend is running on port 8080."
        actionLabel="Retry"
        onAction={() => refetch()}
      />
    );
  }

  const total = data.total_applications;
  const pct = Math.round(data.response_rate * 100);
  const avgReply =
    data.avg_days_to_first_reply > 0 ? data.avg_days_to_first_reply.toFixed(1) : "—";
  const activePipeline =
    (data.by_status.replied ?? 0) +
    (data.by_status.interview ?? 0) +
    (data.by_status.offer ?? 0);

  return (
    <div className="mx-auto max-w-3xl space-y-6">
      <header className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <h1 className="page-title">Overview</h1>
          <p className="page-lead">
            Totals from synced mail. Use <strong className="font-semibold text-foreground">Sync</strong>{" "}
            when you expect new messages.
          </p>
        </div>
        <div className="flex flex-wrap gap-2">
          <Link
            to="/board"
            className={cn(
              "inline-flex h-8 items-center justify-center gap-2 rounded-md px-3 text-xs font-semibold",
              "bg-accent text-accent-foreground shadow-board hover:brightness-110",
            )}
          >
            <Columns3 className="h-4 w-4" aria-hidden />
            Board
          </Link>
          <Link
            to="/applications"
            className={cn(
              "inline-flex h-8 items-center justify-center gap-2 rounded-md border border-border px-3 text-xs font-semibold",
              "bg-surface-elevated text-foreground hover:bg-surface-muted",
            )}
          >
            <List className="h-4 w-4" aria-hidden />
            All rows
          </Link>
        </div>
      </header>

      <section className="panel overflow-hidden shadow-board-lg">
        <div className="flex flex-wrap items-center justify-between gap-3 bg-rail px-4 py-3 text-rail-foreground">
          <p className="font-data text-[11px] font-medium uppercase tracking-[0.2em] text-rail-muted">
            Pipeline snapshot
          </p>
          <p className="font-data text-xs tabular-nums text-rail-foreground">
            <span className="text-accent">{activePipeline}</span> past applied
            {total > 0 ? (
              <>
                {" "}
                · <span className="text-rail-muted">{total}</span> total
              </>
            ) : null}
          </p>
        </div>

        <div className="grid divide-y divide-border sm:grid-cols-3 sm:divide-x sm:divide-y-0">
          <MetricCell label="Applications" value={String(total)} />
          <MetricCell label="Response rate" value={`${pct}%`} hint="Reached replied or later" />
          <MetricCell label="Avg days to reply" value={avgReply} />
        </div>

        <div className="border-t border-border bg-surface-muted/40 px-4 py-2">
          <h2 className="font-data text-[11px] font-medium uppercase tracking-widest text-foreground-muted">
            Status mix
          </h2>
        </div>

        {total === 0 ? (
          <p className="px-4 py-8 text-center text-sm text-foreground-muted">
            No applications yet. Connect Gmail and run Sync.
          </p>
        ) : (
          <ul className="divide-y divide-border">
            {BOARD_STATUSES.map((s) => {
              const count = data.by_status[s] ?? 0;
              const share = total > 0 ? Math.round((count / total) * 100) : 0;
              return (
                <StatusRow key={s} status={s} count={count} share={share} />
              );
            })}
          </ul>
        )}
      </section>
    </div>
  );
}

function MetricCell({
  label,
  value,
  hint,
}: {
  label: string;
  value: string;
  hint?: string;
}) {
  return (
    <div className="px-4 py-4">
      <p className="text-sm text-foreground-muted">{label}</p>
      <p className="mt-1 font-data text-3xl font-semibold tabular-nums tracking-tight text-foreground">
        {value}
      </p>
      {hint ? <p className="mt-1 text-[11px] leading-snug text-foreground-muted">{hint}</p> : null}
    </div>
  );
}

function StatusRow({
  status,
  count,
  share,
}: {
  status: ApplicationStatus;
  count: number;
  share: number;
}) {
  return (
    <li>
      <Link
        to={`/applications?status=${status}`}
        className="group block px-4 py-3 transition hover:bg-surface-muted/50"
      >
        <div className="flex items-center justify-between gap-4">
          <span className="flex min-w-0 items-center gap-3">
            <span className={cn("h-2.5 w-2.5 shrink-0 rounded-sm", STATUS_DOT[status])} aria-hidden />
            <span className="text-sm font-semibold text-foreground">{STATUS_LABELS[status]}</span>
          </span>
          <span className="flex shrink-0 items-baseline gap-2 font-data tabular-nums">
            <span className="text-lg font-semibold text-foreground">{count}</span>
            <span className="text-xs text-foreground-muted">{share}%</span>
          </span>
        </div>
        <div
          className="mt-2 h-1.5 overflow-hidden rounded-full bg-border/80"
          role="presentation"
        >
          <div
            className={cn("h-full rounded-full transition-[width]", STATUS_DOT[status])}
            style={{ width: `${Math.max(share, count > 0 ? 4 : 0)}%` }}
          />
        </div>
      </Link>
    </li>
  );
}
