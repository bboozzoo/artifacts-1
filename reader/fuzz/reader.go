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

		if up.Type == "" {
			panic("update type empty")
		}
		fls := w.GetUpdateFiles()
		if len(fls) == 0 {
			panic("no update files")
		}

		for _, f := range fls {
			if f.Checksum == nil || len(f.Checksum) == 0 {
				panic("no checksum")
			}

			if f.Path == "" {
				panic("empty path?")
			}

			if f.Size == 0 {
				panic("zero size")
			}
		}

		m := w.GetMetadata()
		if m == nil || len(*m) == 0 {
			panic("metadata nil or empty")
		}

		for k, v := range *m {
			if k == "" {
				panic("metadata key empty")
			}
			if v == nil {
				panic("metadata value nil")
			}
		}
	}

	// try reading again
	workers, err = ar.Read()
	if err != nil {
		panic("tried to read again and got no error")
	}
	return 1
}
