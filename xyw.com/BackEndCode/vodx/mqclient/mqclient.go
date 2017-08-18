package mqclient

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"sync"

	"github.com/streadway/amqp"

	"xutils/xdebug"
)

type MQClient struct {
	ConnString string
	pConn      *amqp.Connection

	lock     *sync.Mutex
	disposed bool
}

func OpenDefault() (pClient *MQClient, err error) {
	user, pass, ip, port := "guest", "guest", "localhost", 5672
	return Open(user, pass, ip, port)
}

func Open(user, pass, host string, port int) (pClient *MQClient, err error) {
	conn_str := fmt.Sprintf("amqp://%s:%s@%s:%d/", user, pass, host, port)
	return Connect(conn_str)
}

func Connect(conn_str string) (pClient *MQClient, err error) {
	conn, err := amqp.Dial(conn_str)
	xdebug.LogError(err)
	if err != nil {
		err = errors.New("Failed to connect to RabbitMQ, ConnString=" + conn_str)
		return pClient, err
	}

	pClient = &MQClient{ConnString: conn_str, pConn: conn, lock: &sync.Mutex{}}
	runtime.SetFinalizer(pClient, onClosed)
	return pClient, err
}

func (p *MQClient) Close() {
	p.lock.Lock()
	defer p.lock.Unlock()

	onClosed(p)
}

func (p *MQClient) Initialized() bool {
	p.lock.Lock()
	defer p.lock.Unlock()

	return (p.pConn != nil)
}

func (p *MQClient) Disposed() bool {
	p.lock.Lock()
	defer p.lock.Unlock()

	return p.disposed
}

func onClosed(p *MQClient) {
	if !p.disposed {
		p.pConn.Close()
		p.pConn = nil

		p.disposed = true
	}
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		str := fmt.Sprintf("初始化： %s /%d\n", fun.Name(), line)
		log.Printf(str)
	}
}
