package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
)

// DailyConversation stores conversation data for a day
type DailyConversation struct {
	ID        int             `json:"id" db:"id"`
	UserID    string          `json:"user_id" db:"user_id"`
	Date      string          `json:"date" db:"date"`
	ChatJID   string          `json:"chat_jid" db:"chat_jid"`
	Contact   string          `json:"contact" db:"contact"`
	Messages  json.RawMessage `json:"messages" db:"messages"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
}

// DailyWebhookPayload is the structure sent to the webhook
type DailyWebhookPayload struct {
	InstanceID    string                   `json:"instance_id"`
	Date          string                   `json:"date"`
	Conversations []ConversationData       `json:"conversations"`
	EnviarPara    string                   `json:"enviar_para"`
}

// ConversationData represents a conversation with a contact
type ConversationData struct {
	Contact  string        `json:"contact"`
	Messages []interface{} `json:"messages"`
}

// Initialize daily message sender cron job
func (s *server) initDailyMessageSender() {
	c := cron.New(cron.WithLocation(getBrasiliaLocation()))

	// Schedule for 18:00 (6 PM) Brazil time
	_, err := c.AddFunc("0 18 * * *", func() {
		log.Info().Msg("Starting daily message delivery at 18:00 Brasilia time")
		s.sendDailyMessages()
	})

	if err != nil {
		log.Error().Err(err).Msg("Failed to schedule daily message sender")
		return
	}

	c.Start()
	log.Info().Msg("Daily message sender cron job initialized")
}

// Get Brasilia timezone location
func getBrasiliaLocation() *time.Location {
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		log.Warn().Err(err).Msg("Failed to load America/Sao_Paulo timezone, using UTC-3")
		return time.FixedZone("BRT", -3*60*60)
	}
	return loc
}

// Send daily messages for all instances
func (s *server) sendDailyMessages() {
	today := time.Now().In(getBrasiliaLocation()).Format("2006-01-02")

	// Get all active instances
	rows, err := s.db.Query("SELECT DISTINCT id, destination_number FROM users WHERE jid != ''")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get users for daily message delivery")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var instanceID, destinationNumber string
		if err := rows.Scan(&instanceID, &destinationNumber); err != nil {
			log.Error().Err(err).Msg("Failed to scan user row")
			continue
		}

		// Send daily messages for this instance
		if err := s.sendDailyMessagesForInstance(instanceID, today, destinationNumber); err != nil {
			log.Error().
				Err(err).
				Str("instance_id", instanceID).
				Msg("Failed to send daily messages for instance")
		}
	}
}

// Send daily messages for a specific instance
func (s *server) sendDailyMessagesForInstance(instanceID, date, destinationNumber string) error {
	// Get all conversations for today
	var query string
	if s.db.DriverName() == "postgres" {
		query = `
			SELECT chat_jid, sender_jid, message_type, text_content, media_link, timestamp, datajson
			FROM message_history
			WHERE user_id = $1 AND DATE(timestamp) = $2
			ORDER BY chat_jid, timestamp ASC`
	} else {
		query = `
			SELECT chat_jid, sender_jid, message_type, text_content, media_link, timestamp, datajson
			FROM message_history
			WHERE user_id = ? AND DATE(timestamp) = ?
			ORDER BY chat_jid, timestamp ASC`
	}

	rows, err := s.db.Query(query, instanceID, date)
	if err != nil {
		return fmt.Errorf("failed to query messages: %w", err)
	}
	defer rows.Close()

	// Group messages by conversation
	conversations := make(map[string][]interface{})
	for rows.Next() {
		var chatJID, senderJID, messageType, textContent, mediaLink, timestamp, dataJson string
		if err := rows.Scan(&chatJID, &senderJID, &messageType, &textContent, &mediaLink, &timestamp, &dataJson); err != nil {
			log.Error().Err(err).Msg("Failed to scan message row")
			continue
		}

		message := map[string]interface{}{
			"sender_jid":   senderJID,
			"message_type": messageType,
			"text_content": textContent,
			"media_link":   mediaLink,
			"timestamp":    timestamp,
		}

		if dataJson != "" {
			var jsonData interface{}
			if err := json.Unmarshal([]byte(dataJson), &jsonData); err == nil {
				message["data"] = jsonData
			}
		}

		conversations[chatJID] = append(conversations[chatJID], message)
	}

	// If no messages today, skip
	if len(conversations) == 0 {
		log.Info().
			Str("instance_id", instanceID).
			Str("date", date).
			Msg("No messages to send for today")
		return nil
	}

	// Build payload
	conversationsData := make([]ConversationData, 0, len(conversations))
	for contact, messages := range conversations {
		conversationsData = append(conversationsData, ConversationData{
			Contact:  contact,
			Messages: messages,
		})
	}

	payload := DailyWebhookPayload{
		InstanceID:    instanceID,
		Date:          date,
		Conversations: conversationsData,
		EnviarPara:    destinationNumber,
	}

	// Send to webhook
	return sendToWebhook(FIXED_WEBHOOK_URL, payload)
}

// Send payload to webhook
func sendToWebhook(webhookURL string, payload interface{}) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req = req.WithContext(ctx)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	log.Info().
		Str("webhook_url", webhookURL).
		Int("status", resp.StatusCode).
		Msg("Successfully sent daily messages to webhook")

	return nil
}

// Handler for manual test send
func (s *server) handleManualDailySend(w http.ResponseWriter, r *http.Request) {
	txtid := r.Context().Value("userinfo").(Values).Get("Id")
	
	if txtid == "" {
		log.Info().Msg("Unauthorized access attempt")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	// Get instance ID from query params or use the user's ID
	instanceID := r.URL.Query().Get("instance_id")
	if instanceID == "" {
		instanceID = txtid
	}

	// Get date from query params or use today
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		dateStr = time.Now().In(getBrasiliaLocation()).Format("2006-01-02")
	}

	// Get destination number
	var destinationNumber string
	err := s.db.QueryRow("SELECT destination_number FROM users WHERE id = ?", instanceID).Scan(&destinationNumber)
	if err != nil {
		log.Error().Err(err).Str("instance_id", instanceID).Msg("Failed to get destination number")
		destinationNumber = ""
	}

	log.Info().
		Str("instance_id", instanceID).
		Str("date", dateStr).
		Str("user_id", txtid).
		Msg("Manual daily send triggered")

	// Send daily messages for this instance
	err = s.sendDailyMessagesForInstance(instanceID, dateStr, destinationNumber)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send manual daily messages")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Daily messages sent successfully",
		"instance_id": instanceID,
		"date": dateStr,
	})
}
