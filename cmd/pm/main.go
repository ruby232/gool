package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var config *Config

func main() {
	var err error
	config, err = loadConfig()
	if err != nil {
		//@todo: Improve error handling and config creation
		fmt.Println("Error loading config:", err)
		fmt.Println("Please config your projects root directory in ~/.config/gool/config.json")
		os.Exit(1)
	}

	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	action := os.Args[1]
	fmt.Printf("Argumento %s\n", action)
	switch action {
	case "cd":
		cd()
	case "ls":
		ls()
	default:
		return fmt.Errorf("Unknown action %s \n", action)
	}
	return nil
}

func ls() {
	projects := getProjects(true)
	for _, project := range projects {
		fmt.Println(project.Name)
		fmt.Println(project.Type)
		fmt.Println("***")
	}
}

func cd() {

}

// Misc functions
type Config struct {
	ProjectsRootDir string `json:"ProjectsRootDir"`
}
type Project struct {
	Name      string
	ShortName string
	Path      string
	Type      string
	IsDdev    bool
}

func getProjects(cache bool) []Project {

	var projects []Project
	if cache {
		projects, _ := loadProjectsFromCache()
		if projects != nil {
			return projects
		}
	}

	projectsRootDir := config.ProjectsRootDir
	projectsDirs := getDirProjects(projectsRootDir)
	for _, projectDir := range projectsDirs {
		name := strings.Replace(projectDir, projectsRootDir, "", -1)
		projects = append(projects, Project{
			Name:      name,
			ShortName: filepath.Base(projectDir),
			Path:      projectDir,
			Type:      getProjectType(projectDir),
			IsDdev:    fileExists(filepath.Join(projectDir, ".ddev")),
		})
	}

	_ = saveCache(projects, "projects.json")
	return projects
}

func loadProjectsFromCache() ([]Project, error) {
	file, errFile := loadCache("projects.json")
	if errFile != nil {
		return nil, errFile
	}
	var projects []Project
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&projects)
	return projects, err
}

//func saveProjectsToFile(projects []Project, filename string) error {
//	file, err := os.Create(filename)
//	if err != nil {
//		return err
//	}
//	defer file.Close()
//
//	encoder := json.NewEncoder(file)
//	return encoder.Encode(projects)
//}

func getProjectType(projectDir string) string {
	composerFile := filepath.Join(projectDir, "composer.json")
	if fileExists(composerFile) {
		content, err := os.ReadFile(composerFile)
		if err == nil {
			if strings.Contains(string(content), "drupal/core-recommended") {
				return "Drupal"
			}
			if strings.Contains(string(content), "laravel") {
				return "Laravel"
			}
		}
		return "Composer"
	}

	packageJsonFile := filepath.Join(projectDir, "package.json")
	if fileExists(packageJsonFile) {
		content, err := os.ReadFile(packageJsonFile)
		if err == nil {
			if strings.Contains(string(content), "astro") {
				return "Astro"
			}
			if strings.Contains(string(content), "react") {
				return "React"
			}
		}
		return "Node"
	}

	phpFile := filepath.Join(projectDir, "index.php")
	if fileExists(phpFile) {
		content, err := os.ReadFile(phpFile)
		if err == nil {
			if strings.Contains(string(content), "DRUPAL_ROOT") {
				return "Drupal 7"
			}
		}
		return "PHP"
	}

	goFile := filepath.Join(projectDir, "main.go")
	if fileExists(goFile) {
		return "Go"
	}

	rustFile := filepath.Join(projectDir, "Cargo.toml")
	if fileExists(rustFile) {
		return "Rust"
	}

	reqFile := filepath.Join(projectDir, "requirements.txt")
	if fileExists(reqFile) {
		return "Python"
	}

	return "Unknown"
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func getDirProjects(projectsDir string) []string {
	var projects []string

	err := filepath.WalkDir(projectsDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() && d.Name() == ".idea" {
			parentDir := filepath.Dir(path)
			projects = append(projects, parentDir)
			return filepath.SkipDir
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path %q: %v\n", projectsDir, err)
	}

	return projects
}

func loadConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(homeDir, ".config", "gool", "config.json")
	if fileExists(configPath) == false {
		os.MkdirAll(filepath.Join(homeDir, ".config", "gool"), 0755)
		_, _ = os.Create(configPath)
	}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func loadCache(cacheName string) (*os.File, error) {
	cacheFile, err := getCacheDir(cacheName)
	if err != nil {
		return nil, err
	}

	return os.Open(cacheFile)
}

func saveCache(data any, cacheName string) error {
	cacheFile, err := getCacheDir(cacheName)
	if err != nil {
		return err
	}

	file, err := os.Create(cacheFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(data)
}

func getCacheDir(cacheName string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	cacheDir := filepath.Join(homeDir, ".cache", "gool")
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		err = os.MkdirAll(cacheDir, 0755)
		if err != nil {
			return "", err
		}
	}

	return filepath.Join(cacheDir, cacheName), nil
}
