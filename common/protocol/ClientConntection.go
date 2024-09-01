package protocol

type ClientInitializationCommand struct {
	Nickname string
}

func NewClientInitializationCommand(nickname string) ClientInitializationCommand {
	return ClientInitializationCommand{Nickname: nickname}
}
