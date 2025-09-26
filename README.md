## Stack Probe Experiments

This project contains two simple programs to demonstrate how **stack limits** affect
native threads vs. Go goroutines vs. Node.js recursion.

- `main.go` (Go): runs recursion inside a goroutine. Go uses **dynamic stacks** that
  grow until ~1 GB per goroutine, independent of OS thread `ulimit` stack size.
- `stack_probe_node.js` (Node.js): runs recursion in the JS main thread. V8 enforces
  a **fixed call stack** limit, configurable via `--stack-size`.

---

## 1. Run the Go version locally
```bash
# build and run
go run main.go
```     

Expected output (when crashing):
```
start go recurse test
>>> CRASH at depth=XXXXX | StackInuse=YYYYKB | HeapAlloc=ZZZZKB
runtime: goroutine stack exceeds 1000000000-byte limit
fatal error: stack overflow
```

## 2. Run the Node.js version locally
```bash
# default stack size (~1 MB)
node stack_probe_node.js

# limit stack size to 256 KB
node --stack-size=256 stack_probe_node.js

# larger stack size (e.g. 2 MB)
node --stack-size=2048 stack_probe_node.js
```
Expected output (when crashing):

```bash
start node recurse test
>>> CRASH at depth=1798 | rss=60364KB | heapUsed=11964KB | external=1441KB
Maximum call stack size exceeded
```

## 3. Run the Go version in Docker with limited thread stack
Dockerfile is provided. Build it:
```bash
docker build -t stack-probe-go .
``` 

Run it with limited stack size (e.g. 256 KB):
```bash
docker run --rm -it --ulimit stack=256:256 stack-probe-go 
```
You will see that Go goroutines still grow stacks dynamically up to ~1 GB before crashing,
not being stopped by the OS thread stack limit.

```bash
minhha@LE11-D5482:thread-test$ docker run --rm --ulimit stack=262144:262144 go-stack-test
start go recurse test
depth=10000 | StackInuse=65856KB | HeapAlloc=93KB
depth=20000 | StackInuse=131392KB | HeapAlloc=93KB
depth=30000 | StackInuse=262464KB | HeapAlloc=93KB
depth=40000 | StackInuse=262464KB | HeapAlloc=93KB
depth=50000 | StackInuse=524608KB | HeapAlloc=93KB
depth=60000 | StackInuse=524608KB | HeapAlloc=93KB
depth=70000 | StackInuse=524608KB | HeapAlloc=93KB
depth=80000 | StackInuse=524608KB | HeapAlloc=93KB
depth=90000 | StackInuse=524608KB | HeapAlloc=93KB
runtime: goroutine stack exceeds 1000000000-byte limit
runtime: sp=0xc0200e1230 stack=[0xc0200e0000, 0xc0400e0000]
fatal error: stack overflow

runtime stack:
runtime.throw({0x49f1cc?, 0xc000093e01?})
        /usr/local/go/src/runtime/panic.go:1023 +0x5c fp=0xc000093e18 sp=0xc000093de8 pc=0x433afc
runtime.newstack()
        /usr/local/go/src/runtime/stack.go:1103 +0x5bd fp=0xc000093fc8 sp=0xc000093e18 pc=0x44d37d
runtime.morestack()
        /usr/local/go/src/runtime/asm_amd64.s:616 +0x7a fp=0xc000093fd0 sp=0xc000093fc8 pc=0x4610fa

goroutine 18 gp=0xc000082540 m=4 mp=0xc000080008 [running]:
```

The logs show that even when Docker enforces a thread stack limit of 256 KB, Go goroutines are not constrained by it.  
