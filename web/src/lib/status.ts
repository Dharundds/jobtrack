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
  applied: "bg-status-applied/15 text-status-applied border-status-applied/30",
  replied: "bg-status-replied/15 text-status-replied border-status-replied/30",
  interview: "bg-status-interview/15 text-status-interview border-status-interview/30",
  offer: "bg-status-offer/15 text-status-offer border-status-offer/30",
  rejected: "bg-status-rejected/15 text-status-rejected border-status-rejected/30",
};
