// Package sysinfo собирает данные о состоянии сервера (Linux, через /proc).
package sysinfo

import (
	"bufio"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type LoadAvg struct {
	One     float64 `json:"1m"`
	Five    float64 `json:"5m"`
	Fifteen float64 `json:"15m"`
}

type Memory struct {
	TotalMB     uint64 `json:"total_mb"`
	AvailableMB uint64 `json:"available_mb"`
	UsedMB      uint64 `json:"used_mb"`
}

type Info struct {
	Timestamp     time.Time `json:"timestamp"`
	Hostname      string    `json:"hostname"`
	OS            string    `json:"os"`
	Arch          string    `json:"arch"`
	GoVersion     string    `json:"go_version"`
	NumCPU        int       `json:"num_cpu"`
	Goroutines    int       `json:"goroutines"`
	UptimeSeconds *float64  `json:"uptime_seconds,omitempty"`
	LoadAvg       *LoadAvg  `json:"load_avg,omitempty"`
	Memory        *Memory   `json:"memory,omitempty"`
}

// Collect возвращает снимок состояния на текущий момент.
// for test
// Поля из /proc опциональны: если файла нет (не Linux), они просто не попадут в JSON.
func Collect() Info {
	host, _ := os.Hostname()
	info := Info{
		Timestamp:  time.Now().UTC(),
		Hostname:   host,
		OS:         runtime.GOOS,
		Arch:       runtime.GOARCH,
		GoVersion:  runtime.Version(),
		NumCPU:     runtime.NumCPU(),
		Goroutines: runtime.NumGoroutine(),
	}

	if b, err := os.ReadFile("/proc/uptime"); err == nil {
		if v, ok := ParseUptime(string(b)); ok {
			info.UptimeSeconds = &v
		}
	}
	if b, err := os.ReadFile("/proc/loadavg"); err == nil {
		if v, ok := ParseLoadAvg(string(b)); ok {
			info.LoadAvg = &v
		}
	}
	if f, err := os.Open("/proc/meminfo"); err == nil {
		defer f.Close()
		if v, ok := ParseMeminfo(bufio.NewScanner(f)); ok {
			info.Memory = &v
		}
	}
	return info
}

func ParseUptime(s string) (float64, bool) {
	fields := strings.Fields(s)
	if len(fields) < 1 {
		return 0, false
	}
	v, err := strconv.ParseFloat(fields[0], 64)
	return v, err == nil
}

func ParseLoadAvg(s string) (LoadAvg, bool) {
	fields := strings.Fields(s)
	if len(fields) < 3 {
		return LoadAvg{}, false
	}
	var vals [3]float64
	for i := 0; i < 3; i++ {
		v, err := strconv.ParseFloat(fields[i], 64)
		if err != nil {
			return LoadAvg{}, false
		}
		vals[i] = v
	}
	return LoadAvg{One: vals[0], Five: vals[1], Fifteen: vals[2]}, true
}

func ParseMeminfo(sc *bufio.Scanner) (Memory, bool) {
	var total, avail uint64
	var gotTotal, gotAvail bool

	for sc.Scan() {
		fields := strings.Fields(sc.Text()) // "MemTotal: 16384 kB"
		if len(fields) < 2 {
			continue
		}
		kb, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			continue
		}
		switch fields[0] {
		case "MemTotal:":
			total, gotTotal = kb, true
		case "MemAvailable:":
			avail, gotAvail = kb, true
		}
	}
	if !gotTotal || !gotAvail || avail > total {
		return Memory{}, false
	}
	return Memory{
		TotalMB:     total / 1024,
		AvailableMB: avail / 1024,
		UsedMB:      (total - avail) / 1024,
	}, true
}
