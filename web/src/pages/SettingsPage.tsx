import { useEffect, useState } from "react";
import {
  useCreateRule,
  useDeleteRule,
  usePatchRule,
  usePatchSettings,
  useRules,
  useSettings,
} from "@/api/hooks";
import type { Rule } from "@/api/types";
import { Button } from "@/components/ui/Button";
import { EmptyState } from "@/components/ui/EmptyState";
import { Skeleton } from "@/components/ui/Skeleton";
import { useToast } from "@/components/ui/Toast";

export function SettingsPage() {
  const { data: settings, isLoading: settingsLoading } = useSettings();
  const patchSettings = usePatchSettings();
  const { data: rulesData, isLoading: rulesLoading } = useRules();
  const createRule = useCreateRule();
  const patchRule = usePatchRule();
  const deleteRule = useDeleteRule();
  const { toast } = useToast();

  const [gmailQuery, setGmailQuery] = useState("");
  const [syncMinutes, setSyncMinutes] = useState(30);

  const [newField, setNewField] = useState<Rule["field"]>("subject");
  const [newPattern, setNewPattern] = useState("");
  const [newEvent, setNewEvent] = useState("application_sent");

  useEffect(() => {
    if (!settings) return;
    setGmailQuery(settings.gmail_query);
    setSyncMinutes(settings.sync_interval_minutes);
  }, [settings]);

  const saveSettings = () => {
    patchSettings.mutate(
      {
        gmail_query: gmailQuery,
        sync_interval_minutes: syncMinutes,
      },
      {
        onSuccess: () => toast("Settings saved"),
        onError: (e) => toast(e.message, "error"),
      },
    );
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
      <header>
        <h1 className="page-title">Settings</h1>
        <p className="page-lead">Gmail query, sync interval, and classification rules.</p>
      </header>

      <section className="panel space-y-4 p-4">
        <h2 className="text-sm font-medium">Sync</h2>
        {settingsLoading ? (
          <Skeleton className="h-24" />
        ) : (
          <>
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
            <Button onClick={saveSettings} disabled={patchSettings.isPending}>
              Save sync settings
            </Button>
          </>
        )}
      </section>

      <section className="panel space-y-4 p-4">
        <h2 className="text-sm font-medium">Classification rules</h2>
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
