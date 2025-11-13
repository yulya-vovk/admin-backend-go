package handlers

import (
	"admin-api/models"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type ServicesAPI struct {
	db *gorm.DB
}

func NewServicesAPI(db *gorm.DB) *ServicesAPI {
	return &ServicesAPI{
		db: db,
	}
}

// GetServices godoc
// @Summary      Получить список услуг
// @Description  Возвращает все услуги из БД
// @Tags         services
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {array}  models.Services
// @Failure      400 {object} map[string]string "Неверный запрос"
// @Failure      500  {object}  map[string]string
// @Router       /services [get]

func (s *ServicesAPI) GetServices(w http.ResponseWriter, r *http.Request) {
	var services []models.Services
	if err := s.db.Find(&services).Error; err != nil {
		http.Error(w, "Ошибка БД", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}

// CreateService godoc
// @Summary      Создать новую услугу (только для админа)
// @Description  Этот эндпоинт используется только в админке.
// @Tags         admin-services
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        service body models.Services true "Новая услуга"
// @Success      201 {object} models.Services
// @Failure      400 {object} map[string]string "Неверный запрос"
// @Failure      500 {object} map[string]string "Ошибка БД"
// @Router       /services [post]

func (s *ServicesAPI) CreateService(w http.ResponseWriter, r *http.Request) {
	var newService models.Services
	if err := json.NewDecoder(r.Body).Decode(&newService); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}
	if newService.Eng == "" || newService.Title == "" || newService.Src == "" || newService.Prices == "" || newService.Text == "" {
		http.Error(w, "Все поля обязательны", http.StatusBadRequest)
		return
	}
	if err := s.db.Create(&newService).Error; err != nil {
		http.Error(w, "Ошибка БД", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newService)
}

// UpdateService godoc
// @Summary      Обновить данные услуги
// @Description  Обновляет поля услуги по ID. Возвращает обновлённую запись.
// @Tags         services
// @Accept       json
// @Produce      json
// @Param        id   path int               true "ID услуги"
// @Param        service body models.Services true "Обновлённые данные"
// @Success      200 {object} models.Services
// @Failure      400 {object} map[string]string "Неверный ID или JSON"
// @Failure      404 {object} map[string]string "Услуга не найдена"
// @Failure      500 {object} map[string]string "Ошибка БД"
// @Router       /services/{id} [put]

func (s *ServicesAPI) UpdateService(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if !strings.HasPrefix(path, "/services/") {
		http.Error(w, "Неверный URL", http.StatusBadRequest)
		return
	}

	idStr := path[len("/services/"):]

	serviceId, err := strconv.Atoi(idStr)
	if err != nil || serviceId <= 0 {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	var updatedService models.Services
	if err := json.NewDecoder(r.Body).Decode(&updatedService); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	result := s.db.Model(&models.Services{}).Where("id = ?", serviceId).Updates(map[string]interface{}{
		"eng":    updatedService.Eng,
		"title":  updatedService.Title,
		"src":    updatedService.Src,
		"prices": updatedService.Prices,
		"text":   updatedService.Text,
	})
	if result.Error != nil {
		http.Error(w, "Ошибка БД", http.StatusInternalServerError)
		return
	}
	if result.RowsAffected == 0 {
		http.Error(w, "Услуга не найдена", http.StatusNotFound)
		return
	}

	var service models.Services
	s.db.Take(&service, serviceId)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(service)
}
