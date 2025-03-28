package madden

import (
	"net/http"
	"strings"

	"github.comm/kevinlucasklein/madden-discord-bot/pkg/utils"
)

// Service handles Madden Companion App exports
type Service struct {
	DataDir string
	logger  *utils.Logger
}

// NewService creates a new Madden service instance
func NewService(dataDir string) *Service {
	return &Service{
		DataDir: dataDir,
		logger:  &utils.Logger{}, // This will be replaced with a real logger
	}
}

// SetLogger sets the logger for the service
func (s *Service) SetLogger(logger *utils.Logger) {
	s.logger = logger
}

// RegisterRoutes sets up HTTP routes for the Madden service
func (s *Service) RegisterRoutes(mux *http.ServeMux, exportPath string) {
	// Apply CORS middleware to our handlers

	// Handle the base export path
	mux.HandleFunc(exportPath, utils.AllowCORS(s.ExportHandler))

	// Handle all nested paths under export as well (Madden Companion App uses nested paths)
	// This wildcard handler will catch paths like /export/ps5/123456/week/reg/1/schedules
	mux.HandleFunc("/", utils.AllowCORS(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request path starts with the export path
		if strings.HasPrefix(r.URL.Path, exportPath+"/") {
			s.logger.Info("Handling nested export path: %s", r.URL.Path)
			s.ExportHandler(w, r)
			return
		}

		// For the root path and other non-export paths, use the status handler
		if r.URL.Path == "/" {
			s.StatusHandler(w, r)
			return
		}

		// For unhandled paths, return 404
		http.NotFound(w, r)
	}))
}
