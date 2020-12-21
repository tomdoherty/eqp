package eqp

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
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
			"name",
			"namespace",
			"pattern",
			"pod_name",
		},
	)
)

// Run starts the polling
func Run() {
	log.Println("Starting eqp")

	c := Config{}
	if err := c.loadConfig(os.Getenv("CONFIG_FILE")); err != nil {
		log.Fatalf("Error loading config: %s", err)
	}

	insecure, err := strconv.ParseBool(getEnvWithDefault("VERIFY_TLS", c.Insecure))
	if err != nil {
		log.Fatalf("invalid value for insecure/VERIFY_TLS: %s", err)
	}

	cfg := elasticsearch.Config{
		Addresses: []string{
			getEnvWithDefault("ELASTICSEARCH_HOST", c.URL),
		},
		Username: getEnvWithDefault("ELASTICSEARCH_USER", c.Username),
		Password: getEnvWithDefault("ELASTICSEARCH_PASSWORD", c.Password),
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
		for _, matcher := range c.Matches {
			fmt.Println(matcher.Name)

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
			if err != nil {
				log.Fatalf("error creating template: %s", err)
			}
			type queryTemplate struct {
				Seconds string
				Pattern string
			}
			q := queryTemplate{
				c.Frequency,
				matcher.Pattern,
			}

			var query bytes.Buffer
			if err = tmpl.Execute(&query, q); err != nil {
				log.Fatalf("error templating query: %s", err)
			}

			res, err = es.Search(
				es.Search.WithContext(context.Background()),
				es.Search.WithIndex(matcher.Index),
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
				logScrapeErrorMatches.WithLabelValues(matcher.Name, namespace, matcher.Pattern, podname).Set(r.Hits.Total.Value)
			}
		}
		sleep, err := time.ParseDuration(c.Frequency)
		if err != nil {
			log.Fatalf("Failed parsing FREQUENCY: %s -  %s", c.Frequency, err)
		}
		time.Sleep(sleep)

	}
}

func getEnvWithDefault(env, fallback string) string {
	if value, ok := os.LookupEnv(env); ok {
		return value
	}
	return fallback
}
