package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"

	"github.com/g-laliotis/convertbox/internal/config"
	"github.com/g-laliotis/convertbox/internal/llm"
	"github.com/g-laliotis/convertbox/internal/logger"
	"github.com/g-laliotis/convertbox/internal/media"
	"github.com/g-laliotis/convertbox/internal/tts"
)

func main() {
	// Load environment variables
	_ = godotenv.Load()

	// Parse command line flags
	topic := flag.String("topic", "", "Video topic/title (required)")
	output := flag.String("out", "build/final.mp4", "Output video path")
	test := flag.Bool("test", false, "Run quick test mode")
	flag.Parse()

	if *topic == "" {
		fmt.Println("Usage: convertbox --topic \"Your video topic\"")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Initialize services
	cfg := config.Load()
	log := logger.New()
	
	log.Info("🎬 Starting Convertbox for %s", cfg.ChannelName)
	log.Info("Topic: %s", *topic)

	// Create build directory
	if err := os.MkdirAll("build", 0755); err != nil {
		log.Error("Failed to create build directory: %v", err)
		os.Exit(1)
	}

	ctx := context.Background()

	// Initialize services
	llmService := llm.NewService(cfg, log)
	ttsService := tts.NewService(cfg, log)
	mediaService := media.NewService(cfg, log)

	// Step 1: Generate script
	log.Info("Step 1/5: Generating script...")
	script, err := llmService.GenerateScript(ctx, *topic)
	if err != nil {
		log.Error("Script generation failed: %v", err)
		os.Exit(1)
	}

	scriptPath := "build/script.txt"
	if err := os.WriteFile(scriptPath, []byte(script), 0644); err != nil {
		log.Error("Failed to save script: %v", err)
		os.Exit(1)
	}
	log.Success("Script saved (%d chars)", len(script))

	// Step 2: Generate narration
	log.Info("Step 2/5: Synthesizing narration...")
	narrationPath := "build/narration.wav"
	if err := ttsService.Synthesize(ctx, script, narrationPath); err != nil {
		log.Error("TTS synthesis failed: %v", err)
		os.Exit(1)
	}
	log.Success("Narration synthesized")

	// Step 3: Create background video
	log.Info("Step 3/5: Creating background...")
	backgroundPath := "build/background.mp4"
	duration := 65 * time.Second
	if *test {
		duration = 10 * time.Second // Faster for testing
	}
	
	if err := mediaService.CreateBackground(ctx, backgroundPath, duration); err != nil {
		log.Error("Background creation failed: %v", err)
		os.Exit(1)
	}
	log.Success("Background created")

	// Step 4: Generate subtitles
	log.Info("Step 4/5: Generating subtitles...")
	subtitlesPath := "build/subtitles.srt"
	if err := mediaService.GenerateSubtitles(narrationPath, script, subtitlesPath); err != nil {
		log.Error("Subtitle generation failed: %v", err)
		os.Exit(1)
	}
	log.Success("Subtitles generated")

	// Step 5: Render final video
	log.Info("Step 5/5: Rendering final video...")
	renderCfg := media.RenderConfig{
		VideoInputs: []string{backgroundPath},
		Narration:   narrationPath,
		CaptionsSRT: subtitlesPath,
		Output:      *output,
	}

	// Add logo if exists
	if _, err := os.Stat("assets/logos/logo.png"); err == nil {
		renderCfg.Logo = "assets/logos/logo.png"
		log.Info("Using logo overlay")
	}

	// Add music if exists
	musicFiles := []string{
		"assets/music/background.mp3",
		"assets/music/background.wav",
	}
	for _, musicFile := range musicFiles {
		if _, err := os.Stat(musicFile); err == nil {
			renderCfg.Music = musicFile
			log.Info("Using background music: %s", musicFile)
			break
		}
	}

	if err := mediaService.RenderVideo(ctx, renderCfg); err != nil {
		log.Error("Video rendering failed: %v", err)
		os.Exit(1)
	}

	log.Success("🎉 Video generated successfully!")
	log.Info("Output: %s", filepath.ToSlash(*output))
	log.Info("Ready to upload to %s! 🚀", cfg.ChannelName)
}