package transport

import (
	"bufio"
	"net"
)

type ConnWrapper struct {
	conn   net.Conn
	reader *bufio.Reader
}

func New(conn net.Conn) *ConnWrapper {
	reader := bufio.NewReader(conn)

	return &ConnWrapper{
		conn:   conn,
		reader: reader,
	}
}

func (c *ConnWrapper) WriteMessage(msg string) error {
	buf := []byte(msg)
	buf = append(buf, '\n')

	_, err := c.conn.Write([]byte(buf))
	return err
}

func (c *ConnWrapper) ReadMessage() (string, error) {
	l, _, err := c.reader.ReadLine()
	if err != nil {
		return "", err
	}
	return string(l), nil
}

func (c *ConnWrapper) RemoteAddr() string {
	return c.conn.RemoteAddr().String()
}

func (c *ConnWrapper) Close() error {
	return c.conn.Close()
}
