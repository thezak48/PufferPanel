package pufferpanel

type Task struct {
	Name         string                    `json:"name"`
	CronSchedule string                    `json:"cronSchedule"`
	Description  string                    `json:"description,omitempty"`
	Operations   []ConditionalMetadataType `json:"operations" binding:"required"`
} //@name Task
