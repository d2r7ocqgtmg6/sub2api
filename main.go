package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	defaultPort    = 8080
	defaultHost    = "0.0.0.0"
	appVersion     = "1.0.0"
)

// Config holds the application configuration
type Config struct {
	Host    string
	Port    int
	Debug   bool
	Token   string
}

func main() {
	cfg := parseConfig()

	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	router := setupRouter(cfg)

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	log.Printf("sub2api v%s starting on %s", appVersion, addr)

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// parseConfig reads configuration from flags and environment variables
func parseConfig() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.Host, "host", getEnvString("HOST", defaultHost), "Host to listen on")
	flag.IntVar(&cfg.Port, "port", getEnvInt("PORT", defaultPort), "Port to listen on")
	flag.BoolVar(&cfg.Debug, "debug", getEnvBool("DEBUG", false), "Enable debug mode")
	flag.StringVar(&cfg.Token, "token", os.Getenv("API_TOKEN"), "API token for authentication")
	flag.Parse()

	return cfg
}

// setupRouter initializes the Gin router with all routes
func setupRouter(cfg *Config) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"version": appVersion,
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		v1.GET("/sub", handleSubscription(cfg))
	}

	return router
}

// handleSubscription returns a handler for subscription conversion requests
func handleSubscription(cfg *Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		subURL := c.Query("url")
		if subURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "url parameter is required"})
			return
		}

		// Token validation if configured
		if cfg.Token != "" {
			token := c.Query("token")
			if token != cfg.Token {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"url": subURL, "status": "received"})
	}
}

// getEnvString returns an environment variable value or a default
func getEnvString(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// getEnvInt returns an environment variable as int or a default
func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

// getEnvBool returns an environment variable as bool or a default
func getEnvBool(key string, defaultVal bool) bool {
	if val := os.Getenv(key); val != "" {
		if b, err := strconv.ParseBool(val); err == nil {
			return b
		}
	}
	return defaultVal
}
