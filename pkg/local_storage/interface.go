package store

type LocalStorage interface {
	Set(data *Data, userID int64)
	Read(userID int64) (*Data, bool)
	Delete(userID int64)
}
