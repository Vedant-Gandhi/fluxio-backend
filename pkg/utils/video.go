package utils

import (
	"fmt"
	"math/rand"
	"slices"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

var (

	// Declare it here to reduce creation costs for the GC.
	validVidMimes = []string{
		"video/mp4",        // MP4 format
		"video/webm",       // WebM format
		"video/ogg",        // Ogg video
		"video/3gpp",       // 3GPP mobile video
		"video/quicktime",  // MOV format
		"video/x-msvideo",  // AVI format
		"video/mpeg",       // MPEG format
		"video/x-matroska", // MKV format
		"video/mp2t",       // MPEG Transport Stream (.ts files)
	}
)

func CreateURLSafeVideoSlug(title string) (slug string) {
	// 1. Convert to lowercase and trim spaces
	slug = strings.ToLower(strings.TrimSpace(title))

	// 2. Normalize and transform Unicode characters
	// First normalize using NFD to separate characters and combining marks
	t := norm.NFD.String(slug)

	// 3. Remove non-alphanumeric characters and convert spaces to hyphens
	var result strings.Builder
	var lastWasHyphen bool

	for _, r := range t {
		switch {
		case unicode.IsLetter(r) || unicode.IsNumber(r):
			// Keep letters and numbers
			result.WriteRune(r)
			lastWasHyphen = false
		case unicode.IsSpace(r) || unicode.IsPunct(r) || unicode.IsSymbol(r):
			// Replace spaces, punctuation, and symbols with hyphens
			if !lastWasHyphen {
				result.WriteRune('-')
				lastWasHyphen = true
			}
			// Skip all other characters including combining marks (accents)
		}
	}

	// 4. Convert back to clean string and remove any leading/trailing hyphens
	slug = strings.Trim(result.String(), "-")

	// 5. Limit length
	maxLength := 50
	if len(slug) > maxLength {
		// Cut at word boundary if possible
		if idx := strings.LastIndex(slug[:maxLength], "-"); idx > 0 {
			slug = slug[:idx]
		} else {
			slug = slug[:maxLength]
		}
	}

	// 6. Create a local random source
	source := rand.NewSource(time.Now().UnixNano())
	localRand := rand.New(source)
	randomNum := localRand.Intn(10000)

	// 7. Generate unique suffix using nanosecond timestamp
	timestamp := time.Now().UnixNano()
	uniqueSuffix := fmt.Sprintf("%x%04d", timestamp%0x10000, randomNum)

	// 8. Build the final slug
	if len(slug) > 0 {
		slug = slug + "-" + uniqueSuffix
		return
	}

	// Fallback if title was empty or contained only special chars
	slug = "video-" + uniqueSuffix
	return
}

func CreateURLSafeThumbnailFileName(videoID string, timestamp string) (fileName string) {
	// 1. Convert to lowercase and trim spaces
	fileName = strings.ToLower(fmt.Sprintf("%s-%s", strings.TrimSpace(videoID), strings.TrimSpace(timestamp)))

	// 2. Normalize and transform Unicode characters
	// First normalize using NFD to separate characters and combining marks
	t := norm.NFD.String(fileName)

	// 3. Remove non-alphanumeric characters and convert spaces to hyphens
	var result strings.Builder
	var lastWasHyphen bool

	for _, r := range t {
		switch {
		case unicode.IsLetter(r) || unicode.IsNumber(r):
			// Keep letters and numbers
			result.WriteRune(r)
			lastWasHyphen = false
		case unicode.IsSpace(r) || unicode.IsPunct(r) || unicode.IsSymbol(r):
			// Replace spaces, punctuation, and symbols with hyphens
			if !lastWasHyphen {
				result.WriteRune('-')
				lastWasHyphen = true
			}
			// Skip all other characters including combining marks (accents)
		}
	}

	// 4. Convert back to clean string and remove any leading/trailing hyphens
	fileName = strings.Trim(result.String(), "-")

	// 5. Limit length
	maxLength := 50
	if len(fileName) > maxLength {
		// Cut at word boundary if possible
		if idx := strings.LastIndex(fileName[:maxLength], "-"); idx > 0 {
			fileName = fileName[:idx]
		} else {
			fileName = fileName[:maxLength]
		}
	}

	return "thumbnail-" + fileName
}

func CheckVideoMimeTypeValidity(mimeType string) (valid bool) {
	valid = slices.Contains(validVidMimes, mimeType)
	return
}
