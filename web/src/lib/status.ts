import type { ApplicationStatus } from "@/api/types";

export const STATUS_LABELS: Record<ApplicationStatus, string> = {
  applied: "Applied",
  replied: "Replied",
  interview: "Interview",
  offer: "Offer",
  rejected: "Rejected",
};

export const BOARD_STATUSES: ApplicationStatus[] = [
  "applied",
  "replied",
  "interview",
  "offer",
  "rejected",
];

export const STATUS_COLORS: Record<ApplicationStatus, string> = {
  applied:
    "border-status-applied/40 bg-status-applied/12 text-status-applied dark:bg-status-applied/15",
  replied:
    "border-status-replied/40 bg-status-replied/12 text-status-replied dark:bg-status-replied/15",
  interview:
    "border-status-interview/40 bg-status-interview/12 text-status-interview dark:bg-status-interview/15",
  offer: "border-status-offer/40 bg-status-offer/12 text-status-offer dark:bg-status-offer/15",
  rejected:
    "border-status-rejected/40 bg-status-rejected/12 text-status-rejected dark:bg-status-rejected/15",
};

export const STATUS_DOT: Record<ApplicationStatus, string> = {
  applied: "bg-status-applied",
  replied: "bg-status-replied",
  interview: "bg-status-interview",
  offer: "bg-status-offer",
  rejected: "bg-status-rejected",
};
