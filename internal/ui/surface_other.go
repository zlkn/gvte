//go:build !linux && !freebsd && !netbsd && !openbsd

package ui

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

func getSurfaceDescriptor(glfwWin *glfw.Window) *wgpu.SurfaceDescriptor {
	return nil
}
