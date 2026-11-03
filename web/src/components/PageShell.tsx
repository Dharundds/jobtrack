import { cn } from "@/lib/utils";

type PageShellProps = {
  title: string;
  description?: string;
  children?: React.ReactNode;
};

export function PageShell({ title, description, children }: PageShellProps) {
  return (
    <div className="space-y-2">
      <div>
        <h1 className="text-2xl font-semibold tracking-tight text-foreground">
          {title}
        </h1>
        {description ? (
          <p className="mt-1 text-sm text-foreground-muted">{description}</p>
        ) : null}
      </div>
      {children ?? (
        <div
          className={cn(
            "mt-6 rounded-xl border border-border bg-surface-elevated p-8",
            "text-sm text-foreground-muted",
          )}
        >
          This view will be built in a later step.
        </div>
      )}
    </div>
  );
}
