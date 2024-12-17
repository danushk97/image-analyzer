package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateImagesMetadataTable, downCreateImagesMetadataTable)
}

func upCreateImagesMetadataTable(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`CREATE TABLE images_metadata (
    	id UUID PRIMARY KEY,
    	user_id UUID NOT NULL,
    	filename VARCHAR(255) NOT NULL,
    	file_type VARCHAR(50) NOT NULL,
    	file_size BIGINT NOT NULL,
    	width INT NOT NULL,
    	height INT NOT NULL,
    	status VARCHAR(50) NOT NULL,
    	analysis_result TEXT,
    	created_at BIGINT NOT NULL,
    	updated_at BIGINT NOT NULL
	);`)

	return err
}

func downCreateImagesMetadataTable(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`DROP TABLE IF EXISTS images_metadata`)

	return err
}
