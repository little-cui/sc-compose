package model

type Compose struct {
	Version       string                   `json:"version",yaml:"version"`
	Microservices map[string]*Microservice `json:"microservices",yaml:"microservices"`
}

type Microservice struct {
	Name             string      `json:"name",yaml:"-"`
	Version          string      `json:"version",yaml:"version"`
	DependsOn        []string    `json:"dependsOn,omitempty",yaml:"dependsOn,omitempty"`
	Instances        []*Instance `json:"instances,omitempty",yaml:"instances:omitempty"`
	InstanceReplicas int64       `json:"instanceReplicas,omitempty",yaml:"instanceReplicas:omitempty"`
}

type Instance struct {
	Port int16 `json:"port,omitempty",yaml:"port,omitempty"`
}
