package fuzz

import (
	"bytes"

	areader "github.com/mendersoftware/mender-artifact/reader"
)

func Fuzz(data []byte) int {
	dr := bytes.NewReader(data)
	ar := areader.NewReader(dr)
	workers, err := ar.Read()
	if err != nil {
		return 0
	}

	if len(workers) == 0 {
		return 1
	}

	for _, w := range workers {
		up := w.GetUpdateType()
		if up == nil {
			panic("update type nil")
		}
		fls := w.GetUpdateFiles()
		if len(fls) == 0 {
			panic("no update files")
		}

		m := w.GetMetadata()
		if m == nil {
			panic("metadata nil")
		}
	}

	return 1
}
