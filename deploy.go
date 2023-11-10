package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
)

func (req *funcUpdateReq) deploy(kubeconfig string) (string, error) {
	var output strings.Builder
	// service manifest
	svcFile := path.Join(repoPath, "func", req.name+`.yaml`)
	namespace := envOrDefault("KUBE_NAMESPACE", "func")

	configBytes, err := base64.StdEncoding.DecodeString(kubeconfig)
	if err != nil {
		return output.String(), fmt.Errorf("fail to decode kubeconfig: %w", err)
	}

	cfgFile := "/tmp/config"

	err = os.WriteFile(cfgFile, configBytes, 0644)
	if err != nil {
		return output.String(), fmt.Errorf("fail to write kubeconfig: %w", err)
	}

	cmd := exec.Command("kubectl", "-n", namespace, "apply", "-f", svcFile)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "KUBECONFIG="+cfgFile)
	cmd.Stdout = io.MultiWriter(os.Stdout, &output)
	cmd.Stderr = io.MultiWriter(os.Stderr, &output)

	err = cmd.Run()
	if err != nil {
		return output.String(), fmt.Errorf("kubectl apply error: %w", err)
	}

	return output.String(), nil
}
