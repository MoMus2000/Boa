# Makefile for Python Project

# Define variables
PYTHON = python3
PIP = pip
VENV = venv
VENV_DIR = .env

# Default target: create virtual environment and install dependencies
all: install

# Create virtual environment
$(VENV_DIR)/bin/activate: requirements.txt
	$(PYTHON) -m venv $(VENV_DIR)
	$(VENV_DIR)/bin/pip install --upgrade pip
	$(VENV_DIR)/bin/pip install -r requirements.txt

# Install dependencies from requirements.txt
install: $(VENV_DIR)/bin/activate

# Run tests (assuming you have pytest installed)
test:
	@echo "Running tests..."
	./test.sh

# Run a specific Python script
run:
	$(VENV_DIR)/bin/python main.py

# Clean up virtual environment
clean:
	@echo "Cleaning up..."
	rm -rf __pycache__

# Check the installed dependencies
freeze:
	@echo "Freezing dependencies..."
	$(VENV_DIR)/bin/pip freeze > requirements.txt
