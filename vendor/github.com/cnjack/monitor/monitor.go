package monitor

import (
	"fmt"
	"math"
	"runtime"
	"time"
)

type SysStatus struct {
	NumGoroutine int

	// General statistics.
	MemAllocated string // bytes allocated and still in use
	MemTotal     string // bytes allocated (even if freed)
	MemSys       string // bytes obtained from system (sum of XxxSys below)
	Lookups      uint64 // number of pointer lookups
	MemMallocs   uint64 // number of mallocs
	MemFrees     uint64 // number of frees

	// Main allocation heap statistics.
	HeapAlloc    string // bytes allocated and still in use
	HeapSys      string // bytes obtained from system
	HeapIdle     string // bytes in idle spans
	HeapInuse    string // bytes in non-idle span
	HeapReleased string // bytes released to the OS
	HeapObjects  uint64 // total number of allocated objects

	// Low-level fixed-size structure allocator statistics.
	//	Inuse is bytes used now.
	//	Sys is bytes obtained from system.
	StackInuse  string // bootstrap stacks
	StackSys    string
	MSpanInuse  string // mspan structures
	MSpanSys    string
	MCacheInuse string // mcache structures
	MCacheSys   string
	BuckHashSys string // profiling bucket hash table
	GCSys       string // GC metadata
	OtherSys    string // other system allocations

	// Garbage collector statistics.
	NextGC       string // next run in HeapAlloc time (bytes)
	LastGC       string // last run in absolute time (ns)
	PauseTotalNs string
	PauseNs      string // circular buffer of recent GC pause times, most recent at [(NumGC+255)%256]
	NumGC        uint32
}

func Monitor() *SysStatus {
	sysStatus := &SysStatus{}
	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)
	sysStatus.NumGoroutine = runtime.NumGoroutine()

	sysStatus.MemAllocated = fileSize(int64(m.Alloc))
	sysStatus.MemTotal = fileSize(int64(m.TotalAlloc))
	sysStatus.MemSys = fileSize(int64(m.Sys))
	sysStatus.Lookups = m.Lookups
	sysStatus.MemMallocs = m.Mallocs
	sysStatus.MemFrees = m.Frees

	sysStatus.HeapAlloc = fileSize(int64(m.HeapAlloc))
	sysStatus.HeapSys = fileSize(int64(m.HeapSys))
	sysStatus.HeapIdle = fileSize(int64(m.HeapIdle))
	sysStatus.HeapInuse = fileSize(int64(m.HeapInuse))
	sysStatus.HeapReleased = fileSize(int64(m.HeapReleased))
	sysStatus.HeapObjects = m.HeapObjects

	sysStatus.StackInuse = fileSize(int64(m.StackInuse))
	sysStatus.StackSys = fileSize(int64(m.StackSys))
	sysStatus.MSpanInuse = fileSize(int64(m.MSpanInuse))
	sysStatus.MSpanSys = fileSize(int64(m.MSpanSys))
	sysStatus.MCacheInuse = fileSize(int64(m.MCacheInuse))
	sysStatus.MCacheSys = fileSize(int64(m.MCacheSys))
	sysStatus.BuckHashSys = fileSize(int64(m.BuckHashSys))
	sysStatus.GCSys = fileSize(int64(m.GCSys))
	sysStatus.OtherSys = fileSize(int64(m.OtherSys))

	sysStatus.NextGC = fileSize(int64(m.NextGC))
	sysStatus.LastGC = fmt.Sprintf("%.1fs", float64(time.Now().UnixNano()-int64(m.LastGC))/1000/1000/1000)
	sysStatus.PauseTotalNs = fmt.Sprintf("%.1fs", float64(m.PauseTotalNs)/1000/1000/1000)
	sysStatus.PauseNs = fmt.Sprintf("%.3fs", float64(m.PauseNs[(m.NumGC+255)%256])/1000/1000/1000)
	sysStatus.NumGC = m.NumGC
	return sysStatus
}

func logn(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}

func fileSize(s int64) string {
	sizes := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	return humanateBytes(uint64(s), 1024, sizes)
}

func humanateBytes(s uint64, base float64, sizes []string) string {
	if s < 10 {
		return fmt.Sprintf("%dB", s)
	}
	e := math.Floor(logn(float64(s), base))
	suffix := sizes[int(e)]
	val := float64(s) / math.Pow(base, math.Floor(e))
	f := "%.0f"
	if val < 10 {
		f = "%.1f"
	}

	return fmt.Sprintf(f+"%s", val, suffix)
}
