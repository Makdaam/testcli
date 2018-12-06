package main

import (
	"flag"
	"fmt"

	appsv1 "github.com/openshift/client-go/apps/clientset/versioned/typed/apps/v1"
	projectsv1 "github.com/openshift/client-go/project/clientset/versioned/typed/project/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeconfig = flag.String("kubeconfig", "", "absolute path to kubeconfig")
)

func doubleReplicas(appClientset *appsv1.AppsV1Client, namespace string) error {
	//doubles the number of replicas in all Deployment Configs in a namespace

	dcList, err := appClientset.DeploymentConfigs(namespace).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, dc := range dcList.Items {
		replicaCount := dc.Spec.Replicas
		fmt.Printf("Deployment Config: %s old replicas: %d new replicas: %d\n", dc.ObjectMeta.Name, replicaCount, replicaCount*2)
		dc.Spec.Replicas = replicaCount * 2
		_, err := appClientset.DeploymentConfigs(namespace).Update(&dc)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {

	flag.Parse()
	fmt.Println("Config path: ", *kubeconfig)
	fmt.Println("Increasing replica counts 2x in all DCs")
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	appClientset, err := appsv1.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	projectClientset, err := projectsv1.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	projectList, err := projectClientset.Projects().List(metav1.ListOptions{})
	for _, project := range projectList.Items {
		fmt.Printf("Project: %s\n", project.ObjectMeta.Name)
		doubleReplicas(appClientset, project.ObjectMeta.Name)
		if err != nil {
			panic(err.Error())
		}
	}
}
