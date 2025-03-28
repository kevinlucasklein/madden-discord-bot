package madden

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.comm/kevinlucasklein/madden-discord-bot/pkg/utils"
)

// ExportHandler handles Madden Companion App export requests
func (s *Service) ExportHandler(w http.ResponseWriter, r *http.Request) {
	// Log detailed information about the request
	s.logger.Info("Received request: Method=%s, URL=%s, RemoteAddr=%s, Content-Type=%s",
		r.Method, r.URL.Path, r.RemoteAddr, r.Header.Get("Content-Type"))

	// Log query parameters and headers for debugging
	s.logger.Debug("Request Query Params: %v", r.URL.Query())
	for name, values := range r.Header {
		s.logger.Debug("Header %s: %s", name, strings.Join(values, ", "))
	}

	// Always return a 200 OK status first thing to match Madden's expectations
	w.WriteHeader(http.StatusOK)

	// If it's just a GET request with no body, return a helpful message
	if r.Method == http.MethodGet {
		s.logger.Info("GET request received, sending status message")
		fmt.Fprintf(w, "Madden Companion Export endpoint is ready. Send data to this URL.")
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.Error("Error reading request body: %v", err)
		fmt.Fprintf(w, "Received request but could not read body: %v", err)
		return
	}
	defer r.Body.Close()

	s.logger.Debug("Received data of size %d bytes", len(body))

	// If body is empty, return a simple success
	if len(body) == 0 {
		s.logger.Warn("Empty request body received")
		fmt.Fprintf(w, "Request received with empty body. Endpoint is working.")
		return
	}

	// Extract metadata from the URL path
	pathMetadata := extractPathMetadata(r.URL.Path)
	s.logger.Debug("Extracted path metadata: %v", pathMetadata)

	// Process the export data
	result, err := s.ProcessExport(body, pathMetadata)
	if err != nil {
		s.logger.Error("Error processing export: %v", err)
		fmt.Fprintf(w, "Data received but could not be processed: %v", err)
		return
	}

	// Return success response
	s.logger.Info("Successfully processed export data to %s", result)
	fmt.Fprintf(w, "Data received and saved successfully")
}

// StatusHandler provides a simple status page for the service
func (s *Service) StatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	s.logger.Debug("Status page requested from %s", r.RemoteAddr)
	fmt.Fprintf(w, "Madden Companion Export Service is running\n")
	fmt.Fprintf(w, "Send your Madden Companion App exports to this server's export endpoint\n")
	fmt.Fprintf(w, "Example URL: http://your-server-ip:8080/export\n")
}

// PathMetadata contains metadata extracted from the URL path
type PathMetadata struct {
	Platform   string
	LeagueID   string
	ExportType string
	SeasonType string
	WeekNumber string
	DataType   string
}

// extractPathMetadata extracts metadata from the URL path
// Expected formats:
// - /export/platform/leagueId/week/seasonType/weekNumber/dataType (for weekly data)
// - /export/platform/leagueId/dataType (for league data like leagueteams, standings)
func extractPathMetadata(path string) PathMetadata {
	metadata := PathMetadata{}

	// Split the path into components
	parts := strings.Split(path, "/")

	// Remove empty parts
	var cleanParts []string
	for _, part := range parts {
		if part != "" {
			cleanParts = append(cleanParts, part)
		}
	}

	// Parse the components based on position
	if len(cleanParts) > 0 {
		// First part should be "export"
		if cleanParts[0] == "export" && len(cleanParts) > 1 {
			// Extract platform (ps5, xbox, etc.)
			if len(cleanParts) > 1 {
				metadata.Platform = cleanParts[1]
			}

			// Extract league ID
			if len(cleanParts) > 2 {
				metadata.LeagueID = cleanParts[2]
			}

			// Check if it's a week-based export
			if len(cleanParts) > 3 && cleanParts[3] == "week" {
				metadata.ExportType = "week"

				// Extract season type (reg, pre, post)
				if len(cleanParts) > 4 {
					metadata.SeasonType = cleanParts[4]
				}

				// Extract week number
				if len(cleanParts) > 5 {
					metadata.WeekNumber = cleanParts[5]
				}

				// Extract data type (schedules, team, etc.)
				if len(cleanParts) > 6 {
					metadata.DataType = cleanParts[6]
				}
			} else if len(cleanParts) > 3 {
				// It's another type of export like /export/ps5/12345/leagueteams
				metadata.ExportType = cleanParts[3]
			}
		}
	}

	return metadata
}

// ProcessExport handles the actual processing of the export data
func (s *Service) ProcessExport(data []byte, metadata PathMetadata) (string, error) {
	var jsonData interface{}

	// Try to parse as JSON
	if err := json.Unmarshal(data, &jsonData); err != nil {
		// If not valid JSON, store as raw text
		s.logger.Warn("Received non-JSON data, saving as raw text: %v", err)

		// Log a preview of the data
		preview := string(data)
		if len(preview) > 100 {
			preview = preview[:100] + "..."
		}
		s.logger.Debug("Data preview: %s", preview)

		// Ensure data directory exists
		if err := utils.EnsureDirectoryExists(s.DataDir); err != nil {
			return "", fmt.Errorf("failed to create data directory: %w", err)
		}

		// Create a timestamped filename for raw data
		timestamp := time.Now().Format("20060102-150405")
		filename := filepath.Join(s.DataDir, fmt.Sprintf("madden_raw_%s.txt", timestamp))

		// Save as raw text
		if err := utils.SaveRawToFile(filename, data); err != nil {
			return "", fmt.Errorf("failed to save raw data: %w", err)
		}

		s.logger.Info("Saved raw data to %s", filename)
		return filename, nil
	}

	// Handle various JSON formats
	exportType := "unknown"

	// First try to get the export type from the JSON data
	if jsonMap, ok := jsonData.(map[string]interface{}); ok {
		if typeVal, hasType := jsonMap["exportType"]; hasType && typeVal != nil {
			exportType = fmt.Sprintf("%v", typeVal)
		} else if typeVal, hasType := jsonMap["type"]; hasType && typeVal != nil {
			exportType = fmt.Sprintf("%v", typeVal)
		} else {
			// Use metadata from the URL if available
			if metadata.DataType != "" {
				// DataType is already set (like "schedules", "passing", etc.)
				exportType = metadata.DataType
			} else if metadata.ExportType != "" && metadata.ExportType != "week" {
				// Use the export type from the path (like "leagueteams", "standings", etc.)
				exportType = metadata.ExportType
			} else {
				// Get the last part of the URL path as a fallback
				parts := strings.Split(metadata.Platform, "/")
				if len(parts) > 0 && parts[len(parts)-1] != "" {
					exportType = parts[len(parts)-1]
				} else {
					exportType = "object"
				}
			}
		}

		// Log some object keys for debugging
		keys := make([]string, 0, len(jsonMap))
		for k := range jsonMap {
			keys = append(keys, k)
		}
		s.logger.Debug("JSON object keys: %v", keys)
	} else if jsonArr, ok := jsonData.([]interface{}); ok {
		// It's an array
		if metadata.DataType != "" {
			exportType = metadata.DataType
		} else {
			exportType = fmt.Sprintf("array_%d", len(jsonArr))
		}
	} else {
		// Some other JSON value
		if metadata.DataType != "" {
			exportType = metadata.DataType
		} else {
			exportType = "unknown"
		}
	}

	s.logger.Debug("Processing export of type: %s", exportType)

	// Ensure data directory exists
	if err := utils.EnsureDirectoryExists(s.DataDir); err != nil {
		return "", fmt.Errorf("failed to create data directory: %w", err)
	}

	// Build a more descriptive filename using metadata
	var filenameParts []string

	// Add metadata components if available
	if metadata.Platform != "" {
		filenameParts = append(filenameParts, metadata.Platform)
	}

	if metadata.LeagueID != "" {
		filenameParts = append(filenameParts, "league_"+metadata.LeagueID)
	}

	if metadata.SeasonType != "" && metadata.WeekNumber != "" {
		filenameParts = append(filenameParts, metadata.SeasonType+"_week_"+metadata.WeekNumber)
	}

	if exportType != "" {
		filenameParts = append(filenameParts, exportType)
	}

	// Add timestamp
	timestamp := time.Now().Format("20060102-150405")

	// Create the filename
	var filenameBase string
	if len(filenameParts) > 0 {
		filenameBase = strings.Join(filenameParts, "_")
	} else {
		filenameBase = "madden_export_" + exportType
	}

	filename := filepath.Join(s.DataDir, fmt.Sprintf("%s_%s.json", filenameBase, timestamp))

	// Save the data with pretty formatting
	if err := utils.SaveJSONToFile(filename, jsonData); err != nil {
		return "", fmt.Errorf("failed to save data: %w", err)
	}

	s.logger.Info("Saved %s export data to %s", exportType, filename)
	return filename, nil
}
