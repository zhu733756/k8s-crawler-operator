/*
 * @Author: your name
 * @Date: 2020-09-25 17:29:07
 * @LastEditTime: 2020-10-15 11:32:05
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \k8s-crawler-operator\api\v1beta1\croncrawlerjob_types.go
 */
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

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CronCrawlerJobSpec defines the desired state of CronCrawlerJob
type CronCrawlerJobSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of CronCrawlerJob. Edit CronCrawlerJob_types.go to remove/update
	Publisher CronCrawlerJobPublisher     `json:"publisher"`
	Consumer  CronCrawlerJobConsumer      `json:"consumer"`
	Volumes   []corev1.Volume             `json:"volumes,omitempty"`
	Env       []corev1.EnvVar             `json:"env,omitempty"`
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	JobType   string                      `json:"jobType,omitempty"`
	Schedule  string                      `json:"schedule"`
	Duration  string                      `json:"duration"`
}

// Publisher
type CronCrawlerJobPublisher struct {
	PublisherTTL *int32      `json:"publisherTTL,omitempty"`
	Containers   []Container `json:"containers,omitempty"`
}

// Consumer
type CronCrawlerJobConsumer struct {
	Parallelism    *int32      `json:"parallelism"`
	ConsumerTTL    *int32      `json:"consumerTTL,omitempty"`
	Containers     []Container `json:"containers,omitempty"`
	InitContainers []Container `json:"initContainers,omitempty"`
}

//initContainer
type Container struct {
	Name         string               `json:"name,omitempty"`
	Image        string               `json:"image,omitempty"`
	Command      []string             `json:"command,omitempty"`
	Args         []string             `json:"args,omitempty"`
	VolumeMounts []corev1.VolumeMount `json:"volumeMounts,omitempty"`
}

// CronCrawlerJobStatus defines the observed state of CronCrawlerJob
type CronCrawlerJobStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	PublisherSucceed     int32    `json:"publisherSucceed,omitempty"`
	ConsumerSucceed      int32    `json:"consumerSucceed,omitempty"`
	PublisherActiveNodes []string `json:"publisherActiveNodes,omitempty"`
	ConsumerActiveNodes  []string `json:"consumerActiveNodes,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// CronCrawlerJob is the Schema for the croncrawlerjobs API
type CronCrawlerJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CronCrawlerJobSpec   `json:"spec,omitempty"`
	Status CronCrawlerJobStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CronCrawlerJobList contains a list of CronCrawlerJob
type CronCrawlerJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CronCrawlerJob `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CronCrawlerJob{}, &CronCrawlerJobList{})
}
