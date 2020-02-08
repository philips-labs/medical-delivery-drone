package drone

import (
	"context"
	"log"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"

	"github.com/philips-labs/medical-delivery-drone/video"
)

// Connect connects to the drone and starts the video
func Connect(ctx context.Context, converter *video.Converter) (<-chan []byte, error) {
	drone := tello.NewDriver("8890")

	initDrone := func() {
		_ = drone.On(tello.ConnectedEvent, func(data interface{}) {
			log.Println("Connected to Tello.")
			_ = drone.StartVideo()
			_ = drone.SetVideoEncoderRate(tello.VideoBitRateAuto)
			_ = drone.SetExposure(0)

			gobot.Every(30*time.Millisecond, func() {
				_ = drone.StartVideo()
			})
		})

		drone.On(tello.VideoFrameEvent, func(data interface{}) {
			pkt := data.([]byte)
			if _, err := converter.Write(pkt); err != nil {
				log.Println(err)
			}
		})
	}

	robot := gobot.NewRobot("tello",
		[]gobot.Connection{},
		[]gobot.Device{drone},
		initDrone,
	)

	// calling Start(false) lets the Start routine return immediately without an additional blocking goroutine
	robot.Start(false)

	return converter.Start(ctx)
}
