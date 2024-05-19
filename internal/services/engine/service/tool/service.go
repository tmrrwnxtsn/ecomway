package tool

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"

	perror "github.com/tmrrwnxtsn/ecomway/internal/pkg/error"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type Repository interface {
	All(ctx context.Context, userID int64) ([]*model.Tool, error)
	GetOne(ctx context.Context, id string, userID int64, externalMethod string) (*model.Tool, error)
	Update(ctx context.Context, tool *model.Tool) error
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) All(ctx context.Context, userID int64) ([]*model.Tool, error) {
	tools, err := s.repository.All(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(tools) > 1 {
		sort.SliceStable(tools, func(i, j int) bool {
			return tools[i].UpdatedAt.UTC().After(tools[j].UpdatedAt.UTC())
		})
	}

	return tools, nil
}

func (s *Service) EditOne(ctx context.Context, id string, userID int64, externalMethod, name string) (*model.Tool, error) {
	tool, err := s.repository.GetOne(ctx, id, userID, externalMethod)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, perror.NewInternal().WithCode(
				perror.CodeObjectNotFound,
			).WithDescription(
				fmt.Sprintf("tool with id %q, userID %v and external method %q not found", id, userID, externalMethod),
			)
		}
		return nil, err
	}

	if tool.Name == name {
		return tool, nil
	}

	tool.Name = name

	if err = s.repository.Update(ctx, tool); err != nil {
		return nil, err
	}

	return tool, nil
}

func (s *Service) RemoveOne(ctx context.Context, id string, userID int64, externalMethod string) error {
	tool, err := s.repository.GetOne(ctx, id, userID, externalMethod)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return perror.NewInternal().WithCode(
				perror.CodeObjectNotFound,
			).WithDescription(
				fmt.Sprintf("tool with id %q, userID %v and external method %q not found", id, userID, externalMethod),
			)
		}
		return err
	}

	if tool.Removed() {
		return nil
	}

	// TODO: добавить обработку ActionSource, чтобы различать удаление администратора и юзера
	tool.Status = model.ToolStatusRemovedByClient

	return s.repository.Update(ctx, tool)
}

func (s *Service) RecoverOne(ctx context.Context, id string, userID int64, externalMethod string) error {
	tool, err := s.repository.GetOne(ctx, id, userID, externalMethod)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return perror.NewInternal().WithCode(
				perror.CodeObjectNotFound,
			).WithDescription(
				fmt.Sprintf("tool with id %q, userID %v and external method %q not found", id, userID, externalMethod),
			)
		}
		return err
	}

	if tool.Status != model.ToolStatusRemovedByAdministrator {
		if tool.Status == model.ToolStatusPendingRecovery {
			return nil
		}

		return perror.NewInternal().WithCode(
			perror.CodeUnresolvedStatusConflict,
		).WithDescription(
			fmt.Sprintf("cannot recover tool with status %v", tool.Status),
		)
	}

	tool.Status = model.ToolStatusPendingRecovery

	return s.repository.Update(ctx, tool)
}
