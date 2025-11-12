#!/bin/bash

# Download free background music and images for Convertbox

echo "📥 Downloading free assets for Convertbox..."

# Create directories
mkdir -p assets/music assets/images

# Download royalty-free background music (Creative Commons)
echo "🎵 Downloading background music..."
curl -L "https://www.soundjay.com/misc/beep-07a.wav" -o assets/music/background.wav 2>/dev/null || echo "⚠️  Music download failed - add your own to assets/music/"

# Download free background images from Unsplash (AI/Tech themed)
echo "🖼️  Downloading background images..."
curl -L "https://images.unsplash.com/photo-1518709268805-4e9042af2176?w=1080&h=1920&fit=crop" -o assets/images/tech1.jpg 2>/dev/null || echo "⚠️  Image download failed"
curl -L "https://images.unsplash.com/photo-1555949963-aa79dcee981c?w=1080&h=1920&fit=crop" -o assets/images/tech2.jpg 2>/dev/null || echo "⚠️  Image download failed"
curl -L "https://images.unsplash.com/photo-1485827404703-89b55fcc595e?w=1080&h=1920&fit=crop" -o assets/images/tech3.jpg 2>/dev/null || echo "⚠️  Image download failed"

echo "✅ Asset download complete!"
echo "📁 Check assets/music/ and assets/images/ folders"