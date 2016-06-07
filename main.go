package main

import (
  "fmt"
  // "bufio"
  "log"
  "os/exec"
  // "net/http"
  "github.com/gorilla/websocket"
  "github.com/kr/pty"
)

func main() {
  c, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/?type=terminal", nil)
  if err != nil {
    log.Fatal("dial:", err)
  }

  defer c.Close()

  cmd := exec.Command("bash")
  f, err := pty.Start(cmd)
  if err != nil {
      panic(err)
  }

  go func() {
    // f.Write([]byte("ls\n"))
    // f.Write([]byte{4}) // EOT
    for {
      _, command, err := c.ReadMessage()
      if err != nil {
        log.Println("read:", err)
        return
      }
      f.Write(command)
    }
  }()


  // scanner := bufio.NewScanner(f)
  // for scanner.Scan() {
  //   data := scanner.Bytes()
  //   // fmt.Println(">>", data, "<<")
  //   c.WriteMessage(websocket.TextMessage, data)
  // }
  buf := make([]byte, 1024)
  for {
    size, err := f.Read(buf)
    if err != nil {
      fmt.Println("error", err)
    }
    c.WriteMessage(websocket.TextMessage, buf[:size])
  }
}
