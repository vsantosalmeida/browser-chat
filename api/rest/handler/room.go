package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/vsantosalmeida/browser-chat/api/rest/presenter"
	"github.com/vsantosalmeida/browser-chat/usecase/room"

	"github.com/gorilla/mux"
)

type RoomHandler struct {
	useCase room.UseCase
}

func NewRoomHandler(useCase room.UseCase) *RoomHandler {
	return &RoomHandler{
		useCase: useCase,
	}
}

func (h *RoomHandler) HandleListRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.useCase.ListRooms()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	output := presenter.MapEntityToExternalRooms(rooms)

	b, err := json.Marshal(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

func (h *RoomHandler) HandleListMessages(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	value, ok := params["id"]
	if !ok {
		http.Error(w, "empty room id", http.StatusBadRequest)
		return
	}

	roomID, err := strconv.Atoi(value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mgs, err := h.useCase.ListMessages(roomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	output := presenter.MapEntityToExternalMessages(mgs)

	b, err := json.Marshal(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

func (h *RoomHandler) HandleCreateRoom(w http.ResponseWriter, r *http.Request) {
	id, err := h.useCase.CreateRoom()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	output := presenter.CreateRoomOutput{
		ID: id,
	}

	b, err := json.Marshal(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(b)
}
