package request

import (
	"bytes"
	"fmt"
	"io"
)

const (
	CRLF                       = "\r\n"
	PROTOCOL_ERROR             = "Protocol error"
	INITIALIZED    parserState = "INITIALIZED"
	DONE           parserState = "DONE"
)

type parserState string

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	state       parserState
}

func checkUppercase(s string) bool {
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			continue
		}
		return false
	}
	return true
}

func parseRequestLine(requestLine []byte) (*RequestLine, int, error) {
	requestLineIndex := bytes.Index(requestLine, []byte(CRLF))
	readLenReqLine := requestLineIndex + len(CRLF)
	if requestLineIndex == -1 {
		return nil, 0, nil
	}
	request_line := requestLine[:requestLineIndex]
	request_line_part := bytes.Split(request_line, []byte(" "))
	if len(request_line_part) != 3 {
		return nil, 0, fmt.Errorf(PROTOCOL_ERROR)
	}
	// method should be uppercase
	method := string(request_line_part[0])
	if !checkUppercase(method) {
		return nil, 0, fmt.Errorf(PROTOCOL_ERROR)
	}

	// httpVersion should be 1.1
	requestTarget := request_line_part[1]
	httpVersion := request_line_part[2]
	httpVersionParts := bytes.SplitN(httpVersion, []byte("/"), 2)
	protocol := httpVersionParts[0]
	version := httpVersionParts[1]
	if string(protocol) != "HTTP" || string(version) != "1.1" {
		return nil, 0, fmt.Errorf(PROTOCOL_ERROR)
	}

	return &RequestLine{
		HttpVersion:   string(version),
		RequestTarget: string(requestTarget),
		Method:        method,
	}, readLenReqLine, nil
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
	for {
		switch r.state {
		case INITIALIZED:
			requestLine, readLen, err := parseRequestLine(data[read:])
			if err != nil {
				return 0, err
			}
			if readLen == 0 {
				return 0, nil
			}
			r.RequestLine = *requestLine
			r.state = DONE
			read += readLen
		case DONE:
			return read, nil
		default:
			return 0, fmt.Errorf("unkown state")
		}
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := &Request{
		state: INITIALIZED,
	}

	readStart := 0
	buffer := make([]byte, 1024)
	for request.state != DONE {
		buffLen, err := reader.Read(buffer[readStart:])
		if err != nil {
			return nil, err
		}
		readStart += buffLen
		parseLen, err := request.parse(buffer[:readStart])
		if err != nil {
			return nil, err
		}
		copy(buffer, buffer[parseLen:])
		readStart -= parseLen

	}

	return request, nil
}
