package camera

import (
	"math"
	"rapidengine/configuration"
	"rapidengine/input"

	"github.com/go-gl/mathgl/mgl32"
)

type Camera3D struct {
	Speed       float32
	Sensitivity float64

	Position  mgl32.Vec3
	UpAxis    mgl32.Vec3
	FrontAxis mgl32.Vec3

	Pitch float32
	Yaw   float32
	Roll  float32

	MouseX float64
	MouseY float64

	MouseLastX float64
	MouseLastY float64

	FirstMouse bool

	View   mgl32.Mat4
	Config *configuration.EngineConfig
}

func NewCamera3D(position mgl32.Vec3, speed float32, config *configuration.EngineConfig) *Camera3D {
	return &Camera3D{
		Position:    position,
		UpAxis:      mgl32.Vec3{0, 1, 0},
		FrontAxis:   mgl32.Vec3{0, 0, -1},
		Speed:       speed,
		Sensitivity: 0.2,
		Yaw:         0,
		Pitch:       0,
		Config:      config,
	}
}

func (camera3D *Camera3D) Look(delta float64) {
	camera3D.View = mgl32.LookAtV(
		camera3D.Position,
		camera3D.Position.Add(camera3D.FrontAxis),
		mgl32.HomogRotate3D(camera3D.Roll, camera3D.FrontAxis).Mul4x1(mgl32.Vec4{0, 1, 0, 1.0}).Vec3(),
	)
}

//  --------------------------------------------------
//  Movement
//  --------------------------------------------------

func (camera3D *Camera3D) DefaultControls(inputs *input.Input) {
	if inputs.Keys["w"] {
		camera3D.MoveForward()
	}
	if inputs.Keys["s"] {
		camera3D.MoveBackward()
	}
	if inputs.Keys["a"] {
		camera3D.MoveLeft()
	}
	if inputs.Keys["d"] {
		camera3D.MoveRight()
	}
	if inputs.Keys["space"] {
		camera3D.MoveUp()
	}
	if inputs.Keys["shift"] {
		camera3D.MoveDown()
	}

	camera3D.ProcessMouse(inputs.MouseX, inputs.MouseY, inputs.LastMouseX, inputs.LastMouseY)

}

func (camera3D *Camera3D) ProcessMouse(mouseX, mouseY, lastMouseX, lastMouseY float64) {
	xOffset := (mouseX - lastMouseX) * camera3D.Sensitivity
	yOffset := (mouseY - lastMouseY) * camera3D.Sensitivity
	camera3D.Yaw += float32(xOffset)
	camera3D.Pitch -= float32(yOffset)
	if camera3D.Pitch > 89 {
		camera3D.Pitch = 89
	}
	if camera3D.Pitch < -89 {
		camera3D.Pitch = -89
	}
	camera3D.FrontAxis = CalculateDirection(camera3D.Pitch, camera3D.Yaw).Normalize()
	camera3D.FrontAxis = mgl32.HomogRotate3D(camera3D.Roll, camera3D.FrontAxis).Mul4x1(camera3D.FrontAxis.Vec4(1.0)).Vec3()
}

func CalculateDirection(pitch, yaw float32) mgl32.Vec3 {
	return mgl32.Vec3{
		float32(math.Cos(float64(mgl32.DegToRad(pitch))) * math.Cos(float64(mgl32.DegToRad(yaw)))),
		float32(math.Sin(float64(mgl32.DegToRad(pitch)))),
		float32(math.Cos(float64(mgl32.DegToRad(pitch))) * math.Sin(float64(mgl32.DegToRad(yaw)))),
	}
}

func (camera3D *Camera3D) MoveForward() {
	camera3D.Position = camera3D.Position.Add(camera3D.FrontAxis.Mul(camera3D.Speed))
}

func (camera3D *Camera3D) MoveBackward() {
	camera3D.Position = camera3D.Position.Sub(camera3D.FrontAxis.Mul(camera3D.Speed))
}

func (camera3D *Camera3D) MoveUp() {
	camera3D.Position = camera3D.Position.Add(camera3D.UpAxis.Mul(camera3D.Speed))
}

func (camera3D *Camera3D) MoveDown() {
	camera3D.Position = camera3D.Position.Sub(camera3D.UpAxis.Mul(camera3D.Speed))
}

func (camera3D *Camera3D) MoveLeft() {
	camera3D.Position = camera3D.Position.Sub(camera3D.FrontAxis.Cross(camera3D.UpAxis).Normalize().Mul(camera3D.Speed))
}

func (camera3D *Camera3D) MoveRight() {
	camera3D.Position = camera3D.Position.Add(camera3D.FrontAxis.Cross(camera3D.UpAxis).Normalize().Mul(camera3D.Speed))
}

func (camera3D *Camera3D) ChangeRoll(r float32) {
	camera3D.Roll += r
}

//  --------------------------------------------------
//  Setters
//  --------------------------------------------------

func (camera3D *Camera3D) ChangeYaw(y float32) {
	camera3D.Yaw += y
	camera3D.FrontAxis = CalculateDirection(camera3D.Pitch, camera3D.Yaw).Normalize()
}

func (camera3D *Camera3D) ChangePitch(p float32) {
	camera3D.Pitch += p
	camera3D.FrontAxis = CalculateDirection(camera3D.Pitch, camera3D.Yaw).Normalize()
}

func (camera3D *Camera3D) SetPosition(x, y, z float32) {
	camera3D.Position = mgl32.Vec3{x, y, z}
}

func (camera3D *Camera3D) SetSpeed(s float32) {
	camera3D.Speed = s
}

func (camera3D *Camera3D) SetSmoothSpeed(s float32) {
	//camera3D.SmoothSpeed = s
}

func (camera3D *Camera3D) Shake(duration float64, strength float32) {
	//camera3D.SmoothSpeed = s
}

//  --------------------------------------------------
//  Getters
//  --------------------------------------------------

func (camera3D *Camera3D) GetFirstViewIndex() *float32 {
	return &camera3D.View[0]
}

func (camera3D *Camera3D) GetStaticView() mgl32.Mat4 {
	return mgl32.LookAtV(
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{0, 0, 0}.Add(mgl32.Vec3{0, 0, -1}),
		mgl32.Vec3{0, 1, 0},
	)
}

func (camera3D *Camera3D) GetPosition() (float32, float32, float32) {
	return camera3D.Position.X(), camera3D.Position.Y(), camera3D.Position.Z()
}
