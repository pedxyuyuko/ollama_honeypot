package api

import (
	"encoding/json"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type PullConfig struct {
	Speed         int
	Variance      float64
	WavePeriod    float64
	WaveAmplitude float64
}

func loadPullConfig() PullConfig {
	config := PullConfig{
		Speed:         1048576, // Default 1MB/s
		Variance:      0.2,     // Default 20%
		WavePeriod:    5.0,     // Default 5 seconds
		WaveAmplitude: 0.5,     // Default 50%
	}

	if speedStr := os.Getenv("DOWNLOAD_SPEED"); speedStr != "" {
		if speed, err := strconv.Atoi(speedStr); err == nil {
			config.Speed = speed
		}
	}

	if varianceStr := os.Getenv("DOWNLOAD_SPEED_VARIANCE"); varianceStr != "" {
		if variance, err := strconv.ParseFloat(varianceStr, 64); err == nil {
			config.Variance = variance
		}
	}

	if wavePeriodStr := os.Getenv("DOWNLOAD_SPEED_WAVE_PERIOD"); wavePeriodStr != "" {
		if wavePeriod, err := strconv.ParseFloat(wavePeriodStr, 64); err == nil {
			config.WavePeriod = wavePeriod
		}
	}

	if waveAmplitudeStr := os.Getenv("DOWNLOAD_SPEED_WAVE_AMPLITUDE"); waveAmplitudeStr != "" {
		if waveAmplitude, err := strconv.ParseFloat(waveAmplitudeStr, 64); err == nil {
			config.WaveAmplitude = waveAmplitude
		}
	}

	return config
}

func PullHandler(c *gin.Context) {
	var req struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	// Parse model name and tag
	var modelName, tag string
	parts := strings.Split(req.Name, ":")
	modelName = parts[0]
	if len(parts) > 1 {
		tag = parts[1]
	} else {
		tag = "latest"
	}

	// Check if model exists in registry and get manifest
	var repo string
	if strings.Contains(modelName, "/") {
		repo = modelName
	} else {
		repo = "library/" + modelName
	}
	url := "https://registry.ollama.ai/v2/" + repo + "/manifests/" + tag
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		c.JSON(404, gin.H{"error": "model not found"})
		return
	}

	var manifest struct {
		Layers []struct {
			Digest string `json:"digest"`
			Size   int64  `json:"size"`
		} `json:"layers"`
	}
	if err := json.NewDecoder(io.TeeReader(resp.Body, io.Discard)).Decode(&manifest); err != nil {
		c.JSON(404, gin.H{"error": "invalid manifest"})
		return
	}
	resp.Body.Close()

	// Set headers for streaming
	c.Header("Content-Type", "application/x-ndjson")
	c.Writer.WriteHeader(200)

	config := loadPullConfig()
	interval := 100 * time.Millisecond

	// Channel for progress messages
	ch := make(chan map[string]interface{})
	go func() {
		defer close(ch)
		startTime := time.Now()

		// Simulate pulling manifest
		ch <- map[string]interface{}{"status": "pulling manifest"}
		time.Sleep(200 * time.Millisecond)

		// Simulate pulling layers from manifest
		for _, layer := range manifest.Layers {
			completed := int64(0)
			total := layer.Size
			for completed < total {
				elapsed := time.Since(startTime).Seconds()
				wave := math.Cos(2*math.Pi*elapsed/config.WavePeriod) * config.WaveAmplitude * float64(config.Speed)
				randomError := (rand.Float64()*2 - 1) * config.Variance * float64(config.Speed)
				currentSpeed := float64(config.Speed) + wave + randomError
				if currentSpeed < 0 {
					currentSpeed = 0
				}
				increment := int64(currentSpeed * float64(interval) / float64(time.Second))
				if increment <= 0 {
					increment = 1
				}
				if completed+increment > total {
					increment = total - completed
				}
				completed += increment
				ch <- map[string]interface{}{
					"status":    "pulling layers",
					"digest":    layer.Digest,
					"total":     total,
					"completed": completed,
				}
				time.Sleep(interval)
			}
		}

		// Success
		ch <- map[string]interface{}{"status": "success"}
	}()

	// Stream the messages
	for msg := range ch {
		data, _ := json.Marshal(msg)
		c.Writer.Write(data)
		c.Writer.Write([]byte("\n"))
		if flusher, ok := c.Writer.(interface{ Flush() }); ok {
			flusher.Flush()
		}
	}
}
