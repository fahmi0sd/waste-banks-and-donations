package report

type Repository interface {
	TransactionSummary(filter Filter) ([]TransactionSummary, error)

	// dipakai untuk otorisasi & pembatasan lokasi milik admin
	RoleAndLocation(userID int) (role string, locationID *int, err error)
}
