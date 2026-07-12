package calculator

import (
	"fmt"
	"log/slog"
)

type Service interface {
	Simulate(req SimulateRequest) (SimulateResult, error)
}

type service struct {
	logger *slog.Logger
	repo   Repository
}

func NewService(logger *slog.Logger, repo Repository) Service {
	return &service{logger: logger, repo: repo}
}

func (s *service) Simulate(req SimulateRequest) (SimulateResult, error) {
	result := SimulateResult{Items: make([]SimulateItemResult, 0, len(req.Items))}

	for _, item := range req.Items {
		name, price, found, err := s.repo.ActivePrice(item.CategoryID)
		if err != nil {
			s.logger.Error("failed to fetch active price", "error", err, "category_id", item.CategoryID)
			return SimulateResult{}, fmt.Errorf("gagal mengambil harga kategori id %d", item.CategoryID)
		}
		if !found {
			return SimulateResult{}, fmt.Errorf("kategori id %d tidak ditemukan atau belum punya harga aktif", item.CategoryID)
		}

		subtotal := price * item.Weight
		result.Items = append(result.Items, SimulateItemResult{
			CategoryID:   item.CategoryID,
			CategoryName: name,
			Weight:       item.Weight,
			PricePerUnit: price,
			Subtotal:     subtotal,
		})
		result.TotalRupiah += subtotal
	}

	return result, nil
}
