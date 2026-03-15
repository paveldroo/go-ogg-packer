package ogg_test

import (
	"encoding/binary"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/paveldroo/go-ogg-packer/ogg"
	"github.com/paveldroo/go-ogg-packer/opus"
	"github.com/paveldroo/go-ogg-packer/tests/testutil"
)

// TestPacker compares OGG packer output against reference OGG files.
// References must be generated before running: task generate-ref
// Or run everything together: task test
func TestPacker(t *testing.T) {
	tests := []struct {
		name       string
		fileBase   string
		channels   int
		sampleRate int
	}{
		{
			name:       "8k 1ch",
			fileBase:   "8k_1ch",
			channels:   1,
			sampleRate: 8000,
		},
		{
			name:       "12k 1ch",
			fileBase:   "12k_1ch",
			channels:   1,
			sampleRate: 12000,
		},
		{
			name:       "16k 1ch",
			fileBase:   "16k_1ch",
			channels:   1,
			sampleRate: 16000,
		},
		{
			name:       "24k 1ch",
			fileBase:   "24k_1ch",
			channels:   1,
			sampleRate: 24000,
		},
		{
			name:       "48k 1ch",
			fileBase:   "48k_1ch",
			channels:   1,
			sampleRate: 48000,
		},
		{
			name:       "8k 2ch",
			fileBase:   "8k_2ch",
			channels:   2,
			sampleRate: 8000,
		},
		{
			name:       "12k 2ch",
			fileBase:   "12k_2ch",
			channels:   2,
			sampleRate: 12000,
		},
		{
			name:       "16k 2ch",
			fileBase:   "16k_2ch",
			channels:   2,
			sampleRate: 16000,
		},
		{
			name:       "24k 2ch",
			fileBase:   "24k_2ch",
			channels:   2,
			sampleRate: 24000,
		},
		{
			name:       "48k 2ch",
			fileBase:   "48k_2ch",
			channels:   2,
			sampleRate: 48000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packer, err := ogg.New(uint8(tt.channels), uint32(tt.sampleRate), testutil.TestSerialNo)
			if err != nil {
				t.Fatalf("create ogg packer: %s", err.Error())
			}

			cfg := opus.Config{
				SampleRate:  tt.sampleRate,
				NumChannels: tt.channels,
				FrameSize:   time.Duration(opus.FrameSize) * time.Millisecond,
			}
			frameSizeSamples := opus.FrameSizeSamples(cfg)

			opusFilename := fmt.Sprintf("../tests/testdata/opus_raw/%s.opus_raw", tt.fileBase)
			rawOpusData := testutil.RawOpusPackets(t, opusFilename)
			for _, packet := range rawOpusData {
				if err := packer.AddChunk(packet, false, frameSizeSamples); err != nil {
					t.Fatalf("send opus chunk to packer: %s", err.Error())
				}
			}

			oggData, err := packer.FlushPages()
			if err != nil {
				t.Fatalf("read all pages from packer: %s", err.Error())
			}

			refFilename := fmt.Sprintf("../tests/testdata/want/ogg/%s.ogg", tt.fileBase)
			refData, err := os.ReadFile(refFilename)
			if err != nil {
				t.Fatalf("open reference file: %s", err.Error())
			}

			if !reflect.DeepEqual(refData, oggData) {
				t.Fatal("base data and test data are not equal")
			}
		})
	}
}

func TestPacker_EOS(t *testing.T) {
	// 48kHz * 60ms * 1ch / 1000 = 2880 total samples
	frameSizeSamples := 2880

	t.Run("FlushPages patches EOS on last page", func(t *testing.T) {
		packer, err := ogg.New(1, 48000, testutil.TestSerialNo)
		if err != nil {
			t.Fatalf("create ogg packer: %s", err)
		}

		opusFilename := "../tests/testdata/opus_raw/48k_1ch.opus_raw"
		packets := testutil.RawOpusPackets(t, opusFilename)
		for _, packet := range packets {
			if err := packer.AddChunk(packet, false, frameSizeSamples); err != nil {
				t.Fatalf("add chunk: %s", err)
			}
		}

		oggData, err := packer.FlushPages()
		if err != nil {
			t.Fatalf("read pages: %s", err)
		}

		assertLastPageHasEOS(t, oggData)
	})

	t.Run("explicit EOS via AddChunk", func(t *testing.T) {
		packer, err := ogg.New(1, 48000, testutil.TestSerialNo)
		if err != nil {
			t.Fatalf("create ogg packer: %s", err)
		}

		opusFilename := "../tests/testdata/opus_raw/48k_1ch.opus_raw"
		packets := testutil.RawOpusPackets(t, opusFilename)
		for i, packet := range packets {
			eos := i == len(packets)-1
			if err := packer.AddChunk(packet, eos, frameSizeSamples); err != nil {
				t.Fatalf("add chunk: %s", err)
			}
		}

		oggData, err := packer.ReadPages()
		if err != nil {
			t.Fatalf("read pages: %s", err)
		}

		assertLastPageHasEOS(t, oggData)
	})
}

func TestPacker_AutoSamplesCount(t *testing.T) {
	tests := []struct {
		name       string
		fileBase   string
		channels   int
		sampleRate int
	}{
		{name: "48k 1ch", fileBase: "48k_1ch", channels: 1, sampleRate: 48000},
		{name: "48k 2ch", fileBase: "48k_2ch", channels: 2, sampleRate: 48000},
		{name: "16k 1ch", fileBase: "16k_1ch", channels: 1, sampleRate: 16000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packer, err := ogg.New(uint8(tt.channels), uint32(tt.sampleRate), testutil.TestSerialNo)
			if err != nil {
				t.Fatalf("create ogg packer: %s", err)
			}

			opusFilename := fmt.Sprintf("../tests/testdata/opus_raw/%s.opus_raw", tt.fileBase)
			rawOpusData := testutil.RawOpusPackets(t, opusFilename)
			for _, packet := range rawOpusData {
				if err := packer.AddChunk(packet, false, -1); err != nil {
					t.Fatalf("add chunk: %s", err)
				}
			}

			oggData, err := packer.FlushPages()
			if err != nil {
				t.Fatalf("read pages: %s", err)
			}

			refFilename := fmt.Sprintf("../tests/testdata/want/ogg/%s.ogg", tt.fileBase)
			refData, err := os.ReadFile(refFilename)
			if err != nil {
				t.Fatalf("open reference file: %s", err)
			}

			if !reflect.DeepEqual(refData, oggData) {
				t.Fatal("auto samples count: output does not match reference")
			}
		})
	}
}

func TestPacker_Streaming(t *testing.T) {
	tests := []struct {
		name       string
		fileBase   string
		channels   int
		sampleRate int
	}{
		{name: "48k 1ch", fileBase: "48k_1ch", channels: 1, sampleRate: 48000},
		{name: "16k 1ch", fileBase: "16k_1ch", channels: 1, sampleRate: 16000},
		{name: "48k 2ch", fileBase: "48k_2ch", channels: 2, sampleRate: 48000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := opus.Config{
				SampleRate:  tt.sampleRate,
				NumChannels: tt.channels,
				FrameSize:   time.Duration(opus.FrameSize) * time.Millisecond,
			}
			frameSizeSamples := opus.FrameSizeSamples(cfg)

			packer, err := ogg.New(uint8(tt.channels), uint32(tt.sampleRate), testutil.TestSerialNo)
			if err != nil {
				t.Fatalf("create ogg packer: %s", err)
			}

			opusFilename := fmt.Sprintf("../tests/testdata/opus_raw/%s.opus_raw", tt.fileBase)
			packets := testutil.RawOpusPackets(t, opusFilename)

			// Read BOS + Tags pages produced by New().
			headerData, err := packer.ReadPages()
			if err != nil {
				t.Fatalf("read header pages: %s", err)
			}

			var result []byte
			result = append(result, headerData...)

			// Verify header pages have no EOS.
			headerPages := parseOggPages(t, headerData)
			for i, p := range headerPages {
				if p.headerType&0x04 != 0 {
					t.Errorf("header page %d has EOS flag set", i)
				}
			}

			// Stream each packet: AddChunk + ReadPages.
			for i, packet := range packets {
				if err := packer.AddChunk(packet, false, frameSizeSamples); err != nil {
					t.Fatalf("add chunk %d: %s", i, err)
				}

				if i < len(packets)-1 {
					chunk, err := packer.ReadPages()
					if err != nil {
						t.Fatalf("read pages at chunk %d: %s", i, err)
					}
					// Mid-stream pages must not have EOS.
					midPages := parseOggPages(t, chunk)
					for j, p := range midPages {
						if p.headerType&0x04 != 0 {
							t.Errorf("mid-stream chunk %d page %d has EOS flag set", i, j)
						}
					}
					result = append(result, chunk...)
				} else {
					// Last packet: FlushPages to set EOS.
					chunk, err := packer.FlushPages()
					if err != nil {
						t.Fatalf("flush pages: %s", err)
					}
					result = append(result, chunk...)
				}
			}

			// Compare concatenated streaming output to batch reference.
			refFilename := fmt.Sprintf("../tests/testdata/want/ogg/%s.ogg", tt.fileBase)
			refData, err := os.ReadFile(refFilename)
			if err != nil {
				t.Fatalf("open reference file: %s", err)
			}

			if !reflect.DeepEqual(refData, result) {
				t.Fatalf("streaming output does not match batch reference (got %d bytes, want %d)", len(result), len(refData))
			}
		})
	}
}

func TestPacker_RFC7845_GranulePositions(t *testing.T) {
	tests := []struct {
		name       string
		fileBase   string
		channels   int
		sampleRate int
	}{
		{name: "48k 1ch", fileBase: "48k_1ch", channels: 1, sampleRate: 48000},
		{name: "16k 1ch", fileBase: "16k_1ch", channels: 1, sampleRate: 16000},
	}

	const samplesPerPacketAt48k int64 = 2880 // 60ms * 48kHz

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := opus.Config{
				SampleRate:  tt.sampleRate,
				NumChannels: tt.channels,
				FrameSize:   time.Duration(opus.FrameSize) * time.Millisecond,
			}
			frameSizeSamples := opus.FrameSizeSamples(cfg)

			packer, err := ogg.New(uint8(tt.channels), uint32(tt.sampleRate), testutil.TestSerialNo)
			if err != nil {
				t.Fatalf("create ogg packer: %s", err)
			}

			opusFilename := fmt.Sprintf("../tests/testdata/opus_raw/%s.opus_raw", tt.fileBase)
			packets := testutil.RawOpusPackets(t, opusFilename)
			for _, packet := range packets {
				if err := packer.AddChunk(packet, false, frameSizeSamples); err != nil {
					t.Fatalf("add chunk: %s", err)
				}
			}

			oggData, err := packer.FlushPages()
			if err != nil {
				t.Fatalf("flush pages: %s", err)
			}

			pages := parseOggPages(t, oggData)
			if len(pages) < 3 {
				t.Fatalf("expected at least 3 pages, got %d", len(pages))
			}

			// Page 0: BOS — granule must be 0.
			if pages[0].headerType&0x02 == 0 {
				t.Error("page 0: BOS flag not set")
			}
			if pages[0].granule != 0 {
				t.Errorf("page 0 (BOS): granule = %d, want 0", pages[0].granule)
			}
			if pages[0].pageSeq != 0 {
				t.Errorf("page 0: pageSeq = %d, want 0", pages[0].pageSeq)
			}

			// OpusHead pre-skip field (bytes 10-11) must match PreSkip.
			if len(pages[0].body) < 12 {
				t.Fatalf("page 0 body too short for OpusHead: %d bytes", len(pages[0].body))
			}
			headerPreSkip := binary.LittleEndian.Uint16(pages[0].body[10:12])
			if headerPreSkip != ogg.PreSkip {
				t.Errorf("OpusHead pre-skip = %d, want %d", headerPreSkip, ogg.PreSkip)
			}

			// Page 1: Tags — granule must be 0.
			if pages[1].granule != 0 {
				t.Errorf("page 1 (Tags): granule = %d, want 0", pages[1].granule)
			}
			if pages[1].pageSeq != 1 {
				t.Errorf("page 1: pageSeq = %d, want 1", pages[1].pageSeq)
			}
			if pages[1].headerType&(0x02|0x04) != 0 {
				t.Errorf("page 1 (Tags): unexpected BOS/EOS flags: 0x%02x", pages[1].headerType)
			}

			// Page 2: first audio — granule must be PreSkip + samplesPerPacketAt48k.
			wantGranule := int64(ogg.PreSkip) + samplesPerPacketAt48k
			if pages[2].granule != wantGranule {
				t.Errorf("page 2 (first audio): granule = %d, want %d", pages[2].granule, wantGranule)
			}
			if pages[2].pageSeq != 2 {
				t.Errorf("page 2: pageSeq = %d, want 2", pages[2].pageSeq)
			}

			// Pages 3..N: granule increments by samplesPerPacketAt48k, pageSeq increments by 1.
			for i := 3; i < len(pages); i++ {
				wantGranule += samplesPerPacketAt48k
				if pages[i].granule != wantGranule {
					t.Errorf("page %d: granule = %d, want %d", i, pages[i].granule, wantGranule)
				}
				if pages[i].pageSeq != uint32(i) {
					t.Errorf("page %d: pageSeq = %d, want %d", i, pages[i].pageSeq, i)
				}
				// Monotonicity.
				if pages[i].granule <= pages[i-1].granule {
					t.Errorf("page %d: granule %d not strictly greater than page %d granule %d",
						i, pages[i].granule, i-1, pages[i-1].granule)
				}
			}
		})
	}
}

func TestPacker_RFC7845_RoundTrip(t *testing.T) {
	const sampleRate = 48000
	const channels = 1

	cfg := opus.Config{
		SampleRate:  sampleRate,
		NumChannels: channels,
		FrameSize:   time.Duration(opus.FrameSize) * time.Millisecond,
	}
	frameSizeSamples := opus.FrameSizeSamples(cfg)

	packer, err := ogg.New(uint8(channels), uint32(sampleRate), testutil.TestSerialNo)
	if err != nil {
		t.Fatalf("create ogg packer: %s", err)
	}

	opusFilename := "../tests/testdata/opus_raw/48k_1ch.opus_raw"
	packets := testutil.RawOpusPackets(t, opusFilename)
	for _, packet := range packets {
		if err := packer.AddChunk(packet, false, frameSizeSamples); err != nil {
			t.Fatalf("add chunk: %s", err)
		}
	}

	oggData, err := packer.FlushPages()
	if err != nil {
		t.Fatalf("flush pages: %s", err)
	}

	// Decode OGG back to PCM.
	pcm := testutil.PCMFromOgg(t, oggData, sampleRate, channels)

	// Each opus packet decodes to frameSizeSamples/channels samples per channel.
	// The decoder outputs all packets including the pre-skip region.
	// Expected decoded samples = numPackets * samplesPerFrame.
	expectedSamples := len(packets) * frameSizeSamples
	if len(pcm) != expectedSamples {
		t.Errorf("decoded PCM length = %d, want %d (packets=%d, samplesPerFrame=%d)",
			len(pcm), expectedSamples, len(packets), frameSizeSamples)
	}
}

type parsedPage struct {
	headerType byte
	granule    int64
	serial     uint32
	pageSeq    uint32
	body       []byte
}

func parseOggPages(t *testing.T, data []byte) []parsedPage {
	t.Helper()

	var pages []parsedPage
	off := 0
	for off < len(data) {
		if off+27 > len(data) {
			t.Fatalf("truncated page header at offset %d", off)
		}
		if string(data[off:off+4]) != "OggS" {
			t.Fatalf("expected OggS at offset %d, got %q", off, data[off:off+4])
		}

		headerType := data[off+5]
		granule := int64(binary.LittleEndian.Uint64(data[off+6 : off+14]))
		serial := binary.LittleEndian.Uint32(data[off+14 : off+18])
		pageSeq := binary.LittleEndian.Uint32(data[off+18 : off+22])
		nsegs := int(data[off+26])

		if off+27+nsegs > len(data) {
			t.Fatalf("truncated segment table at offset %d", off)
		}

		bodySize := 0
		for i := 0; i < nsegs; i++ {
			bodySize += int(data[off+27+i])
		}

		bodyStart := off + 27 + nsegs
		if bodyStart+bodySize > len(data) {
			t.Fatalf("truncated body at offset %d", off)
		}

		body := make([]byte, bodySize)
		copy(body, data[bodyStart:bodyStart+bodySize])

		pages = append(pages, parsedPage{
			headerType: headerType,
			granule:    granule,
			serial:     serial,
			pageSeq:    pageSeq,
			body:       body,
		})

		off = bodyStart + bodySize
	}

	return pages
}

func assertLastPageHasEOS(t *testing.T, data []byte) {
	t.Helper()

	lastPageStart := -1
	for i := len(data) - 4; i >= 0; i-- {
		if string(data[i:i+4]) == "OggS" {
			lastPageStart = i
			break
		}
	}
	if lastPageStart < 0 {
		t.Fatal("no OggS page found in output")
	}
	if data[lastPageStart+5]&0x04 == 0 { // 0x04 = EOS flag per RFC 3533
		t.Fatal("last page does not have EOS flag set")
	}
}
