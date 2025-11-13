// @title           Админка API
// @version         1.0
// @description     REST API для административной панели базы отдыха.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@localhost

// @license.name   MIT License
// @license.url    https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization
// @description                 Опционально: можно добавить JWT в будущем

package main

import (
	"admin-api/handlers"
	"admin-api/internal/db"
	"admin-api/models"
	"log"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "admin-api/docs" // важно: инициализация docs

	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	dsn := "host=localhost user=admin password=adminpass dbname=admin_api port=5432 sslmode=disable"
	dbConn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Не удалось подключиться к БД:", err)
	}
	db.Set(dbConn)

	err = dbConn.AutoMigrate(&models.Services{}, &models.Gallery{}, &models.Docs{}, &models.Contacts{})
	if err != nil {
		log.Fatal("Ошибка миграции:", err)
	}

	servicesAPI := handlers.NewServicesAPI(dbConn)
	galleryAPI := handlers.NewGalleryAPI(dbConn)
	contactsAPI := handlers.NewContactsAPI(dbConn)
	docsAPI := handlers.NewDocsAPI(dbConn)

	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	http.HandleFunc("GET /services", handlers.WithCORS(servicesAPI.GetServices))
	http.HandleFunc("POST /services", handlers.WithCORS(servicesAPI.CreateService))
	http.HandleFunc("PUT /services/{id}", handlers.WithCORS(servicesAPI.UpdateService))

	http.HandleFunc("GET /gallery", handlers.WithCORS(galleryAPI.GetGallery))
	http.HandleFunc("POST /gallery", handlers.WithCORS(galleryAPI.UploadGalleryFiles))
	http.HandleFunc("PUT /gallery/{id}", handlers.WithCORS(galleryAPI.UpdateGallery))
	http.HandleFunc("DELETE /gallery/{id}", handlers.WithCORS(galleryAPI.DeleteGallery))

	http.HandleFunc("GET /contacts", handlers.WithCORS(contactsAPI.GetContacts))
	http.HandleFunc("PUT /contacts", handlers.WithCORS(contactsAPI.UpdateContacts))

	http.HandleFunc("GET /docs", handlers.WithCORS(docsAPI.GetDocs))
	http.HandleFunc("POST /docs", handlers.WithCORS(docsAPI.UploadDocsFiles))
	http.HandleFunc("PUT /docs/{id}", handlers.WithCORS(docsAPI.UpdateDocs))
	http.HandleFunc("DELETE /docs/{id}", handlers.WithCORS(docsAPI.DeleteDocs))

	// Swagger UI
	http.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
