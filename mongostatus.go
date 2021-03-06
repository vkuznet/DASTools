package main

// DASTools: mongostatus returns string with basic MongoDB stats
// Copyright (c) 2018 - Valentin Kuznetsov <vkuznet AT gmail dot com>

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func info() string {
	goVersion := runtime.Version()
	tstamp := time.Now()
	return fmt.Sprintf("git={{VERSION}} go=%s date=%s", goVersion, tstamp)
}

func main() {
	var version bool
	flag.BoolVar(&version, "version", false, "Show version")
	var port int
	flag.IntVar(&port, "port", 8230, "MongoDB port")
	flag.Parse()
	if version {
		fmt.Println(info())
		return
	}
	uri := fmt.Sprintf("mongodb://localhost:%d", port)
	session, err := mgo.Dial(uri)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer session.Close()

	r := bson.M{}
	if err := session.DB("admin").Run(bson.D{{"serverStatus", 1}}, &r); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	proc, _ := r["process"]
	pid, _ := r["pid"]
	ver, _ := r["version"]
	host, _ := r["host"]
	c, _ := r["connections"]
	con := c.(bson.M)
	cons := con["totalCreated"]
	s, _ := r["ok"]
	status := fmt.Sprintf("%v", s)
	if s == 1. {
		status = "ok"
	}
	fmt.Printf("Process:%s, PID:%d, Version:%s, Host:%s, Connections:%d, Status:%s\n", proc, pid, ver, host, cons, status)
}
