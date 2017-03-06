package cbor

import (
	"encoding/binary"
	"math"
)

func (d *Encoder) emitLen(majorByte byte, length int) {
	d.emitMajorPlusLen(majorByte, uint64(length))
}

func (d *Encoder) emitMajorPlusLen(majorByte byte, v uint64) {
	if v <= 0x17 {
		d.w.writen1(majorByte + byte(v))
	} else if v <= math.MaxUint8 {
		d.w.writen2(majorByte+0x18, uint8(v))
	} else if v <= math.MaxUint16 {
		d.w.writen1(majorByte + 0x19)
		d.spareBytes = d.spareBytes[:2]
		binary.BigEndian.PutUint16(d.spareBytes, uint16(v))
	} else if v <= math.MaxUint32 {
		d.w.writen1(majorByte + 0x1a)
		d.spareBytes = d.spareBytes[:4]
		binary.BigEndian.PutUint32(d.spareBytes, uint32(v))
	} else { // if v <= math.MaxUint64 {
		d.w.writen1(majorByte + 0x1b)
		d.spareBytes = d.spareBytes[:8]
		binary.BigEndian.PutUint64(d.spareBytes, v)
		d.w.writeb(d.spareBytes)
	}
}

func (d *Encoder) encodeNull() {
	d.w.writen1(cborSigilNil)
}

func (d *Encoder) encodeString(s string) {
	bs := []byte(s)
	d.emitMajorPlusLen(cborMajorString, uint64(len(bs)))
	d.w.writeb(bs)
}

func (d *Encoder) encodeBool(b bool) {
	if b {
		d.w.writen1(cborSigilTrue)
	} else {
		d.w.writen1(cborSigilFalse)
	}
}

func (d *Encoder) encodeInt64(v int64) {
	if v >= 0 {
		d.emitMajorPlusLen(cborMajorUint, uint64(v))
	} else {
		d.emitMajorPlusLen(cborMajorNegInt, uint64(-1-v))
	}
}

func (d *Encoder) encodeUint64(v uint64) {
	d.emitMajorPlusLen(cborMajorUint, v)
}
