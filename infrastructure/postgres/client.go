package postgres

import (
	"auth/config/pg"
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

const (
	_defaultConnAttempts               = 10
	_defaultConnTimeout  time.Duration = 10 * time.Second
)

var (
	ErrNoChange = errors.New("no changes applied")
)

type Client struct {
	Builder squirrel.StatementBuilderType
	Pool    *pgxpool.Pool
	cfg     pg.Config
}

func NewClient(config pg.Config) (*Client, error) {
	client := &Client{cfg: config}

	connAttempts := config.RetryConnectionAttempts
	connTimeout := config.RetryConnectionTimeout

	if connAttempts == 0 {
		connAttempts = _defaultConnAttempts
	}
	if connTimeout == 0 {
		connTimeout = _defaultConnTimeout
	}

	client.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	connectionString := client.connectionString()

	poolConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		fmt.Printf("couldn't parse postgres connection string\n")
		return nil, err
	}

	for connAttempts > 0 {
		client.Pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err == nil {
			err = client.Pool.Ping(context.TODO())
			if err == nil {
				break
			}
		}

		fmt.Printf("failed to connect to postgres\n")
		fmt.Printf("Postgres client is trying to connect, attempts left: %d\n", connAttempts)
		<-time.After(connTimeout)

		connAttempts--
	}

	if err != nil {
		fmt.Println("couldn't connect to postgres")
		return nil, err
	}

	return client, nil
}

func (c *Client) connectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.cfg.User,
		c.cfg.Password,
		c.cfg.Host,
		c.cfg.Port,
		c.cfg.Database)
}

func (c *Client) MigrateUp() error {
	m, err := migrate.New(c.cfg.MigrationsPath, c.connectionString())
	if err != nil {
		return fmt.Errorf("failed to create migration handler: %w", err)
	}

	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return ErrNoChange
		} else {
			return fmt.Errorf("failed to migrate up: %w", err)
		}
	}

	return nil
}

func (c *Client) MigrateDown() error {
	m, err := migrate.New(c.cfg.MigrationsPath, c.connectionString())
	if err != nil {
		return fmt.Errorf("failed to create migration handler: %w", err)
	}

	if err = m.Down(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return ErrNoChange
		} else {
			return fmt.Errorf("failed to migrate down: %w", err)
		}
	}

	return nil
}

func (c *Client) MigrateForce(version int) error {
	m, err := migrate.New(c.cfg.MigrationsPath, c.connectionString())
	if err != nil {
		return fmt.Errorf("failed to create migration handler: %w", err)
	}

	if err = m.Force(version); err != nil {
		return fmt.Errorf("failed to force migration: %w", err)
	}

	return nil
}

func (c *Client) Close() {
	if c.Pool != nil {
		c.Pool.Close()
	}
}
