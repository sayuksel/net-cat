package pkg

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func (s *Server) Create(port string, maxConn int) error {
	serv, err := net.Listen("tcp", port)
	if err != nil {
		return (err)
	}
	s.Server = serv
	s.MaxConnections = maxConn
	s.Connections = make(map[net.Conn]string, maxConn)
	s.UsedNames = make(map[string]bool, maxConn)
	return nil
}

// add can connect fuction
func (s *Server) AddConnection(conn net.Conn, name string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for name == "" {
		return errors.New("Name cant be empty")
	}
	// add checking for  can connect and  check max connection length PLUS ERROR HANDLING

	if s.MaxConnections != 0 && len(s.Connections) >= s.MaxConnections {
		return fmt.Errorf("The room is full [%v]", conn.RemoteAddr())
	}
	if s.UsedNames[name] {
		return fmt.Errorf("Name '%s' is Exist [%v]", name, conn.RemoteAddr())
	}
	// Add the connection to the map
	s.UsedNames[name] = true
	s.Connections[conn] = name
	return nil
}

// go handle
func (s *Server) Handle(conn net.Conn) {

	welcomeLines, err := ReadFile()
	for _, line := range welcomeLines {
		fmt.Fprint(conn, line)
	}
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	name, err := bufio.NewReader(conn).ReadString('\n')
	//fmt.Fprintf(conn, name)

	if err != nil {
		log.Fatal(err)
	}

	name = strings.Replace(name, "\n", "", 1)

	err = s.AddConnection(conn, name)
	if err != nil {
		fmt.Fprint(conn, err.Error())
		conn.Close()
		return
	}
	//// if two names are the same check function

	s.Chat(conn)

	s.Leave(conn)
}
func (s *Server) Chat(conn net.Conn) {
	s.LoadHistory(conn)
	message := Fmessage(s, conn, "", ModeJoinChat)
	s.SendMessage(conn, message)
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			break
		}
		message = Fmessage(s, conn, message, ModeSendMessage)
		s.SendMessage(conn, message)
		s.SaveMessage(message)
	}

}

func (s *Server) SendMessage(conn net.Conn, message string) {
	time := time.Now().Format(TimeDefault)
	if message == "" {
		//check this again
		fmt.Fprintf(conn, PatternSending, time, s.Connections[conn])
		return
	}
	//sending messages to all but current
	sMessage := fmt.Sprint("\n%s", ColorYellow, ColorReset, message)
	s.mutex.Lock()
	for con := range s.Connections {
		if con != conn {
			fmt.Fprint(con, sMessage)
		}
		//tobe removed?
		fmt.Fprintf(con, PatternSending, time, s.Connections[con])
	}
	s.mutex.Unlock()
}

func (s *Server) LoadHistory(conn net.Conn) {
	for _, m := range s.AllMessages {
		fmt.Fprintf(conn, m)
	}
}

func ReadFile() ([]string, error) {
	//read file
	// Open the file
	file, err := os.Open("files/welcome.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// Create a new scanner for reading line by line
	scanner := bufio.NewScanner(file)

	// Read the file line by line
	var line []string
	for scanner.Scan() {
		line = append(line, scanner.Text()+"\r\n")

		// Do   something with the line
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
	//fmt.Println(line)
	return line, nil
}

func (s *Server) SaveMessage(message string) {
	s.mutex.Lock()
	s.AllMessages = append(s.AllMessages, message)
	s.mutex.Unlock()
}

// message formating

func Fmessage(serv *Server, conn net.Conn, message string, mode int) string {
	serv.mutex.Lock()
	name := serv.Connections[conn]
	defer serv.mutex.Unlock()
	switch mode {
	case ModeSendMessage:
		if message == "\n" {
			return ""
		}
		time := time.Now().Format(TimeDefault)
		message := fmt.Sprintf(PatternMessage, time, name, message)
		return message
	case ModeLeftChat:
		message := fmt.Sprintf(ColorYellow+PatternLeftChat+ColorReset, name)
		return message
	case ModeJoinChat:
		message := fmt.Sprintf(ColorYellow+PatternJoinChat+ColorReset, name)
		return message
	}
	return message
}
func (s *Server) Leave(conn net.Conn) {
	s.mutex.Lock()
	delete(s.UsedNames, s.Connections[conn])
	delete(s.Connections, conn)
	log.Printf("Connect %v was left", conn.RemoteAddr())
	s.mutex.Unlock()
}
