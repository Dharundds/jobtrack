import { useState } from "react";
import { Bot, Loader2, RefreshCw, Search, Sun, Moon } from "lucide-react";
import { AuthAccountButton } from "@/components/AuthAccountButton";
import { LLMClassifyDialog } from "@/components/LLMClassifyDialog";
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
  const [llmDialogOpen, setLlmDialogOpen] = useState(false);

  const syncing = Boolean(sync?.running || triggerSync.isPending);
  const llmClassifying = Boolean(sync?.llm_classify_running);
  const busy = syncing || llmClassifying;

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
    <header
      className={cn(
        "sync-scan sticky top-0 z-30 shrink-0 border-b border-border bg-surface-elevated/95 backdrop-blur-sm",
        busy && "sync-scan-active",
      )}
    >
      <LLMClassifyDialog open={llmDialogOpen} onClose={() => setLlmDialogOpen(false)} />
      <div className="flex h-14 items-center gap-3 px-4 md:px-6">
        <div className="relative min-w-0 flex-1 md:max-w-md">
          <Search
            className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-foreground-muted"
            aria-hidden
          />
          <input
            type="search"
            value={searchQuery}
            onChange={(e) => onSearchChange(e.target.value)}
            placeholder="Filter company or role…"
            aria-label="Search applications"
            className={cn(
              "h-9 w-full rounded-md border border-border bg-surface pl-9 pr-3 text-sm",
              "placeholder:text-foreground-muted",
            )}
          />
        </div>

        <div className="flex items-center gap-2">
          <AuthAccountButton />

          <div className="hidden flex-col items-end font-data text-[10px] leading-tight text-foreground-muted lg:flex">
            {llmClassifying ? (
              <span className="text-accent">
                LLM{" "}
                {(sync?.classify_total ?? 0) > 0
                  ? `· ${sync?.classify_done ?? 0}/${sync?.classify_total ?? 0} · model ${sync?.llm_reviewed ?? 0}`
                  : "ACTIVE"}
              </span>
            ) : syncing ? (
              <span className="text-accent">
                SYNC{" "}
                {sync?.sync_phase === "fetch" && sync.processed != null
                  ? `· ${sync.processed} msgs`
                  : sync?.sync_phase === "classify" && (sync.classify_total ?? 0) > 0
                    ? `· ${sync.classify_done ?? 0}/${sync.classify_total} rules`
                    : sync?.sync_phase === "group" && (sync.classify_total ?? 0) > 0
                      ? `· group ${sync.classify_done ?? 0}/${sync.classify_total} · ${sync.applications_linked ?? 0} apps`
                      : "ACTIVE"}
              </span>
            ) : sync?.last_synced_at ? (
              <span>
                UPD {formatRelative(sync.last_synced_at)}
                {sync.classifier ? ` · ${sync.classifier.toUpperCase()}` : ""}
              </span>
            ) : (
              <span>NO SYNC YET</span>
            )}
            {sync?.last_error ? (
              <span className="max-w-[9rem] truncate text-status-rejected">SYNC ERR</span>
            ) : null}
            {sync?.llm_classify_error ? (
              <span className="max-w-[9rem] truncate text-status-rejected">LLM ERR</span>
            ) : null}
          </div>

          <Button
            variant="secondary"
            size="md"
            onClick={() => setLlmDialogOpen(true)}
            disabled={busy}
            aria-label="Classify with LLM"
          >
            <Bot className="h-4 w-4" aria-hidden />
            <span className="hidden sm:inline">LLM</span>
          </Button>

          <Button
            variant="primary"
            size="md"
            onClick={onSync}
            disabled={busy}
            aria-label="Sync now"
          >
            {syncing ? (
              <Loader2 className="h-4 w-4 animate-spin" aria-hidden />
            ) : (
              <RefreshCw className="h-4 w-4" aria-hidden />
            )}
            <span className="hidden sm:inline">Sync</span>
          </Button>

          <button
            type="button"
            onClick={toggle}
            className={cn(
              "inline-flex h-9 w-9 items-center justify-center rounded-md border border-border",
              "bg-surface text-foreground-muted hover:text-foreground",
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
      </div>
    </header>
  );
}
