package trainer

import (
	"context"
	"time"

	"encore.dev/storage/sqldb"
)

// AssignWorkoutToClient assigns a workout to a client
//
//encore:api private method=POST path=/trainer/assign-workout
func AssignWorkoutToClient(ctx context.Context, params *AssignWorkoutRequest) (*AssignedWorkout, error) {
	// Start a transaction
	tx, err := sqldb.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Insert the assigned workout
	var assignedWorkout AssignedWorkout
	err = tx.QueryRow(ctx, `
		INSERT INTO assigned_workouts (trainee_id, workout_id, assigned_by, due_date, completed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, false, NOW(), NOW())
		RETURNING id, trainee_id, workout_id, assigned_by, assigned_at, due_date, completed, completed_at, created_at, updated_at
	`,
		params.ClientID,
		params.WorkoutID,
		params.TrainerID,
		params.DueDate,
	).Scan(
		&assignedWorkout.ID,
		&assignedWorkout.TraineeID,
		&assignedWorkout.WorkoutID,
		&assignedWorkout.AssignedBy,
		&assignedWorkout.AssignedAt,
		&assignedWorkout.DueDate,
		&assignedWorkout.Completed,
		&assignedWorkout.CompletedAt,
		&assignedWorkout.CreatedAt,
		&assignedWorkout.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &assignedWorkout, nil
}

// AssignWorkoutRequest contains the data needed to assign a workout to a client
type AssignWorkoutRequest struct {
	TrainerID string     `json:"trainer_id"`
	ClientID  string     `json:"client_id"`
	WorkoutID string     `json:"workout_id"`
	DueDate   *time.Time `json:"due_date,omitempty"`
}

// AssignedWorkout represents a workout assigned to a client
type AssignedWorkout struct {
	ID          string     `json:"id"`
	TraineeID   string     `json:"trainee_id"`
	WorkoutID   string     `json:"workout_id"`
	AssignedBy  *string    `json:"assigned_by,omitempty"`
	AssignedAt  time.Time  `json:"assigned_at"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Completed   bool       `json:"completed"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
