package main

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"

	"github.com/paveldroo/go-ogg-packer/tests/testutil"
)

const wavSource = "tests/testdata/input/yanka.wav"

var rates = []struct {
	sampleRate int
	fileBase   string
}{
	{8000, "8k_1ch"},
	{12000, "12k_1ch"},
	{16000, "16k_1ch"},
	{24000, "24k_1ch"},
	{48000, "48k_1ch"},
}

func main() {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		log.Fatal("ffmpeg is required for reference generation but not found in PATH")
	}

	fmt.Println("=== generating PCM sources from WAV ===")
	for _, r := range rates {
		dst := fmt.Sprintf("tests/testdata/input/%s.pcm", r.fileBase)
		cmd := exec.Command("ffmpeg", "-y",
			"-i", wavSource,
			"-ac", "1",
			"-ar", strconv.Itoa(r.sampleRate),
			"-f", "s16le",
			dst,
		)
		if out, err := cmd.CombinedOutput(); err != nil {
			log.Fatalf("ffmpeg %s: %s\n%s", r.fileBase, err, out)
		}
		fmt.Printf("generated PCM: %s\n", dst)
	}

	fmt.Println("=== generating opus_raw fixtures ===")
	for _, r := range rates {
		src := fmt.Sprintf("tests/testdata/input/%s.pcm", r.fileBase)
		dst := fmt.Sprintf("tests/testdata/opus_raw/%s.opus_raw", r.fileBase)
		if err := testutil.GenerateOpusRaw(src, dst, r.sampleRate); err != nil {
			log.Fatalf("generate opus raw %s: %s", r.fileBase, err)
		}
	}

	fmt.Println("=== generating ogg references ===")
	for _, r := range rates {
		src := fmt.Sprintf("tests/testdata/opus_raw/%s.opus_raw", r.fileBase)
		dst := fmt.Sprintf("tests/testdata/want/ogg/%s.ogg", r.fileBase)
		if err := testutil.GenerateOggRef(src, dst, 1, r.sampleRate); err != nil {
			log.Fatalf("generate ogg ref %s: %s", r.fileBase, err)
		}
	}

	fmt.Println("=== generating packer references ===")
	for _, r := range rates {
		src := fmt.Sprintf("tests/testdata/input/%s.pcm", r.fileBase)
		dst := fmt.Sprintf("tests/testdata/want/packer/%s.pcm", r.fileBase)
		if err := testutil.GeneratePackerRef(src, dst, r.sampleRate); err != nil {
			log.Fatalf("generate packer ref %s: %s", r.fileBase, err)
		}
	}

	fmt.Println("=== done ===")
}
