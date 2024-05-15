package transactionsgrpc

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
		// TODO: проверка ошибки о несуществовании счета или категории

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

	// TODO: проверка существования категории

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
		// TODO: проверка ошибки о несуществовании категории

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
	req *transv1.EmptyRequest,
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
