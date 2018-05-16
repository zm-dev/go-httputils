package httputils

import (
	"strings"
	"net/http"
)

func ExpectsJson(r *http.Request) bool {
	return (IsAJAX(r) && ! IsPJAX(r) && AcceptsAnyContentType(r)) || WantsJson(r)
}

func IsAJAX(r *http.Request) bool {
	return IsXmlHttpRequest(r)
}

func IsXmlHttpRequest(r *http.Request) bool {
	return "XMLHttpRequest" == r.Header.Get("X-Requested-With")
}

func IsPJAX(r *http.Request) bool {
	return "true" == r.Header.Get("X-PJAX")
}

func AcceptsAnyContentType(r *http.Request) bool {
	accept := strings.SplitN(r.Header.Get("Accept"), ",", 2)[0]
	return "" == accept || "*/*" == accept || "*" == accept
}

func WantsJson(r *http.Request) bool {
	accept := strings.SplitN(r.Header.Get("Accept"), ",", 2)[0]
	for _, s := range []string{"/json", "+json"} {
		if strings.Contains(accept, s) {
			return true
		}
	}
	return false
}

func PreviousURL(r *http.Request) string {
	referer := r.Header.Get("referer")
	// todo: getPreviousUrlFromSession
	if referer != "" {
		return referer
	} else {
		return "/"
	}
}
