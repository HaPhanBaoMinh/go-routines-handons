package main

import (
	"fmt"
	"runtime"
)

var depth int

func recurse(n int) {
	var buf [1024]byte
	_ = buf

	depth = n

	if n%10000 == 0 {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		fmt.Printf("depth=%d | StackInuse=%dKB | HeapAlloc=%dKB\n",
			n, m.StackInuse/1024, m.HeapAlloc/1024)
	}

	recurse(n + 1)
}

func main() {
	fmt.Println("start go recurse test")

	done := make(chan struct{})
	go func() {
		recurse(1)
		close(done)
	}()
	<-done
}
