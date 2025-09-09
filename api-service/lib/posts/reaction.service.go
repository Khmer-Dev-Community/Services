package posts

import (
	"context"
	"fmt"

	"github.com/Khmer-Dev-Community/Services/api-service/utils"
)

type ReactionService interface {
	CreateReaction(ctx context.Context, req Reaction, authorID uint) (*Reaction, error)
}

type Service struct {
	repo ReactionRepository
}

// NewReactionService creates a new instance of ReactionService.
func NewReactionService(repo ReactionRepository) ReactionService {
	return &Service{repo: repo}
}

func (s *Service) CreateReaction(ctx context.Context, req Reaction, authorID uint) (*Reaction, error) {
	utils.InfoLog(req, "(s *Service) CreateReaction")
	if req.ReactionType == "" {
		err := fmt.Errorf("reaction content cannot be empty")
		utils.InfoLog(err, "Create Reaction")
		return nil, err
	}
	reaction, err := s.repo.CreateOrUpdate(ctx, &req)
	if err != nil {
		utils.ErrorLog(err.Error(), "Create Reaction")
		return nil, err
	}

	return reaction, nil
}
