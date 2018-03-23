package ps

import (
	"encoding/hex"
)

var Colors map[string]Color = map[string]Color{
	"Gray":  &RGB{128, 128, 128},
	"White": &RGB{255, 255, 255},
}

// Color represents a color.
type Color interface {
	RGB() [3]int
}

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

// Color is a color in RGB format.
type RGB struct {
	Red   int
	Green int
	Blue  int
}

// RGB returns the color in RGB format.
func (r RGB) RGB() [3]int {
	return [3]int{r.Red, r.Green, r.Blue}
}

// Hex is a color in hexidecimal format.
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

// Stroke represents a layer stroke effect.
type Stroke struct {
	Size float32
	Color
}
