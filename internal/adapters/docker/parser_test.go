package docker

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetect(t *testing.T) {
	a := NewAdapter()
	for _, c := range []string{"docker", "docker.exe", "/usr/bin/docker", "docker-compose"} {
		if !a.Detect(c, nil) {
			t.Errorf("Detect(%q)=false", c)
		}
	}
}
func TestParseFixtures(t *testing.T) {
	cases := map[string]string{"daemon": "DaemonError", "permission": "PermissionError", "image": "ImageError", "pull_denied": "RegistryError", "manifest": "ImageError", "container": "ContainerError", "conflict": "ContainerError", "port": "NetworkError", "dockerfile": "BuildError", "solve": "BuildError", "copy": "BuildError", "network": "NetworkError", "compose": "ComposeError", "compose_service": "ComposeError", "unauthorized": "RegistryError", "unknown": "DockerError"}
	for name, kind := range cases {
		t.Run(name, func(t *testing.T) {
			d, e := os.ReadFile(filepath.Join("testdata", name+".txt"))
			if e != nil {
				t.Fatal(e)
			}
			i, ok := Parse(string(d))
			if !ok || i.Kind != kind || i.Source != "docker" {
				t.Fatalf("%#v %v", i, ok)
			}
		})
	}
}
