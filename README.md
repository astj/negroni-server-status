# negroni-server-status

A [negroni](https://github.com/urfave/negroni) middleware to monitor HTTP request stats like [`Plack::Middleware::ServerStatus::Lite`](https://metacpan.org/pod/Plack::Middleware::ServerStatus::Lite)

## Current Limitations

Comparing to original `Plack::Middleware::ServerStatus::Lite`, this middleware has following limitations.
- IdleWorker value is always empty
  - `net/http` server does not limit total numbers of workers, so this field is meaningless
- Always returns response in JSON formats
- Response does not contain individual processing requests (`stats` field)

The lattar two are not implemented because I don't need them currently, but patches are welcome.

## Usage

See example/main.go.

```go
package main

import (
	"fmt"
	"net/http"

	nss "github.com/astj/negroni-server-status"
	"github.com/urfave/negroni"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "This is from internal")
	})

	n := negroni.New()

	n.Use(nss.NewMiddleware("/server-status"))

	n.UseHandler(mux)

	n.Run(":3000")
}
```

In this case, a response from `http://localhost:3000/server-status` is somewhat like this:
```json
{"Uptime":"60","TotalAccesses":"67","TotalKbytes":"1","BusyWorkers":"1","IdleWorkers":"","stats":[]}
```

| key         | type of value | meaning                                                                |
|-------------|---------------|------------------------------------------------------------------------|
| Uptime      | string        | Elapsed seconds from this server starts                                |
| TotalAccess | string        | Number of requests processed by this server                            |
| TotalKbytes | string        | Sum of response body size processed by this server, in KB units        |
| BusyWorkers | string        | Number of currently processing requests                                |
| IdleWorkers | string        | Just empty. Available for compatibility with original Plack middleware |
| stats       | []            | Should be information about current connections (not implemented)      |

## License

MIT
