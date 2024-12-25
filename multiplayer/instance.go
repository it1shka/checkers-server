package multiplayer

const MULTIPLAYER_SERVER_PORT = ":3056"

var server = NewServer(MULTIPLAYER_SERVER_PORT)

func GetServer() *Server {
	return server
}
