
-- First, create the users table if it doesn't exist (should match admin service)
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Create the trainee profile table
CREATE TABLE trainee_profiles (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    age INTEGER,
    height_cm DECIMAL(5,2),
    weight_kg DECIMAL(5,2),
    fitness_level VARCHAR(20) CHECK (fitness_level IN ('BEGINNER', 'INTERMEDIATE', 'ADVANCED')),
    medical_conditions TEXT,
    injuries TEXT,
    preferences TEXT, -- JSON array of preferences
    goals TEXT, -- JSON array of fitness goals
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_trainee_profile_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create trainer-trainee relationship table
CREATE TABLE trainer_trainee_relationships (
    id BIGSERIAL PRIMARY KEY,
    trainer_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    trainee_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    is_active BOOLEAN DEFAULT TRUE,
    start_date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    end_date TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(trainer_id, trainee_id)
);

-- Create workout templates
CREATE TABLE workout_templates (
    id BIGSERIAL PRIMARY KEY,
    trainer_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    duration_minutes INTEGER,
    difficulty VARCHAR(20) CHECK (difficulty IN ('BEGINNER', 'INTERMEDIATE', 'ADVANCED')),
    is_public BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_workout_template_trainer FOREIGN KEY (trainer_id) REFERENCES users(id) ON DELETE SET NULL
);

-- Create exercises table
CREATE TABLE exercises (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    muscle_group VARCHAR(255),
    equipment_required VARCHAR(255),
    video_url VARCHAR(512),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Create workout exercises (junction table)
CREATE TABLE workout_exercises (
    id BIGSERIAL PRIMARY KEY,
    workout_id BIGINT NOT NULL REFERENCES workout_templates(id) ON DELETE CASCADE,
    exercise_id BIGINT NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
    sets INTEGER,
    reps INTEGER,
    duration_seconds INTEGER,
    notes TEXT,
    order_index INTEGER NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_workout_exercise_workout FOREIGN KEY (workout_id) REFERENCES workout_templates(id) ON DELETE CASCADE,
    CONSTRAINT fk_workout_exercise_exercise FOREIGN KEY (exercise_id) REFERENCES exercises(id) ON DELETE CASCADE
);

-- Create assigned workouts table
CREATE TABLE assigned_workouts (
    id BIGSERIAL PRIMARY KEY,
    trainee_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    workout_id BIGINT NOT NULL REFERENCES workout_templates(id) ON DELETE CASCADE,
    assigned_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    assigned_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    due_date TIMESTAMPTZ,
    completed BOOLEAN DEFAULT FALSE,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_assigned_workout_trainee FOREIGN KEY (trainee_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_assigned_workout_workout FOREIGN KEY (workout_id) REFERENCES workout_templates(id) ON DELETE CASCADE,
    CONSTRAINT fk_assigned_workout_assigner FOREIGN KEY (assigned_by) REFERENCES users(id) ON DELETE SET NULL
);

-- Create workout logs table
CREATE TABLE workout_logs (
    id BIGSERIAL PRIMARY KEY,
    assigned_workout_id BIGINT REFERENCES assigned_workouts(id) ON DELETE SET NULL,
    trainee_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    workout_id BIGINT NOT NULL REFERENCES workout_templates(id) ON DELETE CASCADE,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ,
    duration_minutes INTEGER,
    notes TEXT,
    rating SMALLINT CHECK (rating BETWEEN 1 AND 5),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_workout_log_assigned_workout FOREIGN KEY (assigned_workout_id) REFERENCES assigned_workouts(id) ON DELETE SET NULL,
    CONSTRAINT fk_workout_log_trainee FOREIGN KEY (trainee_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_workout_log_workout FOREIGN KEY (workout_id) REFERENCES workout_templates(id) ON DELETE CASCADE
);

-- Create exercise logs table
CREATE TABLE exercise_logs (
    id BIGSERIAL PRIMARY KEY,
    workout_log_id BIGINT NOT NULL REFERENCES workout_logs(id) ON DELETE CASCADE,
    exercise_id BIGINT NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
    sets_completed INTEGER,
    reps_completed INTEGER,
    weight_kg DECIMAL(5,2),
    duration_seconds INTEGER,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_exercise_log_workout_log FOREIGN KEY (workout_log_id) REFERENCES workout_logs(id) ON DELETE CASCADE,
    CONSTRAINT fk_exercise_log_exercise FOREIGN KEY (exercise_id) REFERENCES exercises(id) ON DELETE CASCADE
);

-- Create progress metrics table
CREATE TABLE progress_metrics (
    id BIGSERIAL PRIMARY KEY,
    trainee_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    metric_type VARCHAR(50) NOT NULL CHECK (metric_type IN ('WEIGHT', 'BODY_FAT', 'WAIST_CIRCUMFERENCE', 'CHEST_CIRCUMFERENCE', 'ARM_CIRCUMFERENCE', 'THIGH_CIRCUMFERENCE')),
    value DECIMAL(7,2) NOT NULL,
    measured_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_progress_metric_trainee FOREIGN KEY (trainee_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create progress photos table
CREATE TABLE progress_photos (
    id BIGSERIAL PRIMARY KEY,
    trainee_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    photo_url VARCHAR(512) NOT NULL,
    taken_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    angle VARCHAR(50) CHECK (angle IN ('FRONT', 'SIDE', 'BACK')),
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_progress_photo_trainee FOREIGN KEY (trainee_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create messages table
CREATE TABLE messages (
    id BIGSERIAL PRIMARY KEY,
    sender_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    receiver_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    read_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_message_sender FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_message_receiver FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create indexes for better performance
CREATE INDEX idx_trainee_profiles_user_id ON trainee_profiles(user_id);
CREATE INDEX idx_trainer_trainee_relationships_trainer ON trainer_trainee_relationships(trainer_id);
CREATE INDEX idx_trainer_trainee_relationships_trainee ON trainer_trainee_relationships(trainee_id);
CREATE INDEX idx_assigned_workouts_trainee ON assigned_workouts(trainee_id);
CREATE INDEX idx_workout_logs_trainee ON workout_logs(trainee_id);
CREATE INDEX idx_workout_logs_workout ON workout_logs(workout_id);
CREATE INDEX idx_exercise_logs_workout_log ON exercise_logs(workout_log_id);
CREATE INDEX idx_progress_metrics_trainee ON progress_metrics(trainee_id);
CREATE INDEX idx_progress_photos_trainee ON progress_photos(trainee_id);
CREATE INDEX idx_messages_sender_receiver ON messages(sender_id, receiver_id);
CREATE INDEX idx_messages_created_at ON messages(created_at);