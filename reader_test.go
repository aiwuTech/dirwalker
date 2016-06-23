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
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func mkTmpDir() {
	os.Mkdir("tmp", 0777)
	for i := 0; i < 10; i++ {
		f, err := os.Create(fmt.Sprintf("./tmp/%v.tmp", i))
		if err != nil {
			panic(err)
		}
		f.Close()
	}
}

func rmTmpDir() {
	os.RemoveAll("./tmp")
}

func TestFileInfoString(t *testing.T) {
	eqFi := &FileInfo{
		Name: "testing",
		Hash: "testing",
		Size: 0,
	}

	assert.Equal(t, "testing,testing,0", eqFi.String())

	neFi := &FileInfo{
		Name: "testing",
		Hash: "testing",
		Size: 1,
	}

	assert.NotEqual(t, "testing,testing,0", neFi.String())
}

func TestFileSha1(t *testing.T) {
	mkTmpDir()
	defer rmTmpDir()

	for i := 0; i < 10; i++ {
		assert.Equal(t, "da39a3ee5e6b4b0d3255bfef95601890afd80709", FileSha1(fmt.Sprintf("./tmp/%v.tmp", i)))
	}
}

func TestReadDir(t *testing.T) {
	mkTmpDir()
	defer rmTmpDir()

	dir := "./tmp"
	ignores := []string{}

	infos, err := ReadDir(dir, ignores)
	assert.Nil(t, err, fmt.Sprintf("should be nil, but: %v", err))
	if assert.NotNil(t, infos) {
		assert.Equal(t, 10, len(infos))
	}
}

func TestReadDir2Writer(t *testing.T) {
	mkTmpDir()
	defer rmTmpDir()

	dir := "./tmp"
	ignores := []string{}
	c := make(chan *FileInfo, 100)

	err := ReadDir2Writer(dir, ignores, c)
	assert.Nil(t, err, fmt.Sprintf("should be nil, but: %v", err))
	assert.Equal(t, 10, len(c))
}
