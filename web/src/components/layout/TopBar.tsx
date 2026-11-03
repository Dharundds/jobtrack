import { Loader2, RefreshCw, Search, Sun, Moon } from "lucide-react";
import { AuthAccountButton } from "@/components/AuthAccountButton";
import { Button } from "@/components/ui/Button";
import { useSyncStatus, useTriggerSync } from "@/api/hooks";
import { useToast } from "@/components/ui/Toast";
import { formatRelative } from "@/lib/dates";
import { useTheme } from "@/lib/theme";
import { cn } from "@/lib/utils";

type TopBarProps = {
  searchQuery: string;
  onSearchChange: (value: string) => void;
};

export function TopBar({ searchQuery, onSearchChange }: TopBarProps) {
  const { resolved, toggle } = useTheme();
  const { data: sync } = useSyncStatus();
  const triggerSync = useTriggerSync();
  const { toast } = useToast();

  const onSync = () => {
    triggerSync.mutate(undefined, {
      onSuccess: (res) => {
        if (res.message?.toLowerCase().includes("already running")) {
          toast("Sync already running");
        } else {
          toast("Sync started");
        }
      },
      onError: (e) => toast(e.message, "error"),
    });
  };

  return (
    <header className="sticky top-0 z-30 flex h-14 shrink-0 items-center gap-3 border-b border-border bg-surface/80 px-4 backdrop-blur md:px-6">
      <div className="relative min-w-0 flex-1 md:max-w-md">
        <Search
          className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-foreground-muted"
          aria-hidden
        />
        <input
          type="search"
          value={searchQuery}
          onChange={(e) => onSearchChange(e.target.value)}
          placeholder="Search applications…"
          aria-label="Search applications"
          className={cn(
            "h-9 w-full rounded-lg border border-border bg-surface-elevated pl-9 pr-3 text-sm",
            "placeholder:text-foreground-muted",
          )}
        />
      </div>

      <div className="flex items-center gap-2">
        <AuthAccountButton />

        <div className="hidden flex-col items-end text-[10px] leading-tight text-foreground-muted lg:flex">
          {sync?.running ? (
            <span className="text-accent">Syncing…</span>
          ) : sync?.last_synced_at ? (
            <span>Last synced {formatRelative(sync.last_synced_at)}</span>
          ) : (
            <span>Not synced yet</span>
          )}
          {sync?.last_error ? (
            <span className="max-w-[8rem] truncate text-status-rejected">Error</span>
          ) : null}
        </div>

        <Button
          variant="secondary"
          size="sm"
          onClick={onSync}
          disabled={triggerSync.isPending || sync?.running}
          aria-label="Sync now"
        >
          {sync?.running || triggerSync.isPending ? (
            <Loader2 className="h-4 w-4 animate-spin" aria-hidden />
          ) : (
            <RefreshCw className="h-4 w-4" aria-hidden />
          )}
          <span className="hidden sm:inline">Sync now</span>
        </Button>

        <button
          type="button"
          onClick={toggle}
          className={cn(
            "inline-flex h-9 w-9 items-center justify-center rounded-lg border border-border",
            "bg-surface-elevated text-foreground-muted hover:text-foreground",
          )}
          aria-label={resolved === "dark" ? "Switch to light mode" : "Switch to dark mode"}
        >
          {resolved === "dark" ? (
            <Sun className="h-4 w-4" aria-hidden />
          ) : (
            <Moon className="h-4 w-4" aria-hidden />
          )}
        </button>
      </div>
    </header>
  );
}
