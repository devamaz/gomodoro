package main

import (
	"testing"
	"time"
)

func TestTimerStateTransitions(t *testing.T) {
	timer := Timer{
		duration: 25 * time.Minute,
		mode:     "FOCUS",
		state:    StateRunning,
	}

	if timer.state != StateRunning {
		t.Errorf("Expected initial state to be Running, got %v", timer.state)
	}

	timer.state = StatePaused
	if timer.state != StatePaused {
		t.Errorf("Expected state to be Paused, got %v", timer.state)
	}

	timer.state = StateStopped
	if timer.state != StateStopped {
		t.Errorf("Expected state to be Stopped, got %v", timer.state)
	}
}

func TestSessionInitialization(t *testing.T) {
	session := Session{
		soundEnabled:            true,
		notificationsEnabled:    true,
		focusMinutes:            25,
		shortBreakMinutes:       5,
		longBreakMinutes:        15,
		sessionsBeforeLongBreak: 4,
	}

	if session.focusMinutes != 25 {
		t.Errorf("Expected focusMinutes to be 25, got %d", session.focusMinutes)
	}

	if session.shortBreakMinutes != 5 {
		t.Errorf("Expected shortBreakMinutes to be 5, got %d", session.shortBreakMinutes)
	}

	if session.longBreakMinutes != 15 {
		t.Errorf("Expected longBreakMinutes to be 15, got %d", session.longBreakMinutes)
	}

	if session.sessionsBeforeLongBreak != 4 {
		t.Errorf("Expected sessionsBeforeLongBreak to be 4, got %d", session.sessionsBeforeLongBreak)
	}
}

func TestSessionCounting(t *testing.T) {
	session := Session{}

	session.focusCount = 3
	session.breakCount = 3

	if session.focusCount != 3 {
		t.Errorf("Expected focusCount to be 3, got %d", session.focusCount)
	}

	if session.breakCount != 3 {
		t.Errorf("Expected breakCount to be 3, got %d", session.breakCount)
	}

	session.focusCount++
	if session.focusCount != 4 {
		t.Errorf("Expected focusCount to be 4, got %d", session.focusCount)
	}
}

func TestLongBreakCalculation(t *testing.T) {
	tests := []struct {
		focusCount              int
		sessionsBeforeLongBreak int
		expectedLongBreak       bool
	}{
		{1, 4, false},
		{2, 4, false},
		{3, 4, false},
		{4, 4, true},
		{8, 4, true},
		{3, 3, true},
	}

	for _, test := range tests {
		result := test.focusCount%test.sessionsBeforeLongBreak == 0
		if result != test.expectedLongBreak {
			t.Errorf("focusCount=%d, sessionsBeforeLongBreak=%d: expected long break %v, got %v",
				test.focusCount, test.sessionsBeforeLongBreak, test.expectedLongBreak, result)
		}
	}
}

func TestTimerDuration(t *testing.T) {
	duration := 25 * time.Minute
	timer := Timer{
		duration:  duration,
		remaining: duration,
	}

	if timer.duration != duration {
		t.Errorf("Expected duration to be %v, got %v", duration, timer.duration)
	}

	if timer.remaining != duration {
		t.Errorf("Expected remaining to be %v, got %v", duration, timer.remaining)
	}
}

func TestTotalFocusTime(t *testing.T) {
	session := Session{}
	focusDuration := 25 * time.Minute

	session.focusCount = 3
	session.totalFocusTime = focusDuration * time.Duration(session.focusCount)

	expectedTotal := 75 * time.Minute
	if session.totalFocusTime != expectedTotal {
		t.Errorf("Expected totalFocusTime to be %v, got %v", expectedTotal, session.totalFocusTime)
	}
}

func TestTotalBreakTime(t *testing.T) {
	session := Session{}
	breakDuration := 5 * time.Minute

	session.breakCount = 2
	session.totalBreakTime = breakDuration * time.Duration(session.breakCount)

	expectedTotal := 10 * time.Minute
	if session.totalBreakTime != expectedTotal {
		t.Errorf("Expected totalBreakTime to be %v, got %v", expectedTotal, session.totalBreakTime)
	}
}

func TestSoundNotificationSettings(t *testing.T) {
	session := Session{
		soundEnabled:         false,
		notificationsEnabled: false,
	}

	if session.soundEnabled {
		t.Error("Expected soundEnabled to be false")
	}

	if session.notificationsEnabled {
		t.Error("Expected notificationsEnabled to be false")
	}

	session.soundEnabled = true
	session.notificationsEnabled = true

	if !session.soundEnabled {
		t.Error("Expected soundEnabled to be true")
	}

	if !session.notificationsEnabled {
		t.Error("Expected notificationsEnabled to be true")
	}
}
