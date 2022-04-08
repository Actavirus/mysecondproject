package main

import (
    "fmt"
    "net/http"
    "strconv"
    "html/template"
    "errors"
    "github.com/snippetbox/pkg/models"
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
        app.notFound(w) // Использование помощника notFound() // Страница не найдена.
        return
    }

    // Вызываем метода Get из модели Snipping для извлечения данных для
	// конкретной записи на основе её ID. Если подходящей записи не найдено,
	// то возвращается ответ 404 Not Found (Страница не найдена).
    s, err := app.snippets.Get(id)
    if err != nil {
        if errors.Is(err, models.ErrNoRecord){
            app.notFound(w)
        } else {
            app.serverError(w, err)
        }
        return
    }

    // Отображаем весь вывод на странице.
    fmt.Fprintf(w, "%v", s)
}

// Меняем сигнатуру обработчика createSnippet, чтобы он определялся как метод
// структуры *application.
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.Header().Set("Allow", http.MethodPost)
        app.clientError(w, http.StatusMethodNotAllowed)  // Используем помощник
        return
    }

    // Создаем несколько переменных, содержащих тестовые данные. Мы удалим их позже.
    title := "История про улитку"
    content := "Улитка выползла из раковины,\nвытянула рожки,\nи опять подобрала их."
    expires := "7"

    // Передаем данные в метод SnippetModel.Insert(), получая обратно
	// ID только что созданной записи в базу данных.
    id, err := app.snippets.Insert(title, content, expires)
    if err != nil {
        app.serverError(w, err)
        return
    }

    // Перенаправляем пользователя на соответствующую страницу заметки.
    http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)

    // w.Write([]byte("Создание новой заметки..."))
}