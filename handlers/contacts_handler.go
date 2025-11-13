package handlers

import (
	"admin-api/models"
	"encoding/json"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

type ContactsAPI struct {
	db *gorm.DB
}

func NewContactsAPI(db *gorm.DB) *ContactsAPI {
	return &ContactsAPI{
		db: db,
	}
}

// GetContacts godoc
// @Summary      Получить контактную информацию
// @Description  Возвращает запись с контактами
// @Tags         contacts
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        contact body models.Contacts true "Обновлённые данные"
// @Success      200  {array}  models.Contacts
// @Failure      400 {object} map[string]string "Неверный запрос"
// @Failure      500  {object}  map[string]string
// @Router       /contacts [get]

func (c *ContactsAPI) GetContacts(w http.ResponseWriter, r *http.Request) {
	var items []models.Contacts
	if err := c.db.Find(&items).Error; err != nil {
		http.Error(w, "Ошибка БД", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// UpdateContacts godoc
// @Summary      Обновить контактную информацию
// @Description  Создаёт или обновляет единственную запись с контактами
// @Tags         contacts
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param 		 contact body models.Contacts true "Обновлённые данные"
// @Success      200  {array}  models.Contacts
// @Failure      400 {object} map[string]string "Неверный запрос"
// @Failure      500 {object} map[string]string "Ошибка БД"
// @Router       /contacts [put]

func (c *ContactsAPI) UpdateContacts(w http.ResponseWriter, r *http.Request) {

	var updatedContacts models.Contacts
	if err := json.NewDecoder(r.Body).Decode(&updatedContacts); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	if updatedContacts.Address == "" || updatedContacts.Phone == "" || updatedContacts.Email == "" {
		http.Error(w, "Обязательные поля не могут быть пустыми", http.StatusBadRequest)
		return
	}

	var existing models.Contacts
	if err := c.db.First(&existing).Error; err != nil {
		newContact := models.Contacts{
			Address:           updatedContacts.Address,
			Phone:             updatedContacts.Phone,
			Email:             updatedContacts.Email,
			Website:           updatedContacts.Website,
			WorkSchedule:      updatedContacts.WorkSchedule,
			SocialMediaVK:     updatedContacts.SocialMediaVK,
			SocialMediaYa:     updatedContacts.SocialMediaYa,
			SocialMediaTwoGis: updatedContacts.SocialMediaTwoGis,
		}
		if err := c.db.Create(&newContact).Error; err != nil {
			http.Error(w, "Ошибка БД при создании", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(newContact)
		return
	}
	result := c.db.Model(&existing).Where("1=1").Updates(updatedContacts)
	if result.Error != nil {
		fmt.Println("c.db.Model.Updates: ", result.Error)
		http.Error(w, "Ошибка БД при обновлении", http.StatusInternalServerError)
		return
	}

	c.db.Take(&existing)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existing)
}
