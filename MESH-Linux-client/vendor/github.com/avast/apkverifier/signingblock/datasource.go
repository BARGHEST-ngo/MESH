package signingblock

import (
	"errors"
	"io"
)

type dataSource interface {
	chunkCount() int64
	length() int64
	writeTo(w io.Writer, offset, size int64) error
}

type dataSourceApk struct {
	file       io.ReadSeeker
	start, end int64
}

func (se *dataSourceApk) chunkCount() int64 {
	return (se.end - se.start + maxChunkSize - 1) / maxChunkSize
}

func (se *dataSourceApk) writeTo(w io.Writer, offset, size int64) error {
	if offset > se.end || offset > se.end-se.start {
		return errors.New("Out of bounds offset")
	} else if size > se.end-se.start || offset+size > se.end-se.start {
		return errors.New("Out of bounds size")
	}

	if _, err := se.file.Seek(se.start+offset, io.SeekStart); err != nil {
		return err
	}

	_, err := io.CopyN(w, se.file, size)
	return err
}

func (se *dataSourceApk) length() int64 {
	return se.end - se.start
}

type dataSourceBytes struct {
	data []byte
}

func (se *dataSourceBytes) chunkCount() int64 {
	return (int64(len(se.data)) + maxChunkSize - 1) / maxChunkSize
}

func (se *dataSourceBytes) writeTo(w io.Writer, offset, size int64) error {
	if offset >= int64(len(se.data)) {
		return errors.New("Out of bounds offset")
	} else if size > int64(len(se.data)) || offset+size > int64(len(se.data)) {
		return errors.New("Out of bounds size")
	}
	_, err := w.Write(se.data[offset : offset+size])
	return err
}

func (se *dataSourceBytes) length() int64 {
	return int64(len(se.data))
}

type dataSourceChained struct {
	sources   []dataSource
	totalSize int64
}

func newChainedDataSource(sources ...dataSource) dataSource {
	res := &dataSourceChained{
		sources: sources,
	}

	for _, cnt := range sources {
		res.totalSize += cnt.length()
	}
	return res
}

func (se *dataSourceChained) chunkCount() int64 {
	return (se.totalSize + maxChunkSize - 1) / maxChunkSize
}

func (se *dataSourceChained) writeTo(w io.Writer, offset, size int64) error {
	if offset >= se.totalSize {
		return errors.New("Out of bounds offset")
	} else if size > se.totalSize || offset+size > se.totalSize {
		return errors.New("Out of bounds size")
	}

	for _, src := range se.sources {
		if offset >= src.length() {
			offset -= src.length()
			continue
		}

		remaining := src.length() - offset
		if remaining >= size {
			return src.writeTo(w, offset, size)
		}

		if err := src.writeTo(w, offset, remaining); err != nil {
			return err
		}
		size -= remaining
		offset = 0
	}
	return nil
}

func (se *dataSourceChained) length() int64 {
	return se.totalSize
}

type byteWriter struct {
	dest   []byte
	offset int
}

func (w *byteWriter) Write(p []byte) (n int, err error) {
	if w.offset >= len(w.dest) {
		return 0, io.EOF
	}

	n = len(w.dest) - w.offset
	if n >= len(p) {
		n = len(p)
	} else {
		err = io.EOF
	}

	copy(w.dest[w.offset:], p[:n])
	w.offset += n
	return
}
