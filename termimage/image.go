package termimage

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"strings"
	"time"

	"github.com/mattn/go-isatty"
	"github.com/mattn/go-sixel"
	"github.com/nfnt/resize"
	"golang.org/x/crypto/ssh/terminal"
)

func checkIterm() bool {
	return strings.HasPrefix(os.Getenv("TERM_PROGRAM"), "iTerm")
}

func checkSixel() bool {
	if isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		return true
	}
	s, err := terminal.MakeRaw(1)
	if err != nil {
		return false
	}
	defer terminal.Restore(1, s)
	_, err = os.Stdout.Write([]byte("\x1b[c"))
	if err != nil {
		return false
	}
	defer os.Stdout.SetReadDeadline(time.Time{})

	var b [100]byte
	n, err := os.Stdout.Read(b[:])
	if err != nil {
		return false
	}
	var supportedTerminals = []string{
		"\x1b[?62;", // VT240
		"\x1b[?63;", // wsltty
		"\x1b[?64;", // mintty
		"\x1b[?65;", // RLogin
	}
	supported := false
	for _, supportedTerminal := range supportedTerminals {
		if bytes.HasPrefix(b[:n], []byte(supportedTerminal)) {
			supported = true
			break
		}
	}
	if !supported {
		return false
	}

	sb := b[6:n]
	n = bytes.IndexByte(sb, 'c')
	if n != -1 {
		sb = sb[:n]
	}
	for _, t := range bytes.Split(sb, []byte(";")) {
		if len(t) == 1 && t[0] == '4' {
			return true
		}
	}
	return false
}

func iTermEncode(out io.Writer, w, h uint, img image.Image) error {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return err
	}

	fmt.Fprint(out, "\033]1337;")
	fmt.Fprintf(out, "File=inline=1")
	fmt.Fprintf(out, ";width=%dpx", w)
	fmt.Fprintf(out, ";height=%dpx", h)
	fmt.Fprint(out, ":")
	fmt.Fprintf(out, "%s", base64.StdEncoding.EncodeToString(buf.Bytes()))
	fmt.Fprint(out, "\a\n")
	return nil
}

func Render(out io.Writer, r io.Reader, w uint) error {
	img, _, err := image.Decode(r)
	if err != nil {
		return err
	}
	h := uint(float64(img.Bounds().Dx()) / float64(img.Bounds().Dy()) * float64(w))
	img = resize.Resize(w, h, img, resize.Lanczos3)

	if checkIterm() {
		return iTermEncode(out, w, h, img)
	} else if checkSixel() {
		return sixel.NewEncoder(out).Encode(img)
	}
	return nil
}
