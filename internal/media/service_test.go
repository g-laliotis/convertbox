package media

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/g-laliotis/convertbox/internal/config"
	"github.com/g-laliotis/convertbox/internal/logger"
)

func TestService_CreateBackground(t *testing.T) {
	cfg := &config.Config{VideoWidth: 1080, VideoHeight: 1920}
	log := logger.New()
	service := NewService(cfg, log)

	tmpFile := "test_bg.mp4"
	defer os.Remove(tmpFile)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := service.CreateBackground(ctx, tmpFile, 5*time.Second)
	if err != nil {
		t.Fatalf("CreateBackground failed: %v", err)
	}

	if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
		t.Fatal("Background video was not created")
	}
}

func TestService_GenerateSubtitles(t *testing.T) {
	cfg := &config.Config{}
	log := logger.New()
	service := NewService(cfg, log)

	// Create a dummy audio file for testing
	audioFile := "test_audio.wav"
	srtFile := "test_subs.srt"
	defer os.Remove(audioFile)
	defer os.Remove(srtFile)

	// Create minimal WAV file (1 second of silence)
	wavData := []byte{
		0x52, 0x49, 0x46, 0x46, 0x24, 0x08, 0x00, 0x00, 0x57, 0x41, 0x56, 0x45, 0x66, 0x6d, 0x74, 0x20,
		0x10, 0x00, 0x00, 0x00, 0x01, 0x00, 0x02, 0x00, 0x22, 0x56, 0x00, 0x00, 0x88, 0x58, 0x01, 0x00,
		0x04, 0x00, 0x10, 0x00, 0x64, 0x61, 0x74, 0x61, 0x00, 0x08, 0x00, 0x00,
	}
	// Add 1 second of silence (44100 samples * 2 channels * 2 bytes)
	silence := make([]byte, 44100*2*2)
	wavData = append(wavData, silence...)
	
	if err := os.WriteFile(audioFile, wavData, 0644); err != nil {
		t.Fatalf("Failed to create test audio: %v", err)
	}

	script := "Hello world. This is a test. How are you?"
	err := service.GenerateSubtitles(audioFile, script, srtFile)
	if err != nil {
		t.Fatalf("GenerateSubtitles failed: %v", err)
	}

	if _, err := os.Stat(srtFile); os.IsNotExist(err) {
		t.Fatal("SRT file was not created")
	}
}

func TestService_SplitSentences(t *testing.T) {
	cfg := &config.Config{}
	log := logger.New()
	service := NewService(cfg, log)

	tests := []struct {
		input    string
		expected int
	}{
		{"Hello world.", 1},
		{"Hello world. How are you?", 2},
		{"One. Two! Three?", 3},
		{"No punctuation", 1},
	}

	for _, test := range tests {
		result := service.splitSentences(test.input)
		if len(result) != test.expected {
			t.Errorf("splitSentences(%q) = %d sentences, want %d", test.input, len(result), test.expected)
		}
	}
}