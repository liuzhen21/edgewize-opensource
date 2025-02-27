package server

import (
	"github.com/emicklei/go-restful/v3"
	"k8s.io/klog/v2"
)

func healthCheckHandler(request *restful.Request, response *restful.Response) {
	response.WriteEntity("OK")
}

func getResourceList(request *restful.Request, response *restful.Response) {
	klog.Infof("get resource list")
	resourceList := ResourceList{
		Kind:         "APIResourceList",
		ApiVersion:   "v1",
		GroupVersion: "ota.edgewize.io/v1alpha1",
		Resources:    []APIResource{},
	}

	resourceList.Resources = append(resourceList.Resources, APIResource{
		Name:         "nodes",
		SingularName: "node",
		Namespaced:   false,
		Kind:         "Node",
		Verbs:        []string{"get", "list", "delete"},
	})

	resourceList.Resources = append(resourceList.Resources, APIResource{
		Name:         "nodes/tags",
		SingularName: "",
		Namespaced:   false,
		Kind:         "Node",
		Verbs:        []string{"update", "patch", "get"},
	})

	resourceList.Resources = append(resourceList.Resources, APIResource{
		Name:         "registers",
		SingularName: "register",
		Namespaced:   false,
		Kind:         "Register",
		Verbs:        []string{"create"},
	})

	response.WriteEntity(resourceList)
}
