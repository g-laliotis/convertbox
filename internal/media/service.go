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

func (s *Service) CreateBackground(ctx context.Context, outPath string, duration time.Duration) error {
	s.logger.Info("Creating animated background (%v duration)", duration)

	sec := int(duration.Seconds())
	filter := fmt.Sprintf(
		"life=s=%dx%d:mold=10:ratio=0.08:seed=42,format=yuv420p,eq=contrast=1.2:brightness=0.03,boxblur=2",
		s.config.VideoWidth, s.config.VideoHeight,
	)

	cmd := exec.CommandContext(ctx, "ffmpeg", "-y",
		"-f", "lavfi", "-t", fmt.Sprintf("%d", sec),
		"-i", "color=c=black:s=1920x1080",
		"-f", "lavfi", "-t", fmt.Sprintf("%d", sec),
		"-i", "nullsrc=s=1920x1080",
		"-filter_complex", filter,
		"-s", fmt.Sprintf("%dx%d", s.config.VideoWidth, s.config.VideoHeight),
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

	sentences := s.splitSentences(script)
	if len(sentences) == 0 {
		sentences = []string{strings.TrimSpace(script)}
	}

	segmentDuration := dur / time.Duration(len(sentences))
	var srt strings.Builder
	currentTime := time.Duration(0)

	for i, sentence := range sentences {
		nextTime := currentTime + segmentDuration
		if i == len(sentences)-1 {
			nextTime = dur
		}

		fmt.Fprintf(&srt, "%d\n%s --> %s\n%s\n\n",
			i+1,
			s.formatSRTTime(currentTime),
			s.formatSRTTime(nextTime),
			strings.TrimSpace(sentence),
		)
		currentTime = nextTime
	}

	return os.WriteFile(outPath, []byte(srt.String()), 0644)
}

func (s *Service) RenderVideo(ctx context.Context, cfg RenderConfig) error {
	s.logger.Info("Rendering final video")

	args := []string{"-y"}
	
	// Add video inputs
	for _, v := range cfg.VideoInputs {
		args = append(args, "-stream_loop", "-1", "-t", "65", "-i", v)
	}
	
	// Add narration
	args = append(args, "-i", cfg.Narration)
	
	// Add music if provided
	if cfg.Music != "" {
		args = append(args, "-i", cfg.Music)
	}
	
	// Add logo if provided
	if cfg.Logo != "" {
		args = append(args, "-i", cfg.Logo)
	}

	// Build filter complex
	var fc strings.Builder
	fmt.Fprintf(&fc, "[0:v]scale=%d:%d,setsar=1:1,format=yuv420p[v0];", s.config.VideoWidth, s.config.VideoHeight)
	
	if cfg.Music != "" {
		fc.WriteString("[2:a]aformat=fltp:44100:stereo,volume=0.5[music];[1:a]anull[narr];")
		fc.WriteString("[music][narr]sidechaincompress=threshold=0.12:ratio=10:attack=5:release=200[aout];")
	} else {
		fc.WriteString("[1:a]anull[aout];")
	}
	
	if cfg.Logo != "" {
		fmt.Fprintf(&fc, "[v0][3]overlay=W-w-%d:H-h-%d:format=auto[vout];", s.config.LogoMargin, s.config.LogoMargin)
	} else {
		fc.WriteString("[v0]null[vout];")
	}

	args = append(args,
		"-filter_complex", fc.String(),
		"-map", "[vout]", "-map", "[aout]",
		"-c:v", "libx264", "-preset", s.config.VideoPreset, "-crf", fmt.Sprintf("%d", s.config.VideoCRF),
		"-c:a", "aac", "-b:a", "192k",
		"-shortest",
		"-vf", fmt.Sprintf("subtitles=%s", cfg.CaptionsSRT),
		cfg.Output,
	)

	cmd := exec.CommandContext(ctx, "ffmpeg", args...)
	return cmd.Run()
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