package trainee

import (
	"context"
	"fmt"
)

// GetWorkoutByID retrieves a workout by its ID
//
//encore:api private method=GET path=/workouts/:workoutID
func GetWorkoutByID(ctx context.Context, workoutID string) (*Workout, error) {
	var workout Workout
	err := db.QueryRow(ctx, `
		SELECT id, trainer_id, name, description, duration, difficulty, created_at, updated_at
		FROM workout_templates
		WHERE id = $1
	`, workoutID).Scan(
		&workout.ID,
		&workout.TrainerID,
		&workout.Name,
		&workout.Description,
		&workout.Duration,
		&workout.Difficulty,
		&workout.CreatedAt,
		&workout.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get workout: %v", err)
	}

	return &workout, nil
}
