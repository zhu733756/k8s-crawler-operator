// +build !ignore_autogenerated

/*
Copyright 2020 zhu733756.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1

import (
	corev1 "k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Container) DeepCopyInto(out *Container) {
	*out = *in
	if in.Command != nil {
		in, out := &in.Command, &out.Command
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Args != nil {
		in, out := &in.Args, &out.Args
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.VolumeMounts != nil {
		in, out := &in.VolumeMounts, &out.VolumeMounts
		*out = make([]corev1.VolumeMount, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Container.
func (in *Container) DeepCopy() *Container {
	if in == nil {
		return nil
	}
	out := new(Container)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CronCrawlerJob) DeepCopyInto(out *CronCrawlerJob) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CronCrawlerJob.
func (in *CronCrawlerJob) DeepCopy() *CronCrawlerJob {
	if in == nil {
		return nil
	}
	out := new(CronCrawlerJob)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CronCrawlerJob) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CronCrawlerJobConsumer) DeepCopyInto(out *CronCrawlerJobConsumer) {
	*out = *in
	if in.Parallelism != nil {
		in, out := &in.Parallelism, &out.Parallelism
		*out = new(int32)
		**out = **in
	}
	if in.ConsumerTTL != nil {
		in, out := &in.ConsumerTTL, &out.ConsumerTTL
		*out = new(int32)
		**out = **in
	}
	if in.Containers != nil {
		in, out := &in.Containers, &out.Containers
		*out = make([]Container, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.InitContainers != nil {
		in, out := &in.InitContainers, &out.InitContainers
		*out = make([]Container, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CronCrawlerJobConsumer.
func (in *CronCrawlerJobConsumer) DeepCopy() *CronCrawlerJobConsumer {
	if in == nil {
		return nil
	}
	out := new(CronCrawlerJobConsumer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CronCrawlerJobList) DeepCopyInto(out *CronCrawlerJobList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]CronCrawlerJob, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CronCrawlerJobList.
func (in *CronCrawlerJobList) DeepCopy() *CronCrawlerJobList {
	if in == nil {
		return nil
	}
	out := new(CronCrawlerJobList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CronCrawlerJobList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CronCrawlerJobPublisher) DeepCopyInto(out *CronCrawlerJobPublisher) {
	*out = *in
	if in.PublisherTTL != nil {
		in, out := &in.PublisherTTL, &out.PublisherTTL
		*out = new(int32)
		**out = **in
	}
	if in.Containers != nil {
		in, out := &in.Containers, &out.Containers
		*out = make([]Container, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CronCrawlerJobPublisher.
func (in *CronCrawlerJobPublisher) DeepCopy() *CronCrawlerJobPublisher {
	if in == nil {
		return nil
	}
	out := new(CronCrawlerJobPublisher)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CronCrawlerJobSpec) DeepCopyInto(out *CronCrawlerJobSpec) {
	*out = *in
	in.Publisher.DeepCopyInto(&out.Publisher)
	in.Consumer.DeepCopyInto(&out.Consumer)
	if in.Volumes != nil {
		in, out := &in.Volumes, &out.Volumes
		*out = make([]corev1.Volume, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = make([]corev1.EnvVar, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.Resources.DeepCopyInto(&out.Resources)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CronCrawlerJobSpec.
func (in *CronCrawlerJobSpec) DeepCopy() *CronCrawlerJobSpec {
	if in == nil {
		return nil
	}
	out := new(CronCrawlerJobSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CronCrawlerJobStatus) DeepCopyInto(out *CronCrawlerJobStatus) {
	*out = *in
	if in.PublisherActiveNodes != nil {
		in, out := &in.PublisherActiveNodes, &out.PublisherActiveNodes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ConsumerActiveNodes != nil {
		in, out := &in.ConsumerActiveNodes, &out.ConsumerActiveNodes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CronCrawlerJobStatus.
func (in *CronCrawlerJobStatus) DeepCopy() *CronCrawlerJobStatus {
	if in == nil {
		return nil
	}
	out := new(CronCrawlerJobStatus)
	in.DeepCopyInto(out)
	return out
}
