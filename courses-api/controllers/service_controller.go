package controllers

import (
	"courses-api/models"
	"courses-api/views"
	"encoding/json"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type ServiceController struct{}

func NewServiceController() *ServiceController {
	return &ServiceController{}
}

func (c *ServiceController) GetServices(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("docker", "ps", "--format", "{{.Names}}\t{{.Status}}\t{{.Ports}}")
	output, err := cmd.Output()
	if err != nil {
		views.JSON(w, views.Response{
			Status: http.StatusInternalServerError,
			Error:  "Failed to get service instances",
		})
		return
	}

	services := map[string][]models.ServiceInstance{}
	lines := strings.Split(string(output), "\n")
	
	for _, line := range lines {
		if line == "" {
			continue
		}
		
		parts := strings.Split(line, "\t")
		if len(parts) < 3 {
			continue
		}

		name := parts[0]
		status := "running"
		if !strings.Contains(strings.ToLower(parts[1]), "up") {
			status = "stopped"
		}
		
		serviceType := ""
		if strings.Contains(name, "courses-api") {
			serviceType = "Courses API"
		} else if strings.Contains(name, "users-api") {
			serviceType = "Users API"
		} else if strings.Contains(name, "search-api") {
			serviceType = "Search API"
		} else {
			continue
		}

		instance := models.ServiceInstance{
			ID:        name,
			Name:      name,
			Status:    status,
			Health:    "healthy",
			URL:       getServiceURL(name),
			CreatedAt: time.Now().Format(time.RFC3339),
		}

		services[serviceType] = append(services[serviceType], instance)
	}

	response := []models.ServiceGroup{}
	for serviceName, instances := range services {
		group := models.ServiceGroup{
			Name:         serviceName,
			Instances:    instances,
			MaxInstances: 3,
		}
		response = append(response, group)
	}

	views.JSON(w, views.Response{
		Status: http.StatusOK,
		Data:   response,
	})
}

func (c *ServiceController) AddInstance(w http.ResponseWriter, r *http.Request) {
	var request struct {
		ServiceName string `json:"serviceName"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		views.JSON(w, views.Response{
			Status: http.StatusBadRequest,
			Error:  "Invalid request body",
		})
		return
	}

	serviceName := strings.ToLower(strings.Replace(request.ServiceName, " ", "-", -1))
	cmd := exec.Command("docker-compose", "up", "-d", "--scale", serviceName+"=+1")
	if err := cmd.Run(); err != nil {
		views.JSON(w, views.Response{
			Status: http.StatusInternalServerError,
			Error:  "Failed to create service instance",
		})
		return
	}

	views.JSON(w, views.Response{
		Status: http.StatusOK,
		Data:   "Instance created successfully",
	})
}

func (c *ServiceController) RemoveInstance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	instanceID := vars["id"]
	if instanceID == "" {
		views.JSON(w, views.Response{
			Status: http.StatusBadRequest,
			Error:  "Instance ID is required",
		})
		return
	}

	cmd := exec.Command("docker", "rm", "-f", instanceID)
	if err := cmd.Run(); err != nil {
		views.JSON(w, views.Response{
			Status: http.StatusInternalServerError,
			Error:  "Failed to remove service instance",
		})
		return
	}

	views.JSON(w, views.Response{
		Status: http.StatusOK,
		Data:   "Instance removed successfully",
	})
}

func getServiceURL(name string) string {
	if strings.Contains(name, "courses-api") {
		return "http://courses-api:8002"
	} else if strings.Contains(name, "users-api") {
		return "http://users-api:8001"
	} else if strings.Contains(name, "search-api") {
		return "http://search-api:8003"
	}
	return ""
} 