package drone

import (
	"github.com/philips-labs/medical-delivery-drone/drone/actions"
	"gobot.io/x/gobot/platforms/dji/tello"
)

// Move moves the drone using the given action and speed
type Move struct {
	Action actions.Action
	Speed  int
}

var (
	TakeoffMove   = Move{Action: actions.Takeoff}
	LandingMove   = Move{Action: actions.Land}
	BackFlipMove  = Move{Action: actions.Backflip}
	FrontFlipMove = Move{Action: actions.FrontFlip}
	LeftFlipMove  = Move{Action: actions.LeftFlip}
	RightFlipMove = Move{Action: actions.RightFlip}
)

// Mover moves the drone using the given Move
type Mover func(driver *tello.Driver, move Move)

func performMove(drone *tello.Driver, move Move) {
	switch move.Action {
	case actions.Takeoff:
		_ = drone.TakeOff()
	case actions.Land:
		_ = drone.Land()
	case actions.Up:
		_ = drone.Up(move.Speed)
	case actions.Down:
		_ = drone.Down(move.Speed)
	case actions.Left:
		_ = drone.Left(move.Speed)
	case actions.Right:
		_ = drone.Right(move.Speed)
	case actions.Forward:
		_ = drone.Forward(move.Speed)
	case actions.Backward:
		_ = drone.Backward(move.Speed)
	case actions.RotateRight:
		_ = drone.Clockwise(move.Speed)
	case actions.RotateLeft:
		_ = drone.CounterClockwise(move.Speed)
	case actions.Backflip:
		_ = drone.BackFlip()
	case actions.FrontFlip:
		_ = drone.FrontFlip()
	case actions.LeftFlip:
		_ = drone.LeftFlip()
	case actions.RightFlip:
		_ = drone.RightFlip()
	case actions.Hover:
		drone.Hover()
	}
}
