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
	"github.com/adjust/bufio"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
	"os"
)

func TestWriteFileInfos(t *testing.T) {
	infos := []*FileInfo{}
	file := "./file.out"
	defer os.Remove("./file.out")

	fileBuffer := bufio.NewBufferString("文件名,文件哈希,文件大小\n")

	for i := 0; i < 10; i++ {
		info := &FileInfo{
			Name: fmt.Sprintf("%v.tmp", i),
			Hash: "da39a3ee5e6b4b0d3255bfef95601890afd80709",
			Size: 0,
		}
		infos = append(infos, info)

		fileBuffer.WriteString(fmt.Sprintf("%s\n", info))
	}

	err := WriteFileInfos(infos, file)
	assert.Nil(t, err, fmt.Sprintf("should be nil, but: %v", err))

	fileBytes, err := ioutil.ReadFile(file)
	assert.Nil(t, err, fmt.Sprintf("should be nil, but: %v", err))

	assert.Equal(t, fileBuffer.Bytes(), fileBytes)
}

func TestWriteFileFromReader(t *testing.T)  {
	c := make(chan *FileInfo, 100)
	file := "./file.out"
	defer os.Remove("./file.out")

	fileBuffer := bufio.NewBufferString("文件名,文件哈希,文件大小\n")

	go func() {
		for i := 0; i < 10; i++ {
			info := &FileInfo{
				Name: fmt.Sprintf("%v.tmp", i),
				Hash: "da39a3ee5e6b4b0d3255bfef95601890afd80709",
				Size: 0,
			}
			c <- info

			fileBuffer.WriteString(fmt.Sprintf("%s\n", info))
		}

		close(c)
	}()

	err := WriteFileFromReader(c, file)
	assert.Nil(t, err, fmt.Sprintf("should be nil, but: %v", err))

	fileBytes, err := ioutil.ReadFile(file)
	assert.Nil(t, err, fmt.Sprintf("should be nil, but: %v", err))

	assert.Equal(t, fileBuffer.Bytes(), fileBytes)
}