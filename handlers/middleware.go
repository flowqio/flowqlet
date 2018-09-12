package handlers

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/flowqio/flowqlet/service"
	log "github.com/sirupsen/logrus"
)

var PREFIX_PATH = []string{"/login", "/images", "/img", "/api/v1", "/static/js", "/static/css", "/static/images"}

var ADMIN_PREFIX_PATH = []string{"/admin"}

// Gzip Compression
type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

//Gzip func provide all gzip (expect WebSocket or other Accept-Encoding not include gzip request)
func GZip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//expect WebSocket Upgrade and other Accept-Encoding
		if r.Header.Get("Connection") == "Upgrade" {
			next.ServeHTTP(w, r)
			return
		}

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzw := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		next.ServeHTTP(gzw, r)
	})
}

func matchPath(path string) bool {
	log.Debugf("user path check %s", path)
	for _, v := range PREFIX_PATH {
		if strings.HasPrefix(path, v) {
			return true
		}
	}
	return false
}

func matchAdminPath(path string) bool {
	log.Debugf("admin path check %s", path)
	for _, v := range ADMIN_PREFIX_PATH {
		if strings.HasPrefix(path, v) {
			return true
		}
	}
	return false
}

//SessionCheck provide normal request and api session check
func SignCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO
		//need API Token Check

		log.Debugf("%s  %s ", r.Method, r.URL.Path)

		if r.URL.Path == "/" || r.URL.Path == "/api/v1" {
			next.ServeHTTP(w, r)
			return
		}
		// if r.Header.Get("X-Auth-Security-Code") == "" || r.Header.Get("X-Auth-TS") == "" || r.Header.Get("X-Auth-Nonce") == "" {
		// 	http.Error(w, "Token required", http.StatusUnauthorized)
		// 	return
		// }

		if service.Signature(r.Header.Get("X-Auth-Security-Code"), r.Header.Get("X-Auth-Nonce"), r.Header.Get("X-Auth-TS")) {
			next.ServeHTTP(w, r)
			return
		}

		http.Error(w, "Token required", http.StatusUnauthorized)

		// if matchPath(r.URL.Path) {
		// 	log.Debugf("match %s", r.URL.Path)
		// 	next.ServeHTTP(w, r)
		// 	return
		// }

		// requireAdmin := matchAdminPath(r.URL.Path)
		// session, _ := service.SessionStore().Get(r, "user")

		// if _, ok := session.Values["user"]; !ok {
		// 	http.Redirect(w, r, "/login", 302)
		// 	return
		// }

		// if requireAdmin {
		// 	user, _ := session.Values["user"]
		// 	if user.(model.User).Role != "admin" {
		// 		log.Errorf("%s have no admin permission", user.(model.User).Email)
		// 		http.Redirect(w, r, "/static/", 302)
		// 		return

		// 	}
		// }

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		// next.ServeHTTP(w, r)
	})
}

//HTTPError  is error output warpper
// func HTTPError(w http.ResponseWriter, code int, message string) {
// 	w.WriteHeader(code)
// 	w.Write([]byte(message))
// }
