#Запуск приложения
run:
	go run ./cmd/api/main.go

#Создание миграций
createmigrate:
	migrate create -ext sql -dir migrations -seq create_users_table

# Приминить миграцию к БД(postgres)
migrateToDB:
	migrate -path ./migrations -database "postgres://postgres:qwerty@localhost:5432/MovieDB?sslmode=disable" up