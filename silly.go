package godist

import (
	"io"
	"net/http"
	"strconv"
	"sync"
)

type sillyJob struct {
	writer    io.Writer
	request   *http.Request
	debug     bool
	response  int
	done      chan bool
	waitCount int
	mutex     sync.Mutex
}

func NewSillyJob(w io.Writer, r *http.Request) *sillyJob {
	return &sillyJob{writer: w, request: r, done: make(chan bool)}
}

func (self *sillyJob) Parse() error {
	return nil
}

func (self *sillyJob) Expand() ([]Task, error) {
	tasks := make([]Task, 0)
	tasks = append(tasks, &sillyTask{job: self, a: 1, b: 20})
	tasks = append(tasks, &sillyTask{job: self, a: 300, b: 4000})
	self.waitCount = len(tasks)
	return tasks, nil
}

func (self *sillyJob) Integrate(task Task) Job {
	// Integrate must be protected by mutex because more than one
	// Task can try to Integrate at same time.
	self.mutex.Lock()
	defer self.mutex.Unlock()

	self.response += task.(*sillyTask).c

	self.waitCount--
	if self.waitCount > 0 {
		return nil
	}
	return self
}

func (self *sillyJob) Respond() (int, error) {
	defer func() {
		self.done <- true
	}()
	return self.writer.Write([]byte(strconv.Itoa(self.response)))
}

func (self *sillyJob) Wait() {
	<-self.done
}

type sillyTask struct {
	job     *sillyJob
	a, b, c int
}

func (self *sillyTask) Perform() error {
	self.c = self.a + self.b
	return nil
}

func (self *sillyTask) Integrate() Job {
	return self.job.Integrate(self)
}
