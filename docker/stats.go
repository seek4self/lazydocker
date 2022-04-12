package docker

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"sync"

	"github.com/docker/docker/api/types"
)

const cache = 60

type Usage struct {
	MUsed      uint64    // memory used
	MAvailable uint64    // memory available
	Memory     []float64 // memory usage
	CPU        []float64 // cpu usage
}

type ContainerStats struct {
	Usage

	l     *sync.Mutex
	name  string
	index int
	stop  chan struct{}
}

var cstats = ContainerStats{
	l:    &sync.Mutex{},
	stop: make(chan struct{}),
}

func (s *ContainerStats) clear() {
	close(s.stop)
	s.CPU = make([]float64, cache)
	s.Memory = make([]float64, cache)
	s.MAvailable = 0
	s.MUsed = 0
	s.index = cache - 1
	s.stop = make(chan struct{})
}

func Stats(container string) Usage {
	if cstats.name == container {
		return cstats.getUsage()
	}
	cstats.clear()
	stats, err := cli.ContainerStats(context.Background(), container, true)
	if err != nil {
		panic(err)
	}
	cstats.name = container
	go cstats.parseStats(stats.Body)
	// buf, err := io.ReadAll(stats.Body)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(string(buf))
	return cstats.getUsage()
}

func (s *ContainerStats) getUsage() Usage {
	s.l.Lock()
	defer s.l.Unlock()
	usage := Usage{
		Memory: make([]float64, cache),
		CPU:    make([]float64, cache),
	}
	usage.MUsed = s.MUsed
	usage.MAvailable = s.MAvailable
	index := s.index
	for i := 0; i < cache; i++ {
		index = (index + 1) % cache
		usage.Memory[i] = s.Memory[index]
		usage.CPU[i] = s.CPU[index]
	}
	return usage
}

func (s *ContainerStats) parseStats(body io.ReadCloser) {
	r := bufio.NewReader(body)
	defer body.Close()
	for {
		select {
		case <-s.stop:
			return
		default:
		}
		bytes, err := r.ReadBytes('\n')
		if err == io.EOF {
			return
		}
		s.l.Lock()
		s.index = (s.index + 1) % cache
		var stats types.StatsJSON
		_ = json.Unmarshal(bytes, &stats)
		cstats.MUsed = stats.MemoryStats.Usage - stats.MemoryStats.Stats["cache"]
		cstats.MAvailable = stats.MemoryStats.Limit
		cstats.Memory[s.index] = float64(cstats.MUsed) / float64(cstats.MAvailable) * 100.0
		cstats.CPU[s.index] = float64(stats.CPUStats.CPUUsage.TotalUsage-stats.PreCPUStats.CPUUsage.TotalUsage) /
			float64(stats.CPUStats.SystemUsage-stats.PreCPUStats.SystemUsage) *
			float64(stats.CPUStats.OnlineCPUs) * 100.0
		s.l.Unlock()
		// fmt.Println(s)
	}
}
