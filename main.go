package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var redisClient *redis.Client
var ctx = context.Background()

// ログイン情報
type User struct {
	ID       string `json:"login_id"`
	Password string `json:"password"`
}

// ログインIDとパスワード
type LoginRequest struct {
	LoginID  string `json:"loginID"`
	Password string `json:"password"`
}

// 登録Item エンティティ
type Item struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

func main() {
	// Redisの設定
	redisClient = redis.NewClient(&redis.Options{
		Addr: "go_api_server-redis-1:6379",
		DB:   0,
	})

	router := mux.NewRouter()

	// エンドポイントの設定
	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/items", createItemHandler).Methods("POST")
	router.HandleFunc("/items/{id}", getItemHandler).Methods("GET")
	router.HandleFunc("/items/{id}", updateItemHandler).Methods("PUT")
	router.HandleFunc("/items/{id}", deleteItemHandler).Methods("DELETE")

	log.Println("Server started on :8000")
	log.Fatal(http.ListenAndServe(":8000", router))
} //end of main

// ユーザー認証用のハンドラ
func userAuthHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println("Failed to decode user:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	authenticated, err := authenticateUser(user)
	if err != nil {
		log.Println("Authentication error:", err)
		http.Error(w, "Authentication error", http.StatusInternalServerError)
		return
	}

	if !authenticated {
		log.Println("User not authenticated:", user.ID)
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// 認証が成功した場合の処理
	log.Println("User authenticated:", user.ID)
}

// //テスト用のHello World表示用。
// func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("Hello World"))
// }

// ユーザーの認証(認証成功時:true, 認証失敗時:fail)
func authenticateUser(user User) (bool, error) {
	// データベースからユーザー情報を取得する
	storedPassword, err := redisClient.Get(user.ID).Result()
	// Redisからの情報取得時にエラーがある場合
	if err != nil {
		return false, err
	}
	// パスワードの検証
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password))
	if err != nil {
		// パスワードが一致しない場合
		return false, nil
	}
	// パスワードが一致した場合
	return true, nil
}

// HTTPリクエストを受け取り、ユーザーのログインを処理
func loginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received a login request")

	var loginRequest LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		log.Println("Invalid request body:", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	storedPassword, err := redisClient.Get(loginRequest.LoginID).Result()

	if err == redis.Nil {
		// ユーザーが存在しない場合、新しいユーザーを作成する
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(loginRequest.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Error hashing password:", err)
			http.Error(w, "error hashing password", http.StatusInternalServerError)
			return
		}

		err = redisClient.Set(loginRequest.LoginID, hashedPassword, 10*time.Minute).Err()
		if err != nil {
			log.Println("Error storing new user:", err)
			http.Error(w, "error storing new user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User created successfully. Please login."))
	} else if err != nil {
		log.Println("Error retrieving user:", err)
		http.Error(w, "error retrieving user", http.StatusInternalServerError)
		return
	} else {
		// パスワードを検証する
		err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(loginRequest.Password))
		if err != nil {
			log.Println("Login error:", err)
			http.Error(w, "login error", http.StatusUnauthorized)
		} else {
			w.Write([]byte("Login successful"))
		}
	}
}

// Create item
func createItemHandler(w http.ResponseWriter, r *http.Request) {
	var item Item
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		log.Println("Failed to decode item:", err)
		http.Error(w, "error creating item", http.StatusBadRequest)
		return
	}

	err = redisClient.Set(item.ID, item.Value, 0).Err()
	if err != nil {
		log.Println("Failed to store item:", err)
		http.Error(w, "error storing item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

// Read item
func getItemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	value, err := redisClient.Get(id).Result()
	if err == redis.Nil {
		http.Error(w, "item not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Println("Failed to retrieve item:", err)
		http.Error(w, "error retrieving item", http.StatusInternalServerError)
		return
	}

	item := Item{ID: id, Value: value}
	json.NewEncoder(w).Encode(item)
}

// Update item
func updateItemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var item Item
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		log.Println("Failed to decode the request body for item update:", err)
		http.Error(w, "failed to decode the request body for item update", http.StatusBadRequest)
		return
	}

	err = redisClient.Set(id, item.Value, 0).Err()
	if err != nil {
		log.Println("Failed to update item:", err)
		http.Error(w, "error updating item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(item)
}

// Delete item
func deleteItemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := redisClient.Del(id).Err()
	if err != nil {
		log.Println("Failed to delete item:", err)
		http.Error(w, "error deleting item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
