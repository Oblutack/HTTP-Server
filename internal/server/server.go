package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
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
		return fmt.Errorf("Unable to run the listener: %w", err)
	}

	defer listener.Close()

	log.Printf("Server is listening at address %s", s.listenAddr)

	for{

		conn, err := listener.Accept()

		if err != nil {
			log.Printf("Error while fetching the connection: v%", err)
			continue
		}

		//Added concurrency
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("Fetched new connection from %s", conn.RemoteAddr())

	request, err := s.parseRequest(conn)
	if err != nil {
		log.Printf("Unable to parse the request: %v", err)
		return
	}

	log.Printf("Method: %s, Path: %s, Version: %s", request.Method, request.Path, request.Version)

	response := s.buildResponse(request)

	_, err = conn.Write([]byte(response))
	if err != nil {
		log.Printf("Error while sending the response: %v", err)
	}

	log.Println("Response sucessfully sent.")
}


func (s *Server) parseRequest(conn net.Conn) (*Request, error) {
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("error while reading from the connection: %w", err)
	}

	rawRequest := string(buffer[:n])
	lines := strings.Split(rawRequest, "\r\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("fetched empty request")
	}

	requestLine := lines[0]
	parts := strings.Split(requestLine, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("fetched incorrect request line: %s", requestLine)
	}

	request := &Request{
		Method:  parts[0],
		Path:    parts[1],
		Version: parts[2],
	}

	return request, nil
}

func (s *Server) buildResponse(request *Request) string {
	var statusLine string
	var filePath string
	var contentType string

	switch request.Path {
	case "/":
		filePath = "static/index.html"
		statusLine = "HTTP/1.1 200 OK"
	case "/about":
		filePath = "static/about.html"
		statusLine = "HTTP/1.1 200 OK"
	case "/style.css":
		filePath = "static/style.css"
		statusLine = "HTTP/1.1 200 OK"
	default:
		filePath = "static/404.html"
		statusLine = "HTTP/1.1 404 Not Found"
	}

	ext := filepath.Ext(filePath)

	switch ext {
	case ".css":
		contentType = "text/css"
	case ".js":
		contentType = "application/javascript"
	case ".png":
		contentType = "image/png"
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	default:
		contentType = "text/html"
	}

	body, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Error whilst reading the file %s: %v", filePath, err)
		filePath = "static/404.html" 
		statusLine = "HTTP/1.1 404 Not Found"
		contentType = "text/html"
		body, _ = os.ReadFile(filePath)
	}

	response := fmt.Sprintf("%s\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s", 
		statusLine, contentType, len(body), body)

	return response
}
