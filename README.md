# Task Manager CLI in Go 📟

`go-task` is a lightweight command‑line task manager written in Go and backed by SQLite. It lets you capture, organize, and track your tasks with priorities, tags, and completion status directly from your terminal.

## What you can do

- Add tasks with title, description, priority, and tags.
- View tasks in a readable table, sorted by priority.
- Update task details as your work evolves.
- Mark tasks as completed to track progress.
- Delete tasks you no longer need.
- Search tasks quickly using keywords.

## Getting Started

### Prerequisites

- Go 1.21.2 or later
- SQLite

### Installation

1. **Clone the repository**:

   ```sh
   git clone https://github.com/dulanhewage/go-task.git
   cd go-task
   ```

2. **Install dependencies**:

   ```sh
   go mod tidy
   ```

3. **Build the project**:

   ```sh
   go build -o go-task
   ```

### Configuration

You can configure the database file path by creating a `config.json` file in the project root:

```json
{
  "db_file": "my_tasks.db"
}
```

If `config.json` is missing or `db_file` is not specified, it defaults to `tasks.db`.

### Run tests

Run all tests:

```sh
go test ./cmd/...
```

Run a specific test file:

```sh
go test ./cmd/add_test.go
```

Run tests with verbose output:

```sh
go test -v ./cmd/...
```

### Usage

After building the project, you can use the CLI tool by running the generated binary:

```sh
./go-task --help
```

### Examples

Here is an example of how to use the Task Manager CLI:

1. **Add a new task (Simple)**:

```sh
./go-task add --title "Buy groceries" --description "Milk, Bread, Eggs"
```

2. **Add a task with Priority and Tags**:

```sh
./go-task add --title "Fix server bug" --description "Crash on startup" --priority High --tags "work,urgent,dev"
```

_Note: Priorities can be `High` (H), `Medium` (M), or `Low` (L). Tags are comma-separated._

3. **List all tasks**:

```sh
./go-task list
```

_Displays tasks in a table format, sorted by Priority (High -> Low)._

4. **Update a task**:

```sh
./go-task update --id 1 --title "Buy groceries and fruits" --description "Milk, Bread, Eggs, Apples"
```

5. **Mark a task as completed**:

```sh
./go-task complete --id 1
```

6. **Search tasks by keyword**:

```sh
./go-task search --keyword "groceries"
```

7. **Delete a task**:

```sh
./go-task delete --id 1
```
