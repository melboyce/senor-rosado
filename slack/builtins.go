package slack

// BuiltinHelpRegister ...
func BuiltinHelpRegister() (re string, help string) {
	re = `^help$`
	help = "`help` asistencia"
	return
}

// BuiltinHelpRespond ...
func BuiltinHelpRespond(conn *Conn, m *Message) {
	r := Reply{}
	for i, cmd := range conn.Commands {
		if i > 0 {
			r.Text += "\n"
		}
		r.Text += cmd.Help
	}
	conn.ReplyTo(m, r)
}
