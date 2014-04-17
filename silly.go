package godist

import (
	"io"
	"strconv"
	"sync"
)

// sillyJob is a concrete data type that demonstrates how to use
// godist Workflow.
type sillyJob struct {
	// data for the job
	writer   io.Writer
	a, b     int
	debug    bool
	response int

	// required for godist workflow to handle
	done      chan bool
	waitCount int
	mutex     sync.Mutex
}

func newSillyJob(a, b int, w io.Writer) *sillyJob {
	return &sillyJob{a: a, b: b, writer: w, done: make(chan bool)}
}

// Parse will process the input parameters
func (self *sillyJob) Parse() error {
	return nil
}

// Expand will multiplex a single Job into the required Task
// structures for a given job.
func (self *sillyJob) Expand() ([]Task, error) {
	tasks := make([]Task, 0)
	tasks = append(tasks, &sillyTask{job: self, a: self.a, b: self.b})
	tasks = append(tasks, &sillyTask{job: self, a: 300, b: 4000})
	self.waitCount = len(tasks)
	return tasks, nil
}

// Integrate receives and integrates the results from a Task into the
// Job's results.
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

// Respond prepares the job for and sends the response to caller.
func (self *sillyJob) Respond() (int, error) {
	defer func() {
		self.done <- true
	}()
	return self.writer.Write([]byte(strconv.Itoa(self.response)))
}

// Wait bloocks until the Job is done and the response has been sent
// to the caller.
func (self *sillyJob) Wait() {
	<-self.done
}

////////////////////////////////////////////////////////////////

// sillyTask is a concrete data type that demonstrates how to use
// godist Workflow.
type sillyTask struct {
	job     *sillyJob
	a, b, c int
}

// Perform does the required computational load for a Task.
func (self *sillyTask) Perform() error {
	self.c = self.a + self.b
	return nil
}

// Integrate allows the Task to tell it's parent Job that it is
// complete, at which time the Job will integrate the Task results. It
// returns the result from the Job's Integrate method to signal to
// godist Workflow whether it is the last Task to be integrated into
// the Job results or further Tasks remain.
func (self *sillyTask) Integrate() Job {
	return self.job.Integrate(self)
}
