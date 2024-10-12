package sessionstorage

import (
	"database/sql"
	"log"
)

func EnsureSessionsTableExists(db *sql.DB) {

	checkQuery := `
        SELECT name
    FROM sqlite_master WHERE type='table' AND name='sessions';
        SELECT name
    FROM sqlite_master WHERE type='index' AND name='sessions_expiry_idx';
    `

	rows, err := db.Query(checkQuery)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	if !rows.Next() {

		query := `
            CREATE table [sessions] (
        token char(43) primary key,
        data BLOB NOT NULL,
        expiry TIMESTAMP(6) NOT NULL
        );  
            CREATE INDEX sessions_expiry_idx ON sessions (expiry);
        `
		log.Println("Creating table and index for sessions")
		_, err := db.Exec(query)

		if err != nil {
			log.Fatal(err)
		}
		log.Println("Table and index sessions created")
	} else {
		log.Println("Table and index sessions already exists")
	}

}
