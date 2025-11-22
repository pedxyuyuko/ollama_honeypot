package api

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type ModelDetails struct {
	Format            string   `json:"format"`
	Family            string   `json:"family"`
	Families          []string `json:"families"`
	ParameterSize     string   `json:"parameter_size"`
	QuantizationLevel string   `json:"quantization_level"`
}

type Layer struct {
	Digest string `json:"digest"`
	Size   int64  `json:"size"`
}

type Model struct {
	Name       string       `json:"name"`
	ModifiedAt string       `json:"modified_at"`
	Size       int64        `json:"size"`
	Digest     string       `json:"digest"`
	Details    ModelDetails `json:"details"`
	Layers     []Layer      `json:"layers"`
}

var Models = make(map[string]Model)

func LoadModels() error {
	file, err := os.Open(MockPath + "/tags.json")
	if err != nil {
		return err
	}
	defer file.Close()
	var data struct {
		Models []Model `json:"models"`
	}
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return err
	}
	for _, model := range data.Models {
		Models[model.Name] = model
	}
	return nil
}

var MockPath = "./mock"

// parseModelName parses a model name string into model name and tag.
// If no tag is provided, defaults to "latest".
func parseModelName(name string) (modelName, tag string) {
	parts := strings.Split(name, ":")
	modelName = parts[0]
	if len(parts) > 1 {
		tag = parts[1]
	} else {
		tag = "latest"
	}
	return
}

// streamNDJSON sets up NDJSON streaming response and streams messages from the provided channel generator.
func streamNDJSON(c *gin.Context, generator func() <-chan map[string]interface{}) {
	c.Header("Content-Type", "application/x-ndjson")
	c.Writer.WriteHeader(200)
	ch := generator()
	for msg := range ch {
		data, _ := json.Marshal(msg)
		_, _ = c.Writer.Write(data)
		_, _ = c.Writer.Write([]byte("\n"))
		if flusher, ok := c.Writer.(interface{ Flush() }); ok {
			flusher.Flush()
		}
	}
}
