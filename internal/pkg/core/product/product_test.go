package product

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/maximgoltsov/botproject/internal/pkg/core/product/models"
	repoMock "github.com/maximgoltsov/botproject/internal/pkg/repository/mocks"
)

func TestUpsertProduct(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	ctx := context.Background()
	in := models.Product{
		Id:      1,
		Title:   "Title",
		Price:   100,
		Type_Id: 1,
	}

	resp := in.Id

	repo := repoMock.NewMockProduct(ctl)
	repo.EXPECT().UpsertProduct(ctx, in).Return(resp, nil).Times(1)

	useCases := New(repo)
	productId, err := useCases.UpsertProduct(ctx, in)

	require.NoError(t, err)
	require.Equal(t, resp, productId)
}

func TestGetProduct_EmptyTitle(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	ctx := context.Background()
	in := models.Product{
		Id:      1,
		Title:   "",
		Price:   100,
		Type_Id: 1,
	}

	repo := repoMock.NewMockProduct(ctl)

	useCases := New(repo)
	_, err := useCases.UpsertProduct(ctx, in)

	require.EqualError(
		t,
		err,
		errors.Wrap(ErrValidation, "field: [title] cannot be empty").Error(),
	)
}

func TestDeleteProduct(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	ctx := context.Background()
	in := uint64(1)

	repo := repoMock.NewMockProduct(ctl)
	repo.EXPECT().DeleteProductById(ctx, in).Return(nil).Times(1)

	useCases := New(repo)
	err := useCases.DeleteProductById(ctx, in)

	require.NoError(t, err)
}

func TestGetProduct(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	ctx := context.Background()
	in := uint64(1)

	resp := models.Product{
		Id:      1,
		Title:   "Title",
		Price:   100,
		Type_Id: 1,
	}

	repo := repoMock.NewMockProduct(ctl)
	repo.EXPECT().GetProduct(ctx, in).Return(resp, nil).Times(1)

	useCases := New(repo)
	product, err := useCases.GetProduct(ctx, in)

	require.NoError(t, err)
	require.Equal(t, resp, product)
}

func TestGetProducts(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	ctx := context.Background()
	fakeLimit := uint64(0)
	fakeOffset := uint64(10)
	fakeDesc := false

	resp := []models.Product{
		{
			Id:      1,
			Title:   "Title 1",
			Price:   100,
			Type_Id: 1,
		},
		{
			Id:      2,
			Title:   "Title 2",
			Price:   200,
			Type_Id: 2,
		},
	}

	repo := repoMock.NewMockProduct(ctl)
	repo.EXPECT().
		GetProducts(ctx, fakeLimit, fakeOffset, fakeDesc).
		Return(resp).Times(1)

	useCases := New(repo)
	products := useCases.GetProducts(ctx, fakeLimit, fakeOffset, fakeDesc)

	require.Equal(t, resp, products)
}
