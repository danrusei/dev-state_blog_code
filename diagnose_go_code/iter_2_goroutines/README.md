# Adding Goroutines

Run Time:

```bash
$ time ./iter_2_goroutines
There are 169257, out of 1708337, valid Luhn numbers. 
United States has the biggest # of visitors, with 717217 of hits. 
Europe is the continent with most unique countries that accessed the site more than 1000 times. It has 33 unique countries. 

real	0m2,261s
user	0m15,618s
sys	    0m0,365s
```

Benchmarks:

```bash
$ go test -bench GetStatistics -cpuprofile cpu.pprof
goos: linux
goarch: amd64
pkg: github.com/Danr17/dev-state_blog_code/tree/master/diagnose_go_code/iter_2_goroutines
BenchmarkGetStatistics-8   	       1	2465210319 ns/op
PASS
ok  	github.com/Danr17/dev-state_blog_code/tree/master/diagnose_go_code/iter_2_goroutines	2.615s

########

$ go test -bench GetStatistics -memprofile mem.pprof -benchmem
goos: linux
goarch: amd64
pkg: github.com/Danr17/dev-state_blog_code/tree/master/diagnose_go_code/iter_2_goroutines
BenchmarkGetStatistics-8   	       1	2257857809 ns/op	2817558408 B/op	25624869 allocs/op
PASS
ok  	github.com/Danr17/dev-state_blog_code/tree/master/diagnose_go_code/iter_2_goroutines	2.262s

```

CPU Profile:  
![iter 2 cpu](imgs/iter2_cpu.png "Iter 2 CPU")

MEM Profile:  
![iter 2 mem](imgs/iter2_mem.png "Iter 2 MEM")

Trace:  
![iter 2 trace](imgs/iter2_trace.png "Iter 2 Trace")

Goroutine analysis:  
![iter 2 trace](imgs/iter2_goroutines.png "Iter 2 Goroutines")
