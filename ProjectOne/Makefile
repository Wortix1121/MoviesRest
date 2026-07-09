#Запуск приложения
run:
	go run ./cmd/api/main.go

#Создание миграций
migrateDB:
	migrate create -ext sql -dir migrations -seq create_users_table
