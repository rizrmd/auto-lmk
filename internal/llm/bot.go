package llm

import (
	"context"
	"fmt"
	"log/slog"
)

// Bot handles car sales conversations
type Bot struct {
	provider Provider
	// carRepo, leadRepo will be injected
}

// NewBot creates a new conversation bot
func NewBot(provider Provider) *Bot {
	return &Bot{
		provider: provider,
	}
}

// ProcessMessage processes incoming message and returns bot response
func (b *Bot) ProcessMessage(ctx context.Context, tenantID int, senderPhone, messageText string, isSales bool) (string, error) {
	slog.Info("processing message", "tenant_id", tenantID, "sender", senderPhone, "is_sales", isSales)

	// TODO: Implement conversation processing
	// 1. Load conversation history
	// 2. Build system prompt (different for sales vs customer)
	// 3. Add user message
	// 4. Call LLM with functions
	// 5. If function call -> execute and call LLM again
	// 6. Return response
	// 7. Store conversation

	systemPrompt := b.buildSystemPrompt(isSales)
	slog.Debug("system prompt", "prompt", systemPrompt)

	return "", fmt.Errorf("bot processing not yet implemented")
}

// buildSystemPrompt creates system prompt based on user type
func (b *Bot) buildSystemPrompt(isSales bool) string {
	if isSales {
		return `Anda adalah asisten penjualan mobil untuk tim sales internal.

Anda membantu sales mengakses informasi inventory mobil dengan cepat.

Kemampuan Anda:
- Mencari mobil berdasarkan kriteria (merek, harga, transmisi, dll)
- Memberikan detail lengkap mobil termasuk spesifikasi
- Membuat catatan lead dari customer

Gunakan bahasa Indonesia yang profesional dan efisien.`
	}

	return `Anda adalah asisten penjualan mobil yang ramah dan membantu customer mencari mobil yang sesuai.

Kemampuan Anda:
- Mencari mobil berdasarkan budget dan preferensi customer
- Memberikan informasi detail mobil termasuk harga dan spesifikasi
- Membantu customer membandingkan pilihan
- Mengatur jadwal untuk melihat mobil atau test drive

Gunakan bahasa Indonesia yang ramah, natural, dan helpful.
Pahami istilah automotive Indonesia seperti: matic (automatic), bensin (gasoline), OTR (On The Road price).`
}

// GetAvailableFunctions returns functions available to LLM
func (b *Bot) GetAvailableFunctions(isSales bool) []Function {
	functions := []Function{
		{
			Name:        "searchCars",
			Description: "Cari mobil yang tersedia berdasarkan kriteria",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"brand": map[string]interface{}{
						"type":        "string",
						"description": "Merek mobil (Toyota, Honda, Mitsubishi, dll)",
					},
					"model": map[string]interface{}{
						"type":        "string",
						"description": "Model mobil (Avanza, CR-V, Xpander, dll)",
					},
					"max_price": map[string]interface{}{
						"type":        "integer",
						"description": "Harga maksimal dalam Rupiah",
					},
					"transmission": map[string]interface{}{
						"type":        "string",
						"enum":        []string{"manual", "automatic"},
						"description": "Jenis transmisi (matic/automatic atau manual)",
					},
					"fuel_type": map[string]interface{}{
						"type":        "string",
						"enum":        []string{"bensin", "diesel", "hybrid", "electric"},
						"description": "Jenis bahan bakar",
					},
				},
			},
		},
		{
			Name:        "getCarDetails",
			Description: "Dapatkan detail lengkap sebuah mobil termasuk foto",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"car_id": map[string]interface{}{
						"type":        "integer",
						"description": "ID mobil",
					},
				},
				"required": []string{"car_id"},
			},
		},
	}

	if isSales {
		// Sales get additional function to create leads
		functions = append(functions, Function{
			Name:        "createLead",
			Description: "Buat lead baru dari customer yang tertarik",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"phone_number": map[string]interface{}{
						"type":        "string",
						"description": "Nomor telepon customer",
					},
					"name": map[string]interface{}{
						"type":        "string",
						"description": "Nama customer",
					},
					"interested_car_id": map[string]interface{}{
						"type":        "integer",
						"description": "ID mobil yang diminati",
					},
				},
				"required": []string{"phone_number"},
			},
		})
	}

	return functions
}
