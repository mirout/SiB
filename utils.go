package main

import (
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

const FILE_SIZE = 1 << 20

type BlockFileReader struct {
	Name       string
	blockCount int
	blockNum   int
	read       int
	file       []byte
}

func Open(filename string) (*BlockFileReader, error) {
	files, err := os.ReadDir(filename)
	if err != nil {
		return nil, err
	}
	return &BlockFileReader{Name: filename, blockNum: 0, blockCount: len(files), read: 0}, nil
}

func (r *BlockFileReader) openBlock() error {
	var err error
	r.file, err = ioutil.ReadFile(r.Name + "/" + strconv.FormatUint(uint64(r.blockNum), 16) + ".dat")
	if err != nil {
		return err
	}

	if r.blockNum+1 == r.blockCount {
		var i int
		for i = len(r.file) - 1; r.file[i] == 0; i-- {
		}
		r.file = r.file[:i+1]
	}

	return nil
}

func (r *BlockFileReader) Read(bytes []byte) (int, error) {
	if r.file != nil && FILE_SIZE-r.read == 0 {
		r.file = nil
		r.blockNum += 1
		r.read = 0
	}
	if r.file == nil {
		if r.blockCount == r.blockNum {
			return 0, io.EOF
		}
		if err := r.openBlock(); err != nil {
			return 0, err
		}
	}

	count := copy(bytes, r.file[r.read:])
	r.read += count

	if count == 0 {
		return 0, io.EOF
	}

	return count, nil
}

type BlockFileWriter struct {
	Name     string
	blockNum uint
	written  int
	file     *os.File
}

func Create(filename string) (*BlockFileWriter, error) {
	if err := os.Mkdir(filename, os.ModePerm); err != nil {
		return nil, err
	}
	return &BlockFileWriter{Name: filename, blockNum: 0, written: 0, file: nil}, nil
}

func (w *BlockFileWriter) createBlock() error {
	var err error
	w.file, err = os.Create(w.Name + "/" + strconv.FormatUint(uint64(w.blockNum), 16) + ".dat")
	return err
}

func (w *BlockFileWriter) Write(b []byte) (int, error) {
	if w.file != nil && FILE_SIZE-w.written == 0 {
		if err := w.file.Close(); err != nil {
			return 0, err
		}
		w.file = nil
		w.blockNum += 1
		w.written = 0
	}
	if w.file == nil {
		if err := w.createBlock(); err != nil {
			return 0, err
		}
	}

	canWrite := FILE_SIZE - w.written

	if canWrite < len(b) {
		if n, err := w.file.Write(b[:canWrite]); err != nil {
			w.written += n
			return n, err
		}
		w.written += canWrite
		return w.Write(b[canWrite:])
	} else {
		if n, err := w.file.Write(b); err != nil {
			w.written += n
			return n, err
		}
		w.written += len(b)
		return len(b), nil
	}
}

func (w *BlockFileWriter) Close() error {
	defer w.file.Close()
	if diff := FILE_SIZE - w.written; diff != 0 {
		_, err := w.file.Write(make([]byte, diff))
		if err != nil {
			return err
		}
	}
	return nil
}
