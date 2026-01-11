package main

import "fmt"

// Device - интерфейс устройства (Implementor)
type Device interface {
	Print(data string)
}

// Monitor - конкретное устройство
type Monitor struct{}

func (m *Monitor) Print(data string) {
	fmt.Printf("Displaying on monitor: %s\n", data)
}

// Printer - конкретное устройство
type Printer struct{}

func (p *Printer) Print(data string) {
	fmt.Printf("Printing to paper: %s\n", data)
}

// Output - абстракция вывода
type Output interface {
	Render(data string)
}

// BaseOutput - базовая абстракция
type BaseOutput struct {
	device Device
}

// TextOutput - расширенная абстракция для текста
type TextOutput struct {
	BaseOutput
}

func NewTextOutput(device Device) *TextOutput {
	return &TextOutput{
		BaseOutput: BaseOutput{device: device},
	}
}

func (t *TextOutput) Render(data string) {
	t.device.Print("Text: " + data)
}

// ImageOutput - расширенная абстракция для изображений
type ImageOutput struct {
	BaseOutput
}

func NewImageOutput(device Device) *ImageOutput {
	return &ImageOutput{
		BaseOutput: BaseOutput{device: device},
	}
}

func (i *ImageOutput) Render(data string) {
	i.device.Print(fmt.Sprintf("Image: [Binary data: %s]", data))
}

func main() {
	fmt.Println("=== Bridge Pattern ===")

	monitor := &Monitor{}
	printer := &Printer{}

	textOnMonitor := NewTextOutput(monitor)
	textOnPrinter := NewTextOutput(printer)

	textOnMonitor.Render("Hello, world!")
	textOnPrinter.Render("Hello, world!")

	imageOnMonitor := NewImageOutput(monitor)
	imageOnMonitor.Render("101010101")
	fmt.Println()
}
