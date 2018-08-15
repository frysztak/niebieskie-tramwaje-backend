package main

import (
	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/docker/ctx"
	"github.com/docker/libcompose/project"
	"github.com/docker/libcompose/project/options"
	"golang.org/x/net/context"
	"log"
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

func createDockerClient() (DockerClient, error) {
	var client DockerClient
	client.context = &ctx.Context{
		Context: project.Context{
			ComposeFiles: []string{"neo4j.yml"},
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
