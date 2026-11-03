import { Link, useSearchParams } from "react-router-dom";
import { useApplications } from "@/api/hooks";
import type { ApplicationStatus } from "@/api/types";
import { Badge } from "@/components/ui/Badge";
import { EmptyState } from "@/components/ui/EmptyState";
import { Skeleton } from "@/components/ui/Skeleton";
import { Button } from "@/components/ui/Button";
import { useAppOutletContext } from "@/lib/layout-context";
import { formatDate, formatRelative } from "@/lib/dates";
import { STATUS_COLORS, STATUS_LABELS } from "@/lib/status";
import { cn } from "@/lib/utils";

export function ApplicationsListPage() {
  const { searchQuery } = useAppOutletContext();
  const [params, setParams] = useSearchParams();
  const status = params.get("status") ?? "";
  const awaiting = params.get("awaiting_reply") === "1";

  const { data, isLoading, isError, refetch } = useApplications({
    status: status || undefined,
    q: searchQuery || undefined,
    awaiting_reply: awaiting || undefined,
  });

  const apps = data?.applications ?? [];

  const setFilter = (key: string, value: string) => {
    const next = new URLSearchParams(params);
    if (value) next.set(key, value);
    else next.delete(key);
    setParams(next);
  };

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap items-end justify-between gap-4">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">Applications</h1>
          <p className="text-sm text-foreground-muted">Search and filter your pipeline.</p>
        </div>
        <div className="flex flex-wrap gap-2">
          <select
            value={status}
            onChange={(e) => setFilter("status", e.target.value)}
            className="h-9 rounded-lg border border-border bg-surface-elevated px-2 text-sm"
            aria-label="Filter by status"
          >
            <option value="">All statuses</option>
            {(Object.keys(STATUS_LABELS) as ApplicationStatus[]).map((s) => (
              <option key={s} value={s}>
                {STATUS_LABELS[s]}
              </option>
            ))}
          </select>
          <Button
            variant={awaiting ? "primary" : "secondary"}
            size="sm"
            onClick={() => setFilter("awaiting_reply", awaiting ? "" : "1")}
          >
            Awaiting my reply
          </Button>
        </div>
      </div>

      {isLoading ? (
        <Skeleton className="h-64 w-full" />
      ) : isError ? (
        <EmptyState
          title="Could not load applications"
          description="Check that the backend is running and try again."
          actionLabel="Retry"
          onAction={() => refetch()}
        />
      ) : apps.length === 0 ? (
        <EmptyState
          title="No applications yet"
          description="Connect Gmail and run Sync now to import job-related emails."
          actionLabel="Open settings"
          onAction={() => (window.location.href = "/settings")}
        />
      ) : (
        <div className="overflow-x-auto rounded-xl border border-border">
          <table className="w-full min-w-[640px] text-left text-sm">
            <thead className="border-b border-border bg-surface-muted text-foreground-muted">
              <tr>
                <th className="px-4 py-3 font-medium">Company</th>
                <th className="px-4 py-3 font-medium">Role</th>
                <th className="px-4 py-3 font-medium">Status</th>
                <th className="px-4 py-3 font-medium">Last activity</th>
                <th className="px-4 py-3 font-medium">Applied</th>
                <th className="px-4 py-3 font-medium">Reply?</th>
              </tr>
            </thead>
            <tbody>
              {apps.map((app) => (
                <tr key={app.id} className="border-b border-border last:border-0 hover:bg-surface-muted/50">
                  <td className="px-4 py-3">
                    <Link to={`/applications/${app.id}`} className="font-medium hover:text-accent">
                      {app.company || "Unknown company"}
                    </Link>
                  </td>
                  <td className="px-4 py-3 text-foreground-muted">{app.role || "—"}</td>
                  <td className="px-4 py-3">
                    <Badge className={cn(STATUS_COLORS[app.status])}>
                      {STATUS_LABELS[app.status]}
                    </Badge>
                  </td>
                  <td className="px-4 py-3 text-foreground-muted">
                    {formatRelative(app.last_inbound_at ?? app.updated_at)}
                  </td>
                  <td className="px-4 py-3 text-foreground-muted">{formatDate(app.applied_at)}</td>
                  <td className="px-4 py-3">
                    {app.awaiting_my_reply ? (
                      <Badge className="border-amber-500/30 bg-amber-500/10 text-amber-700 dark:text-amber-300">
                        Awaiting
                      </Badge>
                    ) : (
                      "—"
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
