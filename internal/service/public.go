package service

import (
	"context"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"template/internal/model/entities"
	"template/internal/repository"
	"template/pkg/auth"
	"template/pkg/config"
	"template/pkg/customerr"
	"template/pkg/database/cached"
	"template/pkg/log"
	"time"
)

type publicService struct {
	userRepo        repository.User
	timeoutDuration time.Duration
	session         cached.Session
	jwt             auth.JWTUtil
	logger          *log.Logs
}

func InitPublicService(
	userRepo repository.User,
	session cached.Session,
	jwt auth.JWTUtil,
	logger *log.Logs,
) Public {
	return publicService{
		userRepo:        userRepo,
		timeoutDuration: time.Duration(viper.GetInt(config.TimeOut)) * time.Millisecond,
		session:         session,
		jwt:             jwt,
		logger:          logger,
	}
}

func (p publicService) CreateUser(ctx context.Context, userCreate entities.UserCreate) (string, string, error) {
	var err error

	userID, err := p.userRepo.Create(ctx, userCreate)
	if err != nil {
		p.logger.Error(err.Error())
		return "", "", err
	}

	userToken := p.jwt.CreateToken(userID)

	userSessionID, err := p.session.Set(ctx, cached.SessionData{
		User: entities.User{
			UserBase: entities.UserBase{
				Name: userCreate.Name,
			},
			ID: userID,
		},
		LoginTimeStamp: time.Now(),
	})
	if err != nil {
		p.logger.Error(err.Error())
		return "", "", err
	}

	return userToken, userSessionID, nil
}

func (p publicService) LoginUser(ctx context.Context, userLogin entities.UserCreate) (string, string, error) {
	userID, hashedPwd, err := p.userRepo.GetHashedPassword(ctx, userLogin.Name.String)
	if err != nil {
		p.logger.Error(err.Error())
		return "", "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(userLogin.Password))
	if err != nil {
		p.logger.Error(err.Error())
		return "", "", err
	}

	userSessionUUID, err := p.session.GetUUID(ctx, strconv.Itoa(userID))
	if err != nil {
		p.logger.Error(err.Error())
		return "", "", err
	}
	if userSessionUUID != "" {
		userSessionUUID, err = p.session.UpdateKey(ctx, userSessionUUID, userID)
		if err != nil {
			p.logger.Error(err.Error())
			return "", "", err
		}
	} else {
		userSessionUUID, err = p.session.Set(ctx, cached.SessionData{
			User: entities.User{
				UserBase: entities.UserBase{
					Name: userLogin.Name,
				},
				ID: userID,
			},
			LoginTimeStamp: time.Now(),
		})
	}

	jwtToken := p.jwt.CreateToken(userID)

	return jwtToken, userSessionUUID, nil
}

func (p publicService) Refresh(ctx context.Context, sessionID string, span trace.Span) (string, string, error) {
	span.AddEvent(CallToRedis)
	userData, err := p.session.Get(ctx, sessionID)
	if err != nil {
		return "", "", err
	}
	if userData.ID == 0 {
		return "", "", customerr.UserNotFound
	}

	span.AddEvent(CallToRedis)
	newSessionID, err := p.session.UpdateKey(ctx, sessionID, userData.ID)
	if err != nil {
		return "", "", err
	}

	span.AddEvent(CreateToken)
	userToken := p.jwt.CreateToken(userData.ID)

	return userToken, newSessionID, nil
}
