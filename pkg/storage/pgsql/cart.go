package pgsql

import (
	"database/sql"
	"fmt"
	cart "github.com/amelendres/go-shopping-cart/pkg"

)

// CartRepository manages the operations with the database
type CartRepository struct {
	db    *sql.DB
	*Queries
}

func NewCartRepository(db *sql.DB) cart.Repository {
	return &CartRepository{db: db, Queries: New(db)}
}

// ExecTx executes a function within a database transaction
func (r *CartRepository) execTx(fn func(*Queries) error) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rollback err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

const getCartSQL = `SELECT id, buyer_id FROM cart WHERE id = $1`
const getCartProductsSQL = `SELECT product_id as id, name, price, qty as units FROM product_line WHERE cart_id = $1`

// Get returns one cart by id.
func (r *CartRepository) Get(ID string) (*cart.Cart, error) {
	c := &cart.Cart{}
	row := r.db.QueryRow(getCartSQL, ID)
	err := row.Scan(&c.ID, &c.BuyerID)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(getCartProductsSQL, ID)
	if err != nil {
		return c, nil
	}
	defer rows.Close()

	for rows.Next() {
		p := &cart.Product{}
		err = rows.Scan(&p.ID, &p.Name, &p.Price, &p.Units)
		if err != nil {
			//break
			return c, err
		}
		c.AddProduct(p.ID, p.Name, float64(p.Price), int(p.Units))
	}

	return c, nil
}

// Save stores the cart.
func (r *CartRepository) Save(c cart.Cart) error {

	if r.exists(c.ID.String()) == true {
		return r.updateCartTx(c)
	}

	return r.createCartTx(c)
}

const existsCartSql = `SELECT c.ID FROM cart c WHERE c.ID = $1`
func (r *CartRepository) exists(ID string) bool {
	var cartId string
	row := r.db.QueryRow(existsCartSql, ID)
	err := row.Scan(&cartId)
	if err != nil {
		return false
	}
	return true
}

func (r *CartRepository) createCartTx(c cart.Cart) error {

	err := r.execTx(func(q *Queries) error {
		var err error

		err = q.insertCart(c)
		if err != nil {
			return err
		}

		for _, p := range c.Products {
			err = q.insertProduct(c.ID.String(), p)
			if err != nil {
				return err
			}
		}
		return err
	})
	return err
}

func (r *CartRepository) updateCartTx(c cart.Cart) error {

	err := r.execTx(func(q *Queries) error {
		var err error

		err = q.resetCart(c.ID.String())
		if err != nil {
			return err
		}

		for _, p := range c.Products {
			err = q.insertProduct(c.ID.String(), p)
			if err != nil {
				return err
			}
		}
		return err
	})
	return err
}



const insertCartSql = `
	INSERT INTO cart (id, buyer_id)
		VALUES ($1, $2)
	`
func (q *Queries) insertCart(c cart.Cart ) error {
	_, err := q.db.Exec(insertCartSql, c.ID, c.BuyerID)
	return err
}

const insertProductSql = `
	INSERT INTO product_line (product_id, cart_id, name, price, qty)
		VALUES ($1, $2, $3, $4, $5)
	`
func (q *Queries) insertProduct(cartID string, p cart.Product) error {
	_, err := q.db.Exec(insertProductSql, p.ID, cartID, p.Name, p.Price, p.Units)
	return err
}

const resetCartSql = `DELETE FROM product_line WHERE cart_id = $1`
func (q *Queries) resetCart(cartID string) error {
	_, err := q.db.Exec(resetCartSql, cartID)
	return err
}