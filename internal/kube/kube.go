package kube

import (
	"context"
	"k8s.io/client-go/rest"

	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	clientcmd "k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"log"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"

	"os"
	"path/filepath"
)

func getCubeConfig(devMode bool) (*rest.Config, error) {
	if devMode {
		var kubeconfigFile = os.Getenv("kubeconfigPath")
		kubeConfigPath := filepath.Join(kubeconfigFile)
		klog.Infof("Using kubeconfig: %s\n", kubeConfigPath)

		kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
		if err != nil {
			klog.Error("error getting Kubernetes config: %v\n", err)
			os.Exit(1)
		}

		return kubeConfig, nil
	} else {
		config, err := rest.InClusterConfig()
		if err != nil {
			return nil, err
		}

		return config, nil
	}
}

func createKubeConfig() (*kubernetes.Clientset, client.Client) {

	config, err := getCubeConfig(devMode)
	if err != nil {
		klog.Info("Config init error...", err)
		os.Exit(1)
	}
	forConfig, err := kubernetes.NewForConfig(config)
	c, _ := client.New(config, client.Options{})
	if err != nil {
		klog.Info("Config init error...", err)
		os.Exit(1)
	}
	return forConfig, c
}

func createJobFunc(clientset *kubernetes.Clientset, gitRepoName string, ciJobName string) {
	println("create job " + gitRepoName + "  " + ciJobName)

	kJobName := gitRepoName + "-" + ciJobName
	dockerStartCommand := "run job command" //todo
	jobFromImage := "run job command"       //todo

	jobs := clientset.BatchV1().Jobs("default")
	var backOffLimit int32 = 0

	jobSpec := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      kJobName,
			Namespace: "ciJobs",
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:    kJobName,
							Image:   jobFromImage,
							Command: strings.Split(dockerStartCommand, " "),
						},
					},
					RestartPolicy: v1.RestartPolicyNever,
				},
			},
			BackoffLimit: &backOffLimit,
		},
	}

	_, err := jobs.Create(context.TODO(), jobSpec, metav1.CreateOptions{})
	if err != nil {
		log.Fatalln("Failed to create K8s job.")
	}

}
