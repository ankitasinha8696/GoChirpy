# GoChirpy
HTTP server in Go

This project implements a simple HTTP server in Go with several endpoints and functionalities:

Readiness Handler:

Endpoint: GET /api/healthz
Responds with a status code 200 and a body text "OK" to indicate the server is ready.
Metrics and Reset Handlers:

Metrics Handler:
Endpoint: GET /admin/metrics
Displays the number of times the file server has been accessed.
Reset Handler:
Endpoint: POST /admin/reset
Resets the file server access counter to 0.
File Server:

Serves files from the current directory at the /app/ path.
Includes middleware to increment a counter each time a file is accessed.
Chirp Validation:

Endpoint: POST /api/validate_chirp
Accepts a JSON payload with a body field.
Validates the chirp to ensure it is not longer than 140 characters.
Cleanses the chirp body by replacing certain profane words with "****".
Responds with a JSON object indicating whether the chirp is valid and the cleaned chirp body.
