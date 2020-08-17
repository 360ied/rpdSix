package keepalive

import (
	"fmt"
	"github.com/ztrue/tracerr"
	"net/http"
)

func KeepAlive() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "I'm alive")

		if err != nil {
			tracerr.PrintSourceColor(err)
		}
	})

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		panic(err)
	}
}
