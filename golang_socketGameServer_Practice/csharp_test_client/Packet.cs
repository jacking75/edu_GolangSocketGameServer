using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace csharp_test_client
{
    struct PacketData
    {
        public Int16 DataSize;
        public Int16 PacketID;
        public SByte Type;
        public byte[] BodyData;
    }

    public class PacketDump
    {
        public static string Bytes(byte[] byteArr)
        {
            StringBuilder sb = new StringBuilder("[");
            for (int i = 0; i < byteArr.Length; ++i)
            {
                sb.Append(byteArr[i] + " ");
            }
            sb.Append("]");
            return sb.ToString();
        }
    }
    

    public class ErrorNtfPacket
    {
        public ERROR_CODE Error;

        public bool FromBytes(byte[] bodyData)
        {
            Error = (ERROR_CODE)BitConverter.ToInt16(bodyData, 0);
            return true;
        }
    }
    

    public class LoginReqPacket
    {
        byte[] UserID = new byte[PacketDef.MAX_USER_ID_BYTE_LENGTH];
        byte[] UserPW = new byte[PacketDef.MAX_USER_PW_BYTE_LENGTH];

        public void SetValue(string userID, string userPW)
        {
            Encoding.UTF8.GetBytes(userID).CopyTo(UserID, 0);
            Encoding.UTF8.GetBytes(userPW).CopyTo(UserPW, 0);
        }

        public byte[] ToBytes()
        {
            List<byte> dataSource = new List<byte>();
            dataSource.AddRange(UserID);
            dataSource.AddRange(UserPW);
            return dataSource.ToArray();
        }
    }

    public class LoginResPacket
    {
        public Int16 Result;

        public bool FromBytes(byte[] bodyData)
        {
            Result = BitConverter.ToInt16(bodyData, 0);
            return true;
        }
    }


    public class RoomEnterReqPacket
    {
        int RoomNumber;
        public void SetValue(int roomNumber)
        {
            RoomNumber = roomNumber;
        }

        public byte[] ToBytes()
        {
            List<byte> dataSource = new List<byte>();
            dataSource.AddRange(BitConverter.GetBytes(RoomNumber));
            return dataSource.ToArray();
        }
    }

    public class RoomEnterResPacket
    {
        public Int16 Result;
        public Int64 MasterUserUniqueId;

        public bool FromBytes(byte[] bodyData)
        {
            Result = BitConverter.ToInt16(bodyData, 0);
            MasterUserUniqueId = BitConverter.ToInt64(bodyData, 2);
            return true;
        }
    }

    public class RoomUserListNtfPacket
    {
        public int UserCount = 0;
        public List<Int64> UserUniqueIdList = new List<Int64>();
        public List<string> UserIDList = new List<string>();

        public bool FromBytes(byte[] bodyData)
        {
            var readPos = 0;
            var userCount = (SByte)bodyData[readPos];
            ++readPos;

            for (int i = 0; i < userCount; ++i)
            {
                var uniqeudId = BitConverter.ToInt64(bodyData, readPos);
                readPos += 8;

                var idlen = (SByte)bodyData[readPos];
                ++readPos;

                var id = Encoding.UTF8.GetString(bodyData, readPos, idlen);
                readPos += idlen;

                UserUniqueIdList.Add(uniqeudId);
                UserIDList.Add(id);
            }

            UserCount = userCount;
            return true;
        }
    }

    public class RoomNewUserNtfPacket
    {
        public Int64 UserUniqueId;
        public string UserID;

        public bool FromBytes(byte[] bodyData)
        {
            var readPos = 0;

            UserUniqueId = BitConverter.ToInt64(bodyData, readPos);
            readPos += 8;

            var idlen = (SByte)bodyData[readPos];
            ++readPos;

            UserID = Encoding.UTF8.GetString(bodyData, readPos, idlen);
            readPos += idlen;

            return true;
        }
    }


    public class RoomChatReqPacket
    {
        Int16 MsgLen;
        byte[] Msg;//= new byte[PacketDef.MAX_USER_ID_BYTE_LENGTH];

        public void SetValue(string message)
        {
            Msg = Encoding.UTF8.GetBytes(message);
            MsgLen = (Int16)Msg.Length;
        }

        public byte[] ToBytes()
        {
            List<byte> dataSource = new List<byte>();
            dataSource.AddRange(BitConverter.GetBytes(MsgLen));
            dataSource.AddRange(Msg);
            return dataSource.ToArray();
        }
    }

    public class RoomChatResPacket
    {
        public Int16 Result;
        
        public bool FromBytes(byte[] bodyData)
        {
            Result = BitConverter.ToInt16(bodyData, 0);
            return true;
        }
    }

    public class RoomChatNtfPacket
    {
        public Int64 UserUniqueId;
        public string Message;

        public bool FromBytes(byte[] bodyData)
        {
            UserUniqueId = BitConverter.ToInt64(bodyData, 0);

            var msgLen = BitConverter.ToInt16(bodyData, 8);
            byte[] messageTemp = new byte[msgLen];
            Buffer.BlockCopy(bodyData, 8 + 2, messageTemp, 0, msgLen);
            Message = Encoding.UTF8.GetString(messageTemp);
            return true;
        }
    }
    
    public class RoomWhisperReqPacket
    {
        public Int64 ReceiveUserUniqueId;
        Int16 MsgLen;
        byte[] Msg;//= new byte[PacketDef.MAX_USER_ID_BYTE_LENGTH];

        public void SetValue(Int64 ReceiveUser,string message)
        {
            ReceiveUserUniqueId = ReceiveUser;
            Msg = Encoding.UTF8.GetBytes(message);
            MsgLen = (Int16)Msg.Length;
        }

        public byte[] ToBytes()
        {
            List<byte> dataSource = new List<byte>();
            dataSource.AddRange(BitConverter.GetBytes(ReceiveUserUniqueId));
            dataSource.AddRange(BitConverter.GetBytes(MsgLen));
            dataSource.AddRange(Msg);
            return dataSource.ToArray();
        }
    }

    public class RoomWhisperResPacket
    {
        public Int16 Result;
        
        public bool FromBytes(byte[] bodyData)
        {
            Result = BitConverter.ToInt16(bodyData, 0);
            return true;
        }
    }

    public class RoomWhisperNtfPacket
    {
        public Int64 SendUserUniqueId;
        public Int64 ReceiveUserUniqueId;
        public string Message;

        public bool FromBytes(byte[] bodyData)
        {
            SendUserUniqueId = BitConverter.ToInt64(bodyData, 0);
            ReceiveUserUniqueId = BitConverter.ToInt64(bodyData, 8);
            
            var msgLen = BitConverter.ToInt16(bodyData, 16);
            byte[] messageTemp = new byte[msgLen];
            Buffer.BlockCopy(bodyData, 8 + 8 + 2, messageTemp, 0, msgLen);
            Message = Encoding.UTF8.GetString(messageTemp);
            return true;
        }
    }
    
     public class RoomLeaveResPacket
    {
        public Int16 Result;
        
        public bool FromBytes(byte[] bodyData)
        {
            Result = BitConverter.ToInt16(bodyData, 0);
            return true;
        }
    }

    public class RoomLeaveUserNtfPacket
    {
        public Int64 UserUniqueId;

        public bool FromBytes(byte[] bodyData)
        {
            UserUniqueId = BitConverter.ToInt64(bodyData, 0);
            return true;
        }
    }


    
    public class RoomRelayNtfPacket
    {
        public Int64 UserUniqueId;
        public byte[] RelayData;

        public bool FromBytes(byte[] bodyData)
        {
            UserUniqueId = BitConverter.ToInt64(bodyData, 0);

            var relayDataLen = bodyData.Length - 8;
            RelayData = new byte[relayDataLen];
            Buffer.BlockCopy(bodyData, 8, RelayData, 0, relayDataLen);
            return true;
        }
    }


    public class PingRequest
    {
        public Int16 PingNum;

        public byte[] ToBytes()
        {
            return BitConverter.GetBytes(PingNum);
        }

    }
    
    // 게임 시작
    public class GameStartReqPacket
    {
    }
    
    public class GameStartResPacket
    {
        public Int16 Result;
        
        public bool FromBytes(byte[] bodyData)
        {
            Result = BitConverter.ToInt16(bodyData, 0);
            return true;
        }
    }
    
    public class GameStartNtfPacket
    {
        public Int64 statusChangeCompletionMillSec;
        
        public bool FromBytes(byte[] bodyData)
        {
            statusChangeCompletionMillSec = BitConverter.ToInt64(bodyData, 0);
            return true;
        }
    }
    
    //게임 베팅
    public class GameBattingReqPacket
    {
        public sbyte SelectSide;
        
        public void SetValue(sbyte Select)
        {
            SelectSide = Select;
        }

        public byte[] ToBytes()
        {
            List<byte> dataSource = new List<byte>();
            dataSource.AddRange(BitConverter.GetBytes(SelectSide));
            return dataSource.ToArray();
        }
    }
    
    public class GameBattingResPacket
    {
        public Int16 Result;
        
        public bool FromBytes(byte[] bodyData)
        {
            Result = BitConverter.ToInt16(bodyData, 0);
            return true;
        }
    }

    public class GameBattingNtfPacket
    {
        public UInt64 RoomUserUniqueId;
        public sbyte SelectSide;
        
        public bool FromBytes(byte[] bodyData)
        {
            RoomUserUniqueId = BitConverter.ToUInt64(bodyData, 0);
            SelectSide = (SByte)bodyData[8];
            return true;
        }
    }
    
    
    //게임 결과 통보
    public class GameResultNtfPacket
    {
        public sbyte[] CardsBanker = new sbyte[3];
        public sbyte[] CardsPlayer = new sbyte[3];
        public sbyte PlayerScore;
        public sbyte BankerScore;
        public sbyte Result;
        public Int64 statusChangeCompletionMillSec;

        public bool FromBytes(byte[] bodyData)
        {
            CardsBanker[0] = (SByte)bodyData[0];
            CardsBanker[1] = (SByte)bodyData[1];
            CardsBanker[2] = (SByte)bodyData[2];
            CardsPlayer[0] = (SByte)bodyData[3];
            CardsPlayer[1] = (SByte)bodyData[4];
            CardsPlayer[2] = (SByte)bodyData[5];
            PlayerScore = (SByte)bodyData[6];
            BankerScore = (SByte)bodyData[7];
            Result = (SByte)bodyData[8];
            statusChangeCompletionMillSec = BitConverter.ToInt64(bodyData, 9);
            return true;
        }
    }
    
    //방장 변경 알림
    public class RoomMasterChangeNtfPacket
    {
        public UInt64 MasterUserUniqueId;

        public bool FromBytes(byte[] bodyData)
        {
            MasterUserUniqueId = BitConverter.ToUInt64(bodyData, 0);
            return true;
        }
    }
    
    ///[개별 상황 알려줌]
    public class RoomUserInfoNtfPacket
    {
        public int Dollar;
        public int Plays;
        public int Win;
        public int Lose;
        
        public bool FromBytes(byte[] bodyData)
        {
            Dollar = BitConverter.ToInt32(bodyData, 0);
            Plays = BitConverter.ToInt32(bodyData, 4);
            Win = BitConverter.ToInt32(bodyData, 8);
            Lose = BitConverter.ToInt32(bodyData, 12);
            return true;
        }
    }
}
