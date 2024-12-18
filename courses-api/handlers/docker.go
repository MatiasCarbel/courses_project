package handlers

import (
	"context"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type ContainerInfo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Image   string `json:"image"`
	Created int64  `json:"created"`
	State   string `json:"state"`
	Ports   []Port `json:"ports"`
}

type Port struct {
	IP          string `json:"ip"`
	PrivatePort uint16 `json:"privatePort"`
	PublicPort  uint16 `json:"publicPort"`
	Type        string `json:"type"`
}

func GetContainers(w http.ResponseWriter, r *http.Request) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create Docker client"})
		return
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to list containers"})
		return
	}

	var containerInfos []ContainerInfo
	for _, container := range containers {
		ports := make([]Port, len(container.Ports))
		for i, p := range container.Ports {
			ports[i] = Port{
				IP:          p.IP,
				PrivatePort: p.PrivatePort,
				PublicPort:  p.PublicPort,
				Type:        p.Type,
			}
		}

		containerInfos = append(containerInfos, ContainerInfo{
			ID:      container.ID[:12],
			Name:    container.Names[0][1:], // Remove leading slash
			Status:  container.Status,
			Image:   container.Image,
			Created: container.Created,
			State:   container.State,
			Ports:   ports,
		})
	}

	jsonResponse(w, http.StatusOK, containerInfos)
} 