package config

import (
	"os"
	"strconv"
)

type Config struct {
	// LLM Configuration
	OllamaModel string
	OllamaHost  string

	// TTS Configuration
	TTSEngine    string
	CoquiModel   string
	ESpeakVoice  string
	ESpeakSpeed  int

	// Video Configuration
	VideoWidth   int
	VideoHeight  int
	VideoCRF     int
	VideoPreset  string
	LogoMargin   int

	// Branding
	ChannelName string
}

func Load() *Config {
	return &Config{
		OllamaModel:  getEnv("OLLAMA_MODEL", "mistral"),
		OllamaHost:   getEnv("OLLAMA_HOST", "http://localhost:11434"),
		TTSEngine:    getEnv("TTS_ENGINE", "coqui"),
		CoquiModel:   getEnv("COQUI_MODEL", "tts_models/en/vctk/vits"),
		ESpeakVoice:  getEnv("ESPEAK_VOICE", "en-us"),
		ESpeakSpeed:  getEnvInt("ESPEAK_SPEED", 160),
		VideoWidth:   getEnvInt("VIDEO_WIDTH", 1080),
		VideoHeight:  getEnvInt("VIDEO_HEIGHT", 1920),
		VideoCRF:     getEnvInt("VIDEO_CRF", 18),
		VideoPreset:  getEnv("VIDEO_PRESET", "veryfast"),
		LogoMargin:   getEnvInt("LOGO_MARGIN", 40),
		ChannelName:  getEnv("CHANNEL_NAME", "AI Unboxed by UnboxGio"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}