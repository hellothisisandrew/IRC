//Rawmessages converts messages to pure byte format
//and visversa taking ircPacket structs and converting them into byte length
package Messages

import (
	"encoding/binary"
	"errors"
)

//CreateRaw function takes an ircPacket and creates a raw data in []byte slice format
//Pass by pointer to save space and time
func CreateRaw(packet *ircPacket) ([]byte,error){

	//Check if it is a valid packet
	if err:= packet.validLen(); !err {
		return nil, errors.New("invalid length")
	}


	opcode := make([]byte, 4) //Make byte slice for opcode
	length := make([]byte, 4) //Make byte slice for length
	final := make([]byte,(8 + packet.header.length))

	binary.BigEndian.PutUint32(opcode,packet.header.opcode) //set uint32 to byte type
	binary.BigEndian.PutUint32(length,packet.header.length) //

	//copy everything over into raw bytes
	copy(final[:4], opcode)
	copy(final[4:8], length)
	copy(final[8:],packet.payload)

	return final, nil
}
//CreateIrc function takes raw data and makes an ircPacket
//The length must already be known before using this function
//The raw []byte must already have divided up so there is only one packet
//Contained within it
func CreateIrc(raw []byte) (*ircPacket,error){
	opcode := binary.BigEndian.Uint32(raw[:4])
	length := binary.BigEndian.Uint32(raw[4:8])
	payload := raw[8:]

	//Check length before creating ircPacket
	if len(payload)	!= int(length) {
		return nil, errors.New("invalid raw data length issue")
	}

	packet := ircPacket{}	//create the packet
	packet.CreatePacket(Header{opcode: opcode,length: length})	//Create the header
	packet.appendPayload(payload)		//add to the payload

	return &packet, nil
}



