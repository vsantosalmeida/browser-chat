# Browser Chat

Application to simulate a simple browser chat.

### Details
Two applications to handle the requirements:
- chat-api
- chat-bot

**chat-api** handles the http and websocket connections, and communicate with chat-bot.

**chat-bot** receive commands from chat-api process it and send back the result.

#### Technologies
- MySQL
- Rabbitmq
- Docker

### Project structure
```bash
    .
    ├── api
    │   ├── midleware #auth and CORS middlewares
    │   ├── rest
    │   │   ├── handler #request handlers
    │   │   └── presenter #dto
    │   └── websocket #server
    ├── build #docker container build 
    │   ├── chat-api 
    │   └── chat-bot
    ├── cmd #app startup
    │   ├── chatapi
    │   └── chatbot
    ├── config
    ├── docker-compose.yml
    ├── entity #database entities
    ├── infrastructure
    │   ├── broker
    │   └── repository
    ├── pkg #shared packages
    │   ├── auth #JWT auth
    │   └── stooq #API client
    ├── templates #HTML page
    └── usecase #business rules
        ├── chatbot
        ├── room
        └── user
```

### How to run:
- Create a `.env` file in the project root, use the `.env-example` to check what are the required env vars.
- Run (if you don't have `docker` and `docker-compose` installed, follow the installation guides: [Docker](https://docs.docker.com/engine/install/) and [docker-compose](https://docs.docker.com/compose/install/))
    ```
  docker-compose up
    ```
- Wait for **chat-api** and **chat-bot** log the startup
    ```
  chat-bot | {"fields":{},"level":"info","timestamp":"2023-09-05T16:17:09.679204706Z","message":"chatbot started"}
  chat-api | {"fields":{},"level":"info","timestamp":"2023-09-05T16:17:09.813343883Z","message":"server started"}
    ```
### Using the chat:
Before start using the chat in your browser, it's required to set some configs in the chat-api.
1. Create some chat rooms
   - Request(no payload required)
   ```
    POST localhost:8080/rooms
    ```
   - You can confirm the created rooms in :
   ```
    GET localhost:8080/rooms
    ```
2. Create users
    - Request
   ```
    POST localhost:8080/users
    {
      "username": "your-user",
      "password": "your-pass"
    }
    ```
    - You can confirm the created users in :
   ```
    GET localhost:8080/users
    ```
3. In your browser go to `localhost:3000` and start using the UI
    - Use your user credentials to login and start send and receive messages

### Running tests:
   ```
    make test
   ```

### Other endpoints
- Login
   ```
    POST localhost:8080/users/login
    {
      "username": "your-user",
      "password": "your-pass"
    }
    ```
- Chat room messages
   ```
    GET localhost:8080/rooms/{id}/messages
   ```
### Websocket
- Connect, you need to log in and use the returned token to connect
   ```
    ws://localhost:8080/ws?bearer={token}
   ```
- Join chat room
   ```
    {
      "action": "joinRoom",
      "payload": {
          "roomID": {id}
      }
    }
   ```
- Send message
   ```
  {
    "action": "sendMessage",
    "payload": {
        "message": "hello world",
        "from": "your-user",
      }
  }
   ```
- Chatbot command
   ```
  {
    "action": "chatbotCommand",
    "payload": {
      "roomID": 1,
        "from": "your-user",
        "commandName": "stock",
        "command": "amzn.us"
      }
  }
   ```