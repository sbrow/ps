package ps

import (
	"encoding/hex"
	// "fmt"
)

// Color is an interface for color objects, allowing colors to be
// used in various formats.
//
// RGB is the default format for everything.
type Color interface {
	RGB() [3]int  // The color in RGB format.
	Hex() []uint8 // The color in hexadecimal format.
}

// Compare determines which of two colors is "brighter".
func Compare(a, b Color) Color {
	A := a.RGB()
	B := b.RGB()
	Aavg := (A[0] + A[1] + A[2]) / 3
	Bavg := (B[0] + B[1] + B[2]) / 3
	if Aavg > Bavg {
		return a
	}
	return b
}

// RGB is a color in RGB format. It fulfills the Color interface.
type RGB struct {
	Red   int
	Green int
	Blue  int
}

// RGB returns the color in RGB format.
func (r RGB) RGB() [3]int {
	return [3]int{r.Red, r.Green, r.Blue}
}

// TODO: Implement RGB.Hex()
func (r RGB) Hex() []uint8 {
	return make([]uint8, 6)
}

// Hex is a color in hexadecimal format. It fulfills the Color interface.
type Hex []uint8

func (h Hex) RGB() [3]int {
	src := []byte(h)
	dst := make([]byte, hex.DecodedLen(len(src)))
	_, err := hex.Decode(dst, src)
	if err != nil {
		panic(err)
	}
	return [3]int{int(dst[0]), int(dst[1]), int(dst[2])}
}

func (h Hex) Hex() []uint8 {
	return h
}

// Stroke represents a layer stroke effect.
type Stroke struct {
	Size float32
	Color
}

// func (s *Stroke) String() string {
// 	return fmt.Sprintf("%vpt %v", s.Size, s.Color.RGB())
// }
