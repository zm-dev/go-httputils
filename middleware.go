package httputils

import (
	"net/http"
)

//type Middleware interface {
//	ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
//}
//
//type MiddlewareFunc func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
//
//func (fn MiddlewareFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
//	fn(rw, r, next)
//}

type APPMiddleware interface {
	ServeHTTPWithError(rw http.ResponseWriter, r *http.Request, next AppHandleFunc) HTTPError
}

type APPMiddlewareFunc func(rw http.ResponseWriter, r *http.Request, next AppHandleFunc) HTTPError

func (fn APPMiddlewareFunc) ServeHTTPWithError(rw http.ResponseWriter, r *http.Request, next AppHandleFunc) HTTPError {
	return fn(rw, r, next)
}

type middleware struct {
	appMiddleware APPMiddleware
	next          *middleware
}

func (m middleware) ServeHTTPWithError(rw http.ResponseWriter, r *http.Request) HTTPError {
	return m.appMiddleware.ServeHTTPWithError(rw, r, m.next.ServeHTTPWithError)
}

type Handler struct {
	mws        []APPMiddleware
	middleware middleware
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	Warp(h).ServeHTTP(w, r)
}

func (h *Handler) ServeHTTPWithError(w http.ResponseWriter, r *http.Request) HTTPError {
	return h.middleware.ServeHTTPWithError(w, r)
}

func (h *Handler) Use(mw APPMiddleware) {
	if mw == nil {
		panic("APPMiddleware cannot be nil")
	}
	h.mws = append(h.mws, mw)
	h.middleware = buildMiddleware(h.mws)
}

func (h *Handler) UseFunc(mwFunc APPMiddlewareFunc) {
	h.Use(mwFunc)
}

func New(appMiddleware ...APPMiddleware) *Handler {
	mw := buildMiddleware(appMiddleware)
	return &Handler{middleware: mw}
}

func buildMiddleware(appMw []APPMiddleware) middleware {
	var next middleware
	if len(appMw) == 0 {
		return voidMiddleware()
	} else if len(appMw) > 1 {
		next = buildMiddleware(appMw[1:])
	} else {
		next = voidMiddleware()
	}
	return middleware{appMw[0], &next}
}

func WarpToMiddleware(appHandle AppHandler) APPMiddleware {
	return APPMiddlewareFunc(func(rw http.ResponseWriter, r *http.Request, next AppHandleFunc) HTTPError {
		err := appHandle.ServeHTTPWithError(rw, r)
		if err != nil {
			return err
		}
		return next(rw, r)
	})
}

func WarpFuncToMiddleware(appHandle AppHandleFunc) APPMiddleware {
	return WarpToMiddleware(appHandle)
}

func voidMiddleware() middleware {
	return middleware{
		APPMiddlewareFunc(func(rw http.ResponseWriter, r *http.Request, next AppHandleFunc) HTTPError {
			return nil
		}),
		&middleware{},
	}
}
