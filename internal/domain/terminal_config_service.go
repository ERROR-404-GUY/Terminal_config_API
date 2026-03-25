package domain

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type TerminalRepository interface {
	CreateTerminal(ctx context.Context, terminal *Terminal) error
	GetTerminalByTID(ctx context.Context, tid string) (*Terminal, error)
	UpdateTerminal(ctx context.Context, terminal *Terminal) error
	UpdateRefundAllowed(ctx context.Context, tid string, allowed bool) error
	DeleteTerminal(ctx context.Context, tid string) error
	ListTerminals(ctx context.Context) ([]*Terminal, error)
}

type TerminalConfigService interface {
	CreateTerminal(ctx context.Context, terminal *Terminal) error
	CreateRandomTerminal(ctx context.Context) (*Terminal, error)
	GetTerminal(ctx context.Context, tid string) (*Terminal, error)
	UpdateTerminal(ctx context.Context, terminal *Terminal) error
	UpdateRefundAllowed(ctx context.Context, tid string, allowed bool) error
	DeleteTerminal(ctx context.Context, tid string) error
	ListTerminals(ctx context.Context) ([]*Terminal, error)
}

type terminalConfigService struct {
	repo TerminalRepository
}

type Terminal struct {
	ID             string `bson:"_id,omitempty" json:"id,omitempty"`
	TID            string `bson:"tid" json:"tid"`
	SerialNumber   string `bson:"serial_number" json:"serial_number"`
	RefundAllowed  bool   `bson:"refund_allowed" json:"refund_allowed"`
	ProductName    string `bson:"product_name" json:"product_name"`
	ActivationCode string `bson:"activation_code" json:"activation_code"`
}

func NewTerminalConfigService(repo TerminalRepository) TerminalConfigService {
	rand.Seed(time.Now().UnixNano())
	return &terminalConfigService{
		repo: repo,
	}
}

func (s *terminalConfigService) CreateTerminal(ctx context.Context, terminal *Terminal) error {
	if terminal.TID == "" {
		terminal.TID = GenerateTID()
	}
	terminal.ID = terminal.TID
	if terminal.ProductName == "" {
		terminal.ProductName = DefaultProductName()
	}
	if terminal.ActivationCode == "" {
		terminal.ActivationCode = GenerateActivationCode()
	}

	return s.repo.CreateTerminal(ctx, terminal)
}

func (s *terminalConfigService) CreateRandomTerminal(ctx context.Context) (*Terminal, error) {
	term := &Terminal{
		TID:            GenerateTID(),
		SerialNumber:   GenerateSerialNumber(),
		RefundAllowed:  rand.Intn(2) == 0,
		ProductName:    DefaultProductName(),
		ActivationCode: GenerateActivationCode(),
	}
	term.ID = term.TID

	if err := s.repo.CreateTerminal(ctx, term); err != nil {
		return nil, err
	}
	return term, nil
}

func GenerateTID() string {
	return fmt.Sprintf("TID-%d-%06d", time.Now().UnixNano(), rand.Intn(1000000))
}

func DefaultProductName() string {
	return "POS Terminal"
}

func GenerateActivationCode() string {
	return fmt.Sprintf("ACT-%06d", rand.Intn(1000000))
}

func GenerateSerialNumber() string {
	return fmt.Sprintf("SN-%d", rand.Intn(900000)+100000)
}

func (s *terminalConfigService) GetTerminal(ctx context.Context, tid string) (*Terminal, error) {
	return s.repo.GetTerminalByTID(ctx, tid)
}

func (s *terminalConfigService) UpdateTerminal(ctx context.Context, terminal *Terminal) error {
	return s.repo.UpdateTerminal(ctx, terminal)
}

func (s *terminalConfigService) UpdateRefundAllowed(ctx context.Context, tid string, allowed bool) error {
	return s.repo.UpdateRefundAllowed(ctx, tid, allowed)
}

func (s *terminalConfigService) DeleteTerminal(ctx context.Context, tid string) error {
	return s.repo.DeleteTerminal(ctx, tid)
}

func (s *terminalConfigService) ListTerminals(ctx context.Context) ([]*Terminal, error) {
	return s.repo.ListTerminals(ctx)
}
