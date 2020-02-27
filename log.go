package gologger

import (
	"fmt"
	"sync"
	"time"
)

// Level defines log levels.
type Level uint8

// Debug levels
const (
	Disabled Level = iota

	PanicLevel

	FatalLevel

	ErrorLevel

	WarnLevel

	InfoLevel

	DebugLevel
)

type Logger struct {
	level        Level
	namedActions map[string]*protectAct
	levelActions map[Level][]string
	mux          sync.Mutex
	wg           sync.WaitGroup
}

type LogMsg struct {
	Level      Level
	Time       time.Time
	Parameters interface{}
}

var instance *Logger

type protectAct struct {
	mux      *sync.Mutex // защищает запуск функции
	function LogAction
}

type LogAction func(logMsg LogMsg)

var mu sync.Mutex

// Init Инициализует логгер и возвращает ссылку на объект
func Init(level Level) *Logger {
	mu.Lock()
	defer mu.Unlock()

	if instance == nil {
		instance = new(Logger)
		instance.level = level
		instance.namedActions = make(map[string]*protectAct)
		instance.levelActions = make(map[Level][]string)

	}
	return instance
}

// Warning invoke actions
func (l *Logger) Warning(parameter interface{}) {
	if l.level < WarnLevel {
		return
	}
	l.invokeActions(WarnLevel, parameter)
	return
}

// Debug invoke actions
func (l *Logger) Debug(parameter interface{}) {
	if l.level < DebugLevel {
		return
	}
	l.invokeActions(DebugLevel, parameter)
	return
}

func GetLevelName(level Level) string {
	switch level {
	case DebugLevel:
		return "debug"
	case WarnLevel:
		return "warning"
	default:
		return "unkown"
	}
}

// Вызов действий ассинхронно
func (l *Logger) invokeActions(level Level, parameter interface{}) {
	actions := l.levelActions[level]
	for _, actionName := range actions {
		l.wg.Add(1)
		protectaction := l.namedActions[actionName]
		logMsg := LogMsg{
			Level:      level,
			Time:       time.Now(),
			Parameters: parameter,
		}
		go func(logMsg LogMsg) {
			protectaction.mux.Lock()
			defer protectaction.mux.Unlock()
			defer l.wg.Done()

			protectaction.function(logMsg)
		}(logMsg)
	}
}

// AddAction Добавление действия
func (l *Logger) AddAction(name string, action LogAction) {
	protect := new(protectAct)
	protect.function = action
	protect.mux = new(sync.Mutex)
	l.namedActions[name] = protect
}

func (l *Logger) SetLogLevelActions(level Level, actions []string) error {
	for _, val := range actions {
		if _, ok := l.namedActions[val]; !ok {
			err := fmt.Errorf("Action with name %s doesn't exists", val)
			return err
		}
	}
	l.levelActions[level] = actions
	return nil
}

// Ожидания выполнения логгирования
func (l *Logger) WaitLogsDone() {
	l.wg.Wait()
}
