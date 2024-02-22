package v1alpha1_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jakobmoellerdev/splitsmart/api/v1alpha1/REST"
	"github.com/jakobmoellerdev/splitsmart/service"
	"github.com/jakobmoellerdev/splitsmart/service/mock"
	json "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAPI_IsReady(t *testing.T) {
	t.Parallel()
	_, assertions, api := SetupAPITest(t)

	req := emptyRequest(http.MethodGet)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	services := APIV1Alpha1()

	controller := gomock.NewController(t)
	mockAccounts := mock.NewMockAccounts(controller)
	services.Accounts = mockAccounts
	check := service.HealthCheck(func(ctx context.Context) (string, bool) {
		return "mock-accounts", true
	})
	mockAccounts.EXPECT().HealthCheck().Return(check).Times(1)

	testAPI(t, api, assertions, req, services.IsReady, verifyIsReady)
}

func verifyIsReady(assert *assert.Assertions, rec *httptest.ResponseRecorder) {
	assert.Equal(http.StatusOK, rec.Code)

	res := REST.HealthAggregation{}

	assert.NoError(json.NewDecoder(rec.Body).Decode(&res))
	assert.NotEmpty(res.Health)
	assert.Equal(REST.Up, res.Health)

	assert.NotNil(res.Components)
	assert.NotEmpty(res.Components)
	assert.Len(*res.Components, 1)

	for _, component := range *res.Components {
		assert.Equal(REST.Up, component.Health)
		assert.NotEmpty(component.Name)
	}
}
