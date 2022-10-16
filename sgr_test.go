package ansisgr_test

import (
	"fmt"
	"strings"
	"testing"

	ansisgr "github.com/ktr0731/go-ansisgr"
)

func TestIterator(t *testing.T) {
	t.Parallel()

	type want struct {
		r              rune
		foreground     int
		foregroundMode ansisgr.Mode
		background     int
		backgroundMode ansisgr.Mode

		bold, dim, italic, underline, blink, reverse, invisible, strikethrough bool
	}

	cases := []struct {
		name string
		in   string
		want []want
	}{
		{
			name: "no colors",
			in:   "foo",
			want: []want{{r: 'f'}, {r: 'o'}, {r: 'o'}},
		},
		{
			name: "16 colors",
			in:   "foo\x1b[30;41mbar\x1b[90;101mbaz",
			want: []want{
				{r: 'f'},
				{r: 'o'},
				{r: 'o'},
				{r: 'b', foreground: 30, foregroundMode: ansisgr.Mode16, background: 41, backgroundMode: ansisgr.Mode16},
				{r: 'a', foreground: 30, foregroundMode: ansisgr.Mode16, background: 41, backgroundMode: ansisgr.Mode16},
				{r: 'r', foreground: 30, foregroundMode: ansisgr.Mode16, background: 41, backgroundMode: ansisgr.Mode16},
				{r: 'b', foreground: 90, foregroundMode: ansisgr.Mode16, background: 101, backgroundMode: ansisgr.Mode16},
				{r: 'a', foreground: 90, foregroundMode: ansisgr.Mode16, background: 101, backgroundMode: ansisgr.Mode16},
				{r: 'z', foreground: 90, foregroundMode: ansisgr.Mode16, background: 101, backgroundMode: ansisgr.Mode16},
			},
		},
		{
			name: "16 colors which contain a broken sequence",
			in:   "foo\x1b[30;441mbar\x1b[90;101mbaz",
			want: []want{
				{r: 'f'},
				{r: 'o'},
				{r: 'o'},
				{r: 'b', foreground: 30, foregroundMode: ansisgr.Mode16},
				{r: 'a', foreground: 30, foregroundMode: ansisgr.Mode16},
				{r: 'r', foreground: 30, foregroundMode: ansisgr.Mode16},
				{r: 'b', foreground: 90, foregroundMode: ansisgr.Mode16, background: 101, backgroundMode: ansisgr.Mode16},
				{r: 'a', foreground: 90, foregroundMode: ansisgr.Mode16, background: 101, backgroundMode: ansisgr.Mode16},
				{r: 'z', foreground: 90, foregroundMode: ansisgr.Mode16, background: 101, backgroundMode: ansisgr.Mode16},
			},
		},
		{
			name: "256 colors",
			in:   "foo\x1b[38;5;117;48;5;104mbar",
			want: []want{
				{r: 'f'},
				{r: 'o'},
				{r: 'o'},
				{r: 'b', foreground: 117, foregroundMode: ansisgr.Mode256, background: 104, backgroundMode: ansisgr.Mode256},
				{r: 'a', foreground: 117, foregroundMode: ansisgr.Mode256, background: 104, backgroundMode: ansisgr.Mode256},
				{r: 'r', foreground: 117, foregroundMode: ansisgr.Mode256, background: 104, backgroundMode: ansisgr.Mode256},
			},
		},
		{
			name: "16 colors by 256 colors",
			in:   "foo\x1b[38;5;1;48;5;2mbar",
			want: []want{
				{r: 'f'},
				{r: 'o'},
				{r: 'o'},
				{r: 'b', foreground: 1, foregroundMode: ansisgr.Mode256, background: 2, backgroundMode: ansisgr.Mode256},
				{r: 'a', foreground: 1, foregroundMode: ansisgr.Mode256, background: 2, backgroundMode: ansisgr.Mode256},
				{r: 'r', foreground: 1, foregroundMode: ansisgr.Mode256, background: 2, backgroundMode: ansisgr.Mode256},
			},
		},
		{
			name: "RGB colors",
			in:   "foo\x1b[38;2;153;255;153;48;2;153;0;51mbar",
			want: []want{
				{r: 'f'},
				{r: 'o'},
				{r: 'o'},
				{r: 'b', foreground: 0x99ff99, foregroundMode: ansisgr.ModeRGB, background: 0x990033, backgroundMode: ansisgr.ModeRGB},
				{r: 'a', foreground: 0x99ff99, foregroundMode: ansisgr.ModeRGB, background: 0x990033, backgroundMode: ansisgr.ModeRGB},
				{r: 'r', foreground: 0x99ff99, foregroundMode: ansisgr.ModeRGB, background: 0x990033, backgroundMode: ansisgr.ModeRGB},
			},
		},
		{
			name: "bold",
			in:   "\x1b[1ma",
			want: []want{{r: 'a', bold: true}},
		},
		{
			name: "dim",
			in:   "\x1b[2ma",
			want: []want{{r: 'a', dim: true}},
		},
		{
			name: "italic",
			in:   "\x1b[3ma",
			want: []want{{r: 'a', italic: true}},
		},
		{
			name: "underline",
			in:   "\x1b[4ma",
			want: []want{{r: 'a', underline: true}},
		},
		{
			name: "blink",
			in:   "\x1b[5ma",
			want: []want{{r: 'a', blink: true}},
		},
		{
			name: "reverse",
			in:   "\x1b[7ma",
			want: []want{{r: 'a', reverse: true}},
		},
		{
			name: "invisible",
			in:   "\x1b[8ma",
			want: []want{{r: 'a', invisible: true}},
		},
		{
			name: "strikethrough",
			in:   "\x1b[9ma",
			want: []want{{r: 'a', strikethrough: true}},
		},
		{
			name: "bold reset",
			in:   "\x1b[1;22ma",
			want: []want{{r: 'a'}},
		},
		{
			name: "dim reset",
			in:   "\x1b[2;22ma",
			want: []want{{r: 'a'}},
		},
		{
			name: "italic reset",
			in:   "\x1b[3;23ma",
			want: []want{{r: 'a'}},
		},
		{
			name: "underline reset",
			in:   "\x1b[4;24ma",
			want: []want{{r: 'a'}},
		},
		{
			name: "blink reset",
			in:   "\x1b[5;25ma",
			want: []want{{r: 'a'}},
		},
		{
			name: "reverse reset",
			in:   "\x1b[7;27ma",
			want: []want{{r: 'a'}},
		},
		{
			name: "invisible reset",
			in:   "\x1b[8;28ma",
			want: []want{{r: 'a'}},
		},
		{
			name: "strikethrough reset",
			in:   "\x1b[9;29ma",
			want: []want{{r: 'a'}},
		},
		{
			name: "all reset",
			in:   "\x1b[1;30;40;0ma",
			want: []want{{r: 'a', foregroundMode: ansisgr.Mode16, backgroundMode: ansisgr.Mode16}},
		},
		{
			name: "foreground default",
			in:   "\x1b[31ma\x1b[39mb",
			want: []want{{r: 'a', foreground: 31, foregroundMode: ansisgr.Mode16}, {r: 'b', foregroundMode: ansisgr.Mode16}},
		},
		{
			name: "background default",
			in:   "\x1b[41ma\x1b[49mb",
			want: []want{{r: 'a', background: 41, backgroundMode: ansisgr.Mode16}, {r: 'b', backgroundMode: ansisgr.Mode16}},
		},
		{
			name: "multiple attributes",
			in:   "\x1b[1;3ma",
			want: []want{{r: 'a', bold: true, italic: true}},
		},
		{
			name: "color with attributes",
			in:   "\x1b[32;1;3ma",
			want: []want{{r: 'a', foreground: 32, foregroundMode: ansisgr.Mode16, bold: true, italic: true}},
		},
		{
			name: "broken sequence",
			in:   "a\x1b",
			want: []want{{r: 'a'}},
		},
		{
			name: "broken sequence which doesn't have [",
			in:   "a\x1bb",
			want: []want{{r: 'a'}},
		},
		{
			name: "invalid 16 colors",
			in:   "a\x1b[6mb",
			want: []want{{r: 'a'}, {r: 'b'}},
		},
		{
			name: "invalid 256 colors which has a non-digit value",
			in:   "a\x1b[38;5;xmb\x1b[48;5;ymc",
			want: []want{{r: 'a'}, {r: 'm'}, {r: 'b'}, {r: 'm'}, {r: 'c'}},
		},
		{
			name: "invalid 256 colors which doesn't have id",
			in:   "a\x1b[38;5mb\x1b[48;5mc",
			want: []want{{r: 'a'}, {r: 'b'}, {r: 'c'}},
		},
		{
			name: "invalid RGB colors which has a non-digit value",
			in:   "a\x1b[38;2;10;20;xmb\x1b[48;2;30;y;40mc",
			want: []want{{r: 'a'}, {r: 'm'}, {r: 'b'}, {r: ';'}, {r: '4'}, {r: '0'}, {r: 'm'}, {r: 'c'}},
		},
		{
			name: "invalid RGB colors which doesn't have red, green, and blue",
			in:   "a\x1b[38;2mb\x1b[48;2mc",
			want: []want{{r: 'a'}, {r: 'b'}, {r: 'c'}},
		},
		{
			name: "invalid next sequence",
			in:   "a\x1b[38;3;1;38;5;80mb",
			want: []want{{r: 'a'}, {r: 'b', foreground: 80, foregroundMode: ansisgr.Mode256, bold: true}},
		},
		{
			name: "invalid 256 color value",
			in:   "a\x1b[38;5;300mb",
			want: []want{{r: 'a'}, {r: 'b'}},
		},
		{
			name: "invalid RGB color value",
			in:   "a\x1b[38;2;300;1;20;48;2;300;10;20mb",
			want: []want{{r: 'a'}, {r: 'b', bold: true}},
		},
		{
			name: "has style and attrs repeatedly",
			in:   "\x1b[1;2;38;5;200;3;48;2;70;80;90ma",
			want: []want{
				{
					r:              'a',
					foreground:     200,
					foregroundMode: ansisgr.Mode256,
					background:     0x46505a,
					backgroundMode: ansisgr.ModeRGB,
					bold:           true,
					dim:            true,
					italic:         true,
				},
			},
		},
		{
			name: "has style and attrs repeatedly 2",
			in:   "\x1b[1;31;40ma\x1b[0;4;38;5;45;48;2;0;51;102mb",
			want: []want{
				{
					r:              'a',
					foreground:     31,
					foregroundMode: ansisgr.Mode16,
					background:     40,
					backgroundMode: ansisgr.Mode16,
					bold:           true,
				},
				{
					r:              'b',
					foreground:     45,
					foregroundMode: ansisgr.Mode256,
					background:     0x003366,
					backgroundMode: ansisgr.ModeRGB,
					underline:      true,
				},
			},
		},
		{
			name: "has style and attrs repeatedly 3",
			in:   "\x1b[38;2;10;38;500;32;20ma",
			want: []want{
				{
					r:              'a',
					foreground:     32,
					foregroundMode: ansisgr.Mode16,
				},
			},
		},
		{
			name: "no any parameters",
			in:   "a\x1b[mb",
			want: []want{{r: 'a'}, {r: 'b', foregroundMode: ansisgr.Mode16, backgroundMode: ansisgr.Mode16}},
		},
		{
			name: "foo",
			in:   "\x1b[38ma",
			want: []want{{r: 'a'}},
		},
	}

	for _, c := range cases {
		c := c

		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			iter := ansisgr.NewIterator(c.in)

			for i := 0; ; i++ {
				r, style, ok := iter.Next()
				if !ok {
					break
				}

				if i == len(c.want) {
					t.Fatalf("%s: i exceeded len(c.want)", string(r))
				}

				if want, got := c.want[i].r, r; want != got {
					t.Errorf("want rune %c, but got %c", want, got)
				}

				color, ok := style.Foreground()
				if want, got := c.want[i].foreground, color.Value(); want != got {
					t.Errorf("want foreground %d, but got %d", want, got)
				}
				if want, got := c.want[i].foregroundMode, color.Mode(); want != got {
					t.Errorf("want foreground mode %d, but got %d", want, got)
				}
				if c.want[i].foregroundMode == ansisgr.ModeRGB {
					r, g, b := color.RGB()
					if want := c.want[i].foreground & 0xff0000 >> 16; want != r {
						t.Errorf("want red %x, but got %x", want, r)
					}
					if want := c.want[i].foreground & 0x00ff00 >> 8; want != g {
						t.Errorf("want green %x, but got %x", want, r)
					}
					if want := c.want[i].foreground & 0x0000ff; want != b {
						t.Errorf("want blue %x, but got %x", want, r)
					}
				}

				color, ok = style.Background()
				if want, got := c.want[i].background, color.Value(); want != got {
					t.Errorf("want background %d, but got %d", want, got)
				}
				if want, got := c.want[i].backgroundMode, color.Mode(); want != got {
					t.Errorf("want background mode %d, but got %d", want, got)
				}

				if c.want[i].backgroundMode == ansisgr.ModeRGB {
					r, g, b := color.RGB()
					if want := c.want[i].background & 0xff0000 >> 16; want != r {
						t.Errorf("want red %x, but got %x", want, r)
					}
					if want := c.want[i].background & 0x00ff00 >> 8; want != g {
						t.Errorf("want green %x, but got %x", want, r)
					}
					if want := c.want[i].background & 0x0000ff; want != b {
						t.Errorf("want blue %x, but got %x", want, r)
					}
				}

				if want, got := c.want[i].bold, style.Bold(); want != got {
					t.Errorf("want Bold %t, but got %t", want, got)
				}
				if want, got := c.want[i].dim, style.Dim(); want != got {
					t.Errorf("want Dim %t, but got %t", want, got)
				}
				if want, got := c.want[i].italic, style.Italic(); want != got {
					t.Errorf("want Italic %t, but got %t", want, got)
				}
				if want, got := c.want[i].underline, style.Underline(); want != got {
					t.Errorf("want Underline %t, but got %t", want, got)
				}
				if want, got := c.want[i].blink, style.Blink(); want != got {
					t.Errorf("want Blink %t, but got %t", want, got)
				}
				if want, got := c.want[i].reverse, style.Reverse(); want != got {
					t.Errorf("want Reverse %t, but got %t", want, got)
				}
				if want, got := c.want[i].invisible, style.Invisible(); want != got {
					t.Errorf("want Invisible %t, but got %t", want, got)
				}
				if want, got := c.want[i].strikethrough, style.Strikethrough(); want != got {
					t.Errorf("want Strikethrough %t, but got %t", want, got)
				}
			}
		})
	}
}

func ExampleNewIterator() {
	in := "a\x1b[1;31;40mb\x1b[0;4;38;5;45;48;2;0;51;102mc"
	iter := ansisgr.NewIterator(in)

	var got []string

	for {
		r, style, ok := iter.Next()
		if !ok {
			break
		}

		var s []string
		s = append(s, fmt.Sprintf("rune: %c", r))
		if color, ok := style.Foreground(); ok {
			switch color.Mode() {
			case ansisgr.Mode16, ansisgr.Mode256:
				s = append(s, fmt.Sprintf("foreground: %d", color.Value()))
			case ansisgr.ModeRGB:
				red, green, blue := color.RGB()
				s = append(s, fmt.Sprintf("foreground: %d (%d, %d, %d)", color.Value(), red, green, blue))
			}
		}
		if color, ok := style.Background(); ok {
			switch color.Mode() {
			case ansisgr.Mode16, ansisgr.Mode256:
				s = append(s, fmt.Sprintf("background: %d", color.Value()))
			case ansisgr.ModeRGB:
				red, green, blue := color.RGB()
				s = append(s, fmt.Sprintf("background: 0x%06x (%d, %d, %d)", color.Value(), red, green, blue))
			}
		}
		if style.Bold() {
			s = append(s, "bold")
		}

		got = append(got, strings.Join(s, ", "))
	}

	fmt.Println(strings.Join(got, "\n"))

	// Output:
	// rune: a
	// rune: b, foreground: 31, background: 40, bold
	// rune: c, foreground: 45, background: 0x003366 (0, 51, 102)
}

func FuzzIterator(f *testing.F) {
	f.Add("Lorem ipsum dolor sit amet, consectetur adipiscing elit")
	f.Add("Sed eget dui libero.\nVivamus tempus, magna nec mollis convallis, ipsum justo tincidunt ligula, ut varius est mi id nisl.\nMorbi commodo turpis risus, nec vehicula leo auctor sit amet.\nUt imperdiet suscipit massa ac vehicula.\nInterdum et malesuada fames ac ante ipsum primis in faucibus.\nPraesent ligula orci, facilisis pulvinar varius eget, iaculis in erat.\nProin pellentesque arcu sed nisl consectetur tristique.\nQuisque tempus blandit dignissim.\nPhasellus dignissim sollicitudin mauris, sed gravida arcu luctus tincidunt.\nNunc rhoncus sed eros vel molestie.\nAenean sodales tortor eu libero rutrum, et lobortis orci scelerisque.\nPraesent sollicitudin, nunc ut consequat commodo, risus velit consectetur nibh, quis pretium nunc elit et erat.")
	f.Add("foo\x1b[31;1;44;0;90;105;38;5;12;48;5;226;38;2;10;20;30;48;2;200;100;50mbar")

	f.Fuzz(func(t *testing.T, s string) {
		iter := ansisgr.NewIterator(s)
		for {
			_, _, ok := iter.Next()
			if !ok {
				break
			}
		}
	})
}
