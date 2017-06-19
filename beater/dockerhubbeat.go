package beater

import (
	"encoding/json"
	"fmt"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"
	"net/http"
	"time"

	"github.com/FriggaHel/dockerhubbeat/config"
)

type Dockerhubbeat struct {
	done   chan struct{}
	config config.Config
	client publisher.Client
}

/*
  "user": "library",
  "name": "traefik",
  "namespace": "library",
  "repository_type": "image",
  "status": 1,
  "description": "Tr\u00e6f\u026ak, a modern reverse proxy",
  "is_private": false,
  "is_automated": false,
  "can_edit": false,
  "star_count": 174,
  "pull_count": 9176629,
  "last_updated": "2017-06-16T17:02:31.412870Z",
  "build_on_cloud": null,
  "has_starred": false,
  "full_description": "descr",
  "affiliation": null,
  "permissions": {
    "read": true,
    "write": false,
    "admin": false
  }
}
*/

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

		logp.Info(bt.config.Repository)
		data, _ := bt.FetchData()

		event := common.MapStr{
			"@timestamp":       common.Time(time.Now()),
			"type":             b.Name,
			"repository.name":  bt.config.Repository,
			"repository.stars": data.StarCount,
			"repository.pulls": data.PullCount,
		}
		bt.client.PublishEvent(event)
		logp.Info("Event sent")
	}
}

func (bt *Dockerhubbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}

func (bt *Dockerhubbeat) FetchData() (DockerhubData, error) {
	res := DockerhubData{}
	r, err := http.Get("https://hub.docker.com/v2/repositories/library/traefik/")
	if err != nil {
		return res, err
	}
	err = json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		return res, err
	}
	r.Body.Close()

	fmt.Println(res)
	return res, nil
}
