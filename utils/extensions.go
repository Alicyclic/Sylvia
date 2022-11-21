package utils

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"regexp"
	"strings"

	"github.com/cenkalti/dominantcolor"
	"github.com/thoas/go-funk"
	"github.com/valyala/fasthttp"
)

type ImageDominateColor struct {
	Hex string `json:"hex"`
	RGB int    `json:"rgb"`
}

func GrabDominateColor(url string) color.RGBA {
	img := DecodeImageFromURL(url)
	dominantcolor := dominantcolor.Find(img)
	if funk.IsEmpty(dominantcolor) {
		return color.RGBA{0, 0, 0, 255}
	}
	return dominantcolor
}

func DominateColorToInt(c color.RGBA) int {
	return int(c.R)<<16 | int(c.G)<<8 | int(c.B)
}

func GetHexFromRGB(r, g, b int) string {
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

func GetColorInfo(url string) ImageDominateColor {
	dominantcolor := GrabDominateColor(url)
	return ImageDominateColor{
		RGB: DominateColorToInt(dominantcolor),
		Hex: GetHexFromRGB(int(dominantcolor.R), int(dominantcolor.G), int(dominantcolor.B)),
	}
}

func CheckIfError(err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in CheckIfError", r)
		}
	}()
}

// Apparently, this works?
// TODO: JPEG doesn't work?
func IsImage(check string) bool {
	regexString := "[^\\s]+(.*?)\\.(jpg|jpeg|png|gif|jpeg|JPG|PNG|GIF)$"
	regexp := regexp.MustCompile(regexString)
	return regexp.MatchString(check)
}

func DecodeImageFromURL(url string) image.Image {
	_, resp, err := fasthttp.Get(nil, url)
	CheckIfError(err)
	img, _, err := image.Decode(strings.NewReader(string(resp)))
	CheckIfError(err)
	return img
}
