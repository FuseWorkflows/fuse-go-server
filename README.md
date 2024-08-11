## Fuse

This project is a Golang backend for a platform designed to help YouTubers manage their video editing workflows. It allows YouTubers to:

- **Manage Multiple Channels:** Add multiple YouTube channels with their API keys.
- **Draft Videos:** Draft multiple videos within each channel.
- **Handle Iterations:** Manage multiple iterations of each video, with editors uploading different versions.
- **Get AI Suggestions:** Obtain AI-powered suggestions for video titles, descriptions, chapters, thumbnails, and keywords.
- **Provide Feedback:** Leave notes and feedback on each iteration.
- **Automatic Upload:** Approve iterations for automatic upload to YouTube with AI-suggested metadata (editable by the YouTuber).

### Technologies Used

- **Golang:** The backend is built using the Go programming language.
- **Chi Router:** Used for handling HTTP requests and routing.
- **PostgreSQL:** Database for storing user, channel, video, iteration, and editor data.
- **JWT Authentication:** Implemented for secure user authentication and authorization.
- **YouTube API:** Integrated to automatically upload videos to YouTube channels.
- **AI Service:** Connects to an AI service (placeholder in the code) for generating metadata suggestions.

### Project Structure

```
├── main.go
├── routes
│   └── routes.go
├── middleware
│   └── auth.go
├── handlers
│   ├── channel.go
│   ├── editor.go
│   ├── iteration.go
│   ├── user.go
│   ├── video.go
│   └── auth.go
├── models
│   ├── ai_suggestions.go
│   ├── channel.go
│   ├── editor.go
│   ├── iteration.go
│   ├── user.go
│   ├── video.go
│   └── note.go
├── database
│   ├── migrations
│   │   ├── 1_init.up.sql
│   │   ├── 1_init.down.sql
│   │   ├── 2_create_users_table.up.sql
│   │   ├── 2_create_users_table.down.sql
│   │   ├── 3_create_channels_table.up.sql
│   │   ├── 3_create_channels_table.down.sql
│   │   ├── 4_create_videos_table.up.sql
│   │   ├── 4_create_videos_table.down.sql
│   │   ├── 5_create_iterations_table.up.sql
│   │   ├── 5_create_iterations_table.down.sql
│   │   ├── 6_create_editors_table.up.sql
│   │   ├── 6_create_editors_table.down.sql
│   │   └── 7_create_video_editor_table.up.sql
│   └── db.go
├── utils
│   ├── ai.go
│   └── youtube.go
├── config
│   └── config.go
├── .env
└── .gitignore
```

### Getting Started

1. **Prerequisites:**

   - Go installed on your system (version 1.11 or later).
   - PostgreSQL installed and running.
   - A YouTube Data API key.

2. **Clone the Repository:**

   ```bash
   git clone https://github.com/FuseWorkflows/fuse-go-server.git
   cd fuse-go-server
   ```

3. **Create a `.env` file:**

   - Copy the provided .env.template file:
     ```bash
     cp .env.template .env
     ```

   - Open the .env file and replace the placeholder values with your actual configuration:
     ```
     DB_HOST=your_database_host
     DB_PORT=your_database_port
     DB_USER=your_database_user
     DB_PASSWORD=your_database_password
     DB_NAME=your_database_name
     JWT_KEY=your_secret_jwt_key
     PORT=8080
     AI_SERVICE=http://your_ai_service_url
     YOUTUBE_API_KEY=your_youtube_api_key
     ```

4. **Run Database Migrations:**

   - **(Replace with your chosen migration tool)**
   - Use a migration tool (e.g., `migrate`, `goose`) to create and apply the database migrations. For example, with `migrate`:
     ```bash
     migrate -database "postgres://your_database_user:your_database_password@your_database_host:your_database_port/your_database_name?sslmode=disable" -path database/migrations up
     ```

5. **Build and Run the Server:**
   ```bash
   go run main.go
   ```
   The server should start listening on the specified port (default: 8080).

### API Endpoints

This project provides various API endpoints. For detailed documentation of the endpoints, please refer to the comments within the `handlers` and `routes` packages.

### AI Service

The AI service used in this project is a placeholder. You need to replace the `AI_SERVICE` environment variable in the `.env` file with the actual URL of your chosen AI service that provides metadata suggestions.

### Contributions

Contributions are welcome! Please submit pull requests for any improvements or bug fixes.

<!-- ### License -->

<!-- This project is licensed under the [MIT License](LICENSE). -->
