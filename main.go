// Copyright © 2017 Marc Vandenbosch
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/bombjack73/goji/cmd"
)

func main() {
	//todo make log and var dir if doesn't exits and add log option //!! elevated privilege
	//os.MkdirAll("/var/run/goji", os.ModePerm)
	pidfile := "/var/run/goji/goji.pid"
	if _, err := os.Stat(pidfile); err == nil {
		fmt.Println("goji already running")
		pid, _ := ioutil.ReadFile(pidfile)
		fmt.Println("To stop it, type: kill -2", string(pid))
		os.Exit(0)
	}
	ioutil.WriteFile(pidfile, []byte(strconv.Itoa(os.Getpid())), os.ModePerm) //todo update to allow correct place on Windows ans MacOS

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		<-signals
		os.Remove(pidfile)
		fmt.Println("byebye")
		os.Exit(0)
	}()
	cmd.Execute()
	os.Remove(pidfile)
}
