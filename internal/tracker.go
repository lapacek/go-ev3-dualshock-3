package internal

import (
	"github.com/ev3go/ev3dev"
	"github.com/sirupsen/logrus"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/joystick"
)

// represents an absolute value of an approximate max speed of a motor
const motorAbsMaxSpeed = 100

// ev3dev related constants
const (
	// inputs
	ev3OutAPortName = "ev3-ports:outA"
	ev3OutBPortName = "ev3-ports:outB"

	// devices
	ev3LargeMotorName = "lego-ev3-l-motor"

	// actions
	ev3BreakActionName = "break"
)

// dualshock3 controller driver related constants
//
// * you can see more here ../config/dualshock3.json
//
const (
	joystickRightYAxisName = "right_y"
	joystickLeftYAxisName = "left_y"
)

type workT func()

type Tracker struct {

	name string
	work workT

	joystick        *joystick.Driver
	joystickAdaptor *joystick.Adaptor
	robot           *gobot.Robot

	// ev3 physical devices
	outA *ev3dev.TachoMotor
	outB *ev3dev.TachoMotor
}

func NewTracker(name string) *Tracker {
	t := Tracker{}
	t.name = name

	return &t
}

func (t *Tracker) Run() {

	if !t.open() {
		logrus.Fatal("Component is not opened.")
	}

	logrus.Debug("Starting...")
	defer logrus.Debug("Stoped")

	t.robot = gobot.NewRobot(
		t.name,
		[]gobot.Connection{t.joystickAdaptor},
		[]gobot.Device{t.joystick},
		t.work,
	)
  
	err := t.robot.Start()
	if err != nil {
		logrus.Errorf("Error occured, err(%v)", err)
	}
}

func (t *Tracker) open() bool {

	logrus.Debug("Opening...")

	if ! t.initMotors() {
		return false
	}

	t.initJoystick()
	t.initWork()

	logrus.Debug("Opened")

	return true
}

func (t *Tracker) initMotors() bool {

	outA, err := ev3dev.TachoMotorFor(ev3OutAPortName, ev3LargeMotorName)
	if err != nil {
		logrus.Errorf("Failed to find right large motor on outA: %v", err)

		return false
	}

	err = outA.SetStopAction(ev3BreakActionName).Err()
	if err != nil {
		logrus.Errorf("Failed to set brake stop for large motor on outA: %v", err)

		return false
	}

	outB, err := ev3dev.TachoMotorFor(ev3OutBPortName, ev3LargeMotorName)
	if err != nil {
		logrus.Errorf("Failed to find left large motor on outB: %v", err)

		return false
	}

	err = outB.SetStopAction(ev3BreakActionName).Err()
	if err != nil {
		logrus.Errorf("Failed to set brake stop for left large motor on outB: %v", err)

		return false
	}

	return true
}

func (t *Tracker) initJoystick() {

	t.joystickAdaptor = joystick.NewAdaptor()
	t.joystick = joystick.NewDriver(t.joystickAdaptor,
		"../config/dualshock3.json",
	)
}

func (t *Tracker) initWork() {

	t.work = func() {

		logrus.Debug("Working...")
		defer logrus.Debug("Work stopped.")

		var err error

		err = t.joystick.On(t.joystick.Event(joystickRightYAxisName), t.handleRightStickAction)
		if err != nil {
			logrus.Errorf("Failed to register a joystick right_y event handler, err(%v)", err)
		}

		err = t.joystick.On(t.joystick.Event(joystickLeftYAxisName), t.handleLeftStickAction)
		if err != nil {
			logrus.Errorf("Failed to register a joystick right_y event handler, err(%v)", err)
		}
	}
}

func (t *Tracker) handleRightStickAction(data interface{}) {

	logrus.Tracef("Joystick event received, right y(%v)", data)

	handleStickEvent(t.outB, data)
}

func (t *Tracker) handleLeftStickAction(data interface{}) {

	logrus.Tracef("Joystick event received, right y(%v)", data)

	handleStickEvent(t.outA, data)
}

func (t *Tracker) close() {

	err := t.robot.Stop()
	if err != nil {
		logrus.Errorf("Failed to stop a robot, err(%v)", err)
	}
}