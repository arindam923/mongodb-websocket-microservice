package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Message struct {
	SenderID   string    `json:"sender_id"`
	ReceiverID string    `json:"receiver_id"`
	Content    string    `json:"content"`
	Timestamp  time.Time `json:"timestamp"`
}

func getMongoDBCollection() *mongo.Collection {

	uri := os.Getenv("MONGODB_URI")

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}
	collection := client.Database("test").Collection("message")
	return collection
}

func saveToMongoDB(msg Message, collection *mongo.Collection) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, msg)
	return err
}

func handleWebSocket(conn *websocket.Conn, collection *mongo.Collection, receiverId string) {
	defer conn.Close()

	for {
		// Read message from the client
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
			return
		}

		// Print the received message
		fmt.Printf("Received message:\nSender ID: %s\nReceiver ID: %s\nContent: %s\nTimestamp: %s\n\n",
			msg.SenderID, msg.ReceiverID, msg.Content, msg.Timestamp.Format(time.RFC3339))

		err = saveToMongoDB(msg, collection)
		if err != nil {
			log.Println(err)
			return
		}

		if msg.ReceiverID == receiverId {
			response := Message{
				SenderID:   "Server",
				ReceiverID: msg.SenderID,
				Content:    "Received your message!",
				Timestamp:  time.Now(),
			}
			err = conn.WriteJSON(response)
			if err != nil {
				log.Println(err)
				return
			}

			// Print the sent response
			fmt.Printf("Sent response:\nSender ID: %s\nReceiver ID: %s\nContent: %s\nTimestamp: %s\n\n",
				response.SenderID, response.ReceiverID, response.Content, response.Timestamp.Format(time.RFC3339))
		}
	}
}

func getMessagesByReceiverId(w http.ResponseWriter, r *http.Request) {
	receiverID := r.URL.Query().Get("receiver_id")
	limit := r.URL.Query().Get("limit")

	limitInt := 100
	if limit != "" {
		fmt.Sscanf(limit, "%d", &limitInt)
	}

	collection := getMongoDBCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{"receiver_id", receiverID}}
	options := options.Find().SetLimit(int64(limitInt))

	curr, err := collection.Find(ctx, filter, options)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error fetching messages", http.StatusInternalServerError)
		return
	}

	var messages []Message
	if err = curr.All(context.TODO(), &messages); err != nil {
		panic(err)
	}

	for _, result := range messages {
		res, _ := bson.MarshalExtJSON(result, false, false)
		fmt.Println(string(res))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	receiverId := r.URL.Query().Get("receiver_id")
	if receiverId == "" {
		log.Println("Receiver ID is missing")
		conn.Close()
		return
	}

	handleWebSocket(conn, getMongoDBCollection(), receiverId)
}

func main() {
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/messages", getMessagesByReceiverId)

	port := 8080
	log.Printf("Server started on :%d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
