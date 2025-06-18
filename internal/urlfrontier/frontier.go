package urlfrontier

import (
	"sync"
)

type Frontier struct {
	queue       *Queue
	visitedURLs map[string]struct{}
	mu          sync.Mutex
}

func NewFrontier(bufferSize int) *Frontier {
	return &Frontier{
		queue:       NewQueue(bufferSize),
		visitedURLs: make(map[string]struct{}),
	}
}

func (f *Frontier) Add(task CrawlTask) bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	if _, seen := f.visitedURLs[task.URL]; seen {
		return false
	}

	f.visitedURLs[task.URL] = struct{}{}
	f.queue.Enqueue(task)
	return true
}

func (f *Frontier) GetNext() CrawlTask {
	return f.queue.Dequeue()
}

func (f *Frontier) QueueSize() int {
	return f.queue.Size()
}
