package serverstatus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/urfave/negroni"
)

type SsRequest struct {
	// TBD. What should we support?
}

type SsMiddleware struct {
	startTimeUnix int64
	busyWorkers   int64
	totalAccesses int64
	totalBytes    int64
	stats         []SsRequest
	mutex         *sync.Mutex
}

type ssRequestJsonResponse struct {
	// TBD
}

type ssJsonResponse struct {
	Uptime        string                  `json:"Uptime"`
	TotalAccesses string                  `json:"TotalAccesses"`
	TotalKbytes   string                  `json:"TotalKbytes"`
	BusyWorkers   string                  `json:"BusyWorkers"`
	IdleWorkers   string                  `json:"IdleWorkers"`
	Stats         []ssRequestJsonResponse `json:"stats"`
}

// Middleware is a struct that has a ServeHTTP method
func NewMiddleware() *SsMiddleware {
	return &SsMiddleware{startTimeUnix: time.Now().Unix(), totalAccesses: 0, totalBytes: 0, busyWorkers: 0, mutex: new(sync.Mutex)}
}

func (s *SsMiddleware) HandleServerStatus(w http.ResponseWriter, req *http.Request) {
	if strings.Contains(req.URL.RawQuery, "json") {
		stats := ssJsonResponse{
			Uptime:        fmt.Sprintf("%d", time.Now().Unix()-s.startTimeUnix),
			TotalAccesses: fmt.Sprintf("%d", s.totalAccesses),
			TotalKbytes:   fmt.Sprintf("%d", s.totalBytes/1024),
			BusyWorkers:   fmt.Sprintf("%d", s.busyWorkers),
			IdleWorkers:   "0", // XXX it's infinity!
			Stats:         make([]ssRequestJsonResponse, 0),
		}
		res, _ := json.Marshal(stats)
		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
	} else {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("this is from middleware"))
	}
}

// ServeHTTP implements negroni.Handler interface
func (s *SsMiddleware) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	s.mutex.Lock()
	s.busyWorkers++
	s.mutex.Unlock()

	// Trap request when it is (GET|HEAD) /server-status
	if req.URL.Path == "/server-status" && (req.Method == "GET" || req.Method == "HEAD") {
		s.HandleServerStatus(w, req)
	} else {
		next(w, req)
	}

	res := w.(negroni.ResponseWriter)
	s.mutex.Lock()
	s.totalAccesses++
	s.busyWorkers--
	s.totalBytes += int64(res.Size())
	s.mutex.Unlock()
}
