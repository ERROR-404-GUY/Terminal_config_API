package domain

import "terminal_config/internal/ports/mongo"

type TerminalConfigService interface {
}

type terminalConfigService struct {
	repo mongo.TerminalConfigRepository
}

func NewTerminalConfigService(repo mongo.TerminalConfigRepository) TerminalConfigService {
	return &terminalConfigService{
		repo: repo,
	}
}

// Define methods for business logic related to terminal configurations
