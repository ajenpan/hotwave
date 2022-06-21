package main

import (
	"bytes"
	"fmt"
	"runtime"
)

var (
	Name       string = "unknown"
	Version    string = "unknown"
	GitCommit  string = "unknown"
	BuildAt    string = "unknown"
	BuildBy    string = runtime.Version()
	RunnningOS string = runtime.GOOS + "/" + runtime.GOARCH
)

func shortVersion() string {
	return Version
}

func longVersion() string {
	buf := bytes.NewBuffer(nil)
	fmt.Println(buf, "project:", Name)
	fmt.Println(buf, "version:", Version)
	fmt.Println(buf, "git commit:", GitCommit)
	fmt.Println(buf, "build at:", BuildAt)
	fmt.Println(buf, "build by:", BuildBy)
	fmt.Fprintln(buf, "Running OS/Arch:", RunnningOS)
	return buf.String()
}
