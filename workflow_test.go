package godist

import (
	"bytes"
	"testing"
)

const (
	NumRoutines = 5
)

func TestWorkflowSubmitWithCallback(t *testing.T) {
	w := NewBasicWorkflow(NumRoutines)
	defer w.Quit()
	var buf bytes.Buffer
	j := newSillyJob(1, 20, &buf)
	done := make(chan bool)
	w.SubmitWithCallback(j, func() { done <- true })
	<-done

	expected := "4321"
	actual := buf.String()
	if expected != actual {
		t.Errorf("Expected: %#v; Actual: %#v\n", expected, actual)
	}
}

func TestWorkflowSubmitAndWait(t *testing.T) {
	w := NewBasicWorkflow(NumRoutines)
	defer w.Quit()
	var buf bytes.Buffer
	j := newSillyJob(5, 60, &buf)
	w.SubmitAndWait(j)

	expected := "4365"
	actual := buf.String()
	if expected != actual {
		t.Errorf("Expected: %#v; Actual: %#v\n", expected, actual)
	}
}
