package backup

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/MaxxtonGroup/backup-validator/pkg/runtime"
	"github.com/dustin/go-humanize"
)

type ElasticsearchBackupProvider struct {
	runtimeProvider runtime.RuntimeProvider
}

type ElasticsearchSnapshotResponse struct {
	Snapshots []*ElasticsearchSnapshot `json:"snapshots"`
}

type ElasticsearchSnapshot struct {
	Snapshot  string   `json:"snapshot"`
	StartTime string   `json:"start_time"`
	State     string   `json:"state"`
	Indices   []string `json:"indices"`
}

type ElastcisearchRestoreOptions struct {
	Indices string `json:"indices"`
}

type ElasticsearchRestore struct {
	Index          string `json:"string"`
	Type           string `json:"type"`
	Stage          string `json:"stage"`
	Repository     string `json:"repository"`
	Snapshot       string `json:"snapshot"`
	BytesTotal     string `json:"bytes_total"`
	BytesRecovered string `json:"bytes_recovered"`
}

func (p ElasticsearchBackupProvider) Restore(testName string, dir string, snapshot *Snapshot, importOptions []string) error {
	log.Printf("[%s] Restoring backup %s...\n", testName, snapshot.Name)

	restoreOptions := &ElastcisearchRestoreOptions{}
	for _, option := range importOptions {
		p := strings.Split(option, "=")
		key := p[0]
		value := strings.Join(p[1:], ",")
		switch key {
		case "indices":
			restoreOptions.Indices = value
		}
	}
	if restoreOptions.Indices != "" {
		restoreOptions.Indices = snapshot.Time.Add(-(24 * time.Hour)).Format(restoreOptions.Indices)
	}

	restoreOptionsString, err := json.Marshal(&restoreOptions)
	if err != nil {
		return err
	}

	output, err := p.runtimeProvider.Exec(testName, "curl", "--output", "/dev/stdout", "--write-out", "%{http_code}", "-X", "POST", "http://localhost:9200/_snapshot/backup/"+snapshot.Name+"/_restore", "-H", "Content-Type: application/json", "-d", string(restoreOptionsString))
	if err != nil {
		return err
	}
	if !strings.HasSuffix(*output, "200") {
		return fmt.Errorf("the requested URL returned error: %s", *output)
	}

	// Wait for recovery to complete
	var previousProgress string
	for {
		time.Sleep(5 * time.Second)
		output, err := p.runtimeProvider.Exec(testName, "curl", "http://localhost:9200/_cat/recovery?format=json")
		if err != nil {
			return err
		}

		esRestores := []*ElasticsearchRestore{}
		err = json.Unmarshal([]byte(*output), &esRestores)
		if err != nil {
			return err
		}

		bytesTotal := float64(0)
		bytesRecovered := float64(0)
		done := true
		for _, restore := range esRestores {
			if restore.Repository == "backup" && restore.Snapshot == snapshot.Name {
				bTotal, err := humanize.ParseBytes(restore.BytesTotal)
				if err != nil {
					return err
				}
				bytesTotal += float64(bTotal)
				bRecovered, err := humanize.ParseBytes(restore.BytesRecovered)
				if err != nil {
					return err
				}
				bytesRecovered += float64(bRecovered)
				if restore.Stage != "done" {
					done = false
				}
			}
		}
		newProgress := fmt.Sprintf("%0.2f%%", bytesRecovered/bytesTotal*100)
		if newProgress != previousProgress {
			log.Printf("[%s] Restoring %s", testName, newProgress)
			previousProgress = newProgress
		}
		if done {
			break
		}
	}
	return nil
}

func (p ElasticsearchBackupProvider) ListSnapshots(testName string, dir string) ([]*Snapshot, error) {
	log.Printf("[%s] List snapshots...\n", testName)
	output, err := p.runtimeProvider.Exec(testName, "curl", "-X", "GET", "http://localhost:9200/_snapshot/backup/_all")
	if err != nil {
		return nil, err
	}

	esSnapshotResponse := &ElasticsearchSnapshotResponse{}
	err = json.Unmarshal([]byte(*output), &esSnapshotResponse)
	if err != nil {
		return nil, err
	}

	snapshots := make([]*Snapshot, 0)
	for _, esSnapshot := range esSnapshotResponse.Snapshots {
		if esSnapshot.State == "SUCCESS" {
			startTime, err := time.Parse(time.RFC3339, esSnapshot.StartTime)
			if err != nil {
				return nil, err
			}
			snapshots = append(snapshots, &Snapshot{
				Time:      startTime,
				Name:      esSnapshot.Snapshot,
				Databases: esSnapshot.Indices,
			})
		}
	}

	sort.Slice(snapshots, func(i, j int) bool {
		return snapshots[i].Time.Before(snapshots[j].Time)
	})
	return snapshots, nil
}

func NewElasticsearchBackupProvider(runtimeProvider runtime.RuntimeProvider) ElasticsearchBackupProvider {
	elasticsearchBackupProvider := ElasticsearchBackupProvider{
		runtimeProvider: runtimeProvider,
	}
	return elasticsearchBackupProvider
}
