/** Deep link into Gmail web for a thread or single message. */
export function gmailUrl(opts: { threadId?: string; gmailId?: string }): string | undefined {
  const threadId = opts.threadId?.trim();
  const gmailId = opts.gmailId?.trim();
  if (threadId) {
    return `https://mail.google.com/mail/u/0/#inbox/${encodeURIComponent(threadId)}`;
  }
  if (gmailId) {
    return `https://mail.google.com/mail/u/0/#all/${encodeURIComponent(gmailId)}`;
  }
  return undefined;
}
