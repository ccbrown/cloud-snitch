package model

import (
	"crypto/rand"

	"github.com/jxskiss/base62"
)

type Id string

func (id Id) String() string {
	return string(id)
}

func NewId(namespace string) Id {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		panic(err)
	}
	s := base62.EncodeToString(buf)
	if len(s) > 22 {
		s = s[:22]
	}
	return Id(namespace + "-" + s)
}
