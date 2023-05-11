package kube

import (
	"context"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	"log"
)

func CreateJobFunc(clientset *kubernetes.Clientset, gitRepoName string, ciJobName string, gitHost string, arch string, dockerFileFolder string, pushAddress string, argsVars map[string]string) {

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
		"--no-cache",
		"--frontend",
		"dockerfile.v0",
		"--opt",
		"platform=" + arch,
		"--local",
		"context=/workspace",
		"--local",
		"dockerfile=" + dockerFileFolder,
	}
	for key, value := range argsVars {
		params := []string{
			"--opt",
			"build-arg:" + key + "=" + value}
		args = append(args, params...)
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
							Args:         []string{},
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
