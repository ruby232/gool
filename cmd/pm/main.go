package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
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
	projects := getProjects()
	for _, project := range projects {
		fmt.Println(project.Name)
		fmt.Println(project.Type)
		fmt.Println("***")
	}
}

func cd() {

}

// Misc functions
type Project struct {
	Name      string
	ShortName string
	Path      string
	Type      string
}

func getProjects() []Project {
	projectsRootDir := getProjectsRootDir()
	projectsDirs := getDirProjects(projectsRootDir)
	var projects []Project
	for _, projectDir := range projectsDirs {
		name := strings.Replace(projectDir, projectsRootDir, "", -1)
		projects = append(projects, Project{
			Name:      name,
			ShortName: filepath.Base(projectDir),
			Path:      projectDir,
			Type:      getProjectType(projectDir),
		})
	}
	return projects
}

func getProjectsRootDir() string {
	return "/home/natsu/Proyectos/"
}

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
