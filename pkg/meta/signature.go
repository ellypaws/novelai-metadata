package meta

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"image"
	"strings"
)

const VerifyKeyHex = "Y2JcQAOhLwzwSDUJPNgL04nS0Tbqm7cSRc4xk0vRMic="

func IsNovelAI(metadata Metadata) (bool, error) {
	if metadata.Comment == nil {
		return false, nil
	}

	if metadata.Comment.SignedHash == nil {
		return false, nil
	}

	if metadata.raw == nil {
		return false, fmt.Errorf("raw is nil")
	}

	signedHash, err := base64.StdEncoding.DecodeString(*metadata.Comment.SignedHash)
	if err != nil {
		return false, fmt.Errorf("failed to decode signed_hash: %w", err)
	}

	verifyKey, err := base64.StdEncoding.DecodeString(VerifyKeyHex)
	if err != nil {
		return false, fmt.Errorf("failed to decode verify key: %w", err)
	}

	var imageAndComment bytes.Buffer
	imageAndComment.Write(rgbaImageBytes(metadata.raw.image.(*image.NRGBA)))

	removeSignedHashField(metadata.raw.comment)

	buf := bytes.Buffer{}
	imageBytes := rgbaImageBytes(metadata.raw.image.(*image.NRGBA))

	buf.Write(imageBytes)
	buf.Write([]byte(*metadata.raw.comment))

	if !ed25519.Verify(verifyKey, buf.Bytes(), signedHash) {
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

func rgbaImageBytes(rgbaImg *image.NRGBA) []byte {
	rgbBytes := make([]byte, 0, 3*rgbaImg.Rect.Dx()*rgbaImg.Rect.Dy())
	for y := rgbaImg.Rect.Min.Y; y < rgbaImg.Rect.Max.Y; y++ {
		for x := rgbaImg.Rect.Min.X; x < rgbaImg.Rect.Max.X; x++ {
			pix := rgbaImg.Pix[(y-rgbaImg.Rect.Min.Y)*rgbaImg.Stride+(x-rgbaImg.Rect.Min.X)*4 : (y-rgbaImg.Rect.Min.Y)*rgbaImg.Stride+(x-rgbaImg.Rect.Min.X)*4+4]
			rgbBytes = append(rgbBytes, pix[0], pix[1], pix[2])
		}
	}
	return rgbBytes
}
