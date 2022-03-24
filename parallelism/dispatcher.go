package parallelism

import (
	"github.com/google/uuid"
	"sync"
)

// worker is an internal struct, the worker is what handles the work for a given job. The WorkHandler and Job interfaces
// must be implemented by the user.
type worker struct {
	id          uuid.UUID       // worker id
	jobChannel  jobChannel      // a channel to receive a job, a job represents a unit of work
	jobQueue    jobQueue        // shared between all workers
	quit        chan struct{}   // a channel to quit working
	waitGroup   *sync.WaitGroup // waitgroup reference, this should be the Dispatcher's waitgroup so that the Dispatcher can wait for all work to complete
	workHandler WorkHandler     // handler interface, this will handle the work to be done
}

// WorkHandler is the user facing interface that does the work
type WorkHandler interface {
	HandleJob(job Job)
}

// Job is the user facing interface that describes the work to be done
type Job interface {
	GetData() interface{}
}

type jobChannel chan Job
type jobQueue chan chan Job

// newWorker returns a new worker
func newWorker(jobChan jobChannel, queue jobQueue, quit chan struct{}, workHandler WorkHandler, waitGroup *sync.WaitGroup) *worker {
	return &worker{
		id:          uuid.New(),
		jobChannel:  jobChan,
		jobQueue:    queue,
		quit:        quit,
		workHandler: workHandler,
		waitGroup:   waitGroup,
	}
}

// start starts a worker, this means the worker will listen on its channels for jobs or a quit signal
func (wr *worker) start() {
	go func() {
		for {
			// when available, put the jobChannel again on the JobPool
			// and wait to receive a job
			wr.jobQueue <- wr.jobChannel
			select {
			case job := <-wr.jobChannel:
				func() {
					defer wr.waitGroup.Done()
					wr.workHandler.HandleJob(job)
				}()
			case <-wr.quit:
				// a signal on this channel means someone triggered
				// a shutdown for this worker
				close(wr.jobChannel)
				return
			}
		}
	}()
}

// stop closes the quit channel on the worker.
func (wr *worker) stop() {
	close(wr.quit)
}

// NewDispatcher returns a new Dispatcher. Its main job is to receive a job and share it on the WorkPool
// WorkPool is the link between the Dispatcher and all the workers as
// the WorkPool of the Dispatcher is common JobPool for all the workers
func NewDispatcher(parallelism int, workHandler WorkHandler) *Dispatcher {
	return &Dispatcher{
		workers:     make([]*worker, parallelism),
		jobChannel:  make(jobChannel),
		jobQueue:    make(jobQueue),
		workHandler: workHandler,
		waitGroup:   new(sync.WaitGroup),
	}
}

// Dispatcher is the link between the client and the workers
type Dispatcher struct {
	workers     []*worker  // this is the list of workers that Dispatcher tracks
	jobChannel  jobChannel // client submits job to this channel
	jobQueue    jobQueue   // this is the shared JobPool between the workers
	workHandler WorkHandler
	waitGroup   *sync.WaitGroup
}

// Start creates pool of workers, and starts each worker
func (d *Dispatcher) Start() *Dispatcher {
	l := len(d.workers)
	for i := 1; i <= l; i++ {
		// all workers share the Dispatcher's waitgroup
		wrk := newWorker(make(jobChannel), d.jobQueue, make(chan struct{}), d.workHandler, d.waitGroup)
		wrk.start()
		d.workers = append(d.workers, wrk)
	}
	go d.process()
	return d
}

// process listens to a job submitted on jobChannel and
// relays it to the WorkPool. The WorkPool is shared between
// the workers.
func (d *Dispatcher) process() {
	for {
		select {
		case job := <-d.jobChannel: // listen to any submitted job on the jobChannel
			// wait for a worker to submit jobChannel to jobQueue
			// note that this jobQueue is shared among all workers.
			// Whenever there is an available jobChannel on jobQueue pull it
			jobChan := <-d.jobQueue
			// Once a jobChan is available, send the submitted Job on this jobChannel
			jobChan <- job
		}
	}
}

// Submit is how a job is submitted to the Dispatcher, jobs will be handled by a worker
func (d *Dispatcher) Submit(job Job) {
	d.waitGroup.Add(1)
	d.jobChannel <- job
}

// Wait will wait until all work is completed. This is accomplished by sharing the Dispatcher's waitgroup
// with workers.
func (d *Dispatcher) Wait() {
	d.waitGroup.Wait()
}
