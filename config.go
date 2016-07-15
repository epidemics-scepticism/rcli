package main
/*
$ go doc github.com/yawning/ricochet
package ricochet // import "github.com/yawning/ricochet"

Package ricochet implements the Ricochet chat protocol.

const MessageMaxCharacters = 2000 ...
const ContactReqMessageMaxCharacters ...
const PublicKeyBits = 1024 ...
const ContactStateOffline ContactState = iota ...
var ErrAlreadyExists = errors.New("contact already exists") ...
var ErrMessageSize = errors.New("chat message too large")
func NewEndpoint(cfg *EndpointConfig) (e *Endpoint, err error)
type ContactRequest struct { ... }
type ContactState int
type ContactStateChange struct { ... }
type Endpoint struct { ... }
type EndpointConfig struct { ... }
type IncomingMessage struct { ... }

$ go doc github.com/yawning/ricochet.EndpointConfig
type EndpointConfig struct {
	TorControlPort *bulb.Conn
	PrivateKey     *rsa.PrivateKey

	KnownContacts       []string
	BlacklistedContacts []string
	PendingContacts     map[string]*ContactRequest

	LogWriter io.Writer
}
    EndpointConfig is a Ricochet endpoint configuration.

$ go doc github.com/yawning/ricochet.Endpoint
type Endpoint struct {
	sync.Mutex

	EventChan <-chan interface{}

	// Has unexported fields.
}
    Endpoint is a active Ricochet client/server instance.


func NewEndpoint(cfg *EndpointConfig) (e *Endpoint, err error)
func (e *Endpoint) AddContact(hostname string, requestData *ContactRequest) error
func (e *Endpoint) BlacklistContact(hostname string, set bool) error
func (e *Endpoint) RemoveContact(hostname string) error
func (e *Endpoint) SendMsg(hostname, message string) error
*/
import (
	"os"
	"io/ioutil"
	"crypto/rsa"
	"crypto/rand"
	"encoding/json"

	"github.com/yawning/bulb"
	"github.com/yawning/ricochet"
)

var (
	configFile = os.Getenv("HOME") + "/.rcli.json"
)

type RCLIConfig struct {
	Key *rsa.PrivateKey
	Contacts []string
	Blacklist []string
	Pending map[string]*ricochet.ContactRequest
	e *ricochet.Endpoint
}

func (r *RCLIConfig) generate() error {
	if key, e := rsa.GenerateKey(rand.Reader, ricochet.PublicKeyBits); e != nil {
		return e
	} else {
		r.Key = key
	}
	r.Pending = make(map[string]*ricochet.ContactRequest)
	r.Contacts = make([]string, 0)
	r.Blacklist = make([]string, 0)
	return r.Save()
}

func (r *RCLIConfig) Save() error {
	out, e := json.Marshal(r)
	if e != nil {
		return e
	}
	if e := ioutil.WriteFile(configFile, out, 0600); e != nil {
		return e
	}
	return nil
}

func (r *RCLIConfig) Load() (*ricochet.Endpoint, error) {
	if _, e := os.Stat(configFile); e != nil {
		if os.IsNotExist(e) {
			if e := r.generate(); e != nil {
				return nil, e
			}
		}
	} else {
		in, e := ioutil.ReadFile(configFile)
		if e != nil {
			return nil, e
		}
		if e := json.Unmarshal(in, r); e != nil {
			return nil, e
		}
	}
	c,e := bulb.Dial("tcp", *cHost + ":" + *cPort)
	if e != nil {
		return nil, e
	}
	c.Authenticate("")
	endcfg := &ricochet.EndpointConfig{
		c,
		r.Key,
		r.Contacts,
		r.Blacklist,
		r.Pending,
		os.Stdout,
	}
	if r.e, e = ricochet.NewEndpoint(endcfg); e != nil {
		return nil, e
	} else {
		return r.e, nil
	}
}
