package meta

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"strings"
	"sync"
)

const VerifyKeyHex = "Y2JcQAOhLwzwSDUJPNgL04nS0Tbqm7cSRc4xk0vRMic="

var (
	verifyKey []byte
	once      sync.Once
)

func initVerifyKey() {
	var err error
	verifyKey, err = base64.StdEncoding.DecodeString(VerifyKeyHex)
	if err != nil {
		panic(fmt.Sprintf("failed to decode verify key: %v", err))
	}
}

func (metadata *Metadata) IsNovelAI() (bool, error) {
	if metadata.Comment == nil {
		return false, nil
	}

	if metadata.Comment.SignedHash == nil {
		return false, nil
	}

	if metadata.raw == nil {
		return false, errors.New("raw is nil")
	}

	signedHash, err := base64.StdEncoding.DecodeString(*metadata.Comment.SignedHash)
	if err != nil {
		return false, fmt.Errorf("failed to decode signed_hash: %w", err)
	}

	once.Do(initVerifyKey)

	removeSignedHashField(metadata.raw.comment)

	bin := rgbaImageBytes(metadata.raw.image.(*image.NRGBA))
	bin.Grow(len(*metadata.raw.comment))
	bin.Write([]byte(*metadata.raw.comment))

	if !ed25519.Verify(verifyKey, bin.Bytes(), signedHash) {
		return false, nil
	}

	return true, nil
}

func removeSignedHashField(comment *string) {
	signedHashIndex := strings.LastIndex(*comment, `"signed_hash":`)
	if signedHashIndex == -1 {
		return
	}

	*comment = strings.TrimRight((*comment)[:signedHashIndex], ` ,}`)
	*comment = *comment + `}`
}

func rgbaImageBytes(rgbaImg *image.NRGBA) *bytes.Buffer {
	var bin bytes.Buffer
	bin.Grow(3 * rgbaImg.Rect.Dx() * rgbaImg.Rect.Dy())

	for y := rgbaImg.Rect.Min.Y; y < rgbaImg.Rect.Max.Y; y++ {
		for x := rgbaImg.Rect.Min.X; x < rgbaImg.Rect.Max.X; x++ {
			stride := (y-rgbaImg.Rect.Min.Y)*rgbaImg.Stride + (x-rgbaImg.Rect.Min.X)*4
			pix := rgbaImg.Pix[stride : stride+4]
			bin.Write(pix[:3])
		}
	}

	return &bin
}
