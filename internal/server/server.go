package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Server struct {
	listenAddr string
}

type Request struct {
	Method string
	Path string
	Version string
	Headers map[string]string 
    Body    string
}

type TimeResponse struct {
    CurrentTime string `json:"currentTime"`
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
	buffer := make([]byte, 2048)
	n, err := conn.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("error while reading from the connection: %w", err)
	}

	rawRequest := string(buffer[:n])

	parts := strings.SplitN(rawRequest, "\r\n\r\n", 2)
	headerBlock := parts[0]

	var body string
	if len(parts) >1 {
		body = parts[1]
	}

	headerLines := strings.Split(headerBlock, "\r\n")
	if len(headerLines) == 0{
		return nil, fmt.Errorf("fetched empty request")
	}

	requestLine := headerLines[0]

	requestLineParts := strings.Split(requestLine, " ")
	if len(requestLineParts) != 3{
		return nil, fmt.Errorf("fetched incorrect request line: %s", requestLine)
	}

	headers := make(map[string]string)

	for _, line := range headerLines[1:]{
		if line == ""{
			continue
		}
		headerParts := strings.SplitN(line, ": ", 2)
		if len (headerParts) == 2{
			headers[headerParts[0]] = headerParts[1]
		}
	}


	request := &Request{
		Method: requestLineParts[0],
		Path: requestLineParts[1],
		Version: requestLineParts[2],
		Headers: headers,
		Body: body,
	}

	return request, nil
}


func (s *Server) buildResponse(request *Request) string {
	var statusLine string
	var body []byte
	contentType := "text/html" 

	if request.Path == "/api/time" && request.Method == "GET" {
		statusLine = "HTTP/1.1 200 OK"
		contentType = "application/json" 

		data := TimeResponse{
			CurrentTime: time.Now().Format(time.RFC3339),
		}

		jsonBody, err := json.Marshal(data)
		if err != nil {
			log.Printf("Error whilst running the JSON: %v", err)
			statusLine = "HTTP/1.1 500 Internal Server Error"
			body = []byte(`{"error": "Internal Server Error"}`)
		} else {
			body = jsonBody
		}

		return fmt.Sprintf("%s\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s",
			statusLine, contentType, len(body), body)
	}


	if request.Path == "/contact" && request.Method == "POST" {
		log.Printf("Fetched the form data: %s", request.Body)
		statusLine = "HTTP/1.1 200 OK"
		htmlBody := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Thanks!</title>
    <link rel="stylesheet" href="/style.css">
</head>
<body>
    <div class="container">
        <h1>Thanks!</h1>
        <p>Your message has been successfully delivered.</p>
        <a href="/">Back to the home page</a>
    </div>
</body>
</html>`
		body = []byte(htmlBody)
		return fmt.Sprintf("%s\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s",
			statusLine, contentType, len(body), body)
	}

	var filePath string
	if request.Path == "/" {
		filePath = "static/index.html"
	} else {
		filePath = "static" + request.Path
	}
	ext := filepath.Ext(filePath)
	switch ext {
	case ".css":
		contentType = "text/css"
	case ".js":
		contentType = "application/javascript"
	default:
		contentType = "text/html"
	}
	fileBody, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("File not found: %s. Serving the 404 page", filePath)
		statusLine = "HTTP/1.1 404 Not Found"
		contentType = "text/html"
		body, _ = os.ReadFile("static/404.html")
	} else {
		statusLine = "HTTP/1.1 200 OK"
		body = fileBody
	}

	return fmt.Sprintf("%s\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s",
		statusLine, contentType, len(body), body)
}
