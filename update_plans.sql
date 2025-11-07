-- Script para atualizar os planos no banco de dados existente
-- Execute este script no seu banco de dados (PostgreSQL ou SQLite)

-- Atualizar limites dos planos
UPDATE plans SET max_instances = 2, trial_days = 0, price = 0.00 WHERE id = 1;
UPDATE plans SET max_instances = 8, price = 47.00 WHERE id = 2;
UPDATE plans SET max_instances = 20, price = 97.00 WHERE id = 3;

-- Remover expiração de subscriptions no plano gratuito
UPDATE user_subscriptions SET expires_at = NULL WHERE plan_id = 1;

-- Verificar as alterações
SELECT * FROM plans;
