package pkg

import (
	"net"
	"sync"
)

type Server struct {
	Server         net.Listener        // connections
	Connections    map[net.Conn]string // names
	UsedNames      map[string]bool     // names is used or not
	MaxConnections int                 // to be changed on server limits
	AllMessages    []string            // history of messages
	mutex          sync.Mutex          // the mutex
}

const (
	// WelcomMessage =
	WelcomMessage = "Welcome to TCP-Chat!\n         _nnnn_\n        dGGGGMMb\n       @p~qp~~qMb\n       M|@||@) M|\n       @,----.JM|\n      JS^\\__/  qKL\n     dZP        qKRb\n    dZP          qKKb\n   fZP            SMMb\n   HZM            MMMM\n   FqM            MMMM\n __| \".        |\\dS\"qML\n |    `.       | `' \\Zq\n_)      \\.___.,|     .'\n\\____   )MMMMMP|   .'\n     `-'       `--'\n[ENTER YOUR NAME]: "
	// [time][name]:
	PatternSending = "[%v][%v]:"
	// [time][name]: message
	PatternMessage  = "[%v][%s]:%s"
	PatternJoinChat = "%s has joined our chat...\n"
	PatternLeftChat = "%s has left our chat...\n"
)

// MessageModes
const (
	// iota from 1 in sequence
	ModeJoinChat = iota
	ModeSendMessage
	ModeLeftChat
)

// TimePatterns
const (
	TimeDefault = "2006-01-02 15:04:05"
)

// Color Patterns
const (
	ColorReset  = "\u001b[0m"
	ColorYellow = "\u001b[33m"

	BgColorRed  = "\u001b[41m"
	BgColorGray = "\u001b[47;1m"
)
