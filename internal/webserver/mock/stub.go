package mock

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func ParseFile(filename string) []byte {
	file, e := os.ReadFile(filename)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	return file
}

func StubHandlerCurrent(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		log.Println("Received " + r.Method + " request on Current stub handler. Returning mocked information.")
		w.Header().Set("content-type", "application/json")
		output := ParseFile("./assignment-2/internal/res/current.json")
		fmt.Fprint(w, string(output))
		break
	default:
		http.Error(w, "Method not supported", http.StatusNotImplemented)
	}
}

func StubHandlerHistory(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		log.Println("Received " + r.Method + " request on History stub handler. Returning mocked information.")
		w.Header().Set("content-type", "application/json")
		output := ParseFile("./assignment-2/internal/res/history.json")
		fmt.Fprint(w, string(output))
		break
	default:
		http.Error(w, "Method not supported", http.StatusNotImplemented)
	}
}
