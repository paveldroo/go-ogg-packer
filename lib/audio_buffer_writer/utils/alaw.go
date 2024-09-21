package utils

func S16ToALaw(samples []int16) []uint8 {
	res := make([]uint8, len(samples))
	for i := range samples {
		res[i] = s16ToALaw(samples[i])
	}
	return res
}

func s16ToALaw(sample int16) uint8 {
	var mask uint8
	if sample >= 0 {
		mask = 0xD5
	} else {
		mask = 0x55
		sample = -sample
		if sample > 0x7fff {
			sample = 0x7fff
		}
	}
	var res uint8
	if sample < 256 {
		res = uint8(sample >> 4)
	} else {
		seg := seg(sample)
		res = uint8((seg << 4) | ((sample >> (seg + 3)) & 0x0f))
	}
	return res ^ mask
}

func seg(sample int16) int16 {
	var res int16 = 1
	sample >>= 8
	if sample&0xf0 != 0 {
		sample >>= 4
		res += 4
	}
	if sample&0x0c != 0 {
		sample >>= 2
		res += 2
	}
	if sample&0x02 != 0 {
		res += 1
	}
	return res
}
