package docker

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"lazydocker/views"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/guptarohit/asciigraph"
)

const cache = 60

type Usage struct {
	MUsed      uint64    // memory used
	MAvailable uint64    // memory available
	Memory     []float64 // memory usage
	CPU        []float64 // cpu usage
}

func newUsage() *Usage {
	return &Usage{
		Memory: make([]float64, cache),
		CPU:    make([]float64, cache),
	}
}

func (u *Usage) clean() {
	for i := 0; i < cache; i++ {
		u.Memory[i] = 0
		u.CPU[i] = 0
	}
	u.MUsed = 0
	u.MAvailable = 0
}

type ContainerStats struct {
	*Usage

	l     *sync.Mutex
	name  string
	body  io.ReadCloser
	index int
	out   *Usage // output data
}

var cstats = ContainerStats{
	l:     &sync.Mutex{},
	Usage: newUsage(),
	out:   newUsage(),
}

func (s *ContainerStats) clear() {
	if s.body != nil {
		s.body.Close()
	}
	s.Usage.clean()
	s.out.clean()
	s.index = cache - 1
}

func Stats(container string) []byte {
	if cstats.name == container {
		return cstats.plot()
	}
	cstats.clear()
	stats, err := cli.ContainerStats(context.Background(), container, true)
	if err != nil {
		panic(err)
	}
	cstats.name = container
	cstats.body = stats.Body
	go cstats.parseStats()
	return cstats.plot()
}

func (s *ContainerStats) getUsage() {
	s.l.Lock()
	defer s.l.Unlock()
	s.out.MUsed = s.MUsed
	s.out.MAvailable = s.MAvailable
	index := s.index
	for i := 0; i < cache; i++ {
		index = (index + 1) % cache
		s.out.Memory[i] = s.Memory[index]
		s.out.CPU[i] = s.CPU[index]
	}
	// fmt.Println("          ", usage.CPU)
	return
}

func (s *ContainerStats) plot() []byte {
	s.getUsage()
	buf := bufPool.Get().(*bytes.Buffer)
	defer resetBuf(buf)
	buf.WriteString(asciigraph.Plot(s.out.CPU,
		asciigraph.Caption(fmt.Sprintf("cpu usage %.2f%%", s.out.CPU[cache-1])),
		asciigraph.Offset(5),
		asciigraph.Height((views.TerminalHeight-10)/2),
	))
	buf.WriteRune('\n')
	buf.WriteRune('\n')
	buf.WriteString(asciigraph.Plot(s.out.Memory,
		asciigraph.Caption(fmt.Sprintf("memory usage %.2f%%", s.out.Memory[cache-1])),
		asciigraph.Offset(5),
		asciigraph.Height((views.TerminalHeight-10)/2),
	))
	return buf.Bytes()
}

func (s *ContainerStats) parseStats() {
	r := bufio.NewReader(s.body)
	for {
		bytes, err := r.ReadBytes('\n')
		if err != nil {
			return
		}
		s.l.Lock()
		s.index = (s.index + 1) % cache
		var stats types.StatsJSON
		_ = json.Unmarshal(bytes, &stats)
		cstats.MUsed = stats.MemoryStats.Usage - stats.MemoryStats.Stats["cache"]
		cstats.MAvailable = stats.MemoryStats.Limit
		if cstats.MAvailable == 0 {
			cstats.Memory[s.index] = 0
		} else {
			cstats.Memory[s.index] = float64(cstats.MUsed) / float64(cstats.MAvailable) * 100.0
		}

		cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage - stats.PreCPUStats.CPUUsage.TotalUsage)
		sysCPUDelta := float64(stats.CPUStats.SystemUsage - stats.PreCPUStats.SystemUsage)
		if sysCPUDelta == 0 {
			cstats.CPU[s.index] = 0
		} else {
			cstats.CPU[s.index] = cpuDelta / sysCPUDelta * float64(stats.CPUStats.OnlineCPUs) * 100.0
		}
		if cstats.CPU[s.index] > 100 {
			cstats.CPU[s.index] = 100
		}
		// fmt.Println(cstats)
		s.l.Unlock()
	}
}
