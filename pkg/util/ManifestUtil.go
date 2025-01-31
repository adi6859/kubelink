package util

import (
	"github.com/devtron-labs/common-lib/utils/k8s/commonBean"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func GetPorts(manifest *unstructured.Unstructured, gvk schema.GroupVersionKind) []int64 {
	ports := make([]int64, 0)
	if gvk.Kind == commonBean.ServiceKind {
		ports = getPortsFromService(manifest, ports)
	}
	if gvk.Kind == commonBean.EndPointsSlice {
		ports = getPortsFromEndPointsSlice(manifest, ports)
	}
	if gvk.Kind == commonBean.EndpointsKind {
		ports = getPortsFromEndpointsKind(manifest, ports)
	}
	return ports
}

func getPortsFromService(manifest *unstructured.Unstructured, ports []int64) []int64 {
	if manifest.Object["spec"] != nil {
		spec := manifest.Object["spec"].(map[string]interface{})
		if spec["ports"] != nil {
			portList := spec["ports"].([]interface{})
			for _, portItem := range portList {
				if portItem.(map[string]interface{}) != nil {
					_portNumber := portItem.(map[string]interface{})["port"]
					portNumber := _portNumber.(int64)
					if portNumber != 0 {
						ports = append(ports, portNumber)
					}
				}
			}
		}
	}
	return ports
}

func getPortsFromEndPointsSlice(manifest *unstructured.Unstructured, ports []int64) []int64 {
	if manifest.Object["ports"] != nil {
		endPointsSlicePorts := manifest.Object["ports"].([]interface{})
		for _, val := range endPointsSlicePorts {
			_portNumber := val.(map[string]interface{})["port"]
			portNumber := _portNumber.(int64)
			if portNumber != 0 {
				ports = append(ports, portNumber)
			}
		}
	}
	return ports
}

func getPortsFromEndpointsKind(manifest *unstructured.Unstructured, ports []int64) []int64 {
	if manifest.Object["subsets"] != nil {
		subsets := manifest.Object["subsets"].([]interface{})
		for _, subset := range subsets {
			subsetObj := subset.(map[string]interface{})
			if subsetObj != nil {
				portsIfs := subsetObj["ports"].([]interface{})
				for _, portsIf := range portsIfs {
					portsIfObj := portsIf.(map[string]interface{})
					if portsIfObj != nil {
						port := portsIfObj["port"].(int64)
						ports = append(ports, port)
					}
				}
			}
		}
	}
	return ports
}
