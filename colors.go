package ps

import "encoding/hex"

// Some basic colors.
var (
	ColorBlack Color = RGB{0, 0, 0}
	ColorGray  Color = RGB{128, 128, 128}
	ColorWhite Color = RGB{255, 255, 255}
)

// Color is an interface for color objects, allowing colors to be
// used in various formats.
//
// RGB is the default format for everything.
type Color interface {
	RGB() [3]int  // The color in RGB format.
	Hex() []uint8 // The color in hexadecimal format.
}

// Compare returns the brighter of a and b.
func Compare(a, b Color) Color {
	A := a.RGB()
	B := b.RGB()
	avgA := (A[0] + A[1] + A[2]) / 3
	avgB := (B[0] + B[1] + B[2]) / 3
	if avgA > avgB {
		return a
	}
	return b
}

// RGB is a color format. It implements the Color interface.
type RGB struct {
	Red   int
	Green int
	Blue  int
}

// RGB returns the color in RGB format.
func (r RGB) RGB() [3]int {
	return [3]int{r.Red, r.Green, r.Blue}
}

// Hex returns the color converted to hexadecimal format.
func (r RGB) Hex() []uint8 {
	src := []uint8{uint8(r.Red), uint8(r.Green), uint8(r.Blue)}
	hex := make([]byte, hex.EncodedLen(len(src)))
	return hex
}

// Hex is a color in hexadecimal format.
// It satisfies the Color interface.
type Hex []uint8

// RGB returns the Hex value converted to RGB
func (h Hex) RGB() [3]int {
	src := []byte(h)
	dst := make([]byte, hex.DecodedLen(len(src)))
	_, err := hex.Decode(dst, src)
	if err != nil {
		panic(err)
	}
	return [3]int{int(dst[0]), int(dst[1]), int(dst[2])}
}

// Hex returns the hex value of the number,
// to satisfy the Color interface.
func (h Hex) Hex() []uint8 {
	return h
}

// Stroke represents a layer stroke effect.
type Stroke struct {
	Size float32
	Color
}
