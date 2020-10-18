package controllers

import (
	cronjobv1 "github.com/zhu733756/k8s-crawler-operator/api/v1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// new a publisher job
func (r *CronCrawlerJobReconciler) NewPublisherJob(croncrawlerjob *cronjobv1.CronCrawlerJob, publisherName string) *batchv1.Job {
	labels := LabelsForCronCrawlerJob(croncrawlerjob.Name, "publisher", "job")
	//selector := &metav1.LabelSelector{MatchLabels: labels}
	//manualSelector := true
	return &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "batch/v1",
			Kind:       "Job",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      publisherName,
			Namespace: croncrawlerjob.Namespace,
			Labels:    labels,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(croncrawlerjob, schema.GroupVersionKind{
					Group:   v1.SchemeGroupVersion.Group,
					Version: v1.SchemeGroupVersion.Version,
					Kind:    "CronCrawlerJob",
				}),
			},
		},
		Spec: batchv1.JobSpec{
			// TTLSecondsAfterFinished: croncrawlerjob.Spec.Publisher.PublisherTTL,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Volumes:       r.GetVolumes(croncrawlerjob),
					Containers:    r.NewPublisherContainers(croncrawlerjob),
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
			//	ManualSelector: &manualSelector,
			//	Selector:       selector,
		},
	}

}

// new a publisher cronjob
func (r *CronCrawlerJobReconciler) NewPublisherCronJob(croncrawlerjob *cronjobv1.CronCrawlerJob, publisherName string) *batchv1beta1.CronJob {
	labels := LabelsForCronCrawlerJob(croncrawlerjob.Name, "publisher", "cronjob")
	//selector := &metav1.LabelSelector{MatchLabels: labels}
	//manualSelector := true
	return &batchv1beta1.CronJob{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "batch/v1beta1",
			Kind:       "CronJob",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      publisherName,
			Namespace: croncrawlerjob.Namespace,
			Labels:    labels,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(croncrawlerjob, schema.GroupVersionKind{
					Group:   v1.SchemeGroupVersion.Group,
					Version: v1.SchemeGroupVersion.Version,
					Kind:    "CronCrawlerJob",
				}),
			},
		},
		Spec: batchv1beta1.CronJobSpec{
			Schedule: croncrawlerjob.Spec.Schedule,
			JobTemplate: batchv1beta1.JobTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: batchv1.JobSpec{
					TTLSecondsAfterFinished: croncrawlerjob.Spec.Publisher.PublisherTTL,
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: labels,
						},
						Spec: corev1.PodSpec{
							Volumes:       r.GetVolumes(croncrawlerjob),
							Containers:    r.NewPublisherContainers(croncrawlerjob),
							RestartPolicy: corev1.RestartPolicyNever,
						},
					},
					//Selector:       selector,
					//ManualSelector: &manualSelector,
				},
			},
		},
	}
}

// new a job for consumers
func (r *CronCrawlerJobReconciler) NewConsumersJob(croncrawlerjob *cronjobv1.CronCrawlerJob, consumerName string) *batchv1.Job {
	labels := LabelsForCronCrawlerJob(croncrawlerjob.Name, "consumer", "job")
	//selector := &metav1.LabelSelector{MatchLabels: labels}
	//manualSelector := true
	return &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "batch/v1",
			Kind:       "Job",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      consumerName,
			Namespace: croncrawlerjob.Namespace,
			Labels:    labels,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(croncrawlerjob, schema.GroupVersionKind{
					Group:   v1.SchemeGroupVersion.Group,
					Version: v1.SchemeGroupVersion.Version,
					Kind:    "CronCrawlerJob",
				}),
			},
		},
		Spec: batchv1.JobSpec{
			TTLSecondsAfterFinished: croncrawlerjob.Spec.Consumer.ConsumerTTL,
			Parallelism:             croncrawlerjob.Spec.Consumer.Parallelism,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Volumes:        r.GetVolumes(croncrawlerjob),
					Containers:     r.NewConsumerContainers(croncrawlerjob, croncrawlerjob.Spec.Consumer.Containers, "containers"),
					InitContainers: r.NewConsumerContainers(croncrawlerjob, croncrawlerjob.Spec.Consumer.InitContainers, "initContainers"),
					RestartPolicy:  corev1.RestartPolicyNever,
				},
			},
			//		ManualSelector: &manualSelector,
			//		Selector:       selector,
		},
	}

}

// new a cronjob for consumers
func (r *CronCrawlerJobReconciler) NewConsumersCronJob(croncrawlerjob *cronjobv1.CronCrawlerJob, consumerName string) *batchv1beta1.CronJob {
	labels := LabelsForCronCrawlerJob(croncrawlerjob.Name, "consumer", "cronjob")
	//	selector := &metav1.LabelSelector{MatchLabels: labels}
	//	manualSelector := true
	return &batchv1beta1.CronJob{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "batch/v1beta1",
			Kind:       "CronJob",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      consumerName,
			Namespace: croncrawlerjob.Namespace,
			Labels:    labels,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(croncrawlerjob, schema.GroupVersionKind{
					Group:   v1.SchemeGroupVersion.Group,
					Version: v1.SchemeGroupVersion.Version,
					Kind:    "CronCrawlerJob",
				}),
			},
		},
		Spec: batchv1beta1.CronJobSpec{
			Schedule: croncrawlerjob.Spec.Schedule,
			JobTemplate: batchv1beta1.JobTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: batchv1.JobSpec{
					TTLSecondsAfterFinished: croncrawlerjob.Spec.Consumer.ConsumerTTL,
					Parallelism:             croncrawlerjob.Spec.Consumer.Parallelism,
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: labels,
						},
						Spec: corev1.PodSpec{
							Volumes:        r.GetVolumes(croncrawlerjob),
							Containers:     r.NewConsumerContainers(croncrawlerjob, croncrawlerjob.Spec.Consumer.Containers, "containers"),
							InitContainers: r.NewConsumerContainers(croncrawlerjob, croncrawlerjob.Spec.Consumer.InitContainers, "initContainers"),
							RestartPolicy:  corev1.RestartPolicyNever,
						},
					},
					//	Selector:       selector,
					//	ManualSelector: &manualSelector,
				},
			},
		},
	}
}

func (r *CronCrawlerJobReconciler) GetVolumes(croncrawlerjob *cronjobv1.CronCrawlerJob) []corev1.Volume {
	volumes := croncrawlerjob.Spec.Volumes
	if len(volumes) == 0 {
		return []corev1.Volume{
			corev1.Volume{
				Name:         "code",
				VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}},
			},
			corev1.Volume{
				Name:         "run-list",
				VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{LocalObjectReference: corev1.LocalObjectReference{Name: "crawler-run-list"}}},
			},
		}
	}
	return volumes
}

// build containers of the publisher
func (r *CronCrawlerJobReconciler) NewPublisherContainers(croncrawlerjob *cronjobv1.CronCrawlerJob) []corev1.Container {
	containers := croncrawlerjob.Spec.Publisher.Containers
	if len(containers) == 0 {
		volumeMount := []corev1.VolumeMount{
			corev1.VolumeMount{
				Name:      "run-list",
				MountPath: "/run_list",
			},
		}
		return []corev1.Container{corev1.Container{
			Name:            "crawler-publisher",
			Image:           "zhu733756/crawler-publisher:v1",
			Resources:       croncrawlerjob.Spec.Resources,
			Args:            []string{},
			VolumeMounts:    volumeMount,
			Command:         []string{"/bin/sh", "/app/prepare.sh"},
			ImagePullPolicy: corev1.PullIfNotPresent,
			Env:             croncrawlerjob.Spec.Env,
		},
		}
	}
	newContainers := make([]corev1.Container, len(containers))
	for index, container := range containers {
		newContainers[index] = corev1.Container{
			Name:            container.Name,
			Image:           container.Image,
			Resources:       croncrawlerjob.Spec.Resources,
			Args:            container.Args,
			VolumeMounts:    container.VolumeMounts,
			Command:         container.Command,
			ImagePullPolicy: corev1.PullIfNotPresent,
			Env:             croncrawlerjob.Spec.Env,
		}
	}
	return newContainers
}

// build containers of the Consumers
func (r *CronCrawlerJobReconciler) NewConsumerContainers(croncrawlerjob *cronjobv1.CronCrawlerJob, containers []cronjobv1.Container, containersType string) []corev1.Container {
	if len(containers) == 0 {
		volumeMounts := []corev1.VolumeMount{
			corev1.VolumeMount{
				Name:      "code",
				MountPath: "/app/code",
			},
		}
		if containersType == "initContainers" {
			return []corev1.Container{corev1.Container{
				Name:            "pull-code",
				Image:           "zhu733756/pull-code:v1",
				Resources:       croncrawlerjob.Spec.Resources,
				Args:            []string{},
				VolumeMounts:    volumeMounts,
				Command:         []string{"/bin/sh", "/app/pull.sh"},
				ImagePullPolicy: corev1.PullIfNotPresent,
				Env:             croncrawlerjob.Spec.Env,
			},
			}
		}
		return []corev1.Container{corev1.Container{
			Name:            "crawler-runner",
			Image:           "zhu733756/crawler-runner:v1",
			Resources:       croncrawlerjob.Spec.Resources,
			Args:            []string{},
			VolumeMounts:    volumeMounts,
			Command:         []string{"/bin/bash", "/app/getJob.sh"},
			ImagePullPolicy: corev1.PullIfNotPresent,
			Env:             croncrawlerjob.Spec.Env,
		},
		}
	}
	newContainers := make([]corev1.Container, len(containers))
	for index, container := range containers {
		newContainers[index] = corev1.Container{
			Name:            container.Name,
			Image:           container.Image,
			Resources:       croncrawlerjob.Spec.Resources,
			Args:            container.Args,
			VolumeMounts:    container.VolumeMounts,
			Command:         container.Command,
			ImagePullPolicy: corev1.PullIfNotPresent,
			Env:             croncrawlerjob.Spec.Env,
		}
	}
	return newContainers
}
