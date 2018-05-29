package httputils

import (
	"net/http"
	"html/template"
	"encoding/json"
	"go.uber.org/zap"
	"io"
)

var Logger *zap.Logger

func SetZapLogger(logger *zap.Logger) {
	Logger = logger
}

type AppHandler interface {
	ServeHTTPWithError(w http.ResponseWriter, r *http.Request) HTTPError
}

type AppHandleFunc func(w http.ResponseWriter, r *http.Request) HTTPError

func (fn AppHandleFunc) ServeHTTPWithError(w http.ResponseWriter, r *http.Request) HTTPError {
	return fn(w, r)
}

func Warp(appHandle AppHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := appHandle.ServeHTTPWithError(w, r)
		if err != nil {
			if f := err.ErrHandlerFunc(); f != nil {
				f(w, r, err)
			} else {
				DefaultErrorHandleFunc(w, r, err)
			}
		}
	})
}

func WarpFunc(f AppHandleFunc) http.Handler {
	return Warp(f)
}

// todo: 不优雅
var errorTpl *template.Template

func SetErrorTpl(tpl string) {
	errorTpl = template.Must(template.ParseFiles("template/error.html"))
}

func DefaultErrorHandleFunc(w http.ResponseWriter, r *http.Request, err HTTPError) {
	// expectsJson := ExpectsJson(r)
	//if expectsJson {
	//	err.AddHeader("Content-Type", "application.pb/json")
	//} else {
	//	err.AddHeader("Content-Type", "text/html")
	//}

	//for key, vals := range err.Headers() {
	//	for _, val := range vals {
	//		w.Header().Add(key, val)
	//	}
	//}

	if err.StatusCode() >= http.StatusInternalServerError {
		Logger.Error("出了一个大错", zap.Error(err), zap.NamedError("InsideError", err.InsideError()))
	}

	// todo 这里不优雅
	if ExpectsJson(r) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(err.StatusCode())
		b, err := json.Marshal(err)
		if err != nil {
			Logger.Error("json.Marshal() 失败", zap.Error(err))
			w.Write([]byte(err.Error()))
		}
		w.Write(b)
	} else if errorTpl != nil {
		w.Header().Add("Content-Type", "text/html")
		w.WriteHeader(err.StatusCode())
		if err := errorTpl.Execute(w, err); err != nil {
			Logger.Warn("errorTpl.Execute(w, err) 执行失败!", zap.Error(err))
		}
	} else {
		w.WriteHeader(err.StatusCode())
		io.WriteString(w, err.Error())
	}
}
