package products

import "github.com/jackc/pgx/v5/pgtype"

type CreateProductParams struct {
	ID           int64              `json:"id"`
	Name         string             `json:"name"`
	PriceInCents int32              `json:"price_in_cents"`
	Quantity     int32              `json:"quantity"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
}
