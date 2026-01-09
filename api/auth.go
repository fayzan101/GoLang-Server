package api

import (
	"net/http"
	"encoding/json"
	"myapp/internal"
	"golang.org/x/crypto/bcrypt"
	"log"
)
// ForgotPasswordHandler initiates password reset (stub: just responds OK)
func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"message": "Invalid request payload",
		})
		return
	}
	// Generate a dummy token (replace with secure token in production)
	token := "dummy-reset-token" // TODO: generate a secure token and store it
	err := internal.SendResetEmail(req.Email, token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"message": "Failed to send reset email",
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"message": "Password reset link sent",
	})
}

// ResetPasswordHandler resets the user's password
func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email       string `json:"email"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"message": "Invalid request payload",
		})
		return
	}
	var user User
	if err := internal.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"message": "User not found",
		})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"message": "Failed to hash password",
		})
		return
	}
	user.Password = string(hashedPassword)
	if err := internal.DB.Save(&user).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"message": "Failed to reset password",
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"message": "Password reset successful",
	})
}



// ListUsersHandler returns users filtered by type
func ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Type string `json:"type"`
	}
	if r.Method == http.MethodPost {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status": "error",
				"message": "Invalid request payload",
			})
			return
		}
	}
	var users []User
	query := internal.DB
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}
	if err := query.Find(&users).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"message": "Failed to fetch users",
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data": users,
	})
}

// RegistrationRequest represents the registration payload
type RegistrationRequest struct {
	FullName   string `json:"full_name"`
	DOB        string `json:"dob"`
	University string `json:"university"`
	Semester   string `json:"semester"`
	Program    string `json:"program"`
	RollNo     string `json:"roll_no"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Type       string `json:"type"` // student or teacher
}

// User is the GORM model for users table
type User struct {
	ID         uint   `gorm:"primaryKey"`
	FullName   string
	DOB        string
	University string
	Semester   string
	Program    string
	RollNo     string
	Email      string
	Password   string
	Type       string
}

// LoginRequest represents the login payload
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterHandler handles user registration
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"message": "Invalid request payload",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"message": "Failed to hash password",
		})
		return
	}

	user := User{
		FullName:   req.FullName,
		DOB:        req.DOB,
		University: req.University,
		Semester:   req.Semester,
		Program:    req.Program,
		RollNo:     req.RollNo,
		Email:      req.Email,
		Password:   string(hashedPassword),
		Type:       req.Type,
	}
	if err := internal.DB.Create(&user).Error; err != nil {
		log.Printf("DB error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"message": "Failed to register user",
		})
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"message": "User registered successfully",
	})
}

// LoginHandler handles user login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"message": "Invalid request payload",
		})
		return
	}

	var user User
	// Find user by email (username)
	if err := internal.DB.Where("email = ?", req.Username).First(&user).Error; err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"message": "Invalid username or password",
		})
		return
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"message": "Invalid username or password",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"message": "Login successful",
	})
}
