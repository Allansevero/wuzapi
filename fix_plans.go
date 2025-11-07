package main

import (
	"fmt"
	"log"
)

// FixPlansInDatabase atualiza os planos no banco de dados existente
func (s *server) FixPlansInDatabase() error {
	log.Println("Iniciando atualização dos planos...")

	// Verificar se a tabela plans existe
	var count int
	err := s.db.Get(&count, "SELECT COUNT(*) FROM plans")
	if err != nil {
		return fmt.Errorf("tabela plans não existe ou erro ao acessar: %w", err)
	}

	if count == 0 {
		log.Println("Nenhum plano encontrado no banco.")
		return nil
	}

	// Atualizar plano Gratuito (ID 1)
	_, err = s.db.Exec(`
		UPDATE plans 
		SET max_instances = 2, trial_days = 0, price = 0.00 
		WHERE id = 1`)
	if err != nil {
		return fmt.Errorf("erro ao atualizar plano Gratuito: %w", err)
	}
	log.Println("✓ Plano Gratuito atualizado: 2 instâncias, R$0,00")

	// Atualizar plano Pro (ID 2)
	_, err = s.db.Exec(`
		UPDATE plans 
		SET max_instances = 8, price = 47.00 
		WHERE id = 2`)
	if err != nil {
		return fmt.Errorf("erro ao atualizar plano Pro: %w", err)
	}
	log.Println("✓ Plano Pro atualizado: 8 instâncias, R$47,00")

	// Atualizar plano Analista (ID 3)
	_, err = s.db.Exec(`
		UPDATE plans 
		SET max_instances = 20, price = 97.00 
		WHERE id = 3`)
	if err != nil {
		return fmt.Errorf("erro ao atualizar plano Analista: %w", err)
	}
	log.Println("✓ Plano Analista atualizado: 20 instâncias, R$97,00")

	// Remover expiração de subscriptions no plano gratuito
	_, err = s.db.Exec(`
		UPDATE user_subscriptions 
		SET expires_at = NULL 
		WHERE plan_id = 1`)
	if err != nil {
		return fmt.Errorf("erro ao remover expiração do plano gratuito: %w", err)
	}
	log.Println("✓ Expiração removida das subscriptions no plano gratuito")

	// Mostrar planos atualizados
	var plans []Plan
	err = s.db.Select(&plans, "SELECT * FROM plans ORDER BY id")
	if err != nil {
		return fmt.Errorf("erro ao buscar planos atualizados: %w", err)
	}

	log.Println("\nPlanos após atualização:")
	for _, plan := range plans {
		log.Printf("  ID %d: %s - %d instâncias - R$%.2f\n", 
			plan.ID, plan.Name, plan.MaxInstances, plan.Price)
	}

	return nil
}
