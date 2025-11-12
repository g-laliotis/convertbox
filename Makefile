.PHONY: build run demo clean install deps

# Build the application
build:
	go build -o bin/convertbox ./cmd/convertbox

# Run with custom topic
run:
	go run ./cmd/convertbox --topic "$(TOPIC)"

# Demo with sample topic
demo:
	@echo "🚀 Running Convertbox demo..."
	@ollama pull mistral || echo "⚠️  Ollama not available, continuing..."
	@$(MAKE) run TOPIC="5 Hidden AI Websites That Will Blow Your Mind in 2025"

# Install dependencies
deps:
	go mod download
	go mod tidy

# Clean build artifacts
clean:
	rm -rf build/ bin/

# Install system dependencies (macOS)
install-macos:
	brew install ffmpeg espeak-ng ollama
	pip3 install TTS

# Install system dependencies (Ubuntu/Debian)
install-ubuntu:
	sudo apt update
	sudo apt install -y ffmpeg espeak-ng python3-pip
	pip3 install TTS
	curl -fsSL https://ollama.com/install.sh | sh

# Setup project
setup:
	cp .env.example .env
	mkdir -p build assets/music assets/logos
	@echo "✅ Project setup complete!"
	@echo "📝 Edit .env file to customize settings"
	@echo "🎵 Add music files to assets/music/"
	@echo "🏷️  Add logo to assets/logos/logo.png"