package slack

// Loop is the main chat loop. It pulls messages from Slack and pushes them
// onto a channel for processing by commandProcessor.
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
