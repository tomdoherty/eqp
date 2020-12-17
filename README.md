# eqp - elasticsearch query as prometheus metrics

[![](https://goreportcard.com/badge/github.com/tomdoherty/eqp)](https://goreportcard.com/report/github.com/tomdoherty/eqp)

Queries Elasticsearch for patterns and exposes them as metrics to
prometheus

With this you can generate Alertmanager notifications based on log entries

## configuration

`eqp` is configured using the below environment variables

| Option                 | Description                  | Default                                               |
| :--------------------: | :--------------------------: | :---------------------------------------------------: |
| ELASTICSEARCH_HOST     | Elastic URL                  | https://elasticsearch-es-http.elastic-system.svc:9200 |
| ELASTICSEARCH_USER     | Elastic User                 | elastic                                               |
| ELASTICSEARCH_PASSWORD | Elastic Pass                 | changeme                                              |
| ELASTICSEARCH_INDEX    | Elastic Index                | log-syslog-serverlog-kubernetes-*                     |
| MATCHES                | Comma seperated search terms | ruok,error,WARN                                       |
| FREQUENCY              | Frequency to check           | 20s                                                   |
| VERIFY_TLS             | Skip TLS validation          | false                                                 |


When running inside Kubernetes you can create a secret that works with
the provided manifest as follows:

```shell
kubectl create secret generic elastic-user --from-literal password=changeme
```

The below metrics are provided:

```shell
# HELP eqp_log_scrape_matches The total number of patterns matched
# TYPE eqp_log_scrape_matches gauge
eqp_log_scrape_matches{namespace="kafka",pattern="ruok",pod_name="kafka-zookeeper-0"} 11
eqp_log_scrape_matches{namespace="kafka",pattern="ruok",pod_name="kafka-zookeeper-1"} 11
eqp_log_scrape_matches{namespace="kafka",pattern="ruok",pod_name="kafka-zookeeper-2"} 11
```
