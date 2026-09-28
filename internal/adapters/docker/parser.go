package docker

import (
	"github.com/duck-ahiru-Z/DevDuck/internal/model"
	"strings"
)

func Parse(stderr string) (model.ErrorInfo, bool) {
	lower := strings.ToLower(stderr)
	kind := "DockerError"
	patterns := []struct{ kind, phrase string }{
		{"DaemonError", "cannot connect to the docker daemon"}, {"DaemonError", "error during connect"}, {"DaemonError", "failed to connect to the docker api"},
		{"PermissionError", "permission denied while trying to connect"}, {"PermissionError", "got permission denied while trying"},
		{"RegistryError", "unauthorized: authentication required"}, {"RegistryError", "pull access denied"}, {"RegistryError", "requested access to the resource is denied"},
		{"ImageError", "unable to find image"}, {"ImageError", "repository does not exist"}, {"ImageError", "manifest unknown"}, {"ImageError", "no matching manifest"},
		{"ContainerError", "no such container"}, {"ContainerError", "is already in use"},
		{"NetworkError", "port is already allocated"}, {"NetworkError", "address already in use"},
		{"NetworkError", "network "}, {"NetworkError", "non-overlapping ipv4 address pool"},
		{"ComposeError", "no configuration file provided"}, {"ComposeError", "undefined network"}, {"ComposeError", "undefined service"}, {"ComposeError", "invalid compose project"},
		{"BuildError", "failed to solve"}, {"BuildError", "failed to read dockerfile"}, {"BuildError", "dockerfile: no such file"}, {"BuildError", "copy failed"}, {"BuildError", "failed to compute cache key"},
	}
	for _, item := range patterns {
		if strings.Contains(lower, item.phrase) {
			kind = item.kind
			break
		}
	}
	message := specificMessage(stderr, kind)
	if message == "" {
		return model.ErrorInfo{}, false
	}
	return model.ErrorInfo{Source: "docker", Kind: kind, Message: message, Raw: stderr}, true
}

func specificMessage(stderr, kind string) string {
	// Prefer the actionable daemon/build failure over incidental client warnings.
	for _, line := range strings.Split(stderr, "\n") {
		line = strings.TrimSpace(line)
		lower := strings.ToLower(line)
		if strings.Contains(lower, "cannot connect to the docker daemon") || strings.Contains(lower, "failed to connect to the docker api") || strings.Contains(lower, "error during connect") {
			return line
		}
	}
	for _, line := range strings.Split(stderr, "\n") {
		line = strings.TrimSpace(line)
		lower := strings.ToLower(line)
		if line == "" {
			continue
		}
		if strings.Contains(lower, "cannot connect") || strings.Contains(lower, "permission denied") || strings.Contains(lower, "unable to find image") || strings.Contains(lower, "pull access denied") || strings.Contains(lower, "manifest unknown") || strings.Contains(lower, "no such container") || strings.Contains(lower, "already in use") || strings.Contains(lower, "allocated") || strings.Contains(lower, "failed to solve") || strings.Contains(lower, "failed to read dockerfile") || strings.Contains(lower, "copy failed") || strings.Contains(lower, "no configuration file") || strings.Contains(lower, "undefined service") || strings.Contains(lower, "undefined network") || strings.Contains(lower, "unauthorized") || strings.Contains(lower, "requested access") || (kind == "NetworkError" && strings.Contains(lower, "network")) {
			return line
		}
	}
	return strings.TrimSpace(strings.Split(stderr, "\n")[0])
}
