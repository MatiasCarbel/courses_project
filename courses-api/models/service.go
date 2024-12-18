package models

type ServiceInstance struct {
    ID        string `json:"id"`
    Name      string `json:"name"`
    Status    string `json:"status"`
    Health    string `json:"health"`
    URL       string `json:"url"`
    CreatedAt string `json:"createdAt"`
}

type ServiceGroup struct {
    Name         string            `json:"name"`
    Instances    []ServiceInstance `json:"instances"`
    MaxInstances int              `json:"maxInstances"`
} 