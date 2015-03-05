package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/leafi/blunt/common"
)

var valloc chan int = make(chan int)
var vfree chan int = make(chan int)
var vstop chan struct{} = make(chan struct{})

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

func InitValloc() {

	nextPostRender <- func() {
		gl.GenVertexArrays(1, &vao)
		gl.BindVertexArray(vao)

		gl.GenBuffers(1, &vbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

		vertices = []float32{
			0.0,
			0.5,
			0.5,
			-0.5,
			-0.5,
			-0.5,
		}

		gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.DYNAMIC_DRAW)

		vshader, _ := Asset("assets/sprite.vert")
		fshader, _ := Asset("assets/sprite.frag")
		var err error
		prog, err = common.NewProgram(string(vshader), string(fshader))

		if err != nil {
			panic(err)
		}

		gl.UseProgram(prog)

		attrPos = uint32(gl.GetAttribLocation(prog, gl.Str("position"+"\x00")))
		// vvv captures current VBO (to the VAO for this attribute) vvv
		gl.VertexAttribPointer(attrPos, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))

		gl.EnableVertexAttribArray(attrPos)

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

func RenderValloc() {
	if !ready {
		return
	}

	// start at 0, render 3 verts == 1 tri
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
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
