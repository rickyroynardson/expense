package auth

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/rickyroynardson/expense/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	config     *utils.Config
	repository *AuthRepository
}

func NewHandler(config *utils.Config, repository *AuthRepository) *AuthHandler {
	return &AuthHandler{
		config:     config,
		repository: repository,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var body RegisterRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		if errors.Is(err, io.EOF) {
			utils.RespondJSON(c, http.StatusBadRequest, "Invalid request body", nil)
			return
		}
		utils.RespondJSON(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.RespondJSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	body.Password = string(hashedPassword)

	if err := h.repository.Register(c.Request.Context(), body); err != nil {
		utils.RespondJSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.RespondJSON(c, http.StatusCreated, "Register success", nil)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var body LoginRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		if errors.Is(err, io.EOF) {
			utils.RespondJSON(c, http.StatusBadRequest, "Invalid request body", nil)
			return
		}
		utils.RespondJSON(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	user, err := h.repository.GetUserByEmail(c.Request.Context(), body.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			utils.RespondJSON(c, http.StatusUnauthorized, "Invalid email or password", nil)
			return
		}
		utils.RespondJSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(body.Password)); err != nil {
		utils.RespondJSON(c, http.StatusUnauthorized, "Invalid email or password", nil)
		return
	}

	accessToken, err := utils.GenerateJWT(user.ID, h.config.JwtSecret)
	if err != nil {
		utils.RespondJSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	refreshToken, err := utils.GenerateRefresh()
	if err != nil {
		utils.RespondJSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	err = h.repository.CreateRefreshToken(c.Request.Context(), user.ID, refreshToken)
	if err != nil {
		utils.RespondJSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	c.SetCookie("access_token", accessToken, 15*60, "/", h.config.CookieDomain, true, true)
	c.SetCookie("refresh_token", refreshToken, 30*24*60*60, "/", h.config.CookieDomain, true, true)

	utils.RespondJSON(c, http.StatusOK, "Login success", gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var token string
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		authorizationToken, _ := utils.GetAuthorizationToken(c.Request.Header)
		token = authorizationToken
	} else {
		token = cookie
	}

	if token == "" {
		utils.RespondJSON(c, http.StatusUnauthorized, "invalid token", nil)
		return
	}

	err = h.repository.DeleteRefreshToken(c.Request.Context(), token)
	if err != nil {
		utils.RespondJSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	c.SetCookie("access_token", "", -1, "/", h.config.CookieDomain, true, true)
	c.SetCookie("refresh_token", "", -1, "/", h.config.CookieDomain, true, true)

	utils.RespondJSON(c, http.StatusOK, "Logout success", nil)
}
