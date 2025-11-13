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

type DocsAPI struct {
	db *gorm.DB
}

func NewDocsAPI(db *gorm.DB) *DocsAPI {
	return &DocsAPI{
		db: db,
	}
}

// GetDocs godoc
// @Summary      Получить список документов
// @Description  Возвращает весь список документов из БД
// @Tags         docs
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Success      200  {array}  models.Docs
// @Failure      500  {object}  map[string]string
// @Router       /docs [get]

func (d *DocsAPI) GetDocs(w http.ResponseWriter, r *http.Request) {
	var items []models.Docs
	if err := d.db.Find(&items).Error; err != nil {
		http.Error(w, "Ошибка БД", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// UpdateDocs godoc
// @Summary      Обновить данные документа
// @Description  Обновляет поля документа по ID. Возвращает обновлённую запись.
// @Tags         docs
// @Accept       mpfd
// @Produce      json
// @Param        id   path int               true "ID документа"
// @Param        file formData file         true "Новый файл"
// @Success      200 {object} models.Docs
// @Failure      400 {object} map[string]string "Неверный ID или JSON"
// @Failure      404 {object} map[string]string "Документ не найден"
// @Failure      500 {object} map[string]string "Ошибка БД"
// @Router       /docs/{id} [put]

func (d *DocsAPI) UpdateDocs(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if !strings.HasPrefix(path, "/docs/") {
		http.Error(w, "Неверный URL", http.StatusBadRequest)
		return
	}

	idStr := path[len("/docs/"):]
	docsId, err := strconv.Atoi(idStr)
	if err != nil || docsId <= 0 {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Ошибка парсинга формы", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["file"]
	if len(files) == 0 {
		http.Error(w, "Файл не найден в запросе", http.StatusBadRequest)
		return
	}
	if len(files) > 1 {
		http.Error(w, "Можно обновить только один файл за раз", http.StatusBadRequest)
		return
	}

	fileHeader := files[0]
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

	result := d.db.Model(&models.Docs{}).
		Where("id = ?", docsId).
		Updates(models.Docs{
			Name: fileHeader.Filename,
			File: "/uploads/" + filename,
		})
	if result.Error != nil {
		http.Error(w, "Ошибка БД", http.StatusInternalServerError)
		return
	}
	if result.RowsAffected == 0 {
		http.Error(w, "Документ не найден", http.StatusNotFound)
		return
	}

	var item models.Docs
	d.db.Take(&item, docsId)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)

}

// DeleteDocs godoc
// @Summary      Удалить документ
// @Description  Удаляет документ по ID
// @Tags         docs
// @Accept       json
// @Produce      json
// @Param        id   path int               true "ID документа"
// @Success      204 "Успешно удалено"
// @Failure      400 {object} map[string]string "Неверный ID или JSON"
// @Failure      404 {object} map[string]string "Документ не найден"
// @Failure      500 {object} map[string]string "Ошибка БД"
// @Router       /docs/{id} [delete]

func (d *DocsAPI) DeleteDocs(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if !strings.HasPrefix(path, "/docs/") {
		http.Error(w, "Неверный URL", http.StatusBadRequest)
		return
	}

	idStr := path[len("/docs/"):]

	docsId, err := strconv.Atoi(idStr)
	if err != nil || docsId <= 0 {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}
	result := d.db.Delete(&models.Docs{}, docsId)

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

// UploadDocsFiles godoc
// @Summary      Создать новый документ
// @Description  Создает новый документ
// @Tags         docs
// @Accept       mpfd
// @Produce      json
// @Param        files formData file true "Документ создан"
// @Success      200 {object} models.Docs
// @Failure      400 {object} map[string]string "Неверный ID или JSON"
// @Failure      500 {object} map[string]string "Ошибка БД"
// @Router       /docs [post]

func (d *DocsAPI) UploadDocsFiles(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(10 << 20)

	files := r.MultipartForm.File["files"]

	var uploadedItems []models.Docs

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

		docsItem := models.Docs{
			Name: fileHeader.Filename,
			File: "/uploads/" + filename,
		}
		if err := d.db.Create(&docsItem).Error; err != nil {
			http.Error(w, "Ошибка БД", http.StatusInternalServerError)
			return
		}

		uploadedItems = append(uploadedItems, docsItem)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(uploadedItems)
}
