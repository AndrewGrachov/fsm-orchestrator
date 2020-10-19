package storage

import (
	"time"
)

type Status string

var (
	Initial    Status = "initial"
	Processing Status = "processing"
	Completed  Status = "completed"
	Failed     Status = "failed"
)

// Update operations must reference this fields by their json tag
type Object struct {
	ID          string    `bson:"_id" json:"id"`
	CreatedAt   time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time `bson:"updatedAt" json:"updatedAt"`
	CompletedAt time.Time `bson:"completedAt" json:"completedAt"`
	Error       string    `bson:"error" json:"error"`
	ObjectDTO
}

type ObjectDTO struct {
	CurrentStep  string                 `bson:"currentStep" json:"currentStep"`
	CommandGraph string                 `bson:"commandGraph" json:"commandGraph"`
	Status       Status                 `bson:"status" json:"status"`
	Params       map[string]interface{} `bson:"params" json:"params"`
}

type KV map[string]interface{}

// Storage provides easy to provide minimalistic approach to abstract persistent storage.
// Update operation receives map of fields which corresponds to object's json field tags by name.
type Storage interface {
	Create(obj ObjectDTO) (Object, error)
	FindById(id string) (Object, error)
	UpdateById(id string, update KV) error
}

// Repository provides utilitarian wrappings for internal usage
type Repository struct {
	Storage
}

func (r *Repository) CreateJob(obj ObjectDTO) (Object, error) {
	// can add validator or something like that here
	return r.Create(obj)
}

func (r *Repository) StartJob(id string) error {
	data := KV{
		"status": Processing,
	}

	return r.UpdateById(id, data)
}

func (r *Repository) FailJob(id string, err error) error {
	data := KV{
		"status": Failed,
		"error":  err.Error(),
	}

	return r.UpdateById(id, data)
}

func (r *Repository) CompleteJob(id string) error {
	data := KV{
		"status": Completed,
	}

	return r.UpdateById(id, data)
}