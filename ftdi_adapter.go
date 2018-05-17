package eltee

/**
  Mac OS wants the D2XX Helper thing installed
  http://www.ftdichip.com/Drivers/D2XX.htm
*/

/*
#cgo pkg-config: libftdi1
#include <ftdi.h>
#include <stdio.h>

char fake_buff[512];

const char* fake_write(struct ftdi_context *ftdi, const unsigned char *buf, int size)
{
    sprintf(fake_buff, "%d %d %d %d %d %d", buf[0], buf[1], buf[2], buf[3], buf[4], buf[5]);

    return fake_buff;
}
*/
import "C"
import (
	"fmt"
	"time"
)

type FtdiVersion struct {
	Major int
	Minor int
	Micro int

	Version  string
	Snapshot string
}

// Provides the libftdi version information in a go struct. This is a super simple pass through
// that is at least useful to test whether the library is available and working at all.
func FtdiLibVersion() FtdiVersion {

	out := FtdiVersion{}

	info := C.ftdi_get_library_version()

	out.Major = int(info.major)
	out.Minor = int(info.minor)
	out.Micro = int(info.micro)

	out.Version = C.GoString(info.version_str)
	out.Snapshot = C.GoString(info.snapshot_str)

	return out
}

/////////////////////

type FtdiStatus int

const (
	Stopped FtdiStatus = iota
	Closed
	Open
)

type FtdiContext struct {
	ctx    *C.struct_ftdi_context
	status FtdiStatus

	// The currentDMX is 513 bytes long because it is exactly what we will output
	// all the time. There is a leading 0 byte and then a full DMX frame of data
	currentDMX []byte

	// This is a list of potentially partial buffers which will be used to update
	// the output right before it is written on it's own clock cycle.
	pendingDMX chan []byte

	lastOpenAttempt time.Time
	lastWriteAt     time.Time
}

func NewFtdiContext() *FtdiContext {
	fc := &FtdiContext{
		status: Stopped,

		currentDMX: make([]byte, 513),
		pendingDMX: make(chan []byte, 80), // 2 seconds or so worth...
	}

	fc.ctx = C.ftdi_new()
	// Check for null???

	// // Some initial data
	// fc.currentDMX[1] = 1
	// fc.currentDMX[2] = 2
	// fc.currentDMX[3] = 3
	// fc.currentDMX[4] = 4

	return fc
}

func (fc *FtdiContext) Start() error {
	if fc == nil {
		return fmt.Errorf("Can not start nil FtdiContext")
	}

	if fc.status != Stopped {
		return fmt.Errorf("Context is already started")
	}

	fc.status = Closed
	go fc.goOpener()

	return nil
}

func (fc *FtdiContext) WriteDmx(dmx []byte) {
	if fc == nil {
		return
	}

	if fc.status != Open {
		return
	}

	fc.pendingDMX <- dmx
}

////////////////////////////////////////////////

func (fc *FtdiContext) errorString() string {
	return C.GoString(C.ftdi_get_error_string(fc.ctx))
}

const REOPEN_TIME time.Duration = time.Second * 5

// Go routine which attempts to open the interface
func (fc *FtdiContext) goOpener() {

	for {
		now := time.Now()
		reopenAt := fc.lastOpenAttempt.Add(REOPEN_TIME)

		if reopenAt.Before(now) {
			fc.lastOpenAttempt = now
			// Attempt to reopen
			if fc.attemptOpen() {
				// Cool! Start the pump routine
				go fc.goFramePump()

				// Done with this go routine
				break
			}
		} else {
			// Sleep until it is time
			toSleep := reopenAt.Sub(now)
			time.Sleep(toSleep)
		}
	}
}

// Make an attempt at opening the first device. Log the error and just
// return true or false
func (fc *FtdiContext) attemptOpen() bool {

	log.Info("FTDI: Open: Attempting to open first ftdi device")

	// Start by getting a list of devices so that we can use the first device
	var list *C.struct_ftdi_device_list
	count := C.ftdi_usb_find_all(fc.ctx, &list, 0, 0)
	log.Infof("Found %v devices", count)

	switch {
	case count == 0:
		log.Warning("FTDI: Open: No devices were found")
		return false

	case count < 0:
		log.Warningf("FTDI: Open: ftdi_usb_find_all() %v", fc.errorString())
		return false
	}

	defer C.ftdi_list_free(&list)

	result := C.ftdi_set_interface(fc.ctx, C.INTERFACE_A)
	if result != 0 {
		log.Warningf("FTDI: Open: ftdi_set_interface() %v", fc.errorString())
		return false
	}

	log.Debug("FTDI: Open: Attempting to open the first device")

	result = C.ftdi_usb_open_dev(fc.ctx, list.dev)

	if result != 0 {
		log.Warningf("FTDI: Open: ftdi_usb_open_dev() %v", fc.errorString())
		return false
	}

	// Reset it
	result = C.ftdi_usb_reset(fc.ctx)
	if result != 0 {
		log.Warningf("FTDI: Open: ftdi_usb_reset() %v", fc.errorString())
		return false
	}

	// Baud Rate
	result = C.ftdi_set_baudrate(fc.ctx, 250000)
	if result != 0 {
		log.Warningf("FTDI: Open: ftdi_set_baudrate() %v", fc.errorString())
		return false
	}

	// Line Properties
	// 8 bits, 2 stop bits, no parity
	result = C.ftdi_set_line_property(fc.ctx, C.BITS_8, C.STOP_BIT_2, C.NONE)
	if result != 0 {
		log.Warningf("FTDI: Open: ftdi_set_line_property() %v", fc.errorString())
		return false
	}

	// Flow Control
	result = C.ftdi_setflowctrl(fc.ctx, C.SIO_DISABLE_FLOW_CTRL)
	if result != 0 {
		log.Warningf("FTDI: Open: ftdi_setflowctrl() %v", fc.errorString())
		return false
	}

	// Purge Buffers
	result = C.ftdi_usb_purge_buffers(fc.ctx)
	if result != 0 {
		log.Warningf("FTDI: Open: ftdi_usb_purge_buffers() %v", fc.errorString())
		return false
	}

	// Clear RTS
	result = C.ftdi_setrts(fc.ctx, 0)
	if result != 0 {
		log.Warningf("FTDI: Open: ftdi_setrts() %v", fc.errorString())
		return false
	}

	fc.status = Open
	log.Info("FTDI: Open: Opened the device ok!")
	return true
}

const DURATION_PER_FRAME = time.Second / 30.0

const DMX_BREAK = time.Microsecond * 110
const DMX_MAB = time.Microsecond * 16

// const DURATION_PER_FRAME = time.Second

// Once the device is opened this go routine is started to output frames at
// a hopefully consistent rate
func (fc *FtdiContext) goFramePump() {
	for {
		// log.Debug("goFramePump loop")

		// While there are any frames that should update our current reality do that
		// There is a possibility here that a caller could overwhelm us with frames to
		// update at such a fast rate that we never get around to writing things out.
		// That seems unlikely enough that we aren't going to bother to guard against it.
		// This whole buffering scheme is really a minimal decoupling thing and we
		// don't _expect_ there to be massive overruns. This is just a note for the
		// crazy super performant future though, just in case ;)

		if fc.status != Open {
			log.Warningf("FTDI: Frame pump: Exiting because status is not open")
			return
		}

		// Need a slice that preserves the leading 0
		// log.Debug("Drain any pending frames")
		currentDest := fc.currentDMX[1:]
	Drain:
		for {
			select {
			case frame := <-fc.pendingDMX:
				// Copy frames into that slice directly
				copy(currentDest, frame)
			default:
				break Drain
			}
		}

		// Okay the fc.currentDMX is a valid buffer that we want sent out on the wire,
		// so we do that

		fc.lastWriteAt = time.Now()
		// log.Debug("Writing a frame at %v", fc.lastWriteAt)

		fc.setBreak(true)
		time.Sleep(DMX_BREAK)
		fc.setBreak(false)
		time.Sleep(DMX_MAB)

		b := &fc.currentDMX[0]

		// str := C.GoString(C.fake_write(fc.ctx, (*_Ctype_uchar)(b), 513))
		// log.Debugf("Fake Write> %v", str)

		result := C.ftdi_write_data(fc.ctx, (*_Ctype_uchar)(b), 513)
		if result < 0 {
			// This is bad - throw us into error state
			log.Warningf("FTDI: Frame pump: %v", fc.errorString())
			fc.enterErrorState()
			return
		}

		// Life is grand, keep trucking
		nextFrameAt := fc.lastWriteAt.Add(DURATION_PER_FRAME)
		sleepTime := nextFrameAt.Sub(time.Now())

		time.Sleep(sleepTime)
	}
}

func (fc *FtdiContext) enterErrorState() {
	log.Warningf("FTDI: Entering closed state because of error")
	fc.status = Stopped

	// Set this to now so a regular interval will pass before we try to re-open the device
	// and get everything going again
	fc.lastOpenAttempt = time.Now()

	// Restart the attempt to open goroutine
	fc.Start()
}

// This important little guy is needed to implement the proper serial signaling
// so that the fixtures actually find the frames of data!
func (fc *FtdiContext) setBreak(on bool) bool {
	var result C.int
	if on {
		result = C.ftdi_set_line_property2(fc.ctx, C.BITS_8, C.STOP_BIT_2, C.NONE, C.BREAK_ON)
	} else {
		result = C.ftdi_set_line_property2(fc.ctx, C.BITS_8, C.STOP_BIT_2, C.NONE, C.BREAK_OFF)
	}

	if result < 0 {
		return false
	}

	return true
}
