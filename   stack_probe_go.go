package main

import (
	"fmt"
	"os"
	"runtime"
)

var depth int

func recurse(n int) {
	var buf [1024]byte
	_ = buf

	depth = n
	recurse(n + 1)
}

func main() {
	fmt.Println("start go recurse test")

	defer func() {
		if r := recover(); r != nil {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			fmt.Fprintf(os.Stderr,
				">>> CRASH at depth=%d | StackInuse=%dKB | HeapAlloc=%dKB\n",
				depth, m.StackInuse/1024, m.HeapAlloc/1024)
			fmt.Fprintln(os.Stderr, r)
			os.Exit(1)
		}
	}()

	done := make(chan struct{})
	go func() {
		recurse(1)
		close(done)
	}()
	<-done
}