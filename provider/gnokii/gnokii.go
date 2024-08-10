package gnokii

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/messagebird/sachet"
)

// Config configuration struct for Gnokii Client.
type Config struct {
	DbUsername string `yaml:"db_username"`
	DbPassword string `yaml:"db_password"`
	DbHost     string `yaml:"db_host"`
	DbPort     int    `yaml:"db_port"`
	DbName     string `yaml:"db_name"`
}

var _ (sachet.Provider) = (*Gnokii)(nil)

type Gnokii struct {
	Config
	Db *sql.DB
}

// NewGnokii creates a new Gnokii SMS client.
func NewGnokii(config Config) *Gnokii {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.DbUsername, config.DbPassword, config.DbHost, config.DbPort, config.DbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error while connecting to gnokii DB: %s", err)
	}

	gnokii := &Gnokii{config, db}
	return gnokii
}

// Send send sms to n number of people using bulk sms api.
func (c *Gnokii) Send(message sachet.Message) error {
	for _, recipient := range message.To {
		query := "INSERT INTO `outbox` (`number`, `text`) VALUES (?, ?)"
		_, err := c.Db.ExecContext(context.Background(), query, recipient, message.Text)
		if err != nil {
			return fmt.Errorf("Error while inserting sms into DB: %w", err)
		}
	}
	return nil
}
