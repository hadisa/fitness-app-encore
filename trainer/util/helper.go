package util

import (
	"context"
	"fmt"

	"encore.app/trainee"
)

// ValidateClientBelongsToTrainer checks if a client belongs to a trainer
func ValidateClientBelongsToTrainer(ctx context.Context, clientID, trainerID string) (bool, error) {
	trainers, err := trainee.GetMyTrainers(ctx, clientID)
	if err != nil {
		return false, fmt.Errorf("failed to get trainers for client: %v", err)
	}
	for _, trainer := range trainers.Trainers {
		if trainer.ID == trainerID {
			return true, nil
		}
	}
	return false, nil
}

// ValidateWorkoutBelongsToTrainer checks if a workout belongs to a trainer
func ValidateWorkoutBelongsToTrainer(ctx context.Context, workoutID, trainerID string) (bool, error) {
	workout, err := trainee.GetWorkoutByID(ctx, workoutID)
	if err != nil {
		return false, fmt.Errorf("failed to get workout: %v", err)
	}
	// Check if the workout's trainer ID matches the provided trainer ID
	return workout.TrainerID == trainerID, nil
}
