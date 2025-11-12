package media

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/g-laliotis/convertbox/internal/config"
	"github.com/g-laliotis/convertbox/internal/logger"
)

type Service struct {
	config *config.Config
	logger *logger.Logger
}

type RenderConfig struct {
	VideoInputs []string
	Narration   string
	Music       string
	Logo        string
	CaptionsSRT string
	Output      string
}

func NewService(cfg *config.Config, log *logger.Logger) *Service {
	return &Service{
		config: cfg,
		logger: log,
	}
}

// CreateBackground creates background video from images or generates one
func (s *Service) CreateBackground(ctx context.Context, outPath string, duration time.Duration) error {
	s.logger.Info("Creating background (%v duration)", duration)

	sec := int(duration.Seconds())
	
	// Check for background images first
	imageFiles := []string{
		"assets/images/tech1.jpg",
		"assets/images/tech2.jpg", 
		"assets/images/tech3.jpg",
	}
	
	// Use images if available
	for _, img := range imageFiles {
		if _, err := os.Stat(img); err == nil {
			s.logger.Info("Using background image: %s", img)
			cmd := exec.CommandContext(ctx, "ffmpeg", "-y",
				"-loop", "1", "-i", img,
				"-t", fmt.Sprintf("%d", sec),
				"-vf", "scale=1080:1920:force_original_aspect_ratio=increase,crop=1080:1920,zoompan=z='min(zoom+0.0015,1.5)':d=125",
				"-c:v", "libx264", "-preset", "ultrafast", "-pix_fmt", "yuv420p",
				outPath,
			)
			return cmd.Run()
		}
	}
	
	// Fallback to simple gradient
	s.logger.Info("No images found, creating simple gradient")
	cmd := exec.CommandContext(ctx, "ffmpeg", "-y",
		"-f", "lavfi", "-t", fmt.Sprintf("%d", sec),
		"-i", "color=c=#0f0f23:s=1080x1920",
		"-c:v", "libx264", "-preset", "ultrafast", "-pix_fmt", "yuv420p",
		outPath,
	)
	
	return cmd.Run()
}

func (s *Service) GenerateSubtitles(audioPath, script, outPath string) error {
	s.logger.Info("Generating subtitles")

	dur, err := s.getAudioDuration(audioPath)
	if err != nil {
		return err
	}

	// Clean script and split by words for better sync
	cleanScript := strings.ReplaceAll(script, `"`, "")
	words := strings.Fields(cleanScript)
	if len(words) == 0 {
		return fmt.Errorf("no words found in script")
	}

	// Group words for better readability (3-5 words per subtitle)
	var groups []string
	wordsPerGroup := 4
	for i := 0; i < len(words); i += wordsPerGroup {
		end := i + wordsPerGroup
		if end > len(words) {
			end = len(words)
		}
		groups = append(groups, strings.Join(words[i:end], " "))
	}

	segmentDuration := dur / time.Duration(len(groups))
	var srt strings.Builder
	currentTime := time.Duration(0)

	for i, group := range groups {
		nextTime := currentTime + segmentDuration
		if i == len(groups)-1 {
			nextTime = dur
		}

		fmt.Fprintf(&srt, "%d\n%s --> %s\n%s\n\n",
			i+1,
			s.formatSRTTime(currentTime),
			s.formatSRTTime(nextTime),
			strings.TrimSpace(group),
		)
		currentTime = nextTime
	}

	return os.WriteFile(outPath, []byte(srt.String()), 0644)
}

func (s *Service) RenderVideo(ctx context.Context, cfg RenderConfig) error {
	s.logger.Info("Rendering final video")

	// Step 1: Add subtitles to video
	tempVideo := "build/temp_with_subs.mp4"
	cmd1 := exec.CommandContext(ctx, "ffmpeg", "-y",
		"-i", cfg.VideoInputs[0],
		"-vf", fmt.Sprintf("subtitles=%s", cfg.CaptionsSRT),
		"-c:v", "libx264", "-preset", "fast", "-crf", "20",
		tempVideo,
	)
	if err := cmd1.Run(); err != nil {
		return err
	}

	// Step 2: Add logo if available
	videoWithLogo := tempVideo
	if cfg.Logo != "" {
		videoWithLogo = "build/temp_with_logo.mp4"
		cmd2 := exec.CommandContext(ctx, "ffmpeg", "-y",
			"-i", tempVideo,
			"-i", cfg.Logo,
			"-filter_complex", "[0:v][1:v]overlay=W-w-20:20",
			"-c:v", "libx264", "-preset", "fast", "-crf", "20",
			videoWithLogo,
		)
		if err := cmd2.Run(); err != nil {
			return err
		}
	}

	// Step 3: Add audio
	args := []string{"-y",
		"-i", videoWithLogo,
		"-i", cfg.Narration,
	}
	
	if cfg.Music != "" {
		args = append(args, "-i", cfg.Music,
			"-filter_complex", "[1:a]volume=5.0[narr];[2:a]volume=1.0[music];[narr][music]amix=inputs=2:duration=first",
			"-c:a", "aac", "-b:a", "192k",
		)
	} else {
		args = append(args,
			"-c:a", "aac", "-b:a", "192k",
			"-filter:a", "volume=5.0",
		)
	}
	
	args = append(args, "-c:v", "copy", "-shortest", cfg.Output)
	
	cmd3 := exec.CommandContext(ctx, "ffmpeg", args...)
	return cmd3.Run()
}

func (s *Service) getAudioDuration(path string) (time.Duration, error) {
	out, err := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1", path).Output()
	if err != nil {
		return 0, err
	}
	
	var seconds float64
	fmt.Sscanf(strings.TrimSpace(string(out)), "%f", &seconds)
	return time.Duration(seconds*1000) * time.Millisecond, nil
}

var sentenceRegex = regexp.MustCompile(`(?m)([^.!?]+[.!?])`)

func (s *Service) splitSentences(text string) []string {
	matches := sentenceRegex.FindAllString(text, -1)
	if len(matches) == 0 {
		return []string{text}
	}
	return matches
}

func (s *Service) formatSRTTime(d time.Duration) string {
	h := int(d / time.Hour)
	m := int(d/time.Minute) % 60
	sec := int(d/time.Second) % 60
	ms := int(d/time.Millisecond) % 1000
	return fmt.Sprintf("%02d:%02d:%02d,%03d", h, m, sec, ms)
}