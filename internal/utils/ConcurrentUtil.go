package utils

import "sync"

// ConcurrentRun 并发执行
func ConcurrentRun(n int, function func(no int)) {
	wg := new(sync.WaitGroup)
	for i := 0; i < n; i++ {
		wg.Add(1)
		//并发执行
		go func(m int) {
			defer wg.Done()
			function(m)
		}(i)
	}
	wg.Wait()
}

func test() {
	ConcurrentRun(5, func(n int) {
		println("test")
	})
}
