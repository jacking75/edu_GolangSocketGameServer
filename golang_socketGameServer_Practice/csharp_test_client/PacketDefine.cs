using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace csharp_test_client
{
    class PacketDef
    {
        public const Int16 PACKET_HEADER_SIZE = 5;
        public const int MAX_USER_ID_BYTE_LENGTH = 16;
        public const int MAX_USER_PW_BYTE_LENGTH = 16;
    }

    public enum PACKET_ID : ushort
    {
        PACKET_ID_ECHO = 101,

        // Ping(Heart-beat)
        PACKET_ID_PING_REQ = 201,
        PACKET_ID_PING_RES = 202,

        PACKET_ID_ERROR_NTF = 203,


        // 로그인
        PACKET_ID_LOGIN_REQ = 701,
        PACKET_ID_LOGIN_RES = 702,
                

        PACKET_ID_ROOM_ENTER_REQ = 721,
        PACKET_ID_ROOM_ENTER_RES = 722,
        PACKET_ID_ROOM_USER_LIST_NTF = 723,
        PACKET_ID_ROOM_NEW_USER_NTF = 724,

         PACKET_ID_ROOM_LEAVE_REQ = 726,
         PACKET_ID_ROOM_LEAVE_RES = 727,
         PACKET_ID_ROOM_LEAVE_USER_NTF = 728,

         PACKET_ID_ROOM_CHAT_REQ = 731,
         PACKET_ID_ROOM_CHAT_RES = 732,
         PACKET_ID_ROOM_CHAT_NOTIFY = 733,

         PACKET_ID_ROOM_WHISPER_REQ       = 751,
         PACKET_ID_ROOM_WHISPER_RES       = 752,
         PACKET_ID_ROOM_WHISPER_NOTIFY    = 753,
         
         PACKET_ID_ROOM_RELAY_REQ = 741,
         PACKET_ID_ROOM_RELAY_NTF = 742,
         
         // 바카라
         PACKET_ID_GAME_START_REQ = 761,
         PACKET_ID_GAME_START_RES = 762,
         PACKET_ID_GAME_START_NTF = 763,

         PACKET_ID_GAME_BATTING_REQ = 771,
         PACKET_ID_GAME_BATTING_RES = 772,
         PACKET_ID_GAME_BATTING_NTF = 773,

         PACKET_ID_GAME_RESULT_NTF = 774,
         PACKET_ID_GAME_USER_INFO_NTF = 775,
         
         PACKET_ID_ROOM_MASTER_CHANGE_NTF = 781
    }


    public enum ERROR_CODE : Int16
    {
        ERROR_NONE = 0,
        
        ERROR_CODE_USER_NOT_IN_ROOM = 54,

        ERROR_CODE_USER_MGR_INVALID_USER_UNIQUEID = 112,

        ERROR_CODE_PUBLIC_CHANNEL_IN_USER = 114,

        ERROR_CODE_PUBLIC_CHANNEL_INVALIDE_NUMBER = 115,
        
        ERROR_CODE_NONE = 1,
        

        ERROR_CODE_ROOM_NOT_IN_USER                    = 121,
        ERROR_CODE_ROOM_INVALIDE_NUMBER                = 122,

        ERROR_CODE_ENTER_ROOM_ALREADY = 131,
        ERROR_CODE_ENTER_ROOM_GAME_IN_PROGRESS = 132,
        ERROR_CODE_ENTER_ROOM_INVALID_USER_ID = 133,
        ERROR_CODE_ENTER_ROOM_USER_FULL = 134,
        ERROR_CODE_ENTER_ROOM_DUPLCATION_USER         = 135,
        ERROR_CODE_ENTER_ROOM_INVALID_SESSION_STATE    = 136,
        ERROR_CODE_ENTER_ROOM_AUTO_ROOM_NUMBER = 137,

        ERROR_CODE_LEAVE_ROOM_INTERNAL_INVALID_USER    = 141,

        ERROR_CODE_ROOM_CHAT_CHAT_MSG_LEN    = 151,

        ERROR_CODE_ROOM_RELAY_FAIL_DECPDING    = 161,

        ERROR_CODE_ROOM_GAME_START_INVALID_ROOM_STATE    = 171,
        ERROR_CODE_ROOM_GAME_START_NOT_ENOUGH_MEMBERS    = 172,
        ERROR_CODE_ROOM_GAME_START_NOT_MASTER    = 173,

        ERROR_CODE_ROOM_GAME_BATTING_FAIL_PACKET    = 181,
        ERROR_CODE_ROOM_GAME_BATTING_INVALID_ROOM_STATE    = 182,
        ERROR_CODE_ROOM_GAME_BATTING_INVALID_BAT_SELECT    = 183,
        ERROR_CODE_ROOM_GAME_BATTING_SAME_BAT_SELECT    = 184
    }

    public enum GAMECASE : Int32
    {
        NONE = 0,
        PLAYER = 1,
        BANKER = 2,
        TIE = 3
    }
    
    
    public static class Card
    {
        public static readonly string[] CARD =
        {
            "(Spade)A","(Spade)2","(Spade)3","(Spade)4","(Spade)5","(Spade)6","(Spade)7","(Spade)8","(Spade)9","(Spade)10","(Spade)J","(Spade)Q","(Spade)K",
            "(Diamond)A","(Diamond)2","(Diamond)3","(Diamond)4","(Diamond)5","(Diamond)6","(Diamond)7","(Diamond)8","(Diamond)9","(Diamond)10","(Diamond)J","(Diamond)Q","(Diamond)K",
            "(Club)A","(Club)2","(Club)3","(Club)4","(Club)5","(Club)6","(Club)7","(Club)8","(Club)9","(Club)10","(Club)J","(Club)Q","(Club)K",
            "(Heart)A","(Heart)2","(Heart)3","(Heart)4","(Heart)5","(Heart)6","(Heart)7","(Heart)8","(Heart)9","(Heart)10","(Heart)J","(Heart)Q","(Heart)K"
        };

        public static string get(Int32 index)
        {
            if (index == -1)
            {
                return "NONE";
            }

            return CARD[index];
        }
    }
}
