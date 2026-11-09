package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCalc(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantResult float64
		wantErr    bool
	}{
		{"сложение", `{"a":2,"b":3,"op":"+"}`, http.StatusOK, 5, false},
		{"вычитание", `{"a":2,"b":5,"op":"-"}`, http.StatusOK, -3, false},
		{"умножение", `{"a":4,"b":6,"op":"*"}`, http.StatusOK, 24, false},
		{"деление", `{"a":7,"b":2,"op":"/"}`, http.StatusOK, 3.5, false},
		{"степень", `{"a":2,"b":10,"op":"**"}`, http.StatusOK, 1024, false},
		{"отрицательная степень", `{"a":2,"b":-1,"op":"**"}`, http.StatusOK, 0.5, false},
		{"нули допустимы", `{"a":0,"b":0,"op":"+"}`, http.StatusOK, 0, false},
		{"деление на ноль", `{"a":1,"b":0,"op":"/"}`, http.StatusUnprocessableEntity, 0, true},
		{"переполнение", `{"a":10,"b":10000,"op":"**"}`, http.StatusUnprocessableEntity, 0, true},
		{"неизвестная операция", `{"a":1,"b":2,"op":"%"}`, http.StatusBadRequest, 0, true},
		{"нет поля b", `{"a":1,"op":"+"}`, http.StatusBadRequest, 0, true},
		{"нет op", `{"a":1,"b":2}`, http.StatusBadRequest, 0, true},
		{"битый JSON", `{"a":1,`, http.StatusBadRequest, 0, true},
		{"не int", `{"a":"x","b":2,"op":"+"}`, http.StatusBadRequest, 0, true},
		{"дробное число", `{"a":1.5,"b":2,"op":"+"}`, http.StatusBadRequest, 0, true},
		{"лишнее поле", `{"a":1,"b":2,"op":"+","c":3}`, http.StatusBadRequest, 0, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/calc", strings.NewReader(tc.body))
			rec := httptest.NewRecorder()

			Calc(rec, req)

			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d, body: %s", rec.Code, tc.wantStatus, rec.Body.String())
			}
			if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
				t.Fatalf("Content-Type = %q, want application/json", ct)
			}

			var got map[string]any
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Fatalf("response is not valid JSON: %v", err)
			}

			if tc.wantErr {
				if msg, ok := got["error"].(string); !ok || msg == "" {
					t.Fatalf(`expected non-empty "error", got %v`, got)
				}
				return
			}
			if res, ok := got["result"].(float64); !ok || res != tc.wantResult {
				t.Fatalf("result = %v, want %v", got["result"], tc.wantResult)
			}
		})
	}
}
