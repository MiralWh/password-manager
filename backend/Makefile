# Требуется Go 1.22+ и запущенный Docker Desktop
.PHONY: test coverage cover-html clean test-unit test-integration

# Все тесты
test:
	go test ./... -count=1 -race

# С покрытием
coverage:
	go test ./... -count=1 -coverprofile=coverage.out -covermode=atomic
	@go tool cover -func=coverage.out | tail -1

# HTML-отчёт покрытия
cover-html: coverage
	go tool cover -html=coverage.out -o coverage.html
	@echo "Открой backend/coverage.html в браузере"

# Только unit-тесты (без Docker)
test-unit:
	go test ./internal/crypto/... ./internal/auth/... ./internal/config/... ./internal/docker/... -count=1 -v

# Только integration (нужен Docker)
test-integration:
	go test ./internal/repository/... ./internal/service/... ./internal/handler/... -count=1 -v

clean:
	rm -f coverage.out coverage.html
