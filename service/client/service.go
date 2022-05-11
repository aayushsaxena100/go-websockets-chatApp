package client

import (
	"errors"
	"github.com/aayushsaxena100/testproject/models"
	"github.com/aayushsaxena100/testproject/store"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type service struct {
	store store.Client
}

func New(store store.Client) *service {
	return &service{
		store: store,
	}
}

func (s *service) Register(c *gin.Context, client *models.Client) error {
	if err := validateClientForRegister(c, client); err != nil {
		return err
	}

	existingClient, err := s.store.GetByUsername(c, client.Username)
	if err != nil {
		return err
	} else {
		if existingClient != nil {
			return errors.New("username is taken")
		}
	}

	client, err = s.store.Create(c, client)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) Signin(c *gin.Context, client *models.Client) error {
	if err := validateClientForSignin(c, client); err != nil {
		return err
	}

	existingClient, err := s.store.GetByUsername(c, client.Username)
	if err != nil {
		return err
	} else {
		if existingClient != nil {
			if strings.Compare(existingClient.Password, client.Password) == 0 {
				return nil
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, struct {
					Message string `json:"message"`
				}{
					Message: "invalid credentials",
				})
				return errors.New("invalid credentials")
			}
		}
	}

	return nil
}

func validateClientForRegister(c *gin.Context, client *models.Client) error {
	if client.Username == "" {
		_, _ = c.Writer.Write([]byte("Missing Username"))
		return errors.New("missing username")
	} else if client.Password == "" {
		_, _ = c.Writer.Write([]byte("Missing Password"))
		return errors.New("missing password")
	}

	return nil
}

func validateClientForSignin(c *gin.Context, client *models.Client) error {
	if client.Username == "" {
		_, _ = c.Writer.Write([]byte("Missing Username"))
		return errors.New("missing username")
	} else if client.Password == "" {
		_, _ = c.Writer.Write([]byte("Missing Password"))
		return errors.New("missing password")
	}

	return nil
}
