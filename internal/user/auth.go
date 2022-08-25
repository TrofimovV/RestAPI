package user

import (
	"fmt"
	"net/http"
)

func Auth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "cookie-name")
		if err != nil {
			fmt.Print(err)
		}
		if auth, ok := session.Values["auth"].(bool); !ok || !auth {
			http.Redirect(w, r, "/", http.StatusNotAcceptable)
			return
		}
		handler.ServeHTTP(w, r)
	}

}
