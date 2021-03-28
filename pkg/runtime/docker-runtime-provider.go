package runtime

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type DockerConfig struct {
	Image       string   `yaml:"image"`
	Environment []string `yaml:"environment"`
	ReadyCheck  []string `yaml:"readyCheck"`
}

type DockerRuntimeProvider struct {
	runtime      *Runtime
	dockerConfig DockerConfig
}

type Runtime struct {
	containerID *string
}

// Setup Docker container
func (p DockerRuntimeProvider) Setup(testName string, dir string) error {
	// create command
	log.Printf("[%s] Startup docker container from image: '%s'\n", testName, p.dockerConfig.Image)
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	workDir := filepath.Join(pwd, dir, "workdir")
	args := []string{
		"run", "-d", "--volume=" + workDir + ":/workdir:ro", "-w=/workdir",
	}
	if p.dockerConfig.Environment != nil {
		for _, env := range p.dockerConfig.Environment {
			args = append(args, "-e", env)
		}
	}
	args = append(args, p.dockerConfig.Image)
	cmd := exec.Command("docker", args...)

	// run command
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	stdOut, err := cmd.Output()
	if err != nil {
		return err
	}

	stdErrSlurp, _ := ioutil.ReadAll(stderr)
	log.Printf("%s", stdErrSlurp)

	containerID := strings.TrimSpace(string(stdOut))
	if len(stdOut) == 0 {
		return fmt.Errorf("[%s] Failed to setup Docker Container", testName)
	}

	p.runtime.containerID = &containerID

	// Wait for container to become ready
	upCount := 0
	if p.dockerConfig.ReadyCheck != nil && len(p.dockerConfig.ReadyCheck) > 0 {
		log.Printf("[%s] Wait for ready check\n", testName)
		for {
			_, execErr := p.Exec(testName, p.dockerConfig.ReadyCheck[0], p.dockerConfig.ReadyCheck[1:]...)
			if execErr == nil {
				upCount++
				if upCount >= 3 {
					break
				}
			} else {
				upCount = 0
			}
			time.Sleep(1 * time.Second)
		}
	}

	return nil
}

// Destroy Docker container
func (p DockerRuntimeProvider) Destroy(testName string, dir string) error {
	if p.runtime.containerID != nil {
		// create command
		log.Printf("[%s] Destroy docker container %s\n", testName, *p.runtime.containerID)
		cmd := exec.Command("docker", "rm", "-f", *p.runtime.containerID)

		// rund command
		return cmd.Run()
	} else {
		log.Printf("[%s] Docker containerID is missing for destroy\n", testName)
	}
	return nil
}

func (p DockerRuntimeProvider) Exec(testName string, command string, args ...string) (*string, error) {
	if p.runtime.containerID == nil {
		return nil, fmt.Errorf("[%s] Docker Container isn't created", testName)
	}

	// create command
	cmdArgs := []string{"exec", *p.runtime.containerID, command}
	cmd := exec.Command("docker", append(cmdArgs, args...)...)

	// run command
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	stdErrSlurp, _ := ioutil.ReadAll(stderr)

	stdOutSlurp, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		log.Printf("[%s] %s", testName, stdErrSlurp)
		return nil, err
	}
	output := string(stdOutSlurp)
	return &output, nil
}

func NewDockerRuntimeProvider(dockerConfig DockerConfig) DockerRuntimeProvider {
	runtime := Runtime{}
	dockerRuntimeProvider := DockerRuntimeProvider{
		dockerConfig: dockerConfig,
		runtime:      &runtime,
	}
	return dockerRuntimeProvider
}
