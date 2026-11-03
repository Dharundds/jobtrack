import { Navigate, useSearchParams } from "react-router-dom";
import { Loader2, Mail, Shield } from "lucide-react";
import { useAuthStatus } from "@/api/hooks";
import { cn } from "@/lib/utils";

export function ConnectPage() {
  const { data, isLoading } = useAuthStatus();
  const [params] = useSearchParams();
  const denied = params.get("auth") === "denied";

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-foreground-muted" aria-label="Loading" />
      </div>
    );
  }

  if (data?.connected) {
    return <Navigate to="/board" replace />;
  }

  return (
    <div className="mx-auto flex min-h-[70vh] max-w-lg flex-col justify-center px-4 py-12">
      <div className="rounded-xl border border-border bg-surface-elevated p-8 shadow-sm">
        {denied ? (
          <p className="mb-4 rounded-lg border border-status-rejected/30 bg-status-rejected/10 px-3 py-2 text-sm text-foreground">
            Google sign-in was cancelled. Try again when you&apos;re ready.
          </p>
        ) : null}
        <div className="mb-6 flex h-12 w-12 items-center justify-center rounded-xl bg-accent/10 text-accent">
          <Mail className="h-6 w-6" aria-hidden />
        </div>
        <h1 className="text-2xl font-semibold tracking-tight">Connect Gmail</h1>
        <p className="mt-3 text-sm leading-relaxed text-foreground-muted">
          Jobtrack reads your inbox with read-only access to find job-related emails.
          Your data stays on this machine — no telemetry or third parties besides Google.
        </p>
        <ul className="mt-4 space-y-2 text-sm text-foreground-muted">
          <li className="flex gap-2">
            <Shield className="mt-0.5 h-4 w-4 shrink-0 text-accent" aria-hidden />
            <span>Scope: Gmail read-only. We never send or delete mail.</span>
          </li>
        </ul>
        <a
          href="/auth/login"
          className={cn(
            "mt-8 inline-flex h-11 w-full items-center justify-center rounded-lg bg-accent px-4",
            "text-sm font-semibold text-accent-foreground hover:opacity-95",
          )}
        >
          Connect Gmail
        </a>
        <p className="mt-4 text-center text-xs text-foreground-muted">
          Use this app at{" "}
          <span className="font-medium text-foreground">http://localhost:5173</span> during
          development so the UI talks to the API.
        </p>
      </div>
    </div>
  );
}
