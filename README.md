<h1 align="center">
  <br>
  Go Ogg Packer
  <br>
</h1>
<h4 align="center">PCM to Ogg Chunked Encoder: Streamlined Audio Packaging in Go</h4>
<p align="center">
  <a href="https://pkg.go.dev/github.com/paveldroo/go-ogg-packer"><img src="https://pkg.go.dev/badge/github.com/paveldroo/go-ogg-packer.svg" alt="Go Reference"></a>
  <a href="https://goreportcard.com/report/github.com/paveldroo/go-ogg-packer"><img src="https://goreportcard.com/badge/github.com/paveldroo/go-ogg-packer" alt="Go Report Card"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License: MIT"></a>
</p>
<br>

### Sample rates and channel support
- Only **48000 Hz** sample rate and **1 channel** (mono) supported at the moment. Feel free to add a PR with different audio settings.

### RFCs
- **RFC 6716**: [The Ogg Encapsulation Format Version 0](https://www.ietf.org/rfc/rfc3533.txt)

### Running
Check out [examples](examples) for demonstration of using Go Ogg Packer with Wav files.

### Roadmap
- [x] Use AudioBufferWriter wrapper for CGo ogg_packer to get green tests before implementing native Go code
- [x] Dirty native Go ogg packer implementation + comparison auto tests
- [x] Add better error handling, split `packer` package to several files/packages
- [x] Check linters
- [x] Remove C ogg lib dependency
- [x] Remove direct C opus lib dependency
- [x] Add ogg encoder to codebase, eliminate unmaintained dependency
- [x] Use only opus raw data in tests
- [x] Add examples for Opus and Wav formats
- [x] Add direct PCM to Ogg API with Opus converting under the hood
- [ ] Add opus encoder tests
- [ ] Check lib layout with peers
- [ ] Production-ready release 1.0

### License
MIT License - see [LICENSE](LICENSE) for full text
