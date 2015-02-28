// blunt project main.go
package main

//go:generate go get github.com/jteeuwen/go-bindata/...
//go:generate $GOPATH/bin/go-bindata -pkg $GOPACKAGE assets/

import (
	"fmt"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	//"github.com/go-gl/mathgl/mgl32"
)

var (
	framebufferSizeCh    = make(chan int, 2)
	framebufferSizeDirty bool
	framebufferWidth     int
	framebufferHeight    int
)

func init() {
	// must perform certain functions (e.g. window creation, maintenance) on main thread
	runtime.LockOSThread()
}

func tryPushFramebufferSize() {
	select {
	case framebufferSizeCh <- framebufferWidth:
		framebufferSizeCh <- framebufferHeight
		framebufferSizeDirty = false
	default:
		// Nope, channel's still full.
	}
}

func main() {

	if err := glfw.Init(); err != nil {
		panic(fmt.Sprintf("failed to initialize glfw: %v", err))
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(640, 400, "Cube", nil, nil)

	if err != nil {
		panic(err)
	}

	framebufferWidth, framebufferHeight = window.GetFramebufferSize()
	tryPushFramebufferSize()

	window.SetFramebufferSizeCallback(func(w *glfw.Window, width int, height int) {
		framebufferWidth = width
		framebufferHeight = height
		framebufferSizeDirty = true
		tryPushFramebufferSize()
	})

	go func() {
		// GL stuff must stay on same thread. (May or may not be same thread as
		// GLFW main thread-only stuff. Probably will be.)
		runtime.LockOSThread()

		// do opengl stuff here
		window.MakeContextCurrent()

		// Initialize Glow
		if err := gl.Init(); err != nil {
			panic(err)
		}

		version := gl.GoStr(gl.GetString(gl.VERSION))
		fmt.Println("OpenGL version", version)

		fw, fh := 640, 400

		for !window.ShouldClose() {
			select {
			case fw = <-framebufferSizeCh:
				fh = <-framebufferSizeCh
				// TODO: viewport shit
				fmt.Printf("got %v x %v\n", fw, fh)
			default:
				// No resize events.
			}

			if fw+fh == 1 {
				fmt.Println("okay")
			}

			gl.ClearColor(0.0, 1.0, 1.0, 1.0)
			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
			window.SwapBuffers()
		}
	}()

	for !window.ShouldClose() {
		glfw.PollEvents()

		if framebufferSizeDirty {
			tryPushFramebufferSize()
		}

		time.Sleep(5 * time.Millisecond)
	}

}
