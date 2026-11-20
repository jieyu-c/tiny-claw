package server

import (
	"net/http"
	"time"
)

// PingHandler 提供一个极简的 HTTP ping 接口。
// GET /ping -> 200 {"pong": "2026-11-19T18:30:00Z"}
func PingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"pong":"` + time.Now().UTC().Format(time.RFC3339) + `"}`))
}
