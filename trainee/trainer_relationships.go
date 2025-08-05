package trainee

import (
	"context"
	"fmt"
)

type GetMyTrainersResponse struct {
	Trainers []*Trainer `json:"trainers"`
}

// GetMyTrainers retrieves all trainers for a given trainee
//
//encore:api private method=GET path=/trainee/trainers/:traineeID
func GetMyTrainers(ctx context.Context, traineeID string) (*GetMyTrainersResponse, error) {
	rows, err := db.Query(ctx, `
		SELECT 
			tr.id, 
			tr.user_id, 
			tr.specialization, 
			tr.years_of_experience,
			tr.rating
		FROM trainers tr
		JOIN trainer_trainee_relationships ttr ON tr.user_id = ttr.trainer_id
		WHERE ttr.trainee_id = $1 AND ttr.is_active = true
	`, traineeID)
	if err != nil {
		return nil, fmt.Errorf("failed to query trainers: %v", err)
	}
	defer rows.Close()

	var trainers []*Trainer
	for rows.Next() {
		var t Trainer
		if err := rows.Scan(
			&t.ID,
			&t.UserID,
			&t.Specialization,
			&t.YearsOfExperience,
			&t.Rating,
		); err != nil {
			return nil, fmt.Errorf("failed to scan trainer: %v", err)
		}
		trainers = append(trainers, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating trainers: %v", err)
	}

	return &GetMyTrainersResponse{Trainers: trainers}, nil
}
