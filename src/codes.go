package GoGameServer

const (
  CODE_NOT_SET = "000"
)

// Server Codes
const (
  SERVER_ASK_USERNAME       = "001"
  SERVER_ASK_PASSWORD       = "002"
  SERVER_MESSAGE            = "005"
  SERVER_USER_AUTHENTICATED = "010"
)


// Error Codes
const (
  ERROR_PROTOCOL                   = "100"
  ERROR_AUTH_FAILED                = "101"
  ERROR_AUTH_TIMED_OUT             = "102"
  ERROR_USER_CONNECTED_SECOND_TIME = "103" // For cases when user has already connected, but connects again
)


// User Codes
const (
  USER_USERNAME = "201"
  USER_PASSWORD = "202"
  USER_QUIT     = "210"
)

