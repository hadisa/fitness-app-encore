// Package trainee provides functionality for trainee management and operations.
package trainee

import (
	"context"
	"time"

	"encore.dev/storage/sqldb"
)

// Trainee represents a trainee user in the system
type Trainee struct {
	ID           string    `json:"id"`
	UserID       int       `json:"user_id"`
	Age          int       `json:"age"`
	Height       float64   `json:"height"`
	Weight       float64   `json:"weight"`
	FitnessGoals []string  `json:"fitness_goals"`
	Injuries     []string  `json:"injuries"`
	Preferences  []string  `json:"preferences"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Workout represents a workout plan
type Workout struct {
	ID          string    `json:"id"`
	TrainerID   string    `json:"trainer_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Duration    int       `json:"duration"`
	Difficulty  string    `json:"difficulty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CompletedWorkout represents a workout completed by a trainee
type CompletedWorkout struct {
	ID        string   `json:"id"`
	TraineeID string   `json:"trainee_id"`
	Workout   *Workout `json:"workout"`
	Date      string   `json:"date"`
	Duration  int      `json:"duration"`
	Notes     *string  `json:"notes,omitempty"`
	Rating    *int     `json:"rating,omitempty"`
}

// ProgressPhoto represents a progress photo uploaded by a trainee
type ProgressPhoto struct {
	ID        string    `json:"id"`
	TraineeID string    `json:"trainee_id"`
	URL       string    `json:"url"`
	Date      string    `json:"date"`
	Notes     *string   `json:"notes,omitempty"`
	Angle     string    `json:"angle"`
	CreatedAt time.Time `json:"created_at"`
}

// Trainer represents a trainer in the system
type Trainer struct {
	ID                string   `json:"id"`
	UserID            int      `json:"user_id"`
	Specialization    []string `json:"specialization"`
	YearsOfExperience int      `json:"years_of_experience"`
	Rating            *float64 `json:"rating,omitempty"`
}

// Message represents a message between trainee and trainer
type Message struct {
	ID        string    `json:"id"`
	SenderID  string    `json:"sender_id"`
	Content   string    `json:"content"`
	Timestamp string    `json:"timestamp"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

// UpdateProfileRequest contains the data needed to update a trainee's profile
type UpdateProfileRequest struct {
	TraineeID    string   `json:"trainee_id"`
	Age          *int     `json:"age,omitempty"`
	Height       *float64 `json:"height,omitempty"`
	Weight       *float64 `json:"weight,omitempty"`
	FitnessGoals []string `json:"fitness_goals,omitempty"`
	Injuries     []string `json:"injuries,omitempty"`
	Preferences  []string `json:"preferences,omitempty"`
}

// UpdateProfile updates the trainee's profile
//
//encore:api private method=POST path=/trainee/profile
func UpdateProfile(ctx context.Context, req *UpdateProfileRequest) (*Trainee, error) {
	// Start a transaction
	tx, err := db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Update trainee profile
	_, err = tx.Exec(ctx, `
		UPDATE trainees
		SET 
			age = COALESCE($1, age),
			height = COALESCE($2, height),
			weight = COALESCE($3, weight),
			fitness_goals = COALESCE($4, fitness_goals),
			injuries = COALESCE($5, injuries),
			preferences = COALESCE($6, preferences),
			updated_at = NOW()
		WHERE id = $7
	`,
		req.Age,
		req.Height,
		req.Weight,
		req.FitnessGoals,
		req.Injuries,
		req.Preferences,
		req.TraineeID,
	)
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Return the updated profile
	return GetTraineeByID(ctx, req.TraineeID)
}

// GetTraineeByID retrieves a trainee by ID
func GetTraineeByID(ctx context.Context, traineeID string) (*Trainee, error) {
	var t Trainee
	err := db.QueryRow(ctx, `
		SELECT id, user_id, age, height, weight, fitness_goals, injuries, preferences, created_at, updated_at
		FROM trainees
		WHERE id = $1
	`, traineeID).Scan(
		&t.ID,
		&t.UserID,
		&t.Age,
		&t.Height,
		&t.Weight,
		&t.FitnessGoals,
		&t.Injuries,
		&t.Preferences,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// GetTraineeWorkouts retrieves all workouts for a trainee
func GetTraineeWorkouts(ctx context.Context, traineeID string) ([]*Workout, error) {
	rows, err := db.Query(ctx, `
		SELECT w.id, w.trainer_id, w.name, w.description, w.duration, w.difficulty, w.created_at, w.updated_at
		FROM workouts w
		JOIN trainee_workouts tw ON w.id = tw.workout_id
		WHERE tw.trainee_id = $1
		ORDER BY w.created_at DESC
	`, traineeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workouts []*Workout
	for rows.Next() {
		var w Workout
		if err := rows.Scan(
			&w.ID,
			&w.TrainerID,
			&w.Name,
			&w.Description,
			&w.Duration,
			&w.Difficulty,
			&w.CreatedAt,
			&w.UpdatedAt,
		); err != nil {
			return nil, err
		}
		workouts = append(workouts, &w)
	}
	return workouts, nil
}

// LogWorkout logs a completed workout for a trainee
func LogWorkout(ctx context.Context, traineeID string, workoutID string, duration int, notes *string, rating *int) (*CompletedWorkout, error) {
	// Start a transaction
	tx, err := db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Get workout details
	var workout Workout
	err = tx.QueryRow(ctx, `
		SELECT id, trainer_id, name, description, duration, difficulty
		FROM workouts
		WHERE id = $1
	`, workoutID).Scan(
		&workout.ID,
		&workout.TrainerID,
		&workout.Name,
		&workout.Description,
		&workout.Duration,
		&workout.Difficulty,
	)
	if err != nil {
		return nil, err
	}

	// Log the completed workout
	var completedWorkoutID string
	err = tx.QueryRow(ctx, `
		INSERT INTO completed_workouts (trainee_id, workout_id, date, duration, notes, rating)
		VALUES ($1, $2, NOW(), $3, $4, $5)
		RETURNING id
	`, traineeID, workoutID, duration, notes, rating).Scan(&completedWorkoutID)
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &CompletedWorkout{
		ID:        completedWorkoutID,
		TraineeID: traineeID,
		Workout:   &workout,
		Date:      time.Now().Format(time.RFC3339),
		Duration:  duration,
		Notes:     notes,
		Rating:    rating,
	}, nil
}

// Define the database connection
var db = sqldb.Named("trainee")
