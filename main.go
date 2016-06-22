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
)

const (
	VERSION = "0.0.1"
)

var (
	debug       = kingpin.Flag("debug", "print debug log").Short('d').Default("false").Bool()
	version     = kingpin.Flag("version", "show version of dirwalker").Short('v').Default("false").Bool()
	workPath    = kingpin.Flag("work", "target path").Short('w').Default(getCurrentPath()).String()
	ignorePaths = kingpin.Flag("ignore", "ignore path").Short('i').Strings()
	outFile     = kingpin.Flag("out", "out file path").Short('o').Default("./files.out").String()
)

func main() {
	kingpin.Parse()

	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if *version {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	infos, err := ReadDir(*workPath, *ignorePaths)
	if err != nil {
		logrus.Fatalf("read dir return error: %v", err)
	}

	if err := WriteFileInfos(infos, *outFile); err != nil {
		logrus.Fatalf("write file infos return error: %v", err)
	} else {
		logrus.Infof("file infos writed")
	}
}
