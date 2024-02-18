package v1alpha1_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jakobmoellerdev/splitsmart/api/v1alpha1/REST"
	json "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAPI_IsHealthy(t *testing.T) {
	t.Parallel()
	_, assertions, api := SetupAPITest(t)

	req := emptyRequest(http.MethodGet)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	testAPI(t, api, assertions, req, APIV1Alpha1().IsHealthy, verifyIsHealthy)
}

func verifyIsHealthy(assert *assert.Assertions, rec *httptest.ResponseRecorder) {
	assert.Equal(http.StatusOK, rec.Code)

	res := REST.HealthAggregation{}

	assert.NoError(json.NewDecoder(rec.Body).Decode(&res))
	assert.NotEmpty(res.Health)
	assert.Equal(REST.Up, res.Health)
	assert.Empty(res.Components)
}
