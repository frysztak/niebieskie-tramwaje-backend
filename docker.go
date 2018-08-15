package main

import (
	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/docker/ctx"
	"github.com/docker/libcompose/project"
	"github.com/docker/libcompose/project/options"
	"golang.org/x/net/context"
	"log"
)

func createProject() (project.APIProject, error) {
	c := ctx.Context{
		Context: project.Context{
			ComposeFiles: []string{"neo4j.yml"},
			ProjectName:  "neo4j",
		},
	}

	project, err := docker.NewProject(&c, nil)
	if err != nil {
		log.Print("Failed to created Docker project")
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	err = project.Up(context.Background(), options.Up{})
	return project, err
}
