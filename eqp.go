package eqp

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	logScrapeErrorMatches = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "eqp_log_scrape_matches",
			Help: "The total number of patterns matched",
		},
		[]string{
			"namespace",
			"pattern",
			"pod_name",
		},
	)
)

// Run starts the polling
func Run() {
	log.Println("Starting eqp")
	esIndex := os.Getenv("ELASTICSEARCH_INDEX")
	matches := strings.Split(os.Getenv("MATCHES"), ",")
	seconds := os.Getenv("FREQUENCY")
	verifyTLS := os.Getenv("VERIFY_TLS")
	insecure := false

	if verifyTLS == "false" {
		insecure = true
	}

	prometheus.MustRegister(logScrapeErrorMatches)

	cfg := elasticsearch.Config{
		Addresses: []string{
			os.Getenv("ELASTICSEARCH_HOST"),
		},
		Username: os.Getenv("ELASTICSEARCH_USER"),
		Password: os.Getenv("ELASTICSEARCH_PASSWORD"),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	for {
		for _, match := range matches {
			tmpl, err := template.New("query").Parse(`{
  "query": {
    "bool": {
      "must": [
        {
          "regexp": {
            "log": {
              "value": "{{ .Pattern }}"
            }
          }
        },
        {
          "range": {
            "@timestamp": {
              "gte": "now-{{ .Seconds }}/s",
              "lte": "now/s"
            }
          }
        }
      ]
    }
  }
}`)
			type Match struct {
				Pattern string
				Seconds string
			}
			pattern := Match{match, seconds}

			var query bytes.Buffer
			err = tmpl.Execute(&query, pattern)

			res, err = es.Search(
				es.Search.WithContext(context.Background()),
				es.Search.WithIndex(esIndex),
				es.Search.WithBody(&query),
			)
			if err != nil {
				log.Fatalf("Error getting response: %s", err)
			}
			defer res.Body.Close()

			if res.IsError() {
				var e map[string]interface{}
				if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
					log.Fatalf("Error parsing the response body: %s", err)
				} else {
					log.Fatalf("[%s] %s: %s",
						res.Status(),
						e["error"].(map[string]interface{})["type"],
						e["error"].(map[string]interface{})["reason"],
					)
				}
			}

			var r Response
			if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
				log.Fatalf("Error parsing the response body: %s", err)
			}
			for _, hit := range r.Hits.Hits {
				podname := hit.Source.Kubernetes.PodName
				namespace := hit.Source.Kubernetes.NamespaceName
				logScrapeErrorMatches.WithLabelValues(namespace, match, podname).Set(r.Hits.Total.Value)
			}
		}
		sleep, err := time.ParseDuration(seconds)
		if err != nil {
			log.Fatalf("Failed parsing FREQUENCY: %s -  %s", seconds, err)
		}
		time.Sleep(sleep)
	}
}
