package GoGameServer

import (
  "fmt"
  "net"
  "os"
  "bufio"
)



func ReadLineHelper( b *bufio.Reader ) (p []byte, err os.Error ) {
  if p, err = b.ReadSlice('\n'); err != nil {
    return nil, err
  }
  var i int
  for i = len( p ); i > 0; i-- {
    if c:= p[i-1]; c != ' ' && c != '\r' && c != '\t' && c != '\n' {
      break
    }
    if string(p[i-2:i]) == "\r\n" {
      return p[0:i-2], nil
    } else if string(p[i-1:i]) == "\n" {
      return p[0:i-1], nil
    } else return p[0:i], nil


  }
  return nil, nil
}



func ReadLine( conn *net.Conn ) (string, os.Error) {
  reader := bufio.NewReader( (*conn) )
  line, err := ReadLineHelper( reader )
  if err != nil {
    fmt.Print( "Error encountered when reading a line!\n" )
    return "ERROR", os.NewError( "Couldn't read line" )
  }
  return string(line), nil
}

