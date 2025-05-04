package packer

type Packer struct {
}

func New() *Packer {
	return &Packer{}
}

func AddChunk(pcm []int16) error {
	return nil
}

func Result() ([]byte, error) {
	return nil, nil
}
