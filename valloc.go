package main

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/leafi/blunt/common"
)

var valloc chan int = make(chan int)
var vfree chan int = make(chan int)
var vstop chan bool = make(chan bool)

const highv int = 256000

// vvv what should the limit be? vvv
var nextPostRender chan func() = make(chan func(), 256)

var (
	vao      uint32
	vbo      uint32
	vertices []float32
	prog     uint32
	attrPos  uint32
	ready    bool
)

var (
	propbo        uint32
	props         []Prop
	propsToUpdate []int
)

type int2 struct {
	x int
	y int
}

type Prop struct {
	position mgl32.Vec2
	size     mgl32.Vec2
	// --- std140 ---
	scale float32
	texUV int2
	angle float32
	// --- std140 ---
	tint mgl32.Vec4
	// --- std140 ---
}

func InitValloc() {

	nextPostRender <- func() {
		vshader, _ := Asset("assets/sprite.vert")
		fshader, _ := Asset("assets/sprite.frag")
		var err error
		prog, err = common.NewProgram(string(vshader), string(fshader))

		if err != nil {
			panic(err)
		}

		gl.UseProgram(prog)

		gl.GenVertexArrays(1, &vao)
		gl.BindVertexArray(vao)

		gl.GenBuffers(1, &vbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

		// origin is bottom-left rather than centre b/c we'll always have to fiddle
		// with the offsets on a sprite-by-sprite basis anyway.
		// this way, at least, doing things like rendering regular sized tiles &
		// ui elements will be trivial.
		vertices = []float32{
			0.0,
			1.0,
			0.0,
			0.0,
			1.0,
			0.0,

			1.0,
			0.0,
			1.0,
			1.0,
			0.0,
			1.0,
		}

		gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

		attrPos = uint32(gl.GetAttribLocation(prog, gl.Str("basePosition"+"\x00")))
		// vvv captures current VBO (to the VAO for this attribute) vvv
		gl.VertexAttribPointer(attrPos, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))

		gl.EnableVertexAttribArray(attrPos)

		props := []Prop{
			Prop{
				position: mgl32.Vec2{-0.5, 0.0},
				size:     mgl32.Vec2{1.0, 1.0},
				scale:    0.4,
				tint:     mgl32.Vec4{1.0, 1.0, 1.0, 1.0},
			},
		}

		gl.GenBuffers(1, &propbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, propbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(props)*48, gl.Ptr(props), gl.STREAM_DRAW)

		instanceAttrNames := []string{"position", "size", "scale", "texU", "texV", "angle", "tint"}
		instanceAttrTypes := []uint32{gl.FLOAT, gl.FLOAT, gl.FLOAT, gl.INT, gl.INT, gl.FLOAT, gl.FLOAT}
		instanceAttrMul := []int32{2, 2, 1, 1, 1, 1, 4}

		offset := 0
		attr := uint32(0)
		for i, name := range instanceAttrNames {
			attr = uint32(gl.GetAttribLocation(prog, gl.Str(name+"\x00")))
			fmt.Printf("%v %v (%v)\n", i, name, attr)
			gl.VertexAttribPointer(attr, instanceAttrMul[i], instanceAttrTypes[i], false, 48, gl.PtrOffset(offset))
			gl.EnableVertexAttribArray(attr)
			gl.VertexAttribDivisor(attr, 1)

			// gl.FLOAT, gl.INT are always 4 bytes. so that's easy-ish.
			offset += 4 * int(instanceAttrMul[i])
			attr++
		}

		ready = true
	}

	// texture atlas building:
	// https://www.opengl.org/wiki/Framebuffer#Blitting

	go func() {
		free := make([]int, 1) // free[0] == 0
		high := 1

		f := func() bool {
			select {
			case valloc <- free[0]:
				free = free[1:]
				if len(free) < 1 {
					free = append(free, high)
					high++
				}
				return true
			case x := <-vfree:
				free = append(free, x)
				return true
			case <-vstop:
				return false
			}
		}

		for f() {
		}
	}()
}

/*void FillUniformBuffer() {
	GL.BindBuffer(BufferTarget.UniformBuffer, BufferUBO);
	GL.BufferSubData(BufferTarget.UniformBuffer, (IntPtr)0, (IntPtr)(sizeof(float) * 8), ref UBOData);
	GL.BindBuffer(BufferTarget.UniformBuffer, 0);
}*/

func RenderValloc() {
	if !ready {
		return
	}

	// start at 0, render 3 verts == 1 tri
	//gl.DrawArrays(gl.TRIANGLES, 0, 6)
	gl.DrawArraysInstanced(gl.TRIANGLES, 0, 6, 1)
}

// called from GL thread & coroutine
func SpinValloc() {
	for {
		select {
		case nf := <-nextPostRender:
			nf()
		default:
			return
		}
	}
}
