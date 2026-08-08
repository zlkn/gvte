package renderer

import (
	"gvte/internal/config"
	"gvte/internal/emulator"
	"gvte/internal/ui/font"

	"github.com/rajveermalviya/go-webgpu/wgpu"
)

type Renderer struct {
	FontMgr   *font.FontManager
	Config    *config.Config
	Device    *wgpu.Device
	Queue     *wgpu.Queue
	SwapChain *wgpu.SwapChain
}

func New(fm *font.FontManager, cfg *config.Config, device *wgpu.Device, queue *wgpu.Queue, swapChain *wgpu.SwapChain) *Renderer {
	return &Renderer{
		FontMgr:   fm,
		Config:    cfg,
		Device:    device,
		Queue:     queue,
		SwapChain: swapChain,
	}
}

// Render draws the current terminal state.
func (r *Renderer) Render(st *emulator.State, fm *font.FontManager, cfg *config.Config) {
	// 1. Получаем текстуру для текущего кадра из SwapChain
	nextTexture, err := r.SwapChain.GetCurrentTextureView()
	if err != nil {
		// Если окно свернуто или изменило размер, может быть ошибка.
		// Пока просто игнорируем кадр.
		return
	}
	// Обязательно освобождаем текстуру после использования
	defer nextTexture.Release()

	// 2. Создаем CommandEncoder для записи команд GPU
	encoder, err := r.Device.CreateCommandEncoder(&wgpu.CommandEncoderDescriptor{
		Label: "Main Render Encoder",
	})
	defer encoder.Release()
	bgColor := config.ColorToWgpu(cfg.Colors.Background)
	// 3. Начинаем RenderPass. Здесь мы указываем, что хотим очистить экран.
	renderPass := encoder.BeginRenderPass(&wgpu.RenderPassDescriptor{
		ColorAttachments: []wgpu.RenderPassColorAttachment{
			{
				View:    nextTexture,
				LoadOp:  wgpu.LoadOp_Clear,
				StoreOp: wgpu.StoreOp_Store,
				// Цвет заливки (Темно-серый фон терминала)
				// ClearValue: wgpu.Color{R: 0.1, G: 0.1, B: 0.1, A: 1.0},
				ClearValue: bgColor,
			},
		},
	})

	// TODO: Screen rendering stub (draw text grid cells, cursor, and selection)
	// В будущем здесь будут вызовы: renderPass.SetPipeline(...), renderPass.Draw(...)

	// Завершаем проход рендеринга
	renderPass.End()

	// 4. Завершаем запись команд и получаем CommandBuffer
	// FIXME: silently ignored error
	cmdBuffer, _ := encoder.Finish(nil)
	defer cmdBuffer.Release()

	// 5. Отправляем команды в очередь видеокарты
	r.Queue.Submit(cmdBuffer)

	// 6. Показываем готовый кадр на экране
	r.SwapChain.Present()
}
