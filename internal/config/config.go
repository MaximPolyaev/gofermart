package config

import (
	"flag"
	"os"
	"regexp"
	"sort"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

const (
	defaultRunAddress           = "localhost:8081"
	defaultDatabaseURI          = "host=localhost port=54333 user=admin password=password dbname=gofermart sslmode=disable"
	defaultAccrualSystemAddress = "localhost:8082"
)

type Config struct {
	RunAddress           *string `env:"RUN_ADDRESS"`
	DatabaseURI          *string `env:"DATABASE_URI"`
	AccrualSystemAddress *string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func New() *Config {
	return &Config{}
}

func (c *Config) Parse() error {
	if err := c.loadEnvFiles(); err != nil {
		return err
	}

	if err := env.Parse(c); err != nil {
		return err
	}

	if err := c.parseFlags(); err != nil {
		return err
	}

	return nil
}

func (c *Config) loadEnvFiles() error {
	files, err := c.getFileNames()
	if err != nil {
		return err
	}

	if len(files) != 0 {
		if err := godotenv.Load(files...); err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) getFileNames() ([]string, error) {
	files, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}

	var fileNames []string

	envRegexp := regexp.MustCompile(`^\.env((.+)?)`)

	for _, f := range files {
		fileName := f.Name()

		if !f.IsDir() && envRegexp.MatchString(fileName) {
			fileNames = append(fileNames, fileName)
		}
	}

	sort.Sort(sort.Reverse(sort.StringSlice(fileNames)))

	return fileNames, nil
}

func (c *Config) parseFlags() error {
	runAddr := flag.String("a", defaultRunAddress, "run address")
	accrualAddr := flag.String("r", defaultAccrualSystemAddress, "accrual system address")
	databaseURI := flag.String("d", defaultDatabaseURI, "database uri")

	if c.RunAddress == nil {
		c.RunAddress = runAddr
	}

	if c.DatabaseURI == nil {
		c.DatabaseURI = databaseURI
	}

	if c.AccrualSystemAddress == nil {
		c.AccrualSystemAddress = accrualAddr
	}

	flag.Parse()

	return nil
}
