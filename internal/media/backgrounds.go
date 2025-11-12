package media

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// BackgroundSegment represents a timed background change
type BackgroundSegment struct {
	StartTime time.Duration
	EndTime   time.Duration
	ImagePath string
	Keywords  []string
}

// CreateDynamicBackground creates a video with changing backgrounds based on script content
func (s *Service) CreateDynamicBackground(ctx context.Context, script, outPath string, duration time.Duration) error {
	s.logger.Info("Creating dynamic background based on script content")

	// Analyze script and create segments
	segments := s.analyzeScriptForBackgrounds(script, duration)
	
	if len(segments) == 0 {
		// Fallback to static background
		return s.CreateBackground(ctx, outPath, duration)
	}

	// Create individual background videos for each segment
	var segmentPaths []string
	for i, segment := range segments {
		segmentPath := fmt.Sprintf("build/segment_%d.mp4", i)
		segmentDuration := segment.EndTime - segment.StartTime
		
		if err := s.createSegmentBackground(ctx, segment, segmentPath, segmentDuration); err != nil {
			s.logger.Warning("Failed to create segment %d, using fallback", i)
			// Create fallback segment
			if err := s.createFallbackSegment(ctx, segmentPath, segmentDuration); err != nil {
				return err
			}
		}
		segmentPaths = append(segmentPaths, segmentPath)
	}

	// Concatenate all segments
	if err := s.concatenateSegments(ctx, segmentPaths, outPath); err != nil {
		// If concatenation fails, use first segment as fallback
		if len(segmentPaths) > 0 {
			return exec.CommandContext(ctx, "cp", segmentPaths[0], outPath).Run()
		}
		return err
	}
	return nil
}

func (s *Service) analyzeScriptForBackgrounds(script string, totalDuration time.Duration) []BackgroundSegment {
	// Clean script
	cleanScript := strings.ReplaceAll(script, `"`, "")
	words := strings.Fields(strings.ToLower(cleanScript))
	
	// Define keyword mappings to image categories
	keywordMap := map[string]string{
		"artificial":   "ai",
		"intelligence": "ai", 
		"ai":          "ai",
		"robot":       "ai",
		"machine":     "ai",
		"neural":      "ai",
		"deep":        "ai",
		"learning":    "ai",
		"algorithm":   "tech",
		"data":        "tech",
		"computer":    "tech",
		"digital":     "tech",
		"technology":  "tech",
		"software":    "tech",
		"code":        "tech",
		"programming": "tech",
		"tool":        "tools",
		"tools":       "tools",
		"app":         "tools",
		"application": "tools",
		"platform":    "tools",
		"service":     "tools",
	}

	// Split script into roughly equal time segments
	segmentCount := 3 // 3 background changes per video
	segmentDuration := totalDuration / time.Duration(segmentCount)
	wordsPerSegment := len(words) / segmentCount
	
	var segments []BackgroundSegment
	
	for i := 0; i < segmentCount; i++ {
		startTime := time.Duration(i) * segmentDuration
		endTime := startTime + segmentDuration
		if i == segmentCount-1 {
			endTime = totalDuration
		}

		// Analyze words in this segment
		startWord := i * wordsPerSegment
		endWord := startWord + wordsPerSegment
		if endWord > len(words) {
			endWord = len(words)
		}

		category := s.detectCategory(words[startWord:endWord], keywordMap)
		imagePath := s.findBestImage(category)

		segments = append(segments, BackgroundSegment{
			StartTime: startTime,
			EndTime:   endTime,
			ImagePath: imagePath,
			Keywords:  words[startWord:endWord],
		})
	}

	return segments
}

func (s *Service) detectCategory(words []string, keywordMap map[string]string) string {
	categoryCount := make(map[string]int)
	
	for _, word := range words {
		if category, exists := keywordMap[word]; exists {
			categoryCount[category]++
		}
	}

	// Find most frequent category
	maxCount := 0
	bestCategory := "tech" // default
	for category, count := range categoryCount {
		if count > maxCount {
			maxCount = count
			bestCategory = category
		}
	}

	return bestCategory
}

func (s *Service) findBestImage(category string) string {
	// Look for images in category-specific folders
	categoryPaths := []string{
		fmt.Sprintf("assets/images/%s", category),
		"assets/images/tech", // fallback
		"assets/images",      // general fallback
	}

	for _, dir := range categoryPaths {
		// Check each extension separately
		extensions := []string{"*.jpg", "*.jpeg", "*.png"}
		for _, ext := range extensions {
			if files, err := filepath.Glob(filepath.Join(dir, ext)); err == nil && len(files) > 0 {
				return files[0] // Return first available image
			}
		}
	}

	return "" // No image found
}

func (s *Service) createSegmentBackground(ctx context.Context, segment BackgroundSegment, outPath string, duration time.Duration) error {
	if segment.ImagePath == "" {
		return fmt.Errorf("no image path provided")
	}

	sec := int(duration.Seconds())
	
	// Create zooming/panning effect based on segment position
	zoomEffect := "zoompan=z='min(zoom+0.002,1.8)':d=125:x='iw/2-(iw/zoom/2)':y='ih/2-(ih/zoom/2)'"
	
	cmd := exec.CommandContext(ctx, "ffmpeg", "-y",
		"-loop", "1", "-i", segment.ImagePath,
		"-t", fmt.Sprintf("%d", sec),
		"-vf", fmt.Sprintf("scale=1080:1920:force_original_aspect_ratio=increase,crop=1080:1920,%s", zoomEffect),
		"-c:v", "libx264", "-preset", "ultrafast", "-pix_fmt", "yuv420p",
		outPath,
	)
	
	return cmd.Run()
}

func (s *Service) createFallbackSegment(ctx context.Context, outPath string, duration time.Duration) error {
	sec := int(duration.Seconds())
	
	cmd := exec.CommandContext(ctx, "ffmpeg", "-y",
		"-f", "lavfi", "-t", fmt.Sprintf("%d", sec),
		"-i", "color=c=#0f0f23:s=1080x1920",
		"-c:v", "libx264", "-preset", "ultrafast", "-pix_fmt", "yuv420p",
		outPath,
	)
	
	return cmd.Run()
}

func (s *Service) concatenateSegments(ctx context.Context, segmentPaths []string, outPath string) error {
	// Create concat file
	concatFile := "build/concat.txt"
	var concatContent strings.Builder
	
	for _, path := range segmentPaths {
		fmt.Fprintf(&concatContent, "file '%s'\n", path)
	}
	
	if err := os.WriteFile(concatFile, []byte(concatContent.String()), 0644); err != nil {
		return err
	}

	// Concatenate using FFmpeg
	cmd := exec.CommandContext(ctx, "ffmpeg", "-y",
		"-f", "concat", "-safe", "0", "-i", concatFile,
		"-c", "copy",
		outPath,
	)
	
	return cmd.Run()
}