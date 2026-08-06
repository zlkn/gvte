//go:build (linux && !wayland) || (freebsd && !wayland) || (netbsd && !wayland) || (openbsd && !wayland)

package ui

import (
	"unsafe"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

func getSurfaceDescriptor(glfwWin *glfw.Window) *wgpu.SurfaceDescriptor {
	display := glfw.GetX11Display()
	x11Win := glfwWin.GetX11Window()
	if display == nil || x11Win == 0 {
		return nil
	}
	return &wgpu.SurfaceDescriptor{
		XlibWindow: &wgpu.SurfaceDescriptorFromXlibWindow{
			Display: unsafe.Pointer(display),
			Window:  uint32(x11Win),
		},
	}
}
