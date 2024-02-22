package v1alpha1_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jakobmoellerdev/splitsmart/api/v1alpha1"
	"github.com/jakobmoellerdev/splitsmart/api/v1alpha1/REST"
	"github.com/labstack/echo/v4"
)

func TestOpenAPIHandler_ServeOpenAPI(t *testing.T) {
	t.Parallel()
	log, assertions, api := SetupAPITest(t)
	swagger, err := REST.GetSwagger()

	assertions.NoError(err)

	handler := v1alpha1.NewOpenAPIHandler(log.WithContext(context.Background()), swagger)
	req := emptyRequest(http.MethodGet)
	rec := httptest.NewRecorder()
	ctx := api.NewContext(req, rec)

	for _, contentType := range []string{
		echo.MIMEApplicationJSONCharsetUTF8,
		echo.MIMEApplicationJSON,
		"some-random-content",
		"text/yaml",
	} {
		ctx.Reset(req, rec)
		req.Header.Set(echo.HeaderContentType, contentType)
		assertions.NoError(handler.ServeOpenAPI(ctx))
		assertions.NotEmpty(rec.Body)
	}
}
