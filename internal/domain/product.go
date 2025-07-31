package domain

type Product struct {
	ID          int
	Name        string
	Description string
	PriceBuy    int
	PriceSell   int
	SupplierID  int
	Weight      int
}
