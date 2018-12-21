// Copyright 2015 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License")
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"fmt"
	"time"
	"log"
)

func testAfter(){
	fmt.Println("====1====")
	start := time.Now().Second()
	fmt.Println(start)
	ta := time.After(2 * time.Second)
	fmt.Println("====2====")
	fmt.Println("====3====")
	<- ta
	end := time.Now().Second()
	fmt.Println(end)
	fmt.Println(end - start)
	fmt.Println("====4=====")
}

func testAfterFunc(){
	ch  := make(chan int, 1)
	fmt.Println("====1====")
	start := time.Now().Second()
	fmt.Println(start)
	ta := time.AfterFunc(2 * time.Second, func() {
		fmt.Println("====here====")
	})
	fmt.Println("====2====")
	fmt.Println("====3====")
	<- ta.C
	<- ch
	end := time.Now().Second()
	fmt.Println(end)
	fmt.Println(end - start)
	fmt.Println("====4=====")

	time.Sleep(time.Second * 25)
}

func testTicker(){
	ticker := time.NewTicker(2 * time.Second)
	for{
		select{
		case <- ticker.C:
			log.Printf("the ticker is running!")
		    fmt.Println("=========execute=========")
		default:
			//log.Println("the ticker is not execute!!!")
		}
	}

}


func testTimerChannel(ch chan bool){
	ticker := time.NewTicker(2 * time.Second)
	for {
		select {
		case <- ticker.C:
			log.Printf("the ticker is running!")
			fmt.Println("=========execute=========")
		case <- ch:
			fmt.Println("hello world")
		default:
			// igore
		}
	}
}

func main(){
	// testAfter()
	// testAfterFunc()
	// testTicker()
	ch := make(chan bool, 1)
	go func() {ch <- true}()
	testTimerChannel(ch)
}
