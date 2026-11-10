import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/api/client";
import type {
  Application,
  ApplicationDetail,
  ApplicationStatus,
  ApplicationsListParams,
  AuthStatus,
  CreateRuleInput,
  MergeApplicationsInput,
  PatchApplicationInput,
  PatchRuleInput,
  PatchSettingsInput,
  Rule,
  Settings,
  Stats,
  SyncStatus,
} from "@/api/types";

export const queryKeys = {
  authStatus: ["auth", "status"] as const,
  applications: (p: ApplicationsListParams) => ["applications", p] as const,
  application: (id: number) => ["application", id] as const,
  stats: ["stats"] as const,
  settings: ["settings"] as const,
  rules: ["rules"] as const,
  syncStatus: ["sync", "status"] as const,
};

function appsParamsKey(p: ApplicationsListParams) {
  return {
    status: p.status ?? "",
    q: p.q ?? "",
    awaiting_reply: p.awaiting_reply ?? false,
  };
}

export function useAuthStatus() {
  return useQuery({
    queryKey: queryKeys.authStatus,
    queryFn: () => apiFetch<AuthStatus>("/api/auth/status"),
    staleTime: 30_000,
  });
}

export function useLogout() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: () => apiFetch<{ ok: boolean }>("/auth/logout", { method: "POST" }),
    onSuccess: () => {
      qc.clear();
    },
  });
}

export function useApplications(params: ApplicationsListParams = {}) {
  const key = appsParamsKey(params);
  return useQuery({
    queryKey: queryKeys.applications(key),
    queryFn: () => {
      const sp = new URLSearchParams();
      if (key.status) sp.set("status", key.status);
      if (key.q) sp.set("q", key.q);
      if (key.awaiting_reply) sp.set("awaiting_reply", "true");
      const qs = sp.toString();
      return apiFetch<{ applications: Application[] }>(
        `/api/applications${qs ? `?${qs}` : ""}`,
      );
    },
  });
}

export function useApplication(id: number) {
  return useQuery({
    queryKey: queryKeys.application(id),
    queryFn: () => apiFetch<ApplicationDetail>(`/api/applications/${id}`),
    enabled: id > 0,
  });
}

export function usePatchApplication() {
  const qc = useQueryClient();
  return useMutation<
    Application,
    Error,
    { id: number; body: PatchApplicationInput },
    { prev?: { applications: Application[] }; listKey: ReturnType<typeof queryKeys.applications> }
  >({
    mutationFn: ({ id, body }: { id: number; body: PatchApplicationInput }) =>
      apiFetch<Application>(`/api/applications/${id}`, {
        method: "PATCH",
        body: JSON.stringify(body),
      }),
    onMutate: async ({ id, body }) => {
      const listKey = queryKeys.applications({ status: "", q: "", awaiting_reply: false });
      if (!body.status) return { listKey };
      await qc.cancelQueries({ queryKey: ["applications"] });
      const prev = qc.getQueryData<{ applications: Application[] }>(listKey);
      qc.setQueryData<{ applications: Application[] }>(listKey, (old) => {
        if (!old) return old;
        return {
          applications: old.applications.map((a) =>
            a.id === id ? { ...a, ...body, status: body.status!, status_overridden: true } : a,
          ),
        };
      });
      return { prev, listKey };
    },
    onError: (_err, _vars, ctx) => {
      if (ctx?.prev && ctx.listKey) qc.setQueryData(ctx.listKey, ctx.prev);
    },
    onSuccess: (data) => {
      qc.invalidateQueries({ queryKey: ["applications"] });
      qc.invalidateQueries({ queryKey: queryKeys.application(data.id) });
      qc.invalidateQueries({ queryKey: queryKeys.stats });
    },
  });
}

export function useMergeApplications() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: MergeApplicationsInput) =>
      apiFetch<Application>("/api/applications/merge", {
        method: "POST",
        body: JSON.stringify(body),
      }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["applications"] });
      qc.invalidateQueries({ queryKey: queryKeys.stats });
    },
  });
}

export function useDeleteApplication() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: number) =>
      apiFetch<{ ok: boolean }>(`/api/applications/${id}`, { method: "DELETE" }),
    onSuccess: () => invalidateApplicationQueries(qc),
  });
}

export type SplitMessageInput = {
  company?: string;
  role?: string;
  status?: ApplicationStatus;
};

export type SplitMessageResult = {
  application: Application;
  application_id: number;
  source_deleted: boolean;
};

export function useSplitMessage(applicationId: number) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ messageId, body }: { messageId: number; body?: SplitMessageInput }) =>
      apiFetch<SplitMessageResult>(`/api/applications/${applicationId}/messages/${messageId}/split`, {
        method: "POST",
        body: JSON.stringify(body ?? {}),
      }),
    onSuccess: (data) => {
      invalidateApplicationQueries(qc);
      qc.invalidateQueries({ queryKey: queryKeys.application(data.application_id) });
      qc.invalidateQueries({ queryKey: queryKeys.application(applicationId) });
    },
  });
}

export function useIgnoreMessage(applicationId: number) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (messageId: number) =>
      apiFetch<{ ok: boolean; application_deleted: boolean }>(
        `/api/applications/${applicationId}/messages/${messageId}`,
        { method: "DELETE" },
      ),
    onSuccess: () => {
      invalidateApplicationQueries(qc);
      qc.invalidateQueries({ queryKey: queryKeys.application(applicationId) });
    },
  });
}

export function useSyncStatus(enabled = true) {
  return useQuery({
    queryKey: queryKeys.syncStatus,
    queryFn: () => apiFetch<SyncStatus>("/api/sync/status"),
    refetchInterval: (q) => (q.state.data?.running ? 2000 : false),
    enabled,
  });
}

export function useTriggerSync() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (opts?: { full?: boolean }) =>
      apiFetch<{ ok: boolean; message?: string; status?: SyncStatus }>(
        opts?.full ? "/api/sync?full=true" : "/api/sync",
        {
          method: "POST",
        },
      ),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.syncStatus });
      qc.invalidateQueries({ queryKey: ["applications"] });
      qc.invalidateQueries({ queryKey: queryKeys.stats });
    },
  });
}

export function useStats() {
  return useQuery({
    queryKey: queryKeys.stats,
    queryFn: () => apiFetch<Stats>("/api/stats"),
  });
}

export function useSettings() {
  return useQuery({
    queryKey: queryKeys.settings,
    queryFn: () => apiFetch<Settings>("/api/settings"),
  });
}

export function usePatchSettings() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: PatchSettingsInput) =>
      apiFetch<Settings>("/api/settings", {
        method: "PATCH",
        body: JSON.stringify(body),
      }),
    onSuccess: () => qc.invalidateQueries({ queryKey: queryKeys.settings }),
  });
}

export function useRules() {
  return useQuery({
    queryKey: queryKeys.rules,
    queryFn: () => apiFetch<{ rules: Rule[] }>("/api/rules"),
  });
}

export function useCreateRule() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: CreateRuleInput) =>
      apiFetch<Rule>("/api/rules", { method: "POST", body: JSON.stringify(body) }),
    onSuccess: () => qc.invalidateQueries({ queryKey: queryKeys.rules }),
  });
}

export function usePatchRule() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, body }: { id: number; body: PatchRuleInput }) =>
      apiFetch<Rule>(`/api/rules/${id}`, {
        method: "PATCH",
        body: JSON.stringify(body),
      }),
    onSuccess: () => qc.invalidateQueries({ queryKey: queryKeys.rules }),
  });
}

export function useDeleteRule() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: number) =>
      apiFetch<{ ok: boolean }>(`/api/rules/${id}`, { method: "DELETE" }),
    onSuccess: () => qc.invalidateQueries({ queryKey: queryKeys.rules }),
  });
}

export function invalidateApplicationQueries(qc: ReturnType<typeof useQueryClient>) {
  qc.invalidateQueries({ queryKey: ["applications"] });
  qc.invalidateQueries({ queryKey: queryKeys.stats });
}
