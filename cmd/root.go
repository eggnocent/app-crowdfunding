package cmd

import (
	v1 "app-crowdfunding/delivery/v1"
	"app-crowdfunding/router"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	dbPool  *sqlx.DB
	logger  *logrus.Logger
)

var rootCmd = &cobra.Command{
	Use:   "app-crowdfunding",
	Short: "Application for Crowdfunding",
	Run: func(cmd *cobra.Command, args []string) {
		// Inisialisasi koneksi database
		if err := initDatabase(); err != nil {
			log.Fatalf("failed to connect database: %v", err)
		}
		defer dbPool.Close()

		// Jalankan migrasi
		if err := runMigrations(); err != nil {
			logger.Fatalf("Migration failed: %v", err)
		}

		fmt.Println("Successfully connected to Database and migrations applied!")

		// Inisialisasi server
		r := mux.NewRouter()

		router.Init(dbPool)

		// public
		apiV1 := r.PathPrefix("/api/v1").Subrouter()
		v1.NewLogin(apiV1)
		v1.NewRegistration(apiV1)

		// with authentication
		apiV1Session := r.PathPrefix("/api/v1").Subrouter()
		apiV1Session.Use(router.GetAuthMiddleware())

		v1.NewAPIUser(apiV1Session)
		v1.NewCampaign(apiV1Session)

		// Jalankan server
		port := viper.GetString("server.port")
		fmt.Printf("Starting server on port %s...\n", port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(
		initConfig,
		initLogger,
	)

	// Definisikan flag untuk file konfigurasi
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is file/.config.toml)")
}

// initConfig memuat konfigurasi toml
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("file")
		viper.SetConfigType("toml")
		viper.SetConfigName(".config")
	}

	viper.AutomaticEnv()

	// Membaca file konfigurasi jika ditemukan, jika tidak, fatalf akan memberhentikan app
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file %s", err)
	}
}

// initLogger mengatur logger dengan Logrus
func initLogger() {
	logger = logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)
}

// Mengonfigurasi dan menghubungkan ke database
func initDatabase() error {
	dbURL := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		viper.GetString("database.host"),
		viper.GetInt("database.port"),
		viper.GetString("database.username"),
		viper.GetString("database.password"),
		viper.GetString("database.name"),
		viper.GetString("database.ssl_mode"),
	)

	var err error
	dbPool, err = sqlx.Connect("postgres", dbURL)
	if err != nil {
		return err
	}
	return nil
}

func runMigrations() error {
	migrations := &migrate.FileMigrationSource{
		Dir: "file/migrations",
	}

	// Menjalankan migrasi ke atas
	n, err := migrate.Exec(dbPool.DB, "postgres", migrations, migrate.Up)
	if err != nil {
		return fmt.Errorf("error applying migrations: %w", err)
	}

	logger.Infof("Successfully applied %d migrations!", n)
	return nil
}

func GetDBPool() *sqlx.DB {
	return dbPool
}
