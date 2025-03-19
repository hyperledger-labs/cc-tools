package querysearch

type OrderSort string

const (
	Desc OrderSort = "desc"
	Asc  OrderSort = "asc"
)

type Sort struct {
	Field string
	Type  OrderSort
}
