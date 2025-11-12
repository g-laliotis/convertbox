# 🎨 Assets Guide for AI Unboxed by UnboxGio

## 📁 Directory Structure

```
assets/
├── logos/                    # Channel logos
│   ├── logo.png             # Main logo (recommended: 200x200px)
│   ├── main_logo.png        # Alternative logo
│   └── logo.jpg             # JPEG version
├── banners/                  # Channel banners  
│   ├── channel_logo.png     # Banner logo
│   └── intro_banner.png     # Intro graphics
├── images/
│   ├── ai/                  # AI-related backgrounds
│   │   ├── ai1.jpg          # Neural networks, robots, AI concepts
│   │   ├── ai2.jpg          # Machine learning visuals
│   │   └── ai3.jpg          # Futuristic AI imagery
│   ├── tech/                # General tech backgrounds
│   │   ├── tech1.jpg        # Code, circuits, computers
│   │   ├── tech2.jpg        # Data visualization
│   │   └── tech3.jpg        # Digital interfaces
│   └── tools/               # Software tools backgrounds
│       ├── tools1.jpg       # App interfaces, dashboards
│       ├── tools2.jpg       # Software screenshots
│       └── tools3.jpg       # Digital tools
└── music/                   # Background music
    ├── background.mp3       # Main background track
    └── background.wav       # Alternative format
```

## 🖼️ Image Specifications

### Logos
- **Format**: PNG (transparent) or JPG
- **Size**: 200x200px to 500x500px
- **Aspect**: Square preferred
- **Style**: Clean, readable at small sizes

### Background Images  
- **Format**: JPG or PNG
- **Size**: 1080x1920px (9:16 aspect ratio)
- **Quality**: High resolution for zoom effects
- **Style**: Tech/AI themed, not too busy (text overlay friendly)

### Banners
- **Format**: PNG preferred
- **Size**: 1920x1080px or 1080x1920px
- **Use**: Channel branding, intro/outro graphics

## 🎵 Audio Specifications

### Background Music
- **Format**: MP3 or WAV
- **Length**: 30+ seconds (will loop)
- **Style**: Upbeat, tech-focused, royalty-free
- **Volume**: Medium (will be auto-adjusted)

## 🚀 How It Works

1. **Script Analysis**: Convertbox analyzes your script for keywords
2. **Smart Matching**: 
   - "AI", "artificial intelligence" → `assets/images/ai/`
   - "technology", "digital", "code" → `assets/images/tech/`  
   - "tools", "software", "app" → `assets/images/tools/`
3. **Dynamic Backgrounds**: Changes images 2-3 times per video based on content
4. **Logo Overlay**: Automatically finds and applies your logo
5. **Music Mix**: Blends background music with narration

## 📤 Upload Your Assets

1. **Add your logos** to `assets/logos/`
2. **Add background images** to appropriate category folders
3. **Add background music** to `assets/music/`
4. **Run Convertbox** - it will automatically use your assets!

## 💡 Pro Tips

- Use **high contrast** images for better text readability
- Keep logos **simple** - they appear small in videos
- Choose **royalty-free** music to avoid copyright issues
- **Test different images** to see what works best for your content
- Images with **subtle motion blur** work great with zoom effects

## 🎬 Example Usage

```bash
# After adding your assets:
go run ./cmd/convertbox --topic "5 AI Tools That Will Change Everything"

# Convertbox will:
# ✅ Use AI-themed backgrounds for AI content
# ✅ Switch to tech backgrounds for technical parts  
# ✅ Apply your channel logo
# ✅ Mix your background music
# ✅ Create professional YouTube Short ready for upload!
```

Your AI Unboxed channel will have unique, branded videos every time! 🚀