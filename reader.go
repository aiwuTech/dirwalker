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
	"crypto/sha1"
	"fmt"
	"github.com/Sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type FileInfo struct {
	Name string
	Size int64
	Hash string
}

func (this *FileInfo) String() string {
	return strings.Join([]string{
		this.Name,
		this.Hash,
		fmt.Sprintf("%v", this.Size),
	}, ",")
}

func FileSha1(fpath string) string {
	file, err := os.Open(fpath)
	if err != nil {
		return ""
	}

	sha1N := sha1.New()
	if _, err := io.Copy(sha1N, file); err != nil {
		return ""
	}

	return fmt.Sprintf("%x", sha1N.Sum(nil))
}

func ReadDir(dir string, ignores []string) ([]*FileInfo, error) {
	if dir == "" {
		dir = getCurrentPath()
	}
	dir, _ = filepath.Abs(dir)
	infos := []*FileInfo{}

	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		logrus.Debugf("file path: %s", path)

		if info == nil || err != nil {
			return err
		}

		for _, ignore := range ignores {
			if ignoreMatch(path, ignore) {
				logrus.Debugf("file<%v> match ignore<%v>", path, ignore)
				return nil
			}
		}

		if info.IsDir() {
			return nil
		}

		fi := &FileInfo{
			Name: info.Name(),
			Size: info.Size(),
			Hash: FileSha1(path),
		}

		logrus.Debugf("File: %s", fi)
		infos = append(infos, fi)
		return nil
	}); err != nil {
		return nil, err
	}

	return infos, nil
}

func ReadDir2Writer(dir string, ignores []string, writeChan chan *FileInfo) error {
	if dir == "" {
		dir = getCurrentPath()
	}

	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		logrus.Debugf("file path: %s", path)

		if info == nil || err != nil {
			return err
		}

		for _, ignore := range ignores {
			if ignoreMatch(path, ignore) {
				logrus.Debugf("file<%v> match ignore<%v>", path, ignore)
				return nil
			}
		}

		if info.IsDir() {
			return nil
		}

		fi := &FileInfo{
			Name: info.Name(),
			Size: info.Size(),
			Hash: FileSha1(path),
		}

		logrus.Debugf("File: %s", fi)
		writeChan <- fi

		return nil
	}); err != nil {
		return err
	}
	close(writeChan)
	logrus.Debug("chan has closed")

	return nil
}

// TODO match can be wrong
// fpath can be dir, file
// ignore can be dir, regexp
func ignoreMatch(fpath, ignore string) bool {
	// condition 1: fully match
	if fpath == ignore {
		logrus.Debugln("condition 1 match")
		return true
	}

	// fpath's filepath
	path := filepath.Dir(fpath)

	// condition 2: ignore is dir,but relatively
	// fpath: /xxx/xx/1.jgp ignore: xx
	if !filepath.IsAbs(ignore) {
		pathSep := strings.Split(path, fmt.Sprintf("%v", filepath.Separator))
		if pathSep[len(pathSep)-1] == ignore {
			logrus.Debugln("condition 2 match")
			return true
		}
	}

	// condition 3: ignore is dir, and abs
	if filepath.IsAbs(ignore) {
		if path == filepath.Dir(ignore) {
			logrus.Debugln("condition 3 match")
			return true
		}
	}

	// condition 4: ignore is file, but only filename
	if ignoreFile, err := filepath.Abs(filepath.Join(path, ignore)); err == nil {
		if fpath == ignoreFile {
			logrus.Debugln("condition 4 match")
			return true
		}
	}

	// condition 5: ignore is file, and contains path
	if ignoreFile, err := filepath.Abs(ignore); err == nil {
		if fpath == ignoreFile {
			logrus.Debugln("condition 5 match")
			return true
		}
	}

	// condtion 6: ignore is regexp
	match, _ := regexp.MatchString(ignore, fpath)
	if match {
		logrus.Debugln("condtion 6 match")
		return true
	}

	return false
}
