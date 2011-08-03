package GoGameServer

import (
  "mysql"
  "fmt"
  "os"
)


type AuthenticationReply struct {
  Authenticated bool
  Id            int64
  Username      string
}


func NewAuthenticationReply( authenticated bool,
                             id int64,
                             username string ) *AuthenticationReply {
  r := &AuthenticationReply{ authenticated, id, username }
  return r
}



func Authenticate( username string,
                   passwordHash string,
                   dbConn *mysql.Client ) (reply *AuthenticationReply, err os.Error) {
  // Default return argument
  reply = NewAuthenticationReply( false, -1, "" )

  // Escape input
  username = dbConn.Escape( username )
  password := dbConn.Escape( passwordHash )

  fmt.Printf( "Authenticating user: '%s:%s'\n", username, password )

  err = dbConn.Query( "SELECT * FROM users WHERE nick = '"+username+"' AND password = '"+password+"' limit 1" )
  if err != nil {
    return
  }

  result, err := dbConn.UseResult()
  defer dbConn.FreeResult()
  if err != nil {
    return
  }

  // Fetch the row
  row := result.FetchMap()

  // If we found it the client got the username and password right
  if row != nil {
    id       := row["id"].(int64)
    nick     := row["nick"].(string)

    reply = NewAuthenticationReply( true, id, nick )
    return
  } else {
    err = os.NewError( "Wrong username or password." )
  }

  return
}

