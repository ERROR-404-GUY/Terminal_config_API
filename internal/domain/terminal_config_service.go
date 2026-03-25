package domain

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type TerminalRepository interface {
	CreateTerminal(ctx context.Context, terminal *Terminal) error
	GetTerminalByTID(ctx context.Context, tid string) (*Terminal, error)
	UpdateTerminal(ctx context.Context, terminal *Terminal) error
	DeleteTerminal(ctx context.Context, tid string) error
	ListTerminals(ctx context.Context) ([]*Terminal, error)
}

type TerminalConfigService interface {
	CreateTerminal(ctx context.Context, terminal *Terminal) error
	GetTerminal(ctx context.Context, tid string) (*Terminal, error)
	UpdateTerminal(ctx context.Context, terminal *Terminal) error
	DeleteTerminal(ctx context.Context, tid string) error
	ListTerminals(ctx context.Context) ([]*Terminal, error)
}

type terminalConfigService struct {
	repo TerminalRepository
}

type Terminal struct {
	ID             bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	TID            string        `bson:"tid" json:"tid"`
	SerialNumber   string        `bson:"serial_number" json:"serial_number"`
	RefundAllowed  bool          `bson:"refund_allowed" json:"refund_allowed"`
	ProductName    string        `bson:"product_name" json:"product_name"`
	ActivationCode string        `bson:"activation_code" json:"activation_code"`
}

func NewTerminalConfigService(repo TerminalRepository) TerminalConfigService {
	return &terminalConfigService{
		repo: repo,
	}
}

func (s *terminalConfigService) CreateTerminal(ctx context.Context, terminal *Terminal) error {
	if terminal.TID == "" {
		terminal.TID = GenerateTID()
	}
	if terminal.ProductName == "" {
		terminal.ProductName = DefaultProductName()
	}
	if terminal.ActivationCode == "" {
		terminal.ActivationCode = GenerateActivationCode()
	}

	return s.repo.CreateTerminal(ctx, terminal)
}

func GenerateTID() string {
	return fmt.Sprintf("TID-%d", time.Now().UnixNano())
}

func DefaultProductName() string {
	return "POS Terminal"
}

func GenerateActivationCode() string {
	return fmt.Sprintf("ACT-%06d", rand.Intn(1000000))
}

func (s *terminalConfigService) GetTerminal(ctx context.Context, tid string) (*Terminal, error) {
	return s.repo.GetTerminalByTID(ctx, tid)
}

func (s *terminalConfigService) UpdateTerminal(ctx context.Context, terminal *Terminal) error {
	return s.repo.UpdateTerminal(ctx, terminal)
}

func (s *terminalConfigService) DeleteTerminal(ctx context.Context, tid string) error {
	return s.repo.DeleteTerminal(ctx, tid)
}

func (s *terminalConfigService) ListTerminals(ctx context.Context) ([]*Terminal, error) {
	return s.repo.ListTerminals(ctx)
}
