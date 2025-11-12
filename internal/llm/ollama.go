package llm

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
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

func (s *Service) GenerateScript(ctx context.Context, topic string) (string, error) {
	s.logger.Info("Generating script for topic: %s", topic)

	prompt := s.buildPrompt(topic)
	
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "ollama", "run", s.config.OllamaModel, prompt)
	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ollama execution failed: %w\nStderr: %s", err, errOut.String())
	}

	script := strings.TrimSpace(out.String())
	if script == "" {
		return "", fmt.Errorf("ollama returned empty response")
	}

	s.logger.Success("Script generated successfully (%d characters)", len(script))
	return script, nil
}

func (s *Service) buildPrompt(topic string) string {
	return fmt.Sprintf(`You are a professional YouTube script writer for "%s", a cutting-edge tech channel focused on AI innovations.

TOPIC: %s

REQUIREMENTS:
- Write EXACTLY 140-160 words for ~60 seconds of speech
- Structure: Hook (5s) → Main Content (50s) → CTA (5s)
- Tone: Energetic, curious, authoritative but accessible
- Use short, punchy sentences with natural pauses
- Include specific numbers, facts, or examples when possible
- End with "Don't forget to subscribe for more AI insights!"

STYLE GUIDELINES:
- Start with an attention-grabbing question or bold statement
- Use "you" to directly address viewers
- Avoid technical jargon - explain complex concepts simply
- Create urgency and excitement about AI developments
- Include transition phrases like "But here's the thing..." or "What's even crazier..."

IMPORTANT: Return ONLY the actual script text that will be spoken. Do NOT include:
- Title headers
- Section labels like "(Hook)" or "(Main Content)" or "(CTA)"
- Any formatting or commentary
- Just the pure spoken script text

OUTPUT: Return ONLY the script text, no additional formatting or commentary.`, 
		s.config.ChannelName, topic)
}