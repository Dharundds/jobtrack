package classify

// Progress is reported during the classify/group pass on sync.
type Progress struct {
	Phase               string `json:"phase"` // classify | group
	ClassifyTotal       int    `json:"classify_total"`
	ClassifyDone        int    `json:"classify_done"`
	LLMReviewed         int    `json:"llm_reviewed"`
	ApplicationsLinked  int    `json:"applications_linked"`
}

// ProgressFunc is called from classify workers (must be concurrency-safe).
type ProgressFunc func(Progress)