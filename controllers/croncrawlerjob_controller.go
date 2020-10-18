/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"encoding/json"
	"reflect"
	"time"

	"github.com/go-logr/logr"
	cronjobv1 "github.com/zhu733756/k8s-crawler-operator/api/v1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CronCrawlerJobReconciler ctrls a CronCrawlerJob object
type CronCrawlerJobReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// updated data
type CrawlerUpdatedData struct {
	JobType      string          `json:"jobType"`
	Schedule     string          `json:"schedule,omitempty"`
	Duration     string          `json:"duration,omitempty"`
	PublisherTTL *int32          `json:"publisherTTL,omitempty"`
	ConsumerTTL  *int32          `json:"consumerTTL,omitempty"`
	Parallelism  *int32          `json:"parallelism"`
	Env          []corev1.EnvVar `json:"env,omitempty"`
}

// +kubebuilder:rbac:groups=cronjob.crawler.com,resources=croncrawlerjobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cronjob.crawler.com,resources=croncrawlerjobs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=batch,resources=cronjobs,verbs=get;list;watch;create;update;patch;delete;deletecollection
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete;deletecollection
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch

func (r *CronCrawlerJobReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("croncrawlerjob", req.NamespacedName)
	instance := &cronjobv1.CronCrawlerJob{}

	// 1 判断自定义资源是否存在
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after ctrl request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get CronCrawlerJob")
		return ctrl.Result{}, err
	}

	if instance.DeletionTimestamp != nil {
		return ctrl.Result{}, err
	}

	jobType := instance.Spec.JobType
	schedule := instance.Spec.Schedule

	// 采集类型为定时任务时，必须指定采集频率
	if (jobType == "cronjob") && (schedule == "") {
		log.Error(err, "`schedule` must be confirmed in a cronjob jobtype！")
		return ctrl.Result{}, err
	}

	// 2 自定义资源已经存在，判断新旧任务类型是否一致，
	// 2.1 如果不一致一律删除，注意是普通任务也直接删除，
	// 2.2 如果一致，也就是定时任务了，直接更新

	newSpec := CrawlerUpdatedData{
		JobType:      instance.Spec.JobType,
		Schedule:     instance.Spec.Schedule,
		Env:          instance.Spec.Env,
		Duration:     instance.Spec.Duration,
		PublisherTTL: instance.Spec.Publisher.PublisherTTL,
		ConsumerTTL:  instance.Spec.Consumer.ConsumerTTL,
		Parallelism:  instance.Spec.Consumer.Parallelism,
	}

	duration, _ := time.ParseDuration(instance.Spec.Duration)
	ns := instance.Namespace

	log.Info("instance ns: " + ns + ", req ns:  " + req.NamespacedName.Namespace)
	name := instance.Name
	shortDes := ns + "/" + name
	publisherName := NameForCronCrawlerJob(name, "publisher", jobType)
	consumerName := NameForCronCrawlerJob(name, "consumer", jobType)

	// 判断是不是旧任务入队
	if (instance.Annotations != nil) && (instance.Annotations["status"] == "old") {
		// 反序列化
		oldSpec := CrawlerUpdatedData{}
		err = json.Unmarshal([]byte(instance.Annotations["spec"]), &oldSpec)
		if err != nil {
			return ctrl.Result{}, nil
		}
		// 判断是否有更新
		if !reflect.DeepEqual(oldSpec, newSpec) {
			// job type
			oldJobType := oldSpec.JobType
			if (jobType != oldJobType) || (oldJobType == "job") {
				// deletes the older resource
				if err := r.DeleteOldJobsByLabels(ns, name, oldJobType); err != nil {
					log.Error(err, "Older "+oldJobType+" for "+shortDes+" deletes failed")
					return ctrl.Result{}, err
				}
				log.Info("Older " + oldJobType + " for " + shortDes + " deletes success")
				// resets the current resource
				instance.Annotations = nil
				if err := r.Update(ctx, instance); err != nil {
					return ctrl.Result{}, err
				}
				log.Info("Current " + jobType + "'s annotation sets to nil, would be requeued later")
				return ctrl.Result{Requeue: true, RequeueAfter: duration}, nil
			} else if jobType == "cronjob" {
				// 更新cronjob，并返回
				if err := r.UpdateOldJobs(instance, ns, jobType, publisherName, consumerName); err != nil {
					log.Error(err, "Older "+jobType+" for "+shortDes+" updates failed")
					return ctrl.Result{}, err
				}
				log.Info("Older " + oldJobType + " for " + shortDes + " updates success")
				return ctrl.Result{}, err

			}

		}
	}

	// 3、使用当前instance创建新资源并关联数据返回，只有当获取不到资源时候，才创建资源

	var jobPublisher runtime.Object
	var jobConsumers runtime.Object

	if jobType == "job" {
		jobPublisher = r.NewPublisherJob(instance, publisherName)
		jobConsumers = r.NewConsumersJob(instance, consumerName)
	} else if jobType == "cronjob" {
		jobPublisher = r.NewPublisherCronJob(instance, publisherName)
		jobConsumers = r.NewConsumersCronJob(instance, consumerName)
	} else {
		log.Error(err, "JobType sets wrong")
		return ctrl.Result{}, nil
	}

	if err := r.Get(ctx, types.NamespacedName{Name: publisherName, Namespace: ns}, jobPublisher); err != nil {
		log.Info("Starting to create " + jobType + " for publisher")
		if err := r.Create(ctx, jobPublisher); err != nil {
			log.Error(err, jobType+" for publisher creates failed")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil

	}
	if err := r.Get(ctx, types.NamespacedName{Name: consumerName, Namespace: ns}, jobConsumers); err != nil {
		log.Info("Starting to create " + jobType + " for consumers")
		if err := r.Create(ctx, jobConsumers); err != nil {
			log.Error(err, jobType+" for Consumers creates failed")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// 关联 Annotations
	data, _ := json.Marshal(&newSpec)
	if instance.Annotations != nil {
		instance.Annotations["spec"] = string(data)
		instance.Annotations["status"] = "old"
	} else {
		instance.Annotations = map[string]string{"spec": string(data), "status": "new"}
	}

	if err := r.Update(ctx, instance); err != nil {
		return ctrl.Result{}, nil
	}

	// 更新显示状态
	podList := &corev1.PodList{}
	labels := map[string]string{"crawler-table": name, "jobType": jobType}
	listOption := []client.ListOption{
		client.InNamespace(ns),
		client.MatchingLabels(labels),
	}
	if err := r.List(ctx, podList, listOption...); err != nil {
		return ctrl.Result{}, nil
	}

	cronCrawlerJobStatus := getPodStatus(podList.Items)
	if !reflect.DeepEqual(cronCrawlerJobStatus, instance.Status) {
		log.Info("Starting to update status")
		instance.Status = *cronCrawlerJobStatus
		err := r.Status().Update(ctx, instance)
		if err != nil {
			log.Error(err, "Failed to update CronCrawlerJob status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// delete jobs/conjob list
func (r *CronCrawlerJobReconciler) DeleteOldJobsByLabels(namespace string, name string, jobType string) error {
	ctx := context.Background()
	// gp := int64(duration.Seconds())
	dp := metav1.DeletePropagationForeground
	var deleteObj runtime.Object
	switch jobType {
	case "job":
		deleteObj = &batchv1.Job{}
	case "cronjob":
		deleteObj = &batchv1beta1.CronJob{}
	}

	deleteAllOfOption := []client.DeleteAllOfOption{
		// client.GracePeriodSeconds(gp),
		client.PropagationPolicy(dp),
		client.InNamespace(namespace),
		client.MatchingLabels(map[string]string{"crawler-table": name, "jobType": jobType}),
	}
	if err := r.DeleteAllOf(ctx, deleteObj, deleteAllOfOption...); err != nil {
		return err
	}
	return nil
}

//updates old job/cronjob
func (r *CronCrawlerJobReconciler) UpdateOldJobs(instance *cronjobv1.CronCrawlerJob, ns string, jobType string, publisherName string, consumerName string) error {
	ctx := context.Background()
	switch jobType {
	case "job":
		oldRuntimeObj := &batchv1.Job{}
		if err := r.Get(ctx, types.NamespacedName{Namespace: ns, Name: publisherName}, oldRuntimeObj); err != nil {
			return err
		}
		jobPublisher := r.NewPublisherJob(instance, publisherName)
		oldRuntimeObj.Spec = jobPublisher.Spec
		if err := r.Update(ctx, oldRuntimeObj); err != nil {
			return err
		}
		if err := r.Get(ctx, types.NamespacedName{Namespace: ns, Name: consumerName}, oldRuntimeObj); err != nil {
			return err
		}
		jobConsumers := r.NewConsumersJob(instance, consumerName)
		oldRuntimeObj.Spec = jobConsumers.Spec
		if err := r.Update(ctx, oldRuntimeObj); err != nil {
			return err
		}
	case "cronjob":
		oldRuntimeObj := &batchv1beta1.CronJob{}
		if err := r.Get(ctx, types.NamespacedName{Namespace: ns, Name: publisherName}, oldRuntimeObj); err != nil {
			return err
		}
		jobPublisher := r.NewPublisherCronJob(instance, publisherName)
		oldRuntimeObj.Spec = jobPublisher.Spec
		if err := r.Update(ctx, oldRuntimeObj); err != nil {
			return err
		}
		if err := r.Get(ctx, types.NamespacedName{Namespace: ns, Name: consumerName}, oldRuntimeObj); err != nil {
			return err
		}
		jobConsumers := r.NewConsumersCronJob(instance, consumerName)
		oldRuntimeObj.Spec = jobConsumers.Spec
		if err := r.Update(ctx, oldRuntimeObj); err != nil {
			return err
		}
	}
	return nil
}

func (r *CronCrawlerJobReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).For(&cronjobv1.CronCrawlerJob{}).Complete(r)
}
