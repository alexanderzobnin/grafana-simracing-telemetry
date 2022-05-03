//go:build mage
// +build mage

package main

import (
	"fmt"
	// mage:import
	build "github.com/grafana/grafana-plugin-sdk-go/build"
	"github.com/magefile/mage/mg"
)

// Hello prints a message (shows that you can define custom Mage targets).
func Hello() {
	fmt.Println("hello plugin developer!")
}

func BuildWindows() { //revive:disable-line
	b := build.Build{}
	mg.Deps(b.Windows)
}

// Default configures the default target.
var Default = BuildWindows
