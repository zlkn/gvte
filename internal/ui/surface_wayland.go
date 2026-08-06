//go:build (linux && wayland) || (freebsd && wayland) || (netbsd && wayland) || (openbsd && wayland)

package ui

import (
	"unsafe"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

func getSurfaceDescriptor(glfwWin *glfw.Window) *wgpu.SurfaceDescriptor {
	display := glfw.GetWaylandDisplay()
	surface := glfwWin.GetWaylandWindow()
	if display == nil || surface == nil {
		return nil
	}
	return &wgpu.SurfaceDescriptor{
		WaylandSurface: &wgpu.SurfaceDescriptorFromWaylandSurface{
			Display: unsafe.Pointer(display),
			Surface: unsafe.Pointer(surface),
		},
	}
}
