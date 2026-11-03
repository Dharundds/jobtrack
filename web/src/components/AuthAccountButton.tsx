import { Link } from "react-router-dom";
import { Loader2, LogOut, User } from "lucide-react";
import { useAuthStatus, useLogout } from "@/api/hooks";
import { cn } from "@/lib/utils";

export function AuthAccountButton() {
  const { data, isLoading, isError, refetch } = useAuthStatus();
  const logout = useLogout();

  if (isLoading) {
    return (
      <div
        className="inline-flex h-9 items-center gap-2 rounded-lg border border-border bg-surface-elevated px-3 text-sm text-foreground-muted"
        aria-busy="true"
      >
        <Loader2 className="h-4 w-4 animate-spin" aria-hidden />
        <span className="hidden sm:inline">Checking…</span>
      </div>
    );
  }

  if (isError) {
    return (
      <button
        type="button"
        onClick={() => refetch()}
        className="h-9 rounded-lg border border-border px-3 text-sm text-foreground-muted hover:text-foreground"
      >
        Retry account
      </button>
    );
  }

  if (!data?.connected) {
    return (
      <Link
        to="/connect"
        className={cn(
          "inline-flex h-9 items-center gap-2 rounded-lg border border-accent/30 bg-accent/10 px-3",
          "text-sm font-medium text-accent hover:bg-accent/15",
        )}
      >
        <User className="h-4 w-4" aria-hidden />
        Connect Gmail
      </Link>
    );
  }

  const email = data.email ?? "Connected";

  return (
    <div className="flex items-center gap-1">
      <div
        className={cn(
          "inline-flex h-9 max-w-[12rem] items-center gap-2 rounded-lg border border-border",
          "bg-surface-elevated px-3 text-sm font-medium text-foreground",
        )}
        title={email}
      >
        <User className="h-4 w-4 shrink-0 text-accent" aria-hidden />
        <span className="truncate">{email}</span>
      </div>
      <button
        type="button"
        disabled={logout.isPending}
        onClick={() => logout.mutate()}
        className={cn(
          "inline-flex h-9 w-9 items-center justify-center rounded-lg border border-border",
          "bg-surface-elevated text-foreground-muted hover:text-foreground",
        )}
        aria-label="Disconnect Gmail"
        title="Disconnect"
      >
        {logout.isPending ? (
          <Loader2 className="h-4 w-4 animate-spin" aria-hidden />
        ) : (
          <LogOut className="h-4 w-4" aria-hidden />
        )}
      </button>
    </div>
  );
}
