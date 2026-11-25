import { useEffect, useState } from "react";
import {
  useCreateRule,
  useDeleteRule,
  usePatchRule,
  usePatchSettings,
  useRules,
  useSettings,
  useTestLLM,
  useTriggerSync,
} from "@/api/hooks";
import type { Rule } from "@/api/types";
import { Button } from "@/components/ui/Button";
import { EmptyState } from "@/components/ui/EmptyState";
import { Skeleton } from "@/components/ui/Skeleton";
import { useToast } from "@/components/ui/Toast";

export function SettingsPage() {
  const { data: settings, isLoading: settingsLoading } = useSettings();
  const patchSettings = usePatchSettings();
  const testLLM = useTestLLM();
  const { data: rulesData, isLoading: rulesLoading } = useRules();
  const createRule = useCreateRule();
  const patchRule = usePatchRule();
  const deleteRule = useDeleteRule();
  const triggerSync = useTriggerSync();
  const { toast } = useToast();

  const [gmailQuery, setGmailQuery] = useState("");
  const [syncMinutes, setSyncMinutes] = useState(30);
  const [llmBaseURL, setLlmBaseURL] = useState("");
  const [llmModel, setLlmModel] = useState("");
  const [llmAPIKey, setLlmAPIKey] = useState("");
  const [llmKeySet, setLlmKeySet] = useState(false);
  const [llmConcurrency, setLlmConcurrency] = useState(2);

  const [newField, setNewField] = useState<Rule["field"]>("subject");
  const [newPattern, setNewPattern] = useState("");
  const [newEvent, setNewEvent] = useState("application_sent");

  useEffect(() => {
    if (!settings) return;
    setGmailQuery(settings.gmail_query);
    setSyncMinutes(settings.sync_interval_minutes);
    setLlmBaseURL(settings.llm_base_url ?? "");
    setLlmModel(settings.llm_model ?? "");
    setLlmKeySet(settings.llm_api_key_set ?? false);
    setLlmConcurrency(settings.llm_concurrency ?? 2);
    setLlmAPIKey("");
  }, [settings]);

  const saveSettings = () => {
    const body: Parameters<typeof patchSettings.mutate>[0] = {
      gmail_query: gmailQuery,
      sync_interval_minutes: syncMinutes,
      llm_base_url: llmBaseURL,
      llm_model: llmModel,
      llm_concurrency: llmConcurrency,
    };
    if (llmAPIKey.trim() !== "") {
      body.llm_api_key = llmAPIKey;
    }
    patchSettings.mutate(body, {
      onSuccess: () => toast("Settings saved"),
      onError: (e) => toast(e.message, "error"),
    });
  };

  const addRule = () => {
    if (!newPattern.trim()) {
      toast("Pattern is required", "error");
      return;
    }
    createRule.mutate(
      {
        field: newField,
        pattern: newPattern.trim(),
        event_type: newEvent.trim(),
        enabled: true,
      },
      {
        onSuccess: () => {
          toast("Rule added");
          setNewPattern("");
        },
        onError: (e) => toast(e.message, "error"),
      },
    );
  };

  const rules = rulesData?.rules ?? [];

  return (
    <div className="mx-auto max-w-3xl space-y-8">
      <header className="space-y-3">
        <h1 className="page-title">Settings</h1>
        <p className="page-lead">Local LLM, Gmail sync, and regex classification rules.</p>
        <nav className="flex flex-wrap gap-2 text-xs" aria-label="Settings sections">
          <a href="#llm-settings" className="rounded-md border border-border px-2 py-1 hover:text-accent">
            LLM
          </a>
          <a href="#sync-settings" className="rounded-md border border-border px-2 py-1 hover:text-accent">
            Sync
          </a>
          <a href="#rules-settings" className="rounded-md border border-border px-2 py-1 hover:text-accent">
            Rules
          </a>
        </nav>
      </header>

      <section className="panel space-y-6 p-4">
        {settingsLoading ? (
          <Skeleton className="h-48" />
        ) : (
          <>
            <div id="llm-settings" className="scroll-mt-24 space-y-3">
              <h2 className="text-sm font-semibold text-foreground">LLM classification (LM Studio)</h2>
              <p className="text-xs text-foreground-muted">
                OpenAI-compatible API.                 Sync classifies with rules only; use the LLM button in the top bar to re-review pipeline
                applications (pick how many, between 2 and the pending total).
              </p>
              <label className="block text-xs text-foreground-muted">
                Base URL
                <input
                  value={llmBaseURL}
                  onChange={(e) => setLlmBaseURL(e.target.value)}
                  className="mt-1 h-9 w-full rounded-md border border-border bg-surface px-2 font-mono text-sm"
                  placeholder="http://localhost:1234/v1"
                />
              </label>
              <label className="block text-xs text-foreground-muted">
                Model name
                <input
                  value={llmModel}
                  onChange={(e) => setLlmModel(e.target.value)}
                  className="mt-1 h-9 w-full rounded-md border border-border bg-surface px-2 font-mono text-sm"
                  placeholder="exact id from LM Studio, e.g. qwen3.5-9b"
                />
              </label>
              <label className="block text-xs text-foreground-muted">
                Parallel LLM requests (1–8, default 2)
                <input
                  type="number"
                  min={1}
                  max={8}
                  value={llmConcurrency}
                  onChange={(e) => setLlmConcurrency(Number(e.target.value))}
                  className="mt-1 h-9 w-24 rounded-md border border-border bg-surface px-2 font-mono text-sm"
                />
              </label>
              <label className="block text-xs text-foreground-muted">
                API key {llmKeySet ? "(saved — leave blank to keep)" : "(optional)"}
                <input
                  type="password"
                  value={llmAPIKey}
                  onChange={(e) => setLlmAPIKey(e.target.value)}
                  autoComplete="off"
                  className="mt-1 h-9 w-full rounded-md border border-border bg-surface px-2 font-mono text-sm"
                  placeholder={llmKeySet ? "••••••••" : "optional"}
                />
              </label>
              <div className="flex flex-wrap gap-2">
                <Button
                  variant="secondary"
                  disabled={testLLM.isPending || !llmBaseURL.trim() || !llmModel.trim()}
                  onClick={() =>
                    testLLM.mutate(undefined, {
                      onSuccess: (r) =>
                        toast(
                          r.reachable ? "LLM reachable" : "LLM not reachable — save settings first if you changed them",
                          r.reachable ? "success" : "error",
                        ),
                      onError: (e) => toast(e.message, "error"),
                    })
                  }
                >
                  Test connection
                </Button>
              </div>
            </div>

            <hr className="border-border" />

            <div id="sync-settings" className="scroll-mt-24 space-y-3">
              <h2 className="text-sm font-semibold text-foreground">Gmail sync</h2>
              <label className="block text-xs text-foreground-muted">
                Gmail search query
                <input
                  value={gmailQuery}
                  onChange={(e) => setGmailQuery(e.target.value)}
                  className="mt-1 h-9 w-full rounded-md border border-border bg-surface px-2 font-mono text-sm"
                  placeholder="e.g. newer_than:90d (subject:application OR ...)"
                />
              </label>
              <label className="block text-xs text-foreground-muted">
                Sync interval (minutes, stored in DB; worker uses env at startup)
                <input
                  type="number"
                  min={5}
                  max={1440}
                  value={syncMinutes}
                  onChange={(e) => setSyncMinutes(Number(e.target.value))}
                  className="mt-1 h-9 w-32 rounded-lg border border-border bg-surface-elevated px-2 text-sm"
                />
              </label>
              <p className="text-xs text-foreground-muted">
                After changing the Gmail query, run <strong className="font-medium text-foreground">Full rescan</strong>{" "}
                so older matching mail is listed again (incremental sync only fetches new mail).
              </p>
            </div>

            <hr className="border-border" />

            <div className="flex flex-wrap gap-2">
              <Button onClick={saveSettings} disabled={patchSettings.isPending}>
                Save settings
              </Button>
              <Button
                variant="secondary"
                disabled={triggerSync.isPending}
                onClick={() => {
                  triggerSync.mutate(
                    { full: true },
                    {
                      onSuccess: (res) => {
                        if (res.message === "sync already running") {
                          toast("Sync already running", "error");
                          return;
                        }
                        toast("Full mailbox rescan started");
                      },
                      onError: (e) => toast(e.message, "error"),
                    },
                  );
                }}
              >
                Full rescan
              </Button>
            </div>
          </>
        )}
      </section>

      <section id="rules-settings" className="panel scroll-mt-24 space-y-4 p-4">
        <h2 className="text-sm font-semibold text-foreground">Classification rules</h2>
        {rulesLoading ? (
          <Skeleton className="h-40" />
        ) : rules.length === 0 ? (
          <EmptyState title="No custom rules" description="Seed rules from migration still apply in the backend." />
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full min-w-[520px] text-left text-sm">
              <thead className="text-xs text-foreground-muted">
                <tr>
                  <th className="pb-2 pr-2">Field</th>
                  <th className="pb-2 pr-2">Pattern</th>
                  <th className="pb-2 pr-2">Event</th>
                  <th className="pb-2 pr-2">On</th>
                  <th className="pb-2" />
                </tr>
              </thead>
              <tbody>
                {rules.map((rule) => (
                  <RuleRow
                    key={rule.id}
                    rule={rule}
                    onToggle={(enabled) =>
                      patchRule.mutate(
                        { id: rule.id, body: { enabled } },
                        { onError: (e) => toast(e.message, "error") },
                      )
                    }
                    onDelete={() =>
                      deleteRule.mutate(rule.id, {
                        onSuccess: () => toast("Rule deleted"),
                        onError: (e) => toast(e.message, "error"),
                      })
                    }
                  />
                ))}
              </tbody>
            </table>
          </div>
        )}

        <div className="border-t border-border pt-4">
          <p className="text-xs font-medium text-foreground-muted">Add rule</p>
          <div className="mt-2 flex flex-wrap gap-2">
            <select
              value={newField}
              onChange={(e) => setNewField(e.target.value as Rule["field"])}
              className="h-9 rounded-lg border border-border bg-surface-elevated px-2 text-sm"
            >
              <option value="subject">subject</option>
              <option value="body">body</option>
              <option value="from">from</option>
            </select>
            <input
              value={newPattern}
              onChange={(e) => setNewPattern(e.target.value)}
              placeholder="regex pattern"
              className="h-9 min-w-[10rem] flex-1 rounded-lg border border-border bg-surface-elevated px-2 font-mono text-sm"
            />
            <input
              value={newEvent}
              onChange={(e) => setNewEvent(e.target.value)}
              placeholder="event_type"
              className="h-9 min-w-[8rem] rounded-lg border border-border bg-surface-elevated px-2 text-sm"
            />
            <Button onClick={addRule} disabled={createRule.isPending}>
              Add
            </Button>
          </div>
        </div>
      </section>
    </div>
  );
}

function RuleRow({
  rule,
  onToggle,
  onDelete,
}: {
  rule: Rule;
  onToggle: (enabled: boolean) => void;
  onDelete: () => void;
}) {
  return (
    <tr className="border-t border-border">
      <td className="py-2 pr-2 font-mono text-xs">{rule.field}</td>
      <td className="max-w-[12rem] truncate py-2 pr-2 font-mono text-xs" title={rule.pattern}>
        {rule.pattern}
      </td>
      <td className="py-2 pr-2 text-xs">{rule.event_type}</td>
      <td className="py-2 pr-2">
        <input
          type="checkbox"
          checked={rule.enabled}
          onChange={(e) => onToggle(e.target.checked)}
          aria-label={`Enable rule ${rule.id}`}
        />
      </td>
      <td className="py-2 text-right">
        <Button variant="danger" size="sm" onClick={onDelete}>
          Delete
        </Button>
      </td>
    </tr>
  );
}
