# godist

## Description

`godist` is a distributed worker library for the Go language.

## Implementation

`godist` is implemented as a collection interfaces, and a function to
create a basic Workflow. To use `godist`, implement the required
interfaces for your concrete `Job` and `Task` data structures, and
submit jobs to the concrete `Workflow`.

## Example use

### Creating a Workflow

A basic concrete `Workflow` data structure is returned by invoking
`NewBasicWorkflow` as demonstrated in the following example:

```Go
    howManyWorkersForEachPhase := 5
	w := NewBasicWorkflow(howManyWorkersForEachPhase)
	defer w.Quit()
```

### Submitting Jobs to the Workflow

Once a `Workflow` is created, jobs may be submitted to the `Workflow`
either asynchronously provided with a call-back function, or
synchronously by blocking until the `Job` is complete.

Synchronous example:

```go
	var buf bytes.Buffer
	job := newSillyJob(5, 60, &buf)
	w.SubmitAndWait(job)
```

Asynchronous example.

```go
	var buf bytes.Buffer
	job := newSillyJob(1, 20, &buf)
	done := make(chan bool)
	w.SubmitWithCallback(job, func() { done <- true })
	<-done
```

### Other Examples

See example implementations of Job and Task in
[silly.go](https://github.com/karrick/godist/blob/master/workflow.go),
with included unit testing in
[silly_test.go](https://github.com/karrick/godist/blob/master/silly_test.go).

See example implementation of Workflow in
[workflow.go](https://github.com/karrick/godist/blob/master/workflow.go),
with included unit testing in
[workflow_test.go](https://github.com/karrick/godist/blob/master/workflow_test.go).
