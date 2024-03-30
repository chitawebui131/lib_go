package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/go-sql-driver/mysql"

	//	"github.com/shopspring/decimal"
	"github.com/chitawebui131/lib_go/categories"
	"github.com/chitawebui131/lib_go/user"
)



// Book представляє модель продукту
type Books struct {
	ID            int       `json:"id"`
	Inventar      string    `json:"booksname"`
	Booksname     string    `json:"booksname"`
	Author        string    `json:"author"`
	Year          int       `json:"year"`
	City          string    `json:"city"`
	Publisher     string    `json:"publisher"`
	Department    int       `json:"department"`
	Count_page    int       `json:"count_page"`
	Bbk           string    `json:"bbk"`
	Count         int       `json:"count"`
	Comment       string    `json:"comment"`
	Created_at    time.Time `json:"created_at"`
	Updated_at    time.Time `json:"updated_at"`
}


type Category struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type BookWithCategoryWithoutDates struct {
	BookID           int     `json:"book_id"`
	BookName         string  `json:"book_name"`
	BookDescription  string  `json:"book_description"`
	BookPrice        float64 `json:"book_price"`
	StockQuantity       int     `json:"book_stockQuantity"`
	BookCategoryID   int     `json:"book_category_id"`
	CategoryID          int     `json:"category_id"`
	CategoryName        string  `json:"category_name"`
	CategoryDescription string  `json:"category_description"`
}

// BookService надає методи для роботи з продуктами
type BookService struct {
	DB *sql.DB
}

// GetBooks повертає список усіх продуктів з пагінацією
//GET /api/books?page=1&limit=10
/*
func (s *BookService) GetBooks(w http.ResponseWriter, r *http.Request) {
	// Отримання значень параметрів пагінації
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10 // За замовчуванням 10 елементів на сторінці
	}

	// Розрахунок зсуву (offset) для пагінації
	offset := (page - 1) * limit

	// Вибірка продуктів з бази даних з пагінацією
	rows, err := s.DB.Query("SELECT * FROM books JOIN categories ON books.category_id = categories.id LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		log.Println("Error querying database:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Створення слайсу для зберігання результатів
	var books []Book

	// Зчитування результатів запиту
	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ID, &book.Name, &book.Description, &book.Price, &book.StockQuantity, &book.CategoryID, &book.Created_at, &book.Updated_at ); err != nil {
			log.Println("Error scanning row:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		books = append(books, book)
	}

	// Перевірка наявності помилок під час зчитування
	if err := rows.Err(); err != nil {
		log.Println("Error reading rows:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Відправлення відповіді у форматі JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Кодуємо та виводимо дані у відповідь
	if err := json.NewEncoder(w).Encode(books); err != nil {
		log.Println("Error encoding JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *BookService) GetBooks(w http.ResponseWriter, r *http.Request) {
	// Отримання значень параметрів пагінації
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10 // За замовчуванням 10 елементів на сторінці
	}

	// Розрахунок зсуву (offset) для пагінації
	offset := (page - 1) * limit

	// Вибірка продуктів з бази даних з пагінацією
	query := `
		SELECT books.id AS book_id, books.name AS book_name,
			   books.description AS book_description, books.price AS book_price,
			   books.stock_quantity AS book_stockQuantity, books.category_id AS book_category_id,
			   books.created_at AS created_at, books.updated_at AS updated_at,
			   categories.id AS category_id, categories.name AS category_name,
			   categories.description AS category_description,
			   categories.created_at AS category_created_at, categories.updated_at AS category_updated_at
		FROM books
		JOIN categories ON books.category_id = categories.id
		LIMIT ? OFFSET ?
	`

	rows, err := s.DB.Query(query, limit, offset)
	if err != nil {
		log.Println("Error querying database:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Створення слайсу для зберігання результатів
	var booksWithCategories []BookWithCategory

	// Зчитування результатів запиту
	for rows.Next() {
		var bookWithCategory BookWithCategory
		// Сканування результатів у структуру продукту з категорією
		if err := rows.Scan(
			&bookWithCategory.BookID,
			&bookWithCategory.BookName,
			&bookWithCategory.BookDescription,
			&bookWithCategory.BookPrice,
			&bookWithCategory.StockQuantity,
			&bookWithCategory.BookCategoryID,
			&bookWithCategory.CreatedAt,
			&bookWithCategory.UpdatedAt,
			&bookWithCategory.CategoryID,
			&bookWithCategory.CategoryName,
			&bookWithCategory.CategoryDescription,
			&bookWithCategory.CategoryCreatedAt,
			&bookWithCategory.CategoryUpdatedAt,
		); err != nil {
			log.Println("Error scanning row:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		booksWithCategories = append(booksWithCategories, bookWithCategory)
	}

	// Перевірка наявності помилок під час зчитування
	if err := rows.Err(); err != nil {
		log.Println("Error reading rows:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Відправлення відповіді у форматі JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Кодуємо та виводимо дані у відповідь
	if err := json.NewEncoder(w).Encode(booksWithCategories); err != nil {
		log.Println("Error encoding JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
*/

func (s *BookService) GetBooks(w http.ResponseWriter, r *http.Request) {
	// ... (зберігаємо код пагінації та запиту з бази даних)
	// Отримання значень параметрів пагінації
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10 // За замовчуванням 10 елементів на сторінці
	}

	// Розрахунок зсуву (offset) для пагінації
	offset := (page - 1) * limit
	// Вибірка продуктів з бази даних з пагінацією
	query := `
		SELECT books.id AS book_id, books.name AS book_name, 
			   books.description AS book_description, books.price AS book_price, 
			   books.stock_quantity AS book_stockQuantity, books.category_id AS book_category_id,
			   categories.id AS category_id, categories.name AS category_name,
			   categories.description AS category_description
		FROM books
		LEFT JOIN categories ON books.category_id = categories.id
		LIMIT ? OFFSET ?
	`

	rows, err := s.DB.Query(query, limit, offset)
	if err != nil {
		log.Println("Error querying database:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Створення слайсу для зберігання результатів
	var booksWithCategoriesWithoutDates []BookWithCategoryWithoutDates

	// Зчитування результатів запиту
	for rows.Next() {
		var bookWithCategoryWithoutDates BookWithCategoryWithoutDates
		// Сканування результатів у структуру продукту з категорією без дат
		if err := rows.Scan(
			&bookWithCategoryWithoutDates.BookID,
			&bookWithCategoryWithoutDates.BookName,
			&bookWithCategoryWithoutDates.BookDescription,
			&bookWithCategoryWithoutDates.BookPrice,
			&bookWithCategoryWithoutDates.StockQuantity,
			&bookWithCategoryWithoutDates.BookCategoryID,
			&bookWithCategoryWithoutDates.CategoryID,
			&bookWithCategoryWithoutDates.CategoryName,
			&bookWithCategoryWithoutDates.CategoryDescription,
		); err != nil {
			log.Println("Error scanning row:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		booksWithCategoriesWithoutDates = append(booksWithCategoriesWithoutDates, bookWithCategoryWithoutDates)
	}

	// Перевірка наявності помилок під час зчитування
	if err := rows.Err(); err != nil {
		log.Println("Error reading rows:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Відправлення відповіді у форматі JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Кодуємо та виводимо дані у відповідь
	if err := json.NewEncoder(w).Encode(booksWithCategoriesWithoutDates); err != nil {
		log.Println("Error encoding JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// GetBook повертає інформацію про конкретний продукт за ID
func (s *BookService) GetBook(w http.ResponseWriter, r *http.Request) {
	// Отримання ID продукту з URL-параметра
	bookID := chi.URLParam(r, "id")
	if bookID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Вибірка конкретного продукту з бази даних за ID
	row := s.DB.QueryRow("SELECT * FROM books WHERE id=?", bookID)

	// Створення змінної для зберігання результатів
	var book Book
	//fmt.Println(row)
	// Зчитування результатів запиту
	err := row.Scan(&book.ID, &book.Name, &book.Description, &book.Price, &book.StockQuantity, &book.CategoryID, &book.Created_at, &book.Updated_at)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
		} else {
			log.Println("Error scanning row:", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// Відправлення відповіді у форматі JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Кодуємо та виводимо дані у відповідь
	if err := json.NewEncoder(w).Encode(book); err != nil {
		log.Println("Error encoding JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// CreateBook додає новий продукт
func (s *BookService) CreateBook(w http.ResponseWriter, r *http.Request) {
	// Отримання даних про новий продукт з тіла запиту (JSON)
	var newBook Book
	if err := json.NewDecoder(r.Body).Decode(&newBook); err != nil {
		log.Println("Error decoding JSON:", err)
		w.WriteHeader(http.StatusBadRequest)
		//fmt.Println(&newBook)
		return
	}
	fmt.Println(newBook)

	// Логіка додавання нового продукту до бази даних
	// result, err := s.DB.Exec("INSERT INTO books (name, description, price, stock_quantity, category_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
	// 	newBook.Name, newBook.Description, newBook.Price, newBook.StockQuantity, newBook.CategoryID, time.Now(), time.Now())
	// if err != nil {
	// 	log.Println("Error inserting book into database:", err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	query := `
    INSERT INTO books (name, description, price, stock_quantity, category_id, created_at, updated_at)
    VALUES (?, ?, ?, ?, ?, ?, ?)
`
	result, err := s.DB.Exec(query, newBook.Name, newBook.Description, newBook.Price, newBook.StockQuantity, newBook.CategoryID, time.Now(), time.Now())
	if err != nil {
		log.Println("Error inserting into database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Отримання ID новоствореного продукту
	newBookID, err := result.LastInsertId()
	if err != nil {
		log.Println("Error getting last insert ID:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	newBook.ID = int(newBookID)

	// Відправлення відповіді у форматі JSON з новоствореним продуктом та статусом 201 (Created)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newBook); err != nil {
		log.Println("Error encoding JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// UpdateBook оновлює інформацію про продукт за ID
func (s *BookService) UpdateBook(w http.ResponseWriter, r *http.Request) {
	// Отримання ID продукту з URL-параметра
	bookID := chi.URLParam(r, "id")
	if bookID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Отримання нових даних про продукт з тіла запиту (JSON)
	var updatedBook Book
	if err := json.NewDecoder(r.Body).Decode(&updatedBook); err != nil {
		log.Println("Error decoding JSON:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Логіка оновлення інформації про продукт в базі даних за ID
	// result, err := s.DB.Exec("UPDATE books SET name=?, description=?, price=?, stock_quantity=?, category_id=?, updated_at=?   WHERE id=?",
	// 	updatedBook.Name, updatedBook.Description, updatedBook.Price, updatedBook.StockQuantity, updatedBook.CategoryID, bookID, time.Now())
	query := `
		UPDATE books
		SET
			name = ?,
			description = ?,
			price = ?,
			stock_quantity = ?,
			category_id = ?,
			updated_at = ?
		WHERE id = ?
	`
	result, err := s.DB.Exec(query,
		updatedBook.Name,
		updatedBook.Description,
		updatedBook.Price,
		updatedBook.StockQuantity,
		updatedBook.CategoryID,
		time.Now(),
		bookID,
	)

	if err != nil {
		log.Println("Error updating book in database:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Перевірка, чи існує продукт за вказаним ID
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error checking rows affected:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		// Якщо немає відповідного продукту, відправити HTTP статус 404 (Not Found)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Відправлення відповіді у форматі JSON з оновленим продуктом
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updatedBook); err != nil {
		log.Println("Error encoding JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// DeleteBook видаляє продукт за ID
func (s *BookService) DeleteBook(w http.ResponseWriter, r *http.Request) {
	// Отримання ID продукту з URL-параметра
	bookID := chi.URLParam(r, "id")
	if bookID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Логіка видалення продукту з бази даних за ID
	result, err := s.DB.Exec("DELETE FROM books WHERE id=?", bookID)
	if err != nil {
		log.Println("Error deleting book from database:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Перевірка, чи існує продукт за вказаним ID
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error checking rows affected:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		// Якщо немає відповідного продукту, відправити HTTP статус 404 (Not Found)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Відправлення відповіді з підтвердженням видалення та статусом 204 (No Content)
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	// Ініціалізація роутера
	r := chi.NewRouter()

	// Ініціалізація сервісу продуктів з підключенням до бази даних
	db, err := sql.Open("mysql", "root:usbw@tcp(localhost:3306)/tklib?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	booktService := &BookService{DB: db}
	userSvc := &user.UserService{DB: db}
	catSvc := &categories.CatSetvices{DB: db}

	// Додавання middleware для логування запитів
	r.Use(middleware.Logger)

	// Додавання роутів
	r.Get("/books", bookService.GetBooks)
	r.Get("/books/{id}", bookService.GetBook)
	r.Post("/books", bookService.CreateBook)
	r.Put("/books/{id}", bookService.UpdateBook)
	r.Delete("/books/{id}", bookService.DeleteBook)
	r.Route("/users", func(r chi.Router) {
		r.Get("/", userSvc.GetUsers)
		r.Get("/{id}", userSvc.GetUser)
		r.Post("/", userSvc.CreateUser)
		r.Put("/{id}", userSvc.UpdateUser)
		r.Delete("/{id}", userSvc.DeleteUser)
	})
	r.Route("/cat", func(r chi.Router) {
		r.Get("/", catSvc.GetCats)
		r.Get("/{id}", catSvc.GetCat)
		r.Post("/", catSvc.CreateCat)
		r.Put("/{id}", catSvc.UpdateCat)
		r.Delete("/{id}", catSvc.DeleteCat)
	})

	// Запуск сервера на порту 8080
	port := getPort()
	fmt.Printf("Server is running on :%s...\n", port)
	http.ListenAndServe(":"+port, r)
}

// getPort повертає номер порту для веб-сервера
func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "7000" // За замовчуванням використовуємо 8080
	}
	return port
}
