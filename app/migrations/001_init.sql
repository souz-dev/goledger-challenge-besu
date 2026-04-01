-- Tabela de armazenamento com registro único (id=1)
-- Padrão Singleton: garante apenas 1 linha na tabela

CREATE TABLE IF NOT EXISTS storage (
    id INTEGER PRIMARY KEY DEFAULT 1 CHECK (id = 1),
    value BIGINT NOT NULL DEFAULT 0,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Inserir valor inicial único
INSERT INTO storage (id, value) 
VALUES (1, 0)
ON CONFLICT (id) DO NOTHING;

-- Index para performance (opcional, mas recomendado)
CREATE INDEX IF NOT EXISTS idx_storage_updated_at ON storage(updated_at DESC);
