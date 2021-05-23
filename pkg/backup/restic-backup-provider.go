package backup

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"time"
)

type ResticBackupProvider struct {
	config ResticConfig
}

// Restore Restic snapshot
func (p ResticBackupProvider) Restore(testName string, dir string, snapshot *Snapshot, importOptions []string) error {
	log.Printf("[%s] Restoring backup %s from %s...\n", testName, snapshot.Name, p.config.Repository)

	// store password
	if p.config.Password != nil {
		p.config.PasswordFile = filepath.Join(dir, "password")
		err := ioutil.WriteFile(filepath.Join(dir, "password"), []byte(*p.config.Password), 0600)
		if err != nil {
			return err
		}
		defer os.Remove(p.config.PasswordFile)
	}

	// create command
	cmd := exec.Command("restic", "restore", "--verify", "--repo", p.config.Repository, "--password-file", p.config.PasswordFile, "--target", filepath.Join(dir, "workdir"), "latest")
	env := os.Environ()
	if p.config.Env != nil {
		for key, value := range p.config.Env {
			env = append(env, key+"="+value)
		}
	}
	cmd.Env = env

	// run command
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	slurp, _ := ioutil.ReadAll(stderr)

	if err := cmd.Wait(); err != nil {
		log.Printf("[%s] %s", testName, slurp)
		return err
	}

	return nil
}

// Restore Restic snapshot
func (p ResticBackupProvider) ListSnapshots(testName string, dir string) ([]*Snapshot, error) {
	// store password
	if p.config.Password != nil {
		p.config.PasswordFile = filepath.Join(dir, "password")
		err := ioutil.WriteFile(filepath.Join(dir, "password"), []byte(*p.config.Password), 0600)
		if err != nil {
			return nil, err
		}
		defer os.Remove(p.config.PasswordFile)
	}

	// create command
	cmd := exec.Command("restic", "snapshots", "--json", "--repo", p.config.Repository, "--password-file", p.config.PasswordFile)
	env := os.Environ()
	if p.config.Env != nil {
		for key, value := range p.config.Env {
			env = append(env, key+"="+value)
		}
	}
	cmd.Env = env

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// parse output
	resticSnapshots := make([]*ResticSnapshot, 0)
	err = json.Unmarshal(output, &resticSnapshots)
	if err != nil {
		return nil, err
	}

	snapshots := make([]*Snapshot, 0)
	for _, resticSnapshot := range resticSnapshots {
		snapshots = append(snapshots, &Snapshot{
			Time: resticSnapshot.Time,
			Name: resticSnapshot.ShortId,
		})
	}

	sort.Slice(snapshots, func(i, j int) bool {
		return snapshots[i].Time.Before(snapshots[j].Time)
	})
	return snapshots, nil
}

func NewResticBackupProvider(config ResticConfig) ResticBackupProvider {
	resticBackupProvider := ResticBackupProvider{
		config: config,
	}
	return resticBackupProvider
}

type ResticSnapshot struct {
	Time    time.Time `json:"time"`
	ShortId string    `json:"short_id"`
}
