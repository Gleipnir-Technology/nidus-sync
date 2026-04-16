package types

type Parcel struct {
	APN         string `db:"apn" json:"apn"`
	ID          int32  `db:"id" json:"id"`
	Description string `db:"description" json:"description"`
}
