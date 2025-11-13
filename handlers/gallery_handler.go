package handlers

import (
	"admin-api/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type GalleryAPI struct {
	db *gorm.DB
}

func NewGalleryAPI(db *gorm.DB) *GalleryAPI {
	return &GalleryAPI{
		db: db,
	}
}

// GetGallery godoc
// @Summary      Получить список изображений
// @Description  Возвращает весь список изображений из БД
// @Tags         gallery
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        gallery body models.Gallery true "Обновлённые данные"
// @Success      200  {array}  models.Gallery
// @Failure      500  {object}  map[string]string
// @Router       /gallery [get]

func (g *GalleryAPI) GetGallery(w http.ResponseWriter, r *http.Request) {
	var items []models.Gallery
	if err := g.db.Find(&items).Error; err != nil {
		http.Error(w, "Ошибка БД", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// UpdateGallery godoc
// @Summary      Обновить данные изображения
// @Description  Обновляет поля изображения по ID. Возвращает обновлённую запись.
// @Tags         gallery
// @Accept       json
// @Produce      json
// @Param        id   path int               true "ID документа"
// @Param        gallery body models.Gallery true "Обновлённые данные"
// @Success      200 {object} models.Gallery
// @Failure      400 {object} map[string]string "Неверный ID или JSON"
// @Failure      400 {object} map[string]string "Неверный ID или JSON"
// @Failure      404 {object} map[string]string "Фото не найдено"
// @Failure      500 {object} map[string]string "Ошибка БД"
// @Router       /gallery/{id} [put]

func (g *GalleryAPI) UpdateGallery(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if !strings.HasPrefix(path, "/gallery/") {
		http.Error(w, "Неверный URL", http.StatusBadRequest)
		return
	}

	idStr := path[len("/gallery/"):]

	galleryId, err := strconv.Atoi(idStr)
	if err != nil || galleryId <= 0 {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	var updatedGallery models.Gallery
	if err := json.NewDecoder(r.Body).Decode(&updatedGallery); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	result := g.db.Model(&models.Gallery{}).
		Where("id = ?", galleryId).
		Updates(map[string]interface{}{
			"hidden": updatedGallery.Hidden,
		})
	if result.Error != nil {
		http.Error(w, "Ошибка БД", http.StatusInternalServerError)
		return
	}
	if result.RowsAffected == 0 {
		http.Error(w, "Картинка не найдена", http.StatusNotFound)
		return
	}

	var item models.Gallery
	g.db.Take(&item, galleryId)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)

}

// DeleteGallery godoc
// @Summary      Удалить изображение
// @Description  Удаляет изображение по ID
// @Tags         gallery
// @Accept       json
// @Produce      json
// @Param        id   path int               true "ID изображения"
// @Param        gallery body models.Gallery true "Обновлённые данные"
// @Success      204 "Успешно удалено"
// @Failure      400 {object} map[string]string "Неверный ID или JSON"
// @Failure      404 {object} map[string]string "Фото не найдено"
// @Failure      500 {object} map[string]string "Ошибка БД"
// @Router       /gallery/{id} [delete]

func (g *GalleryAPI) DeleteGallery(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if !strings.HasPrefix(path, "/gallery/") {
		http.Error(w, "Неверный URL", http.StatusBadRequest)
		return
	}

	idStr := path[len("/gallery/"):]

	galleryId, err := strconv.Atoi(idStr)
	if err != nil || galleryId <= 0 {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}
	result := g.db.Delete(&models.Gallery{}, galleryId)

	if result.Error != nil {
		http.Error(w, "Ошибка БД", http.StatusInternalServerError)
		return
	}
	if result.RowsAffected == 0 {
		http.Error(w, "Картинка не найдена", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UploadGalleryFiles godoc
// @Summary      Загрузить фото в галерею
// @Description  Принимает один или несколько файлов и сохраняет их
// @Tags         gallery
// @Accept       mpfd
// @Produce      json
// @Param        files formData file true "Фото для загрузки"
// @Success      200 {array} models.Gallery
// @Failure      400 {object} map[string]string "Неверный запрос"
// @Failure      500 {object} map[string]string "Ошибка БД или записи файла"
// @Router       /gallery [post]

func (g *GalleryAPI) UploadGalleryFiles(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Ошибка при парсинге запроса: "+err.Error(), http.StatusInternalServerError)
		return
	}

	files := r.MultipartForm.File["files"]

	var uploadedItems []models.Gallery

	for _, fileHeader := range files {

		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Не удалось открыть файл", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		ext := filepath.Ext(fileHeader.Filename)
		filename := fmt.Sprintf("upload_%d_%s%s", time.Now().UnixNano(), randomString(6), ext)
		dstPath := filepath.Join("uploads", filename)

		if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
			http.Error(w, "Не удалось создать папку uploads", http.StatusInternalServerError)
			return
		}

		dst, err := os.Create(dstPath)
		if err != nil {
			http.Error(w, "Не удалось создать файл", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "Ошибка записи файла", http.StatusInternalServerError)
			return
		}

		galleryItem := models.Gallery{
			Filename: "/uploads/" + filename,
			Hidden:   false,
		}
		if err := g.db.Create(&galleryItem).Error; err != nil {
			http.Error(w, "Ошибка БД", http.StatusInternalServerError)
			return
		}

		uploadedItems = append(uploadedItems, galleryItem)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(uploadedItems)
}
