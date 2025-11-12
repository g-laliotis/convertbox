# Convertbox 🎬

**Production-ready faceless YouTube Shorts generator for AI Unboxed by UnboxGio**

Generate engaging YouTube Shorts completely offline using local LLM (Ollama), high-quality TTS (Coqui), dynamic backgrounds, and professional video rendering.

## ✨ Features

- 🤖 **Local AI Script Generation** - Uses Ollama (Mistral) for engaging, hook-driven scripts
- 🎙️ **High-Quality TTS** - Coqui TTS with eSpeak fallback for natural narration
- 🎨 **Dynamic Backgrounds** - Animated abstract visuals with customizable effects
- 📝 **Auto Subtitles** - Perfectly timed captions burned into video
- 🎵 **Audio Mixing** - Background music with smart ducking
- 🏷️ **Brand Integration** - Logo overlay and channel branding
- ⚡ **Fast Rendering** - Optimized FFmpeg pipeline for quick exports

## 🚀 Quick Start

### Prerequisites

**macOS:**
```bash
brew install ffmpeg espeak-ng ollama
pip3 install TTS
```

**Ubuntu/Debian:**
```bash
sudo apt update && sudo apt install -y ffmpeg espeak-ng python3-pip
pip3 install TTS
curl -fsSL https://ollama.com/install.sh | sh
```

### Installation

```bash
git clone https://github.com/g-laliotis/convertbox.git
cd convertbox
go mod download
cp .env.example .env
```

### Usage

```bash
# Pull AI model (first time only)
ollama pull mistral

# Generate a video
go run ./cmd/convertbox --topic "5 AI Tools That Will Blow Your Mind in 2025"

# Or use make
make demo
```

## 📁 Project Structure

```
convertbox/
├── cmd/convertbox/          # Main CLI application
├── internal/
│   ├── config/             # Configuration management
│   ├── llm/                # Local LLM integration
│   ├── tts/                # Text-to-speech engines
│   ├── media/              # Video/audio processing
│   └── logger/             # Structured logging
├── assets/
│   ├── music/              # Background music tracks
│   └── logo.png            # Channel logo
├── build/                  # Generated content (gitignored)
└── scripts/                # Utility scripts
```

## 🛠️ Configuration

Edit `.env` to customize:
- AI model settings
- TTS preferences  
- Video quality/format
- Branding elements

## 📊 Output

Generated videos include:
- ✅ Engaging hook (first 5 seconds)
- ✅ Main content (45-50 seconds)
- ✅ Call-to-action outro
- ✅ Burned-in subtitles
- ✅ Background music with ducking
- ✅ Channel logo overlay
- ✅ Optimized for YouTube Shorts (9:16 aspect ratio)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Ollama](https://ollama.ai/) for local LLM inference
- [Coqui TTS](https://github.com/coqui-ai/TTS) for high-quality speech synthesis
- [FFmpeg](https://ffmpeg.org/) for video processing

---

**Built with ❤️ for AI Unboxed by UnboxGio**