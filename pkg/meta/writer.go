package meta

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"image"
	_ "image/png"
	"io"
)

type LSBExtractor struct {
	data       [][][4]uint8
	rows, cols int
	bits       int
	byteBuf    byte
	row, col   int
}

func NewLSBExtractor(img image.Image) *LSBExtractor {
	bounds := img.Bounds()
	cols, rows := bounds.Max.X, bounds.Max.Y
	data := make([][][4]uint8, rows)
	for y := 0; y < rows; y++ {
		data[y] = make([][4]uint8, cols)
		for x := 0; x < cols; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			data[y][x] = [4]uint8{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
		}
	}
	return &LSBExtractor{data: data, rows: rows, cols: cols}
}

func (e *LSBExtractor) extractNextBit() bool {
	if e.row < e.rows && e.col < e.cols {
		bit := e.data[e.row][e.col][3] & 1
		e.bits++
		e.byteBuf <<= 1
		e.byteBuf |= bit
		e.row++
		if e.row == e.rows {
			e.row = 0
			e.col++
		}
		return true
	}
	return false
}

func (e *LSBExtractor) getOneByte() (byte, bool) {
	for e.bits < 8 {
		if !e.extractNextBit() {
			return 0, false
		}
	}
	b := e.byteBuf
	e.bits = 0
	e.byteBuf = 0
	return b, true
}

func (e *LSBExtractor) Read(p []byte) (int, error) {
	n := 0
	for n < len(p) {
		b, ok := e.getOneByte()
		if !ok {
			if n == 0 {
				return 0, io.EOF
			}
			break
		}
		p[n] = b
		n++
	}
	return n, nil
}

func ExtractFromBytes(r io.Reader) (*Metadata, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	return ExtractMetadata(img)
}

func ExtractMetadata(img image.Image) (*Metadata, error) {
	extractor := NewLSBExtractor(img)
	magic := "stealth_pngcomp"
	magicBytes := make([]byte, len(magic))
	_, err := io.ReadFull(extractor, magicBytes)
	if err != nil {
		return nil, err
	}
	if magic != string(magicBytes) {
		return nil, fmt.Errorf("magic number mismatch")
	}

	lenBytes := make([]byte, 4)
	_, err = io.ReadFull(extractor, lenBytes)
	if err != nil {
		return nil, err
	}
	readLen := int(binary.BigEndian.Uint32(lenBytes)) / 8

	jsonData := make([]byte, readLen)
	_, err = io.ReadFull(extractor, jsonData)
	if err != nil {
		return nil, err
	}

	reader, err := gzip.NewReader(bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	decompressedData, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var metadata struct {
		Metadata
		CommentString *string `json:"Comment,omitempty"`
	}
	err = json.Unmarshal(decompressedData, &metadata)
	if err != nil {
		return nil, err
	}

	if metadata.CommentString != nil {
		err = json.Unmarshal([]byte(*metadata.CommentString), &metadata.Metadata.Comment)
		if err != nil {
			return nil, err
		}
		metadata.raw = &raw{
			comment: metadata.CommentString,
			image:   img,
		}
	}

	return &metadata.Metadata, nil
}
