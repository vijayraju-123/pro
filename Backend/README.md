# Task Management System Backend

This is the backend for an AI-powered task management system built with Go (Golang) using the Gin framework. It includes:

- User authentication (JWT-based)
- Task creation, assignment, and tracking
- AI-powered task suggestions (OpenAI/Gemini API)
- Real-time updates via WebSockets
- Deployment on cloud platforms (e.g., Render, Fly.io)

## Tech Stack
- **Backend**: Go (Golang) with Gin/Fiber
- **Database**: PostgreSQL or MongoDB
- **Authentication**: JWT
- **Real-time**: WebSockets
- **AI**: OpenAI/Gemini API

## Setup
1. Install Go 1.21 or higher.
2. Run `go mod tidy` to install dependencies.
3. Set environment variables (e.g., `DB_USER`, `DB_PASSWORD`, `JWT_SECRET`, `OPENAI_KEY`).
4. Run `go run main.go` to start the server.

## Deployment
Deploy on Render, Fly.io, or similar platforms. Ensure environment variables are configured.

## Contributing
Feel free to fork this repository, make improvements, and submit pull requests!