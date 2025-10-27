package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4"
)

func main() {
	// Читаем строку подключения из переменной окружения DATABASE_URL.
	// Это единственно верный, пролетарский способ!
	connString := os.Getenv("DATABASE_URL")

	if connString == "" {
		log.Fatal("Переменная окружения DATABASE_URL не установлена! Саботаж!")
	}

	// Ждем, пока БД будет готова (простой, но эффективный цикл)
	var conn *pgx.Conn
	var err error
	for i := 0; i < 10; i++ {
		conn, err = pgx.Connect(context.Background(), connString)
		if err == nil {
			break
		}
		fmt.Printf("Не могу подключиться к базе, товарищ... Попытка %d/10. Жду 2 секунды.\n", i+1)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Так и не смог подключиться к базе: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	fmt.Println("Успешно подключились к базе данных!")

	NumUsers, err := strconv.Atoi(os.Getenv("NUM_USERS"))
	if err != nil {
		NumUsers = 50
	}
	NumOrganizations, err := strconv.Atoi(os.Getenv("NUM_ORGANIZATIONS"))
	if err != nil {
		NumOrganizations = 20
	}
	NumProjects, err := strconv.Atoi(os.Getenv("NUM_PROJECTS"))
	if err != nil {
		NumProjects = 40
	}
	fmt.Printf("Делаю: %d пользователей, %d организаций, %d проектов!\n", NumUsers, NumOrganizations, NumProjects)
	config := Config{
		NumUsers:         NumUsers,
		NumOrganizations: NumOrganizations,
		NumProjects:      NumProjects,
		NumTags:          len(tags),
	}

	generator := NewGenerator(conn, config)
	if err := generator.Run(); err != nil {
		log.Fatalf("План провалился! Контрреволюция! Ошибка: %v", err)
	}
}
