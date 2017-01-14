package fuzz

import (
	"bytes"

	areader "github.com/mendersoftware/mender-artifact/reader"
)

func Fuzz(data []byte) int {
	dr := bytes.NewReader(data)
	ar := areader.NewReader(dr)
	_, err := ar.Read()
	if err != nil {
		return 0
	}
	return 1
}
