package kubernetes

import (
	"flag"
	"log"
	"os/user"
	"os"
	"strings"
	"fmt"

	"github.com/achilles-git/k8s-job-notify/env"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Client represents the wrapper of kubernetes API client
type Client struct {
	clientset kubernetes.Interface
}

// NewClient returns Client struct
func NewClient() (*Client, error) {
	config, err := getConfig()
	if err != nil {
		return nil, err
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		clientset: clientSet,
	}, nil
}

// Function to create GKE console URL
func CreateUrl() string {
	usr, err := user.Current()
	if err != nil {
		return "user not found"
		os.Exit(1)
	}
	filePath := usr.HomeDir + "/.kube/config"
	context := clientcmd.GetConfigFromFileOrDie(filePath).CurrentContext
	res := strings.Split(context,"_")
	fmt.Println(res)
	url := "https://console.cloud.google.com/kubernetes/job/"+res[2]+"/"+res[3]+"/default/"
	fmt.Println(url)
	return url
}

func getConfig() (config *rest.Config, err error) {
	if env.IsInCluster() {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
		log.Printf("using inClusterConfig")
	} else {
		usr, err := user.Current()
		if err != nil {
			return nil, err
		}

		filePath := usr.HomeDir + "/.kube/config"
		kubeconfig := flag.String("kubeconfig", filePath, "absolute path to file")
		flag.Parse()
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}

// ListJobs returns the list of Jobs
func (c *Client) ListJobs(namespace string) (*batchv1.JobList, error) {
	jobs, err := c.clientset.BatchV1().Jobs(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return jobs, nil
}
