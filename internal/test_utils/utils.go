package utils

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

type TestEnv struct {
	Context context.Context
	DB      *sql.DB
	Cleanup func()
}

func CreateTestEnv() (TestEnv, error) {

	ctx := context.Background()
	testDb, cleanup, err := CreateTestDB(ctx)
	if err != nil {
		return TestEnv{}, err
	}

	testEnv := TestEnv{
		Context: ctx,
		DB:      testDb,
		Cleanup: cleanup,
	}

	return testEnv, nil
}

func startContainer(pool *dockertest.Pool) (string, func(), error) {
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13-alpine",
		Env: []string{
			"POSTGRES_USER=postgres",
			"POSTGRES_PASSWORD=postgres",
			"POSTGRES_DB=everphone",
		},
		NetworkID: "default",
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
	})

	if err != nil {
		return "", nil, err
	}

	cleanup := func() {
		if err := pool.Purge(resource); err != nil {
			log.Printf("Could not purge resource: %s", err)
		}
	}

	postgresHost := os.Getenv("POSTGRES_HOST")
	if postgresHost == "" {
		postgresHost = "localhost"
	}

	connString := fmt.Sprintf("host=%s port=%s user=postgres password=postgres dbname=everphone sslmode=disable", postgresHost, resource.GetPort("5432/tcp"))
	// Wait for the container to be ready
	err = pool.Retry(func() error {
		db, err := sql.Open("postgres", connString)
		if err != nil {
			return err
		}
		return db.Ping()
	})

	if err != nil {
		cleanup()
		return "", nil, err
	}

	return connString, cleanup, nil
}

func CreateTestDB(ctx context.Context) (*sql.DB, func(), error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, err
	}

	connString, cleanup, err := startContainer(pool)
	if err != nil {
		return nil, nil, err
	}

	testDb, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, nil, err
	}

	testDb.SetMaxOpenConns(50)                  // Maximum number of open connections to the database
	testDb.SetMaxIdleConns(10)                  // Maximum number of connections in the idle connection pool
	testDb.SetConnMaxLifetime(30 * time.Minute) // Maximum amount of time a connection may be reused

	// Load and execute the SQL script to create the required tables and insert test data
	script := `CREATE TABLE IF NOT EXISTS employees (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			interests VARCHAR(255)[],
			gift VARCHAR(255),
			returned BOOLEAN DEFAULT false
			);

		CREATE TABLE IF NOT EXISTS gifts (
		 id SERIAL PRIMARY KEY,
		 name VARCHAR(255) NOT NULL,
		categories VARCHAR(255)[],
		istaken BOOLEAN DEFAULT false,
		 returned BOOLEAN DEFAULT false
		);

		INSERT INTO employees (name, interests) VALUES
			('severin', ARRAY['high end audio', 'photography', 'sitting comfortably']),
			('oliver', ARRAY['pets', 'scifi', 'music making']),
			('henry', ARRAY['scifi', 'music', 'drinking', 'motorcross']),
			('karsten', ARRAY['family', 'reading']),
			('mehmet', ARRAY['eating', 'sleeping', 'power lifting']);

		INSERT INTO gifts (name, categories) VALUES
			('netflix card', ARRAY['music', 'film', 'pets', 'scifi']),
			('gold plated water bottle', ARRAY['power lifting', 'triathlons', 'football', 'crossfit', 'handball', 'running']),
			('second hand earplugs', ARRAY['techno', 'music']),
			('6 pack of beer', ARRAY['drinking']),
			('peterborough united football shirt', ARRAY['football']);`

	_, err = testDb.Exec(script)
	if err != nil {
		log.Fatalf("Could not execute SQL script: %s", err)
	}

	return testDb, cleanup, nil
}
