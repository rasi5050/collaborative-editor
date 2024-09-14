package main

import (
    "encoding/json"
    "net/http"
    "time"

    "golang.org/x/crypto/bcrypt"
    "github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("my_secret_key")

type Credentials struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type Claims struct {
    Username string `json:"username"`
    jwt.StandardClaims
}

// Mock user database (In a real app, you would query a real database)
var users = map[string]string{
    "user1": "$2a$12$Dow5qf.q34/7fN8Z6jU/bOaIgzHpBxAG2Jihp8vN6Dy/LQ5zAx76S", // "password"
}

// Hash a password
func hashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

// Check the hash of a password
func checkPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
    var creds Credentials
    err := json.NewDecoder(r.Body).Decode(&creds)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    storedHash, ok := users[creds.Username]
    if !ok || !checkPasswordHash(creds.Password, storedHash) {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

    expirationTime := time.Now().Add(5 * time.Minute)
    claims := &Claims{
        Username: creds.Username,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    http.SetCookie(w, &http.Cookie{
        Name:    "token",
        Value:   tokenString,
        Expires: expirationTime,
    })

    json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func authenticate(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        cookie, err := r.Cookie("token")
        if err != nil {
            if err == http.ErrNoCookie {
                w.WriteHeader(http.StatusUnauthorized)
                return
            }
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        tokenStr := cookie.Value
        claims := &Claims{}

        token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })

        if err != nil {
            if err == jwt.ErrSignatureInvalid {
                w.WriteHeader(http.StatusUnauthorized)
                return
            }
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        if !token.Valid {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }

        next(w, r)
    }
}
