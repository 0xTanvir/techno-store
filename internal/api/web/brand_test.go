package web

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"techno-store/config"
	"techno-store/internal/domain/algo"
	"techno-store/internal/domain/bo"
	"techno-store/internal/infrastructure/datastores/mockdb"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetBrandAPI(t *testing.T) {
	brand := randomBrand()

	// TODO: Add multiple test cases
	// for each error scenarios
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	appConfig, err := config.Parse()
	if err != nil {
		slog.Error("Error parsing config", "cause", err)
	}

	// Get a mockdb datastore instance
	ds := mockdb.GetInstance(ctrl)
	apiService := NewAPIService(*appConfig.Server, ds)

	brandStore := ds.Brand.(*mockdb.MockBrandRepository)
	// build stubs
	brandStore.EXPECT().
		GetBrandByID(gomock.Any(), gomock.Eq(brand.ID)).
		Times(1).
		Return(brand, nil)

	recorder := httptest.NewRecorder()
	url := fmt.Sprintf("/v1/brand/%d", brand.ID)
	req, err := http.NewRequest("GET", url, nil)
	require.NoError(t, err)

	router := gin.Default()
	apiService.InstallRoutes(router)

	router.ServeHTTP(recorder, req)
	// Check response
	require.Equal(t, http.StatusOK, recorder.Code)
}

func randomBrand() bo.Brand {
	return bo.Brand{
		ID:       int64(algo.GenerateRandomInteger(1, 1000)),
		Name:     algo.GenerateRandomString(10),
		StatusID: 1,
	}
}
