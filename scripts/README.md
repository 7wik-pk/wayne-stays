# Image Optimization Scripts

This directory contains scripts for optimizing images for the Wayne Stays platform.

## optimize-images.sh

Comprehensive script that:
- ✅ Converts HEIC files to JPG
- ✅ Resizes images to max 2400px width (maintains aspect ratio)
- ✅ Optimizes JPEG quality to 85%
- ✅ Strips metadata (reduces file size)
- ✅ Creates backups before processing
- ✅ Provides detailed progress output

### Usage

```bash
cd scripts
./optimize-images.sh
```

### Configuration

Edit these variables at the top of the script:

```bash
MAX_WIDTH=2400        # Max width for images
QUALITY=85            # JPEG quality (0-100)
```

### What it does:

1. **Backup**: Creates `img_backup/` with copies of all originals
2. **Convert HEIC**: Finds all `.heic` files and converts to `.jpg`
3. **Resize**: Resizes images wider than 2400px
4. **Optimize**: Strips metadata and optimizes compression
5. **Clean up**: Removes original HEIC files after conversion

### Example Output:

```
🖼️  Wayne Stays Image Optimization
==================================
✅ Using: magick
📦 Creating backup directory...
✅ Backup created at: ../img_backup

🔍 Scanning for images...

📸 Converting HEIC files to JPG...
   Converting: ./arkham_district_gentlemens_club/587092492_....heic → 587092492_....jpg
   ✅ Converted and removed HEIC

🔧 Resizing and optimizing JPG images...
   Resizing: ./wayne_tower/1.jpg (3840px → 2400px)
   ✅ Resized and optimized

==================================
✅ Image Optimization Complete!
==================================
📊 Statistics:
   Total files processed: 31
   HEIC → JPG conversions: 4
   Images resized: 12
   Skipped files: 0
```

## Quick Commands

### Convert only HEIC to JPG (no resize):

```bash
cd img/accommodations
find . -name "*.heic" -exec sh -c 'magick "$1" "${1%.heic}.jpg"' _ {} \;
```

### Resize all images to 2000px width:

```bash
cd img/accommodations
find . -name "*.jpg" -exec magick {} -resize 2000x\> {} \;
```

### Check image dimensions:

```bash
cd img/accommodations
find . -name "*.jpg" -exec identify -format "%f: %wx%h\n" {} \;
```

## Recommended Sizes

| Image Type | Recommended Width | Use Case |
|------------|------------------|----------|
| Hero Images | 1920-2400px | Full-width displays, retina screens |
| Gallery Images | 1200-1600px | Lightbox/modal views |
| Thumbnails | 400-800px | Grid/card views |
| Mobile | Auto-responsive | Let browser scale down from hero |

## Web Optimization Tips

1. **Format**: JPEG for photos, PNG for graphics with transparency
2. **Quality**: 80-85% is sweet spot (good quality, reasonable size)
3. **Dimensions**: Max 2400px wide covers 4K displays
4. **Metadata**: Strip EXIF data to save ~10-20KB per image
5. **Progressive**: Use progressive JPEG for better perceived loading

## Troubleshooting

### ImageMagick not found (Arch Linux):
```bash
sudo pacman -S imagemagick
```

### HEIC support issues (Arch Linux):
```bash
sudo pacman -S libheif
```

### Restore from backup:
```bash
rm -rf img/accommodations
cp -r img_backup img/accommodations
```
