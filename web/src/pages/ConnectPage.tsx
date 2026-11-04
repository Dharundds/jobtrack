import { Navigate, useSearchParams } from "react-router-dom";
import { Loader2, Mail, Shield, Radio } from "lucide-react";
import { useAuthStatus } from "@/api/hooks";
import { cn } from "@/lib/utils";

export function ConnectPage() {
  const { data, isLoading } = useAuthStatus();
  const [params] = useSearchParams();
  const denied = params.get("auth") === "denied";

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-rail">
        <Loader2
          className="h-8 w-8 animate-spin text-accent"
          aria-label="Loading"
        />
      </div>
    );
  }

  if (data?.connected) {
    return <Navigate to="/board" replace />;
  }

  return (
    <div className="flex min-h-screen flex-col bg-rail text-rail-foreground">
      <div className="border-b border-white/10 px-6 py-4">
        <div className="mx-auto flex max-w-lg items-center gap-2">
          <Radio className="h-5 w-5 text-accent" aria-hidden />
          <span className="text-sm font-semibold tracking-wide">jobtrack</span>
        </div>
      </div>
      <div className="mx-auto flex w-full max-w-lg flex-1 flex-col justify-center px-4 py-12">
        <div className="panel border-white/10 bg-surface-elevated text-foreground shadow-board-lg">
          {denied ? (
            <p className="mb-4 rounded-md border border-status-rejected/40 bg-status-rejected/10 px-3 py-2 text-sm">
              Google sign-in was cancelled. Try again when you&apos;re ready.
            </p>
          ) : null}
          <div className="mb-6 flex h-12 w-12 items-center justify-center rounded-md bg-accent/15 text-accent">
            <Mail className="h-6 w-6" aria-hidden />
          </div>
          <h1 className="page-title">Connect Gmail</h1>
          <p className="page-lead">
            Read-only sync finds job-related mail on this machine. Nothing leaves your
            computer except OAuth to Google.
          </p>
          <ul className="mt-6 space-y-3 text-sm text-foreground-muted">
            <li className="flex gap-3">
              <Shield className="mt-0.5 h-4 w-4 shrink-0 text-accent" aria-hidden />
              <span>Scope: Gmail read-only. We never send or delete mail.</span>
            </li>
          </ul>
          <a
            href="/auth/login"
            className={cn(
              "mt-8 inline-flex h-11 w-full items-center justify-center rounded-md bg-accent px-4",
              "text-sm font-semibold text-accent-foreground shadow-board hover:brightness-110",
            )}
          >
            Connect Gmail
          </a>
          <p className="mt-4 text-center font-data text-[11px] text-foreground-muted">
            Dev UI: http://localhost:5173 → API :8080
          </p>
        </div>
      </div>
    </div>
  );
}
