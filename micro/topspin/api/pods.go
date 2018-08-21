package api

import (
	"io/ioutil"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type PodList struct {
	Pods []string `json:"pods"`
}

func GetPods() (*PodList, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	namespace, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		return nil, err
	}

	pods, err := clientset.CoreV1().Pods(string(namespace)).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	podList := new(PodList)

	for _, pod := range pods.Items {
		podList.Pods = append(podList.Pods, pod.ObjectMeta.Name)
	}

	return podList, nil
}
