# MongoDB WebSocket Microservice

This microservice provides real-time messaging functionality using WebSocket and MongoDB for storage. It's a simple messaging system where clients can send and receive messages, and the messages are stored in a MongoDB database.

## Setup

1. **Clone the repository:**

   ```bash
   git clone https://github.com/arindam/mongodb-websocket-microservice.git
    ```
   
2. **Create a .env file in the project root and set your MongoDB URI:**
    ```
    MONGODB_URI=mongodb://localhost:27017
   ```
3. **Install Dependencies**
    ```bash
    go get -u github.com/joho/godotenv
    go get -u github.com/gorilla/websocket
    go get -u go.mongodb.org/mongo-driver/mongo
    ```


# API Endpoints
WebSocket Endpoint: ws://localhost:8080/ws?receiver_id=your_receiver_id

Connect to this endpoint to establish a WebSocket connection for real-time messaging.

Get Messages Endpoint: http://localhost:8080/messages?receiver_id=your_receiver_id&limit=10

Retrieve recent messages for a specific receiver with an optional limit.

# Roadmap
### Here's what we plan to implement in future updates:

- Caching: Implement caching mechanisms for better performance.
- Message Broker: Integrate a message broker for scalable and efficient message handling.
- Real-time Notifications: Add real-time notifications for new messages.
- User Authentication: Implement user authentication for secure messaging.
- Message Encryption: Enhance security by adding message encryption.
- Message Deletion: Allow users to delete their messages.
Stay tuned for updates! Feel free to contribute or provide feedback.

# Contributing
If you'd like to contribute, please fork the repository, create a new branch, make your changes, and submit a pull request. We welcome any enhancements or bug fixes!

License
This project is licensed under the MIT License.

Happy Coding! 🚀📬

