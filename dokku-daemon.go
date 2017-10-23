package main

import (
        "encoding/json"
        "fmt"
        "github.com/takama/daemon"
        "log"
        "net"
        "os"
        "os/exec"
        "os/signal"
        "syscall"
)

const (
        name        = "dokku-daemon"
        description = "Dokku Daemon"

        socketFolder = "/var/run/dokku-daemon"
        socketFile   = "dokku-daemon.sock"
        user         = "dokku"
        group        = "dokku"
        perms        = 0777
)

var socketPath = fmt.Sprintf("%s/%s", socketFolder, socketFile)

var dependencies = []string{}
var logFile os.File
var stdlog, errlog *log.Logger

// Service has embedded daemon
type Service struct {
        daemon.Daemon
}

type Response struct {
        Status string `json:"status"`
        Output string `json:"output"`
}

func (service *Service) Manage() (string, error) {

        usage := "Usage: myservice install | remove | start | stop | status"

        // if received any kind of command, do it
        if len(os.Args) > 1 {
                command := os.Args[1]
                switch command {
                case "install":
                        return installService(service)
                case "remove":
                        return removeService(service)
                case "start":
                        return service.Start()
                case "stop":
                        return service.Stop()
                case "status":
                        return service.Status()
                default:
                        return usage, nil
                }
        }

        // Delete the socket file if it exist
        if _, err := os.Stat(socketPath); err == nil {
                fmt.Println("Socket exist")
                err = os.Remove(socketPath)
                if err != nil {
                        fmt.Println("Socket rm error")
                }
        } else {
                fmt.Println("Socket does not exist")
        }

        interrupt := make(chan os.Signal, 1)
        signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

        // Set up listener for defined socket path
        listener, err := net.Listen("unix", socketPath)
        if err != nil {
                return "Possibly was a problem with the socket creation", err
        }

        listen := make(chan net.Conn, 100)
        go acceptConnection(listener, listen)

        for {
                select {
                case conn := <-listen:
                        go handleClient(conn)
                case killSignal := <-interrupt:
                        stdlog.Println("Got signal:", killSignal)
                        listener.Close()
                        if killSignal == os.Interrupt {
                                return "Daemon was interruped by system signal", nil
                        }
                        return "Daemon was killed", nil
                }
        }

        // never happen, but need to complete code
        return usage, nil
}

func acceptConnection(listener net.Listener, listen chan<- net.Conn) {
        fmt.Println("acceptConnection")
        for {
                conn, err := listener.Accept()
                if err != nil {
                        continue
                }
                listen <- conn
        }
}

func installService(service *Service) (string, error) {
        if _, err := os.Stat(socketFolder); os.IsNotExist(err) {
                os.Mkdir(socketFolder, perms)
                _, err := exec.Command("chown", fmt.Sprintf("%s:%s", user, group), socketFolder).Output()
                if err != nil {
                        errlog.Println("Error while chown socket folder")
                }
        }

        return service.Install()
}

func removeService(service *Service) (string, error) {
        // Remove socket folder and socket file
        if _, err := os.Stat(socketFolder); err == nil {
                err = os.RemoveAll(socketFolder)
        }

        return service.Remove()
}

// Runs the given bash command
func runcmd(cmd string, shell bool) ([]byte, error) {
        var cmdOut, cmdErr bytes.Buffer
        var cmd Cmd
        if shell {
                cmd = exec.Command("bash", "-c", cmd)
        } else {
                cmd = exec.Command("-c", cmd)
        }

        cmd.Stdout = &cmdOut
        cmd.Stderr = &cmdErr

        err := cmd.Run()

        if err != nil {
                return nil, err
        }
        return out, nil
}

func isValidCommand(cmd string) bool {
        return true
}

func handleClient(client net.Conn) {
        for {
                buf := make([]byte, 4*4096)
                numbytes, err := client.Read(buf)
                if numbytes == 0 || err != nil {
                        return
                }

                receivedData := buf[0:numbytes]
                var receivedDataStr = string(receivedData)

                var commandOut []byte
                commandOut, err = runcmd(fmt.Sprintf("dokku %s", receivedDataStr), true)

                var resp Response

                if err != nil {
                        resp = Response{
                                Status: "error",
                                Output: "Command is not valid"}
                } else {
                        resp = Response{
                                Status: "success",
                                Output: string(commandOut)}
                }

                respJson, _ := json.Marshal(resp)
                client.Write(respJson)
                client.Write([]byte("\n"))
        }
}

func init() {
        //logFile, err := os.OpenFile("/var/log/dokku-daemon.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
        //if err != nil {
        //      return
        //}

        stdlog = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
        errlog = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
}

func main() {
        defer logFile.Close()
        srv, err := daemon.New(name, description, dependencies...)
        if err != nil {
                errlog.Println("Error: ", err)
                os.Exit(1)
        }
        service := &Service{srv}
        status, err := service.Manage()
        if err != nil {
                errlog.Println(status, "\nError: ", err)
                os.Exit(1)
        }
        fmt.Println(status)
}
