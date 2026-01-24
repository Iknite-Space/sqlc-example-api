package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/Iknite-Space/sqlc-example-api/db/repo"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal("Error:", err)
	}
}

func run() error {
	ctx := context.Background()

	// Load environment
	if _, err := os.Stat(".env"); err == nil {
		err = godotenv.Load()
		if err != nil {
			return fmt.Errorf("failed to load env file: %w", err)
		}
	}

	// Connect to database
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return fmt.Errorf("DATABASE_URL environment variable is required")
	}

	db, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	querier := repo.New(db)

	// Seed dashboard data
	err = seedDashboardData(ctx, querier)
	if err != nil {
		return fmt.Errorf("failed to seed dashboard data: %w", err)
	}

	// Seed editor config
	err = seedEditorConfig(ctx, querier)
	if err != nil {
		return fmt.Errorf("failed to seed editor config: %w", err)
	}

	fmt.Println("Data seeding completed successfully!")
	return nil
}

func seedDashboardData(ctx context.Context, querier *repo.Queries) error {
	// Read dashboard.json
	data, err := os.ReadFile("mysampledata/dashboard.json")
	if err != nil {
		return fmt.Errorf("failed to read dashboard.json: %w", err)
	}

	var dashboard map[string]interface{}
	if err := json.Unmarshal(data, &dashboard); err != nil {
		return fmt.Errorf("failed to parse dashboard.json: %w", err)
	}

	// Seed docker containers
	if containers, ok := dashboard["dockerContainers"].([]interface{}); ok {
		for _, containerData := range containers {
			container := containerData.(map[string]interface{})

			// Prepare core data
			coreData := map[string]interface{}{
				"lastRun":      container["lastRun"],
				"origin":       container["origin"],
				"disk":         container["disk"],
				"ram":          container["ram"],
				"unused":       container["unused"],
				"highMem":      container["highMem"],
				"root":         container["root"],
				"exposedPorts": container["exposedPorts"],
				"unlimitedMem": container["unlimitedMem"],
			}

			coreDataJSON, err := json.Marshal(coreData)
			if err != nil {
				log.Printf("Failed to marshal core data for container %s: %v", container["name"], err)
				continue
			}
			customFieldsJSON, err := json.Marshal(map[string]interface{}{})
			if err != nil {
				log.Printf("Failed to marshal custom fields for container %s: %v", container["name"], err)
				continue
			}

			params := repo.CreateDockerContainerParams{
				Name:          container["name"].(string),
				Status:        container["status"].(string),
				CoreData:      coreDataJSON,
				CustomFields:  customFieldsJSON,
				SchemaVersion: int32Ptr(1),
			}

			_, err = querier.CreateDockerContainer(ctx, params)
			if err != nil {
				log.Printf("Failed to create container %s: %v", container["name"], err)
			} else {
				fmt.Printf("Created container: %s\n", container["name"])
			}
		}
	}

	// Seed git repositories
	if repos, ok := dashboard["gitRepos"].([]interface{}); ok {
		for _, repoData := range repos {
			repoInfo := repoData.(map[string]interface{})

			// Prepare core data
			coreData := map[string]interface{}{
				"untouched":        repoInfo["untouched"],
				"duplicate":        repoInfo["duplicate"],
				"clonedNeverBuilt": repoInfo["clonedNeverBuilt"],
			}

			coreDataJSON, err := json.Marshal(coreData)
			if err != nil {
				log.Printf("Failed to marshal core data for repo %s: %v", repoInfo["name"], err)
				continue
			}
			customFieldsJSON, err := json.Marshal(map[string]interface{}{})
			if err != nil {
				log.Printf("Failed to marshal custom fields for repo %s: %v", repoInfo["name"], err)
				continue
			}

			params := repo.CreateGitRepoParams{
				Name:          repoInfo["name"].(string),
				CoreData:      coreDataJSON,
				CustomFields:  customFieldsJSON,
				SchemaVersion: int32Ptr(1),
			}

			_, err = querier.CreateGitRepo(ctx, params)
			if err != nil {
				log.Printf("Failed to create repo %s: %v", repoInfo["name"], err)
			} else {
				fmt.Printf("Created repository: %s\n", repoInfo["name"])
			}
		}
	}

	// Seed cache data
	if cacheData, ok := dashboard["cacheData"].(map[string]interface{}); ok {
		for technology, techData := range cacheData {
			if techMap, ok := techData.(map[string]interface{}); ok {
				for cacheType, size := range techMap {
					coreData := map[string]interface{}{
						"size": size,
					}

					coreDataJSON, err := json.Marshal(coreData)
					if err != nil {
						log.Printf("Failed to marshal core data for cache %s/%s: %v", technology, cacheType, err)
						continue
					}
					customFieldsJSON, err := json.Marshal(map[string]interface{}{})
					if err != nil {
						log.Printf("Failed to marshal custom fields for cache %s/%s: %v", technology, cacheType, err)
						continue
					}

					params := repo.CreateCacheDataParams{
						Technology:    technology,
						CacheType:     cacheType,
						CoreData:      coreDataJSON,
						CustomFields:  customFieldsJSON,
						SchemaVersion: int32Ptr(1),
					}

					_, err = querier.CreateCacheData(ctx, params)
					if err != nil {
						log.Printf("Failed to create cache data %s/%s: %v", technology, cacheType, err)
					} else {
						fmt.Printf("Created cache data: %s/%s\n", technology, cacheType)
					}
				}
			}
		}
	}

	// Seed log entries
	if logEntries, ok := dashboard["logEntries"].([]interface{}); ok {
		for _, logData := range logEntries {
			if logStr, ok := logData.(string); ok {
				// Parse log entry format: "LEVEL: Message at TIME"
				parts := strings.SplitN(logStr, ":", 2)
				if len(parts) == 2 {
					level := strings.TrimSpace(parts[0])
					message := strings.TrimSpace(parts[1])

					coreData := map[string]interface{}{
						"originalEntry": logStr,
					}

					coreDataJSON, err := json.Marshal(coreData)
					if err != nil {
						log.Printf("Failed to marshal core data for log entry: %v", err)
						continue
					}
					customFieldsJSON, err := json.Marshal(map[string]interface{}{})
					if err != nil {
						log.Printf("Failed to marshal custom fields for log entry: %v", err)
						continue
					}

					now := time.Now()
					params := repo.CreateLogEntryParams{
						Level:         level,
						Message:       message,
						CoreData:      coreDataJSON,
						CustomFields:  customFieldsJSON,
						SchemaVersion: int32Ptr(1),
						Timestamp:     pgtype.Timestamp{Time: now, Valid: true},
					}

					_, err = querier.CreateLogEntry(ctx, params)
					if err != nil {
						log.Printf("Failed to create log entry: %v", err)
					} else {
						fmt.Printf("Created log entry: %s\n", level)
					}
				}
			}
		}
	}

	// Seed secrets
	if secrets, ok := dashboard["secrets"].([]interface{}); ok {
		for _, secretData := range secrets {
			if secretStr, ok := secretData.(string); ok {
				coreData := map[string]interface{}{
					"location": secretStr,
				}

				coreDataJSON, err := json.Marshal(coreData)
				if err != nil {
					log.Printf("Failed to marshal core data for secret: %v", err)
					continue
				}
				customFieldsJSON, err := json.Marshal(map[string]interface{}{})
				if err != nil {
					log.Printf("Failed to marshal custom fields for secret: %v", err)
					continue
				}

				params := repo.CreateSecretParams{
					Description:   secretStr,
					CoreData:      coreDataJSON,
					CustomFields:  customFieldsJSON,
					SchemaVersion: int32Ptr(1),
				}

				_, err = querier.CreateSecret(ctx, params)
				if err != nil {
					log.Printf("Failed to create secret: %v", err)
				} else {
					fmt.Printf("Created secret: %s\n", secretStr)
				}
			}
		}
	}

	// Seed registry data
	if registryData, ok := dashboard["registryData"].([]interface{}); ok {
		for _, regData := range registryData {
			if regMap, ok := regData.(map[string]interface{}); ok {
				coreData := map[string]interface{}{
					"data": regMap["data"],
					"type": regMap["type"],
				}

				coreDataJSON, err := json.Marshal(coreData)
				if err != nil {
					log.Printf("Failed to marshal core data for registry data: %v", err)
					continue
				}
				customFieldsJSON, err := json.Marshal(map[string]interface{}{})
				if err != nil {
					log.Printf("Failed to marshal custom fields for registry data: %v", err)
					continue
				}

				params := repo.CreateRegistryDataParams{
					Subkey:        regMap["subkey"].(string),
					ValueName:     regMap["value"].(string),
					CoreData:      coreDataJSON,
					CustomFields:  customFieldsJSON,
					SchemaVersion: int32Ptr(1),
				}

				_, err = querier.CreateRegistryData(ctx, params)
				if err != nil {
					log.Printf("Failed to create registry data: %v", err)
				} else {
					fmt.Printf("Created registry data: %s\n", regMap["subkey"])
				}
			}
		}
	}

	// Seed plist data
	if plistData, ok := dashboard["plistData"].([]interface{}); ok {
		for _, plistItem := range plistData {
			if plistMap, ok := plistItem.(map[string]interface{}); ok {
				coreData := map[string]interface{}{
					"value": plistMap["value"],
					"type":  plistMap["type"],
				}

				coreDataJSON, err := json.Marshal(coreData)
				if err != nil {
					log.Printf("Failed to marshal core data for plist data: %v", err)
					continue
				}
				customFieldsJSON, err := json.Marshal(map[string]interface{}{})
				if err != nil {
					log.Printf("Failed to marshal custom fields for plist data: %v", err)
					continue
				}

				params := repo.CreatePlistDataParams{
					Key:           plistMap["key"].(string),
					CoreData:      coreDataJSON,
					CustomFields:  customFieldsJSON,
					SchemaVersion: int32Ptr(1),
				}

				_, err = querier.CreatePlistData(ctx, params)
				if err != nil {
					log.Printf("Failed to create plist data: %v", err)
				} else {
					fmt.Printf("Created plist data: %s\n", plistMap["key"])
				}
			}
		}
	}

	return nil
}

func seedEditorConfig(ctx context.Context, querier *repo.Queries) error {
	// Read editorConfig.json
	data, err := os.ReadFile("mysampledata/editorConfig.json")
	if err != nil {
		return fmt.Errorf("failed to read editorConfig.json: %w", err)
	}

	params := repo.CreateEditorConfigParams{
		Name:       "Default Editor Config",
		ConfigData: data,
	}

	_, err = querier.CreateEditorConfig(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to create editor config: %w", err)
	}

	fmt.Println("Created editor configuration")
	return nil
}

// Helper functions
func int32Ptr(i int32) *int32 {
	return &i
}
