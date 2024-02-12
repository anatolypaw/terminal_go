package framebuffer

import (
	"errors"
	"image"
	"image/draw"
	"os"
	"os/signal"
	"syscall"
	"unsafe"
)

type Framebuffer struct {
	Fi fbFixScreenInfo // Fixed buffer settings.
	Vi fbVarScreenInfo // Variable buffer settings.

	// Framebuffer state and access bits.
	file        *os.File // Framebuffer file descriptor.
	mem         []byte   // mmap'd memory.
	switchState int      // Current switch state.

	// pre-allocated clean screen
	zero []byte
}

// Open the framebuffer
func Open(path string) (fb *Framebuffer, err error) {
	fb = new(Framebuffer)
	fb.switchState = _FB_ACTIVE

	defer func() {
		// Ensure resources are properly cleaned up when things go booboo.
		if err != nil {
			fb.Close()
		}
	}()

	fb.file, err = os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return
	}

	// Fetch fixed buffer information.
	err = ioctl(fb.file.Fd(), _IOGET_FSCREENINFO, unsafe.Pointer(&fb.Fi))
	if err != nil {
		return
	}

	// Fetch variable information.
	err = ioctl(fb.file.Fd(), _IOGET_VSCREENINFO, unsafe.Pointer(&fb.Vi))
	if err != nil {
		return
	}

	// Ensure we are in PACKED_PIXELS mode. Others are useless to us.
	if fb.Fi.typ != _TYPE_PACKED_PIXELS {
		err = errors.New("framebuffer is not in PACKED PIXELS mode")
		return
	}

	if fb.Fi.smemlen == 0 {
		if fb.Fi.ywrapstep == 0 {
			fb.Fi.smemlen = uint32(fb.Vi.xres * fb.Vi.yres * fb.Vi.bitsPerPixel / 8)
		} else {
			fb.Fi.smemlen = uint32(uint32(fb.Fi.ywrapstep) * fb.Vi.yres)
		}
	}

	// mmap the buffer's memory.
	fb.mem, err = syscall.Mmap(int(fb.file.Fd()), 0, int(fb.Fi.smemlen),
		syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		err = errors.New("Framebuffer: Mmap failed: " + err.Error())
		return
	}

	// Create pre-allocated zero-memory.
	// This is used to do fast screen clears.
	fb.zero = make([]byte, len(fb.mem))

	// Move viewport to top-left corner.
	if fb.Vi.xoffset != 0 || fb.Vi.yoffset != 0 {
		vi := fb.Vi.Copy()
		vi.xoffset = 0
		vi.yoffset = 0

		err = ioctl(fb.file.Fd(), _IOPAN_DISPLAY, unsafe.Pointer(vi))
		if err != nil {
			return
		}
	}

	// Clear screen
	fb.Clear()
	go fb.pollSignals()
	return
}

// Close closes the framebuffer and cleans up its resources.
func (c *Framebuffer) Close() (err error) {
	if c.mem != nil {
		syscall.Munmap(c.mem)
		c.mem = nil
	}

	if c.file != nil {
		c.file.Close()
		c.file = nil
	}

	return
}

// File returns the underlying framebuffer file descriptor.
// This can be used in custom IOCTL calls.
//
// Use with caution and do not close it manually.
func (fb *Framebuffer) File() *os.File {
	return fb.file
}

// Image returns the pixel buffer as a image.Image instance.
// Returns nil if something went wrong.
func (fb *Framebuffer) Image() (draw.Image, error) {
	p := fb.mem
	s := int(fb.Fi.ywrapstep)
	if s == 0 {
		panic("No ywrapstep")
	}
	r := image.Rect(0, 0, int(fb.Vi.xres), int(fb.Vi.yres))

	//return &BGRA{Pix: p, Stride: s, Rect: r}, nil
	return &image.RGBA{Pix: p, Stride: s, Rect: r}, nil

}

// Очищаем фреймбуффер копированием нулей
func (c *Framebuffer) Clear() {
	copy(c.mem, c.zero)

}

// pollSignals polls for user signals.
func (c *Framebuffer) pollSignals() {
	signals := make(chan os.Signal, 2)
	signal.Notify(signals, syscall.SIGUSR1, syscall.SIGUSR2)

	for sig := range signals {
		switch sig {
		case syscall.SIGUSR1: // Release
			c.switchState = _FB_REL_REQ

		case syscall.SIGUSR2: // Acquire
			c.switchState = _FB_ACQ_REQ
		}
	}
}
