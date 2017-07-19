package json

import (
	"errors"
	"io"
)

const (
	scratchByteArrayLen = 32
)

var (
	zeroByteSlice = []byte{}[:0:0]
)

var (
	_ quickReader = &quickReaderStream{}
	_ quickReader = &quickReaderSlice{}
)

// quickReader is a hybrid of reader and buffer interfaces with methods giving
// specific attention to the performance needs found in a decoder.
// Implementations cover io.Reader as well as []byte directly.
//
// In particular, it allows:
//
//   - returning byte-slices with zero-copying (you were warned!) when possible
//   - returning byte-slices for short reads which will be reused (you were warned!)
//   - putting a 'track' point in the buffer, and later yielding all those bytes at once
//   - counting the number of bytes read (for use in parser error messages, mainly)
//
// All of these shortcuts mean correct usage is essential to avoid unexpected behaviors,
// but in return allow avoiding many, many common sources of memory allocations in a parser.
//
// Implementations panic on unexpected IO errors.  TODO change this to swallow errors until checked.
type quickReader interface {

	// Read n bytes into a byte slice which may be shared and must not be reused
	// After any additional calls to this reader.
	// readnzc will use the implementation scratch buffer if possible,
	// i.e. n < len(scratchbuf), or may return a view of the []byte being decoded from.
	// Requesting a zero length read will return `zeroByteSlice`, a len-zero cap-zero slice.
	readnzc(n int) []byte

	// Read n bytes into a new byte slice.
	// If zero-copy views into existing buffers are acceptable (e.g. you know you
	// won't later mutate, reference or expose this memory again), prefer `readnzc`.
	// If you already have an existing slice of sufficient size to reuse, prefer `readb`.
	// Requesting a zero length read will return `zeroByteSlice`, a len-zero cap-zero slice.
	readn(n int) []byte

	// Read `len(b)` bytes into the given slice, starting at its beginning,
	// overwriting all values, and disregarding any extra capacity.
	readb(b []byte)

	readn1() uint8
	readn1eof() (v uint8, eof bool)
	unreadn1()
	numread() int // number of bytes read
	track()
	stopTrack() []byte
}

// quickReaderStream is a quickReader that reads off an io.Reader.
// Initialize it by wrapping an ioDecByteScanner around your io.Reader and dumping it in.
// While this implementation does use some internal buffers, it's still advisable
// to use a buffered reader to avoid small reads for any external IO like disk or network.
type quickReaderStream struct {
	br         *readerByteScanner
	scratch    [scratchByteArrayLen]byte // temp byte array re-used internally for efficiency during read.
	n          int                       // num read
	tracking   []byte                    // tracking bytes read
	isTracking bool
}

func (z *quickReaderStream) numread() int {
	return z.n
}

func (z *quickReaderStream) readnzc(n int) (bs []byte) {
	if n == 0 {
		return zeroByteSlice
	}
	if n < len(z.scratch) {
		bs = z.scratch[:n]
	} else {
		bs = make([]byte, n)
	}
	z.readb(bs)
	return
}

func (z *quickReaderStream) readn(n int) (bs []byte) {
	if n == 0 {
		return zeroByteSlice
	}
	bs = make([]byte, n)
	z.readb(bs)
	return
}

func (z *quickReaderStream) readb(bs []byte) {
	if len(bs) == 0 {
		return
	}
	n, err := io.ReadAtLeast(z.br, bs, len(bs))
	z.n += n
	if err != nil {
		panic(err)
	}
	if z.isTracking {
		z.tracking = append(z.tracking, bs...)
	}
}

func (z *quickReaderStream) readn1() (b uint8) {
	b, err := z.br.ReadByte()
	if err != nil {
		panic(err)
	}
	z.n++
	if z.isTracking {
		z.tracking = append(z.tracking, b)
	}
	return b
}

func (z *quickReaderStream) readn1eof() (b uint8, eof bool) {
	b, err := z.br.ReadByte()
	if err == nil {
		z.n++
		if z.isTracking {
			z.tracking = append(z.tracking, b)
		}
	} else if err == io.EOF {
		eof = true
	} else {
		panic(err)
	}
	return
}

func (z *quickReaderStream) unreadn1() {
	err := z.br.UnreadByte()
	if err != nil {
		panic(err)
	}
	z.n--
	if z.isTracking {
		if l := len(z.tracking) - 1; l >= 0 {
			z.tracking = z.tracking[:l]
		}
	}
}

func (z *quickReaderStream) track() {
	if z.tracking != nil {
		z.tracking = z.tracking[:0]
	}
	z.isTracking = true
}

func (z *quickReaderStream) stopTrack() (bs []byte) {
	z.isTracking = false
	return z.tracking
}

// quickReaderSlice implements quickReader by reading a byte slice directly.
// Often this means the zero-copy methods can simply return subslices.
type quickReaderSlice struct {
	b []byte // data
	c int    // cursor
	a int    // available
	t int    // track start
}

func (z *quickReaderSlice) reset(in []byte) {
	z.b = in
	z.a = len(in)
	z.c = 0
	z.t = 0
}

func (z *quickReaderSlice) numread() int {
	return z.c
}

func (z *quickReaderSlice) unreadn1() {
	if z.c == 0 || len(z.b) == 0 {
		panic(errors.New("cannot unread last byte read"))
	}
	z.c--
	z.a++
	return
}

func (z *quickReaderSlice) readnzc(n int) (bs []byte) {
	if n == 0 {
		return zeroByteSlice
	} else if z.a == 0 {
		panic(io.EOF)
	} else if n > z.a {
		panic(io.ErrUnexpectedEOF)
	} else {
		c0 := z.c
		z.c = c0 + n
		z.a = z.a - n
		bs = z.b[c0:z.c]
	}
	return
}

func (z *quickReaderSlice) readn(n int) (bs []byte) {
	if n == 0 {
		return zeroByteSlice
	}
	bs = make([]byte, n)
	z.readb(bs)
	return
}

func (z *quickReaderSlice) readn1() (v uint8) {
	if z.a == 0 {
		panic(io.EOF)
	}
	v = z.b[z.c]
	z.c++
	z.a--
	return
}

func (z *quickReaderSlice) readn1eof() (v uint8, eof bool) {
	if z.a == 0 {
		eof = true
		return
	}
	v = z.b[z.c]
	z.c++
	z.a--
	return
}

func (z *quickReaderSlice) readb(bs []byte) {
	copy(bs, z.readnzc(len(bs)))
}

func (z *quickReaderSlice) track() {
	z.t = z.c
}

func (z *quickReaderSlice) stopTrack() (bs []byte) {
	return z.b[z.t:z.c]
}

// readerByteScanner decorates an `io.Reader` with all the methods to also
// fulfill the `io.ByteScanner` interface.
type readerByteScanner struct {
	r  io.Reader
	l  byte    // last byte
	ls byte    // last byte status. 0: init-canDoNothing, 1: canRead, 2: canUnread
	b  [1]byte // tiny buffer for reading single bytes
}

func (z *readerByteScanner) Read(p []byte) (n int, err error) {
	var firstByte bool
	if z.ls == 1 {
		z.ls = 2
		p[0] = z.l
		if len(p) == 1 {
			n = 1
			return
		}
		firstByte = true
		p = p[1:]
	}
	n, err = z.r.Read(p)
	if n > 0 {
		if err == io.EOF && n == len(p) {
			err = nil // read was successful, so postpone EOF (till next time)
		}
		z.l = p[n-1]
		z.ls = 2
	}
	if firstByte {
		n++
	}
	return
}

func (z *readerByteScanner) ReadByte() (c byte, err error) {
	n, err := z.Read(z.b[:])
	if n == 1 {
		c = z.b[0]
		if err == io.EOF {
			err = nil // read was successful, so postpone EOF (till next time)
		}
	}
	return
}

func (z *readerByteScanner) UnreadByte() (err error) {
	x := z.ls
	if x == 0 {
		err = errors.New("cannot unread - nothing has been read")
	} else if x == 1 {
		err = errors.New("cannot unread - last byte has not been read")
	} else if x == 2 {
		z.ls = 1
	}
	return
}
