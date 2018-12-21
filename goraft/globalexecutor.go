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
package goraft

import (
	"github.com/low901028/goroutinepool"
	"rpcx/log"
	"time"
)


const (
	HEARTBEAT_INTVERAL_MS time.Duration = time.Second * 5   // 心跳间隔
	LEADER_TIMEOUT_MS time.Duration		= time.Second * 15  // Leader超时检查间隔
	RAMDOM_MS time.Duration				= time.Second * 5   //
	TICK_PERIOD_MS time.Duration        = time.Second * 500
	ADDRESS_SERVER_UPDATE_INTVERAL_MS time.Duration = time.Second * 5
)

// TimerThread: schedule executor thread
type Thread struct {
	name 	string        	// 名称
	daemon 	bool		  	// 后台执行
	run 	Runnable        // 任务
	ticker  *time.Ticker     // 定时器
}

type ScheduledExecutorService struct{
	pool goroutinepool.Pool
	run func()
}

func newThread(run Runnable, ticker *time.Ticker,) Thread{
	thread := Thread{
		name: "com.alibaba.nacos.naming.raft.timer",
		daemon: true,
		run: run,
		ticker: ticker,
	}
	//
	return thread
}

// terminal execute：定期执行
func Register(thread *Thread, duration time.Duration) {
	thread.ticker = time.NewTicker(duration * time.Second)
	for{
		select{
		case <- thread.ticker.C:
			log.Infof("the ticker %s is running!", thread.name)
			thread.run.Run()
		default:
			// igore
			//log.Info("the ticker is not execute!!!")
		}

	}
}


func Register1(thread Thread, delay time.Duration) {
	ch := make(chan bool, 1)

	for {
		select {
		case <- ch:
			log.Infof("the ticker %s is running!", thread.name)
			thread.run.Run()
		default:
			// igore
		}
	}

}






