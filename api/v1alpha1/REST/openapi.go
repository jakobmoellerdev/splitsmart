// Package REST provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.0 DO NOT EDIT.
package REST

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Defines values for HealthResult.
const (
	Down HealthResult = "Down"
	Up   HealthResult = "Up"
)

// HealthAggregation defines model for HealthAggregation.
type HealthAggregation struct {
	// Components The different Components of the Server
	Components *[]HealthAggregationComponent `json:"components,omitempty"`

	// Health A Health Check Result
	Health HealthResult `json:"health"`
}

// HealthAggregationComponent defines model for HealthAggregationComponent.
type HealthAggregationComponent struct {
	// Health A Health Check Result
	Health HealthResult `json:"health"`

	// Name The Name of the Component to be Health Checked
	Name string `json:"name"`
}

// HealthResult A Health Check Result
type HealthResult string

// ID An object ID (in the form of UUID)
type ID = openapi_types.UUID

// Registration defines model for Registration.
type Registration struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// RegistrationResult defines model for RegistrationResult.
type RegistrationResult struct {
	// Id An object ID (in the form of UUID)
	Id       *ID    `json:"id,omitempty"`
	Password string `json:"password"`
	Username string `json:"username"`
}

// RegisterJSONRequestBody defines body for Register for application/json ContentType.
type RegisterJSONRequestBody = Registration

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Checks if the Service is Available for Processing Request
	// (GET /health)
	IsHealthy(ctx echo.Context) error
	// Checks if the Service is Operational
	// (GET /ready)
	IsReady(ctx echo.Context) error
	// Register a new User
	// (POST /register)
	Register(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// IsHealthy converts echo context to params.
func (w *ServerInterfaceWrapper) IsHealthy(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.IsHealthy(ctx)
	return err
}

// IsReady converts echo context to params.
func (w *ServerInterfaceWrapper) IsReady(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.IsReady(ctx)
	return err
}

// Register converts echo context to params.
func (w *ServerInterfaceWrapper) Register(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.Register(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/health", wrapper.IsHealthy)
	router.GET(baseURL+"/ready", wrapper.IsReady)
	router.POST(baseURL+"/register", wrapper.Register)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/8xWTW8bNxD9K+y0hwTYapW4BYI9VbXceoEgFST7VPhAcUfacbjkhh9yhWD/e0FyVx+W",
	"ErtBa/QmLcg3b968meFnELpptULlLBSfwYoaGx5/XiOXrp6s1wbX3JFW4WNrdIvGEcYjx3crtMJQm47C",
	"TY2sotUKDSrHLncnmV4xVyNboNmggQzIYRPv/2BwBQV8n+9h855PfkJmBwhdBm7bIhTAjeHb8L+Op5+H",
	"OUfrpYOuy8DgJ08GKyj+HCDuduB6eY8iRvsKlxOFvoVJBoo3eF7RD7zBQcFdWOY0WyJLKOyyRvERK9gx",
	"t86QWp9kGINkTyfa0zqhMzkKyPpjGaDyTcC/bSGDqX5QB9gDlQzK6RlExVJwVk7ZK1Ixy5U2Tcj49rac",
	"voYMwn/uoADv6UySGcxxTdaZL1i25dY+aFOF3ydXvUUzaP918XYnsz3iOQkPyeyFPKZE1VP+KKcB60Wp",
	"dxlYFN6Q2y4CC+xbfEMCJz5ZOtILl5bcktjXonauhS5A4F8uBJNTLc6MiN/JXfslZOCN7K/ZIs/X5Gq/",
	"HAnd5Pf8o142GqVEU+Emt60kZxtuorikVjqNIeW4iMpiwylg9Z9+iQA/9ggBMtx73FZkGdk0lHb4bH61",
	"uGGTWckWLQpakYg1ZKSYblFNZuUIMpAkUNkoehIfJi0XNbK3o/G35ZUvpV7mDSeVvy8vrz4sruKAIycD",
	"+pP8IIMNGpsy27zhsq35m4AQSPOWoICL0Xh0EYvv6liUfD+k1nim0WdoQtNZxoeOF7Hj+zF0fXMzY3Pt",
	"XZzmwdiRSVlBAaVNN7YQPGhbrWwy0tvxeKhcPzd528o+ifzeptZN/v/HuyF57ziJhRcCrV152SdBCq1N",
	"PvdNw80WCoiTzDLabygSGLwx2XCSfCnjOGIzowMWqTWb4yePNow9x9f2cHEE5Nwgr7bPEXaOvIqMem1f",
	"PZCrWYUtqgqVILSMlJC+wur1GZXnMc7/SONdPsF7P48vXpbIb5ykN8gqb1KRBnFjgZ9b9D8Gkbn8cn3D",
	"eEcTh7q2Z2o870+EIit8YLcWDYvFPYh2UtHhFqTRjdb9qpOR/hURj1Zkd7wgnPHY/YdOOrMRn7LSIdkM",
	"fnpxMqXacEnVIyZHJhoKdlDlA89wnxwTl2p49oavJ07xipFj77XgUm7Dm86ZbfiivfvuaJkUeS7DqVpb",
	"V7wbvxvn+0l/t4v6GP9qg2Yb5l5oCMkdViFG2OWo3H599HssUu6ykydaLEu4eDBFWX7QYqVKL7QjtL5p",
	"urvu7wAAAP//dwFnwXEMAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}