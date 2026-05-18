package system

import (
	"sort"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	psnet "github.com/shirou/gopsutil/v3/net"
	psprocess "github.com/shirou/gopsutil/v3/process"
)

type MemInfo struct {
	Used        uint64
	Total       uint64
	UsedPercent float64
}

type DiskInfo struct {
	Path    string
	Used    uint64
	Total   uint64
	Percent float64
}

type NetInfo struct {
	Name     string
	SentRate float64
	RecvRate float64
}

type ProcessInfo struct {
	PID    int32
	Name   string
	CPU    float64
	Memory float32
}

type Metrics struct {
	CPU       []float64
	Memory    MemInfo
	Disks     []DiskInfo
	Network   []NetInfo
	Processes []ProcessInfo
}

type netSnapshot struct {
	sent uint64
	recv uint64
}

var (
	netMu     sync.Mutex
	prevNet   = map[string]netSnapshot{}
	prevNetAt time.Time
)

func Collect() (Metrics, error) {
	cpuPercent, err := cpu.Percent(200*time.Millisecond, true)
	if err != nil {
		return Metrics{}, err
	}

	memStat, err := mem.VirtualMemory()
	if err != nil {
		return Metrics{}, err
	}

	disks, _ := collectDisks()
	network, _ := collectNetwork()
	procs, _ := collectProcesses()

	return Metrics{
		CPU: cpuPercent,
		Memory: MemInfo{
			Used:        memStat.Used,
			Total:       memStat.Total,
			UsedPercent: memStat.UsedPercent,
		},
		Disks:     disks,
		Network:   network,
		Processes: procs,
	}, nil
}

func collectDisks() ([]DiskInfo, error) {
	parts, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	var result []DiskInfo
	seen := map[string]bool{}
	for _, p := range parts {
		if seen[p.Mountpoint] {
			continue
		}
		seen[p.Mountpoint] = true

		usage, err := disk.Usage(p.Mountpoint)
		if err != nil || usage.Total == 0 {
			continue
		}
		result = append(result, DiskInfo{
			Path:    p.Mountpoint,
			Used:    usage.Used,
			Total:   usage.Total,
			Percent: usage.UsedPercent,
		})
	}
	return result, nil
}

func collectNetwork() ([]NetInfo, error) {
	counters, err := psnet.IOCounters(true)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	netMu.Lock()
	defer netMu.Unlock()

	elapsed := now.Sub(prevNetAt).Seconds()
	if elapsed < 0.1 {
		elapsed = 1
	}

	var result []NetInfo
	newSnaps := map[string]netSnapshot{}

	for _, c := range counters {
		if c.BytesSent+c.BytesRecv == 0 {
			continue
		}
		snap := netSnapshot{sent: c.BytesSent, recv: c.BytesRecv}
		newSnaps[c.Name] = snap

		var sentRate, recvRate float64
		if prev, ok := prevNet[c.Name]; ok && !prevNetAt.IsZero() {
			sentRate = float64(c.BytesSent-prev.sent) / elapsed
			recvRate = float64(c.BytesRecv-prev.recv) / elapsed
			if sentRate < 0 {
				sentRate = 0
			}
			if recvRate < 0 {
				recvRate = 0
			}
		}

		result = append(result, NetInfo{
			Name:     c.Name,
			SentRate: sentRate,
			RecvRate: recvRate,
		})
	}

	prevNet = newSnaps
	prevNetAt = now

	return result, nil
}

func collectProcesses() ([]ProcessInfo, error) {
	procs, err := psprocess.Processes()
	if err != nil {
		return nil, err
	}

	var result []ProcessInfo
	for _, p := range procs {
		name, _ := p.Name()
		cpuPct, _ := p.CPUPercent()
		memPct, _ := p.MemoryPercent()
		result = append(result, ProcessInfo{
			PID:    p.Pid,
			Name:   name,
			CPU:    cpuPct,
			Memory: memPct,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CPU > result[j].CPU
	})

	if len(result) > 50 {
		result = result[:50]
	}

	return result, nil
}
