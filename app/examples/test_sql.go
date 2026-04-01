package main

import (
	"context"
	"log"
	"time"

	"challenge-besu/internal/database"
	"challenge-besu/internal/repository"
)

// Exemplo de teste manual da camada SQL
// Execute: cd app && go run examples/test_sql.go
func main() {
	ctx := context.Background()

	// 1. Conectar ao banco
	databaseURL := "postgres://admin:admin123@localhost:5433/challenge_besu?sslmode=disable"
	pool, err := database.NewPostgresPool(ctx, databaseURL)
	if err != nil {
		log.Fatalf("❌ Erro ao conectar: %v", err)
	}
	defer pool.Close()
	log.Println("✅ Conectado ao PostgreSQL")

	// 2. Criar repository SQL
	repo := repository.NewSQLRepository(pool)

	// 3. Testar SaveValue
	testValue := int64(42)
	if err := repo.SaveValue(ctx, testValue); err != nil {
		log.Fatalf("❌ Erro ao salvar: %v", err)
	}
	log.Printf("✅ Valor %d salvo com sucesso", testValue)

	// 4. Testar GetValue
	time.Sleep(100 * time.Millisecond) // Pequeno delay
	retrieved, err := repo.GetValue(ctx)
	if err != nil {
		log.Fatalf("❌ Erro ao buscar: %v", err)
	}
	log.Printf("✅ Valor recuperado: %d", retrieved)

	// 5. Validar
	if retrieved == testValue {
		log.Println("✅ TESTE PASSOU: valores são iguais")
	} else {
		log.Fatalf("❌ TESTE FALHOU: esperado %d, obtido %d", testValue, retrieved)
	}
}
