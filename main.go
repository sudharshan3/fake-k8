package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	"github.com/sirupsen/logrus"
	apps "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func GetRestConfig() (*rest.Config, error) {
	var kubeconfig *string
	home := homedir.HomeDir()
	kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "")

	flag.Parse()

	cfg, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		logrus.WithError(err).Fatal("could not get config")
		return nil, err
	}
	return cfg, nil
}
func Fetcher(namespace string) (*v1.PodList, *v1.NamespaceList, *apps.DeploymentList, *apps.ReplicaSetList, error) {
	config, err := GetRestConfig()
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	Podlister, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	Namespacelister, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	Deploymentlister, err := clientset.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	Replicasetlister, err := clientset.AppsV1().ReplicaSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
		return nil, nil, nil, nil, err
	}
	return Podlister, Namespacelister, Deploymentlister, Replicasetlister, nil
}
func main() {
	namespace := "default"
	Podlister, Namespacelister, Deploymentlister, Replicasetlister, err := Fetcher(namespace)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d namespaces in the cluster\n", len(Namespacelister.Items))
	for k, i := range Namespacelister.Items {
		fmt.Println(k+1, ".", i.Name)
	}
	fmt.Printf("\nThere are %d pods in the %s namespace\n", len(Podlister.Items), namespace)
	for k, i := range Podlister.Items {
		fmt.Println(k+1, ".", i.Name)
	}
	fmt.Printf("\nThere are %d Deployments in the %s namespace\n", len(Deploymentlister.Items), namespace)
	for k, i := range Deploymentlister.Items {
		fmt.Println(k+1, ".", i.Name)
	}
	fmt.Printf("\nThere are %d Replicasets in the %s namespace\n", len(Replicasetlister.Items), namespace)
	for k, i := range Replicasetlister.Items {
		fmt.Println(k+1, ".", i.Name)
	}

}
