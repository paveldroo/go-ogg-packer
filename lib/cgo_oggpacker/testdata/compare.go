package testdata

import (
	"bytes"
	"regexp"
)

func CompareOggAudio(resAudio, refAudio []byte) bool {
	pattern := regexp.MustCompile(`OggS`)
	resPagesBoundaries := pattern.FindAllIndex(resAudio, -1)
	refPagesBoundaries := pattern.FindAllIndex(refAudio, -1)
	if len(resPagesBoundaries) != len(refPagesBoundaries) { // if we have different Ogg pages quantity we can't be equal
		return false
	}

	for i := range resPagesBoundaries {
		if i == 0 { // safe and easy for understanding defence from `out of boundaries`
			continue
		}
		resAudioPage := resAudio[resPagesBoundaries[i-1][0]:resPagesBoundaries[i][0]]
		refAudioPage := refAudio[refPagesBoundaries[i-1][0]:refPagesBoundaries[i][0]]
		if !compareOggPage(resAudioPage, refAudioPage) {
			return false
		}
	}
	return true
}

// compare page segments except serial number and checksum, because serial number generated randomly generated
func compareOggPage(page, referencePage []byte) bool {
	return bytes.Equal(page[:14], referencePage[:14]) && //compare bytes before serial number
		bytes.Equal(page[18:22], referencePage[18:22]) && //page sequence number
		bytes.Equal(page[26:], referencePage[26:]) //compare bytes after checksum
}
