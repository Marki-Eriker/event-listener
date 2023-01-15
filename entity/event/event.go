package event

import (
	"github.com/marki-eriker/event-listener/pkg/encoder"
	uuid "github.com/satori/go.uuid"
	"time"
)

// IncomingPayload - структура событий приходящих от агентов
type IncomingPayload struct {
	EventID    int       `json:"event_id"`
	Created    time.Time `json:"created"`
	SystemName string    `json:"system_name"`
	Message    string    `json:"message"`
}

// Event - основная структура получаемых событий
type Event struct {
	tableName struct{} `pg:"events"`

	ID         uuid.UUID `json:"id" pg:"id,notnull,pk,type:uuid"`
	EventID    int       `json:"event_id" pg:"event_id"`
	Created    time.Time `json:"created" pg:"created"`
	SystemName string    `json:"system_name" pg:"system_name"`
	Message    string    `json:"message,omitempty" pg:"message"`
	Incident   bool      `json:"incident" pg:"incident,default:false"`
}

// FromPayload - Создаст событие готовое к записи в БД.
// Применит шифрование к сообщению.
func FromPayload(data *IncomingPayload, encoder encoder.Encoder) *Event {
	m := encoder.Encode([]byte(data.Message))

	return &Event{
		ID:         uuid.NewV4(),
		EventID:    data.EventID,
		Created:    data.Created,
		SystemName: data.SystemName,
		Message:    string(m),
	}
}
