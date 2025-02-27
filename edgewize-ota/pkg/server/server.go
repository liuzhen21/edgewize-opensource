package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	restful "github.com/emicklei/go-restful/v3"
	"github.com/go-openapi/spec"
	"k8s.io/klog/v2"

	nodeinfodb "github.com/edgewize-io/edgewize/edgewize-ota/pkg/nodeinfodb"
	worker "github.com/edgewize-io/edgewize/edgewize-ota/pkg/worker"
)

const MIME_MERGEPATCH = "application/merge-patch+json"

var Container = restful.DefaultContainer
var tm *worker.TaskManager

// 提供 k8s 的 api 接口
func OtaWebServiceApi() *restful.WebService {
	restful.RegisterEntityAccessor(MIME_MERGEPATCH, restful.NewEntityAccessorJSON(restful.MIME_JSON))

	ws := new(restful.WebService)
	ws.Path("/apis/ota.edgewize.io/v1alpha1").Consumes(restful.MIME_JSON, MIME_MERGEPATCH).Produces(restful.MIME_JSON)

	// 获取资源列表, 用于 k8s apiservice 的注册
	ws.Route(ws.GET("/").To(getResourceList).
		Doc("get resource list").Metadata(restfulspec.KeyOpenAPITags, []string{"k8s-apiservice"}).Operation("getResourceList").
		Returns(http.StatusOK, "OK", ResourceList{}))

	// 健康检查
	ws.Route(ws.GET("/health").To(healthCheckHandler).
		Doc("health check").Metadata(restfulspec.KeyOpenAPITags, []string{"k8s-apiservice"}).Operation("healthCheck").
		Returns(http.StatusOK, "OK", nil))

	NodeWebServiceApi(ws)
	TaskWebServiceApi(ws)
	return ws
}

func Run() {
	nodeInfoDB = nodeinfodb.DefaultNodeInfoDB()
	tm = worker.GetDefaultTaskManager(nodeInfoDB)

	Container.Add(OtaWebServiceApi())

	config := restfulspec.Config{
		WebServices:                   restful.RegisteredWebServices(), // you control what services are visible
		APIPath:                       "/apidocs.json",
		PostBuildSwaggerObjectHandler: enrichSwaggerObject}
	Container.Add(restfulspec.NewOpenAPIService(config))

	// Optionally, you can install the Swagger Service which provides a nice Web UI on your REST API
	// You need to download the Swagger HTML5 assets and change the FilePath location in the config below.
	// Open http://localhost:8082/apidocs/?url=http://localhost:8082/apidocs.json
	cwdDir, _ := os.Getwd()
	http.Handle("/apidocs/", http.StripPrefix("/apidocs/", http.FileServer(http.Dir(filepath.Join(cwdDir, "swagger-ui/dist")))))

	enableCORS()

	klog.Infof("Edge ota service running\n")

	klog.Infof("Edge ota service serving %+v\n", http.ListenAndServe(fmt.Sprintf(":%d", 8082), nil))
	// 暴露https服务
	//klog.Infof("Edge ota service serving %+v\n", http.ListenAndServeTLS(fmt.Sprintf(":%d", 8443), "/etc/edgewize-ota-server/certs/tls.crt", "/etc/edgewize-ota-server/certs/tls.key", nil))

}

func enableCORS() {
	// Optionally, you may need to enable CORS for the UI to work.
	cors := restful.CrossOriginResourceSharing{
		AllowedHeaders: []string{"Content-Type", "Accept"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		CookiesAllowed: false,
		AllowedDomains: []string{"*"},
		Container:      Container}
	Container.Filter(cors.Filter)
}

// 待处理
func enrichSwaggerObject(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "Edgewize",
			Description: "Resource for managing Users",
			Contact: &spec.ContactInfo{
				ContactInfoProps: spec.ContactInfoProps{
					Name:  "Edgewize",
					Email: "edgewize@yunify.com",
					URL:   "https://kubesphere.io",
				},
			},
			License: &spec.License{
				LicenseProps: spec.LicenseProps{
					Name: "MIT",
					URL:  "http://mit.org",
				},
			},
			Version: "1.0.0",
		},
	}
	swo.Tags = []spec.Tag{{TagProps: spec.TagProps{
		Name:        "users",
		Description: "Managing users"}}}
}
