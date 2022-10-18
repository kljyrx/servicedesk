package helper

import (
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"net"
	"os"
)

// SSHCli Cli 连接信息
type SSHCli struct {
	User       string
	Pwd        string
	Addr       string
	Client     *ssh.Client
	Session    *ssh.Session
	LastResult string
}

// Connect 连接对象
func (c *SSHCli) Connect() (*SSHCli, error) {
	config := &ssh.ClientConfig{}
	config.SetDefaults()
	config.User = c.User
	config.Auth = []ssh.AuthMethod{ssh.Password(c.Pwd)}
	config.HostKeyCallback = func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }
	client, err := ssh.Dial("tcp", c.Addr, config)
	if nil != err {
		return c, err
	}
	c.Client = client
	return c, nil
}

// Run 执行shell
func (c *SSHCli) Run(shell string) (string, error) {
	if c.Client == nil {
		if _, err := c.Connect(); err != nil {
			return "", err
		}
	}
	session, err := c.Client.NewSession()
	if err != nil {
		return "", err
	}
	// 关闭会话
	defer session.Close()
	buf, err := session.CombinedOutput(shell)

	c.LastResult = string(buf)
	return c.LastResult, err
}

func (c *SSHCli) RunTerminal(shell string) (string, error) {
	if c.Client == nil {
		if _, err := c.Connect(); err != nil {
			return "", err
		}
	}
	session, err := c.Client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(fd, oldState)

	//session.Stdout = stdout
	//session.Stderr = stderr
	//session.Stdin = os.Stdin

	/*termWidth, termHeight, err := terminal.GetSize(fd)
	if err != nil {
		panic(err)
	}*/
	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // enable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Request pseudo terminal
	if err := session.RequestPty("xterm", 0, 0, modes); err != nil {
		return "", err
	}

	defer session.Close()
	buf, err := session.CombinedOutput(shell)
	c.LastResult = string(buf)
	return c.LastResult, err
}
