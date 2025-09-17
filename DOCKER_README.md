# Viridian City Bank - Dockerized Full Stack Application

A complete banking web application with Go backend and HTML/CSS/JS frontend, now fully containerized!

## 🐳 Docker Setup (Recommended)

### Quick Start
```bash
# Build and run with one command
./run.sh

# Or manually with docker compose
docker compose up --build

# Or build and run separately
./build.sh
docker run -p 8080:8080 viridian-bank:latest
```

### Access the Application
- **Frontend**: http://localhost:8080
- **API**: http://localhost:8080/api
- **Health Check**: http://localhost:8080/health

### Docker Commands
```bash
# Build the image
docker build -t viridian-bank:latest .

# Run the container
docker run -p 8080:8080 viridian-bank:latest

# Run with docker compose (recommended)
docker compose up -d

# View logs
docker compose logs -f

# Stop the application
docker compose down

# Rebuild and restart
docker compose up --build
```

## 🏗️ Architecture

The application now runs as a single container containing:
- **Go Backend**: Serves API endpoints and static files
- **Frontend**: HTML/CSS/JS served by the Go server
- **Database**: SQLite database included in container

### Container Structure
```
/app/
├── viridian-bank-backend    # Go binary
├── index.html              # Main frontend
├── debug_test.html         # Debug page
├── css/                    # Stylesheets
├── js/                     # JavaScript files
├── assets/                 # Images and assets
└── viridian_bank.db       # SQLite database
```

## 🛠️ Development

### Local Development (without Docker)
```bash
# Backend
cd backend
go run .

# Frontend
# Serve frontend files with any HTTP server
python3 -m http.server 8000
# OR
npx serve .
```

### Environment Variables
The application supports these environment variables:
- `PORT`: Server port (default: 8080)
- `GIN_MODE`: Gin framework mode (release/debug)

## 📂 Project Structure
```
viridian-city-bank/
├── Dockerfile              # Main Docker configuration
├── compose.yaml            # Docker Compose setup
├── build.sh               # Build script
├── run.sh                 # Run script
├── index.html             # Frontend entry point
├── css/                   # Frontend styles
├── js/                    # Frontend JavaScript
├── assets/                # Static assets
├── backend/               # Go backend source
│   ├── main.go
│   ├── models.go
│   ├── *.go
│   └── viridian_bank.db
└── n8n-workflows/         # Automation workflows
```

## 🚀 Features
- Complete banking web interface
- User authentication and registration
- Account management and transfers
- Payment requests
- Card management
- Admin functionality
- Real-time transaction processing
- Responsive design

## 🔧 API Endpoints
- `POST /api/register` - User registration
- `POST /api/login` - User login
- `GET /api/account` - Account information
- `GET /api/balance` - Account balance
- `POST /api/transfer` - Money transfer
- `GET /api/transactions` - Transaction history
- And more...

## 🐛 Troubleshooting

### Common Issues
1. **Port already in use**: Change the port mapping in compose.yaml
2. **Build fails**: Ensure Docker has enough memory allocated
3. **Database issues**: The SQLite database is included in the container

### Logs
```bash
# View application logs
docker-compose logs viridian-bank

# Follow logs in real-time
docker-compose logs -f
```

## 📝 Notes
- The frontend now uses relative API URLs (`/api`) instead of `http://localhost:8088/api`
- All static files are served by the Go backend
- Database persistence can be configured with Docker volumes
- The application runs on port 8080 by default
