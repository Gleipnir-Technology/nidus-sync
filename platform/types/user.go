package types

type User struct {
	ID   int32  `db:"id" json:"id"`
	Name string `db:"-" json:"name"`
}
