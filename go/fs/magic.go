package fs

import (
	"bytes"
	"fmt"
)

type Magic []byte

const MagicSize = 4

func (m Magic) String() string {
	items := []struct {
		m    Magic
		name string
	}{
		{MagicEOF, "EOF"},
		{MagicInode, "inode"},
	}

	for _, i := range items {
		if bytes.Equal(i.m, m) {
			return i.name
		}
	}

	return fmt.Sprintf("unknown: %x", []byte(m))
}

var (
	MagicEOF   = Magic{0x8a, 0x9b, 0x0, 0x1}
	MagicInode = Magic{0x8a, 0x9b, 0x0, 0x2}
)
