package config

import (
	"os"
	"strings"

	"github.com/kevinburke/ssh_config"
)

// SSHHost represents a single SSH host configuration
type SSHHost struct {
	Name     string
	Hostname string
	User     string
	Port     string
}

// GetSSHHosts reads the SSH config file and returns a list of hosts
func GetSSHHosts() ([]SSHHost, error) {
	configPath := os.ExpandEnv(strings.Replace("~/.ssh/config", "~", "$HOME", 1))
	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfg, err := ssh_config.Decode(f)
	if err != nil {
		return nil, err
	}

	var hosts []SSHHost

	for _, host := range cfg.Hosts {
		for _, pattern := range host.Patterns {
			// Skip wildcard patterns
			if pattern.String() == "*" {
				continue
			}

			hostName := pattern.String()
			hostname := ssh_config.Get(hostName, "HostName")
			user := ssh_config.Get(hostName, "User")
			port := ssh_config.Get(hostName, "Port")

			if port == "" {
				port = "22" // Default SSH port
			}

			// Only add hosts with a hostname (skip configurations without actual hosts)
			if hostname != "" {
				hosts = append(hosts, SSHHost{
					Name:     hostName,
					Hostname: hostname,
					User:     user,
					Port:     port,
				})
			}
		}
	}

	return hosts, nil
}
