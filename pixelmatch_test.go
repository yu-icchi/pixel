package pixel

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"image"
	"image/png"
	"os"
	"testing"
)

func TestMatch(t *testing.T) {
	file1, err := os.Open("./testdata/4a.png")
	require.NoError(t, err)
	img1, err := png.Decode(file1)
	require.NoError(t, err)

	file2, err := os.Open("./testdata/4b.png")
	require.NoError(t, err)
	img2, err := png.Decode(file2)
	require.NoError(t, err)

	opt := &Option{
		Threshold: 0.1,
		IncludeAA: false,
	}
	outImg := image.NewRGBA(img2.Bounds())
	d := Match(img1, img2, img2.Bounds(), outImg, opt)
	fmt.Println(d)

	out, err := os.Create("./testdata/4diff.png")
	require.NoError(t, err)
	defer out.Close()

	err = png.Encode(out, outImg)
	require.NoError(t, err)
}

func BenchmarkMatch(b *testing.B) {
	file1, err := os.Open("./testdata/4a.png")
	require.NoError(b, err)
	img1, err := png.Decode(file1)
	require.NoError(b, err)

	file2, err := os.Open("./testdata/4b.png")
	require.NoError(b, err)
	img2, err := png.Decode(file2)
	require.NoError(b, err)

	opt := &Option{
		Threshold: 0.1,
		IncludeAA: false,
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		outImg := image.NewRGBA(img2.Bounds())
		Match(img1, img2, img2.Bounds(), outImg, opt)
	}
}
