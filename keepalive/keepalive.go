package keepalive

import (
	"fmt"
	"net/http"
)

func KeepAlive() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "I'm alive")
	})

	panic(http.ListenAndServe(":8080", nil))
}
