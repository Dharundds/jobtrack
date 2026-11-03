package classify

// GroupKey identifies messages in the same application bucket before DB ids exist.
type GroupKey struct {
	ThreadID string
	Company  string
	Role     string
}

// PickApplication chooses an existing application id by thread, then company+role.
func PickApplication(
	threadID string,
	company, role string,
	byThread map[string]int64,
	apps []ApplicationCandidate,
) (int64, bool) {
	if id, ok := byThread[threadID]; ok && threadID != "" {
		return id, true
	}
	for _, a := range apps {
		if SimilarCompanyRole(company, role, a.Company, a.Role) {
			return a.ID, true
		}
	}
	return 0, false
}

type ApplicationCandidate struct {
	ID      int64
	Company string
	Role    string
}
