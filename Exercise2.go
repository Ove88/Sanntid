// Go 1.2
// go run helloworld_go.go

package main

import (
    . "fmt"     // Using '.' to avoid prefixing functions with their package names
                //   This is probably not a good idea for large projects...
    "runtime"
    "time"
)
func someThread1(finished chan bool,shared chan int) {
    for s:=1;s<=1000001;s++{
		i := <-shared
		i = i+1
		
		shared <- i
	}
	finished<-true
}
func someThread2(finished chan bool,shared chan int) {
    for s:=1;s<=1000001;s++{
		i := <-shared
		i = i-1
		
		shared <- i
	}
	finished<-true
}

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())    // I guess this is a hint to what GOMAXPROCS does...
                                          // Try doing the exercise both with and without it!
    int i := 0
	done1 := make(chan bool)
	done2 := make(chan bool)
	shared := make(chan int,5)
	go someThread1(done1,shared)
	go someThread2(done2,shared)
	shared<-i
	<-done1
	<-done2
	Println(i)
}	
