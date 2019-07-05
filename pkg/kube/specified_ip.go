package kube

import (
	"fmt"

	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/selection"
)

func judgeIpSpecDecreasing(target, existence runtime.Object) (bool, error) {
	// last operation not complete
	_, decreasing := getRuntimeObjectAnnotationValue(existence, AnnoKeySpecifiedIPsDecreasing)
	if decreasing {
		return true, nil
	}
	// ip list diff
	targetIPsRaw, _ := getRuntimeObjectAnnotationValue(target, AnnoKeySpecifiedIPs)
	existenceIPsRaw, _ := getRuntimeObjectAnnotationValue(existence, AnnoKeySpecifiedIPs)
	// string same
	if targetIPsRaw == existenceIPsRaw {
		return false, nil
	}
	// parse
	targetIPs, e := parseSpecifiedIPSetsFromString(targetIPsRaw)
	if e != nil { // do not set invalid ip setting
		return false, e
	}
	existenceIPs, _ := parseSpecifiedIPSetsFromString(existenceIPsRaw)
	if len(existenceIPs) == 0 { // no prev ip setting
		return false, nil
	}
	// check
	m := ipSets2AddressMap(targetIPs)
	for i := range existenceIPs {
		for _, ip := range existenceIPs[i].IPs {
			if _, ok := m[ip]; !ok {
				return true, nil
			}
		}
	}
	return false, nil
}

func (c *client) applyIpSpecDecreasing(client *ResourceClient,
	gvk schema.GroupVersionKind, obj, existence runtime.Object) (err error) {
	// parse and check ips
	targetAnnoValue, _ := getRuntimeObjectAnnotationValue(obj, AnnoKeySpecifiedIPs)
	ipSets, err := parseSpecifiedIPSetsFromString(targetAnnoValue)
	if err != nil {
		return
	}
	// set deployment/statefulset anno
	tempObj := existence.DeepCopyObject()
	setRuntimeObjectAnnotationValue(tempObj, AnnoKeySpecifiedIPs, targetAnnoValue)
	setRuntimeObjectAnnotationValue(tempObj, AnnoKeySpecifiedIPsDecreasing, "") // rm on final update
	// do update
	if _, err = client.Update(tempObj); err != nil {
		return
	}
	// delete pods on deleted ips
	namespace := getRuntimeObjectNamespace(obj)
	podClient, err := c.pool.ClientFor(corev1.SchemeGroupVersion.WithKind("Pod"), namespace)
	if err != nil {
		return
	}
	releaseName, _ := getRuntimeObjectLabelValue(obj, LabelsKeyRelease)
	if len(releaseName) == 0 {
		err = fmt.Errorf("get release name empty")
		return
	}
	delList, err := deleteReleasePodsNotInList(podClient, releaseName, ipSets2AddressMap(ipSets))
	for _, podName := range delList {
		glog.Infof("applyIpSpecDecreasing[namespace:%v][release:%s][pod:%] deleted",
			namespace, releaseName, podName)
	}
	return err
}

func deleteReleasePodsNotInList(podClient *ResourceClient,
	releaseName string, ipMap map[string]struct{}) (delList []string, err error) {
	// list pod of release
	labelReq, err := labels.NewRequirement(LabelsKeyRelease, selection.Equals, []string{releaseName})
	if err != nil {
		return
	}
	labelSelector := labels.NewSelector().Add(*labelReq)
	list, err := podClient.List(metav1.ListOptions{LabelSelector: labelSelector.String()})
	if err != nil {
		return
	}
	podList := list.(*corev1.PodList)
	if podList == nil {
		return
	}
	// do delete
	for i := range podList.Items {
		pod := &podList.Items[i]
		podIP := pod.Status.PodIP
		if len(podIP) == 0 {
			continue
		}
		if _, ok := ipMap[podIP]; ok {
			continue
		}
		// del not in list
		if err = podClient.Delete(pod.Name, nil); err != nil && !errors.IsNotFound(err) {
			return
		}
		delList = append(delList, pod.Name)
	}
	return
}
