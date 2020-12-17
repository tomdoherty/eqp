package eqp

import "time"

// Response is a response from an elasticsearch _search query
type Response struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value    float64 `json:"value"`
			Relation string  `json:"relation"`
		} `json:"total"`
		MaxScore float64 `json:"max_score"`
		Hits     []struct {
			Index  string  `json:"_index"`
			Type   string  `json:"_type"`
			ID     string  `json:"_id"`
			Score  float64 `json:"_score"`
			Source struct {
				Stream string `json:"stream"`
				Logtag string `json:"logtag"`
				Log    string `json:"log"`
				Docker struct {
					ContainerID string `json:"container_id"`
				} `json:"docker"`
				Kubernetes struct {
					ContainerName    string `json:"container_name"`
					NamespaceName    string `json:"namespace_name"`
					PodName          string `json:"pod_name"`
					ContainerImage   string `json:"container_image"`
					ContainerImageID string `json:"container_image_id"`
					PodID            string `json:"pod_id"`
					Host             string `json:"host"`
					Labels           struct {
						ControllerRevisionHash         string `json:"controller-revision-hash"`
						AppKubernetesIoInstance        string `json:"app_kubernetes_io/instance"`
						AppKubernetesIoManagedBy       string `json:"app_kubernetes_io/managed-by"`
						AppKubernetesIoName            string `json:"app_kubernetes_io/name"`
						AppKubernetesIoPartOf          string `json:"app_kubernetes_io/part-of"`
						ArgocdArgoprojIoInstance       string `json:"argocd_argoproj_io/instance"`
						StatefulsetKubernetesIoPodName string `json:"statefulset_kubernetes_io/pod-name"`
						StrimziIoCluster               string `json:"strimzi_io/cluster"`
						StrimziIoKind                  string `json:"strimzi_io/kind"`
						StrimziIoName                  string `json:"strimzi_io/name"`
					} `json:"labels"`
					MasterURL       string `json:"master_url"`
					NamespaceID     string `json:"namespace_id"`
					NamespaceLabels struct {
						Name                     string `json:"name"`
						ArgocdArgoprojIoInstance string `json:"argocd_argoproj_io/instance"`
					} `json:"namespace_labels"`
				} `json:"kubernetes"`
				Timestamp time.Time `json:"@timestamp"`
				Tag       string    `json:"tag"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}
