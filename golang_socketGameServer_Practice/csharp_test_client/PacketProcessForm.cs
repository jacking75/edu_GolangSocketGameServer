using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace csharp_test_client
{
    public partial class mainForm
    {
        Dictionary<PACKET_ID, Action<byte[]>> PacketFuncDic = new Dictionary<PACKET_ID, Action<byte[]>>();

        void SetPacketHandler()
        {
            PacketFuncDic.Add(PACKET_ID.PACKET_ID_ECHO, PacketProcess_Echo);
            PacketFuncDic.Add(PACKET_ID.PACKET_ID_ERROR_NTF, PacketProcess_ErrorNotify);
            PacketFuncDic.Add(PACKET_ID.PACKET_ID_LOGIN_RES, PacketProcess_LoginResponse);
            PacketFuncDic.Add(PACKET_ID.PACKET_ID_ROOM_ENTER_RES, PacketProcess_RoomEnterResponse);
            PacketFuncDic.Add(PACKET_ID.PACKET_ID_ROOM_USER_LIST_NTF, PacketProcess_RoomUserListNotify);
            PacketFuncDic.Add(PACKET_ID.PACKET_ID_ROOM_NEW_USER_NTF, PacketProcess_RoomNewUserNotify);
            PacketFuncDic.Add(PACKET_ID.PACKET_ID_ROOM_LEAVE_RES, PacketProcess_RoomLeaveResponse);
            PacketFuncDic.Add(PACKET_ID.PACKET_ID_ROOM_LEAVE_USER_NTF, PacketProcess_RoomLeaveUserNotify);
            PacketFuncDic.Add(PACKET_ID.PACKET_ID_ROOM_CHAT_RES, PacketProcess_RoomChatResponse);            
            PacketFuncDic.Add(PACKET_ID.PACKET_ID_ROOM_CHAT_NOTIFY, PacketProcess_RoomChatNotify);
            PacketFuncDic.Add(PACKET_ID.PACKET_ID_ROOM_WHISPER_RES, PacketProcess_RoomWhisperResponse);
            PacketFuncDic.Add(PACKET_ID.PACKET_ID_ROOM_WHISPER_NOTIFY, PacketProcess_RoomWhisperNotify);
            PacketFuncDic.Add(PACKET_ID.PACKET_ID_ROOM_RELAY_NTF, PacketProcess_RoomRelayNotify);
            
            PacketFuncDic.Add(PACKET_ID.PACKET_ID_GAME_START_RES, PacketProcess_GameStartResponse);
            PacketFuncDic.Add(PACKET_ID.PACKET_ID_GAME_START_NTF, PacketProcess_GameStartNotify);
            
            PacketFuncDic.Add(PACKET_ID.PACKET_ID_GAME_BATTING_RES, PacketProcess_GameBattingResponse);
            PacketFuncDic.Add(PACKET_ID.PACKET_ID_GAME_BATTING_NTF, PacketProcess_GameBattingNotify);
            
            PacketFuncDic.Add(PACKET_ID.PACKET_ID_GAME_RESULT_NTF, PacketProcess_GameResultNotify);
            
            PacketFuncDic.Add(PACKET_ID.PACKET_ID_ROOM_MASTER_CHANGE_NTF, PacketProcess_RoomMasterChangeNotify);
            
            PacketFuncDic.Add(PACKET_ID.PACKET_ID_GAME_USER_INFO_NTF, PacketProcess_RoomGameUserInfoNotify);
            
        }

        void PacketProcess(PacketData packet)
        {
            var packetType = (PACKET_ID)packet.PacketID;
            //DevLog.Write("Packet Error:  PacketID:{packet.PacketID.ToString()},  Error: {(ERROR_CODE)packet.Result}");
            //DevLog.Write("RawPacket: " + packet.PacketID.ToString() + ", " + PacketDump.Bytes(packet.BodyData));

            if (PacketFuncDic.ContainsKey(packetType))
            {
                PacketFuncDic[packetType](packet.BodyData);
            }
            else
            {
                DevLog.Write("Unknown Packet Id: " + packet.PacketID.ToString());
            }         
        }

        void PacketProcess_Echo(byte[] bodyData)
        {
            DevLog.Write($"Echo 받음:  {bodyData.Length}");
        }

        void PacketProcess_ErrorNotify(byte[] bodyData)
        {
            var notifyPkt = new ErrorNtfPacket();
            notifyPkt.FromBytes(bodyData);

            DevLog.Write($"에러 통보 받음:  {notifyPkt.Error}");
        }


        void PacketProcess_LoginResponse(byte[] bodyData)
        {
            var responsePkt = new LoginResPacket();
            responsePkt.FromBytes(bodyData);

            DevLog.Write($"로그인 결과:  {(ERROR_CODE)responsePkt.Result}");
        }


        void PacketProcess_RoomEnterResponse(byte[] bodyData)
        {
            var responsePkt = new RoomEnterResPacket();
            responsePkt.FromBytes(bodyData);

            DevLog.Write($"방 입장 결과:  {(ERROR_CODE)responsePkt.Result}");

            labelMasterUser.Text = Convert.ToString(responsePkt.MasterUserUniqueId);
        }

        void PacketProcess_RoomUserListNotify(byte[] bodyData)
        {
            var notifyPkt = new RoomUserListNtfPacket();
            notifyPkt.FromBytes(bodyData);

            for (int i = 0; i < notifyPkt.UserCount; ++i)
            {
                AddRoomUserList(notifyPkt.UserUniqueIdList[i], notifyPkt.UserIDList[i]);
            }

            DevLog.Write($"방의 기존 유저 리스트 받음");
        }

        void PacketProcess_RoomNewUserNotify(byte[] bodyData)
        {
            var notifyPkt = new RoomNewUserNtfPacket();
            notifyPkt.FromBytes(bodyData);

            AddRoomUserList(notifyPkt.UserUniqueId, notifyPkt.UserID);
            
            DevLog.Write($"방에 새로 들어온 유저 받음");
        }


        void PacketProcess_RoomLeaveResponse(byte[] bodyData)
        {
            var responsePkt = new RoomLeaveResPacket();
            responsePkt.FromBytes(bodyData);

            DevLog.Write($"방 나가기 결과:  {(ERROR_CODE)responsePkt.Result}");
            
            listBoxRoomChatMsg.Items.Clear();
            listBoxRoomUserList.Items.Clear();
            labelMasterUser.Text = "NONE";
        }

        void PacketProcess_RoomLeaveUserNotify(byte[] bodyData)
        {
            var notifyPkt = new RoomLeaveUserNtfPacket();
            notifyPkt.FromBytes(bodyData);

            RemoveRoomUserList(notifyPkt.UserUniqueId);

            DevLog.Write($"방에서 나간 유저 받음");
        }


        void PacketProcess_RoomChatResponse(byte[] bodyData)
        {
            var responsePkt = new RoomChatResPacket();
            responsePkt.FromBytes(bodyData);

            var errorCode = (ERROR_CODE)responsePkt.Result;
            var msg = $"방 채팅 요청 결과:  {(ERROR_CODE)responsePkt.Result}";
            if (errorCode == ERROR_CODE.ERROR_NONE)
            {
                DevLog.Write(msg, LOG_LEVEL.ERROR);
            }
            else
            {
                AddRoomChatMessageList(0, msg);
            }
        }


        void PacketProcess_RoomChatNotify(byte[] bodyData)
        {
            var responsePkt = new RoomChatNtfPacket();
            responsePkt.FromBytes(bodyData);

            AddRoomChatMessageList(responsePkt.UserUniqueId, responsePkt.Message);
        }
        
        void PacketProcess_RoomWhisperResponse(byte[] bodyData)
        {
            var responsePkt = new RoomWhisperResPacket();
            responsePkt.FromBytes(bodyData);
            
            var msg = $"귓속말 요청 결과:  {(ERROR_CODE)responsePkt.Result}";
            DevLog.Write(msg, LOG_LEVEL.ERROR);
        }


        void PacketProcess_RoomWhisperNotify(byte[] bodyData)
        {
            var responsePkt = new RoomWhisperNtfPacket();
            responsePkt.FromBytes(bodyData);

            AddRoomWhisperMessageList(responsePkt.SendUserUniqueId, responsePkt.ReceiveUserUniqueId,  responsePkt.Message);
        }

        void AddRoomWhisperMessageList(Int64 SendUserUniqueId, Int64 ReceiveUserUniqueId, string msgssage)
        {
            var msg = $"{SendUserUniqueId} -> {ReceiveUserUniqueId}:  {msgssage}";

            if (listBoxRoomChatMsg.Items.Count > 512)
            {
                listBoxRoomChatMsg.Items.Clear();
            }

            listBoxRoomChatMsg.Items.Add(msg);
            listBoxRoomChatMsg.SelectedIndex = listBoxRoomChatMsg.Items.Count - 1;
        }
        
        void AddRoomChatMessageList(Int64 userUniqueId, string msgssage)
        {
            var msg = $"{userUniqueId}:  {msgssage}";

            if (listBoxRoomChatMsg.Items.Count > 512)
            {
                listBoxRoomChatMsg.Items.Clear();
            }

            listBoxRoomChatMsg.Items.Add(msg);
            listBoxRoomChatMsg.SelectedIndex = listBoxRoomChatMsg.Items.Count - 1;
        }


        void PacketProcess_RoomRelayNotify(byte[] bodyData)
        {
            var notifyPkt = new RoomRelayNtfPacket();
            notifyPkt.FromBytes(bodyData);

            var stringData = Encoding.UTF8.GetString(notifyPkt.RelayData);
            DevLog.Write($"방에서 릴레이 받음. {notifyPkt.UserUniqueId} - {stringData}");
        }
        
        // 게임 시작
        void PacketProcess_GameStartResponse(byte[] bodyData)
        {
            var responsePkt = new GameStartResPacket();
            responsePkt.FromBytes(bodyData);

            var msg = $"게임 시작 요청 결과:  {(ERROR_CODE)responsePkt.Result}";
            DevLog.Write(msg, LOG_LEVEL.ERROR);
        }
        
        void PacketProcess_GameStartNotify(byte[] bodyData)
        {
            var responsePkt = new GameStartNtfPacket();
            responsePkt.FromBytes(bodyData);

            var timeover = Convert.ToDouble(responsePkt.statusChangeCompletionMillSec);
            DateTime origin = new DateTime(1970, 1, 1, 0, 0, 0, 0);
            var msg = $"GAME: [[TIME]] {DateTime.Now.ToString("mm분 ss초")} 베팅 시작 - {origin.AddSeconds(timeover).ToString("mm분 ss초")}까지 베팅 미완료시 Player 자동 선택";
            
            //방 사람들에게 게임이 시작되었음을 알리고 배팅 하도록? 함수 생성해야함
            AddRoomGameStartMessageList(msg);
        }
        
        void AddRoomGameStartMessageList(string msg)
        {
            var msg1 = "GAME: START";
            var msg2 = "GAME: Please Batting ...";

            if (listBoxRoomChatMsg.Items.Count > 512)
            {
                listBoxRoomChatMsg.Items.Clear();
            }

            listBoxRoomChatMsg.Items.Add(msg1);
            listBoxRoomChatMsg.Items.Add(msg2);
            listBoxRoomChatMsg.Items.Add(msg);
            listBoxRoomChatMsg.SelectedIndex = listBoxRoomChatMsg.Items.Count - 3;

            btn_RoomLeave.Visible = false;
            btnRoomGameStart.Visible = false;
            btn_RoomEnter.Visible = false;
            btnRoomBattingBanker.Visible = true;
            btnRoomBattingPlayer.Visible = true;
            btnRoomBattingTie.Visible = true;
        }
        
        //게임 베팅
        void PacketProcess_GameBattingResponse(byte[] bodyData)
        {
            var responsePkt = new GameBattingResPacket();
            responsePkt.FromBytes(bodyData);
            
            var msg = $"베팅 요청 결과:  {(ERROR_CODE)responsePkt.Result}";
            DevLog.Write(msg, LOG_LEVEL.ERROR);
        }
        
        void PacketProcess_GameBattingNotify(byte[] bodyData)
        {
            var responsePkt = new GameBattingNtfPacket();
            responsePkt.FromBytes(bodyData);

            //방 사람들에게 베팅한거 알려줌
            if (responsePkt.SelectSide != 0)
            {
                AddRoomGameBattingMessageList(responsePkt.RoomUserUniqueId, responsePkt.SelectSide);
            }
            else
            {
                btnRoomBattingBanker.Visible = false;
                btnRoomBattingPlayer.Visible = false;
                btnRoomBattingTie.Visible = false;
            }
            
        }
        
        void AddRoomGameBattingMessageList(UInt64 SendUserUniqueId, int SelectSide)
        {
            var msg = $"GAME: [[BATTING]] {Convert.ToInt32(SendUserUniqueId)}'s choice --> {(GAMECASE)Convert.ToInt32(SelectSide)}";

            if (listBoxRoomChatMsg.Items.Count > 512)
            {
                listBoxRoomChatMsg.Items.Clear();
            }

            listBoxRoomChatMsg.Items.Add(msg);
            listBoxRoomChatMsg.SelectedIndex = listBoxRoomChatMsg.Items.Count - 1;
        }
        
        //게임 결과 통보
        void PacketProcess_GameResultNotify(byte[] bodyData)
        {
            var responsePkt = new GameResultNtfPacket();
            responsePkt.FromBytes(bodyData);
            
            var timeover = Convert.ToDouble(responsePkt.statusChangeCompletionMillSec);
            DateTime origin = new DateTime(1970, 1, 1, 0, 0, 0, 0);
            var msg = $"GAME: [[TIME]] {DateTime.Now.ToString("mm분 ss초")} 게임 결과발표 시작 - {origin.AddSeconds(timeover).ToString("mm분 ss초")} 이후 게임 재개 가능";

            //게임 결과 알려줌
            AddRoomGameResultMessageList(responsePkt.CardsBanker, responsePkt.CardsPlayer, responsePkt.PlayerScore, responsePkt.BankerScore, responsePkt.Result, msg);
        }
        
        void AddRoomGameResultMessageList(sbyte[] CardsBanker, sbyte[] CardsPlayer, sbyte PlayerScore, sbyte BankerScore, sbyte Result, string msg)
        {
            var msgPlayerCards = $"GAME: [[RESULT]] Player's Cards are {Card.get(Convert.ToInt32(CardsPlayer[0]))}, {Card.get(Convert.ToInt32(CardsPlayer[1]))}, {Card.get(Convert.ToInt32(CardsPlayer[2]))} --> Score {Convert.ToInt32(PlayerScore)}";
            var msgBankerCards = $"GAME: [[RESULT]] Banker's Cards are {Card.get(Convert.ToInt32(CardsBanker[0]))}, {Card.get(Convert.ToInt32(CardsBanker[1]))}, {Card.get(Convert.ToInt32(CardsBanker[2]))} --> Score {Convert.ToInt32(BankerScore)}";
            var msgResult = $"GAME: [[WINNER]] {(GAMECASE)Convert.ToInt32(Result)}";
            
            if (listBoxRoomChatMsg.Items.Count > 512)
            {
                listBoxRoomChatMsg.Items.Clear();
            }

            listBoxRoomChatMsg.Items.Add(msgBankerCards);
            listBoxRoomChatMsg.Items.Add(msgPlayerCards);
            listBoxRoomChatMsg.Items.Add(msgResult);
            listBoxRoomChatMsg.Items.Add(msg);
            listBoxRoomChatMsg.SelectedIndex = listBoxRoomChatMsg.Items.Count - 4;
            
            btn_RoomLeave.Visible = true;
            btnRoomGameStart.Visible = true;
            btn_RoomEnter.Visible = true;
            btnRoomBattingBanker.Visible = false;
            btnRoomBattingPlayer.Visible = false;
            btnRoomBattingTie.Visible = false;
        }
        
        // 방장 변경 알림
        void PacketProcess_RoomMasterChangeNotify(byte[] bodyData)
        {
            var responsePkt = new RoomMasterChangeNtfPacket();
            responsePkt.FromBytes(bodyData);
            
            AddRoomMasterChangeMessageList(responsePkt.MasterUserUniqueId);
        }

        void AddRoomMasterChangeMessageList(UInt64 MasterUserUniqueId)
        {
            var msg = $"GAME: Room Master Change - Master is [{MasterUserUniqueId}]";

            if (listBoxRoomChatMsg.Items.Count > 512)
            {
                listBoxRoomChatMsg.Items.Clear();
            }

            listBoxRoomChatMsg.Items.Add(msg);
            listBoxRoomChatMsg.SelectedIndex = listBoxRoomChatMsg.Items.Count - 1;
            
            labelMasterUser.Text = Convert.ToString(MasterUserUniqueId);
        }

        void PacketProcess_RoomGameUserInfoNotify(byte[] bodyData)
        {
            var responsePkt = new RoomUserInfoNtfPacket();
            responsePkt.FromBytes(bodyData);
            
            var msg = $"GAME: [[USER_INFO]] Dollar {responsePkt.Dollar}, Plays {responsePkt.Plays}, Win {responsePkt.Win}, Lose {responsePkt.Lose}";

            AddRoomGameUserInfoMessageList(msg);
        }

        void AddRoomGameUserInfoMessageList(string msg)
        {
            if (listBoxRoomChatMsg.Items.Count > 512)
            {
                listBoxRoomChatMsg.Items.Clear();
            }

            listBoxRoomChatMsg.Items.Add(msg);
            listBoxRoomChatMsg.SelectedIndex = listBoxRoomChatMsg.Items.Count - 1;
        }
    }
}
