package pg

import (
	"context"
	"strings"
	"testing"

	"techno-store/internal/domain/algo"
	"techno-store/internal/domain/bo"

	"github.com/stretchr/testify/require"
)

// Return brand with happy case
func createGoodRandomBrand(t *testing.T) *bo.Brand {
	brandCreate := &bo.Brand{
		Name:     algo.GenerateRandomString(10),
		StatusID: 1,
	}

	err := testStore.Brand.CreateBrand(context.Background(), brandCreate)
	require.NoError(t, err)
	require.NotEmpty(t, brandCreate.ID)

	return brandCreate
}

func TestAddBrand(t *testing.T) {
	createGoodRandomBrand(t)
	brandWithoutName := &bo.Brand{
		StatusID: 1,
	}

	err := testStore.Brand.CreateBrand(context.Background(), brandWithoutName)
	require.Error(t, err)
	require.Empty(t, brandWithoutName.ID)

	brandEmpty := &bo.Brand{}
	err = testStore.Brand.CreateBrand(context.Background(), brandEmpty)
	require.Error(t, err)
	require.Empty(t, brandEmpty.ID)
}

func TestGetBrandByID(t *testing.T) {
	brandCreate := createGoodRandomBrand(t)

	brandGet, err := testStore.Brand.GetBrandByID(context.Background(), brandCreate.ID)
	require.NoError(t, err)
	require.NotEmpty(t, brandGet)

	require.Equal(t, brandCreate.ID, brandGet.ID)
	require.Equal(t, strings.ToLower(brandCreate.Name), brandGet.Name)

	brandEmpty, err := testStore.Brand.GetBrandByID(context.Background(), 0)
	require.Error(t, err)
	require.Equal(t, bo.ErrBrandNotFound, err)
	require.Empty(t, brandEmpty)
}

func TestQueryBrand(t *testing.T) {
	q := algo.GenerateRandomString(10)
	brand1 := &bo.Brand{
		Name:     q,
		StatusID: 1,
	}
	err := testStore.Brand.CreateBrand(context.Background(), brand1)
	require.NoError(t, err)

	createGoodRandomBrand(t)
	createGoodRandomBrand(t)
	createGoodRandomBrand(t)

	happyQuery := bo.BrandQuery{
		Limit:  10,
		Offset: 0,
	}

	paginatedBrandCollection, err := testStore.Brand.ListBrands(context.Background(), happyQuery)
	require.NoError(t, err)
	require.NotEmpty(t, paginatedBrandCollection)
	require.NotEmpty(t, paginatedBrandCollection.Data)
}

func TestUpdateBrand(t *testing.T) {
	brandCreate := createGoodRandomBrand(t)

	updatedName := brandCreate.Name + "updated"
	brandUpdate := bo.BrandUpdate{
		ID:   brandCreate.ID,
		Name: &updatedName,
	}

	err := testStore.Brand.UpdateBrand(context.Background(), brandUpdate)
	require.NoError(t, err)

	brandGet, err := testStore.Brand.GetBrandByID(context.Background(), brandCreate.ID)
	require.NoError(t, err)
	require.NotEmpty(t, brandGet)

	require.Equal(t, brandCreate.ID, brandGet.ID)
	require.Equal(t, updatedName, brandGet.Name)

	brandUpdateEmpty := bo.BrandUpdate{}
	err = testStore.Brand.UpdateBrand(context.Background(), brandUpdateEmpty)
	require.Error(t, err)
}

func TestDeleteBrand(t *testing.T) {
	brandCreate := createGoodRandomBrand(t)

	err := testStore.Brand.DeleteBrand(context.Background(), brandCreate.ID)
	require.NoError(t, err)

	brandGet, err := testStore.Brand.GetBrandByID(context.Background(), brandCreate.ID)
	require.Error(t, err)
	require.Empty(t, brandGet)
	require.Equal(t, bo.ErrBrandNotFound, err)
}
