package godist

import (
	"bytes"
	"testing"
)

func TestWorkflowSubmitWithCallback(t *testing.T) {
	w := NewWorkflow(5)
	defer w.Quit()
	var buf bytes.Buffer
	j := NewSillyJob(&buf, nil)
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
	w := NewWorkflow(3)
	defer w.Quit()
	var buf bytes.Buffer
	j := NewSillyJob(&buf, nil)
	w.SubmitAndWait(j)

	expected := "4321"
	actual := buf.String()
	if expected != actual {
		t.Errorf("Expected: %#v; Actual: %#v\n", expected, actual)
	}
}
