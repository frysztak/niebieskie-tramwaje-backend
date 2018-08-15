package main

import (
	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/docker/ctx"
	"github.com/docker/libcompose/project"
	"github.com/docker/libcompose/project/options"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type DockerClient struct {
	project project.APIProject
	context *ctx.Context
}

func (c DockerClient) up() error {
	return c.project.Up(context.Background(), options.Up{})
}

func (c DockerClient) down() error {
	return c.project.Down(context.Background(), options.Down{})
}

func createComposeFile(mpkPath string) (string, error) {
	dockerPath := filepath.Join(mpkPath, "docker")
	str := strings.Replace(composeTemplate, "{NEO4J}", dockerPath, -1)

	ymlPath := filepath.Join(mpkPath, "neo4j.yml")
	log.Printf("Writing neo4j.yml to %s...", ymlPath)
	err := ioutil.WriteFile(ymlPath, []byte(str), 0644)
	return ymlPath, err
}

func createDirectories(mpkPath string) error {
	dockerPath := filepath.Join(mpkPath, "docker")

	dirs := []string{
		"data", "logs", "conf", "import", "plugins",
	}

	for _, dirName := range dirs {
		dir := filepath.Join(dockerPath, dirName)
		err := os.MkdirAll(dir, os.ModeDir|0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func createDockerClient(mpkPath string) (DockerClient, error) {
	ymlPath, err := createComposeFile(mpkPath)
	if err != nil {
		log.Fatal(err)
	}

	err = createDirectories(mpkPath)
	if err != nil {
		log.Fatal(err)
	}

	var client DockerClient
	client.context = &ctx.Context{
		Context: project.Context{
			ComposeFiles: []string{ymlPath},
			ProjectName:  "neo4j",
		},
	}

	project, err := docker.NewProject(client.context, nil)
	if err != nil {
		log.Print("Failed to created Docker project")
		return client, err
	}

	client.project = project
	return client, err
}
