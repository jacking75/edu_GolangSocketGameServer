# golang_socketGameServer_codelab
>>> Needs refactoring and documentation cleanup

[한국어 문서](README_kr.md)

- A hands-on codelab for building real-time socket game servers with Go.
- You learn how to build a server by typing out each server's original code yourself, one step at a time.
    - As you code, you get an explanation of how and why that code is implemented the way it is.

**There may be bugs**. Finding and fixing them is considered part of the learning process ^^;.


## Purpose
- The goal is to build up the skills and experience needed to write socket servers in Go.
- Aimed at people with little or no prior experience writing socket servers in Go.
- Instead of using Go's raw socket API, we use a "fake" (knock-off?) of goHiperNet, the Go network library built by the central-server team.
    - This library only shares its API with the real goHiperNet — the internal implementation is completely different.
    - It's more than good enough for learning purposes.
    - If you want to learn how to build a server from scratch using Go's socket API directly, ask separately.
- The exercise proceeds step by step, and each step takes a different amount of time.
    - No single step should take more than 3 hours, even the longer ones.
    - Covering too much in one sitting makes it hard to review later.


## Prerequisites
- One laptop per person (Windows or macOS)
- The latest version of the Go SDK
- The latest version of GoLand
- Basic knowledge of Go syntax
    - It's fine if you've never coded before at all.
- (Optional) If you want to test a server by connecting a client to it: Windows + Visual Studio 2022 or later (with the .NET 8.0 SDK)
    - You can still follow the exercise just by reading and building the server code, without a client.


## What's in this repository

This repository is roughly made up of three parts: "① the servers written in Go", "② the fake network library shared by all of the servers", and "③ a C# client you can use to connect to a server and try it out."

| Directory | What it is | Language |
|---|---|---|
| `gohipernetFake` | The network library shared by every server (the goHiperNet knock-off). All the code that talks to raw sockets lives here. | Go |
| `echoServer` | The simplest server. It just echoes back whatever it receives. The first step of the exercise. | Go |
| `chatServer` | A chat server with the concept of rooms. Packet processing happens on a single goroutine. | Go |
| `chatServer2` | An evolution of `chatServer` that processes packets across multiple goroutines. **As currently committed, this one does not build** (see the [chatServer2](#chatserver2) section below). | Go |
| `chatServer_msgpack` | Nearly identical to `chatServer`, but packets are exchanged using msgpack, a JSON-like binary format. | Go |
| `baccaratServer` | `chatServer` with baccarat card-game logic layered on top. | Go |
| `csharp_test_client` | A Windows GUI client for testing connections to `chatServer` / `chatServer2` / `baccaratServer` (room entry, chat, and relay features only). | C# (.NET 8, WinForms) |
| `csharp_test_client_msgpack` | A dedicated test client for `chatServer_msgpack`. | C# (.NET 8, WinForms) |
| `bin` | Sample batch files (`run_xxx.bat`) for running a built server `.exe`, plus a logger config file. | - |
| `thirdparty/SimpleMsgPack.Net` | An open-source library the C# client uses to handle msgpack (with some code modified). | C# |
| `lib_opensource` | Archives of other open-source Go network libraries, kept for reference. Not used in the exercise. | - |

Each Go server directory (`echoServer`, `chatServer`, `chatServer2`, `chatServer_msgpack`, `baccaratServer`) is a **fully independent Go module in its own right** (each has its own `go.mod`). That means you can open any one folder and follow along with just that folder — code that overlaps between servers (session management, room management, etc.) is intentionally copy-pasted into each folder for the sake of the exercise, rather than pulled out into a shared library. The one exception is `gohipernetFake`: all 5 servers pull it in via a `replace` directive in their `go.mod` pointing at the relative path `../gohipernetFake`.


## Learning order

Here's the recommended order. Each step assumes you already understand the code/concepts from the previous one, so it's best not to skip ahead.

1. Finish **[Prerequisites](#prerequisites)**, then read the **[Packet header](#packet-header)** section once to get a feel for the packet format these servers use.
2. Watch the "Exercise goals and approach" video under **[Explainer videos](#explainer-videos)**.
3. **`echoServer`** — start with the simplest server. Watch the "Building an Echo Server" video alongside it.
4. **`chatServer`** — the first "real" game-server structure, with the concept of rooms and login. Watch the "Chat server code walkthrough" video alongside it. **Understanding this server well is a prerequisite for every step after it.**
5. From here, the path splits depending on what you're interested in (the order among these doesn't matter):
    - **Interested in game logic → `baccaratServer`** (chatServer + card-game logic)
    - **Interested in multi-goroutine concurrency → `chatServer2`** (chatServer parallelized across multiple goroutines). Note that the current code doesn't build — read [that section](#chatserver2) first.
    - **Interested in serialization formats → `chatServer_msgpack`** (chatServer + msgpack serialization)


## Packet header
The packet header is 5 bytes total:
- Total packet size (2 bytes, header + body combined) + packet ID (2 bytes) + packet type (1 byte)


## Explainer videos
- [Exercise goals and approach](https://youtu.be/zR_zcY7SXio )
- [Building an Echo Server](https://youtu.be/OSiwcsPAO2o )
- [Chat server code walkthrough](https://youtu.be/2rppKuW-wQg )


## How to build and run the Go servers

`echoServer`, `chatServer`, `chatServer_msgpack`, and `baccaratServer` are all built and run the same way (`chatServer2` is excluded because it currently doesn't build — see [the section below](#chatserver2)). There are two ways to run them.

### Method 1) Run directly from a terminal (simplest)

Go into the server's directory and run `go run .` followed by options (flags) for things like the bind address and max session count. For example, for `echoServer`:

```bash
cd echoServer
go run . -c_IsTcp4Addr=true -c_BindAddress=127.0.0.1:11021 -c_MaxSessionCount=200 -c_MaxPacketSize=1024
```

Once it's running you'll see logs like the ones below, and the process will just sit there waiting for connections (that's normal, not a hang — press `Ctrl+C` to stop it).

```
[ info ] tcpServerStart - Start
2026/08/27 23:53:28 Server Listen ...
```

> **Careful**: if you leave out `-c_MaxSessionCount`, it defaults to 0, and while the server does start, **nobody can connect to it** (internally, the pool of connectable session slots gets initialized to size zero). Always pass the flags as shown above.

The exact flags differ slightly per server, but the batch files in the `bin` directory already have a verified set of options for each server — just copy those.

| Server | Run command (run from inside that server's directory) |
|---|---|
| echoServer | `go run . -c_IsTcp4Addr=true -c_BindAddress=127.0.0.1:11021 -c_MaxSessionCount=200 -c_MaxPacketSize=1024` |
| chatServer | `go run . -c_IsTcp4Addr=true -c_BindAddress=127.0.0.1:11021 -c_MaxSessionCount=200 -c_MaxPacketSize=1024` |
| chatServer_msgpack | `go run . -c_IsTcp4Addr=true -c_BindAddress=127.0.0.1:11021 -c_MaxSessionCount=200 -c_MaxPacketSize=1024` |
| baccaratServer | `go run . -c_IsTcp4Addr=true -c_BindAddress=127.0.0.1:11021 -c_MaxSessionCount=200 -c_MaxPacketSize=1024` |

By default, all of them listen on `127.0.0.1:11021` (port 11021 on your own machine). If you want to run two or more servers at once, just give `-c_BindAddress` a different port number for each.

Instead of `go run .`, you can also build an executable with `go build`. The output binary is automatically named after the directory (building inside `echoServer` produces `echoServer.exe`). The `run_xxx.bat` files in the `bin` directory are meant to run exactly that kind of `.exe`, so if you move your built `.exe` into `bin`, you can use those batch files as-is.

### Method 2) Run from GoLand

Open a server directory (e.g. `echoServer`) as a project in GoLand, create a Run Configuration that runs `main.go`, and enter the flags from the table above (`-c_IsTcp4Addr=true -c_BindAddress=...` and the rest) as-is in **Program arguments**. This behaves the same as running it from a terminal. Use this method if you need debugging (breakpoints, etc.).


## Using the test client

To see with your own eyes whether a server is working correctly, use the C# GUI client included in this repository. You can't test with a text-based tool like telnet — the packets aren't human-readable text, they're the binary format described in [Packet header](#packet-header).

1. Open `csharp_test_client/csharp_test_client.sln` in Visual Studio. (When testing `chatServer_msgpack`, open `csharp_test_client_msgpack/csharp_test_client_msgpack.sln` instead.)
2. Build and run it (F5). A WinForms window appears.
3. The IP field defaults to `0.0.0.0`, so either check the "localhost" checkbox or type `127.0.0.1` yourself. Set Port to match the server (the default is `11021`; if you started the server on a different port, use that one).
4. Click **접속하기** ("Connect") to connect to the server.
5. What you can test next depends on the server:
    - **echoServer**: click "echo 보내기" ("Send echo") and check the log to see whatever you sent come right back.
    - **chatServer / baccaratServer (and chatServer2 too, once you've fixed it up and gotten it running)**: fill in the UserID/PW fields and click **Login** → enter a room number and click **방 입장** ("Enter room") → fill in the chat field and click **chat**. Open two client windows and put them both in the same room, and you'll see a chat message sent from one show up in the other.
    - **baccaratServer**: you can test the chat/room features above with this client, but there's no UI in it for actually playing the baccarat game itself (betting, etc.) — you'd need to add that yourself. See `baccaratServer/docs/*.puml` (PlantUML sequence diagrams) for the game flow.
    - **chatServer_msgpack**: use `csharp_test_client_msgpack` and follow the same steps.


## echoServer
- Directory: echoServer
- Use GoLand to write, build, and debug your Go program.
- Very small in scope.
- Running it: see [the table above](#method-1-run-directly-from-a-terminal-simplest). Test it with the "echo 보내기" ("Send echo") button in `csharp_test_client`.


## chatServer
- Directory: chatServer
- A chat server with the concept of rooms.
- Packet requests are all processed on a single goroutine (thread).
- About 3-4x the size of echoServer.
- Running it: see [the table above](#method-1-run-directly-from-a-terminal-simplest). Test it with `csharp_test_client`: log in → enter a room → chat.

### Additional features implemented
- 1:1 whispers
- Room invitations



## baccaratServer
- Directory: baccaratServer
- An online version of the gambling card game baccarat.
    - Baccarat rules: https://namu.wiki/w/%EB%B0%94%EC%B9%B4%EB%9D%BC
- This is chatServer with baccarat game logic layered on top, so you need to understand chatServer first.
- Running it: see [the table above](#method-1-run-directly-from-a-terminal-simplest). See `baccaratServer/docs/*.puml` for the game flow. There's no betting UI in the client, so you can only visually verify room entry/chat via `csharp_test_client`.

### Additional features implemented
- Game server scale-out
- Integration with an API server (http)
    - Assigning a user to a specific game server
    - Matchmaking

> The "Additional features implemented" items above are exercise goals — they aren't implemented in the committed code yet.



## chatServer2
- Directory: chatServer2
- A chat server with the concept of rooms.
- Packet requests are processed across N goroutines (threads).
    - Because packets are processed on multiple goroutines, you need to be careful about synchronizing shared objects.
- There's a lot of overlap with chatServer's code, so you need to understand chatServer first.

> **This server currently does not build.** Several symbols — `NTELIB_LOG_INFO`, `NTELIB_LOG_ERROR`, `NetLib_IsRunningServer`, `NetLib_GetCurrnetUnixTime`, and more — are used throughout the code but aren't defined anywhere in this repository. So you can't run it directly using the [instructions above](#how-to-build-and-run-the-go-servers). It's still usable for reading through the code and learning "how was this split up across multiple goroutines," but to actually run it you'll first need to fill in those missing symbols.

### Additional features implemented
- Redis integration
- Integration with an API server (http)
    - Logging in through the API server

> The "Additional features implemented" items above are exercise goals — they aren't implemented in the committed code yet.


## chatServer using msgpack
- Directory: chatServer_msgpack
- Client directory: csharp_test_client_msgpack
- The server and client exchange packet data over the network using msgpack.
    - [Go](https://github.com/vmihailenco/msgpack )
    - [C#](https://github.com/ymofen/SimpleMsgPack.Net )
        - There were parts where the data format didn't match the Go library's, so the code was modified.
        - That code lives in the thirdparty/SimpleMsgPack.Net directory.
- Running it: see [the table above](#method-1-run-directly-from-a-terminal-simplest). Test it with `csharp_test_client_msgpack`.


## References
- [YouTube: Learning Golang TCP Socket Server Programming from Open-Source Code](https://youtu.be/boDo8JoyHuo )
