import type { ReactNode } from "react";
import { ExternalLink } from "lucide-react";
import { gmailUrl } from "@/lib/gmail";
import { cn } from "@/lib/utils";

type GmailLinkProps = {
  threadId?: string;
  gmailId?: string;
  className?: string;
  /** Accessible name when children are only an icon */
  label?: string;
  children?: ReactNode;
  showIcon?: boolean;
};

export function GmailLink({
  threadId,
  gmailId,
  className,
  label = "Open in Gmail",
  children,
  showIcon,
}: GmailLinkProps) {
  const href = gmailUrl({ threadId, gmailId });
  if (!href) {
    return <>{children}</>;
  }

  return (
    <a
      href={href}
      target="_blank"
      rel="noopener noreferrer"
      className={cn(
        "inline-flex items-center gap-1 text-accent underline-offset-2 hover:underline",
        className,
      )}
      aria-label={children ? undefined : label}
      title={label}
    >
      {children}
      {showIcon ? <ExternalLink className="h-3.5 w-3.5 shrink-0 opacity-80" aria-hidden /> : null}
    </a>
  );
}
