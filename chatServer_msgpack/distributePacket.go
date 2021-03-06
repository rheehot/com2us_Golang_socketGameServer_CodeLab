package main

import (
	"github.com/vmihailenco/msgpack/v4"
	"go.uber.org/zap"

	. "gohipernetFake"

	"main/connectedSessions"
	"main/protocol"
)

func (server *ChatServer) DistributePacket(sessionIndex int32,
	sessionUniqueId uint64,
	packetData []byte,
	) {
	packetID := protocol.PeekPacketID(packetData)
	bodySize, bodyData := protocol.PeekPacketBody(packetData)
	NTELIB_LOG_DEBUG("DistributePacket", zap.Int32("sessionIndex", sessionIndex), zap.Uint64("sessionUniqueId", sessionUniqueId), zap.Uint16("PacketID", packetID))


	packet := protocol.Packet{Id: packetID}
	packet.UserSessionIndex = sessionIndex
	packet.UserSessionUniqueId = sessionUniqueId
	packet.Id = packetID
	packet.DataSize = bodySize
	packet.Data = make([]byte, packet.DataSize)
	copy(packet.Data, bodyData)

	server.PacketChan <- packet

	NTELIB_LOG_DEBUG("_distributePacket", zap.Int32("sessionIndex", sessionIndex), zap.Uint16("PacketId", packetID))
}


func (server *ChatServer) PacketProcess_goroutine() {
	NTELIB_LOG_INFO("start PacketProcess goroutine")

	for {
		if server.PacketProcess_goroutine_Impl() {
			NTELIB_LOG_INFO("Wanted Stop PacketProcess goroutine")
			break
		}
	}

	NTELIB_LOG_INFO("Stop rooms PacketProcess goroutine")
}

func (server *ChatServer) PacketProcess_goroutine_Impl() bool {
	IsWantedTermination := false  // 여기에서는 의미 없음. 서버 종료를 명시적으로 하는 경우만 유용
	defer PrintPanicStack()


	for {
		packet := <-server.PacketChan
		sessionIndex := packet.UserSessionIndex
		sessionUniqueId := packet.UserSessionUniqueId
		bodySize := packet.DataSize
		bodyData := packet.Data

		if packet.Id == protocol.PACKET_ID_LOGIN_REQ {
			ProcessPacketLogin(sessionIndex, sessionUniqueId, bodySize, bodyData)
		} else if packet.Id == protocol.PACKET_ID_SESSION_CLOSE_SYS {
			ProcessPacketSessionClosed(server,  sessionIndex, sessionUniqueId)
		} else {
			roomNumber, _ := connectedSessions.GetRoomNumber(sessionIndex)
			server.RoomMgr.PacketProcess(roomNumber, packet)
		}
	}

	return IsWantedTermination
}

func ProcessPacketLogin(sessionIndex int32,
	sessionUniqueId uint64,
	bodySize uint16,
	bodyData []byte )  {
	//DB와 연동하지 않으므로 중복 로그인만 아니면 다 성공으로 한다
	var request protocol.LoginReqPacket
	err := msgpack.Unmarshal(bodyData, &request)
	if err != nil {
		_sendLoginResult(sessionIndex, sessionUniqueId, protocol.ERROR_CODE_PACKET_DECODING_FAIL)
		return
	}

	userID := []byte(request.UserID)
	curTime := NetLib_GetCurrnetUnixTime()

	if connectedSessions.SetLogin(sessionIndex, sessionUniqueId, userID, curTime) == false {
		_sendLoginResult(sessionIndex, sessionUniqueId, protocol.ERROR_CODE_LOGIN_USER_DUPLICATION)
		return
	}

	_sendLoginResult(sessionIndex, sessionUniqueId, protocol.ERROR_CODE_NONE)
}

func _sendLoginResult(sessionIndex int32, sessionUniqueId uint64, result int16) {
	var response protocol.LoginResPacket
	response.Result = int64(result)

	bodyData, err := msgpack.Marshal(response)
	if err != nil {
		panic(err)
	}

	sendPacket := protocol.EncodingPacketHeaderInfo(bodyData, uint16(protocol.PACKET_ID_LOGIN_RES), 0)

	NetLibIPostSendToClient(sessionIndex, sessionUniqueId, sendPacket)

	NTELIB_LOG_DEBUG("SendLoginResult", zap.Int32("sessionIndex", sessionIndex), zap.Int16("result", result))
}


func ProcessPacketSessionClosed(server *ChatServer, sessionIndex int32, sessionUniqueId uint64) {
	roomNumber, _ := connectedSessions.GetRoomNumber(sessionIndex)

	if roomNumber > -1 {
		packet := protocol.Packet{
			sessionIndex,
			sessionUniqueId,
			protocol.PACKET_ID_ROOM_LEAVE_REQ,
			0,
			nil,
		}

		server.RoomMgr.PacketProcess(roomNumber, packet)
	}

	connectedSessions.RemoveSession(sessionIndex, true)
}


