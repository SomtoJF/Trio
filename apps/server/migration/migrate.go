package main

import (
	"github.com/somtojf/trio/initializers"
	"github.com/somtojf/trio/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
}

func main() {
	db := initializers.DB

	// Create ENUM type for SenderType
	db.Exec(`
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'sender_type_enum') THEN
				CREATE TYPE sender_type_enum AS ENUM ('User', 'Agent');
			END IF;
		END
		$$;
	`)

	// AutoMigrate other models
	db.AutoMigrate(&models.User{}, &models.Chat{}, &models.Agent{}, &models.AgentMetadata{})

	// Manually create Message table with ENUM type
	db.Exec(`
		CREATE TABLE IF NOT EXISTS messages (
			id SERIAL PRIMARY KEY,
			created_at TIMESTAMP WITH TIME ZONE,
			updated_at TIMESTAMP WITH TIME ZONE,
			deleted_at TIMESTAMP WITH TIME ZONE,
			external_id UUID DEFAULT gen_random_uuid(),
			content TEXT,
			chat_id INTEGER,
			sender_type sender_type_enum,
			sender_id INTEGER
		)
	`)

	// Add indexes and constraints
	db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_messages_deleted_at ON messages(deleted_at);
		ALTER TABLE messages ADD CONSTRAINT fk_messages_chat FOREIGN KEY (chat_id) REFERENCES chats(id);
	`)

	// Manually create GeminiLogs table with ENUM type
	db.Exec(`
		CREATE TABLE IF NOT EXISTS gemini_logs (
			id SERIAL PRIMARY KEY,
			created_at TIMESTAMP WITH TIME ZONE,
			updated_at TIMESTAMP WITH TIME ZONE,
			deleted_at TIMESTAMP WITH TIME ZONE,
			external_id UUID DEFAULT gen_random_uuid(),
			prompt TEXT,
			input_tokens INTEGER,
			output_tokens INTEGER,
			total_tokens INTEGER,
			sender_type sender_type_enum,
			sender_id INTEGER
		)
	`)

	// Add index for deleted_at
	db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_gemini_logs_deleted_at ON gemini_logs(deleted_at);
	`)
}
