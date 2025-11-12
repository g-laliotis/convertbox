package tts

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/g-laliotis/convertbox/internal/config"
	"github.com/g-laliotis/convertbox/internal/logger"
)

type Service struct {
	config *config.Config
	logger *logger.Logger
}

func NewService(cfg *config.Config, log *logger.Logger) *Service {
	return &Service{
		config: cfg,
		logger: log,
	}
}

func (s *Service) Synthesize(ctx context.Context, text, outPath string) error {
	s.logger.Info("Synthesizing speech (%d chars)", len(text))

	if s.config.TTSEngine == "coqui" {
		if err := s.coquiSpeak(ctx, text, outPath); err != nil {
			s.logger.Warning("Coqui failed, falling back to eSpeak: %v", err)
			return s.eSpeak(ctx, text, outPath)
		}
		return nil
	}
	return s.eSpeak(ctx, text, outPath)
}

func (s *Service) coquiSpeak(ctx context.Context, text, outPath string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "tts",
		"--text", text,
		"--model_name", s.config.CoquiModel,
		"--out_path", outPath,
	)
	return cmd.Run()
}

func (s *Service) eSpeak(ctx context.Context, text, outPath string) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "espeak-ng",
		"-v", s.config.ESpeakVoice,
		"-s", fmt.Sprintf("%d", s.config.ESpeakSpeed),
		"-w", outPath,
		text,
	)
	return cmd.Run()
}