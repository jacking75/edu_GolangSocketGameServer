using SimpleMsgPack;
using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace SimpleMsgPackTester
{
    class Program
    {
        static void Test1()
        {
            MsgPack msgpack = new MsgPack();
            msgpack.ForcePathObject("p.name").AsString = "张三一二三四五六七八九十";
            msgpack.ForcePathObject("p.age").AsInteger = 132123456874125;
            msgpack.ForcePathObject("p.datas").AsArray.Add(90);
            msgpack.ForcePathObject("p.datas").AsArray.Add(80);
            msgpack.ForcePathObject("p.datas").AsArray.Add("李四");
            msgpack.ForcePathObject("p.datas").AsArray.Add(3.1415926);
            msgpack.ForcePathObject("Game.iGameID").AsInteger = 1;

            // 可以直接打包文件数据
            // msgpack.ForcePathObject("p.filedata").LoadFileAsBytes("C:\\a.png");

            // 打包成msgPack协议格式数据
            byte[] packData = msgpack.Encode2Bytes();

            FileStream fs = new FileStream("d:\\simplemsgpack.dat", FileMode.Append);
            fs.Write(packData, 0, packData.Length);
            fs.Close();

            //Console.WriteLine("msgpack序列化数据:\n{0}", BytesTools.BytesAsHexString(packData));

            MsgPack unpack_msgpack = new MsgPack();
            // 从msgPack协议格式数据中还原
            unpack_msgpack.DecodeFromBytes(packData);

            System.Console.WriteLine("name:{0}, age:{1}",
                  unpack_msgpack.ForcePathObject("p.name").AsString,
                  unpack_msgpack.ForcePathObject("p.age").AsInteger);

            Console.WriteLine("==================================");
            System.Console.WriteLine("use index property, Length{0}:{1}",
                  unpack_msgpack.ForcePathObject("p.datas").AsArray.Length,
                  unpack_msgpack.ForcePathObject("p.datas").AsArray[0].AsString
                  );

            Console.WriteLine("==================================");
            Console.WriteLine("use foreach statement:");
            foreach (MsgPack item in unpack_msgpack.ForcePathObject("p.datas"))
            {
                Console.WriteLine(item.AsString);
            }

            Console.WriteLine(unpack_msgpack.ForcePathObject("Game.iGameID").AsInteger);

            // unpack filedata 
            //unpack_msgpack.ForcePathObject("p.filedata").SaveBytesToFile("C:\\b.png");
            Console.Read();
        }

        static void Test2()
        {
            MsgPack msgpack = new MsgPack();
            msgpack.AsString = "张三一二三四五六七八九十";

            // 打包成msgPack协议格式数据
            byte[] packData = msgpack.Encode2Bytes();

            FileStream fs = new FileStream("d:\\simplemsgpack11.dat", FileMode.Append);
            fs.Write(packData, 0, packData.Length);
            fs.Close();

        }

        static void Test3()
        {
            MsgPack msgpack = new MsgPack();
            msgpack.SetAsUInt64(UInt64.MaxValue - 1);

            // 打包成msgPack协议格式数据
            byte[] packData = msgpack.Encode2Bytes();

            MsgPack unpack_msgpack = new MsgPack();
            // 从msgPack协议格式数据中还原
            unpack_msgpack.DecodeFromBytes(packData);

            Console.WriteLine(unpack_msgpack.GetAsUInt64());

            Console.Read();
        }
        static void Main(string[] args)
        {

            Test3();
        }
    }
}
