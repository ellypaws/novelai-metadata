package meta

import (
	"compress/gzip"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"image"
	_ "image/png"
	"io"
	"strings"
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

func (e *LSBExtractor) extractNextBit() {
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
	}
}

func (e *LSBExtractor) getOneByte() byte {
	for e.bits < 8 {
		e.extractNextBit()
	}
	b := e.byteBuf
	e.bits = 0
	e.byteBuf = 0
	return b
}

func (e *LSBExtractor) getNextNBytes(n int) []byte {
	bytesList := make([]byte, n)
	for i := 0; i < n; i++ {
		bytesList[i] = e.getOneByte()
	}
	return bytesList
}

func (e *LSBExtractor) read32BitInteger() int {
	bytesList := e.getNextNBytes(4)
	return int(binary.BigEndian.Uint32(bytesList))
}

func ExtractMetadata(r io.Reader) (*Metadata, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}

	extractor := NewLSBExtractor(img)
	magic := "stealth_pngcomp"
	readMagic := string(extractor.getNextNBytes(len(magic)))
	if magic != readMagic {
		return nil, fmt.Errorf("magic number mismatch")
	}

	readLen := extractor.read32BitInteger() / 8
	jsonData := extractor.getNextNBytes(readLen)

	reader, err := gzip.NewReader(strings.NewReader(string(jsonData)))
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
