package agent

import (
	"fmt"

	"log"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
)

var agentInstances = 0

const AgentBasePath = "/agent/simple%d"
const SimpleAgentPinCode = "0000"
const SimpleAgentPassKey uint32 = 1024

func NextAgentPath() dbus.ObjectPath {
	p := dbus.ObjectPath(fmt.Sprintf(AgentBasePath, agentInstances))
	agentInstances += 1
	return p
}

// NewDefaultSimpleAgent return a SimpleAgent instance with default pincode and passcode
func NewDefaultSimpleAgent() *SimpleAgent {
	ag := &SimpleAgent{
		path:    NextAgentPath(),
		passKey: SimpleAgentPassKey,
		pinCode: SimpleAgentPinCode,
	}

	return ag
}

// NewSimpleAgent return a SimpleAgent instance
func NewSimpleAgent() *SimpleAgent {
	ag := &SimpleAgent{
		path: NextAgentPath(),
	}
	return ag
}

// SimpleAgent implement interface Agent1Client
type SimpleAgent struct {
	path    dbus.ObjectPath
	pinCode string
	passKey uint32
}

func (self *SimpleAgent) SetPassKey(passkey uint32) {
	self.passKey = passkey
}

func (self *SimpleAgent) SetPassCode(pinCode string) {
	self.pinCode = pinCode
}

func (self *SimpleAgent) PassKey() uint32 {
	return self.passKey
}

func (self *SimpleAgent) PassCode() string {
	return self.pinCode
}

func (self *SimpleAgent) Path() dbus.ObjectPath {
	return self.path
}

func (self *SimpleAgent) Interface() string {
	return Agent1Interface
}

func (self *SimpleAgent) Release() *dbus.Error {
	return nil
}

func (self *SimpleAgent) RequestPinCode(path dbus.ObjectPath) (string, *dbus.Error) {

	log.Printf("SimpleAgent: RequestPinCode: %s", path)

	adapterID, err := adapter.ParseAdapterID(path)
	if err != nil {
		log.Printf("SimpleAgent.RequestPinCode: Failed to load adapter %s", err)
		return "", dbus.MakeFailedError(err)
	}

	err = SetTrusted(adapterID, path)
	if err != nil {
		log.Printf("SimpleAgent.RequestPinCode SetTrusted failed: %s", err)
		return "", dbus.MakeFailedError(err)
	}

	log.Printf("SimpleAgent: Returning pin code: %s", self.pinCode)
	return self.pinCode, nil
}

func (self *SimpleAgent) DisplayPinCode(device dbus.ObjectPath, pincode string) *dbus.Error {
	log.Println(fmt.Sprintf("SimpleAgent: DisplayPinCode (%s, %s)", device, pincode))
	return nil
}

func (self *SimpleAgent) RequestPasskey(path dbus.ObjectPath) (uint32, *dbus.Error) {

	adapterID, err := adapter.ParseAdapterID(path)
	if err != nil {
		log.Printf("SimpleAgent.RequestPassKey: Failed to load adapter %s", err)
		return 0, dbus.MakeFailedError(err)
	}

	err = SetTrusted(adapterID, path)
	if err != nil {
		log.Printf("SimpleAgent.RequestPassKey: SetTrusted %s", err)
		return 0, dbus.MakeFailedError(err)
	}

	log.Printf("RequestPasskey: returning %d", self.passKey)
	return self.passKey, nil
}

func (self *SimpleAgent) DisplayPasskey(device dbus.ObjectPath, passkey uint32, entered uint16) *dbus.Error {
	log.Printf("SimpleAgent: DisplayPasskey %s, %06d entered %d", device, passkey, entered)
	return nil
}

func (self *SimpleAgent) RequestConfirmation(path dbus.ObjectPath, passkey uint32) *dbus.Error {

	log.Printf("SimpleAgent: RequestConfirmation (%s, %06d)", path, passkey)

	adapterID, err := adapter.ParseAdapterID(path)
	if err != nil {
		log.Printf("SimpleAgent: Failed to load adapter %s", err)
		return dbus.MakeFailedError(err)
	}

	err = SetTrusted(adapterID, path)
	if err != nil {
		log.Printf("Failed to set trust for %s: %s", path, err)
		return dbus.MakeFailedError(err)
	}

	log.Println("SimpleAgent: RequestConfirmation OK")
	return nil
}

func (self *SimpleAgent) RequestAuthorization(device dbus.ObjectPath) *dbus.Error {
	log.Printf("SimpleAgent: RequestAuthorization (%s)", device)
	return nil
}

func (self *SimpleAgent) AuthorizeService(device dbus.ObjectPath, uuid string) *dbus.Error {
	log.Printf("SimpleAgent: AuthorizeService (%s, %s)", device, uuid) // directly authorized
	return nil
}

func (self *SimpleAgent) Cancel() *dbus.Error {
	log.Printf("SimpleAgent: Cancel")
	return nil
}
