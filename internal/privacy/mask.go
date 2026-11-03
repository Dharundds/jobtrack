package privacy

import "strings"

func MaskEmail(email string) string {
	email = strings.TrimSpace(email)
	at := strings.LastIndex(email, "@")
	if at <= 0 || at == len(email)-1 {
		return "***"
	}
	local := email[:at]
	domain := email[at+1:]
	if len(local) == 0 {
		return "***@" + domain
	}
	return string(local[0]) + "***@" + domain
}
