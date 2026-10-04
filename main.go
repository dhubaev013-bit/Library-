package main

import (
	"errors"
	"fmt"
)

func GetPortFromConfig(config map[string]string) (string, error) {
	port, ok := config["PORT"]
	if !ok {
		return "", errors.New("ключ PORT отсутствует в конфигурации")
	}
	return port, nil
}

type Book struct {
	ID       int
	Title    string
	Author   string
	IsIssued bool
	ReaderID int
}

type Reader struct {
	ID     int
	Name   string
	Active bool
}

type Library struct {
	Books   []Book
	Readers []Reader
}

func (b *Book) IssueBook(reader *Reader) error {
	if b.IsIssued {
		return fmt.Errorf("книга '%s' уже выдана", b.Title)
	}
	if !reader.Active {
		return fmt.Errorf("читатель '%s' не активен", reader.Name)
	}
	b.IsIssued = true
	b.ReaderID = reader.ID
	return nil
}

func (b *Book) ReturnBook() error {
	if !b.IsIssued {
		return fmt.Errorf("книга '%s' и так в библиотеке", b.Title)
	}
	b.IsIssued = false
	b.ReaderID = 0
	return nil
}

func (lib *Library) FindBookByID(id int) (*Book, error) {
	for i := range lib.Books {
		if lib.Books[i].ID == id {
			return &lib.Books[i], nil
		}
	}
	return nil, fmt.Errorf("книга с ID %d не найдена", id)
}

func (lib *Library) FindReaderByID(id int) (*Reader, error) {
	for i := range lib.Readers {
		if lib.Readers[i].ID == id {
			return &lib.Readers[i], nil
		}
	}
	return nil, fmt.Errorf("читатель с ID %d не найден", id)
}

func (lib *Library) IssueBookToReader(bookID int, readerID int) error {
	book, err := lib.FindBookByID(bookID)
	if err != nil {
		return err
	}

	reader, err := lib.FindReaderByID(readerID)
	if err != nil {
		return err
	}

	return book.IssueBook(reader)
}

func (lib *Library) ReturnBook(bookID int) error {
	book, err := lib.FindBookByID(bookID)
	if err != nil {
		return err
	}

	return book.ReturnBook()
}

func main() {
	fmt.Println("---Тестируем чтение конфигурации---")

	configWithPort := map[string]string{"PORT": "8080", "HOST": "localhost"}
	configWithoutPort := map[string]string{"HOST": "localhost"}

	port, err := GetPortFromConfig(configWithPort)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println("PORT:", port)
	}

	port, err = GetPortFromConfig(configWithoutPort)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println("PORT:", port)
	}

	myLibrary := Library{
		Books: []Book{
			{ID: 1, Title: "Война и мир", Author: "Толстой"},
			{ID: 2, Title: "Преступление и наказание", Author: "Достоевский"},
		},
		Readers: []Reader{
			{ID: 1, Name: "Иван", Active: true},
			{ID: 2, Name: "Мария", Active: false},
		},
	}

	fmt.Println("\n---Тестируем выдачу книг---")

	err = myLibrary.IssueBookToReader(1, 1)
	if err != nil {
		fmt.Println("Ошибка выдачи:", err)
	} else {
		fmt.Println("Книга успешно выдана")
	}

	book, _ := myLibrary.FindBookByID(1)
	if book != nil {
		fmt.Println("Статус книги после выдачи:", book)
	}

	err = myLibrary.IssueBookToReader(99, 1)
	if err != nil {
		fmt.Println("Ожидаемая ошибка:", err)
	}

	fmt.Println("\n---Тестируем возврат книг---")

	err = myLibrary.ReturnBook(1)
	if err != nil {
		fmt.Println("Ошибка возврата:", err)
	} else {
		fmt.Println("Книга успешно возвращена")
	}

	err = myLibrary.ReturnBook(1)
	if err != nil {
		fmt.Println("Ожидаемая ошибка:", err)
	}

	err = myLibrary.ReturnBook(99)
	if err != nil {
		fmt.Println("Ожидаемая ошибка:", err)
	}
}
