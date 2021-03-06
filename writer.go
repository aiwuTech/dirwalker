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
	"errors"
	"os"
	"github.com/Sirupsen/logrus"
	"fmt"
)

// if file existing, truncate it
// if error happened, stop
func WriteFileInfos(infos []*FileInfo, file string) error {
	if infos == nil {
		return errors.New("empty files")
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString("文件名,文件哈希,文件大小\n")
	for _, info := range infos {
		f.WriteString(fmt.Sprintf("%s\n", info))
	}

	return nil
}

func WriteFileFromReader(writeChan chan *FileInfo, file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString("文件名,文件哈希,文件大小\n")
	for {
		info, ok := <- writeChan
		logrus.Debugf("from chan: %v, ok: %v", info, ok)
		if !ok {
			break
		}

		f.WriteString(fmt.Sprintf("%s\n", info))
	}

	return nil
}