package main

import (
	"fmt"
	"multitenant/app"
)

const (
	VERSION = "0.0.1"
)

func main() {
	// defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	// defer pprof.StopCPUProfile()
	fmt.Println("Starting Multi Tenant service version", VERSION)
	app.Start()
}
