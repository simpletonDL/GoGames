package client

import (
	"encoding/json"
	"fmt"
)

func (g *Game) HandleServerGameState() {
	decoder := json.NewDecoder(g.Conn)

	for {
		err := decoder.Decode(&g.GameState)
		if err != nil {
			fmt.Printf("GameState decoding error: %v\n", err)
		}
	}
}
