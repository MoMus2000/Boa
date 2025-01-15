package main

import (
  _"fmt"
  "net"
)

type Client struct {
  Net net.Conn
}

func init_telem(Ip string) *Client {
  if Ip == "" {
    return nil
  }
  net, err := net.Dial("tcp", Ip)
  if err != nil {}
  return &Client{
    Net: net,
  }
}

func (s *Client) write_to_server(message string) error {
  if s.Net != nil {
    message_bytes := []byte(message)
    _, err := s.Net.Write(message_bytes)
    if err != nil {}
  }
  return nil
}

