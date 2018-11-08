package child

// --------------------------------------------------
// Child2D.go contains Child2D, the basic Object in
// rapidengine. Every game object is either a Child,
// or a copy of a Child.
// --------------------------------------------------

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"rapidengine/camera"
	"rapidengine/configuration"
	"rapidengine/geometry"
	"rapidengine/material"
	"rapidengine/physics"
)

type Child2D struct {
	vertexArray *geometry.VertexArray
	numVertices int32

	primitive string

	shaderProgram uint32
	material      *material.Material

	modelMatrix      mgl32.Mat4
	projectionMatrix mgl32.Mat4

	numCopies      int
	copies         []ChildCopy
	currentCopies  []ChildCopy
	copyingEnabled bool

	specificRenderDistance float32

	X float32
	Y float32

	VX float32
	VY float32

	Gravity float32

	Group          string
	collider       physics.Collider
	mouseCollision func(bool)

	config *configuration.EngineConfig
	//collisioncontrol *CollisionControl
}

func NewChild2D(config *configuration.EngineConfig) Child2D {
	return Child2D{
		modelMatrix:            mgl32.Ident4(),
		projectionMatrix:       mgl32.Ortho2D(-1, 1, -1, 1),
		config:                 config,
		X:                      0,
		Y:                      0,
		VX:                     0,
		VY:                     0,
		Gravity:                0,
		copyingEnabled:         false,
		specificRenderDistance: 0,
	}
}

func (child2D *Child2D) PreRender(mainCamera camera.Camera) {
	child2D.BindChild()

	gl.UniformMatrix4fv(
		gl.GetUniformLocation(child2D.shaderProgram, gl.Str("modelMtx\x00")),
		1, false, &child2D.modelMatrix[0],
	)

	gl.UniformMatrix4fv(
		gl.GetUniformLocation(child2D.shaderProgram, gl.Str("viewMtx\x00")),
		1, false, mainCamera.GetFirstViewIndex(),
	)

	gl.UniformMatrix4fv(
		gl.GetUniformLocation(child2D.shaderProgram, gl.Str("projectionMtx\x00")),
		1, false, &child2D.projectionMatrix[0],
	)

	gl.BindAttribLocation(child2D.shaderProgram, 0, gl.Str("position\x00"))
	gl.BindAttribLocation(child2D.shaderProgram, 1, gl.Str("tex\x00"))

	gl.BindVertexArray(0)
}

func (child2D *Child2D) BindChild() {
	gl.BindVertexArray(child2D.vertexArray.GetID())
	gl.UseProgram(child2D.shaderProgram)
}

func (child2D *Child2D) Update(mainCamera camera.Camera, delta float64, lastFrame float64) {
	//cx, cy, _ := mainCamera.GetPosition()
	child2D.VY -= child2D.Gravity

	/*cols := child2D.collisioncontrol.CheckCollisionWithGroup(child2D, "ground", cx, cy)
	if (cols[3] && child2D.VY < 0) || (cols[1] && child2D.VY > 0) {
		child2D.VY = 0
	}
	if (cols[0] && child2D.VX < 0) || (cols[2] && child2D.VX > 0) {
		child2D.VX = 0
	}*/

	child2D.X += child2D.VX * -float32(delta*30)
	child2D.Y += child2D.VY * float32(delta*30)

	child2D.Render(mainCamera, delta)
}

func (child2D *Child2D) Render(mainCamera camera.Camera, delta float64) {
	sX, sY := ScaleCoordinates(child2D.X, child2D.Y, float32(child2D.config.ScreenWidth), float32(child2D.config.ScreenHeight))
	child2D.modelMatrix = mgl32.Translate3D(sX, sY, 0)

	gl.UniformMatrix4fv(
		gl.GetUniformLocation(child2D.shaderProgram, gl.Str("viewMtx\x00")),
		1, false, mainCamera.GetFirstViewIndex(),
	)

	gl.UniformMatrix4fv(
		gl.GetUniformLocation(child2D.shaderProgram, gl.Str("modelMtx\x00")),
		1, false, &child2D.modelMatrix[0],
	)

	child2D.material.Render(delta, 1)
}

func (child2D *Child2D) RenderCopy(config ChildCopy, mainCamera camera.Camera) {
	sX, sY := ScaleCoordinates(config.X, config.Y, float32(child2D.config.ScreenWidth), float32(child2D.config.ScreenHeight))
	child2D.modelMatrix = mgl32.Translate3D(sX, sY, 0)
	child2D.projectionMatrix = mgl32.Ortho2D(-1, 1, -1, 1)

	gl.UniformMatrix4fv(
		gl.GetUniformLocation(child2D.shaderProgram, gl.Str("viewMtx\x00")),
		1, false, mainCamera.GetFirstViewIndex(),
	)

	gl.UniformMatrix4fv(
		gl.GetUniformLocation(child2D.shaderProgram, gl.Str("modelMtx\x00")),
		1, false, &child2D.modelMatrix[0],
	)

	config.Material.Render(0, config.Darkness)
}

func (child2D *Child2D) CheckCollision(other Child) int {
	return child2D.collider.CheckCollision(child2D.X, child2D.Y, child2D.VX, child2D.VY, other.GetX(), other.GetY(), other.GetCollider())
}

func (child2D *Child2D) CheckCollisionRaw(otherX, otherY float32, otherCollider *physics.Collider) int {
	return child2D.collider.CheckCollision(child2D.X, child2D.Y, child2D.VX, child2D.VY, otherX, otherY, otherCollider)
}

//  --------------------------------------------------
//  Component Attachment
//  --------------------------------------------------

func (child2D *Child2D) AttachTextureCoords(coords []float32) {
	if child2D.vertexArray == nil {
		panic("Cannot attach texture without VertexArray")
	}
	if child2D.shaderProgram == 0 {
		panic("Cannot attach texture without shader program")
	}

	gl.BindVertexArray(child2D.vertexArray.GetID())
	gl.UseProgram(child2D.shaderProgram)
	child2D.vertexArray.AddVertexAttribute(coords, 1, 2)
	gl.BindVertexArray(0)
}

func (child2D *Child2D) AttachTextureCoordsPrimitive() {
	child2D.AttachTextureCoords(geometry.GetPrimitiveCoords(child2D.primitive))
}

func (child2D *Child2D) AttachCollider(x, y, w, h float32) {
	child2D.collider = physics.NewCollider(x, y, w, h)
}

func (child2D *Child2D) AttachVertexArray(vao *geometry.VertexArray, numVertices int32) {
	child2D.vertexArray = vao
	child2D.numVertices = numVertices
}

func (child2D *Child2D) AttachPrimitive(p geometry.Primitive) {
	child2D.primitive = p.GetID()
	child2D.AttachVertexArray(p.GetVAO(), p.GetNumVertices())
	child2D.vertexArray.AddVertexAttribute(geometry.RectNormals, 2, 3)
}

func (child2D *Child2D) AttachMaterial(m *material.Material) {
	child2D.material = m
}

func (child2D *Child2D) AttachShader(s uint32) {
	child2D.shaderProgram = s
}

func (child2D *Child2D) AttachGroup(group string) {
	child2D.Group = group
}

//  --------------------------------------------------
//  Setters
//  --------------------------------------------------

func (child2D *Child2D) SetVelocity(vx, vy float32) {
	child2D.VX = vx
	child2D.VY = vy
}

func (child2D *Child2D) SetPosition(x, y float32) {
	child2D.X = x
	child2D.Y = y
}

func (child2D *Child2D) SetSpecificRenderDistance(d float32) {
	child2D.specificRenderDistance = d
}

//  --------------------------------------------------
//  Getters
//  --------------------------------------------------

func (child2D *Child2D) GetShaderProgram() uint32 {
	return child2D.shaderProgram
}

func (child2D *Child2D) GetVertexArray() *geometry.VertexArray {
	return child2D.vertexArray
}

func (child2D *Child2D) GetNumVertices() int32 {
	return child2D.numVertices
}

func (child2D *Child2D) GetTexture() *uint32 {
	return child2D.material.GetTexture()
}

func (child2D *Child2D) GetCollider() *physics.Collider {
	return &child2D.collider
}

func (child2D *Child2D) GetX() float32 {
	return child2D.X
}

func (child2D *Child2D) GetY() float32 {
	return child2D.Y
}

func (child2D *Child2D) GetSpecificRenderDistance() float32 {
	return child2D.specificRenderDistance
}

//  --------------------------------------------------
//  Copying
//  --------------------------------------------------

func (child2D *Child2D) EnableCopying() {
	child2D.copyingEnabled = true
}

func (child2D *Child2D) DisableCopying() {
	child2D.copyingEnabled = false
}

func (child2D *Child2D) AddCopy(config ChildCopy) {
	child2D.numCopies += 1
	child2D.copies = append(child2D.copies, config)
}

func (child2D *Child2D) GetCopies() *[]ChildCopy {
	return &child2D.copies
}

func (child2D *Child2D) IterCopies(f func(Child, ChildCopy)) {
	for _, copy := range child2D.copies {
		f(child2D, copy)
	}
}

func (child2D *Child2D) GetNumCopies() int {
	return child2D.numCopies
}

func (child2D *Child2D) CheckCopyingEnabled() bool {
	return child2D.copyingEnabled
}

func (child2D *Child2D) AddCurrentCopy(c ChildCopy) {
	child2D.currentCopies = append(child2D.currentCopies, c)
}

func (child2D *Child2D) RemoveCurrentCopies() {
	child2D.currentCopies = []ChildCopy{}
}

func (child2D *Child2D) GetCurrentCopies() []ChildCopy {
	return child2D.currentCopies
}

//  --------------------------------------------------
//  Mouse Collision
//  --------------------------------------------------

func (child2D *Child2D) SetMouseFunc(r func(bool)) {
	child2D.mouseCollision = r
}

func (child2D *Child2D) MouseCollisionFunc(c bool) {
	child2D.mouseCollision(c)
}

func ScaleCoordinates(x, y, sw, sh float32) (float32, float32) {
	return 2*(x/float32(sw)) - 1, 2*(y/float32(sh)) - 1
}