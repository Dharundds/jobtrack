import { formatDistanceToNow, parseISO } from "date-fns";

export function formatRelative(iso?: string) {
  if (!iso) return "—";
  try {
    return formatDistanceToNow(parseISO(iso), { addSuffix: true });
  } catch {
    return "—";
  }
}

export function formatDate(iso?: string) {
  if (!iso) return "—";
  try {
    return parseISO(iso).toLocaleDateString(undefined, {
      month: "short",
      day: "numeric",
      year: "numeric",
    });
  } catch {
    return "—";
  }
}
