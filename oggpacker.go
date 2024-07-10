package oggpacker

type Packer struct {
	stream []byte
}

func NewPacker(sampleRate, numChannels int) (*Packer, error) {
	return &Packer{}, nil
}
