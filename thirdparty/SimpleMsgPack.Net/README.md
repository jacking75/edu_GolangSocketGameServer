SimpleMsgPack.Net
=================

MessagePack implementation for C# / msgpack.org[C#] 

Binary files distributed via the NuGet package [SimpleMsgPack](http://www.nuget.org/packages/SimpleMsgPack/).

It's like JSON but small and fast.

```
unit Owner: D10.Mofen
contact:
       qq:185511468, 
    email:ymofen@diocp.org
	homepage:www.diocp.org
if you find any bug, please contact me!
```

Works with
--------
  .NET Framework 4.x
  
  
### Code Example
```C#

    MsgPack msgpack = new MsgPack();
    msgpack.ForcePathObject("p.name").AsString = "张三";
    msgpack.ForcePathObject("p.age").AsInteger = 25;
    msgpack.ForcePathObject("p.datas").AsArray.Add(90);
    msgpack.ForcePathObject("p.datas").AsArray.Add(80);
    msgpack.ForcePathObject("p.datas").AsArray.Add("李四");
    msgpack.ForcePathObject("p.datas").AsArray.Add(3.1415926);

    // pack file
    msgpack.ForcePathObject("p.filedata").LoadFileAsBytes("C:\\a.png");

    // pack msgPack binary
    byte[] packData = msgpack.Encode2Bytes();

    MsgPack unpack_msgpack = new MsgPack();
	
    // unpack msgpack
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

    // unpack filedata 
    unpack_msgpack.ForcePathObject("p.filedata").SaveBytesToFile("C:\\b.png");
    Console.Read();
