// Package Messages Package messages is for parsing the payload ONLY
package Messages

import "errors"

//This is what will be returned from every
//parse function
type parsed struct {
	text string	//Used for storing Messages
	numerical uint32 //Used for storing numbers
	label string	//Used for storing a label
	list []string //used for storing a list
}

//Parse determines what an ircPacket is and parses its info
func Parse(packet ircPacket) (*parsed, uint32, error) {
	//Check that the length is valid
	if packet.validLen() != true {
		return nil, IRC_ERR_ILLEGALLENGTH, errors.New("invalid payload length") //Return error
	}

	//Go thorugh all the parse functions until you find the correct
	switch OPCODE := packet.opCode(); OPCODE {
	case IRC_OPCODE_ERR:
		return parseError(packet.getPayload())
	case IRC_OPCODE_KEEPALIVE:
		return nil, 0, nil
	case IRC_OPCODE_HELLO:
		return parseHello(packet.getPayload())
	case IRC_OPCODE_LISTROOMS:
		return nil, 0, nil
	case IRC_OPCODE_JOINROOM:
		return parseJoinLeaveRoom(packet.getPayload())
	case IRC_OPCODE_LEAVEROOM:
		parseJoinLeaveRoom(packet.getPayload())
	case IRC_OPCODE_SENDMSG:
		parseSendMessage(packet.getPayload())
	case IRC_OPCODE_SENDPRIVMSG:
		parseSendMessage(packet.getPayload())




	default:
		return nil, IRC_ERR_ILLEGALOPCODE, errors.New("invalid opcode or format")
	}

}



//Parse error makes sure error message is correct
//And gathers the correct info
func parseError(payload []uint8) (*parsed, uint32, error) {
	//Must be 20, 0, 0, some error code
	check := []uint8{20,0,0}

	//compare checks first three
	for index, value := range check {
		if payload[index] != value {
			return nil, IRC_ERR_UNKNOWN, errors.New("error field incorrect")
		}
	}
	//Error field must be between 1 and 9 inclusive
	if payload[3] < 1 || payload[3] > 9 {
		return nil, IRC_ERR_UNKNOWN, errors.New("error field incorrect")
	}

	//Finally return the correct error both in the Unint
	return &parsed{numerical: ERRORCODEBASENUM + uint32(payload[3])}, ERRORCODEBASENUM + uint32(payload[3]), nil
}
//parseHello parses a hello packet can return IRCILLEGALNAME
//hello introduces the client to the server so i contains info
func parseHello(payload []uint8) (*parsed, uint32, error) {

	version := payload[:4] //get the version from the payload
	current := GetCurrentVersion() //get the current version
	idenity, err := extractLabel(payload[4:])

	//Check the idenity and make sure it is valid
	if err != nil {
		return nil, IRC_ERR_ILLEGALNAME, err
	}

	//loop over version adn check against current
	for index, value := range version {
		if current[index] != value {
			return nil, IRC_ERR_WRONGVERSION, errors.New("wrong version")
		}
	}
	//Otherwise return everything as normal
	return &parsed{label: idenity, text: string(version)}, 0, nil
}
//parseJoinLeaveRoom parses a join or leave room opcode packet
func parseJoinLeaveRoom(payload []uint8) (*parsed, uint32, error) {

	roomName, err := extractLabel(payload[:20]) //extract the label

	if err != nil {
		return nil, IRC_ERR_ILLEGALNAME, err
	}

	return &parsed{label: roomName}, 0, nil		//Extract the roomName
}
//parseSendMessage parses send message packets
func parseSendMessage(payload []uint8) (*parsed, uint32, error) {

	target, err := extractLabel(payload[:20])
	text := string(payload[20:])

	//Err from extract label means an illegal name
	if err != nil {
		return nil, IRC_ERR_ILLEGALNAME, err
	}
	//Clear the messages
	err = checkMessage(text)
	err1 := checkTerminators(payload[20:])
	//Check message for error
	if err != || err1 != nil {
		return nil, IRC_ERR_ILLEGALMESSAGE, err
	}


	return &parsed{label: target, text: text}, 0, nil

}

//parseListResp parses the list response packet type
//Set room reqest to true if looking for a list of rooms
func parseListResp(payload []uint8, isRoomRequest bool) (*parsed, uint32, error){
	ident := "" //label identity
	shift := 0 //This shifts where the list starts
	listUR := []string{}

	if isRoomRequest {
		identifer, err := extractLabel(payload[:20]) //extract the label
		if err != nil {
			return *, IRC_ERR_ILLEGALNAME, err
		}
		ident = identifer
		shift = 20
	}
	//Loop over
	listIndex := 0 //for the new list array
	for i := shift; len(payload[shift:]) > i; i += 20 { 		//Check multiples of 20
			list[listIndex] = string(payload[i: i + 20])
			if err := checkLabel(list[listIndex]); err != nil {
				return nil, IRC_ERR_UNKNOWN, errors.New("one or more illegal names in list")
			}
	}
	return &parsed{list: listUR, label: ident}, 0, nil //Return the the identity and the

}




//checkMessage makes sure a message is in protocol bounds
func checkMessage(text string)  error {
	if len(text) > 8000 {
		return  errors.New("message is too long")
	}

	//Loop through label and check not allowed characters
	for i := 0; i < len(text); i++ {
		if chr := int([]rune(text)[i]); (chr < 0x20 || chr > 0x7e ) && (chr != 0x0D || chr != 0x0A) {			//DO NOT convert to uint8 here as it may be cut down to something acceptable
			return errors.New("non standard label character: either UTF-8 or unacceptable ascii")
		}
	}
	//otheriwise return nil
	return nil

}
//checkTerminators in message //This is just for cleanup
func checkTerminators(msg []uint8) error {
	for index, value := range msg {
		if index == len(msg) -1 && value == 0x00 {
			return nil
		}
		if value == 0x00 {
			return errors.New("illegal message")
		}
	}
}

//extractLabel extracts the label from an array and returns its string version
func extractLabel(label []uint8) (string, error) {
	toConvert := []uint8{}

	for index, value := range label {
		if value == 0x00 {	//Check if you hit a null termiantor
			toConvert = label[:index + 1]
			break
		}
	}
	toReturn := string(toConvert) //convert to a string

	if err := checkLabel(toReturn); err != nil {
		return "", err
	}
	return toReturn, nil
}

//Chekc label checks to make sure the label fits
func checkLabel(label string)  error {


	if len(label) < 1 || len(label) > 20 {		//Make sure label is correct length
		return  errors.New("label length error")
	}


	//Create the slice
	//Check for beginning and ending spaces
	//Check for UTF-8 values here rather than elsewhere
	if int([]rune(label)[0]) == 0x20 || int([]rune(label)[len(label) - 1]) == 0x20 {
		return errors.New("label began or ended with space")
	}
	//Loop through label and check
	for i := 0; i < len(label); i++ {
		if chr := int([]rune(label)[i]); chr < 0x20 || chr > 0x7e {			//DO NOT convert to uint8 here as it may be cut down to something acceptable
			return errors.New("non standard label character: either UTF-8 or unacceptable ascii")
		}
	}
	return nil
}