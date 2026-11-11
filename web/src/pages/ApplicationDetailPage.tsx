import { useEffect, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import {
  useApplication,
  useApplications,
  useDeleteApplication,
  useIgnoreMessage,
  useMergeApplications,
  usePatchApplication,
  useSplitMessage,
} from "@/api/hooks";
import type { ApplicationStatus, Message } from "@/api/types";
import { Badge } from "@/components/ui/Badge";
import { Button } from "@/components/ui/Button";
import { EmptyState } from "@/components/ui/EmptyState";
import { Skeleton } from "@/components/ui/Skeleton";
import { useToast } from "@/components/ui/Toast";
import { formatDate, formatRelative } from "@/lib/dates";
import { GmailLink } from "@/components/GmailLink";
import { guessCompanyFromMessage, guessRoleFromMessage } from "@/lib/message-guess";
import { gmailUrl } from "@/lib/gmail";
import { STATUS_COLORS, STATUS_LABELS } from "@/lib/status";
import { cn } from "@/lib/utils";

export function ApplicationDetailPage() {
  const navigate = useNavigate();
  const { id: idParam } = useParams<{ id: string }>();
  const id = Number(idParam);
  const { data, isLoading, isError, refetch } = useApplication(id);
  const { data: allApps } = useApplications({});
  const patch = usePatchApplication();
  const merge = useMergeApplications();
  const deleteApp = useDeleteApplication();
  const splitMsg = useSplitMessage(id);
  const ignoreMsg = useIgnoreMessage(id);
  const { toast } = useToast();

  const app = data?.application;
  const [company, setCompany] = useState("");
  const [role, setRole] = useState("");
  const [notes, setNotes] = useState("");
  const [status, setStatus] = useState<ApplicationStatus>("applied");
  const [mergeSource, setMergeSource] = useState<number | "">("");
  const [splitDraft, setSplitDraft] = useState<{
    messageId: number;
    company: string;
    role: string;
  } | null>(null);

  useEffect(() => {
    if (!app) return;
    setCompany(app.company);
    setRole(app.role);
    setNotes(app.notes);
    setStatus(app.status);
  }, [app]);

  const save = () => {
    if (!app) return;
    patch.mutate(
      {
        id: app.id,
        body: {
          company: company.trim() || undefined,
          role: role.trim() || undefined,
          notes,
          status: status !== app.status ? status : undefined,
        },
      },
      {
        onSuccess: () => toast("Saved"),
        onError: (e) => toast(e.message, "error"),
      },
    );
  };

  const onMerge = () => {
    if (!app || mergeSource === "") return;
    merge.mutate(
      { target_id: app.id, source_ids: [Number(mergeSource)] },
      {
        onSuccess: () => {
          toast("Applications merged");
          setMergeSource("");
        },
        onError: (e) => toast(e.message, "error"),
      },
    );
  };

  const onDeleteListing = () => {
    if (!app) return;
    if (
      !window.confirm(
        "Delete this entire listing? Messages will be hidden and won’t be grouped again on sync.",
      )
    ) {
      return;
    }
    deleteApp.mutate(app.id, {
      onSuccess: () => {
        toast("Listing deleted");
        navigate("/applications");
      },
      onError: (e) => toast(e.message, "error"),
    });
  };

  const openSplit = (m: Message) => {
    setSplitDraft({
      messageId: m.id,
      company:
        guessCompanyFromMessage(m.subject, m.body_excerpt || m.snippet) || company,
      role: guessRoleFromMessage(m.subject, m.body_excerpt || m.snippet) || role,
    });
  };

  const confirmSplit = () => {
    if (!splitDraft || !splitDraft.company.trim()) {
      toast("Company name is required", "error");
      return;
    }
    splitMsg.mutate(
      {
        messageId: splitDraft.messageId,
        body: {
          company: splitDraft.company.trim(),
          role: splitDraft.role.trim() || undefined,
        },
      },
      {
        onSuccess: (res) => {
          toast("Moved to new application");
          setSplitDraft(null);
          navigate(`/applications/${res.application_id}`);
        },
        onError: (e) => toast(e.message, "error"),
      },
    );
  };

  const onRemoveMessage = (m: Message) => {
    if (
      !window.confirm(
        "Remove this message from tracking? It won’t appear in any listing after sync.",
      )
    ) {
      return;
    }
    ignoreMsg.mutate(m.id, {
      onSuccess: (res) => {
        toast("Message removed");
        if (res.application_deleted) {
          navigate("/applications");
        } else {
          refetch();
        }
      },
      onError: (e) => toast(e.message, "error"),
    });
  };

  if (isLoading) return <Skeleton className="h-96 w-full" />;
  if (isError || !app) {
    return (
      <EmptyState
        title="Application not found"
        description="It may have been merged or deleted."
        actionLabel="Back to list"
        onAction={() => navigate("/applications")}
      />
    );
  }

  const mergeCandidates = (allApps?.applications ?? []).filter(
    (a) => a.id !== app.id && a.company.toLowerCase() === app.company.toLowerCase(),
  );

  const timeline = [
    ...(data?.events ?? []).map((e) => ({
      kind: "event" as const,
      at: e.occurred_at,
      label: e.type.replace(/_/g, " "),
      detail: e.details?.meeting_url ? String(e.details.meeting_url) : "",
    })),
    ...(data?.messages ?? []).map((m) => ({
      kind: "message" as const,
      at: m.received_at,
      label: m.direction === "in" ? "Inbound" : "Outbound",
      detail: m.subject || m.snippet,
      message: m,
    })),
  ].sort((a, b) => new Date(b.at).getTime() - new Date(a.at).getTime());

  return (
    <div className="mx-auto max-w-4xl space-y-6">
      <div className="flex flex-wrap items-start justify-between gap-4">
        <div>
          <Link to="/applications" className="text-sm text-foreground-muted hover:text-accent">
            ← Applications
          </Link>
          <h1 className="mt-1 text-2xl font-semibold">{app.company || "Unknown company"}</h1>
          <p className="text-foreground-muted">{app.role || "Role not set"}</p>
          {gmailUrl({ threadId: app.thread_id }) ? (
            <GmailLink threadId={app.thread_id} className="mt-1 text-xs font-normal" showIcon>
              Open thread in Gmail
            </GmailLink>
          ) : null}
          <div className="mt-2 flex flex-wrap gap-2">
            <Badge className={cn(STATUS_COLORS[app.status])}>{STATUS_LABELS[app.status]}</Badge>
            {app.awaiting_my_reply ? (
              <Badge className="border-amber-500/30 bg-amber-500/10 text-amber-700 dark:text-amber-300">
                Awaiting your reply
              </Badge>
            ) : null}
            {app.status_overridden ? (
              <span className="text-xs text-foreground-muted">Status manually set</span>
            ) : null}
          </div>
        </div>
        <div className="flex flex-wrap gap-2">
          <Button variant="danger" onClick={onDeleteListing} disabled={deleteApp.isPending}>
            Delete listing
          </Button>
          <Button variant="secondary" onClick={save} disabled={patch.isPending}>
            {patch.isPending ? "Saving…" : "Save changes"}
          </Button>
        </div>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        <section className="panel space-y-3 p-4">
          <h2 className="text-sm font-medium">Details</h2>
          <label className="block text-xs text-foreground-muted">
            Company
            <input
              value={company}
              onChange={(e) => setCompany(e.target.value)}
              className="mt-1 h-9 w-full rounded-lg border border-border bg-surface-elevated px-2 text-sm text-foreground"
            />
          </label>
          <label className="block text-xs text-foreground-muted">
            Role
            <input
              value={role}
              onChange={(e) => setRole(e.target.value)}
              className="mt-1 h-9 w-full rounded-lg border border-border bg-surface-elevated px-2 text-sm text-foreground"
            />
          </label>
          <label className="block text-xs text-foreground-muted">
            Status
            <select
              value={status}
              onChange={(e) => setStatus(e.target.value as ApplicationStatus)}
              className="mt-1 h-9 w-full rounded-lg border border-border bg-surface-elevated px-2 text-sm text-foreground"
            >
              {(Object.keys(STATUS_LABELS) as ApplicationStatus[]).map((s) => (
                <option key={s} value={s}>
                  {STATUS_LABELS[s]}
                </option>
              ))}
            </select>
          </label>
          <label className="block text-xs text-foreground-muted">
            Notes
            <textarea
              value={notes}
              onChange={(e) => setNotes(e.target.value)}
              rows={4}
              className="mt-1 w-full rounded-lg border border-border bg-surface-elevated px-2 py-2 text-sm text-foreground"
            />
          </label>
          <dl className="grid grid-cols-2 gap-2 text-xs text-foreground-muted">
            <div>
              <dt>Applied</dt>
              <dd className="text-foreground">{formatDate(app.applied_at)}</dd>
            </div>
            <div>
              <dt>Last activity</dt>
              <dd className="text-foreground">
                {formatRelative(app.last_inbound_at ?? app.updated_at)}
              </dd>
            </div>
            {app.next_interview_at ? (
              <div className="col-span-2">
                <dt>Next interview</dt>
                <dd className="text-foreground">{formatDate(app.next_interview_at)}</dd>
              </div>
            ) : null}
          </dl>
        </section>

        <section className="panel space-y-3 p-4">
          <h2 className="text-sm font-medium">Timeline</h2>
          {timeline.length === 0 ? (
            <p className="text-sm text-foreground-muted">No messages or events yet.</p>
          ) : (
            <ul className="max-h-[28rem] space-y-3 overflow-y-auto pr-1">
              {timeline.map((item, i) => (
                <li key={`${item.at}-${i}`} className="border-l-2 border-border pl-3">
                  <p className="text-[10px] uppercase tracking-wide text-foreground-muted">
                    {item.label} · {formatRelative(item.at)}
                  </p>
                  {item.kind === "message" && item.message ? (
                    <GmailLink
                      threadId={item.message.thread_id}
                      gmailId={item.message.gmail_id}
                      className="text-sm font-normal text-foreground"
                    >
                      {item.detail || "(no subject)"}
                    </GmailLink>
                  ) : (
                    <p className="text-sm">{item.detail || "—"}</p>
                  )}
                </li>
              ))}
            </ul>
          )}
        </section>
      </div>

      {mergeCandidates.length > 0 ? (
        <section className="panel p-4">
          <h2 className="text-sm font-medium">Merge duplicate</h2>
          <p className="mt-1 text-xs text-foreground-muted">
            Combine another application with the same company into this one.
          </p>
          <div className="mt-3 flex flex-wrap gap-2">
            <select
              value={mergeSource}
              onChange={(e) =>
                setMergeSource(e.target.value ? Number(e.target.value) : "")
              }
              className="h-9 min-w-[12rem] rounded-lg border border-border bg-surface-elevated px-2 text-sm"
            >
              <option value="">Select application…</option>
              {mergeCandidates.map((c) => (
                <option key={c.id} value={c.id}>
                  #{c.id} — {c.role || "no role"} ({STATUS_LABELS[c.status]})
                </option>
              ))}
            </select>
            <Button variant="secondary" onClick={onMerge} disabled={mergeSource === "" || merge.isPending}>
              Merge into this
            </Button>
          </div>
        </section>
      ) : null}

      <section className="panel p-4">
        <div className="flex items-center justify-between gap-2">
          <h2 className="text-sm font-medium">Messages</h2>
          <Button variant="secondary" size="sm" onClick={() => refetch()}>
            Refresh
          </Button>
        </div>
        <p className="mt-1 text-xs text-foreground-muted">
          Remove unrelated mail, or split a message into its own application listing.
        </p>

        {splitDraft ? (
          <div className="mt-4 rounded-lg border border-accent/30 bg-surface-muted/50 p-3">
            <p className="text-sm font-medium">New application from message</p>
            <div className="mt-2 grid gap-2 sm:grid-cols-2">
              <label className="block text-xs text-foreground-muted">
                Company
                <input
                  value={splitDraft.company}
                  onChange={(e) =>
                    setSplitDraft({ ...splitDraft, company: e.target.value })
                  }
                  className="mt-1 h-9 w-full rounded-lg border border-border bg-surface-elevated px-2 text-sm text-foreground"
                />
              </label>
              <label className="block text-xs text-foreground-muted">
                Role
                <input
                  value={splitDraft.role}
                  onChange={(e) => setSplitDraft({ ...splitDraft, role: e.target.value })}
                  className="mt-1 h-9 w-full rounded-lg border border-border bg-surface-elevated px-2 text-sm text-foreground"
                />
              </label>
            </div>
            <div className="mt-3 flex flex-wrap gap-2">
              <Button size="sm" onClick={confirmSplit} disabled={splitMsg.isPending}>
                Create listing
              </Button>
              <Button variant="secondary" size="sm" onClick={() => setSplitDraft(null)}>
                Cancel
              </Button>
            </div>
          </div>
        ) : null}

        <ul className="mt-3 divide-y divide-border">
          {(data?.messages ?? []).map((m) => (
            <li key={m.id} className="py-3 first:pt-0">
              <div className="flex flex-wrap items-start justify-between gap-2">
                <div className="min-w-0 flex-1">
                  <div className="flex flex-wrap items-center gap-2 text-xs text-foreground-muted">
                    <Badge>{m.direction === "in" ? "In" : "Out"}</Badge>
                    <span>{formatRelative(m.received_at)}</span>
                    {m.event_type ? <span>{m.event_type}</span> : null}
                  </div>
                  <GmailLink
                    threadId={m.thread_id}
                    gmailId={m.gmail_id}
                    className="mt-1 font-medium text-foreground no-underline hover:text-accent hover:underline"
                    showIcon
                  >
                    {m.subject || "(no subject)"}
                  </GmailLink>
                  <p className="mt-1 line-clamp-3 text-sm text-foreground-muted">
                    {m.body_excerpt || m.snippet}
                  </p>
                  <p className="mt-1 text-xs text-foreground-muted">
                    {m.from_addr} → {m.to_addr}
                  </p>
                </div>
                <div className="flex shrink-0 flex-col gap-1 sm:flex-row">
                  <Button
                    variant="secondary"
                    size="sm"
                    onClick={() => openSplit(m)}
                    disabled={splitMsg.isPending || splitDraft !== null}
                  >
                    Split out
                  </Button>
                  <Button
                    variant="danger"
                    size="sm"
                    onClick={() => onRemoveMessage(m)}
                    disabled={ignoreMsg.isPending}
                  >
                    Remove
                  </Button>
                </div>
              </div>
            </li>
          ))}
        </ul>
      </section>
    </div>
  );
}
