package Messages

//This file contains all the operation codes and error codes
//As well as error handling stuff
const IRC_PORT string = "7734" //Standard port for IRC




const IRC_OPCODE_ERR uint32 = 0x10000001
const IRC_OPCODE_KEEPALIVE uint32 = 0x10000002

const IRC_OPCODE_HELLO uint32 = 0x10000003
const IRC_OPCODE_LISTROOMS uint32 = 0x10000004
const IRC_OPCODE_LISTROOMSRESP uint32 = 0x10000005
const IRC_OPCODE_USERSRESP uint32 = 0x10000006
const IRC_OPCODE_JOINROOM uint32 = 0x10000007
const IRC_OPCODE_LEAVEROOM uint32 = 0x10000008
const IRC_OPCODE_SENDMSG uint32 = 0x10000009
const IRC_OPCODE_TELLMSG uint32 = 0x10000010
const IRC_OPCODE_SENDPRIVMSG uint32 = 0x10000011
const IRC_OPCODE_TELLPRIVMSG uint32 = 0x10000012

const IRC_ERR_UNKNOWN uint32 = 0x20000001
const IRC_ERR_ILLEGALOPCODE uint32 = 0x20000002
const IRC_ERR_ILLEGALLENGTH uint32 = 0x20000003
const IRC_ERR_WRONGVERSION uint32 = 0x20000004
const IRC_ERR_NAMEEXISTS uint32 = 0x20000005
const IRC_ERR_ILLEGALNAME uint32 = 0x20000006
const IRC_ERR_ILLEGALMESSAGE uint32 = 0x20000007
const IRC_ERR_TOOMANYUSERS uint32 = 0x20000008
const IRC_ERR_TOOMANYROOMS uint32 = 0x20000009

const ERRORCODEBASENUM uint32 = 0x20000000

//GetCurrentVersion returns the current version of this protocol
func GetCurrentVersion() []uint8 {
	var currentVersion [4]uint8 = [4]uint8{0xFA, 0xCE, 0x0F, 0xF1}
	return currentVersion[:] //return as slice
}
