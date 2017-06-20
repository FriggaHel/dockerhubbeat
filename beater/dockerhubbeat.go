package beater

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/FriggaHel/dockerhubbeat/config"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"
)

type Dockerhubbeat struct {
	done   chan struct{}
	config config.Config
	client publisher.Client
}

type DockerhubPerms struct {
	Read  bool `json:"read"`
	Write bool `json:"write"`
	Admin bool `json:"admin"`
}
type DockerhubData struct {
	User            string    `json:"user"`
	Name            string    `json:"name"`
	Namespace       string    `json:"namespace"`
	RepositoryType  string    `json:"repository_type"`
	Status          int       `json:"status"`
	IsPrivate       bool      `json:"is_private"`
	IsAutomated     bool      `json:"is_automated"`
	CanEdit         bool      `json:"can_edit"`
	StarCount       int       `json:"star_count"`
	PullCount       int       `json:"pull_count"`
	LastUpdated     time.Time `json:"last_updated"`
	BuildOnCloud    string    `json:"build_on_cloud"`
	HasStarred      bool      `json:"has_starred"`
	FullDescription string    `json:full_description""`
	Affiliation     string    `json:"affiliation"`
	Permissions     DockerhubPerms
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Dockerhubbeat{
		done:   make(chan struct{}),
		config: config,
	}
	return bt, nil
}

func (bt *Dockerhubbeat) Run(b *beat.Beat) error {
	logp.Info("dockerhubbeat is running! Hit CTRL-C to stop it.")

	bt.client = b.Publisher.Connect()
	ticker := time.NewTicker(bt.config.Period)
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		for _, repository := range bt.config.Repositories {
			data, err := bt.FetchData(&repository)
			if err != nil {
				logp.Warn(fmt.Sprintf("Unable to fetch repository data [%s]", repository.Name))
				continue
			}

			event := common.MapStr{
				"@timestamp":       common.Time(time.Now()),
				"type":             b.Name,
				"repository.name":  repository.FullName(),
				"repository.stars": data.StarCount,
				"repository.pulls": data.PullCount,
			}
			bt.client.PublishEvent(event)
		}
	}
}

func (bt *Dockerhubbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}

func (bt *Dockerhubbeat) FetchData(rep *config.Repository) (DockerhubData, error) {
	res := DockerhubData{}
	r, err := http.Get(fmt.Sprintf("https://hub.docker.com/v2/repositories/%s", rep.FullName()))
	if err != nil || r.StatusCode != 200 {
		return res, err
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		return res, err
	}

	if res.Name == "" {
		return res, errors.New("Unable to unmarshal repository data")
	}
	return res, nil
}
