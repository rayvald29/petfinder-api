package handlers

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/rayvald29/petfinder-api/database"
	"github.com/rayvald29/petfinder-api/utils"
)

type RegisterInput struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	Username  string `json:"username" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func RegisterHandler(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if input.Email == "" || input.Password == "" || input.Username == "" || input.FirstName == "" || input.LastName == "" {
		c.JSON(400, gin.H{"error": "Todos los campos son requeridos"})
		return
	}

	hashedPassword, err := utils.HashPassword(input.Password)

	if err != nil {
		c.JSON(500, gin.H{"error": "Error al hashear la contraseña"})
		return
	}

	//INSERT
	_, err = database.DB.Exec(context.Background(), "INSERT INTO users (email, password_hash, role_id, username, first_name, last_name) VALUES ($1, $2, $3, $4, $5, $6)", input.Email, hashedPassword, "d12e4b44-049c-4fc0-9f35-5c8b4f45eea9", input.Username, input.FirstName, input.LastName)

	if err != nil {
		println(err.Error())
		c.JSON(500, gin.H{"error": "Error al registrar el usuario"})
		return
	}

	c.JSON(201, gin.H{"message": "Usuario registrado exitosamente"})
}

func LoginHandler(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		println(err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if input.Email == "" || input.Password == "" {
		c.JSON(400, gin.H{"error": "Email y contraseña son requeridos"})
		return
	}

	var userID string
	var email string
	var passwordHash string

	err := database.DB.QueryRow(
		context.Background(),
		"SELECT id, email, password_hash FROM users WHERE email = $1",
		input.Email,
	).Scan(&userID, &email, &passwordHash)

	if err != nil {
		c.JSON(401, gin.H{"error": "Credenciales inválidas"})
		return
	}

	validPassword := utils.CheckPasswordHash(input.Password, passwordHash)

	if !validPassword {
		c.JSON(401, gin.H{"error": "Credenciales inválidas"})
		return
	}

	token, err := utils.GenerateToken(userID, email)

	if err != nil {
		c.JSON(500, gin.H{"error": "Error al generar el token"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Login Successful",
		"token":   token,
	})
}
