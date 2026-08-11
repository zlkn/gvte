package renderer

import (
	_ "embed"
	"encoding/binary"
	"fmt"
	"image/color"
	"math"

	"gvte/internal/config"
	"gvte/internal/emulator"
	"gvte/internal/ui/font"

	"github.com/rajveermalviya/go-webgpu/wgpu"
)

//go:embed shader.wgsl
var bgShaderWGSLCode string

// bgUniformSize is the size of the Uniforms struct in shader.wgsl: one vec4<f32>.
const bgUniformSize = 4 * 4

type Rect struct{ X, Y, W, H float32 }

type Renderer struct {
	FontMgr *font.FontManager
	Config  *config.Config
	Device  *wgpu.Device
	Queue   *wgpu.Queue

	surfaceFormat wgpu.TextureFormat
	srgbTarget    bool

	bgPipeline  *wgpu.RenderPipeline
	bgBindGroup *wgpu.BindGroup
	bgColorBuf  *wgpu.Buffer
	bgUploaded  bool
}

type Frame struct {
	r       *Renderer
	encoder *wgpu.CommandEncoder
	pass    *wgpu.RenderPassEncoder
}

func New(fm *font.FontManager, cfg *config.Config, device *wgpu.Device, queue *wgpu.Queue, surfaceFormat wgpu.TextureFormat) (*Renderer, error) {
	r := &Renderer{
		FontMgr:       fm,
		Config:        cfg,
		Device:        device,
		Queue:         queue,
		surfaceFormat: surfaceFormat,
		srgbTarget:    config.IsSrgbFormat(surfaceFormat),
	}

	if err := r.initBackgroundPipeline(); err != nil {
		r.Release()
		return nil, err
	}
	return r, nil
}

// Release frees the GPU resources owned by the renderer.
func (r *Renderer) Release() {
	if r.bgBindGroup != nil {
		r.bgBindGroup.Release()
		r.bgBindGroup = nil
	}
	if r.bgColorBuf != nil {
		r.bgColorBuf.Release()
		r.bgColorBuf = nil
	}
	if r.bgPipeline != nil {
		r.bgPipeline.Release()
		r.bgPipeline = nil
	}
}

// initBackgroundPipeline builds the pipeline that fills the pane with a solid
// color. Geometry is generated in the vertex shader, so there is no vertex
// buffer; the color arrives through a uniform buffer at group 0, binding 0.
//
// The shader module and both layouts are only needed to construct the pipeline
// and bind group, so they are released before returning.
func (r *Renderer) initBackgroundPipeline() error {
	shader, err := r.Device.CreateShaderModule(&wgpu.ShaderModuleDescriptor{
		Label:          "Background Shader",
		WGSLDescriptor: &wgpu.ShaderModuleWGSLDescriptor{Code: bgShaderWGSLCode},
	})
	if err != nil {
		return fmt.Errorf("failed to compile background shader: %w", err)
	}
	defer shader.Release()

	r.bgColorBuf, err = r.Device.CreateBuffer(&wgpu.BufferDescriptor{
		Label: "Background Color Buffer",
		Size:  bgUniformSize,
		Usage: wgpu.BufferUsage_Uniform | wgpu.BufferUsage_CopyDst,
	})
	if err != nil {
		return fmt.Errorf("failed to create background color buffer: %w", err)
	}

	bindGroupLayout, err := r.Device.CreateBindGroupLayout(&wgpu.BindGroupLayoutDescriptor{
		Label: "Background Bind Group Layout",
		Entries: []wgpu.BindGroupLayoutEntry{{
			Binding:    0,
			Visibility: wgpu.ShaderStage_Fragment,
			Buffer: wgpu.BufferBindingLayout{
				Type:           wgpu.BufferBindingType_Uniform,
				MinBindingSize: bgUniformSize,
			},
		}},
	})
	if err != nil {
		return fmt.Errorf("failed to create background bind group layout: %w", err)
	}
	defer bindGroupLayout.Release()

	r.bgBindGroup, err = r.Device.CreateBindGroup(&wgpu.BindGroupDescriptor{
		Label:  "Background Bind Group",
		Layout: bindGroupLayout,
		Entries: []wgpu.BindGroupEntry{{
			Binding: 0,
			Buffer:  r.bgColorBuf,
			Offset:  0,
			Size:    bgUniformSize,
		}},
	})
	if err != nil {
		return fmt.Errorf("failed to create background bind group: %w", err)
	}

	pipelineLayout, err := r.Device.CreatePipelineLayout(&wgpu.PipelineLayoutDescriptor{
		Label:            "Background Pipeline Layout",
		BindGroupLayouts: []*wgpu.BindGroupLayout{bindGroupLayout},
	})
	if err != nil {
		return fmt.Errorf("failed to create background pipeline layout: %w", err)
	}
	defer pipelineLayout.Release()

	r.bgPipeline, err = r.Device.CreateRenderPipeline(&wgpu.RenderPipelineDescriptor{
		Label:  "Background Pipeline",
		Layout: pipelineLayout,
		Vertex: wgpu.VertexState{
			Module:     shader,
			EntryPoint: "vs_main",
		},
		Fragment: &wgpu.FragmentState{
			Module:     shader,
			EntryPoint: "fs_main",
			Targets: []wgpu.ColorTargetState{{
				Format:    r.surfaceFormat,
				Blend:     &wgpu.BlendState_Replace,
				WriteMask: wgpu.ColorWriteMask_All,
			}},
		},
		Primitive: wgpu.PrimitiveState{
			Topology: wgpu.PrimitiveTopology_TriangleList,
			CullMode: wgpu.CullMode_None,
		},
		// Count must be at least 1; a zero MultisampleState is invalid.
		Multisample: wgpu.MultisampleState{
			Count: 1,
			Mask:  math.MaxUint32,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create background pipeline: %w", err)
	}
	return nil
}

// setBackgroundColor uploads c into the uniform buffer read by fs_main.
func (r *Renderer) setBackgroundColor(c color.Color) error {
	v := config.ColorToRGBA(c, r.srgbTarget)

	// WGSL vec4<f32> is four tightly packed little-endian floats.
	var buf [bgUniformSize]byte
	for i, f := range v {
		binary.LittleEndian.PutUint32(buf[i*4:], math.Float32bits(float32(f)))
	}
	return r.Queue.WriteBuffer(r.bgColorBuf, 0, buf[:])
}

func (r *Renderer) BeginFrame(view *wgpu.TextureView) (*Frame, error) {
	// The uniform is constant until the color scheme changes, so upload it once
	// instead of every frame.
	if !r.bgUploaded {
		if err := r.setBackgroundColor(r.Config.Colors.Background); err != nil {
			return nil, fmt.Errorf("failed to upload background color: %w", err)
		}
		r.bgUploaded = true
	}

	encoder, err := r.Device.CreateCommandEncoder(&wgpu.CommandEncoderDescriptor{
		Label: "Main Render Encoder",
	})
	if err != nil {
		return nil, err
	}

	pass := encoder.BeginRenderPass(&wgpu.RenderPassDescriptor{
		Label: "Main Render Pass",
		ColorAttachments: []wgpu.RenderPassColorAttachment{{
			View:       view,
			LoadOp:     wgpu.LoadOp_Clear,
			StoreOp:    wgpu.StoreOp_Store,
			ClearValue: config.ColorToWgpu(r.Config.Colors.Background, r.srgbTarget),
		}},
	})
	if pass == nil {
		encoder.Release()
		return nil, fmt.Errorf("failed to begin render pass")
	}

	return &Frame{r: r, encoder: encoder, pass: pass}, nil
}

func (f *Frame) DrawPane(st *emulator.State, rect Rect) {
	f.pass.SetViewport(rect.X, rect.Y, rect.W, rect.H, 0, 1)
	f.pass.SetScissorRect(uint32(rect.X), uint32(rect.Y), uint32(rect.W), uint32(rect.H))

	f.drawBackground()
}

// drawBackground fills the current viewport with the two triangles generated by
// vs_main, shaded with the uniform color.
func (f *Frame) drawBackground() {
	f.pass.SetPipeline(f.r.bgPipeline)
	f.pass.SetBindGroup(0, f.r.bgBindGroup, nil)
	f.pass.Draw(6, 1, 0, 0)
}

func (f *Frame) End() error {
	defer f.pass.Release()
	defer f.encoder.Release()

	if err := f.pass.End(); err != nil {
		return err
	}
	cmd, err := f.encoder.Finish(nil)
	if err != nil {
		return err
	}
	defer cmd.Release()

	f.r.Queue.Submit(cmd)
	return nil
}
