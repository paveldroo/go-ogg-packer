package testutil

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"os"
	"time"

	extopus "gopkg.in/hraban/opus.v2"
	extogg "mccoy.space/g/ogg"

	packer "github.com/paveldroo/go-ogg-packer"
	"github.com/paveldroo/go-ogg-packer/ogg"
	"github.com/paveldroo/go-ogg-packer/opus"
)

// GeneratePackerRef runs the full packer pipeline and writes a PCM reference file.
func GeneratePackerRef(sourceFname, refFname string, sampleRate int) error {
	pcmBytes, err := os.ReadFile(sourceFname)
	if err != nil {
		return fmt.Errorf("read pcm source: %w", err)
	}

	sourcePCM := make([]int16, len(pcmBytes)/2)
	for i := range sourcePCM {
		sourcePCM[i] = int16(pcmBytes[2*i]) | int16(pcmBytes[2*i+1])<<8
	}

	cfg := opus.Config{
		SampleRate:  sampleRate,
		NumChannels: opus.NumChannels,
		FrameSize:   time.Duration(opus.FrameSize) * time.Millisecond,
	}

	p, err := packer.New(cfg)
	if err != nil {
		return fmt.Errorf("create packer: %w", err)
	}

	for i := 0; i < len(sourcePCM); i++ {
		end := min(i+2048, len(sourcePCM))
		if err := p.SendPCMChunk(sourcePCM[i:end]); err != nil {
			return fmt.Errorf("send PCM chunk: %w", err)
		}
		i = end
	}

	audioData, err := p.GetResult()
	if err != nil {
		return fmt.Errorf("get result: %w", err)
	}

	pcm := pcmFromOgg(audioData, sampleRate, opus.NumChannels)

	file, err := os.Create(refFname)
	if err != nil {
		return fmt.Errorf("create reference file: %w", err)
	}
	defer file.Close()

	for _, sample := range pcm {
		if err := binary.Write(file, binary.LittleEndian, sample); err != nil {
			return fmt.Errorf("write pcm sample: %w", err)
		}
	}

	fmt.Printf("generated packer ref: %s\n", refFname)
	return nil
}

// GenerateOggRef packs raw opus packets into OGG and writes a reference file.
func GenerateOggRef(opusRawFname, refFname string, channels, sampleRate int) error {
	rawOpus, err := readRawOpusPackets(opusRawFname)
	if err != nil {
		return fmt.Errorf("read opus raw: %w", err)
	}

	cfg := opus.Config{
		SampleRate:  sampleRate,
		NumChannels: channels,
		FrameSize:   time.Duration(opus.FrameSize) * time.Millisecond,
	}
	frameSizeSamples := opus.FrameSizeSamples(cfg)

	p, err := ogg.New(uint8(channels), uint32(sampleRate))
	if err != nil {
		return fmt.Errorf("create ogg packer: %w", err)
	}

	for _, packet := range rawOpus {
		if err := p.AddChunk(packet, false, frameSizeSamples); err != nil {
			return fmt.Errorf("add chunk: %w", err)
		}
	}

	oggData, err := p.ReadPages()
	if err != nil {
		return fmt.Errorf("read pages: %w", err)
	}

	if err := os.WriteFile(refFname, oggData, 0666); err != nil {
		return fmt.Errorf("write ogg file: %w", err)
	}

	fmt.Printf("generated ogg ref: %s\n", refFname)
	return nil
}

// GenerateOpusRaw encodes PCM to opus packets and writes a GOB-encoded file.
func GenerateOpusRaw(pcmSource, outPath string, sampleRate int) error {
	pcmBytes, err := os.ReadFile(pcmSource)
	if err != nil {
		return fmt.Errorf("read pcm: %w", err)
	}

	samples := make([]int16, len(pcmBytes)/2)
	for i := range samples {
		samples[i] = int16(pcmBytes[2*i]) | int16(pcmBytes[2*i+1])<<8
	}

	cfg := opus.Config{
		SampleRate:  sampleRate,
		NumChannels: opus.NumChannels,
		FrameSize:   time.Duration(opus.FrameSize) * time.Millisecond,
	}

	encoder, err := opus.NewEncoder(cfg)
	if err != nil {
		return fmt.Errorf("create encoder: %w", err)
	}

	packets, err := encoder.EncodeWithPadding(samples)
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	if err := gob.NewEncoder(f).Encode(packets); err != nil {
		return fmt.Errorf("gob encode: %w", err)
	}

	fmt.Printf("generated opus raw: %s (%d packets)\n", outPath, len(packets))
	return nil
}

func pcmFromOgg(oggData []byte, sampleRate, numChannels int) []int16 {
	b := bytes.NewBuffer(oggData)
	oggDecoder := extogg.NewDecoder(b)

	opusDecoder, _ := extopus.NewDecoder(sampleRate, numChannels)

	pcmBuffer := make([]int16, opus.FrameSize*sampleRate*numChannels/1000)

	var pcm []int16
	for {
		page, err := oggDecoder.Decode()
		if err != nil {
			break
		}
		for _, packet := range page.Packets {
			n, err := opusDecoder.Decode(packet, pcmBuffer)
			if err != nil {
				continue
			}
			pcm = append(pcm, pcmBuffer[:n]...)
		}
	}
	return pcm
}

func readRawOpusPackets(fname string) ([][]byte, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var data [][]byte
	if err := gob.NewDecoder(f).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}
