package api

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// FileHook is a logrus hook for writing JSON logs to a file
type FileHook struct {
	file *os.File
}

func (h *FileHook) Fire(entry *logrus.Entry) error {
	formatter := &logrus.JSONFormatter{}
	formatted, err := formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = h.file.Write(formatted)
	return err
}

func (h *FileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

var AuditLogger *logrus.Logger

func InitAuditLogger(logPath string) {
	AuditLogger = logrus.New()
	AuditLogger.SetFormatter(&logrus.TextFormatter{})
	if logPath != "" {
		file, err := os.OpenFile(logPath+"/audit.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Printf("Failed to open audit log file: %v", err)
		} else {
			hook := &FileHook{file: file}
			AuditLogger.AddHook(hook)
		}
	}
}

type RootFS struct {
	Type    string   `json:"type"`
	DiffIDs []string `json:"diff_ids"`
}

type ModelDetails struct {
	ModelFormat   string   `json:"model_format"`
	ModelFamily   string   `json:"model_family"`
	ModelFamilies []string `json:"model_families"`
	ModelType     string   `json:"model_type"`
	FileType      string   `json:"file_type"`
	Architecture  string   `json:"architecture"`
	OS            string   `json:"os"`
	RootFS        RootFS   `json:"rootfs"`
}

type Layer struct {
	Digest    string `json:"digest"`
	Size      int64  `json:"size"`
	MediaType string `json:"mediaType"`
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

type Response struct {
	Text      string `json:"text"`
	RepeatMin int    `json:"repeat_min"`
	RepeatMax int    `json:"repeat_max"`
}

var Responses []Response

func LoadResponses() error {
	file, err := os.Open(MockPath + "/response.json")
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()
	return json.NewDecoder(file).Decode(&Responses)
}

func LoadModels() error {
	file, err := os.Open(MockPath + "/tags.json")
	if err != nil {
		if os.IsNotExist(err) {
			// Create the file with empty models if it doesn't exist
			return SaveModels()
		}
		return err
	}
	defer func() { _ = file.Close() }()
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

func SaveModels() error {
	models := make([]Model, 0, len(Models))
	for _, model := range Models {
		models = append(models, model)
	}
	data := struct {
		Models []Model `json:"models"`
	}{Models: models}
	file, err := os.Create(MockPath + "/tags.json")
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	return encoder.Encode(data)
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
