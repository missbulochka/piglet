package transactionsgrpc

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	transv1 "github.com/missbulochka/protos/gen/piglet-transactions"
	validation "piglet-transactions-service/internal/domain/validator"
)

func (s *serverAPI) AddCategory(
	ctx context.Context,
	req *transv1.AddCategoryRequest,
) (*transv1.CategoryResponse, error) {
	cat, err := validation.CategoryValidator(
		"",
		req.GetType(),
		req.GetName(),
		req.GetMandatory(),
	)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid creditals")
	}

	if err = s.categories.CreateCategory(ctx, &cat); err != nil {
		// TODO: добавить проверку ошибки на существование категории

		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &transv1.CategoryResponse{
		Category: &transv1.Category{
			Id:           cat.Id.String(),
			Type:         cat.CategoryType,
			CategoryName: cat.Name,
			Mandatory:    cat.Mandatory,
		},
	}, nil
}

func (s *serverAPI) UpdateCategory(
	ctx context.Context,
	req *transv1.Category,
) (*transv1.CategoryResponse, error) {
	cat, err := validation.CategoryValidator(
		req.GetId(),
		req.GetType(),
		req.GetCategoryName(),
		req.GetMandatory(),
	)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid creditals")
	}

	// HACK: проверить существование категории и сравнить с имеющимися данными

	if err = s.categories.UpdateCategory(ctx, &cat); err != nil {
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &transv1.CategoryResponse{
		Category: &transv1.Category{
			Id:           cat.Id.String(),
			Type:         cat.CategoryType,
			CategoryName: cat.Name,
			Mandatory:    cat.Mandatory,
		},
	}, nil
}

func (s *serverAPI) GetCategory(
	ctx context.Context,
	req *transv1.IdRequest,
) (*transv1.CategoryResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid creditals")
	}

	cat, err := s.categories.GetCategory(ctx, id)
	if err != nil {
		// TODO: добавить проверку ошибки на существование категории

		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &transv1.CategoryResponse{
		Category: &transv1.Category{
			Id:           cat.Id.String(),
			Type:         cat.CategoryType,
			CategoryName: cat.Name,
			Mandatory:    cat.Mandatory,
		},
	}, nil
}

func (s *serverAPI) GetAllCategories(
	ctx context.Context,
	req *emptypb.Empty,
) (*transv1.GetAllCategoriesResponse, error) {
	cat, err := s.categories.GetAllCategories(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	resTrans := ReturnCategories(cat)

	return &transv1.GetAllCategoriesResponse{Category: resTrans}, nil
}

func (s *serverAPI) DeleteCategory(
	ctx context.Context,
	req *transv1.IdRequest,
) (*transv1.SuccessResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid creditals")
	}

	if err = s.categories.DeleteCategory(ctx, id); err != nil {
		return &transv1.SuccessResponse{Success: false}, status.Errorf(codes.Internal, "internal error")
	}

	return &transv1.SuccessResponse{Success: true}, nil
}
