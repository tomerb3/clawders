# Project Guidelines & Workflow

## Directory Structure
- **Python Scripts:** All Python logic and source files MUST be placed in the `src/` directory.
- **Automation:** The root directory should contain a `run.sh` script for container management.
- **Docker:** The root directory should contain a `Dockerfile` for containerization. the scripts and also the tests.

## Build & Execution (run.sh)
Every task must include or update `run.sh` with the following functionality:
- **Build:** `docker build -t project-image .`
- **Run:** `docker run --rm project-image`
- Use this script as the primary entry point for testing the containerized environment.
- no pip install in the linux. only inside the docker container.

## Task Execution Protocol
Before writing any project code, you must follow this mandatory sequence:

### 1. Analysis & Strategy
Analyze the task and define the measurement for success. Do not write implementation code until the testing strategy is defined.

### 2. Testing Specification
Determine the testing type based on the component:
- **Backend Tasks:**
    - Must include Python unit tests (using `pytest` or `unittest`).
    - Focus on API contracts, data processing logic, and edge cases.
    - Integration tests to ensure `src/` modules interact correctly.
- **UI/Frontend Tasks:**
    - Must include functional tests or end-to-end (E2E) scripts.
    - Focus on user flow validation and visual consistency.
    - Use mock data where applicable to isolate the UI.

### 3. Verification Plan
Explicitly state:
- **How** we are going to measure success (e.g., "All tests pass with >80% coverage").
- **How** we are going to test it (e.g., "Run `run.sh` to trigger the internal pytest suite").
- **how** if for example the project will create video... how the AI will check the quality if the result video.
          if the result is app website - how the AI will check the quality of the website and buttons...

### 4. Implementation
Only after the above steps are confirmed, proceed to build the project files and implementation logic.

## Commands
- **Test:** `docker run --rm -it docker-pytest:ver1 src/tests/`
- **Container Control:** `./run.sh`
