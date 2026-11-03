import { Navigate, Outlet, useLocation } from "react-router-dom";
import { Loader2 } from "lucide-react";
import { useAuthStatus } from "@/api/hooks";

export function RequireConnected() {
  const { data, isLoading, isError, refetch } = useAuthStatus();
  const location = useLocation();

  if (isLoading) {
    return (
      <div className="flex min-h-[50vh] items-center justify-center text-foreground-muted">
        <Loader2 className="h-8 w-8 animate-spin" aria-label="Loading" />
      </div>
    );
  }

  if (isError) {
    return (
      <div className="mx-auto max-w-md py-16 text-center">
        <p className="text-sm text-foreground-muted">Could not reach the server.</p>
        <button
          type="button"
          onClick={() => refetch()}
          className="mt-4 rounded-lg bg-accent px-4 py-2 text-sm font-medium text-accent-foreground"
        >
          Retry
        </button>
      </div>
    );
  }

  if (!data?.connected) {
    return <Navigate to="/connect" replace state={{ from: location.pathname }} />;
  }

  return <Outlet />;
}
