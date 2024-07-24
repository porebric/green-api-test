package server

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/porebric/green-api-test/internal/instance"
	instanceInMemory "github.com/porebric/green-api-test/internal/instance/inmemory"
	models2 "github.com/porebric/green-api-test/internal/instance/models"
	"github.com/porebric/green-api-test/internal/messages"
	messagesInMemory "github.com/porebric/green-api-test/internal/messages/inmemory"
	"github.com/porebric/green-api-test/internal/messages/models"
)

const (
	staticDir = "static"
	filesDir  = "files"
)

const (
	idInstanceName      = "idInstance"
	apiTokenName        = "apiTokenInstance"
	phoneNumberName     = "phoneNumber"
	messageBodyName     = "messageBody"
	phoneNumberFileName = "phoneNumberFile"
	fileUrlName         = "fileUrl"
)

type handler struct {
	messagesProvider messages.Provider
	instanceProvider instance.Provider

	dirPath string
}

func NewHandler() *handler {
	return &handler{
		dirPath:          "./assets",
		messagesProvider: messagesInMemory.NewProvider(), // заменить на реальный провайдер в котором будет коннекшн к базе данных
		instanceProvider: instanceInMemory.NewProvider(), // заменить на реальный провайдер в котором будет коннекшн к базе данных
	}
}

func (h *handler) Home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/%s/index.gohtml", h.dirPath, staticDir))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := mainTemplate{
		IDInstanceName:      idInstanceName,
		APITokenName:        apiTokenName,
		PhoneNumberName:     phoneNumberName,
		MessageBodyName:     messageBodyName,
		PhoneNumberFileName: phoneNumberFileName,
		FileUrlName:         fileUrlName,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *handler) GetSettingsHandler(w http.ResponseWriter, r *http.Request) {
	i, ok := h.checkInstance(r.Context(), w, r.URL.Query().Get(idInstanceName), r.URL.Query().Get(apiTokenName))
	if !ok {
		return
	}

	settings, _ := h.instanceProvider.GetSettings(r.Context(), i.GetInstanceId())

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(&settings); err != nil {
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}

func (h *handler) GetStateInstanceHandler(w http.ResponseWriter, r *http.Request) {
	i, ok := h.checkInstance(r.Context(), w, r.URL.Query().Get(idInstanceName), r.URL.Query().Get(apiTokenName))
	if !ok {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(&i); err != nil {
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}

func (h *handler) SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	i, ok := h.checkInstance(r.Context(), w, r.FormValue(idInstanceName), r.FormValue(apiTokenName))
	if !ok {
		return
	}

	m := models.Message{
		Id:         generateMessageID(),
		Phone:      r.FormValue(phoneNumberName),
		Body:       r.FormValue(messageBodyName),
		InstanceId: i.GetInstanceId(),
	}

	_ = h.messagesProvider.SaveMessage(r.Context(), m)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"idMessage": strconv.Itoa(int(m.Id))})
}

func (h *handler) SendFileByUrlHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	i, ok := h.checkInstance(r.Context(), w, r.FormValue(idInstanceName), r.FormValue(apiTokenName))
	if !ok {
		return
	}

	file, header, err := r.FormFile(fileUrlName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	if err = os.MkdirAll(fmt.Sprintf("%s/%s", h.dirPath, filesDir), os.ModePerm); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	messageID := generateMessageID()
	fileExt := filepath.Ext(header.Filename)
	newFileName := fmt.Sprintf("%d%s", messageID, fileExt)

	out, err := os.Create(filepath.Join(fmt.Sprintf("%s/%s", h.dirPath, filesDir), newFileName))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = out.ReadFrom(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	m := models.Message{
		Id:         generateMessageID(),
		Phone:      r.FormValue(phoneNumberFileName),
		Body:       newFileName,
		InstanceId: i.GetInstanceId(),
	}

	_ = h.messagesProvider.SaveMessage(r.Context(), m)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"idMessage": strconv.Itoa(int(messageID))})
}

func (h *handler) checkInstance(ctx context.Context, w http.ResponseWriter, idInstance, apiToken string) (models2.Instance, bool) {
	intInstanceId, err := strconv.ParseInt(idInstance, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return models2.Instance{}, false
	}

	i, err := h.instanceProvider.GetInstance(ctx, intInstanceId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return models2.Instance{}, false
	}

	if i.GetInstanceId() == 0 {
		http.Error(w, "instance not found", http.StatusNotFound)
		return models2.Instance{}, false
	}

	if i.Token != apiToken {
		http.Error(w, "invalid token", http.StatusForbidden)
		return models2.Instance{}, false
	}

	return i, true
}
