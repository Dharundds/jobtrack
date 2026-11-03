/** API types matching Go JSON responses. */

export type ApplicationStatus =
  | "applied"
  | "replied"
  | "interview"
  | "offer"
  | "rejected";

export type AuthStatus = {
  connected: boolean;
  email?: string;
  last_synced_at?: string;
};

export type Application = {
  id: number;
  company: string;
  role: string;
  status: ApplicationStatus;
  source: string;
  applied_at?: string;
  notes: string;
  status_overridden: boolean;
  last_inbound_at?: string;
  last_outbound_at?: string;
  awaiting_my_reply: boolean;
  days_since_last_activity: number;
  next_interview_at?: string;
  thread_id?: string;
  created_at: string;
  updated_at: string;
};

export type Message = {
  id: number;
  direction: "in" | "out";
  from_addr: string;
  to_addr: string;
  subject: string;
  snippet: string;
  body_excerpt: string;
  received_at: string;
  event_type: string;
  thread_id?: string;
  gmail_id?: string;
};

export type TimelineEvent = {
  id: number;
  type: string;
  occurred_at: string;
  message_id?: number;
  details: Record<string, unknown>;
};

export type ApplicationDetail = {
  application: Application;
  messages: Message[];
  events: TimelineEvent[];
};

export type Rule = {
  id: number;
  field: "subject" | "body" | "from";
  pattern: string;
  event_type: string;
  priority: number;
  enabled: boolean;
};

export type SyncStatus = {
  running: boolean;
  last_synced_at?: string;
  last_error?: string;
  processed?: number;
};

export type Stats = {
  total_applications: number;
  by_status: Record<string, number>;
  response_rate: number;
  avg_days_to_first_reply: number;
};

export type Settings = {
  sync_interval_minutes: number;
  gmail_query: string;
};

export type ApiErrorBody = {
  error: {
    code: string;
    message: string;
  };
};

export type ApplicationsListParams = {
  status?: string;
  q?: string;
  awaiting_reply?: boolean;
};

export type PatchApplicationInput = {
  company?: string;
  role?: string;
  status?: ApplicationStatus;
  notes?: string;
};

export type MergeApplicationsInput = {
  target_id: number;
  source_ids: number[];
};

export type PatchSettingsInput = {
  sync_interval_minutes?: number;
  gmail_query?: string;
};

export type CreateRuleInput = {
  field: Rule["field"];
  pattern: string;
  event_type: string;
  priority?: number;
  enabled?: boolean;
};

export type PatchRuleInput = Partial<CreateRuleInput>;
