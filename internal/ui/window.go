package ui

import (
	"fmt"
	"runtime"

	"gvte/internal/config"
	"gvte/internal/emulator"
	"gvte/internal/ui/font"
	"gvte/internal/ui/input"
	"gvte/internal/ui/renderer"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

func init() {
	// GLFW requires the main thread to safely interact with the OS window manager
	runtime.LockOSThread()
}

type Window struct {
	cfg         *config.Config
	state       *emulator.State
	glfwWindow  *glfw.Window
	instance    *wgpu.Instance
	surface     *wgpu.Surface
	adapter     *wgpu.Adapter
	device      *wgpu.Device
	queue       *wgpu.Queue
	swapChain   *wgpu.SwapChain
	renderer    *renderer.Renderer
	fontMgr     *font.FontManager
	inputMapper *input.InputMapper
}

func NewWindow(cfg *config.Config, state *emulator.State) (*Window, error) {
	if err := glfw.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize GLFW: %w", err)
	}

	// Tell GLFW NOT to create an OpenGL context. We are using WebGPU!
	glfw.WindowHint(glfw.ClientAPI, glfw.NoAPI)
	glfwWin, err := glfw.CreateWindow(cfg.InitialWidth, cfg.InitialHeight, "GVTE Terminal", nil, nil)

	if err != nil {
		glfw.Terminate()
		return nil, fmt.Errorf("failed to create GLFW window: %w", err)
	}
	fmt.Println("glfw Cretead window")

	// Backends must be set explicitly: a zero InstanceDescriptor means
	// InstanceBackend_None, which leaves wgpu with no backend to enumerate.
	fmt.Println("Build instance")
	instance := wgpu.CreateInstance(&wgpu.InstanceDescriptor{
		Backends: wgpu.InstanceBackend_Primary | wgpu.InstanceBackend_GL,
	})

	fmt.Println("Get surface descriptor")
	desc := getSurfaceDescriptor(glfwWin)
	if desc == nil {
		return nil, fmt.Errorf("failed to get surface descriptor from GLFW window")
	}
	surface := instance.CreateSurface(desc)

	adapter, err := instance.RequestAdapter(&wgpu.RequestAdapterOptions{
		CompatibleSurface: surface,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to request WebGPU adapter: %w", err)
	}

	device, err := adapter.RequestDevice(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to request WebGPU device: %w", err)
	}

	queue := device.GetQueue()

	width, height := glfwWin.GetFramebufferSize()
	prefFormat := surface.GetPreferredFormat(adapter)

	swapChain, err := device.CreateSwapChain(surface, &wgpu.SwapChainDescriptor{
		Usage:       wgpu.TextureUsage_RenderAttachment,
		Format:      prefFormat,
		Width:       uint32(width),
		Height:      uint32(height),
		PresentMode: wgpu.PresentMode_Fifo, // Fifo включает VSync
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create swap chain: %w", err)
	}

	fontMgr, err := font.NewManager(cfg.Font.Family, cfg.Font.Size)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize font manager: %w", err)
	}

	rnd := renderer.New(fontMgr, cfg, device, queue, swapChain)
	inputMapper := input.NewMapper()

	glfwWin.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		inputMapper.HandleKey(key, action, mods)
	})

	glfwWin.Show()
	glfwWin.Focus()

	return &Window{
		cfg:         cfg,
		state:       state,
		glfwWindow:  glfwWin,
		instance:    instance,
		surface:     surface,
		adapter:     adapter,
		device:      device,
		queue:       queue,
		swapChain:   swapChain,
		renderer:    rnd,
		fontMgr:     fontMgr,
		inputMapper: inputMapper,
	}, nil
}

func (w *Window) Run() error {
	defer func() {
		if w.device != nil {
			w.device.Release()
		}
		if w.adapter != nil {
			w.adapter.Release()
		}
		if w.surface != nil {
			w.surface.Release()
		}
		if w.instance != nil {
			w.instance.Release()
		}
		if w.glfwWindow != nil {
			if w.swapChain != nil {
				w.swapChain.Release() // <--- ДОБАВИТЬ ЭТО
			}
			w.glfwWindow.Destroy()
		}
		glfw.Terminate()
	}()

	for !w.glfwWindow.ShouldClose() {
		glfw.PollEvents()
		w.renderer.Render(w.state, w.fontMgr, w.cfg)
	}

	return nil
}
