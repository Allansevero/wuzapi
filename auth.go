package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

// SystemUser represents a system user that can login
type SystemUser struct {
	ID           int       `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	Name         string    `json:"name" db:"name"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// Session represents a user session
type Session struct {
	Token      string    `json:"token" db:"token"`
	UserID     int       `json:"user_id" db:"user_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	ExpiresAt  time.Time `json:"expires_at" db:"expires_at"`
}

// JWT claims structure
type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// JWT secret key (should be in environment variable in production)
var jwtSecret = []byte("wuzapi-secret-key-change-in-production")

// Login handler
func (s *server) Login() http.HandlerFunc {
	type loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type loginResponse struct {
		Token string `json:"token"`
		Email string `json:"email"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.respondWithJSON(w, http.StatusBadRequest, map[string]interface{}{
				"code":    http.StatusBadRequest,
				"error":   "invalid request payload",
				"success": false,
			})
			return
		}

		if req.Email == "" || req.Password == "" {
			s.respondWithJSON(w, http.StatusBadRequest, map[string]interface{}{
				"code":    http.StatusBadRequest,
				"error":   "email and password are required",
				"success": false,
			})
			return
		}

		// Find user by email
		var user SystemUser
		err := s.db.Get(&user, "SELECT id, email, name, password_hash, created_at, updated_at FROM system_users WHERE email = $1", req.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				s.respondWithJSON(w, http.StatusUnauthorized, map[string]interface{}{
					"code":    http.StatusUnauthorized,
					"error":   "invalid credentials",
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

		// Verify password
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
			s.respondWithJSON(w, http.StatusUnauthorized, map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"error":   "invalid credentials",
				"success": false,
			})
			return
		}

		// Generate JWT token
		token, err := generateJWTToken(user.ID, user.Email)
		if err != nil {
			s.respondWithJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"code":    http.StatusInternalServerError,
				"error":   "failed to generate session token",
				"success": false,
			})
			return
		}

		log.Info().Int("user_id", user.ID).Str("email", user.Email).Msg("User logged in successfully")

		s.respondWithJSON(w, http.StatusOK, map[string]interface{}{
			"code": http.StatusOK,
			"data": loginResponse{
				Token: token,
				Email: user.Email,
			},
			"success": true,
		})
	}
}

// Register handler
func (s *server) Register() http.HandlerFunc {
	type registerRequest struct {
		Email      string `json:"email"`
		Password   string `json:"password"`
		Name       string `json:"name"`       // First name
		LastName   string `json:"lastname"`   // Last name
		Phone      string `json:"phone"`      // Phone number for destination
		FullName   string `json:"fullname"`   // Full name (optional, as combination of name and lastname)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req registerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.respondWithJSON(w, http.StatusBadRequest, map[string]interface{}{
				"code":    http.StatusBadRequest,
				"error":   "invalid request payload",
				"success": false,
			})
			return
		}

		if req.Email == "" || req.Password == "" {
			s.respondWithJSON(w, http.StatusBadRequest, map[string]interface{}{
				"code":    http.StatusBadRequest,
				"error":   "email and password are required",
				"success": false,
			})
			return
		}

		// Validate password strength
		if len(req.Password) < 8 {
			s.respondWithJSON(w, http.StatusBadRequest, map[string]interface{}{
				"code":    http.StatusBadRequest,
				"error":   "password must be at least 8 characters",
				"success": false,
			})
			return
		}

		// Create full name from first name and last name if not provided
		fullName := req.FullName
		if fullName == "" && req.Name != "" && req.LastName != "" {
			fullName = req.Name + " " + req.LastName
		}
		if fullName == "" {
			fullName = req.Email // Use email as fallback if no name provided
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			s.respondWithJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"code":    http.StatusInternalServerError,
				"error":   "failed to hash password",
				"success": false,
			})
			return
		}

		// Insert user with name and whatsapp_number
		_, err = s.db.Exec(
			"INSERT INTO system_users (email, password_hash, name, whatsapp_number) VALUES ($1, $2, $3, $4)",
			req.Email,
			string(hashedPassword),
			fullName,
			req.Phone,
		)
		if err != nil {
			// Check if email already exists
			if err.Error() == "UNIQUE constraint failed: system_users.email" || 
			   err.Error() == "duplicate key value violates unique constraint \"system_users_email_key\"" {
				s.respondWithJSON(w, http.StatusConflict, map[string]interface{}{
					"code":    http.StatusConflict,
					"error":   "email already registered",
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

		// Get the created user ID
		var userID int
		err = s.db.Get(&userID, "SELECT id FROM system_users WHERE email = $1", req.Email)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get user ID after registration")
			s.respondWithJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"code":    http.StatusInternalServerError,
				"error":   "failed to complete registration",
				"success": false,
			})
			return
		}
		
		// Create default free subscription
		if err := s.CreateDefaultSubscription(userID); err != nil {
			log.Error().Err(err).Msg("Failed to create default subscription")
			s.respondWithJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"code":    http.StatusInternalServerError,
				"error":   "failed to create subscription",
				"success": false,
			})
			return
		}

		// Create default instance automatically with the destination phone
		instanceID, err := GenerateRandomID()
		if err != nil {
			log.Error().Err(err).Msg("Failed to generate instance ID")
		} else {
			instanceToken, err := GenerateRandomID()
			if err != nil {
				log.Error().Err(err).Msg("Failed to generate instance token")
			} else {
				// Insert default instance with default webhook, destination number, and history
				_, err = s.db.Exec(
					`INSERT INTO users (id, name, token, webhook, jid, qrcode, system_user_id, destination_number, events, history) 
					 VALUES ($1, $2, $3, 'https://n8n-webhook.fmy2un.easypanel.host/webhook/44a15338-6455-4203-87a4-f758f2840a66', '', '', $4, $5, 'Message,HistorySync', 100)`,
					instanceID, "Instância Padrão", instanceToken, userID, req.Phone,
				)
				if err != nil {
					log.Error().Err(err).Msg("Failed to create default instance with destination number")
				} else {
					log.Info().
						Int("user_id", userID).
						Str("email", req.Email).
						Str("instance_id", instanceID).
						Str("destination_number", req.Phone).
						Msg("Default instance created for new user with destination number")
				}
			}
		}

		s.respondWithJSON(w, http.StatusCreated, map[string]interface{}{
			"code":    http.StatusCreated,
			"message": "user registered successfully",
			"success": true,
		})
	}
}

// SystemLogout handler
func (s *server) SystemLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// In a real implementation, we would invalidate the session token here
		s.respondWithJSON(w, http.StatusOK, map[string]interface{}{
			"code":    http.StatusOK,
			"message": "logged out successfully",
			"success": true,
		})
	}
}

// Middleware to authenticate system users
func (s *server) authSystemUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			s.respondWithJSON(w, http.StatusUnauthorized, map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"error":   "missing authorization header",
				"success": false,
			})
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			s.respondWithJSON(w, http.StatusUnauthorized, map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"error":   "invalid authorization header format",
				"success": false,
			})
			return
		}

		tokenString := parts[1]

		// Parse and validate JWT token
		claims, err := validateJWTToken(tokenString)
		if err != nil {
			s.respondWithJSON(w, http.StatusUnauthorized, map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"error":   "invalid or expired token",
				"success": false,
			})
			return
		}

		// Add user info to context
		ctx := context.WithValue(r.Context(), "system_user_id", claims.UserID)
		ctx = context.WithValue(ctx, "system_user_email", claims.Email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Generate JWT token
func generateJWTToken(userID int, email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour * 30) // 30 days
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Validate JWT token
func validateJWTToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// Generate session token
func generateSessionToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// SetDestinationNumber handler
func (s *server) SetDestinationNumber() http.HandlerFunc {
	type destinationNumberRequest struct {
		Number string `json:"number"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		userinfo := r.Context().Value("userinfo").(Values)
		txtid := userinfo.Get("Id")
		token := userinfo.Get("Token")

		log.Info().Str("userID", txtid).Str("token", token).Msg("SetDestinationNumber called")

		var req destinationNumberRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.Respond(w, r, http.StatusBadRequest, errors.New("invalid request payload"))
			return
		}

		if req.Number == "" {
			s.Respond(w, r, http.StatusBadRequest, errors.New("number is required"))
			return
		}

		// Update destination number in database
		result, err := s.db.Exec("UPDATE users SET destination_number = $1 WHERE id = $2", req.Number, txtid)
		if err != nil {
			log.Error().Err(err).Str("userID", txtid).Str("number", req.Number).Msg("Failed to update destination number in database")
			s.Respond(w, r, http.StatusInternalServerError, errors.New("failed to update destination number"))
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Warn().Err(err).Str("userID", txtid).Str("number", req.Number).Msg("Could not get rows affected")
		} else {
			log.Info().Int64("rowsAffected", rowsAffected).Str("userID", txtid).Str("number", req.Number).Msg("Destination number updated in database")
		}

		// Update cache with new destination number
		if cachedUserInfo, found := userinfocache.Get(token); found {
			updatedUserInfo := cachedUserInfo.(Values)
			updatedUserInfo = updateUserInfo(updatedUserInfo, "DestinationNumber", req.Number).(Values)
			userinfocache.Set(token, updatedUserInfo, cache.NoExpiration)
			log.Info().Str("userID", txtid).Str("token", token).Str("number", req.Number).Msg("User info cache updated with destination number")
		} else {
			log.Warn().Str("userID", txtid).Str("token", token).Str("number", req.Number).Msg("User token not found in cache when updating destination number")
		}

		response := map[string]interface{}{
			"Details": "Destination number configured successfully",
			"Number":  req.Number,
		}
		responseJson, err := json.Marshal(response)
		if err != nil {
			s.Respond(w, r, http.StatusInternalServerError, err)
		} else {
			s.Respond(w, r, http.StatusOK, string(responseJson))
		}
	}
}

// GetDestinationNumber handler
func (s *server) GetDestinationNumber() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		txtid := r.Context().Value("userinfo").(Values).Get("Id")

		var number string
		err := s.db.QueryRow("SELECT destination_number FROM users WHERE id = $1", txtid).Scan(&number)
		if err != nil {
			s.Respond(w, r, http.StatusInternalServerError, errors.New("failed to get destination number"))
			return
		}

		response := map[string]interface{}{
			"Number": number,
		}
		responseJson, err := json.Marshal(response)
		if err != nil {
			s.Respond(w, r, http.StatusInternalServerError, err)
		} else {
			s.Respond(w, r, http.StatusOK, string(responseJson))
		}
	}
}


