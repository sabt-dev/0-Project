# Define the paths to the frontend and backend directories
$frontendPath = "..\frontend"
$backendPath = "..\backend\cmd"

# Define the commands to run
$frontendCmd = "npm run dev"
$backendCmd = "go run main.go"

# Start the frontend server in a new Command Prompt window
Start-Process "cmd.exe" -ArgumentList "/c cd $frontendPath && $frontendCmd"

# Start the backend server in a new Command Prompt window
Start-Process "cmd.exe" -ArgumentList "/c cd $backendPath && $backendCmd"