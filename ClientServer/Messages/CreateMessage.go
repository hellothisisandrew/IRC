package Messages

//This package create messages
// -------------------- Notes ----------------------------------------
	//Big Edian is used thorughout program

	//--I can merge check label and make label label
//---------------------- Notes -------------------------------------

//------  *** makelabel can bereused as check label for now I can make it cleaner later

//MakeLabel takes a label and converts to []uint8
//Labels must be be length 1 - 20 inclusive
//Labels that are less than 20 must end with null character 0x00
//Cannot start or end with space
func MakeLabel(label string) ([]uint8, error){

	if err := checkLabel(label); err != nil {
		return nil ,err
	}

	//Slice must be exactly size 20
	toReturn := make([]uint8,20)
	copy(toReturn,[]byte(label))

	if len(label) < 20 {
		toReturn[len(toReturn) - 1] = 0x00		//Place null terminator
	}

	return toReturn[:20], nil
}

//CreateText Function creates a text payload of uint8
func CreateText(text string) ([]uint8, error) {

	if err := checkMessage(text); err != nil {
		return nil, err
	}
	arrayStr := []uint8(text)	//Convert and return the array


	//If the last Character is not the null terminator make it the null terminator
	if arrayStr[len(text) - 1] != 0 {
		arrayStr = append(arrayStr, 0)
	}
	return arrayStr[:], nil
}
//CreateErrorMessage creates a new errormessage
//errotCode is the specific type of error
func CreateErrorMessage(errorCode uint32) *ircPacket {
	packet := ircPacket{} //Create packet
	packet.CreatePacket(Header{opcode: IRC_OPCODE_ERR, length: 4})//Set headers and create variables

	//When creating the payload we must convert uint32 to[4]uint8
	//We are using hexadecimal and big endian hence the 20 going at index 0
	payload := [...]uint8{20,0,0,uint8(errorCode - ERRORCODEBASENUM)} //Subtract the rest of the number to get least sig digit

	packet.appendPayload(payload[:])// load the payload
	return &packet
}

//CreateKeepAlive creates a new keepalive message
//Length must be 0 and opcode must be IRC_Opcode_error
func CreateKeepAlive() *ircPacket {
	packet := ircPacket{}	//Create Packet
	packet.CreatePacket(Header{opcode: IRC_OPCODE_KEEPALIVE, length: 0})

	return &packet
}

//createHello creates a hello packet
//Sends a field of verMagic which contains
//protocol version
func createHello(label string) (*ircPacket, error) {
	packet := ircPacket{}
	packet.CreatePacket(Header{opcode: IRC_OPCODE_HELLO, length: 24})

	version := GetCurrentVersion()
	ident, err := MakeLabel(label)

	if err != nil {
		return nil, err
	}
	packet.appendPayload(version[:])
	packet.appendPayload(ident[:])
	return &packet, nil
}
//RoomsList  list rooms message That requests a list of rooms.
//Sent by client to server
func RoomsList() *ircPacket {
	packet := ircPacket{}
	packet.CreatePacket(Header{opcode: IRC_OPCODE_LISTROOMS, length: 0}) //Opcode 1
	return &packet
}

//joinChatroom sends a message to join a chatroom
//If a room by that name does not exist it prompts
//The server to create it
func joinChatRoom(roomName string) (*ircPacket, error) {

	packet := ircPacket{} //Creatte the packet
	packet.CreatePacket(Header{IRC_OPCODE_JOINROOM,20}) //Set the headers
	label, err := MakeLabel(roomName)	//Make the label

	if err != nil {
		return nil, err
	}
	packet.appendPayload(label[:]) //Add the label
	return &packet, nil
}
//leaveRoom function create the leave room message
func leaveRoom(roomName string) (*ircPacket, error){
	packet := ircPacket{}
	packet.CreatePacket(Header{IRC_OPCODE_LEAVEROOM, 20})

	label, err := MakeLabel(roomName) //Make the label

	if err != nil {		//Make sure it is in correct label format
		return nil,err
	}

	packet.appendPayload(label) //Add the label to the payload

	return &packet, nil
}

//sendMessage function creates a packet with a text message
// If private is set to true then it is a private message
func sendMessage (messasge string, target string, private bool) (*ircPacket, error) {

	packet := ircPacket{}		//Create the packet
	text, err := CreateText(messasge) //Get the array
	label, err := MakeLabel(target)

	if err != nil {
		return nil, err
	}

	//If private change opcode
	//Length of the text/label is length in header
	if private == true {
		packet.CreatePacket(Header{IRC_OPCODE_SENDPRIVMSG, uint32(len(text) + 20)})
	}else {
		packet.CreatePacket(Header{IRC_OPCODE_SENDMSG, uint32(len(text) + 20)})
	}
	packet.appendPayload(label) //Add the label
	packet.appendPayload(text) //Load the payload

	return &packet, nil
}

//listingResponse function creates the listing response packet
//Pass a room name if listing users in a room
//The list can either be users in a room or listing rooms
func listingResponse(list []string, room string) (*ircPacket, error) {
	packet := ircPacket{}
	payload := make([]uint8,20*len(list), 20*len(list) + 20) //Make a slice with length of list + 20

	//if room is empty then it is rooms request
	if room == "" {
		packet.CreatePacket(Header{IRC_OPCODE_LISTROOMSRESP, uint32(len(list)*20)}) //*20 because each item should be a lebel
	} else {
		packet.CreatePacket(Header{IRC_OPCODE_USERSRESP, uint32(len(list)*20 + 20)})	//+20 is for the identifier
		label, err := MakeLabel(room)	//Make the identfier
		packet.appendPayload(label) //add the label
		//If the label is incorrect
		if err != nil {
			return nil, err
		}
	}

	//Loop over and change string list into a single array of uint8
	for i := 0; i < len(list); i++ {
		item, err := MakeLabel(list[i])
		//Check the label for errors
		if err != nil {
			return nil, err
		}
		payload = append(payload, item... ) //append the item
	}

	packet.appendPayload(payload)	//add the payload

	return &packet, nil
}

//forwardMessage creates packet either from server informing that specfied message was posted to a room
//If sender is not an empty string it is private message to be forwarded to a specific user
func forwardMessage(target string, sender string, message string) (*ircPacket,error) {

	packet := ircPacket{}
	text, err := CreateText(message)// create the message
	recipient, err := MakeLabel(target) //recipent or target

	if err != nil {
		return nil,err
	}

	//If it is a tell message from server string is empty
	if sender == "" {
		packet.CreatePacket(Header{IRC_OPCODE_TELLMSG, uint32(len(text) + 20)} ) //Length needs to include
		packet.appendPayload(recipient)//Add the recipent
		packet.appendPayload(text) //Add the payload

	}else {
		packet.CreatePacket(Header{IRC_OPCODE_TELLPRIVMSG, uint32(len(text) + 40)})
		send, err := MakeLabel(sender) //make the sender label

		//check
		if err != nil {
			return nil,err
		}
		//load all the info in the correct order
		packet.appendPayload(recipient)
		packet.appendPayload(send)
		packet.appendPayload(text)
	}

	return &packet, nil //Return  the packet
}







