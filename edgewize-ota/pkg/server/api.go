package server

import (
	"net/http"

	"github.com/edgewize-io/edgewize/edgewize-ota/pkg/types"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
)

func NodeWebServiceApi(ws *restful.WebService) {
	// node register
	ws.Route(ws.POST("/register").To(installAgentHandler).
		Doc("install agent").Metadata(restfulspec.KeyOpenAPITags, []string{"nodes"}).Operation("installAgent").
		Reads(EdgeInstallAgent{}).
		Returns(http.StatusOK, "OK", nil))

	ws.Route(ws.GET("/nodes").To(listNodes).
		Doc("get nodes status").Metadata(restfulspec.KeyOpenAPITags, []string{"nodes"}).Operation("getNodesStatus").
		Param(ws.QueryParameter("nodeName", "node name").DataType("string")).
		Param(ws.QueryParameter("status", "status").DataType("string")).
		Param(ws.QueryParameter("sortBy", "sort order by").DataType("string")).
		Param(ws.QueryParameter("limit", "limit").DataType("int")).
		Param(ws.QueryParameter("page", "page").DataType("int")).
		Returns(http.StatusOK, "OK", ListResult{}))

	ws.Route(ws.GET("/nodes/{nodeName}").To(getNode).
		Doc("get node").Metadata(restfulspec.KeyOpenAPITags, []string{"nodes"}).Operation("getNode").
		Param(ws.PathParameter("nodeName", "node name").DataType("string")).
		Returns(http.StatusOK, "OK", EdgeNodeInfo{}))

	ws.Route(ws.DELETE("/nodes/{nodeName}").To(delNode).
		Doc("delete node").Metadata(restfulspec.KeyOpenAPITags, []string{"nodes"}).Operation("deleteNode").
		Param(ws.PathParameter("nodeName", "node name").DataType("string")).
		Returns(http.StatusOK, "OK", nil))
}

func TaskWebServiceApi(ws *restful.WebService) {
	ws.Route(ws.POST("/tasks").To(doTask).
		Doc("do task").Metadata(restfulspec.KeyOpenAPITags, []string{"tasks"}).Operation("doTask").
		Reads(types.TaskRequest{}).
		Returns(http.StatusOK, "OK", nil))

	ws.Route(ws.GET("/tasks").To(getTask).
		Doc("get task list").Metadata(restfulspec.KeyOpenAPITags, []string{"tasks"}).Operation("getTask").
		Returns(http.StatusOK, "OK", nil))

	ws.Route(ws.GET("/tasks/{taskName}").To(getTaskConfig).
		Doc("get task config").Metadata(restfulspec.KeyOpenAPITags, []string{"tasks"}).Operation("getTaskConfig").
		Param(ws.PathParameter("taskName", "task name").DataType("string")).
		Returns(http.StatusOK, "OK", nil))
}
