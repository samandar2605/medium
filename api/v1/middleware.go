package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	emailPkg "github.com/post/pkg/email"
	"github.com/post/pkg/utils"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationPayloadKey = "authorization_payload"
)

func (h *handlerV1) AuthMiddleware(c *gin.Context) {
	accessToken := c.GetHeader(authorizationHeaderKey)

	if len(accessToken) == 0 {
		err := errors.New("authorization header is not provided")
		c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	payload, err := utils.VerifyToken(accessToken)


	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	c.Set(authorizationPayloadKey, payload)
	c.Next()
}

func (m *handlerV1) GetAuthPayload(ctx *gin.Context) (*utils.Payload, error) {
	i, exists := ctx.Get(authorizationPayloadKey)
	if !exists {
		return nil, errors.New("kalla qo'yding")
	}

	payload, ok := i.(*utils.Payload)

	if !ok {
		return nil, errors.New("unknown user")
	}
	return payload, nil
}

func (h *handlerV1) sendVerificationCode(key, email string) error {
	code, err := utils.GenerateRandomCode(6)
	if err != nil {
		return err
	}

	err = h.inMemory.SetWithTTL(key+email, code, 1)
	if err != nil {
		return err
	}

	err = emailPkg.SendEmail(h.cfg, &emailPkg.SendEmailRequest{
		To:      []string{email},
		Subject: "Verification email",
		Body: map[string]string{
			"code": code,
		},
		Type: emailPkg.VerificationEmail,
	})

	if err != nil {
		return err
	}

	return nil
}
