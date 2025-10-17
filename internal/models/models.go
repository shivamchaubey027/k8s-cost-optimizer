package models

import "gorm.io/gorm"

type Pod struct {
	gorm.Model
	Name          string `json:"name"`
	Namespace     string `json:"nameSpace"`
	CPURequest    string `json:"cpuRequest"`
	MemoryRequest string `json:"memoryRequest"`
	IpAddress     string `json:"ipAddress"`
}

// {
//     "name": "example-pod",
//     "nameSpace": "default",
//     "cpuRequest": "250m",
//     "memoryRequest": "512Mi",
//     "ipAddress": "192.168.1.100"
// }
