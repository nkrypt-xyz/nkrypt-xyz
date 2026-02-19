package middleware

import "net/http"

// CORS middleware matches the old Express CORS behavior, including nk-crypto-meta.
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers",
			"Origin, X-Requested-With, Content-Type, Accept, Authorization, nk-crypto-meta")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("Ok"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

