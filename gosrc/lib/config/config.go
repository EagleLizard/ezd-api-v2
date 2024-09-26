package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/EagleLizard/jcd-api/gosrc/util/constants"
	"github.com/joho/godotenv"
)

type JcdApiConfigType struct {
	Host    string
	Port    string
	SfsHost string
	SfsPort string
	JcdEnv  string
}

var JcdApiConfig *JcdApiConfigType

var requiredEnvVars = [...]string{
	"JCD_SESSION_SECRET",
	"JCD_ENCRYPTION_SECRET",
	"JCD_JWT_SECRET",
	"JCD_ENV",
}

func init() {
	baseDir := constants.BaseDir()
	dotenvFilePath := filepath.Join(baseDir, ".env")
	err := godotenv.Load(dotenvFilePath)
	if err != nil {
		log.Fatal(err)
	}
	err = validateRequiredVars()
	if err != nil {
		log.Fatal(err)
	}

	cfg := JcdApiConfigType{
		Host:    getEnvVarOrDefault("JCD_API_HOST", "0.0.0.0"),
		Port:    getEnvVarOrDefault("JCD_API_PORT", "4040"),
		SfsHost: getEnvVarOrDefault("SFS_HOST", "0.0.0.0"),
		SfsPort: getEnvVarOrDefault("SFS_PORT", "4041"),
		JcdEnv:  getEnvVar("JCD_ENV"),
	}
	JcdApiConfig = &cfg
}

func getEnvVar(envKey string) string {
	envVal := os.Getenv(envKey)
	/* maybe panic here. Does no default imply required? */
	return envVal
}

func getEnvVarOrDefault(envKey string, defaultVal string) string {
	envVal := os.Getenv(envKey)
	if len(envVal) == 0 {
		return defaultVal
	}
	return envVal
}

func validateRequiredVars() error {
	missingEnvKeys := []string{}
	for _, envKey := range requiredEnvVars {
		envVal := os.Getenv(envKey)
		if len(envVal) == 0 {
			missingEnvKeys = append(missingEnvKeys, envKey)
			// fmt.Fprintf(os.Stderr, "missing required env var: %s\n", envKey)
		}
	}
	if len(missingEnvKeys) > 0 {
		return fmt.Errorf("missing env vars: %s", strings.Join(missingEnvKeys, ", "))
	}
	return nil
}
