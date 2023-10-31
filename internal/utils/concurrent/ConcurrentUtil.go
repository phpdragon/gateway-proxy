package concurrent

import "sync"

// Run 并发执行
func Run(n int, function func(no int)) {
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
	Run(5, func(n int) {
		println("test")
	})
}
