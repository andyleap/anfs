package proto

import (
	"io"
)

type ReadReq struct {
	Key []byte
}

func (d *ReadReq) Size() (s uint64) {

	{
		l := uint64(len(d.Key))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		s += l
	}
	return
}
func (d *ReadReq) Marshal(buf []byte) ([]byte, error) {
	size := d.Size()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	{
		l := uint64(len(d.Key))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+0] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+0] = byte(t)
			i++

		}
		copy(buf[i+0:], d.Key)
		i += l
	}
	return buf[:i+0], nil
}

func (d *ReadReq) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+0] & 0x7F)
			for buf[i+0]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+0]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		if uint64(cap(d.Key)) >= l {
			d.Key = d.Key[:l]
		} else {
			d.Key = make([]byte, l)
		}
		copy(d.Key, buf[i+0:])
		i += l
	}
	return i + 0, nil
}

type ReadResp struct {
	Error string
	Data  []byte
}

func (d *ReadResp) Size() (s uint64) {

	{
		l := uint64(len(d.Error))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		s += l
	}
	{
		l := uint64(len(d.Data))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		s += l
	}
	return
}
func (d *ReadResp) Marshal(buf []byte) ([]byte, error) {
	size := d.Size()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	{
		l := uint64(len(d.Error))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+0] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+0] = byte(t)
			i++

		}
		copy(buf[i+0:], d.Error)
		i += l
	}
	{
		l := uint64(len(d.Data))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+0] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+0] = byte(t)
			i++

		}
		copy(buf[i+0:], d.Data)
		i += l
	}
	return buf[:i+0], nil
}

func (d *ReadResp) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+0] & 0x7F)
			for buf[i+0]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+0]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		d.Error = string(buf[i+0 : i+0+l])
		i += l
	}
	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+0] & 0x7F)
			for buf[i+0]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+0]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		if uint64(cap(d.Data)) >= l {
			d.Data = d.Data[:l]
		} else {
			d.Data = make([]byte, l)
		}
		copy(d.Data, buf[i+0:])
		i += l
	}
	return i + 0, nil
}

type WriteReq struct {
	Key  []byte
	Data []byte
}

func (d *WriteReq) Size() (s uint64) {

	{
		l := uint64(len(d.Key))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		s += l
	}
	{
		l := uint64(len(d.Data))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		s += l
	}
	return
}
func (d *WriteReq) Marshal(buf []byte) ([]byte, error) {
	size := d.Size()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	{
		l := uint64(len(d.Key))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+0] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+0] = byte(t)
			i++

		}
		copy(buf[i+0:], d.Key)
		i += l
	}
	{
		l := uint64(len(d.Data))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+0] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+0] = byte(t)
			i++

		}
		copy(buf[i+0:], d.Data)
		i += l
	}
	return buf[:i+0], nil
}

func (d *WriteReq) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+0] & 0x7F)
			for buf[i+0]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+0]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		if uint64(cap(d.Key)) >= l {
			d.Key = d.Key[:l]
		} else {
			d.Key = make([]byte, l)
		}
		copy(d.Key, buf[i+0:])
		i += l
	}
	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+0] & 0x7F)
			for buf[i+0]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+0]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		if uint64(cap(d.Data)) >= l {
			d.Data = d.Data[:l]
		} else {
			d.Data = make([]byte, l)
		}
		copy(d.Data, buf[i+0:])
		i += l
	}
	return i + 0, nil
}

type WriteResp struct {
	Error string
}

func (d *WriteResp) Size() (s uint64) {

	{
		l := uint64(len(d.Error))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		s += l
	}
	return
}
func (d *WriteResp) Marshal(buf []byte) ([]byte, error) {
	size := d.Size()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	{
		l := uint64(len(d.Error))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+0] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+0] = byte(t)
			i++

		}
		copy(buf[i+0:], d.Error)
		i += l
	}
	return buf[:i+0], nil
}

func (d *WriteResp) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+0] & 0x7F)
			for buf[i+0]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+0]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		d.Error = string(buf[i+0 : i+0+l])
		i += l
	}
	return i + 0, nil
}

type DeleteReq struct {
	Key []byte
}

func (d *DeleteReq) Size() (s uint64) {

	{
		l := uint64(len(d.Key))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		s += l
	}
	return
}
func (d *DeleteReq) Marshal(buf []byte) ([]byte, error) {
	size := d.Size()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	{
		l := uint64(len(d.Key))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+0] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+0] = byte(t)
			i++

		}
		copy(buf[i+0:], d.Key)
		i += l
	}
	return buf[:i+0], nil
}

func (d *DeleteReq) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+0] & 0x7F)
			for buf[i+0]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+0]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		if uint64(cap(d.Key)) >= l {
			d.Key = d.Key[:l]
		} else {
			d.Key = make([]byte, l)
		}
		copy(d.Key, buf[i+0:])
		i += l
	}
	return i + 0, nil
}

type DeleteResp struct {
	Error string
}

func (d *DeleteResp) Size() (s uint64) {

	{
		l := uint64(len(d.Error))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		s += l
	}
	return
}
func (d *DeleteResp) Marshal(buf []byte) ([]byte, error) {
	size := d.Size()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	{
		l := uint64(len(d.Error))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+0] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+0] = byte(t)
			i++

		}
		copy(buf[i+0:], d.Error)
		i += l
	}
	return buf[:i+0], nil
}

func (d *DeleteResp) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+0] & 0x7F)
			for buf[i+0]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+0]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		d.Error = string(buf[i+0 : i+0+l])
		i += l
	}
	return i + 0, nil
}

type ClearReq struct {
}

func (d *ClearReq) Size() (s uint64) {

	return
}
func (d *ClearReq) Marshal(buf []byte) ([]byte, error) {
	size := d.Size()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	return buf[:i+0], nil
}

func (d *ClearReq) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	return i + 0, nil
}

type ClearResp struct {
}

func (d *ClearResp) Size() (s uint64) {

	return
}
func (d *ClearResp) Marshal(buf []byte) ([]byte, error) {
	size := d.Size()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	return buf[:i+0], nil
}

func (d *ClearResp) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	return i + 0, nil
}

type CheckReq struct {
	Key []byte
}

func (d *CheckReq) Size() (s uint64) {

	{
		l := uint64(len(d.Key))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		s += l
	}
	return
}
func (d *CheckReq) Marshal(buf []byte) ([]byte, error) {
	size := d.Size()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	{
		l := uint64(len(d.Key))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+0] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+0] = byte(t)
			i++

		}
		copy(buf[i+0:], d.Key)
		i += l
	}
	return buf[:i+0], nil
}

func (d *CheckReq) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+0] & 0x7F)
			for buf[i+0]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+0]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		if uint64(cap(d.Key)) >= l {
			d.Key = d.Key[:l]
		} else {
			d.Key = make([]byte, l)
		}
		copy(d.Key, buf[i+0:])
		i += l
	}
	return i + 0, nil
}

type CheckResp struct {
	Found bool
}

func (d *CheckResp) Size() (s uint64) {

	s += 1
	return
}
func (d *CheckResp) Marshal(buf []byte) ([]byte, error) {
	size := d.Size()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	{
		if d.Found {
			buf[0] = 1
		} else {
			buf[0] = 0
		}
	}
	return buf[:i+1], nil
}

func (d *CheckResp) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	{
		d.Found = buf[0] == 1
	}
	return i + 1, nil
}

type SweepReq struct {
}

func (d *SweepReq) Size() (s uint64) {

	return
}
func (d *SweepReq) Marshal(buf []byte) ([]byte, error) {
	size := d.Size()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	return buf[:i+0], nil
}

func (d *SweepReq) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	return i + 0, nil
}

type SweepResp struct {
}

func (d *SweepResp) Size() (s uint64) {

	return
}
func (d *SweepResp) Marshal(buf []byte) ([]byte, error) {
	size := d.Size()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	return buf[:i+0], nil
}

func (d *SweepResp) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	return i + 0, nil
}

type Req struct {
	ID      uint64
	Request Request
}

func (d *Req) FramedSize() (s uint64, us uint64) {

	{

		t := d.ID
		for t >= 0x80 {
			t >>= 7
			s++
		}
		s++

	}
	{
		var v uint64
		switch d.Request.(type) {

		case ReadReq:
			v = 0 + 1

		case WriteReq:
			v = 1 + 1

		case DeleteReq:
			v = 2 + 1

		case ClearReq:
			v = 3 + 1

		case CheckReq:
			v = 4 + 1

		case SweepReq:
			v = 5 + 1

		}

		{

			t := v
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		switch tt := d.Request.(type) {

		case ReadReq:

			{
				s += tt.Size()
			}

		case WriteReq:

			{
				s += tt.Size()
			}

		case DeleteReq:

			{
				s += tt.Size()
			}

		case ClearReq:

			{
				s += tt.Size()
			}

		case CheckReq:

			{
				s += tt.Size()
			}

		case SweepReq:

			{
				s += tt.Size()
			}

		}
	}
	l := s
	us = s

	{

		t := l
		for t >= 0x80 {
			t >>= 7
			s++
		}
		s++

	}
	return
}
func (d *Req) Size() (s uint64) {
	s, _ = d.FramedSize()
	return
}

func (d *Req) Marshal(buf []byte) ([]byte, error) {
	size, usize := d.FramedSize()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	{

		t := uint64(usize)

		for t >= 0x80 {
			buf[i+0] = byte(t) | 0x80
			t >>= 7
			i++
		}
		buf[i+0] = byte(t)
		i++

	}
	{

		t := uint64(d.ID)

		for t >= 0x80 {
			buf[i+0] = byte(t) | 0x80
			t >>= 7
			i++
		}
		buf[i+0] = byte(t)
		i++

	}
	{
		var v uint64
		switch d.Request.(type) {

		case ReadReq:
			v = 0 + 1

		case WriteReq:
			v = 1 + 1

		case DeleteReq:
			v = 2 + 1

		case ClearReq:
			v = 3 + 1

		case CheckReq:
			v = 4 + 1

		case SweepReq:
			v = 5 + 1

		}

		{

			t := uint64(v)

			for t >= 0x80 {
				buf[i+0] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+0] = byte(t)
			i++

		}
		switch tt := d.Request.(type) {

		case ReadReq:

			{
				nbuf, err := tt.Marshal(buf[i+0:])
				if err != nil {
					return nil, err
				}
				i += uint64(len(nbuf))
			}

		case WriteReq:

			{
				nbuf, err := tt.Marshal(buf[i+0:])
				if err != nil {
					return nil, err
				}
				i += uint64(len(nbuf))
			}

		case DeleteReq:

			{
				nbuf, err := tt.Marshal(buf[i+0:])
				if err != nil {
					return nil, err
				}
				i += uint64(len(nbuf))
			}

		case ClearReq:

			{
				nbuf, err := tt.Marshal(buf[i+0:])
				if err != nil {
					return nil, err
				}
				i += uint64(len(nbuf))
			}

		case CheckReq:

			{
				nbuf, err := tt.Marshal(buf[i+0:])
				if err != nil {
					return nil, err
				}
				i += uint64(len(nbuf))
			}

		case SweepReq:

			{
				nbuf, err := tt.Marshal(buf[i+0:])
				if err != nil {
					return nil, err
				}
				i += uint64(len(nbuf))
			}

		}
	}
	return buf[:i+0], nil
}

func (d *Req) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)
	usize := uint64(0)

	{

		bs := uint8(7)
		t := uint64(buf[i+0] & 0x7F)
		for buf[i+0]&0x80 == 0x80 {
			i++
			t |= uint64(buf[i+0]&0x7F) << bs
			bs += 7
		}
		i++

		usize = t

	}
	if usize > uint64(len(buf))+i {
		return 0, io.EOF
	}
	{

		bs := uint8(7)
		t := uint64(buf[i+0] & 0x7F)
		for buf[i+0]&0x80 == 0x80 {
			i++
			t |= uint64(buf[i+0]&0x7F) << bs
			bs += 7
		}
		i++

		d.ID = t

	}
	{
		v := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+0] & 0x7F)
			for buf[i+0]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+0]&0x7F) << bs
				bs += 7
			}
			i++

			v = t

		}
		switch v {

		case 0 + 1:
			var tt ReadReq

			{
				ni, err := tt.Unmarshal(buf[i+0:])
				if err != nil {
					return 0, err
				}
				i += ni
			}

			d.Request = tt

		case 1 + 1:
			var tt WriteReq

			{
				ni, err := tt.Unmarshal(buf[i+0:])
				if err != nil {
					return 0, err
				}
				i += ni
			}

			d.Request = tt

		case 2 + 1:
			var tt DeleteReq

			{
				ni, err := tt.Unmarshal(buf[i+0:])
				if err != nil {
					return 0, err
				}
				i += ni
			}

			d.Request = tt

		case 3 + 1:
			var tt ClearReq

			{
				ni, err := tt.Unmarshal(buf[i+0:])
				if err != nil {
					return 0, err
				}
				i += ni
			}

			d.Request = tt

		case 4 + 1:
			var tt CheckReq

			{
				ni, err := tt.Unmarshal(buf[i+0:])
				if err != nil {
					return 0, err
				}
				i += ni
			}

			d.Request = tt

		case 5 + 1:
			var tt SweepReq

			{
				ni, err := tt.Unmarshal(buf[i+0:])
				if err != nil {
					return 0, err
				}
				i += ni
			}

			d.Request = tt

		default:
			d.Request = nil
		}
	}
	return i + 0, nil
}

func (d *Req) Serialize(w io.Writer) error {
	buf, err := d.Marshal(nil)
	if err != nil {
		return err
	}
	_, err = w.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func (d *Req) Deserialize(r io.Reader) error {
	size := uint64(0)
	sbuf := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	bs := uint8(0)
	i := uint64(0)
	_, err := r.Read(sbuf[i : i+1])
	if err != nil {
		return err
	}
	size |= uint64(sbuf[i]&0x7F) << bs
	bs += 7
	for sbuf[i]&0x80 == 0x80 {
		i++
		_, err = r.Read(sbuf[i : i+1])
		if err != nil {
			return err
		}
		size |= uint64(sbuf[i]&0x7F) << bs
		bs += 7
	}
	i++
	buf := make([]byte, size+i)
	copy(buf, sbuf[0:i])
	n := uint64(i)
	size += i
	for n < size && err == nil {
		var nn int
		nn, err = r.Read(buf[n:])
		n += uint64(nn)
	}
	if err != nil {
		return err
	}
	_, err = d.Unmarshal(buf)
	if err != nil {
		return err
	}
	return nil
}

type Resp struct {
	ID       uint64
	Response Response
}

func (d *Resp) FramedSize() (s uint64, us uint64) {

	{

		t := d.ID
		for t >= 0x80 {
			t >>= 7
			s++
		}
		s++

	}
	{
		var v uint64
		switch d.Response.(type) {

		case ReadResp:
			v = 0 + 1

		case WriteResp:
			v = 1 + 1

		case DeleteResp:
			v = 2 + 1

		case ClearResp:
			v = 3 + 1

		case CheckResp:
			v = 4 + 1

		case SweepResp:
			v = 5 + 1

		}

		{

			t := v
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		switch tt := d.Response.(type) {

		case ReadResp:

			{
				s += tt.Size()
			}

		case WriteResp:

			{
				s += tt.Size()
			}

		case DeleteResp:

			{
				s += tt.Size()
			}

		case ClearResp:

			{
				s += tt.Size()
			}

		case CheckResp:

			{
				s += tt.Size()
			}

		case SweepResp:

			{
				s += tt.Size()
			}

		}
	}
	l := s
	us = s

	{

		t := l
		for t >= 0x80 {
			t >>= 7
			s++
		}
		s++

	}
	return
}
func (d *Resp) Size() (s uint64) {
	s, _ = d.FramedSize()
	return
}

func (d *Resp) Marshal(buf []byte) ([]byte, error) {
	size, usize := d.FramedSize()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	{

		t := uint64(usize)

		for t >= 0x80 {
			buf[i+0] = byte(t) | 0x80
			t >>= 7
			i++
		}
		buf[i+0] = byte(t)
		i++

	}
	{

		t := uint64(d.ID)

		for t >= 0x80 {
			buf[i+0] = byte(t) | 0x80
			t >>= 7
			i++
		}
		buf[i+0] = byte(t)
		i++

	}
	{
		var v uint64
		switch d.Response.(type) {

		case ReadResp:
			v = 0 + 1

		case WriteResp:
			v = 1 + 1

		case DeleteResp:
			v = 2 + 1

		case ClearResp:
			v = 3 + 1

		case CheckResp:
			v = 4 + 1

		case SweepResp:
			v = 5 + 1

		}

		{

			t := uint64(v)

			for t >= 0x80 {
				buf[i+0] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+0] = byte(t)
			i++

		}
		switch tt := d.Response.(type) {

		case ReadResp:

			{
				nbuf, err := tt.Marshal(buf[i+0:])
				if err != nil {
					return nil, err
				}
				i += uint64(len(nbuf))
			}

		case WriteResp:

			{
				nbuf, err := tt.Marshal(buf[i+0:])
				if err != nil {
					return nil, err
				}
				i += uint64(len(nbuf))
			}

		case DeleteResp:

			{
				nbuf, err := tt.Marshal(buf[i+0:])
				if err != nil {
					return nil, err
				}
				i += uint64(len(nbuf))
			}

		case ClearResp:

			{
				nbuf, err := tt.Marshal(buf[i+0:])
				if err != nil {
					return nil, err
				}
				i += uint64(len(nbuf))
			}

		case CheckResp:

			{
				nbuf, err := tt.Marshal(buf[i+0:])
				if err != nil {
					return nil, err
				}
				i += uint64(len(nbuf))
			}

		case SweepResp:

			{
				nbuf, err := tt.Marshal(buf[i+0:])
				if err != nil {
					return nil, err
				}
				i += uint64(len(nbuf))
			}

		}
	}
	return buf[:i+0], nil
}

func (d *Resp) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)
	usize := uint64(0)

	{

		bs := uint8(7)
		t := uint64(buf[i+0] & 0x7F)
		for buf[i+0]&0x80 == 0x80 {
			i++
			t |= uint64(buf[i+0]&0x7F) << bs
			bs += 7
		}
		i++

		usize = t

	}
	if usize > uint64(len(buf))+i {
		return 0, io.EOF
	}
	{

		bs := uint8(7)
		t := uint64(buf[i+0] & 0x7F)
		for buf[i+0]&0x80 == 0x80 {
			i++
			t |= uint64(buf[i+0]&0x7F) << bs
			bs += 7
		}
		i++

		d.ID = t

	}
	{
		v := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+0] & 0x7F)
			for buf[i+0]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+0]&0x7F) << bs
				bs += 7
			}
			i++

			v = t

		}
		switch v {

		case 0 + 1:
			var tt ReadResp

			{
				ni, err := tt.Unmarshal(buf[i+0:])
				if err != nil {
					return 0, err
				}
				i += ni
			}

			d.Response = tt

		case 1 + 1:
			var tt WriteResp

			{
				ni, err := tt.Unmarshal(buf[i+0:])
				if err != nil {
					return 0, err
				}
				i += ni
			}

			d.Response = tt

		case 2 + 1:
			var tt DeleteResp

			{
				ni, err := tt.Unmarshal(buf[i+0:])
				if err != nil {
					return 0, err
				}
				i += ni
			}

			d.Response = tt

		case 3 + 1:
			var tt ClearResp

			{
				ni, err := tt.Unmarshal(buf[i+0:])
				if err != nil {
					return 0, err
				}
				i += ni
			}

			d.Response = tt

		case 4 + 1:
			var tt CheckResp

			{
				ni, err := tt.Unmarshal(buf[i+0:])
				if err != nil {
					return 0, err
				}
				i += ni
			}

			d.Response = tt

		case 5 + 1:
			var tt SweepResp

			{
				ni, err := tt.Unmarshal(buf[i+0:])
				if err != nil {
					return 0, err
				}
				i += ni
			}

			d.Response = tt

		default:
			d.Response = nil
		}
	}
	return i + 0, nil
}

func (d *Resp) Serialize(w io.Writer) error {
	buf, err := d.Marshal(nil)
	if err != nil {
		return err
	}
	_, err = w.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func (d *Resp) Deserialize(r io.Reader) error {
	size := uint64(0)
	sbuf := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	bs := uint8(0)
	i := uint64(0)
	_, err := r.Read(sbuf[i : i+1])
	if err != nil {
		return err
	}
	size |= uint64(sbuf[i]&0x7F) << bs
	bs += 7
	for sbuf[i]&0x80 == 0x80 {
		i++
		_, err = r.Read(sbuf[i : i+1])
		if err != nil {
			return err
		}
		size |= uint64(sbuf[i]&0x7F) << bs
		bs += 7
	}
	i++
	buf := make([]byte, size+i)
	copy(buf, sbuf[0:i])
	n := uint64(i)
	size += i
	for n < size && err == nil {
		var nn int
		nn, err = r.Read(buf[n:])
		n += uint64(nn)
	}
	if err != nil {
		return err
	}
	_, err = d.Unmarshal(buf)
	if err != nil {
		return err
	}
	return nil
}
