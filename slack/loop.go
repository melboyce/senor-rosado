package slack

// Loop ...
func Loop(conn *Conn, cmds []Command) {
	mchan := make(chan Message)
	go commandProcessor(conn, mchan, cmds)

	for {
		m, err := conn.Get()
		if err != nil {
			panic(err)
		}

		if m.Type != "message" || m.Text == "" {
			continue
		}

		mchan <- m
	}
}
