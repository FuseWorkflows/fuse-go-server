package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/FuseWorkflows/fuse-go-server/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
	_ "github.com/lib/pq" // postgres driver
)

var ErrNotFound = errors.New("resource not found")

type DB struct {
	*sql.DB
}

func InitDB() (*DB, error) {
	// Get environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Build connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)

	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Ping the database to check the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	return &DB{db}, nil
}

// GetUserByID retrieves a user by ID
func (db *DB) GetUserByID(userID string) (*models.User, error) {
	var user models.User
	err := db.QueryRowContext(context.Background(), "SELECT * FROM users WHERE id = $1", userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Tier,
		&user.Trial,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error fetching user: %w", err)
	}

	return &user, nil
}

// GetUserByEmail retrieves a user by email
func (db *DB) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := db.QueryRowContext(context.Background(), "SELECT * FROM users WHERE email = $1", email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Tier,
		&user.Trial,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error fetching user: %w", err)
	}

	return &user, nil
}

// GetUsers retrieves all users
func (db *DB) GetUsers() ([]models.User, error) {
	var users []models.User
	rows, err := db.QueryContext(context.Background(), "SELECT * FROM users")
	if err != nil {
		return nil, fmt.Errorf("error fetching users: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.Tier,
			&user.Trial,
		); err != nil {
			return nil, fmt.Errorf("error scanning user: %w", err)
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through rows: %w", err)
	}

	return users, nil
}

// CreateUser creates a new user
func (db *DB) CreateUser(user *models.User) (*models.User, error) {
	ctx := context.Background()

	// Generate a new UUID
	user.ID = uuid.New().String()

	// Insert the user with the generated UUID and return the ID for verification
	err := db.QueryRowContext(ctx,
		"INSERT INTO users (id, username, email, password) VALUES ($1, $2, $3, $4) RETURNING id",
		user.ID, user.Username, user.Email, user.Password).Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	// Fetch the user before returning
	createdUser, err := db.GetUserByID(user.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching user: %w", err)
	}
	return createdUser, nil
}

// UpdateUser updates an existing user
func (db *DB) UpdateUser(userID string, user *models.User) (*models.User, error) {
	ctx := context.Background()

	result, err := db.ExecContext(ctx, "UPDATE users SET username = $1, email = $2, password = $3, tier = $4, trial = $5, updated_at = NOW() WHERE id = $6",
		user.Username, user.Email, user.Password, user.Tier, user.Trial, userID)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("error getting rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return nil, ErrNotFound
	}

	// Fetch the user before returning
	updatedUser, err := db.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching user: %w", err)
	}
	return updatedUser, nil
}

// DeleteUser deletes an existing user
func (db *DB) DeleteUser(userID string) error {
	ctx := context.Background()
	result, err := db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// GetChannelByID retrieves a channel by ID
func (db *DB) GetChannelByID(channelID string) (*models.Channel, error) {
	var channel models.Channel
	err := db.QueryRowContext(context.Background(), "SELECT * FROM channels WHERE id = $1", channelID).Scan(
		&channel.ID,
		&channel.Name,
		&channel.API_KEY,
		&channel.Owner.ID,
		&channel.CreatedAt,
		&channel.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error fetching channel: %w", err)
	}

	// Fetch the owner data using the owner ID
	owner, err := db.GetUserByID(channel.Owner.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching owner: %w", err)
	}

	// Assign the fetched owner to the channel
	channel.Owner = *owner // Dereference the owner pointer

	return &channel, nil
}

// GetChannelsByUser retrieves channels by user ID
func (db *DB) GetChannelsByUser(userID string) ([]models.Channel, error) {
	var channels []models.Channel

	// Fetch the owner data using the owner ID
	owner, err := db.GetOwner(userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching owner: %w", err)
	}

	rows, err := db.QueryContext(context.Background(), "SELECT * FROM channels WHERE owner_id = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching channels: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var channel models.Channel
		if err := rows.Scan(
			&channel.ID,
			&channel.Name,
			&channel.API_KEY,
			&channel.Owner.ID,
			&channel.CreatedAt,
			&channel.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning channel: %w", err)
		}

		// Assign the fetched owner to the channel
		channel.Owner = *owner // Dereference the owner pointer

		if err != nil {
			return nil, fmt.Errorf("error fetching owner: %w", err)
		}

		channels = append(channels, channel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through rows: %w", err)
	}

	return channels, nil
}

// CreateChannel creates a new channel
func (db *DB) CreateChannel(channel *models.Channel) (*models.Channel, error) {
	ctx := context.Background()

	// Generate a new UUID
	channel.ID = uuid.New().String()

	// Insert the channel with the generated UUID
	err := db.QueryRowContext(ctx, "INSERT INTO channels (id, name, api_key, owner_id) VALUES ($1, $2, $3, $4) RETURNING id",
		channel.ID, channel.Name, channel.API_KEY, channel.Owner.ID).Scan(&channel.ID)
	if err != nil {
		return nil, fmt.Errorf("error creating channel: %w", err)
	}

	// fecth the channel before returning
	createdChannel, err := db.GetChannelByID(channel.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching channel: %w", err)
	}
	return createdChannel, nil
}

// UpdateChannel updates an existing channel
func (db *DB) UpdateChannel(channelID string, channel *models.Channel) (*models.Channel, error) {
	ctx := context.Background()
	query := "UPDATE channels SET"
	params := []interface{}{}
	paramCounter := 1

	// Dynamically add fields to the query if they are not empty
	if channel.Name != "" {
		query += fmt.Sprintf(" name = $%d,", paramCounter)
		params = append(params, channel.Name)
		paramCounter++
	}
	if channel.API_KEY != "" {
		query += fmt.Sprintf(" api_key = $%d,", paramCounter)
		params = append(params, channel.API_KEY)
		paramCounter++
	}

	query += fmt.Sprintf(" updated_at = $%d", paramCounter)
	params = append(params, time.Now()) // Bind the current time
	paramCounter++

	// Remove trailing comma and add WHERE clause
	query = strings.TrimSuffix(query, ",") + fmt.Sprintf(" WHERE id = $%d", paramCounter)
	params = append(params, channelID)

	// Execute the query
	result, err := db.ExecContext(ctx, query, params...)
	if err != nil {
		return nil, fmt.Errorf("error updating channel: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("error getting rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return nil, ErrNotFound
	}

	// Fetch the updated channel before returning
	updatedChannel, err := db.GetChannelByID(channelID)
	if err != nil {
		return nil, fmt.Errorf("error fetching channel: %w", err)
	}
	return updatedChannel, nil
}

// DeleteChannel deletes an existing channel
func (db *DB) DeleteChannel(channelID string) error {
	ctx := context.Background()
	result, err := db.ExecContext(ctx, "DELETE FROM channels WHERE id = $1", channelID)
	if err != nil {
		return fmt.Errorf("error deleting channel: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// GetVideoByID retrieves a video by ID
func (db *DB) GetVideoByID(videoID string) (*models.Video, error) {
	var video models.Video
	var keywords []byte
	err := db.QueryRowContext(context.Background(), "SELECT * FROM videos WHERE id = $1", videoID).Scan(
		&video.ID,
		&video.Status,
		&video.Resources,
		&video.Title,
		&video.Description,
		&keywords,
		&video.Category,
		&video.PrivacyStatus,
		&video.Channel.ID,
		&video.CreatedAt,
		&video.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error fetching video: %w", err)
	}

	// Convert byte array to string array using pq.Array
	if err := pq.Array(&video.Keywords).Scan(keywords); err != nil {
		return nil, fmt.Errorf("error converting keywords array: %w", err)
	}

	// Fetch the channel data using the channel ID
	channel, err := db.GetChannelByID(video.Channel.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching channel: %w", err)
	}

	// Assign the fetched channel to the video
	video.Channel = *channel // Dereference the channel pointer

	video.Iterations, err = db.GetIterationsByVideo(video.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching iterations: %w", err)
	}

	video.Editors, err = db.GetEditorsByVideo(video.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching editors: %w", err)
	}

	return &video, nil
}

// GetVideosByUser retrieves videos by user ID
func (db *DB) GetVideosByUser(userID string) ([]models.Video, error) {
	var videos []models.Video
	rows, err := db.QueryContext(context.Background(), "SELECT v.* FROM videos v JOIN channels c ON v.channel_id = c.id WHERE c.owner_id = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching videos: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var video models.Video
		var keywords []byte
		if err := rows.Scan(
			&video.ID,
			&video.Status,
			&video.Resources,
			&video.Title,
			&video.Description,
			&keywords,
			&video.Category,
			&video.PrivacyStatus,
			&video.Channel.ID,
			&video.CreatedAt,
			&video.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning video: %w", err)
		}

		// Convert byte array to string array using pq.Array
		if err := pq.Array(&video.Keywords).Scan(keywords); err != nil {
			return nil, fmt.Errorf("error converting keywords array: %w", err)
		}

		// Fetch the channel data using the channel ID
		channel, err := db.GetChannelByID(video.Channel.ID)
		if err != nil {
			return nil, fmt.Errorf("error fetching channel: %w", err)
		}

		// Assign the fetched channel to the video
		video.Channel = *channel // Dereference the channel pointer

		video.Iterations, err = db.GetIterationsByVideo(video.ID)
		if err != nil {
			return nil, fmt.Errorf("error fetching iterations: %w", err)
		}

		video.Editors, err = db.GetEditorsByVideo(video.ID)
		if err != nil {
			return nil, fmt.Errorf("error fetching editors: %w", err)
		}

		videos = append(videos, video)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through rows: %w", err)
	}

	return videos, nil
}

// GetVideosByChannel retrieves videos by channel ID
func (db *DB) GetVideosByChannel(channelID string) ([]models.Video, error) {
	var videos []models.Video
	rows, err := db.QueryContext(context.Background(), "SELECT * FROM videos WHERE channel_id = $1", channelID)
	if err != nil {
		return nil, fmt.Errorf("error fetching videos: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var video models.Video
		var keywords []byte
		if err := rows.Scan(
			&video.ID,
			&video.Status,
			&video.Resources,
			&video.Title,
			&video.Description,
			&keywords,
			&video.Category,
			&video.PrivacyStatus,
			&video.Channel.ID,
			&video.CreatedAt,
			&video.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning video: %w", err)
		}

		// Convert byte array to string array using pq.Array
		if err := pq.Array(&video.Keywords).Scan(keywords); err != nil {
			return nil, fmt.Errorf("error converting keywords array: %w", err)
		}

		// Fetch the channel data using the channel ID
		channel, err := db.GetChannelByID(video.Channel.ID)
		if err != nil {
			return nil, fmt.Errorf("error fetching channel: %w", err)
		}

		// Assign the fetched channel to the video
		video.Channel = *channel // Dereference the channel pointer

		video.Iterations, err = db.GetIterationsByVideo(video.ID)
		if err != nil {
			return nil, fmt.Errorf("error fetching iterations: %w", err)
		}

		video.Editors, err = db.GetEditorsByVideo(video.ID)
		if err != nil {
			return nil, fmt.Errorf("error fetching editors: %w", err)
		}

		videos = append(videos, video)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through rows: %w", err)
	}

	return videos, nil
}

// CreateVideo creates a new video
func (db *DB) CreateVideo(video *models.Video) (*models.Video, error) {
	ctx := context.Background()

	// Format keywords into a PostgreSQL array literal
	keywords := "{}" // Default to an empty array literal
	if len(video.Keywords) > 0 {
		keywords = fmt.Sprintf("{%s}", strings.Join(video.Keywords, ", "))
	}

	video.ID = uuid.New().String()

	// Use QueryRowContext and RETURNING to get the video ID
	err := db.QueryRowContext(ctx, `
		INSERT INTO videos (status, resources, title, description, keywords, category, privacy_status, channel_id) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`,
		video.Status, video.Resources, video.Title, video.Description, keywords, video.Category, video.PrivacyStatus, video.Channel.ID,
	).Scan(&video.ID)

	if err != nil {
		return nil, fmt.Errorf("error creating video: %w", err)
	}

	// Assign editors to the video
	for _, editor := range video.Editors {
		_, err = db.AddEditorToVideo(video.ID, editor.ID)
		if err != nil {
			return nil, fmt.Errorf("error assigning editor to video: %w", err)
		}
	}

	// fecth the channel before returning
	createdVideo, err := db.GetVideoByID(video.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching video: %w", err)
	}
	return createdVideo, nil
}

// UpdateVideo updates an existing video
func (db *DB) UpdateVideo(videoID string, video *models.Video) (*models.Video, error) {
	ctx := context.Background()
	query := "UPDATE videos SET"
	params := []interface{}{}
	paramCounter := 1

	// Dynamically add fields to the query if they are not empty
	if video.Status != "" {
		// Check if the status is valid
		validStatus := false
		for _, s := range []models.Status{models.Pending, models.Published, models.Draft} {
			if video.Status == s {
				validStatus = true
				break
			}
		}

		if validStatus {
			query += fmt.Sprintf(" status = $%d,", paramCounter)
			params = append(params, video.Status)
			paramCounter++
		} else {
			return nil, fmt.Errorf("invalid video status: %s", video.Status)
		}
	}
	if video.Resources != "" {
		query += fmt.Sprintf(" resources = $%d,", paramCounter)
		params = append(params, video.Resources)
		paramCounter++
	}
	if video.Title != "" {
		query += fmt.Sprintf(" title = $%d,", paramCounter)
		params = append(params, video.Title)
		paramCounter++
	}
	if video.Description != "" {
		query += fmt.Sprintf(" description = $%d,", paramCounter)
		params = append(params, video.Description)
		paramCounter++
	}
	if len(video.Keywords) > 0 {
		keywords := fmt.Sprintf("{%s}", strings.Join(video.Keywords, ", ")) // Remove single quotes
		query += fmt.Sprintf(" keywords = $%d,", paramCounter)
		params = append(params, keywords)
		paramCounter++
	}
	if video.Category != "" {
		query += fmt.Sprintf(" category = $%d,", paramCounter)
		params = append(params, video.Category)
		paramCounter++
	}
	if video.PrivacyStatus {
		query += fmt.Sprintf(" privacy_status = $%d,", paramCounter)
		params = append(params, video.PrivacyStatus)
		paramCounter++
	}

	query += fmt.Sprintf(" updated_at = $%d", paramCounter)
	params = append(params, time.Now()) // Bind the current time
	paramCounter++

	// Remove trailing comma and add WHERE clause
	query = strings.TrimSuffix(query, ",") + fmt.Sprintf(" WHERE id = $%d", paramCounter)
	params = append(params, videoID)

	// Execute the query
	result, err := db.ExecContext(ctx, query, params...)
	if err != nil {
		return nil, fmt.Errorf("error updating video: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("error getting rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return nil, ErrNotFound
	}

	// Update editors assigned to the video
	// Get the current editors assigned to the video
	currentEditors, err := db.GetEditorsByVideo(videoID)
	if err != nil {
		return nil, fmt.Errorf("error fetching current editors: %w", err)
	}

	// Compare new editors with current editors
	for _, newEditor := range video.Editors {
		// Check if the new editor is already in the current editors list
		found := false
		for _, currentEditor := range currentEditors {
			if newEditor.ID == currentEditor.ID {
				found = true
				break
			}
		}

		// If the new editor is not found in the current list, add it
		if !found {
			_, err := db.AddEditorToVideo(videoID, newEditor.ID)
			if err != nil {
				return nil, fmt.Errorf("error assigning editor to video: %w", err)
			}
		}
	}

	// Fetch the updated video before returning
	updatedVideo, err := db.GetVideoByID(videoID)
	if err != nil {
		return nil, fmt.Errorf("error fetching video: %w", err)
	}
	return updatedVideo, nil
}

// DeleteVideo deletes an existing video
func (db *DB) DeleteVideo(videoID string) error {
	ctx := context.Background()
	result, err := db.ExecContext(ctx, "DELETE FROM videos WHERE id = $1", videoID)
	if err != nil {
		return fmt.Errorf("error deleting video: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// GetIterationByID retrieves an iteration by ID
func (db *DB) GetIterationByID(iterationID string) (*models.Iteration, error) {
	var iteration models.Iteration
	err := db.QueryRowContext(context.Background(), "SELECT * FROM iterations WHERE id = $1", iterationID).Scan(
		&iteration.ID,
		&iteration.Video.ID,
		&iteration.URL,
		&iteration.Length,
		&iteration.Status,
		&iteration.Notes,
		&iteration.CreatedAt,
		&iteration.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error fetching iteration: %w", err)
	}

	// Fetch the channel data using the channel ID
	video, err := db.GetVideoByID(iteration.Video.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching video: %w", err)
	}

	// Assign the fetched channel to the video
	iteration.Video = *video // Dereference the channel pointer

	return &iteration, nil
}

// GetIterations retrieves all iterations
func (db *DB) GetIterations() ([]models.Iteration, error) {
	var iterations []models.Iteration
	rows, err := db.QueryContext(context.Background(), "SELECT * FROM iterations")
	if err != nil {
		return nil, fmt.Errorf("error fetching iterations: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var iteration models.Iteration
		if err := rows.Scan(
			&iteration.ID,
			&iteration.Video.ID,
			&iteration.URL,
			&iteration.Length,
			&iteration.Status,
			&iteration.Notes,
			&iteration.CreatedAt,
			&iteration.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning iteration: %w", err)
		}

		// Fetch the channel data using the channel ID
		video, err := db.GetVideoByID(iteration.Video.ID)
		if err != nil {
			return nil, fmt.Errorf("error fetching video: %w", err)
		}

		// Assign the fetched channel to the video
		iteration.Video = *video // Dereference the channel pointer

		iterations = append(iterations, iteration)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through rows: %w", err)
	}

	return iterations, nil
}

// GetIterationsByVideo retrieves iterations by video ID
func (db *DB) GetIterationsByVideo(videoID string) ([]models.Iteration, error) {
	var iterations []models.Iteration
	rows, err := db.QueryContext(context.Background(), "SELECT * FROM iterations WHERE video_id = $1", videoID)
	if err != nil {
		return nil, fmt.Errorf("error fetching iterations: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var iteration models.Iteration
		if err := rows.Scan(
			&iteration.ID,
			&iteration.Video.ID,
			&iteration.URL,
			&iteration.Length,
			&iteration.Status,
			&iteration.Notes,
			&iteration.CreatedAt,
			&iteration.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning iteration: %w", err)
		}

		// Fetch the channel data using the channel ID
		video, err := db.GetVideoByID(iteration.Video.ID)
		if err != nil {
			return nil, fmt.Errorf("error fetching video: %w", err)
		}

		// Assign the fetched channel to the video
		iteration.Video = *video // Dereference the channel pointer

		iterations = append(iterations, iteration)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through rows: %w", err)
	}

	return iterations, nil
}

// CreateIteration creates a new iteration
func (db *DB) CreateIteration(iteration *models.Iteration) (*models.Iteration, error) {
	ctx := context.Background()
	result, err := db.ExecContext(ctx, "INSERT INTO iterations (video_id, url, length, status, notes) VALUES ($1, $2, $3, $4, $5)",
		iteration.Video.ID, iteration.URL, iteration.Length, iteration.Status, iteration.Notes)
	if err != nil {
		return nil, fmt.Errorf("error creating iteration: %w", err)
	}

	iterationID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert ID: %w", err)
	}
	iteration.ID = fmt.Sprintf("%d", iterationID)

	return iteration, nil
}

// UpdateIteration updates an existing iteration
func (db *DB) UpdateIteration(iterationID string, iteration *models.Iteration) (*models.Iteration, error) {
	ctx := context.Background()
	result, err := db.ExecContext(ctx, "UPDATE iterations SET video_id = $1, url = $2, length = $3, status = $4, notes = $5, updated_at = NOW() WHERE id = $6",
		iteration.Video.ID, iteration.URL, iteration.Length, iteration.Status, iteration.Notes, iterationID)
	if err != nil {
		return nil, fmt.Errorf("error updating iteration: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("error getting rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return nil, ErrNotFound
	}

	return iteration, nil
}

// DeleteIteration deletes an existing iteration
func (db *DB) DeleteIteration(iterationID string) error {
	ctx := context.Background()
	result, err := db.ExecContext(ctx, "DELETE FROM iterations WHERE id = $1", iterationID)
	if err != nil {
		return fmt.Errorf("error deleting iteration: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// AddNoteToIteration adds a note to an iteration
func (db *DB) AddNoteToIteration(iterationID string, note *models.Note) error {
	ctx := context.Background()
	_, err := db.ExecContext(ctx, "UPDATE iterations SET notes = $1 WHERE id = $2", note.Content, iterationID)
	if err != nil {
		return fmt.Errorf("error adding note to iteration: %w", err)
	}
	return nil
}

// GetEditorByID retrieves an editor by ID
func (db *DB) GetEditorByID(editorID string) (*models.Editor, error) {
	var editor models.Editor
	err := db.QueryRowContext(context.Background(), "SELECT * FROM editors WHERE id = $1", editorID).Scan(
		&editor.ID,
		&editor.Username,
		&editor.Email,
		&editor.Password,
		&editor.CreatedAt,
		&editor.UpdatedAt,
		&editor.Tier,
		&editor.Trial,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error fetching editor: %w", err)
	}

	return &editor, nil
}

// GetEditors retrieves all editors
func (db *DB) GetEditors() ([]models.Editor, error) {
	var editors []models.Editor
	rows, err := db.QueryContext(context.Background(), "SELECT * FROM editors")
	if err != nil {
		return nil, fmt.Errorf("error fetching editors: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var editor models.Editor
		if err := rows.Scan(
			&editor.ID,
			&editor.Username,
			&editor.Email,
			&editor.Password,
			&editor.CreatedAt,
			&editor.UpdatedAt,
			&editor.Tier,
			&editor.Trial,
		); err != nil {
			return nil, fmt.Errorf("error scanning editor: %w", err)
		}

		editors = append(editors, editor)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through rows: %w", err)
	}

	return editors, nil
}

// CreateEditor creates a new editor
func (db *DB) CreateEditor(editor *models.Editor) (*models.Editor, error) {
	ctx := context.Background()
	result, err := db.ExecContext(ctx, "INSERT INTO editors (username, email, password, tier, trial) VALUES ($1, $2, $3, $4, $5)",
		editor.Username, editor.Email, editor.Password, editor.Tier, editor.Trial)
	if err != nil {
		return nil, fmt.Errorf("error creating editor: %w", err)
	}

	editorID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert ID: %w", err)
	}
	editor.ID = fmt.Sprintf("%d", editorID)

	return editor, nil
}

// UpdateEditor updates an existing editor
func (db *DB) UpdateEditor(editorID string, editor *models.Editor) (*models.Editor, error) {
	ctx := context.Background()
	result, err := db.ExecContext(ctx, "UPDATE editors SET username = $1, email = $2, password = $3, tier = $4, trial = $5, updated_at = NOW() WHERE id = $6",
		editor.Username, editor.Email, editor.Password, editor.Tier, editor.Trial, editorID)
	if err != nil {
		return nil, fmt.Errorf("error updating editor: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("error getting rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return nil, ErrNotFound
	}

	return editor, nil
}

// DeleteEditor deletes an existing editor
func (db *DB) DeleteEditor(editorID string) error {
	ctx := context.Background()
	result, err := db.ExecContext(ctx, "DELETE FROM editors WHERE id = $1", editorID)
	if err != nil {
		return fmt.Errorf("error deleting editor: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// AddEditorToVideo assigns an editor to a video
func (db *DB) AddEditorToVideo(videoID string, editorID string) (sql.Result, error) {
	ctx := context.Background()
	return db.ExecContext(ctx, "INSERT INTO video_editor (video_id, editor_id) VALUES ($1, $2)", videoID, editorID)
}

// RemoveEditorsFromVideo removes all editors from a video
func (db *DB) RemoveEditorsFromVideo(videoID string) error {
	ctx := context.Background()
	_, err := db.ExecContext(ctx, "DELETE FROM video_editor WHERE video_id = $1", videoID)
	if err != nil {
		return fmt.Errorf("error removing editors from video: %w", err)
	}
	return nil
}

// GetEditorsByVideo retrieves editors assigned to a video
func (db *DB) GetEditorsByVideo(videoID string) ([]models.Editor, error) {
	var editors []models.Editor
	rows, err := db.QueryContext(context.Background(), "SELECT e.* FROM video_editor ve JOIN editors e ON ve.editor_id = e.id WHERE ve.video_id = $1", videoID)
	if err != nil {
		return nil, fmt.Errorf("error fetching editors: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var editor models.Editor
		if err := rows.Scan(
			&editor.ID,
			&editor.Username,
			&editor.Email,
			&editor.Password,
			&editor.CreatedAt,
			&editor.UpdatedAt,
			&editor.Tier,
			&editor.Trial,
		); err != nil {
			return nil, fmt.Errorf("error scanning editor: %w", err)
		}

		editors = append(editors, editor)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through rows: %w", err)
	}

	return editors, nil
}

// GetOwner retrieves the owner of a video, channel, or iteration
func (db *DB) GetOwner(userID string) (*models.User, error) {
	var user models.User
	err := db.QueryRowContext(context.Background(), "SELECT * FROM users WHERE id = $1", userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Tier,
		&user.Trial,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error fetching user: %w", err)
	}

	return &user, nil
}
