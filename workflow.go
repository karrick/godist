package godist

// Workflow is data structure which allows Job structures to be
// submitted.
type Workflow struct {
	done  chan bool
	queue chan Job
}

// NewWorkflow creates a Workflow data structure, with 'count' number
// of go routines for each phase of Job processing.
func NewWorkflow(count int) *Workflow {
	toParse := make(chan Job)
	toExpand := make(chan Job)
	toPerform := make(chan Task)
	toIntegrate := make(chan Task)
	toRespond := make(chan Job)
	done := make(chan bool)

	self := &Workflow{
		queue: toParse,
		done:  done,
	}

	go func() {
		done := make(chan bool)
		for index := 0; index < count; index++ {
			go func() {
				for job := range toParse {
					if err := job.Parse(); err != nil {
						toRespond <- job
						continue
					}
					toExpand <- job
				}
				done <- true
			}()
		}
		for index := 0; index < count; index++ {
			<-done
		}
		close(toExpand)
	}()

	go func() {
		done := make(chan bool)
		for index := 0; index < count; index++ {
			go func() {
				for job := range toExpand {
					tasks, err := job.Expand()
					if err != nil {
						toRespond <- job
						continue
					}
					for _, task := range tasks {
						toPerform <- task
					}
				}
				done <- true
			}()
		}
		for index := 0; index < count; index++ {
			<-done
		}
		close(toPerform)
	}()

	go func() {
		done := make(chan bool)
		for index := 0; index < count; index++ {
			go func() {
				for task := range toPerform {
					// NOTE: same behavior for error or not (??? change ???)
					if err := task.Perform(); err != nil {
						toIntegrate <- task
						continue
					}
					toIntegrate <- task
				}
				done <- true
			}()
		}
		for index := 0; index < count; index++ {
			<-done
		}
		close(toIntegrate)
	}()

	go func() {
		done := make(chan bool)
		for index := 0; index < count; index++ {
			go func() {
				for task := range toIntegrate {
					job := task.Integrate()
					if job != nil {
						toRespond <- job
					}
				}
				done <- true
			}()
		}
		for index := 0; index < count; index++ {
			<-done
		}
		close(toRespond)
	}()

	go func() {
		done := make(chan bool)
		for index := 0; index < count; index++ {
			go func() {
				for job := range toRespond {
					job.Respond()
				}
				done <- true
			}()
		}
		for index := 0; index < count; index++ {
			<-done
		}
		self.done <- true
	}()

	return self
}

// SubmitWithCallback asynchronously sends a Job to the workflow, and
// calls back the given anonymous function when the Job is completed.
func (self *Workflow) SubmitWithCallback(job Job, fn func()) {
	self.queue <- job
	go func() {
		job.Wait()
		fn()
	}()
}

// SubmitAndWait synchronously sends a Job to the workflow, and waits
// until the job is complete prior to returning to the caller
func (self *Workflow) SubmitAndWait(job Job) {
	self.queue <- job
	job.Wait()
}

// Quit shuts down a workflow, closing each channel in sequence until
// all go routines have stopped.
func (self *Workflow) Quit() {
	close(self.queue)
	<-self.done
}
