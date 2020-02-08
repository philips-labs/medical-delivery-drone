package drone

import (
	"log"
	"math"
	"os"
	"os/signal"

	"github.com/kpeu3i/gods4"
	"github.com/kpeu3i/gods4/led"
	"github.com/kpeu3i/gods4/rumble"
	"gobot.io/x/gobot/platforms/dji/tello"

	"github.com/philips-labs/medical-delivery-drone/drone/actions"
)

const (
	stickRadius = 127
	moveSpeed   = 60
	upDownSpeed = 60
	turnSpeed   = 90
	offsetError = 30
)

func startController(control Mover, drone *tello.Driver) {
	// Find all controllers connected to your machine via USB or Bluetooth
	controllers := gods4.Find()
	if len(controllers) == 0 {
		panic("No connected DS4 controllers found")
	}

	// Select first controller from the list
	controller := controllers[0]

	// Connect to the controller
	err := controller.Connect()
	if err != nil {
		panic(err)
	}

	log.Printf("* Controller #1 | %-10s | Action: actions.%s, connection: %s\n", "Connect", controller, controller.ConnectionType())

	// Disconnect controller when a program is terminated
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		err = drone.Land()
		if err != nil {
			panic(err)
		}
		err = controller.Disconnect()
		if err != nil {
			panic(err)
		}

		log.Printf("* Controller #1 | %-10s | bye!\n", "Disconnect")
	}()

	// Register callback for "BatteryUpdate" event
	controller.On(gods4.EventBatteryUpdate, func(data interface{}) error {
		battery := data.(gods4.Battery)
		log.Printf("* Controller #1 | %-10s | capacity: %v%%, charging: %v, cable: %v\n",
			"Battery",
			battery.Capacity,
			battery.IsCharging,
			battery.IsCableConnected,
		)

		return nil
	})

	// R1 for take off
	controller.On(gods4.EventR1Press, func(data interface{}) error {
		log.Printf("* Controller #1 | %-10s | state: press\n", "R1")
		control(drone, TakeoffMove)
		return nil
	})

	// L1 for land
	controller.On(gods4.EventL1Release, func(data interface{}) error {
		log.Printf("* Controller #1 | %-10s | state: press\n", "L1")
		control(drone, LandingMove)
		return nil
	})

	// R2 for up
	var R2Pressed bool
	controller.On(gods4.EventR2Press, func(data interface{}) error {
		log.Printf("* Controller #1 | %-10s | state: press\n", "R2")
		if !R2Pressed {
			R2Pressed = !R2Pressed
			control(drone, Move{Action: actions.Up, Speed: upDownSpeed})
			log.Printf("go UP value %d ", upDownSpeed)
		}
		return nil
	})

	controller.On(gods4.EventR2Release, func(data interface{}) error {
		log.Printf("* Controller #1 | %-10s | state: release\n", "R2")
		R2Pressed = !R2Pressed
		control(drone, Move{Action: actions.Up})
		log.Printf("go UP value %d ", 0)
		return nil
	})

	// L2 for down
	var L2Pressed bool
	controller.On(gods4.EventL2Press, func(data interface{}) error {
		log.Printf("* Controller #1 | %-10s | state: press\n", "L2")
		if !L2Pressed {
			L2Pressed = !L2Pressed
			control(drone, Move{Action: actions.Down, Speed: upDownSpeed})
			log.Printf("go down value %d ", upDownSpeed)
		}
		return nil
	})

	// L2 for down
	controller.On(gods4.EventL2Release, func(data interface{}) error {
		log.Printf("* Controller #1 | %-10s | state: release\n", "L2")
		L2Pressed = false
		control(drone, Move{Action: actions.Down})
		log.Printf("go down value %d ", 0)
		return nil
	})

	// Register callback for "RightStickMove" event
	controller.On(gods4.EventRightStickMove, func(data interface{}) error {
		stick := data.(gods4.Stick)
		log.Printf("* Controller #1 | %-10s | x: %v, y: %v\n", "RightStick", stick.X, stick.Y)

		moveX := (int(stick.X) - stickRadius) * turnSpeed / stickRadius
		if int(stick.X) <= stickRadius {
			control(drone, Move{Action: actions.RotateLeft, Speed: int(math.Abs(float64(moveX)))})
		} else {
			control(drone, Move{Action: actions.RotateRight, Speed: int(math.Abs(float64(moveX)))})
		}
		return nil
	})
	controller.On(gods4.EventLeftStickMove, func(data interface{}) error {
		stick := data.(gods4.Stick)
		log.Printf("* Controller #1 | %-10s | x: %v, y: %v\n", "LeftStick", stick.X, stick.Y)

		if int(stick.X) > stickRadius-offsetError &&
			int(stick.X) < stickRadius+offsetError &&
			int(stick.Y) > stickRadius-offsetError &&
			int(stick.Y) < stickRadius+offsetError {
			control(drone, Move{Action: actions.Hover})
			return nil
		}
		// do nothing
		moveX := (int(stick.X) - stickRadius) * moveSpeed / stickRadius
		moveY := (int(stick.Y) - stickRadius) * moveSpeed / stickRadius

		absX := int(math.Abs(float64(moveX)))
		absY := int(math.Abs(float64(moveY)))
		var action actions.Action
		if absX > absY {
			if moveX < 0 {
				action = actions.Left
			} else {
				action = actions.Right
			}
			control(drone, Move{Action: action, Speed: absX})
		} else {
			if moveY < 0 {
				action = actions.Forward
			} else {
				action = actions.Backward
			}
			control(drone, Move{Action: action, Speed: absY})
		}
		return nil
	})

	// Trick controls
	controller.On(gods4.EventDPadUpPress, func(data interface{}) error {
		log.Printf("* Controller #1 | %-10s", "Up Press")
		control(drone, FrontFlipMove)
		return nil
	})
	controller.On(gods4.EventDPadDownPress, func(data interface{}) error {
		log.Printf("* Controller #1 | %-10s", "Down Press")
		control(drone, BackFlipMove)
		return nil
	})
	controller.On(gods4.EventDPadLeftPress, func(data interface{}) error {
		log.Printf("* Controller #1 | %-10s", "Left Press")
		control(drone, LeftFlipMove)

		return nil
	})
	controller.On(gods4.EventDPadRightPress, func(data interface{}) error {
		log.Printf("* Controller #1 | %-10s", "Right Press")
		control(drone, RightFlipMove)

		return nil
	})

	// Enable left and right rumble motors
	err = controller.Rumble(rumble.Both())
	if err != nil {
		panic(err)
	}

	// Enable LED (yellow) with flash
	err = controller.Led(led.Yellow().Flash(50, 50))
	if err != nil {
		panic(err)
	}

	// Start listening for controller events
	err = controller.Listen()
	if err != nil {
		panic(err)
	}
}
