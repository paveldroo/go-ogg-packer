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

### Dependencies
You need to install the following system dependencies:
- `libopus`
- `libopusfile`

For example, on Linux you can install them using:
```bash
apt-get install libopus-dev libopusfile-dev
```
For the most up-to-date list of dependencies, feel free to check the `Dockerfile`.

### What is PCM
- [Pulse-code modulation](https://en.wikipedia.org/wiki/Pulse-code_modulation) - universal format to transfer audio data
- You should use appropriate library to convert audio data from your container (WAV, MP3, etc.) to PCM data before using Go Ogg Packer
- Your PCM sample rate and channels count should be supported by this library

### Sample rates and channels support
- Only **48000 Hz** sample rate and **1 channel** (mono) supported at the moment. Feel free to add a PR with different audio settings.

### RFCs
- **RFC 6716**: [The Ogg Encapsulation Format Version 0](https://www.ietf.org/rfc/rfc3533.txt)

### Benchmarks
```
          │ c_ogg.txt  │        native_ogg.txt        │
          │   sec/op   │   sec/op    vs base          │
Packer-14   2.027 ± 2%   2.024 ± 3%  ~ (p=0.971 n=10)

          │  c_ogg.txt   │            native_ogg.txt            │
          │     B/op     │     B/op      vs base                │
Packer-14   282.2Mi ± 0%   355.8Mi ± 0%  +26.09% (p=0.000 n=10)

          │  c_ogg.txt  │           native_ogg.txt           │
          │  allocs/op  │  allocs/op   vs base               │
Packer-14   25.30M ± 0%   25.32M ± 0%  +0.06% (p=0.000 n=10)
```
Although the native Go implementation allocates 26% more space, the difference in overall execution speed is statistically insignificant.

### Running
Check out [examples](examples) for demonstration of using Go Ogg Packer with Wav files.

### License
MIT License - see [LICENSE](LICENSE) for full text
