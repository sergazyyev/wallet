package store

//Store sore for app
type Store interface {
	WalletRepository() WalletRepository
	ReportRepository() ReportRepository
}
