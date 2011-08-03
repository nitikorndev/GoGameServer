package GoGameServer

import (
  "net"
  "container/list"
)


type User struct {
  Id            int64
  Username      string
  Conn          *net.Conn
  UserList      *list.List
  Disconnecting bool
}



func NewUser( id int64, username string, conn *net.Conn, userList *list.List ) *User {
  u := &User{ id, username, conn, userList, false }
  return u
}



func (u *User) Disconnect() {
  if u.Conn != nil {
    (*u.Conn).Close()
    u.Conn = nil
  }
  u.DeleteFromList()
}



func (u *User) Write( msg string ) {
  (*u.Conn).Write( []byte(msg) )
}



func (u *User) Equal( b *User ) bool {
  if u.Id == b.Id {
    return true
  }
  return false
}



func (u *User) DeleteFromList() {
  for e := u.UserList.Front(); e != nil; e = e.Next() {
    b := e.Value.(*User)
    if u.Equal( b ) {
      u.UserList.Remove( e )
      break
    }
  }
}



// If user has old, but active connections, this function drops them.
// (almost the same as DeleteFromList())
func (u *User) DropOldConnections() {
  for e := u.UserList.Front(); e != nil; e = e.Next() {
    b := e.Value.(*User)
    if u.Equal( b ) {
      u.Write( ERROR_USER_CONNECTED_SECOND_TIME+":You have connected to here elsewhere; dropping this connection.\n" )
      u.Disconnecting = true
    }
  }
}



// UserList related:
func GetConnectedUsers( userList *list.List ) *list.List {
  list := list.New()
  for e := userList.Front(); e != nil; e = e.Next() {
    list.PushBack( e.Value.(*User).Username )
  }
  return list
}

