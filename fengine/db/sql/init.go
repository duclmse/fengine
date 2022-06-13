package sql

import (
	"fmt"

	_ "github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"

	"github.com/duclmse/fengine/pkg/logger"
)

// Config define the options that are used when connecting to a Postgres instance
type Config struct {
	Host        string
	Port        string
	User        string
	Pass        string
	Name        string
	SSLMode     string
	SSLCert     string
	SSLKey      string
	SSLRootCert string
}

// Connect method is used to create a connection to the Postgres instance and applies any unapply database migrations.
// A non-nil error is return to indicate failure
func Connect(cfg Config, log logger.Logger) (*sqlx.DB, error) {
	url := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s sslcert=%s sslkey=%s sslrootcert=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Name, cfg.Pass, cfg.SSLMode, cfg.SSLCert, cfg.SSLKey, cfg.SSLRootCert)
	log.Info("db info: %s", url)

	db, err := sqlx.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	applied, err := migrateDB(db)
	if err == nil {
		log.Info("Applied %d migrations!", applied)
		return db, nil
	} else {
		log.Info("Error applying migrations: %s", err.Error())
		return nil, err
	}
}

func migrateDB(db *sqlx.DB) (int, error) {
	up := []string{
		// language=postgresql
		`DO $$BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'var_type') THEN
				CREATE TYPE VAR_TYPE AS ENUM ('i32', 'i64', 'f32', 'f64', 'bool', 'json', 'string', 'binary');
			END IF;
		END$$;`,
		// language=postgresql
		`DO $$BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'entity_type') THEN
				CREATE TYPE ENTITY_TYPE AS ENUM ('shape', 'template', 'thing');
			END IF;
		END$$;`,
		// language=postgresql
		`CREATE TABLE IF NOT EXISTS "entity" (
			"id"            UUID NOT NULL,
			"name"          VARCHAR(255) NOT NULL,
			"type"          ENTITY_TYPE  NOT NULL,
			"description"   VARCHAR(500),
			"project_id"    UUID,
			"base_template" UUID,
			"base_shapes"   UUID[],
			"create_ts" TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
			"update_ts" TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL,
			PRIMARY KEY (id)
		);`,
		// language=postgresql
		`CREATE TABLE IF NOT EXISTS "attribute" (
			"entity_id"    UUID NOT NULL,
			"name"         VARCHAR(255) NOT NULL,
			"type"         VAR_TYPE NOT NULL,
			"from"         UUID,
			"value_i32"    INT4,
			"value_i64"    INT4,
			"value_f32"    FLOAT4,
			"value_f64"    FLOAT8,
			"value_bool"   BOOLEAN,
			"value_json"   JSONB,
			"value_string" TEXT,
			"value_binary" BYTEA,
			"create_ts" TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
			"update_ts" TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL,
			PRIMARY KEY ("entity_id", "name"),
			FOREIGN KEY ("entity_id") REFERENCES entity ("id"),
			FOREIGN KEY ("from") REFERENCES entity ("id")
		);`,
		// language=postgresql
		`CREATE TABLE IF NOT EXISTS "service" (
			"entity_id" UUID NOT NULL,
			"name"      VARCHAR(255) NOT NULL,
			"input"     JSONB,
			"output"    VAR_TYPE,
			"from"      UUID,
			"code"      TEXT,
			"create_ts" TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
			"update_ts" TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL,
			PRIMARY KEY ("entity_id", "name"),
			FOREIGN KEY ("entity_id") REFERENCES entity ("id"),
			FOREIGN KEY ("from") REFERENCES entity ("id")
		);`,
		// language=postgresql
		`CREATE TABLE IF NOT EXISTS "subscription" (
			"entity_id" UUID NOT NULL,
			"name"      VARCHAR(255) NOT NULL,
			"subs_on"   VARCHAR(50),
			"event"     VARCHAR(50),
			"from"      UUID,
			"code"      TEXT,
			"create_ts" TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
			"update_ts" TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL,
			PRIMARY KEY ("entity_id", "name"),
			FOREIGN KEY ("entity_id") REFERENCES entity ("id"),
			FOREIGN KEY ("from") REFERENCES entity ("id")
		);`,
	}
	down := []string{
		`DROP TABLE "method";`,
		`DROP TABLE "attribute";`,
		`DROP TABLE "entity";`,
		`DROP TYPE  "var_type";`,
		`DROP TYPE  "entity_type";`,
		`DROP TYPE  "method_type";`,
	}
	migrations := &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			{Id: "fengine_1", Up: up, Down: down},
		},
	}

	return migrate.Exec(db.DB, "postgres", migrations, migrate.Up)
}
