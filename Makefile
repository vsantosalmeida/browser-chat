BINARY_PATH=./bin
CHAT_API_BINARY_NAME=$(BINARY_PATH)/chat-api.bin
CHATBOT_BINARY_NAME=$(BINARY_PATH)/chat-bot.bin
VERSION=1.0.0


clean:
	@ rm -rf bin/*

build-chat-api:
	@ echo " ---       BUILDING CHAT API     --- "
	@ go build -ldflags "-s -w -X main.version=$(VERSION)" -o $(CHAT_API_BINARY_NAME) cmd/chatapi/main.go
	@ echo " ---      FINISH BUILD       --- "

build-chat-bot:
	@ echo " ---       BUILDING CHATBOT     --- "
	@ go build -ldflags "-s -w -X main.version=$(VERSION)" -o $(CHATBOT_BINARY_NAME) cmd/chatbot/main.go
	@ echo " ---      FINISH BUILD       --- "