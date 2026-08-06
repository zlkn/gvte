package ui

import (
	"log"
	"runtime"

	// "gvte/internal/config"
	// "gvte/internal/emulator"
	// "gvte/internal/ui/font"
	// "gvte/internal/ui/input"
	// "gvte/internal/ui/renderer"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

func init() {
	// GLFW requires the main thread to safely interact with the OS window manager
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalf("Failed to initialize GLFW: %v", err)
	}
	defer glfw.Terminate()

	// 1. Tell GLFW NOT to create an OpenGL context. We are using WebGPU!
	glfw.WindowHint(glfw.ClientAPI, glfw.NoAPI)

	// Create the window
	window, err := glfw.CreateWindow(800, 600, "WebGPU Terminal", nil, nil)
	if err != nil {
		log.Fatalf("Failed to create window: %v", err)
	}

	// 2. Initialize WebGPU Instance
	instance := wgpu.CreateInstance(&wgpu.InstanceDescriptor{
		// Optional: you can specify backends like Vulkan, Metal, DX12 explicitly
	})
	defer instance.Release()

	// 3. THE BRIDGE: Create a WebGPU surface from the GLFW window
	// The rajveermalviya package abstracts away the nasty OS-specific handle extraction
	surface := instance.CreateSurface(window.GetX11Window(), window.GetX11Display())
	// Note: In actual implementation, you use platform-specific getters based on GOOS,
	// e.g., window.GetWin32Window() or window.GetCocoaWindow().

	defer surface.Release()

	// 4. Request the GPU adapter and device to start rendering
	adapter, err := instance.RequestAdapter(&wgpu.RequestAdapterOptions{
		CompatibleSurface: surface,
	})
	if err != nil {
		log.Fatalf("Failed to get adapter: %v", err)
	}
	defer adapter.Release()

	// Setup your Render Loop here using glfw.PollEvents()
	for !window.ShouldClose() {
		glfw.PollEvents()
		// RenderFrame(...)
	}
}
