package pgsql

import (
	cart "github.com/amelendres/go-shopping-cart/pkg"
	"reflect"
	"unsafe"
)

type CartDTO struct {
	ID, BuyerID string
	Products    []ProductDTO
}

type ProductDTO struct {
	ID    string
	Name  string
	Price float64
	Qty int
}

func (c CartDTO) toDomain() *cart.Cart {
	ca := &cart.Cart{}

	id := reflect.ValueOf(ca).Elem().FieldByName("id")
	SetUnexportedField(id, cart.UUID(c.ID))
	SetUnexportedField(reflect.ValueOf(ca).Elem().FieldByName("buyerID"), cart.UUID(c.BuyerID))

	return ca
}

func SetUnexportedField(field reflect.Value, value interface{}) {
	reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).
		Elem().
		Set(reflect.ValueOf(value))
}
