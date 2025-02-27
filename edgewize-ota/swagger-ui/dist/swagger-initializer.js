window.onload = function() {
  const protocol = window.location.protocol; // http: or https:
  const host = window.location.host; // e.g., 127.0.0.1:8080
  const apiPath = "/apidocs.json"; // API 文档的路径
  
  const apiUrl = `${protocol}//${host}${apiPath}`; // 动态拼接完整 URL

  window.ui = SwaggerUIBundle({
    url: apiUrl,
    dom_id: '#swagger-ui',
    deepLinking: true,
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ],
    plugins: [
      SwaggerUIBundle.plugins.DownloadUrl
    ],
    layout: "StandaloneLayout"
  });
};
