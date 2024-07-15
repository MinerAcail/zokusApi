package redis

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/nfnt/resize"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func ConfigureRedis() *redis.Client {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: "", // No password set
		DB:       0,  // Use default DB
	})

	// Test the connection
	fmt.Println("Attempting to connect to Redis...")
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Could			 not connect to Redis: %v", err)
	} else {
		fmt.Println("Successfully connected to Redis!")
	}
	return rdb
}

// ImageToByteSlice converts an image file to a byte slice
func ImageToByteSlice(file io.Reader) ([]byte, error) {
	img, format, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("could not decode image: %v", err)
	}

	var buf bytes.Buffer
	switch format {
	case "jpeg":
		if err := jpeg.Encode(&buf, img, nil); err != nil {
			return nil, fmt.Errorf("could not encode image to JPEG: %v", err)
		}
	case "png":
		if err := png.Encode(&buf, img); err != nil {
			return nil, fmt.Errorf("could not encode image to PNG: %v", err)
		}
	default:
		return nil, fmt.Errorf("unsupported image format: %s", format)
	}

	return buf.Bytes(), nil
}

// ByteSliceToImage converts a byte slice to an image.Image
func ByteSliceToImage(data []byte) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("could not decode byte slice to image: %v", err)
	}
	return img, nil
}

// PassImageToFfmpeg passes an image.Image to FFmpeg for processing

// 	img = resize.Resize(1070, 0, img, resize.Lanczos3)
func PassImageToFfmpeg(img image.Image) error {
	// Resize the image to ensure dimensions are even
	img = resize.Resize(1070, 0, img, resize.Lanczos3)

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return fmt.Errorf("could not encode image to PNG: %v", err)
	}

	// Define the output directory and path
	outputDir := "./output_videos"
	outputPath := filepath.Join(outputDir, "slideshow.mp4")

	// Remove existing output directory if it exists
	if _, err := os.Stat(outputDir); !os.IsNotExist(err) {
		if err := os.RemoveAll(outputDir); err != nil {
			return fmt.Errorf("could not remove existing output directory: %v", err)
		}
	}

	// Create the output directory
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("could not create output directory: %v", err)
	}

	// FFmpeg command to convert image to video
	ffmpegCmd := exec.Command("ffmpeg",
		"-framerate", "1/5",           // 1 image every 5 seconds
		"-f", "image2pipe",            // Input from image2pipe
		"-i", "pipe:0",                // Read input from stdin
		"-vf", "scale=972:768",        // Scale the video
		"-c:v", "libx264",             // Encoder
		"-pix_fmt", "yuv420p",         // Pixel format
		outputPath)                    // Output path

	ffmpegCmd.Stdin = &buf // Set stdin to the image data

	// Run the command
	output, err := ffmpegCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg error: %v, output: %s", err, output)
	}

	return nil
}


// PassByteSliceToRedis stores a byte slice in Redis
func PassByteSliceToRedis(rdb *redis.Client, key string, data []byte) error {
	err := rdb.Set(ctx, key, data, 0).Err()
	if err != nil {
		return fmt.Errorf("could not set byte slice in Redis: %v", err)
	}
	return nil
}

// RetrieveByteSliceFromRedis retrieves a byte slice from Redis
func RetrieveByteSliceFromRedis(rdb *redis.Client, key string) ([]byte, error) {
	data, err := rdb.Get(ctx, key).Bytes()
	if err != nil {
		return nil, fmt.Errorf("could not get byte slice from Redis: %v", err)
	}
	return data, nil
}



// 	var buf bytes.Buffer
// 	if err := png.Encode(&buf, img); err != nil {
// 		return fmt.Errorf("could not encode image to PNG: %v", err)
// 	}

// 	// Define the output directory and path
// 	outputDir := "./output_videos"
// 	outputPath := filepath.Join(outputDir, "output.mp4")

// 	// Remove existing output directory if it exists
// 	if _, err := os.Stat(outputDir); !os.IsNotExist(err) {
// 		if err := os.RemoveAll(outputDir); err != nil {
// 			return fmt.Errorf("could not remove existing output directory: %v", err)
// 		}
// 	}

// 	// Create the output directory
// 	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
// 		return fmt.Errorf("could not create output directory: %v", err)
// 	}

// 	// FFmpeg command to convert image to video
// 	ffmpegCmd := exec.Command("ffmpeg",
// 		"-loop", "1", // Loop the input indefinitely (for the image)
// 		"-i", "-", // Input from stdin (the image)
// 		"-t", "10", // Duration of the video (10 seconds)
// 		"-s", "1070x1070", // Output video size (even dimensions)
// 		"-pix_fmt", "yuv420p", // Pixel format for compatibility
// 		"-vf", "fps=25", // Output frame rate (adjust as needed)
// 		"-vcodec", "libx264", // Video codec
// 		"-preset", "medium", // Encoding preset
// 		outputPath) // Output file path

// 	ffmpegCmd.Stdin = &buf // Set stdin to the image data

// 	// Run the command
// 	output, err := ffmpegCmd.CombinedOutput()
// 	if err != nil {
// 		return fmt.Errorf("ffmpeg error: %v, output: %s", err, output)
// 	}

// 	return nil
// }