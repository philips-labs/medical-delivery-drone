package actions

//go:generate stringer -type=Action

// Action the action you would like the drone to execute
type Action int

const (
	// Takeoff instruct the drone to take off
	Takeoff Action = iota
	// Land instruct the drone to land
	Land
	// Up instruct the drone to go up
	Up
	// Down instruct the drone to go down
	Down
	// Left instruct the drone to go left
	Left
	// Right instruct the drone to go right
	Right
	// Forward instruct the drone to go forward
	Forward
	// Backward instruct the drone to go backward
	Backward
	// RotateRight instruct the drone to rotate clockwise
	RotateRight
	// RotateLeft instruct the drone to rotate counter clockwise
	RotateLeft
	// Hover instruct the drone to hover
	Hover
	// FrontFlip instruct the drone to perform a Backflip
	FrontFlip
	// Backflip instruct the drone to perform a Frontflip
	Backflip
	// LeftFlip instruct the drone to perform a Leftflip
	LeftFlip
	// RightFlip instruct the drone to perform a Rightflip
	RightFlip
)
