/** Best-effort company/role from email subject for split forms. */
export function guessCompanyFromMessage(subject: string, body: string): string {
  const text = `${subject}\n${body}`;
  const patterns = [
    /application was sent to\s+([^,\n]+)/i,
    /thank you for applying to\s+([^,\n!.]+)/i,
    /thanks for applying to\s+([^,\n!.]+)/i,
    /application (?:to|for)\s+([^,\n]+)/i,
    /(?:^|:\s)([A-Za-z0-9][A-Za-z0-9 &.'-]{1,50})\s+is hiring/i,
  ];
  for (const re of patterns) {
    const m = text.match(re);
    if (m?.[1]) {
      return m[1].trim().replace(/\s+/g, " ");
    }
  }
  return "";
}

export function guessRoleFromMessage(subject: string, body: string): string {
  const text = `${subject}\n${body}`;
  const m =
    text.match(/(?:for the|for)\s+(.+?)\s+(?:position|role)\b/i) ||
    text.match(/application:\s*(.+)/i);
  return m?.[1]?.trim() ?? "";
}
