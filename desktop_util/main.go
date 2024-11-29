package main

import (
	"image"
	"image/color"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ncruces/zenity"
)

var percent float32 = 0.055
var max float32 = 65535 * percent
var maxInt16 uint16 = uint16(max)

func main() {
	filePaths, err := zenity.SelectFileMultiple(
		zenity.Title("Select Image Files"),
		zenity.FileFilter{
			Name: "Image Files",
			Patterns: []string{
				"*.jpg",
				"*.jpeg",
				"*.png",
			},
		},
	)
	if err != nil {
		log.Fatalf("Failed to select files: %v", err)
	}

	for _, filePath := range filePaths {
		fileName := filepath.Base(filePath)

		if isImage(fileName) {
			processImage(filePath)
		} else {
			log.Printf("The selected file is not suitable for processing: %s", filePath)
		}
	}
}

func isImage(fileName string) bool {
	lower := strings.ToLower(fileName)
	return strings.HasSuffix(lower, ".jpg") || strings.HasSuffix(lower, ".jpeg") || strings.HasSuffix(lower, ".png")
}

func processImage(inputFile string) {
	outputFile := strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + "_готово_к_печати" + filepath.Ext(inputFile)

	file, err := os.Open(inputFile)
	if err != nil {
		log.Printf("Failed to open input file %s: %v", inputFile, err)
		return
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		log.Printf("Failed to decode image %s: %v", inputFile, err)
		return
	}

	newImg := replaceBlackPixels(img)

	outFile, err := os.Create(outputFile)
	if err != nil {
		log.Printf("Failed to create output file %s: %v", outputFile, err)
		return
	}
	defer outFile.Close()

	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		err = jpeg.Encode(outFile, newImg, &jpeg.Options{Quality: 100})
	case "png":
		err = png.Encode(outFile, newImg)
	default:
		log.Printf("Unsupported image format %s: %s", inputFile, format)
		return
	}

	if err != nil {
		log.Printf("Failed to save output file %s: %v", outputFile, err)
		return
	}

	log.Printf("Image processing completed for %s, saved to %s", inputFile, outputFile)
}

func replaceBlackPixels(img image.Image) image.Image {
	bounds := img.Bounds()
	newImg := image.NewRGBA64(bounds)

	lessBlack := color.RGBA64{R: maxInt16, G: maxInt16, B: maxInt16, A: 65535}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			if uint16(r) <= maxInt16 && uint16(g) <= maxInt16 && uint16(b) <= maxInt16 {
				newImg.SetRGBA64(x, y, lessBlack)
			} else {
				newImg.SetRGBA64(x, y, color.RGBA64{
					R: uint16(r),
					G: uint16(g),
					B: uint16(b),
					A: uint16(a),
				})
			}
		}
	}

	return newImg
}
