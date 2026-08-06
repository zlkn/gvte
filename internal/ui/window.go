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

	// Backends must be set explicitly: a zero InstanceDescriptor means
	// InstanceBackend_None, which leaves wgpu with no backend to enumerate.
	instance := wgpu.CreateInstance(&wgpu.InstanceDescriptor{
		Backends: wgpu.InstanceBackend_Primary | wgpu.InstanceBackend_GL,
	})

	var surface *wgpu.Surface
	if desc := getSurfaceDescriptor(glfwWin); desc != nil {
		surface = instance.CreateSurface(desc)
	}

	var adapter *wgpu.Adapter
	var adapterErr error

	// Strategy 1: Request adapter with surface compatibility
	if surface != nil {
		adapter, adapterErr = instance.RequestAdapter(&wgpu.RequestAdapterOptions{
			CompatibleSurface: surface,
		})
		if adapterErr != nil {
			// Strategy 2: Request adapter with fallback adapter enabled
			adapter, adapterErr = instance.RequestAdapter(&wgpu.RequestAdapterOptions{
				CompatibleSurface:    surface,
				ForceFallbackAdapter: true,
			})
		}
	}

	// Strategy 3: Request any available adapter (high performance / low power fallback)
	if adapter == nil {
		adapter, adapterErr = instance.RequestAdapter(&wgpu.RequestAdapterOptions{
			PowerPreference: wgpu.PowerPreference_LowPower,
		})
		if adapterErr != nil {
			adapter, adapterErr = instance.RequestAdapter(&wgpu.RequestAdapterOptions{
				PowerPreference: wgpu.PowerPreference_HighPerformance,
			})
		}
	}

	if adapter == nil || adapterErr != nil {
		if surface != nil {
			surface.Release()
		}
		instance.Release()
		glfwWin.Destroy()
		glfw.Terminate()
		return nil, fmt.Errorf("failed to request WebGPU adapter: %w", adapterErr)
	}

	device, err := adapter.RequestDevice(nil)
	if err != nil {
		adapter.Release()
		surface.Release()
		instance.Release()
		glfwWin.Destroy()
		glfw.Terminate()
		return nil, fmt.Errorf("failed to request WebGPU device: %w", err)
	}

	queue := device.GetQueue()

	fontMgr, err := font.NewManager(cfg.Font.Family, cfg.Font.Size)
	if err != nil {
		device.Release()
		adapter.Release()
		surface.Release()
		instance.Release()
		glfwWin.Destroy()
		glfw.Terminate()
		return nil, fmt.Errorf("failed to initialize font manager: %w", err)
	}

	rnd := renderer.New(fontMgr, cfg)
	inputMapper := input.NewMapper()

	return &Window{
		cfg:         cfg,
		state:       state,
		glfwWindow:  glfwWin,
		instance:    instance,
		surface:     surface,
		adapter:     adapter,
		device:      device,
		queue:       queue,
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
			w.glfwWindow.Destroy()
		}
		glfw.Terminate()
	}()

	for !w.glfwWindow.ShouldClose() {
		glfw.PollEvents()
		w.renderer.Render(w.state)
	}

	return nil
}
