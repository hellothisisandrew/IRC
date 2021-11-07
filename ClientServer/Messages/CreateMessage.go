package Messages

import "errors"

//This package create messages
// -------------------- Notes ----------------------------------------
//Big Edian is used thorughout program
//---------------------- Notes -------------------------------------

//MakeLabel takes a label and converts to []uint8
//Labels must be be length 1 - 20 inclusive
//Labels that are less than 20 must end with null character 0x00
//Cannot start or end with space
func MakeLabel(label string) ([]uint8, error) {
	var toReturn [20]uint8 //Holds the label

	if len(label) < 1 || len(label) > 20 { //Make sure label is correct length
		return nil, errors.New("label length error")
	}
	//Doing this here because why not
	if len(label) < 20 {
		toReturn[len(toReturn)-1] = 0x00 //Place null terminator
	}

	//Create the slice
	//Check for beginning and ending spaces
	//Check for UTF-8 values here rather than elsewhere
	if int([]rune(label)[0]) == 0x20 || int([]rune(label)[len(label)-1]) == 0x20 {
		return nil, errors.New("label began or ended with space")
	}
	//Loop through label and check
	for i := 0; i < len(label); i++ {
		if chr := int([]rune(label)[i]); chr < 0x20 || chr > 0x7e { //DO NOT convert to uint8 here as it may be cut down to something acceptable
			return nil, errors.New("non standard label character: either UTF-8 or unacceptable ascii")
		} else {
			toReturn[i] = uint8(chr) //Convert to uint8 here
		}
	}
	return toReturn[:20], nil
}

//CreateErrorMessage creates a new errormessage
//errotCode is the specific type of error
func CreateErrorMessage(errorCode uint32) *ircPacket {
	packet := ircPacket{}                                          //Create packet
	packet.CreatePacket(Header{opcode: IRC_OPCODE_ERR, length: 4}) //Set headers and create variables

	//When creating the payload we must convert uint32 to[4]uint8
	//We are using hexadecimal and big endian hence the 20 going at index 0
	payload := [...]uint8{20, 0, 0, uint8(errorCode - 20000000)} //Subtract the rest of the number to get least sig digit

	packet.loadPayload(payload[:]) // load the payload
	return &packet
}

//CreateKeepAlive creates a new keepalive message
//Length must be 0 and opcode must be IRC_Opcode_error
func CreateKeepAlive() *ircPacket {
	packet := ircPacket{} //Create Packet
	packet.CreatePacket(Header{opcode: IRC_OPCODE_KEEPALIVE, length: 0})

	return &packet
}

//CreateHello creates a hello packet
//Sends a field of verMagic which contains
//protocol version
func createHello(label string) (*ircPacket, error) {
	packet := ircPacket{}
	packet.CreatePacket(Header{opcode: IRC_OPCODE_HELLO, length: 24})

	version := [...]uint8{0xFA, 0xCE, 0x0F, 0xF1}
	ident, err := MakeLabel(label)

	if err != nil {
		return nil, err
	}
	packet.loadPayload(append(version[:], ident[:]...))
	return &packet, nil
}

//CreateRoomsList Create the list rooms message
func createRoomsList() *ircPacket {
	packet := ircPacket{}
	packet.CreatePacket(Header{opcode: IRC_OPCODE_LISTROOMS, length: 0})

	return &packet
}
