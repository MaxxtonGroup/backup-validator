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
	DumpLogs    bool     `yaml:"dumpLogs"`
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
	if p.runtime.containerID != nil {
		// cleanup old container
		p.Destroy(testName, dir)
	}

	// create command
	log.Printf("[%s] Startup docker container from image: '%s'", testName, p.dockerConfig.Image)
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Create workdir if not exists
	mntDir := filepath.Join(pwd, dir)
	workDir := filepath.Join(pwd, dir, "workdir")
	_, err = os.Stat(workDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(workDir, os.ModePerm)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	args := []string{
		"run", "-d", "--volume=" + mntDir + ":/mnt/host", "-w=/mnt/host/workdir",
	}
	if p.dockerConfig.Environment != nil {
		for _, env := range p.dockerConfig.Environment {
			args = append(args, "-e", env)
		}
	}
	args = append(args, p.dockerConfig.Image)
	log.Printf("Run: docker %s", strings.Join(args, " "))
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
	if len(stdErrSlurp) > 0 {
		log.Printf("%s", stdErrSlurp)
	}

	containerID := strings.TrimSpace(string(stdOut))
	if containerID == "" {
		return fmt.Errorf("[%s] Failed to setup Docker Container", testName)
	}

	p.runtime.containerID = &containerID

	// Wait for container to become ready
	upCount := 0
	if p.dockerConfig.ReadyCheck != nil && len(p.dockerConfig.ReadyCheck) > 0 {
		log.Printf("[%s] Wait for ready check", testName)
		for {
			_, execErr := p.Exec(testName, p.dockerConfig.ReadyCheck[0], p.dockerConfig.ReadyCheck[1:]...)
			if execErr == nil {
				upCount++
				if upCount >= 5 {
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
		if p.dockerConfig.DumpLogs {
			// dump logs to stdout
			log.Printf("[%s] Dump docker container logs:%s", testName, *p.runtime.containerID)
			logCmd := exec.Command("docker", "logs", *p.runtime.containerID)

			logs, err := logCmd.CombinedOutput()
			if err != nil {
				log.Printf("Failed to get logs: %s", err)
				println("output: " + string(logs))
			} else {
				for _, line := range strings.Split(string(logs), "\n") {
					log.Printf("[%s] logs: %s", testName, line)
				}
			}
		}

		// create command
		log.Printf("[%s] Destroy docker container %s", testName, *p.runtime.containerID)
		cmd := exec.Command("docker", "rm", "-f", *p.runtime.containerID)

		// run command
		return cmd.Run()
	} else {
		log.Printf("[%s] Docker containerID is missing for destroy", testName)
	}
	return nil
}

func (p DockerRuntimeProvider) Exec(testName string, command string, args ...string) (*string, error) {
	return p.execAsUser(testName, nil, command, args...)
}

func (p DockerRuntimeProvider) ExecRoot(testName string, command string, args ...string) (*string, error) {
	rootUID := "0"
	return p.execAsUser(testName, &rootUID, command, args...)
}

func (p DockerRuntimeProvider) execAsUser(testName string, uid *string, command string, args ...string) (*string, error) {
	if p.runtime.containerID == nil {
		return nil, fmt.Errorf("[%s] Docker Container isn't created", testName)
	}

	// create command
	cmdArgs := []string{"exec"}
	if uid != nil {
		cmdArgs = append(cmdArgs, "-u", *uid)
	}
	cmdArgs = append(cmdArgs, *p.runtime.containerID, command)
	// log.Printf("[%s] exec: docker %s\n", testName, strings.Join(append(cmdArgs, args...), " "))
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
		return nil, fmt.Errorf("command [%s, %s] failed: %s", command, strings.Join(args, ", "), err)
	}

	stdOutSlurp, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, err
	}
	stdErrSlurp, _ := ioutil.ReadAll(stderr)
	if err != nil {
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		output := strings.TrimSpace(string(stdOutSlurp))
		if len(output) > 0 {
			log.Printf("[%s] exec: %s", testName, output)
		}

		errOutput := strings.TrimSpace(string(stdErrSlurp))
		if len(errOutput) > 0 {
			log.Printf("[%s] exec: %s", testName, errOutput)
		}
		return nil, fmt.Errorf("command [%s %s] failed: %s", "docker", strings.Join(append(cmdArgs, args...), " "), err)
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
