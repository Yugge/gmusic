package clients

import (
	"net/http"
)

func inCookies(c []*http.Cookie, s []string) bool {
	amt := len(s)
	crit := 0
	for _, w := range s {
	inner:
		for _, v := range c {
			if v.Name == w {
				crit += 1
				break inner
			}
		}
	}
	if crit == amt {
		return true
	}
	return false
}
