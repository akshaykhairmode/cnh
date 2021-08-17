package main

import "net/http"

func main() {

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {

		rw.Header().Add("Content-Type", "text/html")
		rw.WriteHeader(http.StatusOK)

		q := r.URL.Query()
		message := q.Get("message")

		if message != "" {
			rw.Write([]byte(message))
			return
		}

		rw.Write([]byte("Hello World !!!"))
	})

	//NOTE :: Graceful shutdown not there
	http.ListenAndServe(":5000", nil)
}
