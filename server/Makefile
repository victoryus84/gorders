.PHONY: up down restart build clean help

# English: Include environment variables from .env
# Română: Include variabilele de mediu din .env
# Русский: Включает переменные окружения из .env
include .env
export

# English: Build and start the entire stack
# Română: Construiește și pornește întreaga stivă
# Русский: Построить и запустить весь стек
build:
	docker compose -f compose-prod.yml up --build -d

# English: Start the entire stack in detached mode
# Română: Pornește întreaga stivă în mod detached
# Русский: Запуск всего стека в фоновом режиме
up:
	docker compose -f compose-prod.yml up -d

# English: Stop and remove containers and orphans
# Română: Oprește și șterge containerele și rețelele orfane
# Русский: Остановка всего стека и удаление лишних объектов
down:
	docker compose -f compose-prod.yml down --remove-orphans

# English: Full restart (Down then Up)
# Română: Repornire completă (Oprește apoi Pornește)
# Русский: Полная перезагрузка (Down потом Up)
restart: 
	docker compose -f compose-prod.yml down --remove-orphans
	docker compose -f compose-prod.yml up -d

# English: Full cleanup. WARNING: Deletes volumes (data loss!)
# Română: Curățare completă. ATENȚIE: Șterge volumele (pierdere de date!)
# Русский: Полная очистка. ВНИМАНИЕ: Удаляет тома (потеря данных!)
clean:
	docker compose -f compose-prod.yml down --remove-orphans -v
	docker system prune -f

# English: Show help information
# Română: Afișează informații de ajutor
# Русский: Показать справку
help:
	@echo "Commands / Comenzi / Команды:"
	@echo "  make build   - Build and start stack / Construiește și pornește / Построить и запустить"
	@echo "  make up      - Start stack / Pornește stiva / Запуск стека"
	@echo "  make down    - Stop stack / Oprește stiva / Остановка стека"
	@echo "  make restart - Restart stack / Repornește stiva / Перезагрузка"
	@echo "  make clean   - Full cleanup (DATA LOSS!) / Curățare (PIERDERE DATE!) / Очистка (ПОТЕРЯ ДАННЫХ!)"
	@echo "  make help    - Show this help / Ajutor / Справка"