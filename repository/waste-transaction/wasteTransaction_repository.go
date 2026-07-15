package wastetransaction

import (
	"errors"
	"fmt"

	wastetransaction "github.com/fahmi0sd/waste-banks-and-donations/service/waste-transaction"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

var _ wastetransaction.Repository = (*GormRepository)(nil)

func (r *GormRepository) QueueInfo(queueID int) (wastetransaction.QueueInfo, error) {
	type row struct {
		UserID     int
		LocationID int
		Status     string
	}
	var res row
	err := r.db.Table("queue").
		Select("user_id, location_id, status").
		Where("id = ?", queueID).
		Take(&res).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return wastetransaction.QueueInfo{Found: false}, nil
		}
		return wastetransaction.QueueInfo{}, err
	}
	return wastetransaction.QueueInfo{
		UserID: res.UserID, LocationID: res.LocationID, Status: res.Status, Found: true,
	}, nil
}

func (r *GormRepository) CreateWithEffects(input wastetransaction.CreateInput) (wastetransaction.WasteTransaction, error) {
	var result wastetransaction.WasteTransaction

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var queueStatus string
		if err := tx.Raw(`SELECT status FROM queue WHERE id = ? FOR UPDATE`, input.QueueID).
			Scan(&queueStatus).Error; err != nil {
			return err
		}
		if queueStatus != "verified" {
			return errors.New("antrian sudah tidak berstatus verified (kemungkinan sedang/sudah diproses)")
		}

		type detailToInsert struct {
			CategoryID           int
			Weight               float64
			PricePerUnitSnapshot float64
			Subtotal             float64
		}
		var details []detailToInsert
		var totalRupiah float64

		for _, item := range input.Items {
			var priceRow struct{ PricePerUnit float64 }
			err := tx.Table("waste_price").
				Select("price_per_unit").
				Where("category_id = ? AND effective_until IS NULL", item.CategoryID).
				Take(&priceRow).Error
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return fmt.Errorf("kategori id %d tidak ditemukan atau belum punya harga aktif", item.CategoryID)
				}
				return err
			}
			subtotal := priceRow.PricePerUnit * item.Weight
			totalRupiah += subtotal
			details = append(details, detailToInsert{
				CategoryID: item.CategoryID, Weight: item.Weight,
				PricePerUnitSnapshot: priceRow.PricePerUnit, Subtotal: subtotal,
			})
		}

		wt := wastetransaction.WasteTransaction{
			QueueID: &input.QueueID, UserID: input.UserID, LocationID: input.LocationID,
			AdminID: input.AdminID, TotalRupiah: totalRupiah, Status: wastetransaction.StatusCompleted,
		}
		if err := tx.Create(&wt).Error; err != nil {
			return err
		}

		for _, d := range details {
			detailRow := wastetransaction.WasteTransactionDetail{
				TransactionID: wt.ID, CategoryID: d.CategoryID, Weight: d.Weight,
				PricePerUnitSnapshot: d.PricePerUnitSnapshot, Subtotal: d.Subtotal,
			}
			if err := tx.Create(&detailRow).Error; err != nil {
				return err
			}

			err := tx.Exec(`
				INSERT INTO location_inventory (location_id, category_id, weight_kg, updated_at)
				VALUES (?, ?, ?, NOW())
				ON CONFLICT (location_id, category_id)
				DO UPDATE SET weight_kg = location_inventory.weight_kg + EXCLUDED.weight_kg, updated_at = NOW()
			`, input.LocationID, d.CategoryID, d.Weight).Error
			if err != nil {
				return err
			}
		}

		if err := tx.Exec(`
			INSERT INTO wallet (user_id, balance, updated_at) VALUES (?, 0, NOW())
			ON CONFLICT (user_id) DO NOTHING
		`, input.UserID).Error; err != nil {
			return err
		}
		if err := tx.Exec(`
			UPDATE wallet SET balance = balance + ?, updated_at = NOW() WHERE user_id = ?
		`, totalRupiah, input.UserID).Error; err != nil {
			return err
		}

		var walletID int
		if err := tx.Raw(`SELECT id FROM wallet WHERE user_id = ?`, input.UserID).Scan(&walletID).Error; err != nil {
			return err
		}
		if err := tx.Exec(`
			INSERT INTO wallet_transaction (wallet_id, type, amount, reference_type, reference_id, created_at)
			VALUES (?, 'credit_waste', ?, 'waste_transaction', ?, NOW())
		`, walletID, totalRupiah, wt.ID).Error; err != nil {
			return err
		}

		if err := tx.Exec(`UPDATE queue SET status = 'completed' WHERE id = ?`, input.QueueID).Error; err != nil {
			return err
		}

		result = wt
		return nil
	})

	return result, err
}

func (r *GormRepository) GetByID(id int) (wastetransaction.WasteTransaction, error) {
	var wt wastetransaction.WasteTransaction
	if err := r.db.First(&wt, "id = ?", id).Error; err != nil {
		return wastetransaction.WasteTransaction{}, err
	}
	return wt, nil
}

func (r *GormRepository) GetDetailsByTransactionID(transactionID int) ([]wastetransaction.WasteTransactionDetail, error) {
	var details []wastetransaction.WasteTransactionDetail
	err := r.db.Table("waste_transaction_detail AS d").
		Select("d.id, d.transaction_id, d.category_id, c.name AS category_name, d.weight, d.price_per_unit_snapshot, d.subtotal").
		Joins("JOIN waste_category c ON c.id = d.category_id").
		Where("d.transaction_id = ?", transactionID).
		Scan(&details).Error
	return details, err
}

func (r *GormRepository) ListByUser(userID int) ([]wastetransaction.WasteTransaction, error) {
	var list []wastetransaction.WasteTransaction
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&list).Error
	return list, err
}

func (r *GormRepository) ListByLocation(locationID int) ([]wastetransaction.WasteTransaction, error) {
	var list []wastetransaction.WasteTransaction
	err := r.db.Where("location_id = ?", locationID).Order("created_at DESC").Find(&list).Error
	return list, err
}

func (r *GormRepository) CancelWithEffects(transactionID int) (wastetransaction.WasteTransaction, error) {
	var result wastetransaction.WasteTransaction

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var current struct {
			Status      string
			UserID      int
			LocationID  int
			TotalRupiah float64
		}
		if err := tx.Raw(`
			SELECT status, user_id, location_id, total_rupiah
			FROM waste_transaction WHERE id = ? FOR UPDATE
		`, transactionID).Scan(&current).Error; err != nil {
			return err
		}
		if current.Status != wastetransaction.StatusCompleted {
			return fmt.Errorf("transaksi berstatus %s, tidak bisa dibatalkan", current.Status)
		}

		var balance float64
		if err := tx.Raw(`SELECT balance FROM wallet WHERE user_id = ?`, current.UserID).
			Scan(&balance).Error; err != nil {
			return err
		}
		if balance < current.TotalRupiah {
			return errors.New("saldo wallet user sudah tidak cukup untuk dibatalkan (kemungkinan sudah ditarik/didonasikan)")
		}

		var details []struct {
			CategoryID int
			Weight     float64
		}
		if err := tx.Table("waste_transaction_detail").
			Select("category_id, weight").
			Where("transaction_id = ?", transactionID).
			Scan(&details).Error; err != nil {
			return err
		}

		for _, d := range details {
			if err := tx.Exec(`
				UPDATE location_inventory SET weight_kg = weight_kg - ?, updated_at = NOW()
				WHERE location_id = ? AND category_id = ?
			`, d.Weight, current.LocationID, d.CategoryID).Error; err != nil {
				return err
			}
		}

		if err := tx.Exec(`
			UPDATE wallet SET balance = balance - ?, updated_at = NOW() WHERE user_id = ?
		`, current.TotalRupiah, current.UserID).Error; err != nil {
			return err
		}

		var walletID int
		if err := tx.Raw(`SELECT id FROM wallet WHERE user_id = ?`, current.UserID).
			Scan(&walletID).Error; err != nil {
			return err
		}
		if err := tx.Exec(`
			INSERT INTO wallet_transaction (wallet_id, type, amount, reference_type, reference_id, created_at)
			VALUES (?, 'cancellation', ?, 'waste_transaction', ?, NOW())
		`, walletID, -current.TotalRupiah, transactionID).Error; err != nil {
			return err
		}

		if err := tx.Exec(`UPDATE waste_transaction SET status = 'cancelled' WHERE id = ?`, transactionID).Error; err != nil {
			return err
		}

		result = wastetransaction.WasteTransaction{
			ID: transactionID, UserID: current.UserID, LocationID: current.LocationID,
			TotalRupiah: current.TotalRupiah, Status: wastetransaction.StatusCancelled,
		}
		return nil
	})

	return result, err
}
