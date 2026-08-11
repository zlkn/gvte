package ui

import (
	"fmt"
	"runtime"
	"sync/atomic"

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
	fbWidth     int
	fbHeight    int
	dirty       atomic.Bool
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

	prefFormat := surface.GetPreferredFormat(adapter)

	swapChain, err := device.CreateSwapChain(surface, &wgpu.SwapChainDescriptor{
		Usage:       wgpu.TextureUsage_RenderAttachment,
		Format:      prefFormat,
		Width:       uint32(cfg.InitialWidth),
		Height:      uint32(cfg.InitialHeight),
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

	glfwWin.Show()
	glfwWin.Focus()

	win := &Window{
		cfg:         cfg,
		state:       state,
		fbWidth:     cfg.InitialWidth,
		fbHeight:    cfg.InitialHeight,
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
	}

	win.dirty.Store(true)

	glfwWin.SetFramebufferSizeCallback(func(w *glfw.Window, width int, height int) {
		win.onResize(width, height)
	})

	glfwWin.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		inputMapper.HandleKey(key, action, mods)
	})

	return win, nil
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
		// Prevent CPU eating
		glfw.WaitEvents()
		if w.fbWidth > 0 && w.fbHeight > 0 && w.dirty.Swap(false) {
			w.drawFrame()
		}
	}

	return nil
}

func (w *Window) onResize(width, height int) {
	w.fbWidth, w.fbHeight = width, height

	// iconified — nothing valid to create
	if width == 0 || height == 0 {
		return
	}
	w.resizeSwapChain(width, height)

	cols := max(1, width/w.fontMgr.CellWidth)
	rows := max(1, height/w.fontMgr.CellHeight)
	if cols != w.state.Grid.Cols || rows != w.state.Grid.Rows {
		w.state.Grid.Resize(cols, rows) // internal/emulator/grid/grid.go:64
		// TODO: w.pty.Resize(cols, rows) once internal/pty.Start is implemented
	}

	//NOTE Instant redraw instead of damage for eliminate ghoust window effect
	// w.Damage()
	w.drawFrame()
}

func (w *Window) resizeSwapChain(width, height int) {
	if width == 0 || height == 0 {
		return
	}

	if w.swapChain != nil {
		w.swapChain.Release()
	}

	prefFormat := w.surface.GetPreferredFormat(w.adapter)
	// 2. Создаем новый с новыми размерами
	w.swapChain, _ = w.device.CreateSwapChain(w.surface, &wgpu.SwapChainDescriptor{
		Usage:       wgpu.TextureUsage_RenderAttachment,
		Format:      prefFormat,
		Width:       uint32(width),
		Height:      uint32(height),
		PresentMode: wgpu.PresentMode_Fifo, // Или PresentMode_Mailbox (vsync off)
	})
}

func (w *Window) drawFrame() {
	view, err := w.swapChain.GetCurrentTextureView()
	if err != nil {
		width, height := w.glfwWindow.GetFramebufferSize()
		w.resizeSwapChain(width, height)
		return
	}
	defer view.Release()

	frame, err := w.renderer.BeginFrame(view)
	if err != nil {
		return
	}

	width, height := w.glfwWindow.GetFramebufferSize()
	frame.DrawPane(w.state, renderer.Rect{W: float32(width), H: float32(height)})

	if err := frame.End(); err != nil {
		return
	}
	w.swapChain.Present()
}

func (w *Window) Damage() {
	w.dirty.Store(true)
	glfw.PostEmptyEvent()
}
