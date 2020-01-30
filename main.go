package main

import (
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
)

func main() {
	drone := tello.NewDriver("8888")
	work := func() {
		drone.TakeOff()
		<-time.After(2 * time.Second)
		drone.Up(15)
		<-time.After(2 * time.Second)
		drone.BackFlip()
		<-time.After(2 * time.Second)
		drone.Forward(5)
		<-time.After(2 * time.Second)
		drone.BackFlip()
		<-time.After(2 * time.Second)
		drone.Land()
	}
	robot := gobot.NewRobot("tello",
		[]gobot.Connection{},
		[]gobot.Device{drone},
		work,
	)
	robot.Start()
}
