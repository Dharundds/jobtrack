import { cn } from "@/lib/utils";

type ButtonProps = React.ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: "primary" | "secondary" | "ghost" | "danger";
  size?: "sm" | "md";
};

export function Button({
  className,
  variant = "secondary",
  size = "md",
  ...props
}: ButtonProps) {
  return (
    <button
      type="button"
      className={cn(
        "inline-flex items-center justify-center gap-2 rounded-md font-semibold transition-colors",
        "focus-visible:outline-none disabled:cursor-not-allowed disabled:opacity-45",
        size === "sm" ? "h-8 px-3 text-xs" : "h-9 px-4 text-sm",
        variant === "primary" &&
          "bg-accent text-accent-foreground shadow-board hover:brightness-110",
        variant === "secondary" &&
          "border border-border bg-surface-elevated text-foreground hover:bg-surface-muted",
        variant === "ghost" && "text-foreground-muted hover:bg-surface-muted hover:text-foreground",
        variant === "danger" &&
          "border border-status-rejected bg-status-rejected text-white shadow-board hover:bg-status-rejected/90 focus-visible:ring-status-rejected",
        className,
      )}
      {...props}
    />
  );
}
