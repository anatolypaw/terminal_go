// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package framebuffer

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"unsafe"
)

type Framebuffer struct {
	// Backup storage.
	// These hold the initial system state, which will be restored once we shut down.
	origFi   fbFixScreenInfo // Fixed buffer settings.
	origVi   fbVarScreenInfo // Variable buffer settings.
	origR    [256]uint16     // Palette red channel.
	origG    [256]uint16     // Palette green channel.
	origB    [256]uint16     // Palette blue channel.
	origA    [256]uint16     // Palette transparent channel.
	origVT   vtMode          // Virtual terminal mode.
	origVTNo int             // Virtual terminal number.
	origKd   int             // KD mode.

	// Framebuffer state and access bits.
	fb          *os.File // Framebuffer file descriptor.
	mem         []byte   // mmap'd memory.
	dev         string   // name of the device we are using.
	switchState int      // Current switch state.

	// pre-allocated scratchpad values.
	zero []byte
	tmpR [256]uint16
	tmpG [256]uint16
	tmpB [256]uint16
	tmpA [256]uint16
}

// Open opens the framebuffer with the given display mode.
//
// If mode is nil, the default framebuffer mode is used.
//
// The framebuffer is usually initialized to a specific display mode by the
// kernel itself. While this library supplies the means to alter the current
// display mode, this may not always have any effect as a driver can
// choose to ignore your requested values. Besides that, it is generally
// considered safer to use the external `fbset` command for this purpose.
//
// Video modes for the framebuffer require very precise timing values to
// be supplied along with any desired resolution. Doing this incorrectly
// can damage the display. Refer to Canvas.Modes() and Canvas.FindMode()
// for more information. Canvas.CurrentMode() can be used to see which
// mode is actually being used.
func Open(dev string) (c *Framebuffer, err error) {
	c = new(Framebuffer)
	c.dev = dev
	c.origVTNo = 0
	c.switchState = _FB_ACTIVE

	defer func() {
		// Ensure resources are properly cleaned up when things go booboo.
		if err != nil {
			c.Close()
		}
	}()

	// Open the frame buffer.
	c.fb, err = os.OpenFile(c.dev, os.O_RDWR, 0)
	if err != nil {
		return
	}

	// Fetch original fixed buffer information.
	// This will never be changed, but we need the information
	// in various places.
	err = ioctl(c.fb.Fd(), _IOGET_FSCREENINFO, unsafe.Pointer(&c.origFi))
	if err != nil {
		return
	}

	// Fetch original variable information.
	err = ioctl(c.fb.Fd(), _IOGET_VSCREENINFO, unsafe.Pointer(&c.origVi))
	if err != nil {
		return
	}

	// Fetch original color palette if applicable.
	if c.origVi.bitsPerPixel == 8 || c.origFi.visual == _VISUAL_DIRECTCOLOR {
		var cm fb_cmap
		cm.start = 0
		cm.len = 256
		cm.red = unsafe.Pointer(&c.origR[0])
		cm.green = unsafe.Pointer(&c.origG[0])
		cm.blue = unsafe.Pointer(&c.origB[0])
		cm.transp = unsafe.Pointer(&c.origA[0])

		err = ioctl(c.fb.Fd(), _IOGET_CMAP, unsafe.Pointer(&cm))
		if err != nil {
			return
		}
	}

	// Get KD mode
	err = ioctl(c.fb.Fd(), _KDGETMODE, unsafe.Pointer(&c.origKd))
	if err != nil {
		return
	}

	// Get original vt mode
	err = ioctl(c.fb.Fd(), _VT_GETMODE, unsafe.Pointer(&c.origVT))
	if err != nil {
		return

	}

	// Ensure we are in PACKED_PIXELS mode. Others are useless to us.
	if c.origFi.typ != _TYPE_PACKED_PIXELS {
		err = errors.New("Canvas.Open: Framebuffer is not in PACKED PIXELS mode. Unable to continue")
		return
	}

	// If we have a non-standard pixel format, we can't continue.
	// if c.orig_vi.nonstd != 0 {
	// 	err = errors.New("Canvas.Open: Framebuffer uses a non-standard pixel format. This is not supported.")
	// 	return
	// }

	if c.origFi.smemlen == 0 {
		if c.origFi.ywrapstep == 0 {
			c.origFi.smemlen = uint32(c.origVi.xres * c.origVi.yres * c.origVi.bitsPerPixel / 8)
		} else {
			c.origFi.smemlen = uint32(uint32(c.origFi.ywrapstep) * c.origVi.yres)
		}
	}

	// mmap the buffer's memory.
	c.mem, err = syscall.Mmap(int(c.fb.Fd()), 0, int(c.origFi.smemlen),
		syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		err = errors.New("Canvas.Open: Mmap failed: " + err.Error())
		return
	}

	// Create pre-allocated zero-memory.
	// This is used to do fast screen clears.
	c.zero = make([]byte, len(c.mem))

	// Move viewport to top-left corner.
	if c.origVi.xoffset != 0 || c.origVi.yoffset != 0 {
		vi := c.origVi.Copy()
		vi.xoffset = 0
		vi.yoffset = 0

		err = ioctl(c.fb.Fd(), _IOPAN_DISPLAY, unsafe.Pointer(vi))
		if err != nil {
			return
		}
	}

	if c.tty != nil {
		// Switch terminal to graphics mode.
		err = ioctl(c.fb.Fd(), _KDSETMODE, _KD_GRAPHICS)
		if err != nil {
			return
		}

		// Activate the given tty.
		err = c.activateCurrent(c.tty)
		if err != nil {
			return
		}
	}

	// Clear screen
	c.Clear()
	go c.pollSignals()
	return
}

// Close closes the framebuffer and cleans up its resources.
func (c *Framebuffer) Close() (err error) {
	if c.mem != nil {
		syscall.Munmap(c.mem)
		c.mem = nil
	}

	if c.fb != nil {
		c.fb.Close()
		c.fb = nil
	}

	return
}

// File returns the underlying framebuffer file descriptor.
// This can be used in custom IOCTL calls.
//
// Use with caution and do not close it manually.
func (c *Framebuffer) File() *os.File {
	return c.fb
}

// Image returns the pixel buffer as a draw.Image instance.
// Returns nil if something went wrong.
func (c *Framebuffer) Image() (draw.Image, error) {
	mode, err := c.CurrentMode()
	if err != nil {
		return nil, err
	}

	p := c.mem
	s := int(c.origFi.ywrapstep)
	if s == 0 {
		s = mode.Stride()
	}
	r := image.Rect(0, 0, mode.Geometry.XVRes, mode.Geometry.YVRes)

	// Find out which image type we should be returning.
	// This depends on the current pixel format.
	return &image.RGBA{Pix: p, Stride: s, Rect: r}, nil

}

// Clear clears (zeroes) the framebuffer memory.
func (c *Framebuffer) Clear() {
	copy(c.mem, c.zero)
}

// Accelerated returns true if the framebuffer
// currently supports hardware acceleration.
func (c *Framebuffer) Accelerated() bool {
	return c.origFi.accel != _ACCEL_NONE
}

// Buffer provides direct access to the entire memory-mapped pixel buffer.
func (c *Framebuffer) Buffer() []byte {
	return c.mem
}

// setMode sets the given display mode.
// If the mode is nil, this returns without error;
// the call is simply ignored.
func (c *Framebuffer) setMode(dm *DisplayMode) error {
	if dm == nil {
		return nil
	}

	var v fbVarScreenInfo

	err := ioctl(c.fb.Fd(), _IOGET_VSCREENINFO, unsafe.Pointer(&v))
	if err != nil {
		return err
	}

	v.xres = uint32(dm.Geometry.XRes)
	v.yres = uint32(dm.Geometry.YRes)
	v.xresVirtual = uint32(dm.Geometry.XVRes)
	v.yresVirtual = uint32(dm.Geometry.YVRes)
	v.bitsPerPixel = uint32(dm.Geometry.Depth)
	v.pixclock = uint32(dm.Timings.Pixclock)
	v.leftMargin = uint32(dm.Timings.Left)
	v.rightMargin = uint32(dm.Timings.Right)
	v.upperMargin = uint32(dm.Timings.Upper)
	v.lowerMargin = uint32(dm.Timings.Lower)
	v.hsyncLen = uint32(dm.Timings.HSLen)
	v.vsyncLen = uint32(dm.Timings.VSLen)
	v.sync = uint32(dm.Sync)
	v.vmode = uint32(dm.VMode)

	pf := dm.Format
	v.red.length = uint32(pf.RedBits)
	v.red.offset = uint32(pf.RedShift)
	v.red.msb_right = 1

	v.green.length = uint32(pf.GreenBits)
	v.green.offset = uint32(pf.GreenShift)
	v.green.msb_right = 1

	v.blue.length = uint32(pf.BlueBits)
	v.blue.offset = uint32(pf.BlueShift)
	v.blue.msb_right = 1

	v.transparent.length = uint32(pf.AlphaBits)
	v.transparent.offset = uint32(pf.AlphaShift)
	v.transparent.msb_right = 1

	v.xoffset = 0
	v.yoffset = 0

	return ioctl(c.fb.Fd(), _IOPUT_VSCREENINFO, unsafe.Pointer(&v))
}

// CurrentMode returns the current framebuffer display mode.
func (c *Framebuffer) CurrentMode() (*DisplayMode, error) {
	var v fbVarScreenInfo
	var dm DisplayMode

	if ioctl(c.fb.Fd(), _IOGET_VSCREENINFO, unsafe.Pointer(&v)) != nil {
		return nil, errors.New("Canvas.CurrentMode failed")
	}

	dm.Accelerated = c.origFi.accel != _ACCEL_NONE

	dm.Geometry.XRes = int(v.xres)
	dm.Geometry.YRes = int(v.yres)
	dm.Geometry.XVRes = int(v.xresVirtual)
	dm.Geometry.YVRes = int(v.yresVirtual)
	dm.Geometry.Depth = int(v.bitsPerPixel)
	dm.Timings.Pixclock = int(v.pixclock)
	dm.Timings.Left = int(v.leftMargin)
	dm.Timings.Right = int(v.rightMargin)
	dm.Timings.Upper = int(v.upperMargin)
	dm.Timings.Lower = int(v.lowerMargin)
	dm.Timings.HSLen = int(v.hsyncLen)
	dm.Timings.VSLen = int(v.vsyncLen)
	dm.Sync = int(v.sync)
	dm.VMode = int(v.vmode)

	var pf PixelFormat
	pf.RedBits = uint8(v.red.length)
	pf.RedShift = uint8(v.red.offset)
	pf.GreenBits = uint8(v.green.length)
	pf.GreenShift = uint8(v.green.offset)
	pf.BlueBits = uint8(v.blue.length)
	pf.BlueShift = uint8(v.blue.offset)
	pf.AlphaBits = uint8(v.transparent.length)
	pf.AlphaShift = uint8(v.transparent.offset)
	dm.Format = pf

	return &dm, nil
}

// FindMode finds the display mode with the given name.
// Returns nil if it does not exist.
//
// The external `fbset` tool comes with a set of default modes
// which are stored in the file `/etc/fb.modes`. We read this file
// and extract the set of video modes from it. These modes each have
// a name by which they can be identified. When supplying a new
// mode to this function, it should come in the form of this name.
// For example: "1600x1200-76".
//
// New video modes can be added to the `/etc/fb.modes` file.
func (c *Framebuffer) FindMode(name string) *DisplayMode {
	modes, err := c.Modes()
	if err != nil {
		return nil
	}

	for _, m := range modes {
		if strings.EqualFold(m.Name, name) {
			return m
		}
	}

	return nil
}

// Modes returns the list of supported display modes.
// These are read from `/etc/fb.modes`.
// This can be called before the framebuffer has been opened.
func (c *Framebuffer) Modes() ([]*DisplayMode, error) {
	fd, err := os.Open("/etc/fb.modes")
	if err != nil {
		return nil, err
	}

	defer fd.Close()

	return readFBModes(fd)
}

// Palette returns the current framebuffer color palette.
func (c *Framebuffer) Palette() (color.Palette, error) {
	var cm fb_cmap

	cm.start = 0
	cm.len = 256
	cm.red = unsafe.Pointer(&c.tmpR[0])
	cm.green = unsafe.Pointer(&c.tmpG[0])
	cm.blue = unsafe.Pointer(&c.tmpB[0])
	cm.transp = unsafe.Pointer(&c.tmpA[0])

	if ioctl(c.fb.Fd(), _IOGET_CMAP, unsafe.Pointer(&cm)) != nil {
		return nil, errors.New("Canvas.Palette failed")
	}

	s := int(cm.start)
	pal := make(color.Palette, cm.len)

	for i := range pal {
		pal[i] = color.NRGBA{
			uint8(c.tmpR[i+s] >> 8),
			uint8(c.tmpG[i+s] >> 8),
			uint8(c.tmpB[i+s] >> 8),
			uint8(c.tmpA[i+s] >> 8),
		}
	}

	return pal, nil
}

// SetPalette sets the current framebuffer color palette.
func (c *Framebuffer) SetPalette(pal color.Palette) error {
	if len(pal) >= 256 {
		pal = pal[:256]
	}

	for i, clr := range pal {
		r, g, b, a := clr.RGBA()
		c.tmpR[i] = uint16(r >> 16)
		c.tmpG[i] = uint16(g >> 16)
		c.tmpB[i] = uint16(b >> 16)
		c.tmpA[i] = uint16(a >> 16)
	}

	var cm fb_cmap
	cm.start = 0
	cm.len = 256
	cm.red = unsafe.Pointer(&c.tmpR[0])
	cm.green = unsafe.Pointer(&c.tmpG[0])
	cm.blue = unsafe.Pointer(&c.tmpB[0])
	cm.transp = unsafe.Pointer(&c.tmpA[0])

	if ioctl(c.fb.Fd(), _IOPUT_CMAP, unsafe.Pointer(&cm)) != nil {
		return errors.New("Canvas.SetPalette failed")
	}

	return nil
}

func (c *Framebuffer) switchAcquire() {
	if c.tty != nil {
		ioctl(c.tty.Fd(), _VT_RELDISP, _VT_ACKACQ)
	}
	c.switchState = _FB_ACTIVE
}

func (c *Framebuffer) switchRelease() {
	if c.tty != nil {
		ioctl(c.tty.Fd(), _VT_RELDISP, 1)
	}
	c.switchState = _FB_INACTIVE
}

func (c *Framebuffer) switchInit() error {
	if c.tty == nil {
		return nil
	}

	var vm vtMode

	vm.mode = _VT_PROCESS
	vm.waitv = 0
	vm.relsig = int16(syscall.SIGUSR1)
	vm.acqsig = int16(syscall.SIGUSR2)

	return ioctl(c.tty.Fd(), _VT_SETMODE, unsafe.Pointer(&vm))
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

func (c *Framebuffer) activateCurrent(tty *os.File) error {
	var vts vtStat

	err := ioctl(tty.Fd(), _VT_GETSTATE, unsafe.Pointer(&vts))
	if err != nil {
		return err
	}

	err = ioctl(tty.Fd(), _VT_ACTIVATE, int(vts.vActive))
	if err != nil {
		return err
	}

	return ioctl(tty.Fd(), _VT_WAITACTIVE, int(vts.vActive))
}
