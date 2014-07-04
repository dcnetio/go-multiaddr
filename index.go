package multiaddr

import (
	"fmt"
	"strings"
)

type Multiaddr struct {
	Bytes []byte
}

func NewMultiaddr(s string) (*Multiaddr, error) {
	b, err := StringToBytes(s)
	if err != nil {
		return nil, err
	}
	return &Multiaddr{Bytes: b}, nil
}

func (m *Multiaddr) String() (string, error) {
	return BytesToString(m.Bytes)
}

func (m *Multiaddr) Protocols() (ret []*Protocol, err error) {

	// panic handler, in case we try accessing bytes incorrectly.
	defer func() {
		if e := recover(); e != nil {
			ret = nil
			err = e.(error)
		}
	}()

	ps := []*Protocol{}
	b := m.Bytes[:]
	for len(b) > 0 {
		p := ProtocolWithCode(int(b[0]))
		if p == nil {
			return nil, fmt.Errorf("no protocol with code %d", b[0])
		}
		ps = append(ps, p)
		b = b[1+(p.Size/8):]
	}
	return ps, nil
}

func (m *Multiaddr) Encapsulate(o *Multiaddr) *Multiaddr {
	b := make([]byte, len(m.Bytes)+len(o.Bytes))
	b = append(m.Bytes, o.Bytes...)
	return &Multiaddr{Bytes: b}
}

func (m *Multiaddr) Decapsulate(o *Multiaddr) (*Multiaddr, error) {
	s1, err := m.String()
	if err != nil {
		return nil, err
	}

	s2, err := o.String()
	if err != nil {
		return nil, err
	}

	i := strings.LastIndex(s1, s2)
	if i < 0 {
		return nil, fmt.Errorf("%s not contained in %s", s2, s1)
	}
	return NewMultiaddr(s1[:i])
}