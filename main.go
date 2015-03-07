// blunt project main.go
package main

//go:generate go get github.com/jteeuwen/go-bindata/...
//go:generate $GOPATH/bin/go-bindata -pkg $GOPACKAGE assets/

import (
	"fmt"
	"runtime"
	"time"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	//"github.com/go-gl/mathgl/mgl32"
)

var (
	framebufferSizeCh    = make(chan int, 2)
	framebufferSizeDirty bool
	framebufferWidth     int
	framebufferHeight    int
)

var (
	AddGameLoopCh  = make(chan chan time.Duration)
	RmGameLoopCh   = make(chan chan time.Duration)
	swapBuffersCh  = make(chan bool)
	stopGameLoopCh = make(chan bool)
)

func init() {
	// must perform certain functions (e.g. window creation, maintenance) on main thread
	runtime.LockOSThread()
}

func gameLoopHandler() {
	chans := make([]chan time.Duration, 0)

	// TODO: better timekeeping! linux is probably alright but the others aren't >60FPS

	lastTime := time.Now()

	for {

		break2 := false
		for !break2 {
			select {
			case ch := <-AddGameLoopCh:
				chans = append(chans, ch)
			case ch := <-RmGameLoopCh:
				for i, ch2 := range chans {
					if ch2 == ch {
						chans = append(chans[:i], chans[i+1:]...)
					}
				}
			case <-swapBuffersCh:
				break2 = true
			}
		}

		// TODO: ensure t > 0 and < (reasonable maximum)!
		nextTime := time.Now()
		t := nextTime.Sub(lastTime)
		for _, ch := range chans {
			ch <- t
		}
		for _, ch := range chans {
			// maybe we should report channels that take too long to respond?
			<-ch
		}
		swapBuffersCh <- true
		lastTime = nextTime
	}
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
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	//glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

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

		go gameLoopHandler()

		// set up render layers
		InitValloc()

		for !window.ShouldClose() {
			select {
			case fw = <-framebufferSizeCh:
				fh = <-framebufferSizeCh
				// TODO: viewport shit
				fmt.Printf("got %v x %v\n", fw, fh)
				gl.Viewport(0, 0, int32(fw), int32(fh))
			default:
				// No resize events.
			}

			if fw+fh == 1 {
				fmt.Println("okay")
			}

			gl.ClearColor(0.0, 1.0, 1.0, 1.0)
			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

			// sprite render layer
			RenderValloc()

			// post-render tasks for layers (updating buffers probably)
			SpinValloc()

			window.SwapBuffers()
			swapBuffersCh <- true
			<-swapBuffersCh
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
