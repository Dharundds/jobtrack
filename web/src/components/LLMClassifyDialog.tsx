import { useEffect, useState } from "react";
import { createPortal } from "react-dom";
import { Button } from "@/components/ui/Button";
import { useLLMClassifyPending, useTriggerLLMClassify } from "@/api/hooks";
import { useToast } from "@/components/ui/Toast";
import { cn } from "@/lib/utils";

const DEFAULT_BATCH = 10;
const MIN_BATCH = 2;

function defaultCount(pending: number) {
  if (pending < MIN_BATCH) {
    return MIN_BATCH;
  }
  return Math.min(DEFAULT_BATCH, pending);
}

function clampCount(value: number, pending: number) {
  if (pending < MIN_BATCH) {
    return MIN_BATCH;
  }
  const max = pending;
  return Math.min(max, Math.max(MIN_BATCH, value));
}

type LLMClassifyDialogProps = {
  open: boolean;
  onClose: () => void;
};

export function LLMClassifyDialog({ open, onClose }: LLMClassifyDialogProps) {
  const { data, refetch, isFetching } = useLLMClassifyPending(open);
  const trigger = useTriggerLLMClassify();
  const { toast } = useToast();
  const [pending, setPending] = useState(0);
  const [count, setCount] = useState(DEFAULT_BATCH);

  useEffect(() => {
    if (open) {
      void refetch();
    }
  }, [open, refetch]);

  useEffect(() => {
    if (data?.pending != null) {
      setPending(data.pending);
    }
  }, [data?.pending]);

  useEffect(() => {
    if (open && pending > 0) {
      setCount(defaultCount(pending));
    }
  }, [open, pending]);

  if (!open) {
    return null;
  }

  const canClassify = pending >= MIN_BATCH;
  const max = pending;
  const validCount = canClassify ? clampCount(count, pending) : MIN_BATCH;

  const start = () => {
    const limit = clampCount(count, pending);
    trigger.mutate(
      { limit },
      {
        onSuccess: (res) => {
          if (res.message?.toLowerCase().includes("already")) {
            toast("LLM classify already running");
          } else {
            toast("LLM classify started");
          }
          onClose();
        },
        onError: (e) => toast(e.message, "error"),
      },
    );
  };

  return createPortal(
    <div
      className="fixed inset-0 z-[100] flex items-center justify-center bg-black/50 p-4"
      role="dialog"
      aria-modal
      aria-labelledby="llm-classify-title"
      onMouseDown={(e) => {
        if (e.target === e.currentTarget) {
          onClose();
        }
      }}
    >
      <div className="panel w-full max-w-md space-y-4 p-5 shadow-lg">
        <h2 id="llm-classify-title" className="text-sm font-semibold text-foreground">
          Classify with LLM
        </h2>
        <p className="text-xs text-foreground-muted">
          Sync stores all matching Gmail and links your pipeline with rules only. Bulk LLM re-checks{" "}
          <strong className="font-medium text-foreground">applications on the Board/Rows</strong>, not leftover
          unlinked mail.
        </p>
        {data?.unlinked_messages != null && data.unlinked_messages > 0 ? (
          <p className="text-[11px] text-foreground-muted">
            {data.unlinked_messages} synced message{data.unlinked_messages === 1 ? "" : "s"} are in the DB but not
            in the pipeline — they are not included here.
          </p>
        ) : null}
        <p className={cn("font-data text-sm", pending === 0 && "text-foreground-muted")}>
          {isFetching && pending === 0
            ? "Loading…"
            : pending === 0
              ? "All pipeline applications are already LLM-reviewed."
              : pending === 1
                ? "1 application needs LLM review (need at least 2 to run a batch)."
                : `${pending} applications need LLM review.`}
        </p>

        {canClassify ? (
          <label className="block text-xs text-foreground-muted">
            Applications to review, newest activity first ({MIN_BATCH}–{max})
            <input
              type="number"
              min={MIN_BATCH}
              max={max}
              value={count}
              onChange={(e) => setCount(Number(e.target.value))}
              onBlur={() => setCount(clampCount(count, pending))}
              className="mt-1 h-9 w-full rounded-md border border-border bg-surface px-2 font-mono text-sm"
            />
          </label>
        ) : null}

        <div className="flex flex-col gap-2 sm:flex-row sm:justify-end">
          <Button variant="ghost" size="sm" onClick={onClose} disabled={trigger.isPending}>
            Cancel
          </Button>
          <Button
            variant="primary"
            size="sm"
            disabled={!canClassify || trigger.isPending}
            onClick={start}
          >
            Review {validCount}
          </Button>
        </div>
      </div>
    </div>,
    document.body,
  );
}