package user

import (
	"RestAPI/pkg/logging"
	"net/http"
)

func CheckCookie(handler http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		logger := logging.GetLogger()
		session, err := store.Get(r, "cookie-name")
		if err != nil {
			logger.Error(err)
		}
		if auth, ok := session.Values["auth"].(bool); !ok || !auth {
			http.Redirect(w, r, "/", http.StatusNotAcceptable)
			return
		}
		handler.ServeHTTP(w, r)
	}

}
