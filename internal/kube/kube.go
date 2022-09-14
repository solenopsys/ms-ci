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
	"os"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

func СreateKubeConfig(devMode bool) (*kubernetes.Clientset, client.Client) {

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

func CreateJobFunc(clientset *kubernetes.Clientset, gitRepoName string, ciJobName string, gitHost string, arch string, dockerFileFolder string, pushAddress string) {

	const SW_NAMESPACE = "shockwaves"
	jobs := clientset.BatchV1().Jobs(SW_NAMESPACE)
	var backOffLimit int32 = 0

	volumeMount := []v1.VolumeMount{
		{
			Name:      "workspace",
			MountPath: "/workspace",
		},
	}
	args := []string{
		"build",
		"--frontend",
		"dockerfile.v0",
		"--opt",
		"platform=" + arch,
		"--local",
		"context=/workspace",
		"--local",
		"dockerfile=" + dockerFileFolder,
	}
	if pushAddress != "" {
		tail := []string{
			"--output",
			"type=image,name=" + pushAddress + ",push=true"}
		args = append(args, tail...)
	}
	jobSpec := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ciJobName,
			Namespace: SW_NAMESPACE,
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{"container.apparmor.security.beta.kubernetes.io/buildkit": "unconfined"},
				},
				Spec: v1.PodSpec{
					InitContainers: []v1.Container{
						{
							Name:         "prepare",
							Image:        "alpine/git",
							Command:      []string{"git", "clone", "http://admin:root@" + gitHost + "/" + gitRepoName, "/workspace/" + gitRepoName},
							VolumeMounts: volumeMount,
						},
					},
					Containers: []v1.Container{
						{
							Name:    "buildkit",
							Image:   "moby/buildkit:master-rootless",
							Command: []string{"buildctl-daemonless.sh"},
							Env: []v1.EnvVar{{
								Name:  "BUILDKITD_FLAGS",
								Value: "--oci-worker-no-process-sandbox",
							}},
							Args:         args,
							VolumeMounts: volumeMount,
						},
					},
					Volumes: []v1.Volume{
						{Name: "workspace", VolumeSource: v1.VolumeSource{
							EmptyDir: &v1.EmptyDirVolumeSource{},
						}},
					},
					RestartPolicy: v1.RestartPolicyNever,
				},
			},
			BackoffLimit: &backOffLimit,
		},
	}

	_, err := jobs.Create(context.TODO(), jobSpec, metav1.CreateOptions{})
	if err != nil {
		klog.Error(err)
		log.Fatalln("Failed to create K8s job.")
	}

}
