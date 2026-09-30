package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/AsWali/CallBoard/internal/server"
)

type tailscaleStatus struct {
	BackendState string
	Self         struct {
		DNSName string
		UserID  int64
	}
	User map[string]struct{ LoginName string }
}

func tailnetOf(statusJSON []byte) (server.Tailnet, error) {
	var st tailscaleStatus
	if err := json.Unmarshal(statusJSON, &st); err != nil {
		return server.Tailnet{}, fmt.Errorf("can't read tailscale status: %v", err)
	}
	if st.BackendState != "Running" {
		return server.Tailnet{}, fmt.Errorf("Tailscale isn't connected (%s); connect it, then try again", st.BackendState)
	}
	t := server.Tailnet{Host: strings.TrimSuffix(st.Self.DNSName, "."), Login: st.User[strconv.FormatInt(st.Self.UserID, 10)].LoginName}
	if t.Host == "" || t.Login == "" {
		return server.Tailnet{}, fmt.Errorf("tailscale status names no machine or login for this computer")
	}
	return t, nil
}

// tailnetUp puts the page on this machine's tailnet name at the same port, over
// HTTPS, and returns what takes it off again.
func tailnetUp(port int) (server.Tailnet, func(), error) {
	if _, err := exec.LookPath("tailscale"); err != nil {
		return server.Tailnet{}, nil, fmt.Errorf("--tailnet needs the tailscale command on your PATH")
	}
	out, err := exec.Command("tailscale", "status", "--json").Output()
	if err != nil {
		return server.Tailnet{}, nil, fmt.Errorf("tailscale status: %v", err)
	}
	t, err := tailnetOf(out)
	if err != nil {
		return server.Tailnet{}, nil, err
	}
	https := fmt.Sprintf("--https=%d", port)
	if out, err := exec.Command("tailscale", "serve", "--bg", https, fmt.Sprintf("http://127.0.0.1:%d", port)).CombinedOutput(); err != nil {
		return server.Tailnet{}, nil, fmt.Errorf("tailscale serve: %v\n%s", err, strings.TrimSpace(string(out)))
	}
	return t, func() { exec.Command("tailscale", "serve", https, "off").Run() }, nil
}
