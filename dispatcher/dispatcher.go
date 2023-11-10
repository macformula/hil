package dispatcher

import (
	"fmt"
	"io"
	"net/http"
	"time"
)


func Dispatcher() {
	port := 8080

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}

		message := string(body) + " returning..."

		fmt.Printf("Received a request with message: %s\n", message)

		// REPLACE THIS WITH THE ACTUAL TESTS
		time.Sleep(10 * time.Second)

		// Set a custom status code (e.g., 201 Created) and write the response message
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(message))
	})

	fmt.Printf("Listening on :%d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
