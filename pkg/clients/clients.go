package clients

import (
	"net/http"
	"time"
)

type ServiceClients struct {
	BootURL     string
	MetadataURL string
	HTTP        *http.Client
}

func NewServiceClients(bootURL, metadataURL string) *ServiceClients {
	return &ServiceClients{
		BootURL:     bootURL,
		MetadataURL: metadataURL,
		HTTP:        &http.Client{Timeout: 5 * time.Second},
	}
}

// Boot Structures
type BootConfig struct {
	Kernel string `json:"kernel"`
	Params string `json:"params"`
	Image  string `json:"image"`
}

// Metadata Structures
type MetaConfig struct {
	Groups []string `json:"groups"`
}

// FetchBootConfig simulates fetching from boot-service
// In a real implementation, this would hit: GET /boot/v1/bootparameters?host={xname}
func (c *ServiceClients) FetchBootConfig(xname string) (*BootConfig, error) {
	// For MVP/Demo: Returning mock data if service isn't actually running
	// TODO: Replace with real http call:
	// resp, err := c.HTTP.Get(fmt.Sprintf("%s/boot/v1/bootparameters?host=%s", c.BootURL, xname))
    // ... unmarshal logic ...
    
	return &BootConfig{
		Kernel: "vmlinuz-stable",
		Params: "console=tty0",
		Image:  "initrd-stable",
	}, nil
}

// FetchMetaConfig simulates fetching from metadata-service
// In a real implementation, this would hit: GET /cloud-init/v1/groups?host={xname}
func (c *ServiceClients) FetchMetaConfig(xname string) (*MetaConfig, error) {
	return &MetaConfig{
		Groups: []string{"base", "compute", "monitoring"},
	}, nil
}
