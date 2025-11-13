package docs

// This file provides minimal, no-op exported functions with swagger annotations
// so the swag generator picks up operations even when handlers are defined in
// another package. These functions are not used at runtime; they only exist
// to produce documentation.

import "net/http"

// GetServicesDocs godoc
// @Summary      Получить список услуг
// @Description  Возвращает все услуги из БД
// @Tags         services
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {array}  models.Services
// @Failure      400 {object} map[string]string "Неверный запрос"
// @Failure      500  {object}  map[string]string
// @Router       /services [get]
func GetServicesDocs(w http.ResponseWriter, r *http.Request) {}

// CreateServiceDocs godoc
// @Summary      Создать новую услугу (только для админа)
// @Tags         admin-services
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        service body models.Services true "Новая услуга"
// @Success      201 {object} models.Services
// @Failure      400 {object} map[string]string "Неверный запрос"
// @Failure      500 {object} map[string]string "Ошибка БД"
// @Router       /services [post]
func CreateServiceDocs(w http.ResponseWriter, r *http.Request) {}

// UpdateServiceDocs godoc
// @Summary      Обновить данные услуги
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
func UpdateServiceDocs(w http.ResponseWriter, r *http.Request) {}

// Gallery and Docs endpoints
// GetGalleryDocs godoc
// @Summary      Получить список изображений
// @Tags         gallery
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {array}  models.Gallery
// @Router       /gallery [get]
func GetGalleryDocs(w http.ResponseWriter, r *http.Request) {}

// UploadGalleryDocs godoc
// @Summary      Загрузить фото в галерею
// @Tags         gallery
// @Accept       mpfd
// @Produce      json
// @Param        files formData file true "Фото для загрузки"
// @Success      200 {array} models.Gallery
// @Router       /gallery [post]
func UploadGalleryDocs(w http.ResponseWriter, r *http.Request) {}

// UpdateGalleryDocs godoc
// @Summary      Обновить данные изображения
// @Tags         gallery
// @Accept       mpfd
// @Produce      json
// @Param        id   path int               true "ID документа"
// @Param        file formData file          true "Обновлённые данные"
// @Success      200 {object} models.Gallery
// @Router       /gallery/{id} [put]
func UpdateGalleryDocs(w http.ResponseWriter, r *http.Request) {}

// DeleteGalleryDocs godoc
// @Summary      Удалить изображение
// @Tags         gallery
// @Produce      json
// @Param        id   path int               true "ID изображения"
// @Success      204 "Успешно удалено"
// @Router       /gallery/{id} [delete]
func DeleteGalleryDocs(w http.ResponseWriter, r *http.Request) {}

// Contacts endpoints
// GetContactsDocs godoc
// @Summary      Получить контактную информацию
// @Tags         contacts
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {array}  models.Contacts
// @Router       /contacts [get]
func GetContactsDocs(w http.ResponseWriter, r *http.Request) {}

// UpdateContactsDocs godoc
// @Summary      Обновить контактную информацию
// @Tags         contacts
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        contact body models.Contacts true "Обновлённые данные"
// @Success      200 {array} models.Contacts
// @Router       /contacts [put]
func UpdateContactsDocs(w http.ResponseWriter, r *http.Request) {}

// Docs endpoints
// GetDocsDocs godoc
// @Summary      Получить список документов
// @Tags         docs
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {array}  models.Docs
// @Router       /docs [get]
func GetDocsDocs(w http.ResponseWriter, r *http.Request) {}

// UploadDocsDocs godoc
// @Summary      Создать новый документ
// @Tags         docs
// @Accept       mpfd
// @Produce      json
// @Param        files formData file true"Документ создан"
// @Success      200 {object} models.Docs
// @Router       /docs [post]
func UploadDocsDocs(w http.ResponseWriter, r *http.Request) {}

// UpdateDocsDocs godoc
// @Summary      Обновить данные документа
// @Tags         docs
// @Accept       mpfd
// @Produce      json
// @Param        id   path int               true "ID документа"
// @Param        file formData file         true "Новый файл"
// @Success      200 {object} models.Docs
// @Router       /docs/{id} [put]
func UpdateDocsDocs(w http.ResponseWriter, r *http.Request) {}

// DeleteDocsDocs godoc
// @Summary      Удалить документ
// @Tags         docs
// @Produce      json
// @Param        id   path int               true "ID документа"
// @Success      204 "Успешно удалено"
// @Router       /docs/{id} [delete]
func DeleteDocsDocs(w http.ResponseWriter, r *http.Request) {}
