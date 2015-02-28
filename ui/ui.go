// An OpenGL user interface library.
package ui

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/leafi/blunt/common"
)

type Element struct {
	// Click is the channel to communicate on onClick, or nil if not clickable.
	Click chan int
	// Should this be rendered with an opaque backing rectangle?
	Pane bool
	// Will this be an OpenGL FBO? (Likely to be a 4:3 square.)
	RSquare bool
	// Text to display, or nil.
	Text string
	// Is this an X or a v button or something?
	Small bool
}

type Container struct {
	Element
	Ctrls []*Element
}

type WContainer struct {
	Container
}

type HContainer struct {
	Container
}

type Window struct {
	Container
	Start mgl32.Vec2
	End   mgl32.Vec2
}

type Ui struct {
	// OpenGL buffers & shaders

	// current state
	RectStart []mgl32.Vec2
	RectEnd   []mgl32.Vec2

	// state for generatin'
	Windows []*Window
}

var u Ui

/*func Get() *Ui {
	return &u
}*/

// Tells the UI module it needs to rebuild.
func Notify() {

}

func GLInit() {
	var vaoBack uint32
	gl.GenVertexArrays(1, &vaoBack)
	gl.BindVertexArray(vaoBack)

	var vboBack uint32
	gl.GenBuffers(1, &vboBack)
	gl.BindBuffer(gl.ARRAY_BUFFER, vboBack)
	//gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*4, gl.Ptr(cubeVertices), gl.DYNAMIC_DRAW)

	/*	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
		gl.EnableVertexAttribArray(vertAttrib)
		gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))*/

}

func GLRender() {
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

}

// rebuildRects rebuilds layout rectangles & does depth calculations.
func rebuildRects() {

}

// rebuildGL rebuilds OpenGL state based on RectStart, RectEnd.
// Call after rebuildRects.
func rebuildGL() {

}

/*func GiveEventCh() {

}*/
