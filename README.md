# Terminal Emulator (`gvte`)

A modular cross-platform terminal emulator written in Go.

## Project Structure

```text
my-terminal/
├── cmd/
│   └── terminal/
│       └── main.go           # Entry point: configuration, PTY, engine, and UI assembly
├── internal/
│   ├── pty/                  # OS Pseudo-terminal (PTY) interface
│   │   ├── pty.go            # PTY interface and Winsize definition
│   │   ├── pty_unix.go       # Termios / Master-Slave PTY for Linux/macOS
│   │   └── pty_windows.go    # ConPTY initialization for Windows
│   ├── emulator/             # Terminal Engine & State
│   │   ├── parser/           # ANSI / VT100 / xterm escape sequence parser
│   │   ├── grid/             # 2D character matrix (Cell, Line, Scrollback buffer)
│   │   ├── selection/        # Text selection & system clipboard logic
│   │   └── state.go          # State machine (cursor, scroll mode, color palette)
│   ├── ui/                   # Rendering & Application Window (Presentation)
│   │   ├── font/             # Font loading, rasterization & glyph caching
│   │   ├── renderer/         # Character grid rendering (OpenGL/Vulkan/Soft-render)
│   │   ├── input/            # Keypress & mouse mapping to ANSI escape sequences
│   │   └── window.go         # App window lifecycle (creation, resize, context)
│   └── config/               # Configuration parser (color schemes, fonts, hotkeys)
├── assets/                   # Default embedded fonts, window icons
├── go.mod
└── README.md
```

## Packages Breakdown

- **`cmd/terminal/main.go`**: Initializes configuration, spawns PTY, sets up emulator state, and starts the UI window.
- **`internal/pty`**: OS-level pseudo-terminal abstractions (`pty_unix.go` using termios/master-slave and `pty_windows.go` using ConPTY).
- **`internal/emulator`**:
  - `parser/`: Parses incoming VT100 / ANSI escape sequences.
  - `grid/`: Maintains character matrix, lines, attributes, and scrollback history.
  - `selection/`: Manages text selection and clipboard operations.
  - `state.go`: Core state machine holding cursor state, current screen buffer, and titles.
- **`internal/ui`**:
  - `font/`: Handles font rendering, metrics, and glyph rasterization cache.
  - `renderer/`: Draws grid state, cursor, and selections to screen canvas.
  - `input/`: Maps user input into ANSI sequences sent to PTY.
  - `window.go`: Manages main application window loop.
- **`internal/config`**: Reads and manages configuration options.

## Building and Running

```bash
# Run the application
go run ./cmd/terminal

# Build binary
go build -o bin/terminal ./cmd/terminal
```
