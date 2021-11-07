package Messages

// -- This package contains message structs
//---------------------------------------------------  NOTES -----------------------------------------------------------
// --- Finish valid function
//--------------------------------------------------- Notes -----------------------------------------------------------
//Header contains nessesary header information
type Header struct {
	opcode uint32 //Operation code
	length uint32 //How many bytes in length follow header in the message
}

//Valid function determies if the length is correct
//given the opcode
func (h *Header) Valid() bool { // Have it return error

	return true
}

//This is the IRC packet
type ircPacket struct {
	header  Header  //Store the header
	payload []uint8 //Payload
}

//CreatePacket Passes in needed header and setup payload size
func (i *ircPacket) CreatePacket(header Header) {

	i.header = header                        //Set header to header
	i.payload = make([]uint8, header.length) //make a slice payload of size header length
}
func (i *ircPacket) loadPayload(payload []uint8) {
	copy(i.payload, payload)
}

//Valid function determines if the function is valid
//Implment later
/*
func (i *ircPacket) Valid() { }

*/
/*
//ircPacketErr is the error message struct
type ircPacketErr struct {
	header Header
	errorCode uint32
}


//setError function sets the error code
func (i *ircPacketErr) setError(error uint32) {
	i.errorCode = error
}

//ircKeepAlive is the keepalive message struct
//This message should be sent periodically (5 seconds)
//to show that tcp connection is still fine
type ircKeepAlivePacket struct {
	header Header
	opcode uint32
}
*/
