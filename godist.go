package godist

// Job interface represents the functionality that a concrete job data
// structure must provide in order to submit to the workflow.
//
// Parse is responsible for any processing on a Job after
// initialization. It ought return nil if errors were encountered
// during parsing.
//
// Expand expands the Job structure into multiple Task structures, or
// returns a nil slice and an error. The Task structures should have a
// reference to their Job, so the Task results can be Integrated into
// the Job's response.
//
// Integrate receives a Task structure and integrates its results into
// the Job's results. Once the Job has received and integrated results
// from all Task structures, it ought return itself so the Workflow
// can pass the Job to the Respond stage. Prior to receiving all Task
// results, Integrate ought return nil.
//
// Respond processes the Job results and sends the results to the
// caller in the requested format, or sends an error.
//
// Wait should block until the Job has been responded to.
type Job interface {
	Parse() error
	Expand() ([]Task, error)
	Integrate(Task) Job
	Respond() (int, error)
	Wait()
}

// Task interface represents the functionality that a concrete task
// data structure must provide in order to be processed by the
// workflow.
//
// Perform actually performs the work needed by the task, and either
// returns nil or an error, as appropriate.
//
// Integrate invokes the Job's Integrate method above, but is required
// because once the Workflow has performed a Task, it does not have a
// reference to the Task's Job structure. See 'sillyTask'
// implementation of Integrate for an example.
type Task interface {
	Perform() error
	Integrate() Job
}

// Workflow interface allows different implementations of a workflow.
type Workflow interface {
	SubmitWithCallback(Job, func())
	SubmitAndWait(Job)
	Quit()
}
