package handlers

import (
	"net/http"
	"time"

	"task-management-system/Backend/config"
	"task-management-system/Backend/models"
	"task-management-system/Backend/utils"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocket upgrader for real-time updates
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins (modify for production)
	},
}

// CreateTask handles creating a new task with JWT authentication and database storage
func CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task data: " + err.Error()})
		return
	}

	// Validate task fields (e.g., ensure title and description are not empty)
	if task.Title == "" || task.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title and description are required"})
		return
	}

	// Get user ID from JWT token
	userID := utils.GetUserIDFromToken(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid or missing token"})
		return
	}
	task.AssignedTo = uint(userID) // Ensure compatibility with uint in models.Task

	// Save task to database
	query := "INSERT INTO tasks (title, description, assigned_to, status, priority, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	err := config.DB.QueryRow(query, task.Title, task.Description, task.AssignedTo, task.Status, task.Priority, time.Now(), time.Now()).Scan(&task.ID)
	if err != nil {
		log.Printf("Database error creating task: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task: " + err.Error()})
		return
	}

	// Broadcast task update via WebSocket (implement broadcast logic here)
	// Example: broadcastTaskUpdate(task)

	c.JSON(http.StatusCreated, task)
}

// GetTasks retrieves all tasks for the authenticated user
func GetTasks(c *gin.Context) {
	userID := utils.GetUserIDFromToken(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid or missing token"})
		return
	}

	// Retrieve tasks from database
	rows, err := config.DB.Query("SELECT id, title, description, assigned_to, status, priority, created_at, updated_at FROM tasks WHERE assigned_to = $1", userID)
	if err != nil {
		log.Printf("Database error retrieving tasks: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks: " + err.Error()})
		return
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.AssignedTo, &task.Status, &task.Priority, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			log.Printf("Error scanning task row: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing tasks: " + err.Error()})
			return
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Rows error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing tasks: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

// UpdateTask updates an existing task
func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task
	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task data: " + err.Error()})
		return
	}

	userID := utils.GetUserIDFromToken(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid or missing token"})
		return
	}

	// Update task in database
	query := "UPDATE tasks SET title = $1, description = $2, status = $3, priority = $4, updated_at = $5 WHERE id = $6 AND assigned_to = $7"
	result, err := config.DB.Exec(query, task.Title, task.Description, task.Status, task.Priority, time.Now(), id, userID)
	if err != nil {
		log.Printf("Database error updating task: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task: " + err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found or unauthorized"})
		return
	}

	// Broadcast task update via WebSocket (implement broadcast logic here)
	// Example: broadcastTaskUpdate(task)

	c.JSON(http.StatusOK, gin.H{"message": "Task updated"})
}

// DeleteTask deletes a task
func DeleteTask(c *gin.Context) {
	id := c.Param("id")

	userID := utils.GetUserIDFromToken(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid or missing token"})
		return
	}

	// Delete task from database
	query := "DELETE FROM tasks WHERE id = $1 AND assigned_to = $2"
	result, err := config.DB.Exec(query, id, userID)
	if err != nil {
		log.Printf("Database error deleting task: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task: " + err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found or unauthorized"})
		return
	}

	// Broadcast task deletion via WebSocket (implement broadcast logic here)
	// Example: broadcastTaskDeletion(id)

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}

// WSHandler handles WebSocket connections for real-time updates
func WSHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "WebSocket upgrade failed: " + err.Error()})
		return
	}
	defer conn.Close()

	// Get user ID from JWT token for authentication
	userID := utils.GetUserIDFromToken(c)
	if userID == 0 {
		if err := conn.WriteMessage(websocket.TextMessage, []byte("Unauthorized: Invalid or missing token")); err != nil {
			log.Printf("WebSocket write error: %v", err)
		}
		return
	}

	// Handle WebSocket connections for real-time updates
	// Implement logic to broadcast task updates to connected clients
	// Example: Add conn to a list of connections and broadcast updates
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket unexpected close for user %d: %v", userID, err)
			}
			break
		}
		// Broadcast task updates here (e.g., when a task is created/updated/deleted)
	}
}