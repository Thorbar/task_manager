package task

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"task-manager/backend-go/db"
	"task-manager/backend-go/internal/auth"
	"task-manager/backend-go/internal/i18n"
	"time"
)

// TasksRouter routes HTTP methods to appropriate handlers.
func TasksRouter(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetTasksHandler(w, r)
	case http.MethodPost:
		CreateTaskHandler(w, r)
	case http.MethodPut:
		UpdateTaskHandler(w, r)
	case http.MethodDelete:
		DeleteTaskHandler(w, r)
	default:
		http.Error(w, i18n.T("error.method_not_allowed"), http.StatusMethodNotAllowed)
	}
}

// GetTasksHandler returns all tasks for the authenticated user.
func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromAuthHeader(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	rows, err := db.DB.Query(`
		SELECT t.id, t.title, t.description, t.status, t.created_at, t.user_id, u.username
		FROM tasks t
		JOIN users u ON t.user_id = u.id
		WHERE u.id = ?`, userID)
	if err != nil {
		http.Error(w, i18n.T("error.query_tasks_failed"), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(
			&task.ID, &task.Title, &task.Description, &task.Status,
			&task.CreatedAt, &task.UserID, &task.Username,
		); err != nil {
			http.Error(w, i18n.T("error.read_tasks_failed"), http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"tasks": tasks})
}

// CreateTaskHandler creates a new task for the authenticated user.
func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromAuthHeader(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var newTask Task
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		http.Error(w, i18n.T("error.invalid_data"), http.StatusBadRequest)
		return
	}

	createdAt := time.Now()
	query := `
		INSERT INTO tasks (title, description, status, created_at, user_id)
		VALUES (?, ?, 'pending', ?, ?)
	`
	res, err := db.DB.Exec(query, newTask.Title, newTask.Description, createdAt, userID)
	if err != nil {
		http.Error(w, i18n.T("error.create_task_failed"), http.StatusInternalServerError)
		return
	}

	taskID, err := res.LastInsertId()
	if err != nil {
		http.Error(w, i18n.T("error.retrieve_task_id_failed"), http.StatusInternalServerError)
		return
	}

	var created Task
	query = `
		SELECT t.id, t.title, t.description, t.status, t.created_at, t.user_id, u.username
		FROM tasks t
		JOIN users u ON t.user_id = u.id
		WHERE t.id = ?
	`
	err = db.DB.QueryRow(query, taskID).Scan(
		&created.ID,
		&created.Title,
		&created.Description,
		&created.Status,
		&created.CreatedAt,
		&created.UserID,
		&created.Username,
	)
	if err != nil {
		http.Error(w, i18n.T("error.get_task_failed"), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

// UpdateTaskHandler updates an existing task owned by the authenticated user.
func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromAuthHeader(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/tasks/")
	idStr = strings.TrimSuffix(idStr, "/update")
	taskID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, i18n.T("error.invalid_id"), http.StatusBadRequest)
		return
	}

	var updatedTask Task
	if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
		http.Error(w, i18n.T("error.invalid_data"), http.StatusBadRequest)
		return
	}

	query := `
		UPDATE tasks
		SET title = ?, description = ?, status = ?
		WHERE id = ? AND user_id = ?
	`
	_, err = db.DB.Exec(query, updatedTask.Title, updatedTask.Description, updatedTask.Status, taskID, userID)
	if err != nil {
		http.Error(w, i18n.T("error.update_task_failed"), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": i18n.T("message.task_updated")})
}

// DeleteTaskHandler deletes a task owned by the authenticated user.
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromAuthHeader(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, i18n.T("error.invalid_id"), http.StatusBadRequest)
		return
	}
	taskID, err := strconv.Atoi(parts[3])
	if err != nil {
		http.Error(w, i18n.T("error.invalid_id"), http.StatusBadRequest)
		return
	}

	query := `DELETE FROM tasks WHERE id = ? AND user_id = ?`
	res, err := db.DB.Exec(query, taskID, userID)
	if err != nil {
		http.Error(w, i18n.T("error.delete_task_failed"), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, i18n.T("error.task_not_found_or_not_owned"), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": i18n.T("message.task_deleted")})
}

// getUserIDFromAuthHeader extracts and validates the user ID from the Authorization header.
func getUserIDFromAuthHeader(r *http.Request) (int, error) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return 0, &authError{i18n.T("error.token_not_provided")}
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	userID, err := auth.ParseToken(token)
	if err != nil {
		return 0, &authError{i18n.T("error.invalid_token")}
	}

	return userID, nil
}

// authError implements error interface for authentication errors.
type authError struct {
	message string
}

func (e *authError) Error() string {
	return e.message
}
