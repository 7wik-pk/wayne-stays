#!/bin/bash

# Image Optimization Script for Wayne Stays
# Converts HEIC to JPG and resizes all images for web

set -e  # Exit on error

# Configuration (paths are relative to the repo root, wherever the script is run from)
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
IMG_DIR="$REPO_ROOT/img/accommodations"
BACKUP_DIR="$REPO_ROOT/img_backup"
MAX_WIDTH=2400        # Max width for hero images
THUMBNAIL_WIDTH=800   # Width for thumbnails
QUALITY=85            # JPEG quality (0-100)

echo "🖼️  Wayne Stays Image Optimization"
echo "=================================="

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Check if ImageMagick is installed
if ! command -v convert &> /dev/null && ! command -v magick &> /dev/null; then
    echo -e "${RED}❌ ImageMagick not found. Please install it:${NC}"
    echo "   sudo pacman -S imagemagick"
    exit 1
fi

# Use magick if available (ImageMagick 7), otherwise the v6 convert/identify pair
if command -v magick &> /dev/null; then
    CONVERT_CMD="magick"
    IDENTIFY_CMD="magick identify"
else
    CONVERT_CMD="convert"
    IDENTIFY_CMD="identify"
fi

echo -e "${GREEN}✅ Using: $CONVERT_CMD${NC}"

# Create backup directory
if [ ! -d "$BACKUP_DIR" ]; then
    echo -e "${YELLOW}📦 Creating backup directory...${NC}"
    cp -r "$IMG_DIR" "$BACKUP_DIR"
    echo -e "${GREEN}✅ Backup created at: $BACKUP_DIR${NC}"
fi

# Navigate to img directory
cd "$IMG_DIR"

# Statistics
total_files=0
converted_heic=0
resized_images=0
skipped_files=0

echo ""
echo "🔍 Scanning for images..."

# Process HEIC files first
echo ""
echo -e "${YELLOW}📸 Converting HEIC files to JPG...${NC}"
# NUL-separated find output keeps filenames with spaces intact, and process
# substitution (instead of a pipe) keeps the counters in this shell
while IFS= read -r -d '' heic_file; do
    total_files=$((total_files + 1))
    jpg_file="${heic_file%.*}.jpg" # strips .heic or .HEIC

    echo "   Converting: $heic_file → ${jpg_file##*/}"

    # Tested in the if, so a failure is reported instead of tripping set -e
    if $CONVERT_CMD "$heic_file" -quality $QUALITY "$jpg_file"; then
        # Remove original HEIC file after successful conversion
        rm "$heic_file"
        converted_heic=$((converted_heic + 1))
        echo -e "   ${GREEN}✅ Converted and removed HEIC${NC}"
    else
        echo -e "   ${RED}❌ Failed to convert${NC}"
    fi
done < <(find . -type f -iname "*.heic" -print0)

echo ""
echo -e "${YELLOW}🔧 Resizing and optimizing JPG images...${NC}"

# Process all JPG/JPEG files
while IFS= read -r -d '' img_file; do
    total_files=$((total_files + 1))

    # Get current dimensions
    width=$($IDENTIFY_CMD -format "%w" "$img_file" 2>/dev/null || true)
    
    if [ -z "$width" ]; then
        echo "   ⚠️  Skipping invalid file: $img_file"
        skipped_files=$((skipped_files + 1))
        continue
    fi
    
    # Only resize if width exceeds MAX_WIDTH
    if [ "$width" -gt "$MAX_WIDTH" ]; then
        echo "   Resizing: $img_file (${width}px → ${MAX_WIDTH}px)"
        
        # Create temporary file
        temp_file="${img_file}.tmp.jpg"
        
        # Resize image
        if $CONVERT_CMD "$img_file" \
            -resize "${MAX_WIDTH}x>" \
            -quality $QUALITY \
            -strip \
            "$temp_file"; then
            mv "$temp_file" "$img_file"
            resized_images=$((resized_images + 1))
            echo -e "   ${GREEN}✅ Resized and optimized${NC}"
        else
            echo -e "   ${RED}❌ Failed to resize${NC}"
            rm -f "$temp_file"
        fi
    else
        # Still optimize even if not resizing
        echo "   Optimizing: $img_file (${width}px - no resize needed)"
        
        temp_file="${img_file}.tmp.jpg"
        
        if $CONVERT_CMD "$img_file" \
            -quality $QUALITY \
            -strip \
            "$temp_file"; then
            mv "$temp_file" "$img_file"
            echo -e "   ${GREEN}✅ Optimized${NC}"
        else
            rm -f "$temp_file"
        fi
    fi
done < <(find . -type f \( -iname "*.jpg" -o -iname "*.jpeg" \) -print0)

# Summary
echo ""
echo "=================================="
echo -e "${GREEN}✅ Image Optimization Complete!${NC}"
echo "=================================="
echo "📊 Statistics:"
echo "   Total files processed: $total_files"
echo "   HEIC → JPG conversions: $converted_heic"
echo "   Images resized: $resized_images"
echo "   Skipped files: $skipped_files"
echo ""
echo "📁 Original images backed up to: $BACKUP_DIR"
echo ""
echo -e "${YELLOW}💡 Tip: You can delete the backup with:${NC}"
echo "   rm -rf $BACKUP_DIR"
