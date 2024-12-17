package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	cfg "github.com/danushk97/image-analyzer/internal/config"
	_ "github.com/danushk97/image-analyzer/internal/database/migrations"
	"github.com/danushk97/image-analyzer/pkg/env"
	storage "github.com/danushk97/image-analyzer/pkg/storage/sql"
	"github.com/pressly/goose/v3"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
	dir   = flags.String("dir", "./internal/database/migrations",
		"Directory with migration files")
	verbose = flags.Bool("v", false, "Enable verbose mode")
)

func main() {
	env := env.GetEnv()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// load configurations and distribute parts of it in main
	config := cfg.NewConfig(env)

	err := goose.SetDialect(config.Store.SQL.Dialect)
	if err != nil {
		log.Fatalf("could not set dialect, err:%+v", err)
	}

	// storage service is the service for main persistent store
	db, err := storage.NewDb(config.Store.SQL)

	if err != nil {
		log.Fatalf("could not create database, err:%+v", err)
	}

	dbInstance, err := db.GetInstance(ctx).DB()

	if err != nil {
		log.Fatalf("could not get the datinstance,  err:%+v", err)
	}

	run(dbInstance, *dir)
}

func run(db *sql.DB, dir string) {
	flags.Usage = usage
	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Fatalf("error parsing flags: %v", err)
	}
	args := flags.Args()
	if *verbose {
		goose.SetVerbose(true)
	}

	// I.e. no command provided, hence print usage and return.
	if len(args) < 1 {
		cmd := os.Getenv("MIGRATION_CMD")
		if cmd == "" {
			flags.Usage()
			return
		}

		args = append(args, cmd)

	}
	// Prepares command and arguments for goose's run.
	command := args[0]
	arguments := []string{}
	if len(args) > 1 {
		arguments = append(arguments, args[1:]...)
	}

	// Finally, executes the goose's command.
	if err := goose.Run(command, db, dir, arguments...); err != nil {
		log.Fatalf("error parsing flags: %v", err)
	}

}

func usage() {
	flags.PrintDefaults()
	fmt.Println(usageCommands)
}

var usageCommands = `
	Commands:
		up                   Migrate the DB to the most recent version available
		up-to VERSION        Migrate the DB to a specific VERSION
		down                 Roll back the version by 1
		down-to VERSION      Roll back to a specific VERSION
		redo                 Re-run the latest migration
		reset                Roll back all migrations
		status               Dump the migration status for the current DB
		version              Print the current version of the database
		create NAME          Creates new migration file with the current timestamp
		fix                  Apply sequential ordering to migrations
`
