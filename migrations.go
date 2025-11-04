package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Migration struct {
	ID      int
	Name    string
	UpSQL   string
	DownSQL string
}

var migrations = []Migration{
	{
		ID:    1,
		Name:  "initial_schema",
		UpSQL: initialSchemaSQL,
	},
	{
		ID:   2,
		Name: "add_proxy_url",
		UpSQL: `
            -- PostgreSQL version
            DO $$
            BEGIN
                IF NOT EXISTS (
                    SELECT 1 FROM information_schema.columns 
                    WHERE table_name = 'users' AND column_name = 'proxy_url'
                ) THEN
                    ALTER TABLE users ADD COLUMN proxy_url TEXT DEFAULT '';
                END IF;
            END $$;
            
            -- SQLite version (handled in code)
            `,
	},
	{
		ID:    3,
		Name:  "change_id_to_string",
		UpSQL: changeIDToStringSQL,
	},
	{
		ID:    4,
		Name:  "add_s3_support",
		UpSQL: addS3SupportSQL,
	},
	{
		ID:    5,
		Name:  "add_message_history",
		UpSQL: addMessageHistorySQL,
	},
	{
		ID:    6,
		Name:  "add_quoted_message_id",
		UpSQL: addQuotedMessageIDSQL,
	},
	{
		ID:    7,
		Name:  "add_hmac_key",
		UpSQL: addHmacKeySQL,
	},
	{
		ID:    8,
		Name:  "add_data_json",
		UpSQL: addDataJsonSQL,
	},
	{
		ID:    9,
		Name:  "add_system_users",
		UpSQL: addSystemUsersSQL,
	},
	{
		ID:    10,
		Name:  "add_user_id_to_instances",
		UpSQL: addUserIDToInstancesSQL,
	},
	{
		ID:    11,
		Name:  "add_destination_number",
		UpSQL: addDestinationNumberSQL,
	},
	{
		ID:    12,
		Name:  "add_daily_conversations",
		UpSQL: addDailyConversationsSQL,
	},
	{
		ID:    13,
		Name:  "add_subscription_plans",
		UpSQL: addSubscriptionPlansSQL,
	},
	{
		ID:    14,
		Name:  "add_system_user_profile_fields",
		UpSQL: addSystemUserProfileFieldsSQL,
	},
}

const changeIDToStringSQL = `
-- Migration to change ID from integer to random string
DO $$
BEGIN
    -- Only execute if the column is currently integer type
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'users' AND column_name = 'id' AND data_type = 'integer'
    ) THEN
        -- For PostgreSQL
        ALTER TABLE users ADD COLUMN new_id TEXT;
		UPDATE users SET new_id = md5(random()::text || id::text || clock_timestamp()::text);
		ALTER TABLE users DROP CONSTRAINT users_pkey;
        ALTER TABLE users DROP COLUMN id CASCADE;
        ALTER TABLE users RENAME COLUMN new_id TO id;
        ALTER TABLE users ALTER COLUMN id SET NOT NULL;
        ALTER TABLE users ADD PRIMARY KEY (id);
    END IF;
END $$;
`

const initialSchemaSQL = `
-- PostgreSQL version
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'users') THEN
        CREATE TABLE users (
            id TEXT PRIMARY KEY,
            name TEXT NOT NULL,
            token TEXT NOT NULL,
            webhook TEXT NOT NULL DEFAULT '',
            jid TEXT NOT NULL DEFAULT '',
            qrcode TEXT NOT NULL DEFAULT '',
            connected INTEGER,
            expiration INTEGER,
            events TEXT NOT NULL DEFAULT '',
            proxy_url TEXT DEFAULT ''
        );
    END IF;
END $$;

-- SQLite version (handled in code)
`

const addS3SupportSQL = `
-- PostgreSQL version
DO $$
BEGIN
    -- Add S3 configuration columns if they don't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 's3_enabled') THEN
        ALTER TABLE users ADD COLUMN s3_enabled BOOLEAN DEFAULT FALSE;
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 's3_endpoint') THEN
        ALTER TABLE users ADD COLUMN s3_endpoint TEXT DEFAULT '';
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 's3_region') THEN
        ALTER TABLE users ADD COLUMN s3_region TEXT DEFAULT '';
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 's3_bucket') THEN
        ALTER TABLE users ADD COLUMN s3_bucket TEXT DEFAULT '';
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 's3_access_key') THEN
        ALTER TABLE users ADD COLUMN s3_access_key TEXT DEFAULT '';
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 's3_secret_key') THEN
        ALTER TABLE users ADD COLUMN s3_secret_key TEXT DEFAULT '';
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 's3_path_style') THEN
        ALTER TABLE users ADD COLUMN s3_path_style BOOLEAN DEFAULT TRUE;
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 's3_public_url') THEN
        ALTER TABLE users ADD COLUMN s3_public_url TEXT DEFAULT '';
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'media_delivery') THEN
        ALTER TABLE users ADD COLUMN media_delivery TEXT DEFAULT 'base64';
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 's3_retention_days') THEN
        ALTER TABLE users ADD COLUMN s3_retention_days INTEGER DEFAULT 30;
    END IF;
END $$;
`

const addMessageHistorySQL = `
-- PostgreSQL version
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'message_history') THEN
        CREATE TABLE message_history (
            id SERIAL PRIMARY KEY,
            user_id TEXT NOT NULL,
            chat_jid TEXT NOT NULL,
            sender_jid TEXT NOT NULL,
            message_id TEXT NOT NULL,
            timestamp TIMESTAMP NOT NULL,
            message_type TEXT NOT NULL,
            text_content TEXT,
            media_link TEXT,
            UNIQUE(user_id, message_id)
        );
        CREATE INDEX idx_message_history_user_chat_timestamp ON message_history (user_id, chat_jid, timestamp DESC);
    END IF;
    
    -- Add history column to users table if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'history') THEN
        ALTER TABLE users ADD COLUMN history INTEGER DEFAULT 0;
    END IF;
END $$;
`

const addQuotedMessageIDSQL = `
-- PostgreSQL version
DO $$
BEGIN
    -- Add quoted_message_id column to message_history table if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'message_history' AND column_name = 'quoted_message_id') THEN
        ALTER TABLE message_history ADD COLUMN quoted_message_id TEXT;
    END IF;
END $$;
`

const addDataJsonSQL = `
-- PostgreSQL version
DO $$
BEGIN
    -- Add dataJson column to message_history table if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'message_history' AND column_name = 'datajson') THEN
        ALTER TABLE message_history ADD COLUMN datajson TEXT;
    END IF;
END $$;

-- SQLite version (handled in code)
`

// GenerateRandomID creates a random string ID
func GenerateRandomID() (string, error) {
	bytes := make([]byte, 16) // 128 bits
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random ID: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// Initialize the database with migrations
func initializeSchema(db *sqlx.DB) error {
	// Create migrations table if it doesn't exist
	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get already applied migrations
	applied, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Apply missing migrations
	for _, migration := range migrations {
		if _, ok := applied[migration.ID]; !ok {
			if err := applyMigration(db, migration); err != nil {
				return fmt.Errorf("failed to apply migration %d: %w", migration.ID, err)
			}
		}
	}

	return nil
}

func createMigrationsTable(db *sqlx.DB) error {
	var tableExists bool
	var err error

	switch db.DriverName() {
	case "postgres":
		err = db.Get(&tableExists, `
			SELECT EXISTS (
				SELECT 1 FROM information_schema.tables 
				WHERE table_name = 'migrations'
			)`)
	case "sqlite":
		err = db.Get(&tableExists, `
			SELECT EXISTS (
				SELECT 1 FROM sqlite_master 
				WHERE type='table' AND name='migrations'
			)`)
	default:
		return fmt.Errorf("unsupported database driver: %s", db.DriverName())
	}

	if err != nil {
		return fmt.Errorf("failed to check migrations table existence: %w", err)
	}

	if tableExists {
		return nil
	}

	_, err = db.Exec(`
		CREATE TABLE migrations (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	return nil
}

func getAppliedMigrations(db *sqlx.DB) (map[int]struct{}, error) {
	applied := make(map[int]struct{})
	var rows []struct {
		ID   int    `db:"id"`
		Name string `db:"name"`
	}

	err := db.Select(&rows, "SELECT id, name FROM migrations ORDER BY id ASC")
	if err != nil {
		return nil, fmt.Errorf("failed to query applied migrations: %w", err)
	}

	for _, row := range rows {
		applied[row.ID] = struct{}{}
	}

	return applied, nil
}

func applyMigration(db *sqlx.DB, migration Migration) error {
	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if migration.ID == 1 {
		// Handle initial schema creation differently per database
		if db.DriverName() == "sqlite" {
			err = createTableIfNotExistsSQLite(tx, "users", `
                CREATE TABLE users (
                    id TEXT PRIMARY KEY,
                    name TEXT NOT NULL,
                    token TEXT NOT NULL,
                    webhook TEXT NOT NULL DEFAULT '',
                    jid TEXT NOT NULL DEFAULT '',
                    qrcode TEXT NOT NULL DEFAULT '',
                    connected INTEGER,
                    expiration INTEGER,
                    events TEXT NOT NULL DEFAULT '',
                    proxy_url TEXT DEFAULT ''
                )`)
		} else {
			_, err = tx.Exec(migration.UpSQL)
		}
	} else if migration.ID == 2 {
		if db.DriverName() == "sqlite" {
			err = addColumnIfNotExistsSQLite(tx, "users", "proxy_url", "TEXT DEFAULT ''")
		} else {
			_, err = tx.Exec(migration.UpSQL)
		}
	} else if migration.ID == 3 {
		if db.DriverName() == "sqlite" {
			err = migrateSQLiteIDToString(tx)
		} else {
			_, err = tx.Exec(migration.UpSQL)
		}
	} else if migration.ID == 4 {
		if db.DriverName() == "sqlite" {
			// Handle S3 columns for SQLite
			err = addColumnIfNotExistsSQLite(tx, "users", "s3_enabled", "BOOLEAN DEFAULT 0")
			if err == nil {
				err = addColumnIfNotExistsSQLite(tx, "users", "s3_endpoint", "TEXT DEFAULT ''")
			}
			if err == nil {
				err = addColumnIfNotExistsSQLite(tx, "users", "s3_region", "TEXT DEFAULT ''")
			}
			if err == nil {
				err = addColumnIfNotExistsSQLite(tx, "users", "s3_bucket", "TEXT DEFAULT ''")
			}
			if err == nil {
				err = addColumnIfNotExistsSQLite(tx, "users", "s3_access_key", "TEXT DEFAULT ''")
			}
			if err == nil {
				err = addColumnIfNotExistsSQLite(tx, "users", "s3_secret_key", "TEXT DEFAULT ''")
			}
			if err == nil {
				err = addColumnIfNotExistsSQLite(tx, "users", "s3_path_style", "BOOLEAN DEFAULT 1")
			}
			if err == nil {
				err = addColumnIfNotExistsSQLite(tx, "users", "s3_public_url", "TEXT DEFAULT ''")
			}
			if err == nil {
				err = addColumnIfNotExistsSQLite(tx, "users", "media_delivery", "TEXT DEFAULT 'base64'")
			}
			if err == nil {
				err = addColumnIfNotExistsSQLite(tx, "users", "s3_retention_days", "INTEGER DEFAULT 30")
			}
		} else {
			_, err = tx.Exec(migration.UpSQL)
		}
	} else if migration.ID == 5 {
		if db.DriverName() == "sqlite" {
			// Handle message_history table creation for SQLite
			err = createTableIfNotExistsSQLite(tx, "message_history", `
				CREATE TABLE message_history (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					user_id TEXT NOT NULL,
					chat_jid TEXT NOT NULL,
					sender_jid TEXT NOT NULL,
					message_id TEXT NOT NULL,
					timestamp DATETIME NOT NULL,
					message_type TEXT NOT NULL,
					text_content TEXT,
					media_link TEXT,
					UNIQUE(user_id, message_id)
				)`)
			if err == nil {
				// Create index for SQLite
				_, err = tx.Exec(`
					CREATE INDEX IF NOT EXISTS idx_message_history_user_chat_timestamp 
					ON message_history (user_id, chat_jid, timestamp DESC)`)
			}
			if err == nil {
				// Add history column to users table
				err = addColumnIfNotExistsSQLite(tx, "users", "history", "INTEGER DEFAULT 0")
			}
		} else {
			_, err = tx.Exec(migration.UpSQL)
		}
	} else if migration.ID == 6 {
		if db.DriverName() == "sqlite" {
			// Add quoted_message_id column to message_history table for SQLite
			err = addColumnIfNotExistsSQLite(tx, "message_history", "quoted_message_id", "TEXT")
		} else {
			_, err = tx.Exec(migration.UpSQL)
		}
	} else if migration.ID == 7 {
		if db.DriverName() == "sqlite" {
			// Add hmac_key column as BLOB for encrypted data in SQLite
			err = addColumnIfNotExistsSQLite(tx, "users", "hmac_key", "BLOB")
		} else {
			_, err = tx.Exec(migration.UpSQL)
		}
	} else if migration.ID == 8 {
		if db.DriverName() == "sqlite" {
			// Add dataJson column to message_history table for SQLite
			err = addColumnIfNotExistsSQLite(tx, "message_history", "datajson", "TEXT")
		} else {
			_, err = tx.Exec(migration.UpSQL)
		}
	} else if migration.ID == 9 {
		if db.DriverName() == "sqlite" {
			// Create system_users table for SQLite
			err = createTableIfNotExistsSQLite(tx, "system_users", `
				CREATE TABLE system_users (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					email TEXT UNIQUE NOT NULL,
					password_hash TEXT NOT NULL,
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
				)`)
		} else {
			_, err = tx.Exec(migration.UpSQL)
		}
	} else if migration.ID == 10 {
		if db.DriverName() == "sqlite" {
			// Add system_user_id column to users table for SQLite
			err = addColumnIfNotExistsSQLite(tx, "users", "system_user_id", "INTEGER REFERENCES system_users(id) ON DELETE CASCADE")
		} else {
			_, err = tx.Exec(migration.UpSQL)
		}
	} else if migration.ID == 11 {
		if db.DriverName() == "sqlite" {
			// Add destination_number column to users table for SQLite
			err = addColumnIfNotExistsSQLite(tx, "users", "destination_number", "TEXT DEFAULT ''")
		} else {
			_, err = tx.Exec(migration.UpSQL)
		}
	} else if migration.ID == 12 {
		if db.DriverName() == "sqlite" {
			// Create daily_conversations table for SQLite
			err = createTableIfNotExistsSQLite(tx, "daily_conversations", `
				CREATE TABLE daily_conversations (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					user_id TEXT NOT NULL,
					date DATE NOT NULL,
					chat_jid TEXT NOT NULL,
					contact TEXT NOT NULL,
					messages TEXT NOT NULL,
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					UNIQUE(user_id, date, chat_jid)
				)`)
			if err == nil {
				// Create index for SQLite
				_, err = tx.Exec(`
					CREATE INDEX IF NOT EXISTS idx_daily_conversations_user_date 
					ON daily_conversations (user_id, date)`)
			}
		} else {
			_, err = tx.Exec(migration.UpSQL)
		}
	} else if migration.ID == 13 {
		if db.DriverName() == "sqlite" {
			// Create plans table for SQLite
			err = createTableIfNotExistsSQLite(tx, "plans", `
				CREATE TABLE plans (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					name TEXT NOT NULL,
					price REAL NOT NULL,
					max_instances INTEGER NOT NULL,
					trial_days INTEGER DEFAULT 0,
					is_active BOOLEAN DEFAULT 1,
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP
				)`)
			if err == nil {
				// Insert default plans
				_, err = tx.Exec(`
					INSERT OR IGNORE INTO plans (id, name, price, max_instances, trial_days) VALUES
					(1, 'Gratuito', 0.00, 999999, 5),
					(2, 'Pro', 29.00, 5, 0),
					(3, 'Analista', 97.00, 12, 0)`)
			}
			if err == nil {
				// Create user_subscriptions table
				err = createTableIfNotExistsSQLite(tx, "user_subscriptions", `
					CREATE TABLE user_subscriptions (
						id INTEGER PRIMARY KEY AUTOINCREMENT,
						system_user_id INTEGER NOT NULL REFERENCES system_users(id) ON DELETE CASCADE,
						plan_id INTEGER NOT NULL REFERENCES plans(id),
						started_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
						expires_at DATETIME,
						is_active BOOLEAN DEFAULT 1,
						created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
						updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
					)`)
			}
			if err == nil {
				// Create indexes
				_, err = tx.Exec(`
					CREATE INDEX IF NOT EXISTS idx_user_subscriptions_user 
					ON user_subscriptions (system_user_id)`)
			}
			if err == nil {
				_, err = tx.Exec(`
					CREATE INDEX IF NOT EXISTS idx_user_subscriptions_active 
					ON user_subscriptions (system_user_id, is_active, expires_at)`)
			}
			if err == nil {
				// Create subscription_history table
				err = createTableIfNotExistsSQLite(tx, "subscription_history", `
					CREATE TABLE subscription_history (
						id INTEGER PRIMARY KEY AUTOINCREMENT,
						system_user_id INTEGER NOT NULL REFERENCES system_users(id) ON DELETE CASCADE,
						plan_id INTEGER NOT NULL REFERENCES plans(id),
						started_at DATETIME NOT NULL,
						ended_at DATETIME,
						created_at DATETIME DEFAULT CURRENT_TIMESTAMP
					)`)
			}
			if err == nil {
				// Create index
				_, err = tx.Exec(`
					CREATE INDEX IF NOT EXISTS idx_subscription_history_user 
					ON subscription_history (system_user_id)`)
			}
		} else {
			_, err = tx.Exec(migration.UpSQL)
		}
	} else if migration.ID == 14 {
		if db.DriverName() == "sqlite" {
			// Add name and whatsapp_number to system_users table for SQLite
			err = addColumnIfNotExistsSQLite(tx, "system_users", "name", "TEXT DEFAULT ''")
			if err == nil {
				err = addColumnIfNotExistsSQLite(tx, "system_users", "whatsapp_number", "TEXT DEFAULT ''")
			}
		} else {
			_, err = tx.Exec(migration.UpSQL)
		}
	} else {
		_, err = tx.Exec(migration.UpSQL)
	}

	if err != nil {
		return fmt.Errorf("failed to execute migration SQL: %w", err)
	}

	// Record the migration
	if _, err = tx.Exec(`
        INSERT INTO migrations (id, name) 
        VALUES ($1, $2)`, migration.ID, migration.Name); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	return tx.Commit()
}

func createTableIfNotExistsSQLite(tx *sqlx.Tx, tableName, createSQL string) error {
	var exists int
	err := tx.Get(&exists, `
        SELECT COUNT(*) FROM sqlite_master
        WHERE type='table' AND name=?`, tableName)
	if err != nil {
		return err
	}

	if exists == 0 {
		_, err = tx.Exec(createSQL)
		return err
	}
	return nil
}
func sqliteChangeIDType(tx *sqlx.Tx) error {
	// SQLite requires a more complex approach:
	// 1. Create new table with string ID
	// 2. Copy data with new UUIDs
	// 3. Drop old table
	// 4. Rename new table

	// Step 1: Get the current schema
	var tableInfo string
	err := tx.Get(&tableInfo, `
        SELECT sql FROM sqlite_master
        WHERE type='table' AND name='users'`)
	if err != nil {
		return fmt.Errorf("failed to get table info: %w", err)
	}

	// Step 2: Create new table with string ID
	newTableSQL := strings.Replace(tableInfo,
		"CREATE TABLE users (",
		"CREATE TABLE users_new (id TEXT PRIMARY KEY, ", 1)
	newTableSQL = strings.Replace(newTableSQL,
		"id INTEGER PRIMARY KEY AUTOINCREMENT,", "", 1)

	if _, err = tx.Exec(newTableSQL); err != nil {
		return fmt.Errorf("failed to create new table: %w", err)
	}

	// Step 3: Copy data with new UUIDs
	columns, err := getTableColumns(tx, "users")
	if err != nil {
		return fmt.Errorf("failed to get table columns: %w", err)
	}

	// Remove 'id' from columns list
	var filteredColumns []string
	for _, col := range columns {
		if col != "id" {
			filteredColumns = append(filteredColumns, col)
		}
	}

	columnList := strings.Join(filteredColumns, ", ")
	if _, err = tx.Exec(fmt.Sprintf(`
        INSERT INTO users_new (id, %s)
        SELECT gen_random_uuid(), %s FROM users`,
		columnList, columnList)); err != nil {
		return fmt.Errorf("failed to copy data: %w", err)
	}

	// Step 4: Drop old table
	if _, err = tx.Exec("DROP TABLE users"); err != nil {
		return fmt.Errorf("failed to drop old table: %w", err)
	}

	// Step 5: Rename new table
	if _, err = tx.Exec("ALTER TABLE users_new RENAME TO users"); err != nil {
		return fmt.Errorf("failed to rename table: %w", err)
	}

	return nil
}

func getTableColumns(tx *sqlx.Tx, tableName string) ([]string, error) {
	var columns []string
	rows, err := tx.Query(fmt.Sprintf("PRAGMA table_info(%s)", tableName))
	if err != nil {
		return nil, fmt.Errorf("failed to get table info: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var name, typ string
		var notnull int
		var dfltValue interface{}
		var pk int

		if err := rows.Scan(&cid, &name, &typ, &notnull, &dfltValue, &pk); err != nil {
			return nil, fmt.Errorf("failed to scan column info: %w", err)
		}
		columns = append(columns, name)
	}

	return columns, nil
}

func migrateSQLiteIDToString(tx *sqlx.Tx) error {
	// 1. Check if we need to do the migration
	var currentType string
	err := tx.QueryRow(`
        SELECT type FROM pragma_table_info('users')
        WHERE name = 'id'`).Scan(&currentType)
	if err != nil {
		return fmt.Errorf("failed to check column type: %w", err)
	}

	if currentType != "INTEGER" {
		// No migration needed
		return nil
	}

	// 2. Create new table with string ID
	_, err = tx.Exec(`
        CREATE TABLE users_new (
            id TEXT PRIMARY KEY,
            name TEXT NOT NULL,
            token TEXT NOT NULL,
            webhook TEXT NOT NULL DEFAULT '',
            jid TEXT NOT NULL DEFAULT '',
            qrcode TEXT NOT NULL DEFAULT '',
            connected INTEGER,
            expiration INTEGER,
            events TEXT NOT NULL DEFAULT '',
            proxy_url TEXT DEFAULT ''
        )`)
	if err != nil {
		return fmt.Errorf("failed to create new table: %w", err)
	}

	// 3. Copy data with new UUIDs
	_, err = tx.Exec(`
        INSERT INTO users_new
        SELECT
            hex(randomblob(16)),
            name, token, webhook, jid, qrcode,
            connected, expiration, events, proxy_url 
        FROM users`)
	if err != nil {
		return fmt.Errorf("failed to copy data: %w", err)
	}

	// 4. Drop old table
	_, err = tx.Exec(`DROP TABLE users`)
	if err != nil {
		return fmt.Errorf("failed to drop old table: %w", err)
	}

	// 5. Rename new table
	_, err = tx.Exec(`ALTER TABLE users_new RENAME TO users`)
	if err != nil {
		return fmt.Errorf("failed to rename table: %w", err)
	}

	return nil
}

func addColumnIfNotExistsSQLite(tx *sqlx.Tx, tableName, columnName, columnDef string) error {
	var exists int
	err := tx.Get(&exists, `
        SELECT COUNT(*) FROM pragma_table_info(?)
        WHERE name = ?`, tableName, columnName)
	if err != nil {
		return fmt.Errorf("failed to check column existence: %w", err)
	}

	if exists == 0 {
		_, err = tx.Exec(fmt.Sprintf(
			"ALTER TABLE %s ADD COLUMN %s %s",
			tableName, columnName, columnDef))
		if err != nil {
			return fmt.Errorf("failed to add column: %w", err)
		}
	}
	return nil
}

const addHmacKeySQL = `
-- PostgreSQL version - Add encrypted HMAC key column
DO $$
BEGIN
    -- Add hmac_key column as BYTEA for encrypted data
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'hmac_key') THEN
        ALTER TABLE users ADD COLUMN hmac_key BYTEA;
    END IF;
END $$;

-- SQLite version (handled in code)
`

const addSystemUsersSQL = `
-- PostgreSQL version - Create system_users table
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'system_users') THEN
        CREATE TABLE system_users (
            id SERIAL PRIMARY KEY,
            email TEXT UNIQUE NOT NULL,
            password_hash TEXT NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
    END IF;
END $$;

-- SQLite version (handled in code)
`

const addUserIDToInstancesSQL = `
-- PostgreSQL version - Add system_user_id to users (instances) table
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'system_user_id') THEN
        ALTER TABLE users ADD COLUMN system_user_id INTEGER REFERENCES system_users(id) ON DELETE CASCADE;
    END IF;
END $$;

-- SQLite version (handled in code)
`

const addDestinationNumberSQL = `
-- PostgreSQL version - Add destination number field
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'destination_number') THEN
        ALTER TABLE users ADD COLUMN destination_number TEXT DEFAULT '';
    END IF;
END $$;

-- SQLite version (handled in code)
`

const addDailyConversationsSQL = `
-- PostgreSQL version - Create table for daily conversation cache
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'daily_conversations') THEN
        CREATE TABLE daily_conversations (
            id SERIAL PRIMARY KEY,
            user_id TEXT NOT NULL,
            date DATE NOT NULL,
            chat_jid TEXT NOT NULL,
            contact TEXT NOT NULL,
            messages JSONB NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            UNIQUE(user_id, date, chat_jid)
        );
        CREATE INDEX idx_daily_conversations_user_date ON daily_conversations (user_id, date);
    END IF;
END $$;

-- SQLite version (handled in code)
`

const addSubscriptionPlansSQL = `
-- PostgreSQL version - Create subscription plans tables
DO $$
BEGIN
    -- Plans table
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'plans') THEN
        CREATE TABLE plans (
            id SERIAL PRIMARY KEY,
            name TEXT NOT NULL,
            price DECIMAL(10, 2) NOT NULL,
            max_instances INTEGER NOT NULL,
            trial_days INTEGER DEFAULT 0,
            is_active BOOLEAN DEFAULT TRUE,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
        
        -- Insert default plans
        INSERT INTO plans (name, price, max_instances, trial_days) VALUES
        ('Gratuito', 0.00, 999999, 5),
        ('Pro', 29.00, 5, 0),
        ('Analista', 97.00, 12, 0);
    END IF;
    
    -- User subscriptions table
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'user_subscriptions') THEN
        CREATE TABLE user_subscriptions (
            id SERIAL PRIMARY KEY,
            system_user_id INTEGER NOT NULL REFERENCES system_users(id) ON DELETE CASCADE,
            plan_id INTEGER NOT NULL REFERENCES plans(id),
            started_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
            expires_at TIMESTAMP,
            is_active BOOLEAN DEFAULT TRUE,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
        CREATE INDEX idx_user_subscriptions_user ON user_subscriptions (system_user_id);
        CREATE INDEX idx_user_subscriptions_active ON user_subscriptions (system_user_id, is_active, expires_at);
    END IF;
    
    -- Subscription history table
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'subscription_history') THEN
        CREATE TABLE subscription_history (
            id SERIAL PRIMARY KEY,
            system_user_id INTEGER NOT NULL REFERENCES system_users(id) ON DELETE CASCADE,
            plan_id INTEGER NOT NULL REFERENCES plans(id),
            started_at TIMESTAMP NOT NULL,
            ended_at TIMESTAMP,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
        CREATE INDEX idx_subscription_history_user ON subscription_history (system_user_id);
    END IF;
END $$;

-- SQLite version (handled in code)
`

const addSystemUserProfileFieldsSQL = `
-- PostgreSQL version - Add name and whatsapp_number to system_users
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'system_users' AND column_name = 'name') THEN
        ALTER TABLE system_users ADD COLUMN name TEXT DEFAULT '';
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'system_users' AND column_name = 'whatsapp_number') THEN
        ALTER TABLE system_users ADD COLUMN whatsapp_number TEXT DEFAULT '';
    END IF;
END $$;

-- SQLite version (handled in code)
`
