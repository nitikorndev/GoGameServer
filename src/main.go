package GoGameServer

import (
  "mysql"
  "fmt"
  "net"
  "container/list"
  "os"
)


type Reply struct {
  Code string
  Data string
}

func ParseReply( reply string ) (ret *Reply, err os.Error) {
  ret = new( Reply )
  ret.Code = CODE_NOT_SET
  ret.Data = ""
  err = nil

  // Check if the reply is too short
  if len( reply ) < 3 {
    err = os.NewError( "Reply '"+reply+"' was too short." )
    return
  }

  // Snip the reply code
  ret.Code = reply[0:3]

  // If the reply contains some data
  if len( reply ) > 4 {
    // Check if the ':' is missing
    if reply[3] != ':' {
      err = os.NewError( "Reply '"+reply+"' was missing the ':'." )
      return
    }
    // Snip the data
    ret.Data = reply[4:len( reply )]
  }

  return
}



func ClientReceiver( client *User ) {
  for {
    if client.Disconnecting {
      client.Disconnect()
      break
    }

    line, err := ReadLine( client.Conn )
    if err != nil {
      client.Disconnect()
      break
    }

    reply, err := ParseReply( line )
    if err != nil {
      // If errors were encountered, just send a error and wait for next one
      client.Write( ERROR_PROTOCOL+": Protocol error.\n" )
      continue
    }
    fmt.Printf( "<%s>(%s): %s\n", client.Username, reply.Code, reply.Data )

    // Check for the quit command
    if reply.Code == USER_QUIT {
      client.Write( SERVER_MESSAGE+":Bye bye!" )
      client.Disconnect()
      break
    }
  }
}


func HandleNewClient( conn *net.Conn, dbConn *mysql.Client, userList *list.List) {
  (*conn).Write( []byte( SERVER_ASK_USERNAME+":Username: " ) )

  // Get and parse username
  username, err := ReadLine( conn )
  if err != nil {
    (*conn).Close()
    return
  }
  reply, err := ParseReply( username )
  if err != nil {
    (*conn).Close()
    return
  }
  if reply.Code != USER_USERNAME {
    (*conn).Write( []byte( ERROR_PROTOCOL+":Protocol error: expected USER_USERNAME.\n" ) )
    (*conn).Close()
    return
  }
  username = reply.Data


  // Get and parse password
  (*conn).Write( []byte( SERVER_ASK_PASSWORD+":Password(hash): " ) )
  password, err := ReadLine( conn )
  if err != nil {
    (*conn).Close()
    return
  }
  reply, err = ParseReply( password )
  if err != nil {
    (*conn).Close()
    return
  }
  if reply.Code != USER_PASSWORD {
    (*conn).Write( []byte( ERROR_PROTOCOL+":Protocol error: expected USER_PASSWORD.\n" ) )
    (*conn).Close()
    return
  }
  password = reply.Data


  // Authenticate user
  authReply, err := Authenticate( username, password, dbConn )
  if err != nil {
    fmt.Printf( "Error when authenticating: %s\n", err.String() )
    (*conn).Write( []byte( ERROR_AUTH_FAILED+":Authentication failed.\n" ) )
    (*conn).Close()
    return
  }

  // If the authentication failed
  if !authReply.Authenticated {
    fmt.Print( "Dropped connection..\n" )
    (*conn).Write( []byte( ERROR_AUTH_FAILED+":Wrong username or password.\n" ) )
    (*conn).Close()
    return
  }

  // User has been authenticated and logged in
  fmt.Printf( "User '%s' logged in.\n", username )
  (*conn).Write( []byte( SERVER_USER_AUTHENTICATED+":You have logged in.\n" ) )

  user := NewUser( authReply.Id, username, conn, userList )
  // Drop old connections, if any
  user.DropOldConnections()
  userList.PushFront( user )
  go ClientReceiver( user )
}



func main() {
  // Connect to the mysql server
  db, err := mysql.DialUnix( mysql.DEFAULT_SOCKET, "user", "pass123", "gogameserver" )
  if err != nil {
    fmt.Printf( "Error: %s\n", err.String() )
    os.Exit( 1 )
  }
  defer db.Close()

  userList := list.New()

  addr := net.TCPAddr{ net.ParseIP( "127.0.0.1" ), 9440 }
  netListen, err := net.ListenTCP( "tcp", &addr )
  if err != nil {
    os.Exit( 1 )
  }
  defer netListen.Close()

  for {
    fmt.Print( "Connected users:\n" )
    users := GetConnectedUsers( userList )
    for e:= users.Front(); e != nil; e = e.Next() {
      fmt.Printf( "\t%s\n", e.Value.(string) )
    }

    fmt.Print( "\nWaiting for client..\n" )
    conn, err := netListen.Accept()
    if err != nil {
      fmt.Print( "Error encountered when accepting client." )
    }
    fmt.Print( "Accepted client.\n" )

    go HandleNewClient( &conn, db, userList )
  }
}

