package middleware

import (
	"auth/internal/controllers"
	"auth/internal/usecases"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (m *middleware) HandleErrors(c *gin.Context) {
	if len(c.Errors) > 0 {
		err := c.Errors.Last()

		// Common ////////////////////////////////////////////////////////////////////////
		if errors.Is(err, controllers.ErrDataBindError) {
			c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		}

		if errors.Is(err, controllers.ErrAuthRequired) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
			return
		}

		if errors.Is(err, usecases.ErrEntityAlreadyExists) {
			c.AbortWithStatusJSON(http.StatusConflict, err.Error())
			return
		}
		//////////////////////////////////////////////////////////////////////////////////

		// Validation ////////////////////////////////////////////////////////////////////
		if errors.Is(err, usecases.ErrInvalidEntity) {
			c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		}
		///////////////////////////////////////////////////////////////////////////////////

		// Auth ///////////////////////////////////////////////////////////////////////////
		if errors.Is(err, usecases.ErrWrongPassword) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
			return
		}
		if errors.Is(err, usecases.ErrEntityNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, usecases.ErrNotAValidAccessToken) || errors.Is(err, usecases.ErrAccessTokenExpired) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
			return
		}
		if errors.Is(err, usecases.ErrNotAValidRefreshToken) || errors.Is(err, usecases.ErrRefreshTokenExpired) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
			return
		}
		///////////////////////////////////////////////////////////////////////////////////

		fmt.Printf("Unexpected error: %s\n", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, "Internal Server Error")
	}
}
