package redis

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
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

// img = resize.Resize(1070, 0, img, resize.Lanczos3)
func PassImageToFfmpegg(img image.Image) error {
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
		"-framerate", "1/5", // 1 image every 5 seconds
		"-f", "image2pipe", // Input from image2pipe
		"-i", "pipe:0", // Read input from stdin
		"-t", "5", // 5 seconds duration
		"-vf", "scale=-1:1920", // Scale the video
		"-c:v", "libx264", // Video codec
		"-pix_fmt", "yuv420p", // Pixel format
		outputPath) // Output file path

	ffmpegCmd.Stdin = &buf // Set stdin to the image data

	// Run the command
	output, err := ffmpegCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg error: %v, output: %s", err, output)
	}

	return nil
}

func PassImageToFfmpeg(img image.Image) (string, error) {
	// Resize the image to ensure dimensions are even
	img = resize.Resize(1070, 0, img, resize.Lanczos3)

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", fmt.Errorf("could not encode image to PNG: %v", err)
	}

	// Define the output directory and path
	outputDir := "./output_videos"
	outputPath := filepath.Join(outputDir, "slideshow.mp4")

	// Remove existing output directory if it exists
	if _, err := os.Stat(outputDir); !os.IsNotExist(err) {
		if err := os.RemoveAll(outputDir); err != nil {
			return "", fmt.Errorf("could not remove existing output directory: %v", err)
		}
	}

	// Create the output directory
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("could not create output directory: %v", err)
	}

	// FFmpeg command to convert image to video
	ffmpegCmd := exec.Command("ffmpeg",
		"-framerate", "1/5", // 1 image every 5 seconds
		"-f", "image2pipe", // Input from image2pipe
		"-i", "pipe:0", // Read input from stdin
		"-t", "5", // 5 seconds duration
		"-vf", "scale=-1:1920", // Scale the video
		"-c:v", "libx264", // Video codec
		"-pix_fmt", "yuv420p", // Pixel format
		outputPath) // Output file path

	ffmpegCmd.Stdin = &buf // Set stdin to the image data

	// Run the command
	output, err := ffmpegCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ffmpeg error: %v, output: %s", err, output)
	}

	return outputPath, nil
}
func PassImagesToFfmpeg1(imageKeys []string, rdb *redis.Client) (string, error) {
	// Define the output directory and path
	outputDir := "./output_videos"
	outputPath := filepath.Join(outputDir, "slideshow.mp4")

	// Remove existing output directory if it exists
	if _, err := os.Stat(outputDir); !os.IsNotExist(err) {
		if err := os.RemoveAll(outputDir); err != nil {
			return "", fmt.Errorf("could not remove existing output directory: %v", err)
		}
	}

	// Create the output directory
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("could not create output directory: %v", err)
	}

	// Create a buffer to hold image data
	var imgBuffer bytes.Buffer

	// Retrieve images from Redis
	for _, imageKey := range imageKeys {
		imgData, err := rdb.Get(ctx, imageKey).Bytes()
		if err != nil {
			return "", fmt.Errorf("error retrieving image from Redis: %v", err)
		}
		imgBuffer.Write(imgData)
	}

	// FFmpeg command to convert images to video
	ffmpegCmd := exec.Command("ffmpeg",
		"-framerate", "1/5", // 1 image every 5 seconds
		"-f", "image2pipe", // Input from image2pipe
		"-i", "pipe:0", // Read input from stdin
		"-vf", "scale=-1:1920", // Scale the video
		"-c:v", "libx264", // Video codec
		"-pix_fmt", "yuv420p", // Pixel format
		outputPath) // Output file path

	ffmpegCmd.Stdin = &imgBuffer // Set stdin to the image data

	// Run the command
	output, err := ffmpegCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ffmpeg error: %v, output: %s", err, output)
	}

	return outputPath, nil
}
func PassImagesToFfmpeg(imageKeys []string, rdb *redis.Client, scale string) (string, error) {
	// Define the output directory and path
	outputDir := "./output_videos"
	outputPath := filepath.Join(outputDir, "slideshow.mp4")

	// Remove existing output directory if it exists
	if _, err := os.Stat(outputDir); !os.IsNotExist(err) {
		if err := os.RemoveAll(outputDir); err != nil {
			return "", fmt.Errorf("could not remove existing output directory: %v", err)
		}
	}

	// Create the output directory
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("could not create output directory: %v", err)
	}

	// Create a buffer to hold image data
	var imgBuffer bytes.Buffer

	// Retrieve images from Redis
	for _, imageKey := range imageKeys {
		imgData, err := rdb.Get(ctx, imageKey).Bytes()
		if err != nil {
			return "", fmt.Errorf("error retrieving image from Redis: %v", err)
		}
		imgBuffer.Write(imgData)
	}

	// FFmpeg command to convert images to video
	ffmpegCmd := exec.Command("ffmpeg",
		"-framerate", "1/5", // 1 image every 5 seconds
		"-f", "image2pipe", // Input from image2pipe
		"-i", "pipe:0", // Read input from stdin
		// "-vf", "scale=trunc(iw/2)*2:trunc(ih/2)*2", // Ensure dimensions are divisible by 2
		"-vf", scale,
		"-c:v", "libx264", // Video codec
		"-pix_fmt", "yuv420p", // Pixel format
		outputPath) // Output file path

	ffmpegCmd.Stdin = &imgBuffer // Set stdin to the image data

	// Run the command
	output, err := ffmpegCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ffmpeg error: %v, output: %s", err, output)
	}

	return outputPath, nil
}

func processImagesToVideo(imageKeys []string, rdb *redis.Client, scale string, format string, movflags string) ([]byte, error) {
	var imgBuffer bytes.Buffer

	// Retrieve image data from Redis and write to imgBuffer
	for _, imageKey := range imageKeys {
		imgData, err := rdb.Get(ctx, imageKey).Bytes()
		if err != nil {
			return nil, fmt.Errorf("error retrieving image from Redis: %v", err)
		}
		imgBuffer.Write(imgData)
	}

	// Create MKV video
	var mkvBuffer bytes.Buffer
	var ffmpegStderr bytes.Buffer

	ffmpegCmd := exec.Command("ffmpeg",
		"-framerate", "1/5",
		"-f", "image2pipe",
		"-i", "pipe:0",
		"-vf", scale,
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		"-f", "matroska", // Use matroska (MKV) format
		"pipe:1")

	ffmpegCmd.Stdin = &imgBuffer
	ffmpegCmd.Stdout = &mkvBuffer
	ffmpegCmd.Stderr = &ffmpegStderr

	if err := ffmpegCmd.Run(); err != nil {
		return nil, fmt.Errorf("ffmpeg error during MKV generation: %v, stderr: %v", err, ffmpegStderr.String())
	}

	// Convert MKV to MP4
	var mp4Buffer bytes.Buffer
	ffmpegStderr.Reset()

	ffmpegCmd = exec.Command("ffmpeg",
		"-i", "pipe:0",
		"-vf", scale,
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		"-movflags", movflags,
		"-f", format,
		"pipe:1")

	ffmpegCmd.Stdin = &mkvBuffer
	ffmpegCmd.Stdout = &mp4Buffer
	ffmpegCmd.Stderr = &ffmpegStderr

	if err := ffmpegCmd.Run(); err != nil {
		return nil, fmt.Errorf("ffmpeg error during MP4 conversion: %v, stderr: %v", err, ffmpegStderr.String())
	}

	return mp4Buffer.Bytes(), nil
}

func PassImagesToFfmpeg_OVERLAY_Bytes(imageKeys []string, rdb *redis.Client, scale string) ([]byte, error) {
	return processImagesToVideo(imageKeys, rdb, scale, "mp4", "frag_keyframe+empty_moov")
}

func PassImagesToFfmpeg_BG_Bytes(imageKeys []string, rdb *redis.Client, scale string) ([]byte, error) {
	return processImagesToVideo(imageKeys, rdb, scale, "mp4", "frag_keyframe+empty_moov")
}

func CombineVideosFromMemory(bgData, overlayData []byte, outputFormat string, overlayX, overlayY int) ([]byte, error) {
	var outputBuffer bytes.Buffer
	var ffmpegStderr bytes.Buffer

	// Generate UUID for the FIFO name
	overlayFIFOName := fmt.Sprintf("/tmp/overlay_video_%s.fifo", uuid.New().String())

	// Create the FIFO for the overlay video
	if err := syscall.Mkfifo(overlayFIFOName, 0666); err != nil {
		return nil, fmt.Errorf("error creating FIFO for overlay video: %w", err)
	}
	defer os.Remove(overlayFIFOName) // Clean up the FIFO

	// Write overlay data to the FIFO in a separate goroutine
	go func() {
		file, err := os.OpenFile(overlayFIFOName, os.O_WRONLY, 0600)
		if err != nil {
			fmt.Printf("error opening FIFO for writing: %v\n", err)
			return
		}
		defer file.Close()

		if _, err := file.Write(overlayData); err != nil {
			fmt.Printf("error writing to FIFO: %v\n", err)
		}
	}()

	// Create a pipe for background video data
	bgReader, bgWriter := io.Pipe()

	// FFmpeg command to combine videos and output to MKV format
	ffmpegCmd := exec.Command("ffmpeg",
		"-i", "pipe:0", // Background video from stdin
		"-i", overlayFIFOName, // Overlay video from the FIFO
		"-filter_complex", fmt.Sprintf("overlay=%d:%d", overlayX, overlayY),
		"-preset", "fast",
		"-f", "matroska", // Use Matroska format for combining
		"pipe:1") // Output to stdout

	// Set the pipes for FFmpeg
	ffmpegCmd.Stdin = bgReader
	ffmpegCmd.Stdout = &outputBuffer
	ffmpegCmd.Stderr = &ffmpegStderr

	// Start FFmpeg command
	if err := ffmpegCmd.Start(); err != nil {
		return nil, fmt.Errorf("ffmpeg error starting command: %w", err)
	}

	// Write the background data to the pipe in a separate goroutine
	go func() {
		defer bgWriter.Close()
		if _, err := bgWriter.Write(bgData); err != nil {
			fmt.Printf("error writing background data to pipe: %v\n", err)
		}
	}()

	// Wait for FFmpeg command to finish
	if err := ffmpegCmd.Wait(); err != nil {
		return nil, fmt.Errorf("ffmpeg error during video combination: %w, stderr: %v", err, ffmpegStderr.String())
	}

	// Convert MKV to MP4
	var mp4Buffer bytes.Buffer
	ffmpegStderr.Reset()

	ffmpegCmd = exec.Command("ffmpeg",
		"-i", "pipe:0", // Input from the combined MKV video
		"-c:v", "libx264", // Video codec
		"-pix_fmt", "yuv420p", // Pixel format
		"-movflags", "frag_keyframe+empty_moov",
		"-f", "mp4", // Output format
		"pipe:1") // Output to stdout

	ffmpegCmd.Stdin = &outputBuffer
	ffmpegCmd.Stdout = &mp4Buffer
	ffmpegCmd.Stderr = &ffmpegStderr

	if err := ffmpegCmd.Run(); err != nil {
		return nil, fmt.Errorf("ffmpeg error during MP4 conversion: %v, stderr: %v", err, ffmpegStderr.String())
	}

	return mp4Buffer.Bytes(), nil
}

func PassImagessToFfmpeg(images []image.Image) (string, error) {
	// Define the output directory and path
	outputDir := "./output_videos"
	outputPath := filepath.Join(outputDir, "slideshow.mp4")

	// Remove existing output directory if it exists
	if _, err := os.Stat(outputDir); !os.IsNotExist(err) {
		if err := os.RemoveAll(outputDir); err != nil {
			return "", fmt.Errorf("could not remove existing output directory: %v", err)
		}
	}

	// Create the output directory
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("could not create output directory: %v", err)
	}

	// Write images to a temporary file
	tempFile, err := os.CreateTemp("", "images-*.txt")
	if err != nil {
		return "", fmt.Errorf("could not create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	for i, img := range images {
		imgPath := fmt.Sprintf("%s/image_%d.png", outputDir, i)
		imgFile, err := os.Create(imgPath)
		if err != nil {
			return "", fmt.Errorf("could not create image file: %v", err)
		}
		defer imgFile.Close()

		if err := png.Encode(imgFile, img); err != nil {
			return "", fmt.Errorf("could not encode image to PNG: %v", err)
		}

		if _, err := tempFile.WriteString(fmt.Sprintf("file '%s'\n", imgPath)); err != nil {
			return "", fmt.Errorf("could not write to temporary file: %v", err)
		}
	}

	// FFmpeg command to convert images to video
	ffmpegCmd := exec.Command("ffmpeg",
		"-framerate", "1/5", // 1 image every 5 seconds
		"-f", "concat", // Concatenate images
		"-safe", "0",
		"-i", tempFile.Name(), // Read input from the temporary file
		"-vf", "scale=-1:1920", // Scale the video
		"-c:v", "libx264", // Video codec
		"-pix_fmt", "yuv420p", // Pixel format
		outputPath) // Output file path

	// Run the command
	output, err := ffmpegCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ffmpeg error: %v, output: %s", err, output)
	}

	return outputPath, nil
}

func PassImagesInToFfmpegAsOverlay(files []image.Image) (string, error) {
	// Create input directory if it doesn't exist
	inputDir := "./input"
	if err := os.MkdirAll(inputDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("could not create input directory: %v", err)
	}

	// Write images to input directory as JPEG files
	for i, img := range files {
		// Resize the image to ensure dimensions are even
		img = resize.Resize(972, 768, img, resize.Lanczos3)

		// Encode image to JPEG format
		var buf bytes.Buffer
		if err := jpeg.Encode(&buf, img, nil); err != nil {
			return "", fmt.Errorf("could not encode image to JPEG: %v", err)
		}

		// Write image data to file
		imgPath := filepath.Join(inputDir, strconv.Itoa(i)+".jpg")
		if err := os.WriteFile(imgPath, buf.Bytes(), 0644); err != nil {
			return "", fmt.Errorf("could not write image to file: %v", err)
		}

		fmt.Printf("Uploaded file %d of %d\n", i+1, len(files))
	}

	// Define output directory and path
	outputDir := "./output"
	outputPath := filepath.Join(outputDir, "slideshow.mp4")

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("could not create output directory: %v", err)
	}

	// FFmpeg command to create video from JPEG images
	ffmpegCmd := exec.Command("ffmpeg",
		"-framerate", "1/5", // 1 image every 5 seconds
		"-pattern_type", "glob", // Use glob pattern for input files
		"-i", filepath.Join(inputDir, "*.jpg"), // Input files pattern
		"-vf", "scale=972:768", // Scale the video
		"-c:v", "libx264", // Video codec
		"-pix_fmt", "yuv420p", // Pixel format
		outputPath) // Output file path

	// Run FFmpeg command
	output, err := ffmpegCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ffmpeg error: %v, output: %s", err, output)
	}

	fmt.Println("Video created successfully at:", outputPath)

	return outputPath, nil
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

func ProcessImages(files []*graphql.Upload) ([][]byte, error) {
	var retrievedByteSlices [][]byte

	for _, file := range files {
		// Convert image to byte slice
		byteSlice, err := ImageToByteSlice(file.File)
		if err != nil {
			return nil, fmt.Errorf("error converting image to byte slice: %w", err)
		}

		// Configure Redis
		rdb := ConfigureRedis()

		// Generate a unique key for each image
		imageKey := fmt.Sprintf("image_key_%d", len(retrievedByteSlices))

		// Store byte slice in Redis
		if err := PassByteSliceToRedis(rdb, imageKey, byteSlice); err != nil {
			return nil, fmt.Errorf("error storing byte slice in Redis: %w", err)
		}

		// Retrieve byte slice from Redis
		retrievedByteSlice, err := RetrieveByteSliceFromRedis(rdb, imageKey)
		if err != nil {
			return nil, fmt.Errorf("error retrieving byte slice from Redis: %w", err)
		}

		// Append retrieved byte slice to the result slice
		retrievedByteSlices = append(retrievedByteSlices, retrievedByteSlice)
	}

	return retrievedByteSlices, nil
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

func writeTempFile(data []byte, prefix string) (string, error) {
	tmpfile, err := ioutil.TempFile("", prefix)
	if err != nil {
		return "", fmt.Errorf("error creating temp file: %w", err)
	}

	if _, err := tmpfile.Write(data); err != nil {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
		return "", fmt.Errorf("error writing to temp file: %w", err)
	}

	if err := tmpfile.Close(); err != nil {
		os.Remove(tmpfile.Name())
		return "", fmt.Errorf("error closing temp file: %w", err)
	}

	return tmpfile.Name(), nil
}

func CombineVideos(bgData, overlayData []byte, outputFormat string, overlayX, overlayY int) ([]byte, error) {
	var outputBuffer bytes.Buffer
	var ffmpegStderr bytes.Buffer

	// Write background and overlay data to temp files
	bgFile, err := writeTempFile(bgData, "bg_video_")
	if err != nil {
		return nil, err
	}
	defer os.Remove(bgFile)

	overlayFile, err := writeTempFile(overlayData, "overlay_video_")
	if err != nil {
		return nil, err
	}
	defer os.Remove(overlayFile)

	// FFmpeg command
	ffmpegCmd := exec.Command("ffmpeg",
		"-i", bgFile, // Background video
		"-i", overlayFile, // Overlay video
		"-filter_complex", fmt.Sprintf("overlay=%d:%d", overlayX, overlayY),
		"-preset", "fast",
		"-f", outputFormat,
		"pipe:1") // Output to stdout

	ffmpegCmd.Stdout = &outputBuffer
	ffmpegCmd.Stderr = &ffmpegStderr

	// Start FFmpeg command
	if err := ffmpegCmd.Run(); err != nil {
		return nil, fmt.Errorf("ffmpeg error during video combination: %w, stderr: %v", err, ffmpegStderr.String())
	}

	return outputBuffer.Bytes(), nil
}
