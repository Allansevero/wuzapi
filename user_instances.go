package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

// List instances for authenticated system user
func (s *server) ListMyInstances() http.HandlerFunc {
	type instanceStruct struct {
		Id              string         `db:"id" json:"id"`
		Name            string         `db:"name" json:"name"`
		Token           string         `db:"token" json:"token"`
		Jid             string         `db:"jid" json:"jid"`
		Connected       sql.NullBool   `db:"connected" json:"connected"`
		DestinationNumber string       `db:"destination_number" json:"destination_number"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		systemUserID, ok := r.Context().Value("system_user_id").(int)
		if !ok {
			s.respondWithJSON(w, http.StatusUnauthorized, map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"error":   "unauthorized",
				"success": false,
			})
			return
		}

		var instances []instanceStruct
		query := `SELECT id, name, token, jid, connected, destination_number 
				  FROM users 
				  WHERE system_user_id = $1 
				  ORDER BY name`

		err := s.db.Select(&instances, query, systemUserID)
		if err != nil {
			log.Error().Err(err).Int("system_user_id", systemUserID).Msg("Failed to get instances")
			s.respondWithJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"code":    http.StatusInternalServerError,
				"error":   "failed to get instances",
				"success": false,
			})
			return
		}

		// Add real-time connection status
		for i := range instances {
			if clientManager.GetWhatsmeowClient(instances[i].Id) != nil {
				client := clientManager.GetWhatsmeowClient(instances[i].Id)
				// Only mark as connected if both connected AND logged in
				isConnected := client.IsConnected() && client.IsLoggedIn()
				instances[i].Connected = sql.NullBool{Bool: isConnected, Valid: true}
			}
		}

		s.respondWithJSON(w, http.StatusOK, map[string]interface{}{
			"code":    http.StatusOK,
			"data":    instances,
			"success": true,
		})
	}
}

// Create new instance for authenticated system user
func (s *server) CreateMyInstance() http.HandlerFunc {
	type createInstanceRequest struct {
		Name              string `json:"name"`
		DestinationNumber string `json:"destination_number"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		systemUserID, ok := r.Context().Value("system_user_id").(int)
		if !ok {
			log.Warn().Msg("CreateMyInstance: system_user_id not found in context")
			s.respondWithJSON(w, http.StatusUnauthorized, map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"error":   "unauthorized",
				"success": false,
			})
			return
		}

		log.Info().Int("system_user_id", systemUserID).Msg("CreateMyInstance: Request received")

		var req createInstanceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error().Err(err).Msg("CreateMyInstance: Failed to decode request")
			s.respondWithJSON(w, http.StatusBadRequest, map[string]interface{}{
				"code":    http.StatusBadRequest,
				"error":   "invalid request payload",
				"success": false,
			})
			return
		}

		if req.Name == "" {
			s.respondWithJSON(w, http.StatusBadRequest, map[string]interface{}{
				"code":    http.StatusBadRequest,
				"error":   "name is required",
				"success": false,
			})
			return
		}
		
		// Check subscription limits
		canCreate, err := s.CanCreateInstance(systemUserID)
		if err != nil {
			log.Error().Err(err).Int("system_user_id", systemUserID).Msg("Failed to check instance limit")
			s.respondWithJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"code":    http.StatusInternalServerError,
				"error":   "failed to check subscription limit",
				"success": false,
			})
			return
		}
		
		if !canCreate {
			subscription, _ := s.GetActiveSubscription(systemUserID)
			var message string
			if subscription != nil && subscription.ExpiresAt != nil && subscription.ExpiresAt.Before(time.Now()) {
				message = "Seu plano expirou! Assine um dos nossos planos para continuar usando o sistema."
			} else if subscription != nil {
				message = "Você atingiu o limite de instâncias do seu plano. Faça upgrade para criar mais."
			} else {
				message = "Nenhuma assinatura ativa encontrada. Por favor, assine um plano."
			}
			
			s.respondWithJSON(w, http.StatusForbidden, map[string]interface{}{
				"code":    http.StatusForbidden,
				"error":   message,
				"success": false,
			})
			return
		}

		// Generate ID and token automatically
		id, err := GenerateRandomID()
		if err != nil {
			log.Error().Err(err).Msg("Failed to generate random ID")
			s.respondWithJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"code":    http.StatusInternalServerError,
				"error":   "failed to generate instance ID",
				"success": false,
			})
			return
		}

		token, err := GenerateRandomID()
		if err != nil {
			log.Error().Err(err).Msg("Failed to generate token")
			s.respondWithJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"code":    http.StatusInternalServerError,
				"error":   "failed to generate instance token",
				"success": false,
			})
			return
		}

		// Insert instance with auto-generated token
		_, err = s.db.Exec(
			`INSERT INTO users (id, name, token, webhook, jid, qrcode, system_user_id, destination_number, events) 
			 VALUES ($1, $2, $3, '', '', '', $4, $5, 'Message')`,
			id, req.Name, token, systemUserID, req.DestinationNumber,
		)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create instance")
			s.respondWithJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"code":    http.StatusInternalServerError,
				"error":   "failed to create instance",
				"success": false,
			})
			return
		}

		log.Info().
			Int("system_user_id", systemUserID).
			Str("instance_id", id).
			Str("instance_name", req.Name).
			Msg("Instance created successfully with auto-generated token")

		s.respondWithJSON(w, http.StatusCreated, map[string]interface{}{
			"code": http.StatusCreated,
			"data": map[string]interface{}{
				"id":                 id,
				"name":               req.Name,
				"token":              token,
				"destination_number": req.DestinationNumber,
				"message":            "Token gerado automaticamente. Use-o para acessar a API.",
			},
			"success": true,
		})
	}
}

// Update instance for authenticated system user
func (s *server) UpdateMyInstance() http.HandlerFunc {
	type updateInstanceRequest struct {
		Name              string `json:"name"`
		DestinationNumber string `json:"destination_number"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		systemUserID, ok := r.Context().Value("system_user_id").(int)
		if !ok {
			s.respondWithJSON(w, http.StatusUnauthorized, map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"error":   "unauthorized",
				"success": false,
			})
			return
		}

		vars := mux.Vars(r)
		instanceID := vars["id"]

		var req updateInstanceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.respondWithJSON(w, http.StatusBadRequest, map[string]interface{}{
				"code":    http.StatusBadRequest,
				"error":   "invalid request payload",
				"success": false,
			})
			return
		}

		// Verify ownership
		var count int
		err := s.db.Get(&count, "SELECT COUNT(*) FROM users WHERE id = $1 AND system_user_id = $2", instanceID, systemUserID)
		if err != nil || count == 0 {
			s.respondWithJSON(w, http.StatusNotFound, map[string]interface{}{
				"code":    http.StatusNotFound,
				"error":   "instance not found",
				"success": false,
			})
			return
		}

		// Update instance
		_, err = s.db.Exec(
			`UPDATE users SET name = $1, destination_number = $2 WHERE id = $3`,
			req.Name, req.DestinationNumber, instanceID,
		)
		if err != nil {
			log.Error().Err(err).Msg("Failed to update instance")
			s.respondWithJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"code":    http.StatusInternalServerError,
				"error":   "failed to update instance",
				"success": false,
			})
			return
		}

		s.respondWithJSON(w, http.StatusOK, map[string]interface{}{
			"code":    http.StatusOK,
			"message": "instance updated successfully",
			"success": true,
		})
	}
}

// Delete instance for authenticated system user
func (s *server) DeleteMyInstance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		systemUserID, ok := r.Context().Value("system_user_id").(int)
		if !ok {
			s.respondWithJSON(w, http.StatusUnauthorized, map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"error":   "unauthorized",
				"success": false,
			})
			return
		}

		vars := mux.Vars(r)
		instanceID := vars["id"]

		// Verify ownership
		var count int
		err := s.db.Get(&count, "SELECT COUNT(*) FROM users WHERE id = $1 AND system_user_id = $2", instanceID, systemUserID)
		if err != nil || count == 0 {
			s.respondWithJSON(w, http.StatusNotFound, map[string]interface{}{
				"code":    http.StatusNotFound,
				"error":   "instance not found",
				"success": false,
			})
			return
		}

		// Disconnect if connected
		if client := clientManager.GetWhatsmeowClient(instanceID); client != nil {
			if client.IsConnected() {
				client.Logout(r.Context())
			}
			client.Disconnect()
			clientManager.DeleteWhatsmeowClient(instanceID)
		}

		// Delete instance
		_, err = s.db.Exec("DELETE FROM users WHERE id = $1", instanceID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to delete instance")
			s.respondWithJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"code":    http.StatusInternalServerError,
				"error":   "failed to delete instance",
				"success": false,
			})
			return
		}

		s.respondWithJSON(w, http.StatusOK, map[string]interface{}{
			"code":    http.StatusOK,
			"message": "instance deleted successfully",
			"success": true,
		})
	}
}

// Get instance details for authenticated system user
func (s *server) GetMyInstance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		systemUserID, ok := r.Context().Value("system_user_id").(int)
		if !ok {
			s.respondWithJSON(w, http.StatusUnauthorized, map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"error":   "unauthorized",
				"success": false,
			})
			return
		}

		vars := mux.Vars(r)
		instanceID := vars["id"]

		// Verify ownership and get instance
		var instance struct {
			Id                string `db:"id" json:"id"`
			Name              string `db:"name" json:"name"`
			Token             string `db:"token" json:"token"`
			Jid               string `db:"jid" json:"jid"`
			DestinationNumber string `db:"destination_number" json:"destination_number"`
		}

		err := s.db.Get(&instance, `
			SELECT id, name, token, jid, destination_number 
			FROM users 
			WHERE id = $1 AND system_user_id = $2`,
			instanceID, systemUserID)

		if err != nil {
			if err == sql.ErrNoRows {
				s.respondWithJSON(w, http.StatusNotFound, map[string]interface{}{
					"code":    http.StatusNotFound,
					"error":   "instance not found",
					"success": false,
				})
				return
			}
			s.respondWithJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"code":    http.StatusInternalServerError,
				"error":   "database error",
				"success": false,
			})
			return
		}

		// Add connection status
		connected := false
		loggedIn := false
		if client := clientManager.GetWhatsmeowClient(instanceID); client != nil {
			connected = client.IsConnected()
			loggedIn = client.IsLoggedIn()
		}

		response := map[string]interface{}{
			"id":                 instance.Id,
			"name":               instance.Name,
			"token":              instance.Token,
			"jid":                instance.Jid,
			"destination_number": instance.DestinationNumber,
			"connected":          connected,
			"logged_in":          loggedIn,
		}

		s.respondWithJSON(w, http.StatusOK, map[string]interface{}{
			"code":    http.StatusOK,
			"data":    response,
			"success": true,
		})
	}
}

// Helper function to check if user owns instance
func (s *server) userOwnsInstance(instanceID string, systemUserID int) bool {
	var count int
	err := s.db.Get(&count, "SELECT COUNT(*) FROM users WHERE id = $1 AND system_user_id = $2", instanceID, systemUserID)
	return err == nil && count > 0
}

// Get profile for authenticated system user
func (s *server) GetMyProfile() http.HandlerFunc {
	type profileResponse struct {
		ID              int       `db:"id" json:"id"`
		Email           string    `db:"email" json:"email"`
		Name            string    `db:"name" json:"name"`
		WhatsappNumber  string    `db:"whatsapp_number" json:"whatsapp_number"`
		CreatedAt       time.Time `db:"created_at" json:"created_at"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		systemUserID, ok := r.Context().Value("system_user_id").(int)
		if !ok {
			s.respondWithJSON(w, http.StatusUnauthorized, map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"error":   "unauthorized",
				"success": false,
			})
			return
		}

		var profile profileResponse
		err := s.db.Get(&profile, `
			SELECT id, email, 
				   COALESCE(name, '') as name, 
				   COALESCE(whatsapp_number, '') as whatsapp_number, 
				   created_at 
			FROM system_users 
			WHERE id = $1`,
			systemUserID)

		if err != nil {
			if err == sql.ErrNoRows {
				s.respondWithJSON(w, http.StatusNotFound, map[string]interface{}{
					"code":    http.StatusNotFound,
					"error":   "user not found",
					"success": false,
				})
				return
			}
			log.Error().Err(err).Int("system_user_id", systemUserID).Msg("Failed to get profile")
			s.respondWithJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"code":    http.StatusInternalServerError,
				"error":   "database error",
				"success": false,
			})
			return
		}

		s.respondWithJSON(w, http.StatusOK, map[string]interface{}{
			"code":    http.StatusOK,
			"data":    profile,
			"success": true,
		})
	}
}

// Update profile for authenticated system user
func (s *server) UpdateMyProfile() http.HandlerFunc {
	type updateProfileRequest struct {
		Name           string `json:"name"`
		WhatsappNumber string `json:"whatsapp_number"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		systemUserID, ok := r.Context().Value("system_user_id").(int)
		if !ok {
			s.respondWithJSON(w, http.StatusUnauthorized, map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"error":   "unauthorized",
				"success": false,
			})
			return
		}

		var req updateProfileRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.respondWithJSON(w, http.StatusBadRequest, map[string]interface{}{
				"code":    http.StatusBadRequest,
				"error":   "invalid request payload",
				"success": false,
			})
			return
		}

		// Update profile
		_, err := s.db.Exec(
			`UPDATE system_users 
			 SET name = $1, whatsapp_number = $2, updated_at = CURRENT_TIMESTAMP 
			 WHERE id = $3`,
			req.Name, req.WhatsappNumber, systemUserID,
		)
		if err != nil {
			log.Error().Err(err).Int("system_user_id", systemUserID).Msg("Failed to update profile")
			s.respondWithJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"code":    http.StatusInternalServerError,
				"error":   "failed to update profile",
				"success": false,
			})
			return
		}

		s.respondWithJSON(w, http.StatusOK, map[string]interface{}{
			"code":    http.StatusOK,
			"message": "profile updated successfully",
			"success": true,
		})
	}
}
