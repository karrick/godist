package godist

import (
	"bytes"
	"testing"
)

////////////////////////////////////////////////////////////////

func TestSillyJobParse(t *testing.T) {
	dj := newSillyJob(1, 2, nil)
	if err := dj.Parse(); err != nil {
		t.Errorf("Expected: %#v; Actual: %#v\n", nil, err)
	}
}

func TestSillyJobExpand(t *testing.T) {
	job := newSillyJob(3, 4, nil)
	job.debug = true
	job.Parse()

	tasks, _ := job.Expand()

	if job.waitCount != 2 {
		t.Errorf("Expected: %#v; Actual: %#v\n", 2, job.waitCount)
	}

	if len(tasks) != 2 {
		t.Errorf("Expected: %#v; Actual: %#v\n", 2, len(tasks))
	}
	if tasks[0].(*sillyTask).a != 3 {
		t.Errorf("Expected: %#v; Actual: %#v\n", 3, tasks[0].(*sillyTask).a)
	}
	if tasks[0].(*sillyTask).b != 4 {
		t.Errorf("Expected: %#v; Actual: %#v\n", 4, tasks[0].(*sillyTask).b)
	}
}

func TestSillyJobIntegrate(t *testing.T) {
	job := newSillyJob(5, 6, nil)

	t1 := sillyTask{job: job, c: 13}
	t1.Integrate()
	t2 := sillyTask{job: job, c: 42}
	t2.Integrate()

	if job.response != 55 {
		t.Errorf("Expected: %#v; Actual: %#v\n", 55, job.response)
	}
}

func TestSillyJobWriteTo(t *testing.T) {
	var buf bytes.Buffer
	job := newSillyJob(7, 8, &buf)
	job.response = 55

	go func() {
		job.Wait()
	}()
	n, err := job.Respond()

	if err != nil {
		t.Errorf("Expected: %#v; Actual: %#v\n", nil, err)
	}
	if n != 2 {
		t.Errorf("Expected: %#v; Actual: %#v\n", 2, n)
	}
	if buf.String() != "55" {
		t.Errorf("Expected: %#v; Actual: %#v\n", "55", buf.String())
	}
}

func TestSillyTaskPerform(t *testing.T) {
	task := &sillyTask{a: 1, b: 2}
	_ = task.Perform()
	if task.c != 3 {
		t.Errorf("Expected: %#v; Actual: %#v\n", 3, task.c)
	}
}

func TestSillyWorkflow(t *testing.T) {
	var buf bytes.Buffer
	job := newSillyJob(7, 8, &buf)
	job.Parse()
	tasks, _ := job.Expand()
	for _, task := range tasks {
		task.Perform()
		task.Integrate()
	}
	if job.response != 4315 {
		t.Errorf("Expected: %#v; Actual: %#v\n", 4315, job.response)
	}

	go func() {
		job.Wait()
	}()
	n, err := job.Respond()
	if err != nil {
		t.Errorf("Expected: %#v; Actual: %#v\n", nil, err)
	}
	if n != 4 {
		t.Errorf("Expected: %#v; Actual: %#v\n", 4, n)
	}
	if buf.String() != "4315" {
		t.Errorf("Expected: %#v; Actual: %#v\n", "4315", buf.String())
	}
}
