package urlfrontier

type Queue struct {
	tasks chan CrawlTask
}

func NewQueue(bufferSize int) *Queue {
	return &Queue{
		tasks: make(chan CrawlTask, bufferSize),
	}
}

func (q *Queue) Enqueue(task CrawlTask) {
	q.tasks <- task
}

func (q *Queue) Dequeue() CrawlTask {
	return <-q.tasks
}

func (q *Queue) DequeueNonBlocking() (CrawlTask, bool) {
	select {
	case task := <-q.tasks:
		return task, true
	default:
		return CrawlTask{}, false
	}
}

func (q *Queue) Size() int {
	return len(q.tasks)
}
