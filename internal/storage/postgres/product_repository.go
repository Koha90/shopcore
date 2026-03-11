// Package postgres provide methods for monipulate data postgresql.
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"botmanager/internal/domain"
)

// ProductRepository represent product repository.
type ProductRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

// NewProductRepository created a new product repository.
func NewProductRepository(db *sql.DB, logger *slog.Logger) *ProductRepository {
	return &ProductRepository{
		db:     db,
		logger: logger,
	}
}

// Save creates or updates product and all his variants in transaction.
func (r *ProductRepository) Save(ctx context.Context, p *domain.Product) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Error("failed to begin tx", "err", err)
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Insert product if new
	var productID int
	if p.ID() == 0 {
		err := tx.QueryRowContext(ctx,
			`INSERT INTO products (category_id, name, description, image_path, version)
		 	 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			p.CategoryID(), p.Name(), p.Description(), p.ImagePath(), p.Version(),
		).Scan(&productID)
		if err != nil {
			tx.Rollback()
			r.logger.Error("failed to insert product", "err", err)
			return err
		}
		p.SetID(productID)
		r.logger.Info("product created", "id", productID)
	} else {
		_, err := tx.ExecContext(
			ctx,
			`UPDATE products
		 	 SET category_id=$1, name=$2, description=$3, image_path=$4, version=$5
		   WHERE id=$6 AND version=$7`,
			p.CategoryID(),
			p.Name(),
			p.Description(),
			p.ImagePath(),
			p.Version(),
			p.ID(),
			p.Version(),
		)
		if err != nil {
			tx.Rollback()
			r.logger.Error("failed to update product", "err", err)
			return err
		}
		productID = p.ID()
		r.logger.Info("product updated", "id", productID)
	}

	// Insert/update variants
	for _, v := range p.VariantForUpdate() {
		if !v.IsActive() {
			continue
		}

		if v.ID() == 0 {
			// new variant
			var variantID int
			err := tx.QueryRowContext(ctx, `
				INSERT INTO product_variants
				(product_id, pack_size, district_id, price)
				VALUES ($1, $2, $3, $4)
				RETURNING id
			`, productID, v.PackSize(), v.DistrictID(), v.Price(),
			).Scan(&variantID)
			if err != nil {
				tx.Rollback()
				r.logger.Error("failed to insert variant", "productID", productID, "err", err)
				return err
			}
			v.SetID(variantID)
			r.logger.Info("variant created", "id", variantID)
		} else {
			// update existing
			_, err := tx.ExecContext(ctx, `
			  UPDATE product_variants
				SET pack_size=$1, district_id=$2, price=$3
				WHERE id=$4 
			`, v.PackSize(), v.DistrictID(), v.Price(), v.ID())
			if err != nil {
				tx.Rollback()
				r.logger.Error("failed to update variant", "id", v.ID(), "err", err)
				return err
			}

			r.logger.Info("variant updated", "id", v.ID())
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		r.logger.Error("failed to commit tx", "err", err)
		return err
	}

	return nil
}

// ByID load product and its variants.
func (r *ProductRepository) ByID(ctx context.Context, id int) (*domain.Product, error) {
	var (
		categoryID  int
		version     int
		name        string
		description string
		imgPath     sql.NullString
	)

	row := r.db.QueryRowContext(ctx,
		`SELECT id, category_id, name, description, image_path, version
	 	 FROM products WHERE id=$1`, id)
	if err := row.Scan(
		&id,
		&categoryID,
		&name,
		&description,
		&imgPath,
		&version,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrProductNotFound
		}
		r.logger.Error("failed to load product", "id", id, "err", err)
		return nil, err
	}

	var imgPtr *string
	if imgPath.Valid {
		imgPtr = &imgPath.String
	}

	variants, err := r.loadVariants(ctx, id)
	if err != nil {
		r.logger.Error("failed load variants", "err", err)
		return nil, err
	}

	product := domain.NewProductFromDB(
		id,
		&categoryID,
		name,
		description,
		imgPtr,
		version,
		variants,
	)

	return product, nil
}

func (r *ProductRepository) loadVariants(
	ctx context.Context,
	productID int,
) ([]domain.ProductVariant, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, pack_size, district_id, price, archived_at
		FROM product_variants
		WHERE product_id = $1
	`, productID)
	if err != nil {
		r.logger.Error("failed to query variants", "product_id", productID, "err", err)
		return nil, err
	}
	defer rows.Close()

	var variants []domain.ProductVariant

	for rows.Next() {
		var (
			id         int
			packSize   string
			districtID int
			price      int64
			archivedAt sql.NullTime
		)

		if err := rows.Scan(&id, &packSize, &districtID, &price, &archivedAt); err != nil {
			r.logger.Error("failed to scan variant", "product_id", productID, "err", err)
			return nil, err
		}

		var archivedPtr *time.Time
		if archivedAt.Valid {
			archivedPtr = &archivedAt.Time
		}

		v := domain.NewProductVariantFromDB(
			id, packSize, districtID, price, archivedPtr,
		)

		variants = append(variants, *v)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return variants, nil
}
