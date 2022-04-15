package sql

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
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
func Connect(cfg Config) (*sqlx.DB, error) {
	url := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s sslcert=%s sslkey=%s sslrootcert=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Name, cfg.Pass, cfg.SSLMode, cfg.SSLCert, cfg.SSLKey, cfg.SSLRootCert)

	db, err := sqlx.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	applied, err := migrateDB(db)
	if err != nil {
		fmt.Printf("Applied %d migrations!\n", applied)
		return nil, err
	}

	return db, nil
}

func migrateDB(db *sqlx.DB) (int, error) {
	up := []string{
		`CREATE TABLE IF NOT EXISTS pricing_plan (
			id	              VARCHAR(50),
			name              VARCHAR(255) UNIQUE,
			description		  VARCHAR(500),
			cycle             INTEGER,
			rate_type         INTEGER,
			default_price     INTEGER,
			max_number_device INTEGER,
			max_number_msg    INTEGER,
			unit_price    	  INTEGER,
			charging_unit     INTEGER,
			project_id    	  VARCHAR(255),
			PRIMARY KEY (id)
		)`,
		`CREATE TABLE IF NOT EXISTS user_payment (
			user_id 			 VARCHAR(100),
			statement_cycle_m 	 INTEGER,
			statement_cycle_y 	 INTEGER,
			amount 				 INTEGER,
			paid 				 INTEGER,
			payment_ts 			 TIME,
			statement_expired_ts TIME,
			PRIMARY KEY (user_id, statement_cycle_m, statement_cycle_y)
		)`,
	}
	down := []string{
		"DROP TABLE user_payment",
		"DROP TABLE pricing_plan",
	}
	migrations := &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			{Id: "Pricing_1", Up: up, Down: down},
		},
	}

	return migrate.Exec(db.DB, "postgres", migrations, migrate.Up)
}
