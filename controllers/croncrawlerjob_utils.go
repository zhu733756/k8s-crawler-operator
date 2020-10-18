/*
 * @Author: your name
 * @Date: 2020-10-15 09:35:38
 * @LastEditTime: 2020-10-15 11:47:36
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \undefinede:\go\src\github.com\zhu733756\k8s-crawler-operator\croncrawlerjob_utils.go
 */
package controllers

import (
	cronjobv1 "github.com/zhu733756/k8s-crawler-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
)

// job names
func NameForCronCrawlerJob(name string, role string, jobType string) string {
	if role == "publisher" {
		return name + "-" + jobType + "-publisher"
	}
	return name + "-" + jobType + "-consumers"
}

// labels
func LabelsForCronCrawlerJob(name string, role string, jobType string) map[string]string {
	if role == "publisher" {
		return map[string]string{
			"crawler-table": name,
			"role":          "crawler-publisher",
			"jobType":       jobType,
		}
	}
	return map[string]string{
		"crawler-table": name,
		"role":          "crawler-consumers",
		"jobType":       jobType,
	}
}

func getPodStatus(pods []corev1.Pod) *cronjobv1.CronCrawlerJobStatus {
	var publisherSucceed int32
	var consumerSucceed int32
	var publisherActiveNodes []string
	var consumerActiveNodes []string
	for _, pod := range pods {
		role, _ := pod.Labels["role"]
		if role == "crawler-publisher" {
			if pod.Status.Phase == "Running" {
				publisherActiveNodes = append(publisherActiveNodes, pod.Name)
			} else if pod.Status.Phase == "Succeeded" {
				publisherSucceed++
			}
		}
		if role == "crawler-consumers" {
			if pod.Status.Phase == "Running" {
				consumerActiveNodes = append(consumerActiveNodes, pod.Name)
			} else if pod.Status.Phase == "Succeeded" {
				consumerSucceed++
			}
		}
	}
	return &cronjobv1.CronCrawlerJobStatus{
		PublisherSucceed:     publisherSucceed,
		ConsumerSucceed:      consumerSucceed,
		PublisherActiveNodes: publisherActiveNodes,
		ConsumerActiveNodes:  consumerActiveNodes,
	}
}
