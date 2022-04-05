package main

import (
    "fmt"
    "net/http"
    "strconv"
    "html/template"
)
// Обработчик главной страницы.
// Меняем сигнатуры обработчика home, чтобы он определялся как метод
// структуры *application.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        app.notFound(w) // Использование помощника notFound()
        return
    }
    // Инициализируем срез содержащий пути к двум файлам. Обратите внимание, что
	// файл home.page.tmpl должен быть *первым* файлом в срезе.
    files := []string{
        "./ui/html/home.page.tmpl",
        "./ui/html/base.layout.tmpl",
        "./ui/html/footer.partial.tmpl",
    }
    // Используем функцию template.ParseFiles() для чтения файла шаблона.
	// Если возникла ошибка, мы запишем детальное сообщение ошибки и
	// используя функцию http.Error() мы отправим пользователю
	// ответ: 500 Internal Server Error (Внутренняя ошибка на сервере)
    ts, err := template.ParseFiles(files...)
    if err != nil {
        // Поскольку обработчик home теперь является методом структуры application
		// он может получить доступ к логгерам из структуры.
		// Используем их вместо стандартного логгера от Go.
        app.serverError(w, err) // Использование помощника serverError()
        return
    }
    // Затем мы используем метод Execute() для записи содержимого
	// шаблона в тело HTTP ответа. Последний параметр в Execute() предоставляет
	// возможность отправки динамических данных в шаблон.
    err = ts.Execute(w, nil)
    if err != nil {
        // Обновляем код для использования логгера-ошибок
		// из структуры application.
        app.serverError(w, err) // Использование помощника serverError()
    }
}

// Меняем сигнатуру обработчика showSnippet, чтобы он был определен как метод
// структуры *application
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.URL.Query().Get("id"))
    if err != nil || id < 1 {
        app.notFound(w) // Использование помощника notFound()
        return
    }

    fmt.Fprintf(w, "Отображение определенной заметки с ID %d...", id)
}

// Меняем сигнатуру обработчика createSnippet, чтобы он определялся как метод
// структуры *application.
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.Header().Set("Allow", http.MethodPost)
        app.clientError(w, http.StatusMethodNotAllowed)  // Используем помощник
        return
    }

    w.Write([]byte("Создание новой заметки..."))
}