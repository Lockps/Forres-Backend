package token













// package database

// import (
// 	"fmt"
// 	"time"

// 	"github.com/Lockps/Forres-release-version/cmd/function"
// 	"github.com/golang-jwt/jwt/v4"
// )

// type Auth struct {
// 	Issuer        string
// 	Audience      string
// 	Secret        string
// 	TokenExpiry   time.Duration
// 	RefreshExpiry time.Duration
// 	CookieDomain  string
// 	CookiePath    string
// 	CookieName    string
// }

// type jwtUser struct {
// 	ID    int    `json:"ID"`
// 	Fname string `json:"first_name"`
// 	Lname string `json:"last_name"`
// }

// type TokenPair struct {
// 	Token        string `json:"access_Token"`
// 	RefreshToken string `json:"refresh_token"`
// }

// type Claims struct {
// 	jwt.RegisteredClaims
// }

// func (j *Auth) GenerateTokenPair(user *jwtUser) (TokenPair, error) {
// 	//? Create A Token
// 	token := jwt.New(jwt.SigningMethodHS256)

// 	//? Set the Claim
// 	claim := token.Claims.(jwt.MapClaims)
// 	claim["name"] = fmt.Sprintf("%s %s", user.Fname, user.Lname)
// 	claim["sub"] = fmt.Sprint(user.ID)
// 	claim["aud"] = j.Audience
// 	claim["iss"] = j.Issuer
// 	claim["iat"] = time.Now().UTC().Unix()
// 	claim["typ"] = "JWT"

// 	//? Set the expiry for jwt
// 	claim["exp"] = time.Now().UTC().Add(j.TokenExpiry).Unix()

// 	//? Create a signed token
// 	signAccessToken, err := token.SignedString(function.StrToByteSlice(j.Secret))
// 	if err != nil {
// 		return TokenPair{}, err
// 	}

// 	//? Create refresh token and set cliams
// 	refreshToken := jwt.New(jwt.SigningMethodHS256)
// 	refreshTokenClaim := refreshToken.Claims.(jwt.MapClaims)
// 	refreshTokenClaim["sub"] = fmt.Sprint(user.ID)
// 	refreshTokenClaim["iat"] = time.Now().UTC().Unix()

// 	//? Set the expiry for the refresh Token
// 	refreshTokenClaim["exp"] = time.Now().UTC().Add(j.RefreshExpiry).Unix()

// 	//? Create signed refresh Token
// 	signrefreshToken, err := refreshToken.SignedString(function.StrToByteSlice(j.Secret))

// 	if err != nil {
// 		return TokenPair{}, err
// 	}
// 	//? Create TokenPair and populate with signed tokens
// 	var tokenPair = TokenPair{
// 		Token:        signAccessToken,
// 		RefreshToken: signrefreshToken,
// 	}

// 	//? Return TokenPair
// 	return tokenPair, nil
// }
