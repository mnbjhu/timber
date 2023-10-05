package main

import (
	"time"

	"github.com/labstack/gommon/log"
)

func main() {
	a := 1
	for {
		log.Infof("My value = %d", a)
		log.Warnf("My value = %d", a)
		log.Errorf("My value = %d", a)
		a++
		time.Sleep(1 * time.Second)
	}
}
