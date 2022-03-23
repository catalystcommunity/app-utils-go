
# Dispatcher
The dispatcher is a worker pattern that can be used when you want to run n things in parallel, with a blocking submit. Calls to `Submit()` block if all workers are busy, this is intentional and helps to avoid imbalanced load. For example imagine a scenario where you have two dispatchers (read, and write). One has n workers reading from an API (read) which then submits work the the other dispatcher (write) which has  n workers writing to a database. If the read is much faster, you could have a scenario where read runs unbounded and reads far too much data into memory, resulting in a crash. In this scenario, since `Submit()` blocks, the `read` dispatcher will block until the `write` dispatcher has available workers, which prevents the faster `read` side from running unbounded.
### Usage
Implement the `WorkHandler` and `Job` interfaces on your own structs that have any data you need, and pass your handler struct when creating a new dispatcher via `NewDispatcher()`

### Example Dispatcher usage
The below example will generate 15 jobs and submit them to the dispatcher with a parallelism of 5. You'll see 5 jobs running at once, then the next 5 when those are done, etc. The sleep is there to demonstrate that the jobs don't run until there are free workers to run them. The log "queued work" is there to demonstrate that the call to submit does indeed block as we expect. You can copy and paste this into https://go.dev/play/ to experiment.

```go  
package main

import (
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/catalystsquad/app-utils-go/logging"
	"github.com/catalystsquad/app-utils-go/parallelism"
)

// MyHandler imlements the WorkHandler interface. It gets the job data, prints the phrase, and then sleeps.
type MyHandler struct{}

// The Implementation
func (h MyHandler) HandleJob(job parallelism.Job) {
	// assert type to my job data struct type
	data := job.GetData().(MyJobData)
	// do work
	fmt.Println(data.Phrase)
	// sleep to show parallel blocking work
	time.Sleep(2 * time.Second)
}

// MyJob implements the Job interface and has a custom struct for the data my job needs to run, in this case `MyJobData` which has a single string field
type MyJob struct {
	JobData MyJobData
}

// The implementation returns the job data
func (j MyJob) GetData() interface{} {
	return j.JobData
}

// MyJobData is a custom struct for job data
type MyJobData struct {
	Phrase string
}

func main() {
	// generate jobs
	phrases := []string{}
	for i := 0; i < 15; i++ {
		phrases = append(phrases, gofakeit.HackerPhrase())
	}
	// instantiate my handler
	handler := MyHandler{}
	// instantiate and start the dispatcher
	inParallel := 5
	dispatcher := parallelism.NewDispatcher(inParallel, handler).Start()
	// queue work
	for _, phrase := range phrases {
		dispatcher.Submit(MyJob{JobData: MyJobData{Phrase: phrase}})
		logging.Log.Info("queued work")
	}
	// wait for work to complete
	dispatcher.Wait()
	logging.Log.Info("work complete")
}
```
