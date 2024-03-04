package utils

import (
	"context"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"regexp"
	"strings"
	"time"
	"todo-grpc/models"
)

const AuthorizationKey = "authorization"

func ValidateJwtToken(md metadata.MD, config *EnvConfig) (*models.UserClaims, error) {
	authHeaders, ok := md[AuthorizationKey]
	if !ok || len(authHeaders) == 0 {
		return nil, grpc.Errorf(codes.Unauthenticated, "token not present in headers")
	}

	authToken := strings.Split(authHeaders[0], " ")
	if len(authToken) != 2 || authToken[0] != "Bearer" {
		return nil, grpc.Errorf(codes.Unauthenticated, "token not bearer")
	}

	token := authToken[1]
	claims := &models.UserClaims{}
	parseToken, err := jwt.ParseWithClaims(
		token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.jwtSecret), nil
		},
	)
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, grpc.Errorf(codes.Unauthenticated, "invalid token signature")
		}
		return nil, grpc.Errorf(codes.Unauthenticated, "token is expired")
	}
	if !parseToken.Valid || claims.UserID == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "invalid token")
	}
	return claims, nil
}

func GenerateToken(userID string, config *EnvConfig) (string, error) {
	// Create token
	jwtExpirationTime := time.Now().Add(time.Hour * 24).Unix()
	claims := &models.UserClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwtExpirationTime,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	// The signing string should be secret (a generated UUID works too)
	t, err := token.SignedString([]byte(config.jwtSecret))
	if err != nil {
		return "", err
	}

	return t, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetUserNameFromContext(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if authedUsernames := md.Get(AuthedUserIdHex); len(authedUsernames) > 0 {
			return authedUsernames[0]
		}
	}
	return ""
}

func validatePassword(password string) bool {
	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	// Check for at least one digit
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	// Check for at least one of the special characters @$#
	hasSpecialChar := regexp.MustCompile(`[@$#]`).MatchString(password)
	return hasUppercase && hasDigit && hasSpecialChar && len(password) >= 8
}
