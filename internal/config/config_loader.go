package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"slices"

	"github.com/DivyeshMangla/tiet-timetable/internal/model"
	"github.com/DivyeshMangla/tiet-timetable/internal/types"
	"gopkg.in/yaml.v3"
)

const initConfigFlag = "--init-config"

var ErrMissingTimetableURL = errors.New("missing timetable.url in config")

type ConfigLoader struct {
	configFileName string
}

func NewConfigLoader(configFileName string) *ConfigLoader {
	return &ConfigLoader{configFileName: configFileName}
}

func (cl *ConfigLoader) Load() (*Config, error) {
	reader, err := cl.Resolve()
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return cl.Parse(reader)
}

func (cl *ConfigLoader) Resolve() (io.ReadCloser, error) {
	if cl.hasLocalConfig() {
		return cl.openLocalConfig()
	}
	return cl.getBundledConfig()
}

func (cl *ConfigLoader) HandleInitFlag(args []string) (bool, error) {
	if !slices.Contains(args, initConfigFlag) {
		return false, nil
	}
	if cl.hasLocalConfig() {
		fmt.Printf("%s already exists\n", cl.configFileName)
		return true, nil
	}
	return true, cl.copyBundledToLocal()
}

func (cl *ConfigLoader) Parse(in io.Reader) (*Config, error) {
	var root struct {
		Timetable struct {
			URL string `yaml:"url"`
		} `yaml:"timetable"`
		Subjects map[string]struct {
			Name         string `yaml:"name"`
			Abbreviation string `yaml:"abbreviation"`
		} `yaml:"subjects"`
	}
	if err := yaml.NewDecoder(in).Decode(&root); err != nil {
		return nil, fmt.Errorf("failed to decode YAML: %w", err)
	}
	if root.Timetable.URL == "" {
		return nil, ErrMissingTimetableURL
	}
	subjects := make([]model.Subject, 0, len(root.Subjects))
	for code, data := range root.Subjects {
		subjects = append(subjects, model.NewSubject(
			types.SubjectCode(code),
			types.SubjectName(data.Name),
			types.SubjectAbbr(data.Abbreviation),
		))
	}
	return &Config{
		Subjects:     subjects,
		TimetableURL: root.Timetable.URL,
	}, nil
}

func (cl *ConfigLoader) hasLocalConfig() bool {
	_, err := os.Stat(cl.configFileName)
	return err == nil
}

func (cl *ConfigLoader) openLocalConfig() (io.ReadCloser, error) {
	file, err := os.Open(cl.configFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open local config %q: %w", cl.configFileName, err)
	}
	return file, nil
}

func (cl *ConfigLoader) getBundledConfig() (io.ReadCloser, error) {
	file, err := FS.Open("config.yaml")
	if err != nil {
		return nil, fmt.Errorf("bundled config not found: %w", err)
	}
	return file, nil
}

func (cl *ConfigLoader) copyBundledToLocal() error {
	bundled, err := cl.getBundledConfig()
	if err != nil {
		return fmt.Errorf("failed to get bundled config: %w", err)
	}
	defer bundled.Close()

	local, err := os.Create(cl.configFileName)
	if err != nil {
		return fmt.Errorf("failed to create local config %q: %w", cl.configFileName, err)
	}
	defer local.Close()

	if _, err := io.Copy(local, bundled); err != nil {
		return fmt.Errorf("failed to copy config: %w", err)
	}
	fmt.Printf("%s created successfully\n", cl.configFileName)
	return nil
}
