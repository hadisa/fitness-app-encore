// Package trainer provides functionality for trainer management and operations.
package trainer

import (
	"context"
	"fmt"
	"time"

	"encore.dev/storage/sqldb"
)

// Trainer represents a trainer in the system
type Trainer struct {
	ID                string          `json:"id"`
	UserID            int             `json:"user_id"`
	Specialization    []string        `json:"specialization"`
	Certifications    []Certification `json:"certifications"`
	YearsOfExperience int             `json:"years_of_experience"`
	Bio               string          `json:"bio"`
	HourlyRate        float64         `json:"hourly_rate"`
	Rating            *float64        `json:"rating,omitempty"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

// Certification represents a trainer's certification
type Certification struct {
	ID                  string    `json:"id"`
	TrainerID           string    `json:"trainer_id"`
	Name                string    `json:"name"`
	IssuingOrganization string    `json:"issuing_organization"`
	DateIssued          string    `json:"date_issued"`
	CredentialID        *string   `json:"credential_id,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
}

// AvailabilitySlot represents a time slot when the trainer is available
type AvailabilitySlot struct {
	ID          string    `json:"id"`
	TrainerID   string    `json:"trainer_id"`
	DayOfWeek   int       `json:"day_of_week"`
	StartTime   string    `json:"start_time"`
	EndTime     string    `json:"end_time"`
	IsRecurring bool      `json:"is_recurring"`
	IsBooked    bool      `json:"is_booked"`
	CreatedAt   time.Time `json:"created_at"`
}

// Appointment represents a scheduled appointment with a client
type Appointment struct {
	ID        string            `json:"id"`
	ClientID  string            `json:"client_id"`
	TrainerID string            `json:"trainer_id"`
	SlotID    string            `json:"slot_id"`
	Slot      *AvailabilitySlot `json:"slot,omitempty"`
	Status    string            `json:"status"` // e.g., "scheduled", "completed", "cancelled"
	Notes     *string           `json:"notes,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// UpdateProfileRequest contains the data needed to update a trainer's profile
type UpdateProfileRequest struct {
	TrainerID         string   `json:"trainer_id"`
	Specialization    []string `json:"specialization,omitempty"`
	YearsOfExperience *int     `json:"years_of_experience,omitempty"`
	Bio               *string  `json:"bio,omitempty"`
	HourlyRate        *float64 `json:"hourly_rate,omitempty"`
}

// UpdateCertificationsRequest contains the data needed to update a trainer's certifications
type UpdateCertificationsRequest struct {
	TrainerID      string          `json:"trainer_id"`
	Certifications []Certification `json:"certifications"`
}

// UpdateCertificationsResponse is the response from updating trainer certifications
type UpdateCertificationsResponse struct {
	Certifications []Certification `json:"certifications"`
}

// SetAvailabilityRequest contains the data needed to set a trainer's availability
type SetAvailabilityRequest struct {
	TrainerID string             `json:"trainer_id"`
	Slots     []AvailabilitySlot `json:"slots"`
}

// SetAvailabilityResponse is the response from setting trainer availability
type SetAvailabilityResponse struct {
	Slots []AvailabilitySlot `json:"slots"`
}

// CreateAppointmentRequest contains the data needed to create an appointment
type CreateAppointmentRequest struct {
	ClientID string `json:"client_id"`
	SlotID   string `json:"slot_id"`
}

// GetTrainerProfileRequest is the request to get a trainer's profile
type GetTrainerProfileRequest struct {
	// TrainerID string `json:"trainer_id"`
	TrainerID string `json:"-"`
}

// GetTrainerClientsResponse is the response for getting a trainer's clients
type GetTrainerClientsResponse struct {
	ClientIDs []string `json:"client_ids"`
}

// GetTrainerAvailabilityResponse is the response for getting a trainer's availability
type GetTrainerAvailabilityResponse struct {
	Slots []AvailabilitySlot `json:"slots"`
}

// GetTrainerAppointmentsResponse is the response for getting a trainer's appointments
type GetTrainerAppointmentsResponse struct {
	Appointments []Appointment `json:"appointments"`
}

// GetTrainerProfile returns the trainer's profile
//
// encore:api private method=GET path=/trainer/profile/:trainerID
func GetTrainerProfile(ctx context.Context, trainerID string) (*Trainer, error) {
	var trainer Trainer
	err := db.QueryRow(ctx, `
		SELECT id, user_id, specialization, years_of_experience, bio, hourly_rate, rating, created_at, updated_at
		FROM trainers
		WHERE user_id = $1
	`, trainerID).Scan(
		&trainer.ID,
		&trainer.UserID,
		&trainer.Specialization,
		&trainer.YearsOfExperience,
		&trainer.Bio,
		&trainer.HourlyRate,
		&trainer.Rating,
		&trainer.CreatedAt,
		&trainer.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get trainer profile: %w", err)
	}

	// Get certifications
	certRows, err := db.Query(ctx, `
		SELECT id, name, issuing_organization, date_issued, credential_id, created_at
		FROM trainer_certifications
		WHERE trainer_id = $1
	`, trainerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get trainer certifications: %w", err)
	}
	defer certRows.Close()

	var certifications []Certification
	for certRows.Next() {
		var cert Certification
		err := certRows.Scan(
			&cert.ID,
			&cert.Name,
			&cert.IssuingOrganization,
			&cert.DateIssued,
			&cert.CredentialID,
			&cert.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan certification: %w", err)
		}
		cert.TrainerID = trainerID
		certifications = append(certifications, cert)
	}
	trainer.Certifications = certifications

	return &trainer, nil
}

// UpdateTrainerProfile updates the trainer's profile
//
// encore:api private method=POST path=/trainer/profile/update
func UpdateTrainerProfile(ctx context.Context, req *UpdateProfileRequest) (*Trainer, error) {
	// Start a transaction
	tx, err := db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Update trainer profile
	_, err = tx.Exec(ctx, `
		UPDATE trainers
		SET
			specialization = COALESCE($1, specialization),
			years_of_experience = COALESCE($2, years_of_experience),
			bio = COALESCE($3, bio),
			hourly_rate = COALESCE($4, hourly_rate),
			updated_at = NOW()
		WHERE user_id = $5
	`,
		req.Specialization,
		req.YearsOfExperience,
		req.Bio,
		req.HourlyRate,
		req.TrainerID,
	)
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Return the updated profile
	return GetTrainerProfile(ctx, req.TrainerID)
}

// UpdateTrainerCertifications updates a trainer's certifications
//
// encore:api private method=POST path=/trainer/certifications/update
func UpdateTrainerCertifications(ctx context.Context, req *UpdateCertificationsRequest) (*UpdateCertificationsResponse, error) {
	// Start a transaction
	tx, err := db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Delete existing certifications
	_, err = tx.Exec(ctx, `DELETE FROM trainer_certifications WHERE trainer_id = $1`, req.TrainerID)
	if err != nil {
		return nil, err
	}

	// Insert new certifications
	for _, cert := range req.Certifications {
		_, err = tx.Exec(ctx, `
			INSERT INTO trainer_certifications (id, trainer_id, name, issuing_organization, date_issued, credential_id)
			VALUES ($1, $2, $3, $4, $5, $6)
		`,
			cert.ID,
			req.TrainerID,
			cert.Name,
			cert.IssuingOrganization,
			cert.DateIssued,
			cert.CredentialID,
		)
		if err != nil {
			return nil, err
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Return the updated certifications
	return &UpdateCertificationsResponse{
		Certifications: req.Certifications,
	}, nil
}

// SetAvailability sets the trainer's availability
//
// encore:api private method=POST path=/trainer/availability/set
func SetAvailability(ctx context.Context, req *SetAvailabilityRequest) (*SetAvailabilityResponse, error) {
	// Start a transaction
	tx, err := db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Delete existing availability slots
	_, err = tx.Exec(ctx, `DELETE FROM availability_slots WHERE trainer_id = $1`, req.TrainerID)
	if err != nil {
		return nil, err
	}

	// Insert new availability slots
	for _, slot := range req.Slots {
		_, err = tx.Exec(ctx, `
			INSERT INTO availability_slots (id, trainer_id, day_of_week, start_time, end_time, is_recurring, is_booked)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`,
			slot.ID,
			req.TrainerID,
			slot.DayOfWeek,
			slot.StartTime,
			slot.EndTime,
			slot.IsRecurring,
			slot.IsBooked,
		)
		if err != nil {
			return nil, err
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Return the updated availability slots
	return &SetAvailabilityResponse{
		Slots: req.Slots,
	}, nil
}

// CreateAppointment creates a new appointment
//
// encore:api private method=POST path=/trainer/appointments/create
func CreateAppointment(ctx context.Context, req *CreateAppointmentRequest) (*Appointment, error) {
	// Start a transaction
	tx, err := db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Create the appointment
	var appointment Appointment
	err = tx.QueryRow(ctx, `
		INSERT INTO appointments (client_id, slot_id, status)
		VALUES ($1, $2, 'scheduled')
		RETURNING id, client_id, slot_id, status, created_at, updated_at
	`,
		req.ClientID,
		req.SlotID,
	).Scan(
		&appointment.ID,
		&appointment.ClientID,
		&appointment.SlotID,
		&appointment.Status,
		&appointment.CreatedAt,
		&appointment.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Update the slot to mark it as booked
	_, err = tx.Exec(ctx, `
		UPDATE availability_slots
		SET is_booked = true
		WHERE id = $1
	`, req.SlotID)
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &appointment, nil
}

// GetTrainerClientsRequest is the request to get a trainer's clients
type GetTrainerClientsRequest struct {
	TrainerID string `json:"-"` // URL parameter, not in JSON body
}

// GetTrainerClients returns a list of clients for a trainer
//
// encore:api private method=GET path=/trainer/clients/:trainerID
func GetTrainerClients(ctx context.Context, trainerID string) (*GetTrainerClientsResponse, error) {
	rows, err := db.Query(ctx, `
		SELECT DISTINCT client_id
		FROM appointments
		JOIN availability_slots ON appointments.slot_id = availability_slots.id
		WHERE availability_slots.trainer_id = $1
	`, trainerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clientIDs []string
	for rows.Next() {
		var clientID string
		if err := rows.Scan(&clientID); err != nil {
			return nil, err
		}
		clientIDs = append(clientIDs, clientID)
	}

	return &GetTrainerClientsResponse{
		ClientIDs: clientIDs,
	}, nil
}

// GetTrainerAvailabilityRequest is the request to get a trainer's availability
type GetTrainerAvailabilityRequest struct {
	TrainerID string `json:"-"` // URL parameter, not in JSON body
}

// GetTrainerAvailability returns a trainer's availability
//
// encore:api private method=GET path=/trainer/availability/:trainerID
func GetTrainerAvailability(ctx context.Context, trainerID string) (*GetTrainerAvailabilityResponse, error) {
	rows, err := db.Query(ctx, `
		SELECT id, day_of_week, start_time, end_time, is_recurring, is_booked, created_at
		FROM availability_slots
		WHERE trainer_id = $1
	`, trainerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []AvailabilitySlot
	for rows.Next() {
		var slot AvailabilitySlot
		err = rows.Scan(
			&slot.ID,
			&slot.DayOfWeek,
			&slot.StartTime,
			&slot.EndTime,
			&slot.IsRecurring,
			&slot.IsBooked,
			&slot.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		slots = append(slots, slot)
	}

	return &GetTrainerAvailabilityResponse{
		Slots: slots,
	}, nil
}

// GetTrainerAppointmentsRequest is the request to get a trainer's appointments
type GetTrainerAppointmentsRequest struct {
	TrainerID string `json:"-"` // URL parameter, not in JSON body
}

// GetTrainerAppointments returns a trainer's appointments
//
// encore:api private method=GET path=/trainer/appointments/:trainerID
func GetTrainerAppointments(ctx context.Context, trainerID string) (*GetTrainerAppointmentsResponse, error) {
	rows, err := db.Query(ctx, `
		SELECT a.id, a.client_id, a.slot_id, a.status, a.notes, a.created_at, a.updated_at,
		       av.day_of_week, av.start_time, av.end_time, av.is_recurring
		FROM appointments a
		JOIN availability_slots av ON a.slot_id = av.id		
		WHERE a.trainer_id = $1
	`, trainerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get trainer appointments: %w", err)
	}
	defer rows.Close()

	var appointments []Appointment
	for rows.Next() {
		var appt Appointment
		var slot AvailabilitySlot
		err := rows.Scan(
			&appt.ID,
			&appt.ClientID,
			&appt.SlotID,
			&appt.Status,
			&appt.Notes,
			&appt.CreatedAt,
			&appt.UpdatedAt,
			&slot.DayOfWeek,
			&slot.StartTime,
			&slot.EndTime,
			&slot.IsRecurring,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan appointment: %w", err)
		}
		slot.TrainerID = trainerID
		appt.Slot = &slot
		appointments = append(appointments, appt)
	}

	return &GetTrainerAppointmentsResponse{
		Appointments: appointments,
	}, nil
}

// Define the database connection
var db = sqldb.Named("trainer")
