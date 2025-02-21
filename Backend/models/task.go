package models

import "time"

type Task struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	AssignedTo  uint      `json:"assigned_to"` // User ID
	Status      string    `json:"status"`      // e.g., "pending", "in_progress", "completed"
	Priority    string    `json:"priority"`    // e.g., "low", "medium", "high"
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}