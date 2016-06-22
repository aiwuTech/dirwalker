// Copyright 2016 zm@huantucorp.com
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
/*
                   _ooOoo_
                  o8888888o
                  88" . "88
                  (| -_- |)
                  O\  =  /O
               ____/`---'\____
             .'  \\|     |//  `.
            /  \\|||  :  |||//  \
           /  _||||| -:- |||||-  \
           |   | \\\  -  /// |   |
           | \_|  ''\---/''  |   |
           \  .-\__  `-`  ___/-. /
         ___`. .'  /--.--\  `. . __
      ."" '<  `.___\_<|>_/___.'  >'"".
     | | :  `- \`.;`\ _ /`;.`/ - ` : | |
     \  \ `-.   \_ __\ /__ _/   .-` /  /
======`-.____`-.___\_____/___.-`____.-'======
                   `=---='
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
         佛祖保佑       永无BUG
*/
package main

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"time"
	"sync"
	"runtime"
)

const (
	VERSION = "0.0.3"
)

var (
	debug       = kingpin.Flag("debug", "print debug log").Short('d').Default("false").Bool()
	version     = kingpin.Flag("version", "show version of dirwalker").Short('v').Default("false").Bool()
	workPath    = kingpin.Flag("work", "target path, default is current directory").Short('w').String()
	ignorePaths = kingpin.Flag("ignore", "ignore paths, default is nil").Short('i').Strings()
	outFile     = kingpin.Flag("out", "out file path").Short('o').Default("./files.out").String()
	solution    = kingpin.Flag("solution", "there has 2 solutions for this tool, you can use s1 or s2 to choose which solution is using, s1 is sample, and s2 is using channel, default is s1").Short('s').Default("s1").String()
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	kingpin.Parse()

	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if *version {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	beginTime := time.Now()
	logrus.Infof("begin time: %v", beginTime)
	switch *solution {
	case "s1":
		infos, err := ReadDir(*workPath, *ignorePaths)
		if err != nil {
			logrus.Fatalf("read dir return error: %v", err)
		}

		if err := WriteFileInfos(infos, *outFile); err != nil {
			logrus.Fatalf("write file infos return error: %v", err)
		}
	case "s2":
		wg := sync.WaitGroup{}
		writeChan := make(chan *FileInfo, 1000)

		go func() {
			wg.Add(1)
			if err := WriteFileFromReader(writeChan, *outFile); err != nil {
				logrus.Fatalf("write file infos return error: %v", err)
			}
			wg.Done()
		}()

		if err := ReadDir2Writer(*workPath, *ignorePaths, writeChan); err != nil {
			logrus.Fatalf("read dir return error: %v", err)
		}
		wg.Wait()
	}

	logrus.Infof("using solution[%v] cost time: %v", *solution, time.Since(beginTime))
}
