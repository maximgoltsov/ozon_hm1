package api

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	productPkg "github.com/maximgoltsov/botproject/internal/pkg/core/product"
	"github.com/maximgoltsov/botproject/internal/pkg/core/product/models"
	fakeProduct "github.com/maximgoltsov/botproject/internal/pkg/repository/fake"
	repoMock "github.com/maximgoltsov/botproject/internal/pkg/repository/mocks"
	pb "github.com/maximgoltsov/botproject/pkg/api"
)

func TestProduct(t *testing.T) {
	ctx := context.Background()

	repo := fakeProduct.New()
	api := New(repo)
	// Create and check
	inCreate := pb.ProductCreateRequest{
		Title:  "Title",
		Price:  100,
		TypeId: 1,
	}

	createResp, err := api.ProductCreate(ctx, &inCreate)

	require.NoError(t, err)
	require.NotEqual(t, createResp.Id, uint64(0))

	inGet := pb.ProductGetRequest{
		Id: createResp.GetId(),
	}
	getResp, err := api.ProductGet(ctx, &inGet)

	require.NoError(t, err)
	require.Equal(t, getResp.GetTitle(), inCreate.GetTitle())
	require.Equal(t, getResp.GetPrice(), inCreate.GetPrice())
	require.Equal(t, getResp.GetTypeId(), inCreate.GetTypeId())

	// Update and check
	inUpdate := pb.ProductUpdateRequest{
		Id:     createResp.GetId(),
		Title:  "Rose",
		Price:  333,
		TypeId: 1,
	}
	updateResp, err := api.ProductUpdate(ctx, &inUpdate)

	require.NoError(t, err)
	require.Equal(t, updateResp.GetId(), createResp.GetId())

	inGet = pb.ProductGetRequest{
		Id: createResp.GetId(),
	}
	getResp, err = api.ProductGet(ctx, &inGet)

	require.NoError(t, err)
	require.Equal(t, getResp.GetTitle(), inUpdate.GetTitle())
	require.Equal(t, getResp.GetPrice(), inUpdate.GetPrice())
	require.Equal(t, getResp.GetTypeId(), inUpdate.GetTypeId())

	// Delete and check
	inDelete := pb.ProductDeleteRequest{
		Id: createResp.GetId(),
	}
	_, err = api.ProductDelete(ctx, &inDelete)

	require.NoError(t, err)

	inGet = pb.ProductGetRequest{
		Id: createResp.GetId(),
	}
	_, err = api.ProductGet(ctx, &inGet)

	require.Error(t, err)
}

func TestProductCreateSuccess(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	ctx := context.Background()
	in := pb.ProductCreateRequest{
		Title:  "Title",
		Price:  100,
		TypeId: 1,
	}

	repo := repoMock.NewMockProduct(ctl)
	repo.EXPECT().UpsertProduct(ctx, models.Product{
		Id:      0,
		Title:   in.GetTitle(),
		Price:   in.GetPrice(),
		Type_Id: in.GetTypeId(),
	}).Return(uint64(1), nil).Times(1)

	api := New(repo)

	resp, err := api.ProductCreate(ctx, &in)

	require.NoError(t, err)
	require.NotEqual(t, resp.Id, uint64(0))
}

func TestProductCreateError(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	ctx := context.Background()
	in := pb.ProductCreateRequest{
		Title:  "Title",
		Price:  100,
		TypeId: 1,
	}

	repoError := errors.New("some repo error")
	repo := repoMock.NewMockProduct(ctl)
	repo.EXPECT().UpsertProduct(ctx, models.Product{
		Id:      0,
		Title:   in.GetTitle(),
		Price:   in.GetPrice(),
		Type_Id: in.GetTypeId(),
	}).Return(uint64(0), repoError).Times(1)

	api := New(repo)

	_, err := api.ProductCreate(ctx, &in)

	require.EqualError(t, err, "rpc error: code = Internal desc = some repo error")
}

func TestProductCreateErrorValidation(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	ctx := context.Background()
	in := pb.ProductCreateRequest{
		Title:  "Title",
		Price:  100,
		TypeId: 1,
	}

	repo := repoMock.NewMockProduct(ctl)
	repo.EXPECT().UpsertProduct(ctx, models.Product{
		Id:      0,
		Title:   in.GetTitle(),
		Price:   in.GetPrice(),
		Type_Id: in.GetTypeId(),
	}).Return(uint64(0), productPkg.ErrValidation).Times(1)

	api := New(repo)

	_, err := api.ProductCreate(ctx, &in)

	require.EqualError(t, err, "rpc error: code = InvalidArgument desc = invalid data")
}

func TestProductDeleteSuccess(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	ctx := context.Background()
	in := pb.ProductDeleteRequest{
		Id: 1,
	}

	repo := repoMock.NewMockProduct(ctl)
	repo.EXPECT().
		DeleteProductById(ctx, in.GetId()).
		Return(nil).Times(1)

	api := New(repo)

	_, err := api.ProductDelete(ctx, &in)

	require.NoError(t, err)
}

func TestProductDeleteError(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	ctx := context.Background()
	in := pb.ProductDeleteRequest{
		Id: 1,
	}

	repoError := errors.New("some repo error")
	repo := repoMock.NewMockProduct(ctl)
	repo.EXPECT().
		DeleteProductById(ctx, in.GetId()).
		Return(repoError).Times(1)

	api := New(repo)

	_, err := api.ProductDelete(ctx, &in)

	require.EqualError(t, err, "rpc error: code = Internal desc = some repo error")
}

func TestProductDeleteValidationError(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	ctx := context.Background()
	in := pb.ProductDeleteRequest{
		Id: 1,
	}

	repo := repoMock.NewMockProduct(ctl)
	repo.EXPECT().
		DeleteProductById(ctx, in.GetId()).
		Return(productPkg.ErrValidation).Times(1)

	api := New(repo)

	_, err := api.ProductDelete(ctx, &in)

	require.EqualError(t, err, "rpc error: code = InvalidArgument desc = invalid data")
}

func TestProductUpdateSuccess(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	ctx := context.Background()
	in := pb.ProductUpdateRequest{
		Id:     1,
		Title:  "Title",
		Price:  100,
		TypeId: 2,
	}

	repo := repoMock.NewMockProduct(ctl)
	repo.EXPECT().
		UpsertProduct(ctx, models.Product{
			Id:      in.GetId(),
			Title:   in.GetTitle(),
			Price:   in.GetPrice(),
			Type_Id: in.GetTypeId(),
		}).
		Return(in.GetId(), nil).Times(1)

	api := New(repo)

	resp, err := api.ProductUpdate(ctx, &in)

	require.NoError(t, err)
	require.Equal(t, resp.GetId(), in.GetId())
}

func TestProductUpdateError(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	ctx := context.Background()
	in := pb.ProductUpdateRequest{
		Id:     1,
		Title:  "Title",
		Price:  100,
		TypeId: 2,
	}

	repoError := errors.New("some repo error")
	repo := repoMock.NewMockProduct(ctl)
	repo.EXPECT().
		UpsertProduct(ctx, models.Product{
			Id:      in.GetId(),
			Title:   in.GetTitle(),
			Price:   in.GetPrice(),
			Type_Id: in.GetTypeId(),
		}).
		Return(uint64(0), repoError).Times(1)

	api := New(repo)

	_, err := api.ProductUpdate(ctx, &in)

	require.EqualError(t, err, "rpc error: code = Internal desc = some repo error")
}

func TestProductUpdateErrorValidation(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	ctx := context.Background()
	in := pb.ProductUpdateRequest{
		Id:     1,
		Title:  "Title",
		Price:  100,
		TypeId: 2,
	}

	repo := repoMock.NewMockProduct(ctl)
	repo.EXPECT().
		UpsertProduct(ctx, models.Product{
			Id:      in.GetId(),
			Title:   in.GetTitle(),
			Price:   in.GetPrice(),
			Type_Id: in.GetTypeId(),
		}).
		Return(uint64(0), productPkg.ErrValidation).Times(1)

	api := New(repo)

	_, err := api.ProductUpdate(ctx, &in)

	require.EqualError(t, err, "rpc error: code = InvalidArgument desc = invalid data")
}

func TestProductGetError(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	ctx := context.Background()
	in := pb.ProductGetRequest{
		Id: 1,
	}

	repoError := errors.New("some repo error")
	repo := repoMock.NewMockProduct(ctl)
	repo.EXPECT().
		GetProduct(ctx, in.GetId()).
		Return(models.Product{}, repoError).Times(1)

	api := New(repo)

	_, err := api.ProductGet(ctx, &in)

	require.EqualError(t, err, "rpc error: code = Internal desc = some repo error")
}
