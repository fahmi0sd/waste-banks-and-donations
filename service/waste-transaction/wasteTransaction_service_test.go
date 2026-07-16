package wastetransaction_test

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"testing"

	wastetransaction "github.com/fahmi0sd/waste-banks-and-donations/service/waste-transaction"
	"github.com/fahmi0sd/waste-banks-and-donations/service/waste-transaction/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var (
	testDB      *gorm.DB
	dbAvailable bool
)

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func TestMain(m *testing.M) {
	dbName := envOr("TEST_DB_NAME", "banksampah_test")
	if !strings.Contains(strings.ToLower(dbName), "test") {
		fmt.Println("SKIP setup test DB: TEST_DB_NAME harus mengandung kata 'test' (pengaman supaya tidak salah konek ke database asli)")
		os.Exit(m.Run())
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		envOr("TEST_DB_HOST", "localhost"),
		envOr("TEST_DB_PORT", "5432"),
		envOr("TEST_DB_USER", "postgres"),
		envOr("TEST_DB_PASS", "postgres"),
		dbName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: gormlogger.Default.LogMode(gormlogger.Silent)})
	if err == nil {
		if sqlDB, pingErr := db.DB(); pingErr == nil && sqlDB.Ping() == nil {
			if seedErr := seedTestUsers(db); seedErr == nil {
				testDB = db
				dbAvailable = true
			} else {
				fmt.Println("SKIP: gagal seed test database:", seedErr)
			}
		}
	}
	if !dbAvailable {
		fmt.Println("SKIP: test database tidak tersedia (set TEST_DB_HOST/TEST_DB_PORT/TEST_DB_USER/TEST_DB_PASS/TEST_DB_NAME) — test yang butuh pkg.GetUserIdentity akan di-skip")
	}

	os.Exit(m.Run())
}

func seedTestUsers(db *gorm.DB) error {
	if err := db.Exec(`DROP TABLE IF EXISTS users CASCADE`).Error; err != nil {
		return fmt.Errorf("drop table users: %w", err)
	}
	if err := db.Exec(`CREATE TABLE users (id INT PRIMARY KEY, role VARCHAR(20), location_id INT)`).Error; err != nil {
		return fmt.Errorf("create table users: %w", err)
	}
	if err := db.Exec(`INSERT INTO users (id, role, location_id) VALUES
		(1, 'user', NULL),
		(2, 'admin', 1),
		(3, 'admin', 2),
		(4, 'master_admin', NULL)`).Error; err != nil {
		return fmt.Errorf("seed users: %w", err)
	}
	return nil
}

func requireTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	if !dbAvailable {
		t.Skip("test database tidak tersedia, lewati test yang butuh pkg.GetUserIdentity")
	}
	return testDB
}

func silentLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError + 1}))
}

func intPtr(i int) *int { return &i }

func TestCreate_SuksesAdminLokasiSendiri(t *testing.T) {
	db := requireTestDB(t)
	repo := mocks.NewMockRepository(t)
	notifier := mocks.NewMockNotificationSender(t)

	repo.EXPECT().QueueInfo(10).Return(wastetransaction.QueueInfo{UserID: 1, LocationID: 1, Status: "verified", Found: true}, nil)
	created := wastetransaction.WasteTransaction{ID: 99, QueueID: intPtr(10), UserID: 1, LocationID: 1, AdminID: 2, TotalRupiah: 17500, Status: wastetransaction.StatusCompleted}
	repo.EXPECT().CreateWithEffects(mock.MatchedBy(func(in wastetransaction.CreateInput) bool {
		return in.QueueID == 10 && in.UserID == 1 && in.LocationID == 1 && in.AdminID == 2
	})).Return(created, nil)
	repo.EXPECT().GetDetailsByTransactionID(99).Return([]wastetransaction.WasteTransactionDetail{}, nil)
	notifier.EXPECT().NotifyDepositUpdate(1, 17500.0).Return(nil)

	svc := wastetransaction.NewService(silentLogger(), repo, db, notifier)
	result, err := svc.Create(2, 10, []wastetransaction.ItemInput{{CategoryID: 1, Weight: 5}})

	assert.NoError(t, err)
	assert.Equal(t, 99, result.ID)
}

func TestCreate_DitolakBukanAdmin(t *testing.T) {
	db := requireTestDB(t)
	repo := mocks.NewMockRepository(t)

	svc := wastetransaction.NewService(silentLogger(), repo, db, nil)
	_, err := svc.Create(1, 10, []wastetransaction.ItemInput{{CategoryID: 1, Weight: 5}}) // user_id=1 role=user

	assert.Error(t, err)
	repo.AssertNotCalled(t, "CreateWithEffects", mock.Anything)
	repo.AssertNotCalled(t, "QueueInfo", mock.Anything)
}

func TestCreate_DitolakItemsKosong(t *testing.T) {
	db := requireTestDB(t)
	repo := mocks.NewMockRepository(t)

	svc := wastetransaction.NewService(silentLogger(), repo, db, nil)
	_, err := svc.Create(2, 10, []wastetransaction.ItemInput{})

	assert.Error(t, err)
	repo.AssertNotCalled(t, "QueueInfo", mock.Anything)
}

func TestCreate_DitolakBeratNolAtauNegatif(t *testing.T) {
	db := requireTestDB(t)
	repo := mocks.NewMockRepository(t)

	svc := wastetransaction.NewService(silentLogger(), repo, db, nil)
	_, err := svc.Create(2, 10, []wastetransaction.ItemInput{{CategoryID: 1, Weight: 0}})

	assert.Error(t, err)
	repo.AssertNotCalled(t, "QueueInfo", mock.Anything)
}

func TestCreate_DitolakQueueTidakDitemukan(t *testing.T) {
	db := requireTestDB(t)
	repo := mocks.NewMockRepository(t)
	repo.EXPECT().QueueInfo(999).Return(wastetransaction.QueueInfo{Found: false}, nil)

	svc := wastetransaction.NewService(silentLogger(), repo, db, nil)
	_, err := svc.Create(2, 999, []wastetransaction.ItemInput{{CategoryID: 1, Weight: 5}})

	assert.Error(t, err)
	repo.AssertNotCalled(t, "CreateWithEffects", mock.Anything)
}

func TestCreate_DitolakQueueBedaLokasi(t *testing.T) {
	db := requireTestDB(t)
	repo := mocks.NewMockRepository(t)
	repo.EXPECT().QueueInfo(10).Return(wastetransaction.QueueInfo{UserID: 1, LocationID: 2, Status: "verified", Found: true}, nil)

	svc := wastetransaction.NewService(silentLogger(), repo, db, nil)
	_, err := svc.Create(2, 10, []wastetransaction.ItemInput{{CategoryID: 1, Weight: 5}})

	assert.Error(t, err)
	repo.AssertNotCalled(t, "CreateWithEffects", mock.Anything)
}

func TestCreate_DitolakQueueBelumVerified(t *testing.T) {
	db := requireTestDB(t)
	repo := mocks.NewMockRepository(t)
	repo.EXPECT().QueueInfo(10).Return(wastetransaction.QueueInfo{UserID: 1, LocationID: 1, Status: "waiting", Found: true}, nil)

	svc := wastetransaction.NewService(silentLogger(), repo, db, nil)
	_, err := svc.Create(2, 10, []wastetransaction.ItemInput{{CategoryID: 1, Weight: 5}})

	assert.Error(t, err)
	repo.AssertNotCalled(t, "CreateWithEffects", mock.Anything)
}

func TestCreate_NotifikasiGagalTidakMenggagalkanTransaksi(t *testing.T) {
	db := requireTestDB(t)
	repo := mocks.NewMockRepository(t)
	notifier := mocks.NewMockNotificationSender(t)

	repo.EXPECT().QueueInfo(10).Return(wastetransaction.QueueInfo{UserID: 1, LocationID: 1, Status: "verified", Found: true}, nil)
	created := wastetransaction.WasteTransaction{ID: 99, UserID: 1, LocationID: 1, AdminID: 2, TotalRupiah: 17500, Status: wastetransaction.StatusCompleted}
	repo.EXPECT().CreateWithEffects(mock.Anything).Return(created, nil)
	repo.EXPECT().GetDetailsByTransactionID(99).Return([]wastetransaction.WasteTransactionDetail{}, nil)
	notifier.EXPECT().NotifyDepositUpdate(1, 17500.0).Return(errors.New("mailjet down"))

	svc := wastetransaction.NewService(silentLogger(), repo, db, notifier)
	result, err := svc.Create(2, 10, []wastetransaction.ItemInput{{CategoryID: 1, Weight: 5}})

	assert.NoError(t, err, "Create harus tetap sukses meski notifikasi gagal")
	assert.Equal(t, 99, result.ID)
}

func TestGetByID_UserPemilikBolehAkses(t *testing.T) {
	db := requireTestDB(t)
	repo := mocks.NewMockRepository(t)
	tx := wastetransaction.WasteTransaction{ID: 5, UserID: 1, LocationID: 1}
	repo.EXPECT().GetByID(5).Return(tx, nil)
	repo.EXPECT().GetDetailsByTransactionID(5).Return([]wastetransaction.WasteTransactionDetail{}, nil)

	svc := wastetransaction.NewService(silentLogger(), repo, db, nil)
	result, err := svc.GetByID(5, 1)

	assert.NoError(t, err)
	assert.Equal(t, 5, result.ID)
}

func TestGetByID_UserBukanPemilikDitolak(t *testing.T) {
	db := requireTestDB(t)
	repo := mocks.NewMockRepository(t)
	tx := wastetransaction.WasteTransaction{ID: 5, UserID: 1, LocationID: 1}
	repo.EXPECT().GetByID(5).Return(tx, nil)

	svc := wastetransaction.NewService(silentLogger(), repo, db, nil)
	_, err := svc.GetByID(5, 99)

	assert.Error(t, err)
}

func TestGetByID_AdminLokasiSendiriBolehAkses(t *testing.T) {
	db := requireTestDB(t)
	repo := mocks.NewMockRepository(t)
	tx := wastetransaction.WasteTransaction{ID: 5, UserID: 1, LocationID: 1}
	repo.EXPECT().GetByID(5).Return(tx, nil)
	repo.EXPECT().GetDetailsByTransactionID(5).Return([]wastetransaction.WasteTransactionDetail{}, nil)

	svc := wastetransaction.NewService(silentLogger(), repo, db, nil)
	result, err := svc.GetByID(5, 2)

	assert.NoError(t, err)
	assert.Equal(t, 5, result.ID)
}

func TestGetByID_AdminLokasiLainDitolak(t *testing.T) {
	db := requireTestDB(t)
	repo := mocks.NewMockRepository(t)
	tx := wastetransaction.WasteTransaction{ID: 5, UserID: 1, LocationID: 1}
	repo.EXPECT().GetByID(5).Return(tx, nil)

	svc := wastetransaction.NewService(silentLogger(), repo, db, nil)
	result, err := svc.GetByID(5, 3)

	assert.Error(t, err)
	assert.Equal(t, wastetransaction.WasteTransaction{}, result)
}

func TestGetByID_MasterAdminBolehAksesApaSaja(t *testing.T) {
	db := requireTestDB(t)
	repo := mocks.NewMockRepository(t)
	tx := wastetransaction.WasteTransaction{ID: 5, UserID: 1, LocationID: 1}
	repo.EXPECT().GetByID(5).Return(tx, nil)
	repo.EXPECT().GetDetailsByTransactionID(5).Return([]wastetransaction.WasteTransactionDetail{}, nil)

	svc := wastetransaction.NewService(silentLogger(), repo, db, nil)
	result, err := svc.GetByID(5, 4)

	assert.NoError(t, err)
	assert.Equal(t, 5, result.ID)
}

func TestCancel_SuksesAdminLokasiSendiri(t *testing.T) {
	db := requireTestDB(t)
	repo := mocks.NewMockRepository(t)
	tx := wastetransaction.WasteTransaction{ID: 5, LocationID: 1, Status: wastetransaction.StatusCompleted}
	repo.EXPECT().GetByID(5).Return(tx, nil)
	cancelled := wastetransaction.WasteTransaction{ID: 5, LocationID: 1, Status: wastetransaction.StatusCancelled}
	repo.EXPECT().CancelWithEffects(5).Return(cancelled, nil)

	svc := wastetransaction.NewService(silentLogger(), repo, db, nil)
	result, err := svc.Cancel(2, 5)

	assert.NoError(t, err)
	assert.Equal(t, wastetransaction.StatusCancelled, result.Status)
}

func TestCancel_DitolakBukanAdmin(t *testing.T) {
	db := requireTestDB(t)
	repo := mocks.NewMockRepository(t)

	svc := wastetransaction.NewService(silentLogger(), repo, db, nil)
	_, err := svc.Cancel(1, 5)

	assert.Error(t, err)
	repo.AssertNotCalled(t, "GetByID", mock.Anything)
}

func TestCancel_DitolakLokasiBerbeda(t *testing.T) {
	db := requireTestDB(t)
	repo := mocks.NewMockRepository(t)
	tx := wastetransaction.WasteTransaction{ID: 5, LocationID: 1, Status: wastetransaction.StatusCompleted}
	repo.EXPECT().GetByID(5).Return(tx, nil)

	svc := wastetransaction.NewService(silentLogger(), repo, db, nil)
	_, err := svc.Cancel(3, 5)

	assert.Error(t, err)
	repo.AssertNotCalled(t, "CancelWithEffects", mock.Anything)
}

func TestCancel_DitolakStatusBukanCompleted(t *testing.T) {
	db := requireTestDB(t)
	repo := mocks.NewMockRepository(t)
	tx := wastetransaction.WasteTransaction{ID: 5, LocationID: 1, Status: wastetransaction.StatusCancelled}
	repo.EXPECT().GetByID(5).Return(tx, nil)

	svc := wastetransaction.NewService(silentLogger(), repo, db, nil)
	_, err := svc.Cancel(2, 5)

	assert.Error(t, err)
	repo.AssertNotCalled(t, "CancelWithEffects", mock.Anything)
}

func TestLocationHistory_SuksesLokasiSendiri(t *testing.T) {
	db := requireTestDB(t)
	repo := mocks.NewMockRepository(t)
	list := []wastetransaction.WasteTransaction{{ID: 1, LocationID: 1}, {ID: 2, LocationID: 1}}
	repo.EXPECT().ListByLocation(1).Return(list, nil)

	svc := wastetransaction.NewService(silentLogger(), repo, db, nil)
	result, err := svc.LocationHistory(2, 1)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestLocationHistory_DitolakLokasiLain(t *testing.T) {
	db := requireTestDB(t)
	repo := mocks.NewMockRepository(t)

	svc := wastetransaction.NewService(silentLogger(), repo, db, nil)
	_, err := svc.LocationHistory(2, 2)

	assert.Error(t, err)
	repo.AssertNotCalled(t, "ListByLocation", mock.Anything)
}

func TestLocationHistory_DitolakBukanAdmin(t *testing.T) {
	db := requireTestDB(t)
	repo := mocks.NewMockRepository(t)

	svc := wastetransaction.NewService(silentLogger(), repo, db, nil)
	_, err := svc.LocationHistory(1, 1)

	assert.Error(t, err)
	repo.AssertNotCalled(t, "ListByLocation", mock.Anything)
}
