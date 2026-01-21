# gomodoro

A simple, elegant CLI Pomodoro timer written in Go. Boost your productivity with customizable focus sessions and breaks, complete with visual progress tracking and desktop notifications.

## Features

- **Customizable Sessions**: Configure focus duration, short breaks, long breaks, and sessions before long break
- **Visual Progress Bar**: Real-time progress visualization during sessions
- **Pause/Resume**: Press Enter to pause and resume sessions
- **Desktop Notifications**: Get notified when sessions complete (desktop notifications)
- **Sound Alerts**: Audio beeps to signal session start and completion
- **Session Statistics**: Track your productivity with session counts and total time
- **Graceful Shutdown**: Clean exit with Ctrl+C + session statistics summary

## Installation

### Build from source

```bash
git clone https://github.com/devamaz/gomodoro.git
cd gomodoro
go build -o gomodoro
```

### Using Go install

```bash
go install github.com/devamaz/gomodoro@latest
```

## Usage

### Basic usage

```bash
gomodoro
```

This starts a Pomodoro session with default settings:
- Focus: 25 minutes
- Short break: 5 minutes
- Long break: 15 minutes
- Long break every 4 focus sessions

### Custom session settings

```bash
gomodoro -f 30 -b 10 -l 15 -s 3
```

This configures:
- 30 minute focus sessions
- 10 minute short breaks
- 15 minute long breaks
- Long break every 3 focus sessions

### Disable notifications or sound

```bash
gomodoro -sound=false
gomodoro -notify=false
```

### Show help

```bash
gomodoro -h
```

## Command-line Options

| Flag | Description | Default |
|------|-------------|---------|
| `-f` | Focus session duration in minutes | 25 |
| `-b` | Short break duration in minutes | 5 |
| `-l` | Long break duration in minutes | 15 |
| `-s` | Number of focus sessions before long break | 4 |
| `-sound` | Enable sound notifications (true/false) | true |
| `-notify` | Enable desktop notifications (true/false) | true |
| `-h` | Show help message | - |

## Controls

- **Enter** - Pause/Resume current session
- **Ctrl+C** - Exit and show session statistics

## Examples

Start with default settings:
```bash
gomodoro
```

Customize for longer focus periods:
```bash
gomodoro -f 45 -b 10
```

Minimal setup without notifications:
```bash
gomodoro -sound=false -notify=false
```

## Session Statistics

After each session, you'll see a summary of your productivity:

```
📊 Session Statistics:
  Focus Sessions: 1
  Total Focus Time: 25m0s
  Break Sessions: 1
  Total Break Time: 5m0s
----------------------------------------
```

## Development

### Run tests

```bash
go test -v
```

### Build

```bash
go build -o gomodoro
```

## Requirements

- Go 1.21.5 or higher
- Desktop notifications support (depends on OS)

## License

This project is open source and available under the terms specified in the repository.

## Contributing

Contributions are welcome! Feel free to submit issues and pull requests.

## Acknowledgments

Built with ❤️ using Go
