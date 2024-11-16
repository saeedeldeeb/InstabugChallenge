#!/bin/sh

# Generate application key
./main artisan key:generate

# Generate JWT secret key
./main artisan jwt:secret

# Migrate the database
./main artisan migrate

# Start the application
exec ./main 