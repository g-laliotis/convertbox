package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	topic := flag.String("topic", "", "Video topic")
	flag.Parse()
	
	if *topic == "" {
		fmt.Println("Usage: go run ./cmd/simple --topic \"Your topic\"")
		os.Exit(1)
	}

	fmt.Printf("🎬 Creating video: %s\n", *topic)
	os.MkdirAll("build", 0755)

	// 1. Generate script
	fmt.Println("📝 Generating script...")
	script := generateScript(*topic)
	os.WriteFile("build/script.txt", []byte(script), 0644)
	fmt.Printf("✅ Script: %d words\n", len(strings.Fields(script)))

	// 2. Text to speech
	fmt.Println("🎙️ Converting to speech...")
	textToSpeech(script)
	fmt.Println("✅ Speech generated")

	// 3. Create subtitles
	fmt.Println("📝 Creating subtitles...")
	createSubtitles(script)
	fmt.Println("✅ Subtitles created")

	// 4. Random background
	fmt.Println("🖼️ Creating background...")
	createRandomBackground()
	fmt.Println("✅ Background created")

	// 5. Assemble video
	fmt.Println("🎬 Assembling video...")
	assembleVideo()
	fmt.Println("🎉 Video ready: build/final.mp4")
}

func generateScript(topic string) string {
	prompt := fmt.Sprintf("Write a 60-second YouTube script about: %s. Make it energetic for A.I. Unboxed by UnboxGio. 140-160 words. End with 'Subscribe for more A.I. insights!' Always write 'A.I.' with periods. Return only script text.", topic)

	cmd := exec.Command("ollama", "run", "mistral", prompt)
	out, _ := cmd.Output()
	script := strings.TrimSpace(string(out))
	if len(script) == 0 {
		script = fmt.Sprintf("Welcome to A.I. Unboxed by UnboxGio! Today we explore %s. This amazing topic in artificial intelligence is revolutionizing our world. These tools boost productivity and creativity. The future is powered by A.I. Subscribe for more A.I. insights!", topic)
	}
	// Fix AI pronunciation
	script = strings.ReplaceAll(script, " AI ", " A.I. ")
	script = strings.ReplaceAll(script, "AI ", "A.I. ")
	script = strings.ReplaceAll(script, " AI.", " A.I.")
	return script
}

func textToSpeech(text string) {
	// Try Coqui with better model
	cmd := exec.Command("tts", "--text", text, "--model_name", "tts_models/en/ljspeech/tacotron2-DDC", "--out_path", "build/speech.wav")
	if cmd.Run() != nil {
		// Fallback to eSpeak with better settings
		exec.Command("espeak-ng", "-v", "en-us+f3", "-s", "150", "-p", "50", "-w", "build/speech.wav", text).Run()
	}
}

func createSubtitles(script string) {
	// Get actual audio duration
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", "build/speech.wav")
	out, _ := cmd.Output()
	duration := 60.0
	if len(out) > 0 {
		fmt.Sscanf(strings.TrimSpace(string(out)), "%f", &duration)
	}
	
	words := strings.Fields(script)
	wordsPerSub := 3 // Fewer words per subtitle
	
	var srt strings.Builder
	subIndex := 1
	
	for i := 0; i < len(words); i += wordsPerSub {
		end := i + wordsPerSub
		if end > len(words) {
			end = len(words)
		}
		
		// Better timing calculation
		startTime := (float64(i) / float64(len(words))) * duration
		endTime := (float64(end) / float64(len(words))) * duration
		
		fmt.Fprintf(&srt, "%d\n%s --> %s\n%s\n\n",
			subIndex,
			formatSRTTime(startTime),
			formatSRTTime(endTime),
			strings.Join(words[i:end], " "))
		subIndex++
	}
	
	os.WriteFile("build/subtitles.srt", []byte(srt.String()), 0644)
}

func formatSRTTime(seconds float64) string {
	h := int(seconds) / 3600
	m := (int(seconds) % 3600) / 60
	s := int(seconds) % 60
	ms := int((seconds - float64(int(seconds))) * 1000)
	return fmt.Sprintf("%02d:%02d:%02d,%03d", h, m, s, ms)
}

func createRandomBackground() {
	rand.Seed(time.Now().UnixNano())
	
	// Create multiple tech images
	images := []string{
		"build/tech1.jpg", "build/tech2.jpg", "build/tech3.jpg",
	}
	colors := []string{
		"#1a1a2e-#16213e", "#0f3460-#16537e", "#2d1b69-#11998e",
	}
	
	// Create random tech backgrounds
	for i, img := range images {
		exec.Command("magick", "-size", "1080x1920", 
			fmt.Sprintf("gradient:%s", colors[i]), 
			"-blur", "0x8", 
			"-noise", "1", 
			img).Run()
	}
	
	// Pick random image
	selectedImage := images[rand.Intn(len(images))]
	
	// Create video with zoom effect
	exec.Command("ffmpeg", "-y", 
		"-loop", "1", "-i", selectedImage, 
		"-t", "65", 
		"-vf", "zoompan=z='min(zoom+0.002,1.8)':d=125:x='iw/2-(iw/zoom/2)':y='ih/2-(ih/zoom/2)'", 
		"-c:v", "libx264", "-preset", "fast", "-pix_fmt", "yuv420p",
		"build/background.mp4").Run()
}

func assembleVideo() {
	// Create better random music
	rand.Seed(time.Now().UnixNano())
	freqs := []int{220, 330, 440, 550}
	freq1 := freqs[rand.Intn(len(freqs))]
	freq2 := freqs[rand.Intn(len(freqs))]
	
	exec.Command("ffmpeg", "-y", 
		"-f", "lavfi", "-i", fmt.Sprintf("sine=frequency=%d:duration=65", freq1),
		"-f", "lavfi", "-i", fmt.Sprintf("sine=frequency=%d:duration=65", freq2),
		"-filter_complex", "[0:a][1:a]amix=inputs=2:duration=shortest,volume=0.15",
		"build/music.wav").Run()
	
	// Create logo if doesn't exist
	if _, err := os.Stat("build/logo.png"); os.IsNotExist(err) {
		exec.Command("magick", "-size", "200x200", "-background", "transparent", 
			"-fill", "white", "-pointsize", "24", "-gravity", "center", 
			"label:A.I. UNBOXED", "build/logo.png").Run()
	}
	
	// Step 1: Add subtitles
	exec.Command("ffmpeg", "-y", 
		"-i", "build/background.mp4", 
		"-vf", "subtitles=build/subtitles.srt",
		"-c:v", "libx264", "-preset", "fast",
		"build/video_subs.mp4").Run()
	
	// Step 2: Add logo
	exec.Command("ffmpeg", "-y",
		"-i", "build/video_subs.mp4",
		"-i", "build/logo.png",
		"-filter_complex", "[0:v][1:v]overlay=W-w-20:20",
		"-c:v", "libx264", "-preset", "fast",
		"build/video_logo.mp4").Run()
	
	// Step 3: Add audio with better mixing
	exec.Command("ffmpeg", "-y", 
		"-i", "build/video_logo.mp4", 
		"-i", "build/speech.wav", 
		"-i", "build/music.wav", 
		"-filter_complex", "[1:a]volume=4.0[speech];[2:a]volume=1.0[music];[speech][music]amix=inputs=2:duration=first",
		"-c:v", "copy", "-c:a", "aac", "-b:a", "192k",
		"-shortest", "build/final.mp4").Run()
}