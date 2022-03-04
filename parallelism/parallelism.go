package parallelism

import (
	"fmt"
	"reflect"
	"sync"
)

// RunInParallel runs a given function in parallel, and blocking until all executions are complete.
// num is the number of times to execute the function. parallelism is how many concurrent executions you want.
// this function divides num / parallelism and gives the remainder to the first worker.
func RunInParallel(num, parallelism int, theFunc interface{}) {
	each := num / parallelism
	remainder := num % parallelism
	wg := new(sync.WaitGroup)
	for i := 0; i < parallelism; i++ {
		wg.Add(1)
		numGenerate := each
		if i == 0 {
			numGenerate += remainder
		}
		go func(num int) {
			defer wg.Done()
			for i := 0; i < num; i++ {
				funcValue := reflect.ValueOf(theFunc)
				funcValue.Call(nil)
			}
		}(numGenerate)
	}
	wg.Wait()
}

func ExampleRunInParallel() {
	// this will print "hello world" 100 times, running 10 concurrent executions
	RunInParallel(100, 10, func() { fmt.Println("hello world") })
}
