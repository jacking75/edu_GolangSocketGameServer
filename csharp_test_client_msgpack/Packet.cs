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
    

    public class PacketHeader
    {
        static public void SetHeadInfo(byte[] packetData, UInt16 packetId, UInt16 size)
        {
            Buffer.BlockCopy(BitConverter.GetBytes(size), 0, packetData, 10, 2);
            Buffer.BlockCopy(BitConverter.GetBytes(packetId), 0, packetData, 12, 2);
        }
    }

    public class ErrorNtfPacket
    {
        public ERROR_CODE ErrorCode;

        public bool FromBytes(byte[] bodyData)
        {
            var unpack_msgpack = new SimpleMsgPack.MsgPack();
            unpack_msgpack.DecodeFromBytes(bodyData);
            
            ErrorCode = (ERROR_CODE)unpack_msgpack.ForcePathObject("ErrorCode").AsInteger;
            return true;
        }
    }
    

    public class LoginReqPacket
    {
        public string UserID;
        public string UserPW;
                
        public byte[] ToBytes()
        {
            var msgpack = new SimpleMsgPack.MsgPack();
            msgpack.ForcePathObject("UserID").AsString = UserID;
            msgpack.ForcePathObject("UserPW").AsString = UserPW;
            byte[] packData = msgpack.Encode2Bytes();
            return packData;
        }
    }

    public class LoginResPacket
    {
        public Int64 Result;

        public bool FromBytes(byte[] bodyData)
        {
            var unpack_msgpack = new SimpleMsgPack.MsgPack();
            unpack_msgpack.DecodeFromBytes(bodyData);
            Result = unpack_msgpack.ForcePathObject("Result").AsInteger;
            return true;
        }
    }


    public class RoomEnterReqPacket
    {
        public Int64 RoomNumber;
        public void SetValue(int roomNumber)
        {
            RoomNumber = roomNumber;
        }

        public byte[] ToBytes()
        {
            var msgpack = new SimpleMsgPack.MsgPack();
            msgpack.ForcePathObject("RoomNumber").AsInteger = RoomNumber;
            byte[] packData = msgpack.Encode2Bytes();
            return packData;
        }
    }

    public class RoomEnterResPacket
    {
        public Int64 Result;
        public UInt64 RoomUserUniqueId;

        public bool FromBytes(byte[] bodyData)
        {
            var unpack_msgpack = new SimpleMsgPack.MsgPack();
            unpack_msgpack.DecodeFromBytes(bodyData);
            Result = unpack_msgpack.ForcePathObject("Result").AsInteger;
            RoomUserUniqueId = (UInt64)unpack_msgpack.ForcePathObject("RoomUserUniqueId").AsInteger;
            return true;
        }
    }

    public class RoomUserListNtfPacket
    {
        public Int64 UserCount = 0;
        public List<UInt64> UserUniqueIdList = new List<UInt64>();
        public List<string> UserIDList = new List<string>();

        public bool FromBytes(byte[] bodyData)
        {
            var unpack_msgpack = new SimpleMsgPack.MsgPack();
            unpack_msgpack.DecodeFromBytes(bodyData);

            UserCount = unpack_msgpack.ForcePathObject("UserCount").AsInteger;

            foreach (SimpleMsgPack.MsgPack item in unpack_msgpack.ForcePathObject("UniqueId"))
            {
                UserUniqueIdList.Add((UInt64)item.GetAsInteger());
            }

            foreach (SimpleMsgPack.MsgPack item in unpack_msgpack.ForcePathObject("ID"))
            {
                UserIDList.Add(item.AsString);
            }

            return true;
        }
    }

    public class RoomNewUserNtfPacket
    {
        public UInt64 UserUniqueId;
        public string UserID;

        public bool FromBytes(byte[] bodyData)
        {
            var unpack_msgpack = new SimpleMsgPack.MsgPack();
            unpack_msgpack.DecodeFromBytes(bodyData);

            UserUniqueId = (UInt64)unpack_msgpack.ForcePathObject("UniqueId").AsInteger;

            UserID = unpack_msgpack.ForcePathObject("ID").AsString;

            return true;
        }
    }


    public class RoomChatReqPacket
    {
        public string Msg;
                
        public byte[] ToBytes()
        {
            var msgpack = new SimpleMsgPack.MsgPack();
            msgpack.ForcePathObject("Msg").AsString = Msg;
            byte[] packData = msgpack.Encode2Bytes();
            return packData;
        }
    }

    public class RoomChatResPacket
    {
        public Int64 Result;
        
        public bool FromBytes(byte[] bodyData)
        {
            var unpack_msgpack = new SimpleMsgPack.MsgPack();
            unpack_msgpack.DecodeFromBytes(bodyData);
            Result = unpack_msgpack.ForcePathObject("Result").AsInteger;
            return true;
        }
    }

    public class RoomChatNtfPacket
    {
        public UInt64 UserUniqueId;
        public string Msg;

        public bool FromBytes(byte[] bodyData)
        {
            var unpack_msgpack = new SimpleMsgPack.MsgPack();
            unpack_msgpack.DecodeFromBytes(bodyData);
            
            UserUniqueId = (UInt64)unpack_msgpack.ForcePathObject("UserUniqueId").AsInteger;
            Msg = unpack_msgpack.ForcePathObject("Msg").AsString;
            return true;
        }
    }


     public class RoomLeaveResPacket
    {
        public Int64 Result;
        
        public bool FromBytes(byte[] bodyData)
        {
            var unpack_msgpack = new SimpleMsgPack.MsgPack();
            unpack_msgpack.DecodeFromBytes(bodyData);
            Result = unpack_msgpack.ForcePathObject("Result").AsInteger;
            return true;
        }
    }

    public class RoomLeaveUserNtfPacket
    {
        public Int64 UserUniqueId;

        public bool FromBytes(byte[] bodyData)
        {
            var unpack_msgpack = new SimpleMsgPack.MsgPack();
            unpack_msgpack.DecodeFromBytes(bodyData);
            UserUniqueId = unpack_msgpack.ForcePathObject("UserUniqueId").AsInteger;
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
}
