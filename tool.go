package packer

func ResampleLinearInt16(input []int16, inRate, outRate int) []int16 {
	ratio := float64(inRate) / float64(outRate)
	outLen := int(float64(len(input)) / ratio)
	output := make([]int16, outLen)

	for i := 0; i < outLen; i++ {
		srcIdx := float64(i) * ratio
		idx0 := int(srcIdx)
		idx1 := idx0 + 1
		if idx1 >= len(input) {
			idx1 = len(input) - 1
		}
		frac := srcIdx - float64(idx0)
		sample := float64(input[idx0])*(1-frac) + float64(input[idx1])*frac
		output[i] = int16(sample)
	}
	return output
}
