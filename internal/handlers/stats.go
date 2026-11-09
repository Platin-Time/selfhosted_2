package handlers

import (
	"net/http"

	"ci-demo/internal/sysinfo"
)

// Stats возвращает текущее состояние сервера в JSON.
func Stats(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, sysinfo.Collect())
}
