package renderer

import (
	"gvte/internal/config"
	"gvte/internal/emulator"
	"gvte/internal/ui/font"

	"github.com/rajveermalviya/go-webgpu/wgpu"
)

type Rect struct{ X, Y, W, H float32 }

type Renderer struct {
	FontMgr *font.FontManager
	Config  *config.Config
	Device  *wgpu.Device
	Queue   *wgpu.Queue
}

type Frame struct {
	r       *Renderer
	encoder *wgpu.CommandEncoder
	pass    *wgpu.RenderPassEncoder
}

func New(fm *font.FontManager, cfg *config.Config, device *wgpu.Device, queue *wgpu.Queue, swapChain *wgpu.SwapChain) *Renderer {
	return &Renderer{
		FontMgr: fm,
		Config:  cfg,
		Device:  device,
		Queue:   queue,
	}
}

func (r *Renderer) BeginFrame(view *wgpu.TextureView) (*Frame, error) {
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
			ClearValue: config.ColorToWgpu(r.Config.Colors.Background),
		}},
	})

	return &Frame{r: r, encoder: encoder, pass: pass}, nil
}

func (f *Frame) DrawPane(st *emulator.State, rect Rect) {
	f.pass.SetViewport(rect.X, rect.Y, rect.W, rect.H, 0, 1)
	f.pass.SetScissorRect(uint32(rect.X), uint32(rect.Y), uint32(rect.W), uint32(rect.H))

	// f.pass.SetPipeline(f.r.pipeline)
	// f.r.queue.WriteBuffer(...)  // вершины для st.Grid
	// f.pass.Draw(...)
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
