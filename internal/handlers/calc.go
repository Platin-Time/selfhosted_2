package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"ci-demo/internal/calc"
)

// Указатели нужны, чтобы отличить 0 от «поле не передано».
type calcRequest struct {
	A  *int   `json:"a"`
	B  *int   `json:"b"`
	Op string `json:"op"`
}

type calcResponse struct {
	A      int     `json:"a"`
	B      int     `json:"b"`
	Op     string  `json:"op"`
	Result float64 `json:"result"`
}

// Calc принимает {"a":2,"b":3,"op":"+"} и возвращает результат в JSON.
func Calc(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var req calcRequest
	if err := dec.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body: "+err.Error())
		return
	}
	if req.A == nil || req.B == nil || req.Op == "" {
		writeError(w, http.StatusBadRequest, `fields "a", "b" and "op" are required`)
		return
	}

	res, err := calc.Calculate(*req.A, *req.B, req.Op)
	if err != nil {
		status := http.StatusUnprocessableEntity
		if errors.Is(err, calc.ErrUnknownOp) {
			status = http.StatusBadRequest
		}
		writeError(w, status, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, calcResponse{A: *req.A, B: *req.B, Op: req.Op, Result: res})
}
