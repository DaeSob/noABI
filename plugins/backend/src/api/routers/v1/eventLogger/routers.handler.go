package eventLogger

import (
	"cia/api/logHandler"
	"cia/api/preference"
	eventCollector "cia/pkg/eventLogger/eventCollector"
	jMapper "cia/common/utils/json/mapper"
	"cia/common/utils"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

const (
	FieldDB     = "db"
	FieldActors = "actors"
	ModeDev     = "dev"
	ModeLive    = "live"
)

// validateDevModeConstraints validates that db and actors are not set in dev mode
func validateDevModeConstraints(_data map[string]interface{}) error {
	if _, exists := _data[FieldDB]; exists {
		return fmt.Errorf("db setting is not allowed in dev mode")
	}
	if _, exists := _data[FieldActors]; exists {
		return fmt.Errorf("actors setting is not allowed in dev mode")
	}
	return nil
}

// validateModeForModification validates if chain can be modified based on mode
func validateModeForModification(_mode string) error {
	if _mode == ModeLive {
		return fmt.Errorf("cannot modify live mode chain")
	}
	if _mode != ModeDev {
		return fmt.Errorf("invalid mode")
	}
	return nil
}

// updateEventMetaFiles updates event.meta files for collectors when startBlockNumber changes
func updateEventMetaFiles(_alias string, _collectors interface{}) error {
	collectors, ok := _collectors.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid collectors format")
	}

	eventLoggerInfo := preference.GetEventLogger(_alias)
	if eventLoggerInfo == nil {
		return fmt.Errorf("chain not found: %s", _alias)
	}

	for collectorName, v := range collectors {
		var collector preference.TCollector
		if err := mapstructure.Decode(v, &collector); err != nil {
			logHandler.Write("trace", 0, "failed to decode collector", collectorName, "error:", err.Error())
			continue
		}

		// Update event.meta file with startBlockNumber
		// If file doesn't exist, it will be created
		eventCollector.UpdateEventMetaInfo(_alias, collectorName, collector.StartBlockNumber)
		logHandler.Write("trace", 0, "updated event.meta for chain", _alias, "collector", collectorName, "startBlockNumber", collector.StartBlockNumber)
	}

	return nil
}

// GetChains returns all chain information
func GetChains(c *gin.Context) {
	chains := preference.GetEventLoggerChains()
	result := make([]map[string]interface{}, 0)

	for _, alias := range chains {
		eventLogger := preference.GetEventLogger(alias)
		if eventLogger != nil {
			chainMap := convertEventLoggerToMap(eventLogger)
			chainMap["alias"] = alias // Add alias to response
			result = append(result, chainMap)
		}
	}

	c.JSON(http.StatusOK, gin.H{"chains": result})
}

// GetChain returns specific chain information by alias
func GetChain(c *gin.Context) {
	alias := c.Param("alias")

	eventLogger := preference.GetEventLogger(alias)
	if eventLogger == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "chain not found"})
		return
	}

	chainMap := convertEventLoggerToMap(eventLogger)
	chainMap["alias"] = alias // Add alias to response
	c.JSON(http.StatusOK, chainMap)
}

// UpdateOrAddChain updates an existing chain or adds a new chain by alias
func UpdateOrAddChain(c *gin.Context) {
	alias := c.Param("alias")

	// Read request body
	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
		return
	}

	var updateData map[string]interface{}
	if err := jMapper.FromJson(bytes, &updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json format"})
		return
	}

	// Check if chain exists
	eventLogger := preference.GetEventLogger(alias)
	isNewChain := (eventLogger == nil)

	if isNewChain {
		// Add new chain
		// 1. Dev mode constraint validation
		if err := validateDevModeConstraints(updateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 2. Set chainId if not provided (chainId is the actual blockchain chain ID, not alias)
		if updateData["chainId"] == nil {
			// If chainId is not provided, we cannot determine it from alias
			// User must provide chainId in the request body
			c.JSON(http.StatusBadRequest, gin.H{"error": "chainId is required in request body"})
			return
		}

		// 3. Mode must be dev for new chains
		if mode, exists := updateData["mode"]; exists {
			if mode != ModeDev {
				c.JSON(http.StatusBadRequest, gin.H{"error": "new chain must be in dev mode"})
				return
			}
		} else {
			updateData["mode"] = ModeDev // default for new chains
		}

		// 4. Add to memory (use alias as key)
		if err := preference.AddEventLoggerChain(alias, updateData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 5. Initialize and run eventCollector for the new chain
		if err := eventCollector.InitializeChain(alias); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to initialize collector: " + err.Error()})
			return
		}
		if err := eventCollector.RunChain(alias); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to run collector: " + err.Error()})
			return
		}

		// 6. Update event.meta files for collectors
		if updateData["collectors"] != nil {
			if err := updateEventMetaFiles(alias, updateData["collectors"]); err != nil {
				// Log error but don't fail the request
				logHandler.Write("trace", 0, "failed to update event.meta files:", err.Error())
			}
		}

		c.JSON(http.StatusCreated, gin.H{"message": "chain created", "alias": alias})
	} else {
		// Update existing chain
		// 1. Mode validation
		if err := validateModeForModification(eventLogger.Mode); err != nil {
			if eventLogger.Mode == ModeLive {
				c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			return
		}

		// 2. Dev mode constraint validation
		if err := validateDevModeConstraints(updateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 3. Don't allow mode change via update
		if updateData["mode"] != nil && updateData["mode"] != ModeDev {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot change mode from dev"})
			return
		}

		// 4. Check if path or collectors changed (before update)
		oldPath := eventLogger.Path
		oldCollectors := make(map[string]bool)
		for name := range eventLogger.Collectors {
			oldCollectors[name] = true
		}

		// 5. Update in memory (use alias as key)
		if err := preference.UpdateEventLoggerChain(alias, updateData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 6. Get updated eventLogger to check changes
		updatedEventLogger := preference.GetEventLogger(alias)
		if updatedEventLogger == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get updated chain"})
			return
		}

		// 7. Handle path change - reinitialize directories
		if updateData["path"] != nil && updatedEventLogger.Path != oldPath {
			if err := eventCollector.InitializeChain(alias); err != nil {
				logHandler.Write("trace", 0, "failed to reinitialize chain directories after path change:", err.Error())
				// Continue even if directory creation fails
			}
		}

		// 8. Handle collectors change - create directories for new collectors
		if updateData["collectors"] != nil {
			// Create directories for new collectors
			for name := range updatedEventLogger.Collectors {
				if !oldCollectors[name] {
					// New collector added - create directory
					logPath := preference.GetEventLoggerFilePath(alias)
					collectorPath := fmt.Sprintf("%v/%v/%v", logPath, alias, name)
					if err := utils.Mkdir(collectorPath); err != nil {
						logHandler.Write("trace", 0, "failed to create collector directory", collectorPath, "error:", err.Error())
					}
				}
			}

			// Update event.meta files for collectors
			if err := updateEventMetaFiles(alias, updateData["collectors"]); err != nil {
				// Log error but don't fail the request
				logHandler.Write("trace", 0, "failed to update event.meta files:", err.Error())
			}
		}

		c.JSON(http.StatusOK, gin.H{"message": "chain updated", "alias": alias})
	}
}

// RemoveChain deletes a chain by alias
func RemoveChain(c *gin.Context) {
	alias := c.Param("alias")

	// 1. Check if chain exists
	eventLogger := preference.GetEventLogger(alias)
	if eventLogger == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "chain not found"})
		return
	}

	// 2. Mode validation
	if err := validateModeForModification(eventLogger.Mode); err != nil {
		if eventLogger.Mode == ModeLive {
			c.JSON(http.StatusForbidden, gin.H{"error": "cannot delete live mode chain"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	// 3. Stop all workers for this chain
	if err := eventCollector.StopChain(alias); err != nil {
		logHandler.Write("trace", 0, "failed to stop chain workers:", err.Error())
		// Continue with deletion even if stop fails
	}

	// 4. Delete from memory
	if err := preference.DeleteEventLoggerChain(alias); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "chain deleted", "alias": alias})
}

// convertEventLoggerToMap converts TEvetLoggerInfo to map for JSON response
func convertEventLoggerToMap(eventLogger *preference.TEvetLoggerInfo) map[string]interface{} {
	result := make(map[string]interface{})
	result["mode"] = eventLogger.Mode
	result["chainId"] = eventLogger.ChainId
	result["rpc"] = eventLogger.RPC
	result["db"] = eventLogger.DB
	result["path"] = eventLogger.Path
	result["commitBlockCount"] = eventLogger.CommitBlockCount
	result["period"] = eventLogger.Period
	result["logRange"] = eventLogger.LogRange
	result["inputLogging"] = eventLogger.InputLogging
	result["blockLogging"] = eventLogger.BlockLogging
	result["collections"] = eventLogger.Collections
	result["collectors"] = eventLogger.Collectors
	result["actors"] = eventLogger.Actors
	return result
}

