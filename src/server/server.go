// server/server.go
package server

import (
    "fmt"
    "log"
    "net/http"
    "leetwaf/src/security"
)

func Serve() {
    securityMiddleware := security.SecurityMiddleware

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            w.Header().Set("Content-Type", "text/html")
            fmt.Fprint(w, `
	<!html>
	<head>
		<title>Search Page</title>
	</head>
	<body>
		<h1>Simple Search</h1>
		<form action="/search" method="get">
			<input type="text" name="id" placeholder="Enter your search query" />
			<input type="submit" value="Search" />
		</form>
	</body>
	</html>
`)
        } else {
            fmt.Println("MS")
        }
    })

	http.Handle("/search", securityMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("id")
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "Search: %s", query)
	})))
	

    fmt.Println("Server is running on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
