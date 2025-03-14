package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	createDatabase()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, World!")
	})

	http.HandleFunc("/login", login)

	http.HandleFunc("/protected", protectedEndpoint)

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createDatabase() {
	// Abrir ou criar banco de dados SQLite
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Criar tabela de usuários se não existir
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Erro ao criar tabela:", err)
	}

	fmt.Println("Banco de dados e tabela criados com sucesso!")

	// Criar um usuário de teste
	email := "user@example.com"
	password := "senha123"

	// Gerar hash seguro da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Erro ao gerar hash da senha:", err)
	}

	// Inserir usuário no banco
	_, err = db.Exec("INSERT OR IGNORE INTO users (email, password) VALUES (?, ?)", email, string(hashedPassword))
	if err != nil {
		log.Fatal("Erro ao inserir usuário:", err)
	}


	fmt.Println("Usuário de teste criado com sucesso!")
}

func login(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		http.Error(w, "Erro ao conectar ao banco de dados", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var storedEmail, storedPasswordHash string
	err = db.QueryRow("SELECT email, password FROM users WHERE email = ?", data.Email).Scan(&storedEmail, &storedPasswordHash)
	if err == sql.ErrNoRows {
		http.Error(w, "Usuário não encontrado", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Erro ao buscar usuário", http.StatusInternalServerError)
		return
	}

	// Comparar senha usando bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(data.Password))
	if err != nil {
		http.Error(w, "Senha incorreta", http.StatusUnauthorized)
		return
	}

	// Gerar token JWT (precisa implementar)
	token := generateToken(data.Email)

	// Retornar token JWT
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Token{AccessToken: token})
}

func generateToken(email string) string {
	header := "{\"alg\":\"HS256\",\"typ\":\"JWT\"}"
	payload := fmt.Sprintf("{\"sub\":\"%s\"}", email)

	signedJWT := SIGJWT(header, payload, "123")

	jwt := fmt.Sprintf("%s.%s.%s", signedJWT.Header, signedJWT.Payload, signedJWT.Signature)
	response := Token{AccessToken: jwt}
	return response.AccessToken
}

func validateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inválido")
		}
		return []byte("123"), nil
	})
}

func protectedEndpoint(w http.ResponseWriter, r *http.Request) {
	tokenString := extractToken(r)
	if tokenString == "" {
		http.Error(w, "Token ausente", http.StatusUnauthorized)
		return
	}

	token, err := validateToken(tokenString)
	if err != nil || !token.Valid {
		http.Error(w, "Token inválido ou expirado", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Erro ao interpretar token", http.StatusUnauthorized)
		return
	}

	email := claims["sub"].(string)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Acesso autorizado",
		"email":   email,
	})
}

func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 && parts[0] == "Bearer" {
		return parts[1]
	}

	return ""
}
