package report

import (
	"github.com/fahmi0sd/waste-banks-and-donations/service/report"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

var _ report.Repository = (*GormRepository)(nil)

type summaryRow struct {
	LocationID   int
	LocationName string
	Date         string
	TotalTrx     int
	TotalRupiah  float64
}

func (r *GormRepository) TransactionSummary(filter report.Filter) ([]report.TransactionSummary, error) {
	query := r.db.Table("waste_transaction AS t").
		Select(`
			t.location_id,
			l.name AS location_name,
			TO_CHAR(t.created_at, 'YYYY-MM-DD') AS date,
			COUNT(t.id) AS total_trx,
			COALESCE(SUM(t.total_rupiah), 0) AS total_rupiah
		`).
		Joins("JOIN locations l ON l.id = t.location_id").
		Where("t.status = 'completed'")

	if filter.LocationID != nil {
		query = query.Where("t.location_id = ?", *filter.LocationID)
	}
	if filter.DateFrom != nil {
		query = query.Where("t.created_at::date >= ?", *filter.DateFrom)
	}
	if filter.DateTo != nil {
		query = query.Where("t.created_at::date <= ?", *filter.DateTo)
	}

	var rows []summaryRow
	err := query.
		Group("t.location_id, l.name, TO_CHAR(t.created_at, 'YYYY-MM-DD')").
		Order("date DESC, t.location_id ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	list := make([]report.TransactionSummary, 0, len(rows))
	for _, row := range rows {
		list = append(list, report.TransactionSummary{
			LocationID: row.LocationID, LocationName: row.LocationName,
			Date: row.Date, TotalTrx: row.TotalTrx, TotalRupiah: row.TotalRupiah,
		})
	}
	return list, nil
}

func (r *GormRepository) RoleAndLocation(userID int) (string, *int, error) {
	type row struct {
		Role       string
		LocationID *int
	}
	var res row
	err := r.db.Table("users").Select("role, location_id").Where("id = ?", userID).Take(&res).Error
	if err != nil {
		return "", nil, err
	}
	return res.Role, res.LocationID, nil
}
