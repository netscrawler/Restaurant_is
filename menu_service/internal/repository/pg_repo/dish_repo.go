package pgrepo

type dishPgRepo struct{}

func NewDishPgRepo() *dishPgRepo {
	return &dishPgRepo{}
}
