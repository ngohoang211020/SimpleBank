package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"simplebank/token"
	"strings"
)

const (
	authorizationHeaderKey   = "authorization"
	authorizationTypeBearer  = "bearer"
	UnauthorizedTokenMissing = "token is missing"
	UnauthorizedTokenInvalid = "the access token is not valid"
	UnsupportedTokenType     = "unsupported authorization type"
	authorizationPayloadKey  = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New(UnauthorizedTokenMissing)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)

		if len(fields) < 2 {
			err := errors.New(UnauthorizedTokenInvalid)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := errors.New(UnsupportedTokenType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]

		//Check token is valid or not, if valid return payload
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		//Store payload to the context before passing it to next handler
		ctx.Set(authorizationPayloadKey, payload)

		//Forward request to next Handler
		ctx.Next()
	}
}
