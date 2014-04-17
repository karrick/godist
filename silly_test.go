package godist

import (
	"bytes"
	"testing"
)

func TestSillyJobParse(t *testing.T) {
	dj := NewSillyJob(nil, nil)
	if err := dj.Parse(); err != nil {
		t.Errorf("Expected: %#v; Actual: %#v\n", nil, err)
	}
}

func TestSillyJobExpand(t *testing.T) {
	job := NewSillyJob(nil, nil)
	job.debug = true
	job.Parse()

	tasks, _ := job.Expand()

	if job.waitCount != 2 {
		t.Errorf("Expected: %#v; Actual: %#v\n", 2, job.waitCount)
	}

	if len(tasks) != 2 {
		t.Errorf("Expected: %#v; Actual: %#v\n", 2, len(tasks))
	}
	if tasks[0].(*sillyTask).a != 1 {
		t.Errorf("Expected: %#v; Actual: %#v\n", 1, tasks[0].(*sillyTask).a)
	}
	if tasks[0].(*sillyTask).b != 20 {
		t.Errorf("Expected: %#v; Actual: %#v\n", 20, tasks[0].(*sillyTask).b)
	}
}

func TestSillyJobIntegrate(t *testing.T) {
	job := NewSillyJob(nil, nil)

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
	job := NewSillyJob(&buf, nil)
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
	job := NewSillyJob(&buf, nil)
	job.Parse()
	tasks, _ := job.Expand()
	for _, task := range tasks {
		task.Perform()
		task.Integrate()
	}
	if job.response != 4321 {
		t.Errorf("Expected: %#v; Actual: %#v\n", 4321, job.response)
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
	if buf.String() != "4321" {
		t.Errorf("Expected: %#v; Actual: %#v\n", "4321", buf.String())
	}
}
