package cli

import (
	"bytes"
	"fmt"
	"runtime"
)

var DefaultVersion = NewAppInfo()

func NewAppInfo() *AppInfo {
	return &AppInfo{
		Name:       "unknown",
		Version:    "unknown",
		GitCommit:  "unknown",
		BuildAt:    "unknown",
		BuildBy:    "runtime.Version()",
		RunnningOS: runtime.GOOS + "/" + runtime.GOARCH,
	}
}

type AppInfo struct {
	Name       string
	Version    string
	GitCommit  string
	BuildAt    string
	BuildBy    string
	RunnningOS string
}

func (info *AppInfo) ShortVersion() string {
	return info.Version
}

func (info *AppInfo) LongVersion() string {
	buf := bytes.NewBuffer(nil)
	fmt.Println(buf, "project:", info.Name)
	fmt.Println(buf, "version:", info.Version)
	fmt.Println(buf, "git commit:", info.GitCommit)
	fmt.Println(buf, "build at:", info.BuildAt)
	fmt.Println(buf, "build by:", info.BuildBy)
	fmt.Fprintln(buf, "Running OS/Arch:", info.RunnningOS)
	return buf.String()
}
