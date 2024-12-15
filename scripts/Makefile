# Define the backend and frontend start commands
BACKEND_CMD = go run ..\backend\cmd\main.go
FRONTEND_CMD = npm run dev --prefix ..\frontend

# Target to start the backend server
backend:
    @echo "Starting backend server..."
    $(BACKEND_CMD)

# Target to start the frontend server
frontend:
    @echo "Starting frontend server..."
    $(FRONTEND_CMD)

# Target to start both servers concurrently
.PHONY: all
all: backend frontend
# Target to start both servers concurrently using a tool like concurrently
start:
    @echo "Starting both servers..."
    concurrently "$(BACKEND_CMD)" "$(FRONTEND_CMD)"
