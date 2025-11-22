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

type Manifest struct {
	Layers []Layer `json:"layers"`
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

func handleExistingModelPull(c *gin.Context, fullName string) {
	streamNDJSON(c, func() <-chan map[string]interface{} {
		ch := make(chan map[string]interface{})
		go func() {
			defer close(ch)
			ch <- map[string]interface{}{"status": "pulling manifest"}
			time.Sleep(200 * time.Millisecond)
			// Iterate through all layers
			for _, layer := range Models[fullName].Layers {
				ch <- map[string]interface{}{
					"status":    "pulling layers",
					"digest":    layer.Digest,
					"total":     layer.Size,
					"completed": layer.Size,
				}
				time.Sleep(100 * time.Millisecond)
			}
			ch <- map[string]interface{}{"status": "verifying sha256 digest"}
			time.Sleep(500 * time.Millisecond)
			ch <- map[string]interface{}{"status": "writing manifest"}
			time.Sleep(800 * time.Millisecond)
			ch <- map[string]interface{}{"status": "success"}
		}()
		return ch
	})
}

func fetchManifest(repo, tag string) (Manifest, error) {
	url := "https://registry.ollama.ai/v2/" + repo + "/manifests/" + tag
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		return Manifest{}, err
	}
	var manifest Manifest
	if err := json.NewDecoder(io.TeeReader(resp.Body, io.Discard)).Decode(&manifest); err != nil {
		return Manifest{}, err
	}
	_ = resp.Body.Close()
	return manifest, nil
}

func simulatePull(c *gin.Context, manifest Manifest, config PullConfig) {
	interval := 100 * time.Millisecond
	streamNDJSON(c, func() <-chan map[string]interface{} {
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
			ch <- map[string]interface{}{"status": "verifying sha256 digest"}
			ch <- map[string]interface{}{"status": "writing manifest"}
			ch <- map[string]interface{}{"status": "success"}
		}()
		return ch
	})
}

func addModelToDatabase(fullName string, manifest Manifest) {
	var totalSize int64
	for _, layer := range manifest.Layers {
		totalSize += layer.Size
	}
	Models[fullName] = Model{
		Name:       fullName,
		ModifiedAt: time.Now().Format(time.RFC3339),
		Size:       totalSize,
		Digest:     manifest.Layers[0].Digest,
		Details: ModelDetails{
			Format:            "gguf",
			Family:            "unknown",
			Families:          []string{"unknown"},
			ParameterSize:     "unknown",
			QuantizationLevel: "unknown",
		},
		Layers: manifest.Layers,
	}
}

func PullHandler(c *gin.Context) {
	var req struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}
	// Normalize model name to include tag
	fullName := req.Name
	if !strings.Contains(fullName, ":") {
		fullName += ":latest"
	}
	// Parse model name and tag
	modelName, tag := parseModelName(fullName)
	// Check if model already exists
	if _, exists := Models[fullName]; exists {
		handleExistingModelPull(c, fullName)
		return
	}
	// Check if model exists in registry and get manifest
	var repo string
	if strings.Contains(modelName, "/") {
		repo = modelName
	} else {
		repo = "library/" + modelName
	}
	manifest, err := fetchManifest(repo, tag)
	if err != nil {
		c.JSON(404, gin.H{"error": "model not found"})
		return
	}
	config := loadPullConfig()
	simulatePull(c, manifest, config)
	// Add model to database
	addModelToDatabase(fullName, manifest)
}
