package main

import (
	"database/sql"
	"fmt"
	"time"
)

// Plan represents a subscription plan
type Plan struct {
	ID           int       `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Price        float64   `json:"price" db:"price"`
	MaxInstances int       `json:"max_instances" db:"max_instances"`
	TrialDays    int       `json:"trial_days" db:"trial_days"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// UserSubscription represents a user's active subscription
type UserSubscription struct {
	ID           int       `json:"id" db:"id"`
	SystemUserID int       `json:"system_user_id" db:"system_user_id"`
	PlanID       int       `json:"plan_id" db:"plan_id"`
	StartedAt    time.Time `json:"started_at" db:"started_at"`
	ExpiresAt    *time.Time `json:"expires_at" db:"expires_at"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// SubscriptionHistory tracks historical subscriptions
type SubscriptionHistory struct {
	ID           int       `json:"id" db:"id"`
	SystemUserID int       `json:"system_user_id" db:"system_user_id"`
	PlanID       int       `json:"plan_id" db:"plan_id"`
	StartedAt    time.Time `json:"started_at" db:"started_at"`
	EndedAt      *time.Time `json:"ended_at" db:"ended_at"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// UserSubscriptionDetails combines subscription with plan details
type UserSubscriptionDetails struct {
	UserSubscription
	Plan Plan `json:"plan"`
}

// CreateDefaultSubscription creates a free subscription for new users
func (s *server) CreateDefaultSubscription(systemUserID int) error {
	// Get the free plan (ID 1)
	freePlanID := 1
	
	var query string
	if s.db.DriverName() == "postgres" {
		query = `
			INSERT INTO user_subscriptions (system_user_id, plan_id, started_at, is_active)
			VALUES ($1, $2, $3, $4)
			RETURNING id`
	} else {
		query = `
			INSERT INTO user_subscriptions (system_user_id, plan_id, started_at, is_active)
			VALUES (?, ?, ?, ?)`
	}
	
	now := time.Now()
	if s.db.DriverName() == "postgres" {
		var id int
		err := s.db.QueryRow(query, systemUserID, freePlanID, now, true).Scan(&id)
		return err
	}
	
	_, err := s.db.Exec(query, systemUserID, freePlanID, now, true)
	return err
}

// GetActiveSubscription returns the user's active subscription with plan details
func (s *server) GetActiveSubscription(systemUserID int) (*UserSubscriptionDetails, error) {
	var query string
	if s.db.DriverName() == "postgres" {
		query = `
			SELECT 
				us.id, us.system_user_id, us.plan_id, us.started_at, us.expires_at, 
				us.is_active, us.created_at, us.updated_at,
				p.id, p.name, p.price, p.max_instances, p.trial_days, p.is_active, p.created_at
			FROM user_subscriptions us
			JOIN plans p ON us.plan_id = p.id
			WHERE us.system_user_id = $1 AND us.is_active = true
			ORDER BY us.created_at DESC
			LIMIT 1`
	} else {
		query = `
			SELECT 
				us.id, us.system_user_id, us.plan_id, us.started_at, us.expires_at, 
				us.is_active, us.created_at, us.updated_at,
				p.id, p.name, p.price, p.max_instances, p.trial_days, p.is_active, p.created_at
			FROM user_subscriptions us
			JOIN plans p ON us.plan_id = p.id
			WHERE us.system_user_id = ? AND us.is_active = 1
			ORDER BY us.created_at DESC
			LIMIT 1`
	}
	
	var details UserSubscriptionDetails
	err := s.db.QueryRow(query, systemUserID).Scan(
		&details.ID, &details.SystemUserID, &details.PlanID, &details.StartedAt, &details.ExpiresAt,
		&details.IsActive, &details.CreatedAt, &details.UpdatedAt,
		&details.Plan.ID, &details.Plan.Name, &details.Plan.Price, &details.Plan.MaxInstances,
		&details.Plan.TrialDays, &details.Plan.IsActive, &details.Plan.CreatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get active subscription: %w", err)
	}
	
	return &details, nil
}

// UpdateSubscription updates or creates a new subscription
func (s *server) UpdateSubscription(systemUserID, planID int) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	
	// Deactivate current subscription
	var deactivateQuery string
	if s.db.DriverName() == "postgres" {
		deactivateQuery = `
			UPDATE user_subscriptions 
			SET is_active = false, updated_at = $1
			WHERE system_user_id = $2 AND is_active = true`
	} else {
		deactivateQuery = `
			UPDATE user_subscriptions 
			SET is_active = 0, updated_at = ?
			WHERE system_user_id = ? AND is_active = 1`
	}
	
	now := time.Now()
	if _, err := tx.Exec(deactivateQuery, now, systemUserID); err != nil {
		return fmt.Errorf("failed to deactivate current subscription: %w", err)
	}
	
	// Create new subscription
	var insertQuery string
	if s.db.DriverName() == "postgres" {
		insertQuery = `
			INSERT INTO user_subscriptions (system_user_id, plan_id, started_at, is_active)
			VALUES ($1, $2, $3, $4)`
	} else {
		insertQuery = `
			INSERT INTO user_subscriptions (system_user_id, plan_id, started_at, is_active)
			VALUES (?, ?, ?, ?)`
	}
	
	if _, err := tx.Exec(insertQuery, systemUserID, planID, now, true); err != nil {
		return fmt.Errorf("failed to create new subscription: %w", err)
	}
	
	return tx.Commit()
}

// CheckSubscriptionExpired checks if subscription has expired and deactivates it
func (s *server) CheckSubscriptionExpired(systemUserID int) error {
	var query string
	if s.db.DriverName() == "postgres" {
		query = `
			UPDATE user_subscriptions 
			SET is_active = false, updated_at = $1
			WHERE system_user_id = $2 
			AND is_active = true 
			AND expires_at IS NOT NULL 
			AND expires_at < $1`
	} else {
		query = `
			UPDATE user_subscriptions 
			SET is_active = 0, updated_at = ?
			WHERE system_user_id = ? 
			AND is_active = 1 
			AND expires_at IS NOT NULL 
			AND expires_at < ?`
	}
	
	now := time.Now()
	_, err := s.db.Exec(query, now, systemUserID, now)
	return err
}

// GetUserInstanceCount returns the number of instances for a user
func (s *server) GetUserInstanceCount(systemUserID int) (int, error) {
	var count int
	var query string
	
	if s.db.DriverName() == "postgres" {
		query = `SELECT COUNT(*) FROM users WHERE system_user_id = $1`
	} else {
		query = `SELECT COUNT(*) FROM users WHERE system_user_id = ?`
	}
	
	err := s.db.Get(&count, query, systemUserID)
	if err != nil {
		return 0, fmt.Errorf("failed to count instances: %w", err)
	}
	
	return count, nil
}

// GetUserConnectedInstanceCount returns the number of CONNECTED instances for a user
func (s *server) GetUserConnectedInstanceCount(systemUserID int) (int, error) {
	var count int
	var query string
	
	if s.db.DriverName() == "postgres" {
		query = `SELECT COUNT(*) FROM users WHERE system_user_id = $1 AND connected = 1`
	} else {
		query = `SELECT COUNT(*) FROM users WHERE system_user_id = ? AND connected = 1`
	}
	
	err := s.db.Get(&count, query, systemUserID)
	if err != nil {
		return 0, fmt.Errorf("failed to count connected instances: %w", err)
	}
	
	return count, nil
}

// CanCreateInstance checks if user can create more instances based on their plan
// This checks CONNECTED instances only
func (s *server) CanCreateInstance(systemUserID int) (bool, error) {
	// Check and update expired subscription
	if err := s.CheckSubscriptionExpired(systemUserID); err != nil {
		return false, fmt.Errorf("failed to check subscription expiration: %w", err)
	}
	
	// Get active subscription
	subscription, err := s.GetActiveSubscription(systemUserID)
	if err != nil {
		return false, fmt.Errorf("failed to get active subscription: %w", err)
	}
	
	// Se não tem assinatura, criar plano gratuito automaticamente
	if subscription == nil {
		if err := s.CreateDefaultSubscription(systemUserID); err != nil {
			return false, fmt.Errorf("failed to create default subscription: %w", err)
		}
		// Tentar buscar novamente
		subscription, err = s.GetActiveSubscription(systemUserID)
		if err != nil || subscription == nil {
			return false, fmt.Errorf("failed to get subscription after creation")
		}
	}
	
	// Check if subscription is expired
	if subscription.ExpiresAt != nil && subscription.ExpiresAt.Before(time.Now()) {
		// Plano expirado - bloquear usuário
		return false, nil
	}
	
	// Get current CONNECTED instance count
	count, err := s.GetUserConnectedInstanceCount(systemUserID)
	if err != nil {
		return false, fmt.Errorf("failed to get connected instance count: %w", err)
	}
	
	// Check if can create more instances based on connected count
	return count < subscription.Plan.MaxInstances, nil
}

// GetAllPlans returns all active plans
func (s *server) GetAllPlans() ([]Plan, error) {
	var plans []Plan
	var query string
	
	if s.db.DriverName() == "postgres" {
		query = `SELECT id, name, price, max_instances, trial_days, is_active, created_at 
				 FROM plans WHERE is_active = true ORDER BY price ASC`
	} else {
		query = `SELECT id, name, price, max_instances, trial_days, is_active, created_at 
				 FROM plans WHERE is_active = 1 ORDER BY price ASC`
	}
	
	err := s.db.Select(&plans, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get plans: %w", err)
	}
	
	return plans, nil
}

// AddSubscriptionHistory adds a record to subscription history
func (s *server) AddSubscriptionHistory(systemUserID, planID int, startedAt time.Time, endedAt *time.Time) error {
	var query string
	if s.db.DriverName() == "postgres" {
		query = `
			INSERT INTO subscription_history (system_user_id, plan_id, started_at, ended_at)
			VALUES ($1, $2, $3, $4)`
	} else {
		query = `
			INSERT INTO subscription_history (system_user_id, plan_id, started_at, ended_at)
			VALUES (?, ?, ?, ?)`
	}
	
	_, err := s.db.Exec(query, systemUserID, planID, startedAt, endedAt)
	return err
}

// HasValidSubscription checks if user has an active and non-expired subscription
func (s *server) HasValidSubscription(systemUserID int) (bool, error) {
// Check and update expired subscriptions
if err := s.CheckSubscriptionExpired(systemUserID); err != nil {
return false, fmt.Errorf("failed to check subscription expiration: %w", err)
}

// Get active subscription
subscription, err := s.GetActiveSubscription(systemUserID)
if err != nil {
return false, fmt.Errorf("failed to get subscription: %w", err)
}

// No subscription found
if subscription == nil {
return false, nil
}

// Check if subscription is expired
if subscription.ExpiresAt != nil && subscription.ExpiresAt.Before(time.Now()) {
return false, nil
}

// Has valid subscription
return true, nil
}
