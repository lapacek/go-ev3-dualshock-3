package main

import (
	"fmt"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/joystick"
)

func main() {
	joystickAdaptor := joystick.NewAdaptor()
	joystick := joystick.NewDriver(joystickAdaptor,
		//"./platforms/joystick/configs/dualshock3.json",
		"./dualshock3.json",
	)

	work := func() {
		joystick.On(joystick.Event("square_press"), func(data interface{}) {
			fmt.Println("square_press")
		})
		joystick.On(joystick.Event("square_release"), func(data interface{}) {
			fmt.Println("square_release")
		})
		joystick.On(joystick.Event("triangle_press"), func(data interface{}) {
			fmt.Println("triangle_press")
		})
		joystick.On(joystick.Event("triangle_release"), func(data interface{}) {
			fmt.Println("triangle_release")
		})
		joystick.On(joystick.Event("left_x"), func(data interface{}) {
			fmt.Println("left_x", data)
		})
		joystick.On(joystick.Event("left_y"), func(data interface{}) {
			fmt.Println("left_y", data)
		})
		joystick.On(joystick.Event("right_x"), func(data interface{}) {
			fmt.Println("right_x", data)
		})
		joystick.On(joystick.Event("right_y"), func(data interface{}) {
			fmt.Println("right_y", data)
		})
	}

	robot := gobot.NewRobot("joystickBot",
		[]gobot.Connection{joystickAdaptor},
		[]gobot.Device{joystick},
		work,
	)

	robot.Start()
}
