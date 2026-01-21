package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gen2brain/beeep"
)

type TimerState int

const (
	StateStopped TimerState = iota
	StateRunning
	StatePaused
)

type Timer struct {
	duration  time.Duration
	remaining time.Duration
	mode      string
	state     TimerState
	startTime time.Time
	pausedAt  time.Time
}

type Session struct {
	focusCount              int
	totalFocusTime          time.Duration
	breakCount              int
	totalBreakTime          time.Duration
	soundEnabled            bool
	notificationsEnabled    bool
	focusMinutes            int
	shortBreakMinutes       int
	longBreakMinutes        int
	sessionsBeforeLongBreak int
}

func main() {
	focusMinutes := flag.Int("f", 25, "Focus session duration in minutes (default: 25)")
	breakMinutes := flag.Int("b", 5, "Short break duration in minutes (default: 5)")
	longBreakMinutes := flag.Int("l", 15, "Long break duration in minutes (default: 15)")
	sessionsBeforeLongBreak := flag.Int("s", 4, "Number of focus sessions before long break (default: 4)")
	sound := flag.Bool("sound", true, "Enable sound notifications (default: true)")
	notifications := flag.Bool("notify", true, "Enable desktop notifications (default: true)")
	showHelp := flag.Bool("h", false, "Show help message")

	flag.Usage = func() {
		fmt.Printf("gomodoro - CLI Pomodoro Timer\n\n")
		fmt.Printf("Usage: gomodoro [options]\n\n")
		fmt.Printf("Options:\n")
		flag.PrintDefaults()
		fmt.Printf("\nExamples:\n")
		fmt.Printf("  gomodoro                         # Start with default settings\n")
		fmt.Printf("  gomodoro -f 30 -b 10              # 30 min focus, 10 min breaks\n")
		fmt.Printf("  gomodoro -f 20 -b 5 -l 15 -s 3   # Custom long break settings\n")
		fmt.Printf("  gomodoro -sound=false            # Disable sound\n")
		fmt.Printf("  gomodoro -h                      # Show this help\n")
	}

	flag.Parse()

	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	inputChan := make(chan string, 1)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			inputChan <- scanner.Text()
		}
	}()

	timer := Timer{
		duration: time.Duration(*focusMinutes) * time.Minute,
		mode:     "FOCUS",
		state:    StateRunning,
	}

	session := Session{
		soundEnabled:            *sound,
		notificationsEnabled:    *notifications,
		focusMinutes:            *focusMinutes,
		shortBreakMinutes:       *breakMinutes,
		longBreakMinutes:        *longBreakMinutes,
		sessionsBeforeLongBreak: *sessionsBeforeLongBreak,
	}

	go func() {
		<-sigChan
		fmt.Println("\nTimer stopped by user")
		printSessionStats(session)
		os.Exit(0)
	}()

	printHeader(fmt.Sprintf("Starting %s session for %d minutes", timer.mode, session.focusMinutes))
	fmt.Println("Controls: [Enter] to pause/resume, Ctrl+C to quit")
	if session.soundEnabled {
		playBeep()
	}

	runTimer(timer, inputChan)

	session.focusCount++
	session.totalFocusTime += timer.duration

	if session.notificationsEnabled {
		notify("Pomodoro Timer", fmt.Sprintf("%s session completed!", timer.mode))
	}

	printHeader(fmt.Sprintf("%s session completed! Time for a %d minute break", timer.mode, *breakMinutes))
	printSessionStats(session)

	timer.mode = "BREAK"

	if session.focusCount%session.sessionsBeforeLongBreak == 0 {
		timer.duration = time.Duration(session.longBreakMinutes) * time.Minute
		fmt.Printf("🎉 Long break this time!\n")
	} else {
		timer.duration = time.Duration(session.shortBreakMinutes) * time.Minute
	}

	timer.state = StateRunning

	breakDuration := int(timer.duration.Minutes())
	printHeader(fmt.Sprintf("Starting %s session for %d minutes", timer.mode, breakDuration))
	if session.soundEnabled {
		playBeep()
	}
	runTimer(timer, inputChan)

	session.breakCount++
	session.totalBreakTime += timer.duration

	if session.notificationsEnabled {
		notify("Pomodoro Timer", fmt.Sprintf("%s session completed! Great job!", timer.mode))
	}
	printHeader(fmt.Sprintf("%s session completed! Great job!", timer.mode))
	printSessionStats(session)
}

func printHeader(message string) {
	fmt.Printf("\n%s\n", message)
	fmt.Println(strings.Repeat("=", len(message)))
}

func notify(title, message string) {
	err := beeep.Notify(title, message, "")
	if err != nil {
		fmt.Printf("\a")
	}
}

func playBeep() {
	beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
}

func printSessionStats(session Session) {
	fmt.Printf("\n📊 Session Statistics:\n")
	fmt.Printf("  Focus Sessions: %d\n", session.focusCount)
	fmt.Printf("  Total Focus Time: %v\n", session.totalFocusTime)
	fmt.Printf("  Break Sessions: %d\n", session.breakCount)
	fmt.Printf("  Total Break Time: %v\n", session.totalBreakTime)
	fmt.Println(strings.Repeat("-", 40))
}

func runTimer(timer Timer, inputChan chan string) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	timer.remaining = timer.duration
	timer.startTime = time.Now()
	endTime := timer.startTime.Add(timer.duration)

	for timer.remaining > 0 {
		select {
		case <-ticker.C:
			if timer.state == StateRunning {
				timer.remaining = endTime.Sub(time.Now())
				if timer.remaining < 0 {
					timer.remaining = 0
				}

				elapsed := timer.duration - timer.remaining
				progress := int(math.Floor(float64(elapsed) / float64(timer.duration) * 20))
				bar := strings.Repeat("█", progress) + strings.Repeat("░", 20-progress)

				minutes := int(timer.remaining.Minutes())
				seconds := int(timer.remaining.Seconds()) % 60

				status := "▶"
				if timer.state == StatePaused {
					status = "⏸"
				}

				fmt.Printf("\r%s [%s] %02d:%02d [%-20s]", status, timer.mode, minutes, seconds, bar)
			}
		case input := <-inputChan:
			if input == "" {
				if timer.state == StateRunning {
					timer.state = StatePaused
					timer.pausedAt = time.Now()
				} else if timer.state == StatePaused {
					pauseDuration := time.Since(timer.pausedAt)
					endTime = endTime.Add(pauseDuration)
					timer.state = StateRunning
				}
			}
		}
	}
}
