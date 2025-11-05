package server

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type Server struct {
	listenAddr string
}

type Request struct {
	Method string
	Path string
	Version string
}

func NewServer (addr string) *Server{
	return &Server{
		listenAddr: addr,
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.listenAddr)

	if err != nil{
		return fmt.Errorf("Nije moguce pokrenuti listener: %w", err)
	}

	defer listener.Close()

	log.Printf("Server slusa na adresi %s", s.listenAddr)

	for{

		conn, err := listener.Accept()

		if err != nil {
			log.Printf("Greska prilikom prihvatanja konekcije: v%", err)
			continue
		}

		s.handleConnection(conn)
	}
}

func (s *Server) handleConnection (conn net.Conn){

	defer conn.Close()
	log.Printf("Primljena nova konekcije od %s", conn.RemoteAddr())

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)

	if err != nil {
		log.Printf("Greska prilikom citanja sa konekcije. v%", err)
		return
	}

	rawRequest := string(buffer[:n])

	lines := strings.Split(rawRequest, "\r\n")

	if len(lines) == 0 {
		log.Println("Primljen prazan zahtjev.")
		return
	}

	requestLine := lines[0]

	parts := strings.Split(requestLine, " ")
	if len(parts) != 3 {
		log.Printf("Primljena neispravna request line: %s", requestLine)
		return
	}

	request := Request {
		Method: parts[0],
		Path: parts[1],
		Version: parts[2],
	}

	log.Printf("Metoda: %s, Putanja: %s, Verzija: %s", request.Method, request.Path, request.Version)

	
	var response string

	switch request.Path {
	case "/":
		response = "HTTP/1.1 200 OK\r\n\r\nDobrodošli na početnu stranicu!"
	case "/about":
		response = "HTTP/1.1 200 OK\r\n\r\nOvo je naš sjajan Go web server."
	default:
		response = "HTTP/1.1 404 Not Found\r\n\r\nStranica nije pronađena."
	}

	_, err = conn.Write([]byte(response))

	if err != nil {
		log.Printf("Greska prilikom slanja odgovora: v%", err)
	}

	log.Println("Odgovor uspjesno poslat.")
}
