# golang_socketGameServer_codelab

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
- Why use a "game server" as the example?
    - A game server is a classic case of socket programming: a connection stays open while multiple clients exchange state with each other in real time.
    - Once you're dealing with rooms, logins, and several concurrent users interacting at once, you naturally run into problems (concurrency, state management, ordering) that a typical web API server — one request in, one response out, done — rarely forces you to face.


## Prerequisites
- One laptop per person (Windows or macOS)
- The latest version of the Go SDK
- The latest version of GoLand
- Basic knowledge of Go syntax
    - It's fine if you've never coded before at all.
- (Optional) If you want to test a server by connecting a client to it: Windows + Visual Studio 2022 or later (with the .NET 8.0 SDK)
    - You can still follow the exercise just by reading and building the server code, without a client.
    - The client is written in C# (a Windows-only WinForms app), so it's hard to run on macOS. macOS users should focus on the server-side code, and if you really need to test a live connection, either borrow a Windows machine or, as a nice exercise in itself, write your own text- or GUI-based client later.


## Concepts worth knowing beforehand

Here's a quick glossary of terms you'll run into a lot when you first start reading this repo's code. Feel free to skip this if you're already comfortable with them.

- **Socket / TCP connection**: the "connection" a client and server establish to exchange data over the network. Once connected, they can keep exchanging data over it until it's disconnected. Every server in this repo uses TCP.
- **Packet**: one unit of data exchanged between client and server. Each request — "log me in," "send this chat message" — corresponds to one packet. See [Packet header](#packet-header) for the exact format.
- **Packet ID**: a number that says what kind of packet this is (a login request, a chat request, etc.). Each server's `protocol/packetID.go` file has the full list.
- **Serialization / encoding & decoding**: turning an in-memory struct (e.g. `LoginReqPacket{UserID: "abc", PassWD: "123"}`) into a plain byte array you can send over the network (encoding), and turning received bytes back into a struct (decoding). The `chatServer`/`baccaratServer` family does this by hand with custom code; `chatServer_msgpack` uses a standard format called msgpack instead.
- **Goroutine**: Go's lightweight unit of concurrent execution. Similar to a thread in other languages, but cheap enough that you can spin up tens of thousands of them. `chatServer` processes every packet in order on a single goroutine, while `chatServer2` spreads that work across multiple goroutines — this is the biggest structural difference between the two.
- **Session**: the unit the server uses internally to track one connected client. A session is created on connect, and state like "is this client logged in" and "which room is it in" lives on the session. In the code, the `connectedSessions` package handles this.
- **Room**: the logical space where chat or gameplay actually happens. After logging in, a client has to "enter" a specific room number before it can chat or play with anyone else in that room. In the code, the `roomPkg` package handles this.
- **Little-endian**: the byte order this protocol uses to represent numbers. For example, when writing the integer 700 as 2 bytes, the low-order byte comes first, not the high-order one. See the worked example in [Packet header](#packet-header) for a concrete illustration.


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

### What's inside one `chatServer`-family server

Every server except echoServer is organized the same way, folder-wise. Once you know this pattern you can navigate any of them.

| Folder / file | Role |
|---|---|
| `main.go` | The entry point. Reads command-line flags (`-c_BindAddress`, etc.) and starts the server. |
| `chatServer.go` / `baccaratServer.go` (a file named after the server itself) | Initializes the server and registers callbacks with `gohipernetFake` — "call this function when this event happens." |
| `distributePacket.go` | The first place a packet lands once it arrives from a client. It looks at the packet ID and decides whether it's a login, something that should go to a room, etc. — the "dispatcher." |
| `connectedSessions/` | Tracks the state of each connected client (session): whether it's logged in, which room it's currently in, and so on. |
| `protocol/` | Packet formats (structs), the list of packet IDs, the list of error codes, and the struct ↔ byte encoding/decoding code. |
| `roomPkg/` | Manages the list of rooms, and everything that happens inside a room (entering, leaving, chatting, playing, etc.). A file starting with `room_Packet` handles "what the room does when it receives that particular packet." |

`baccaratServer` adds `roomPkg/Game.go` on top of this — the baccarat rules themselves (card shuffling, betting, deciding the winner, etc.).


## Learning order

Here's the recommended order. Each step assumes you already understand the code/concepts from the previous one, so it's best not to skip ahead. Each step also notes what to pay attention to.

1. Finish **[Prerequisites](#prerequisites)**, then read **[Concepts worth knowing beforehand](#concepts-worth-knowing-beforehand)** and **[Packet header](#packet-header)** once to get a feel for the packet format these servers use. It's fine if it doesn't all click yet — it'll make a lot more sense once you come back to it after reading the actual echoServer/chatServer code.
2. Watch the "Exercise goals and approach" video under **[Explainer videos](#explainer-videos)**.
3. **`echoServer`** — start with the simplest server. Watch the "Building an Echo Server" video alongside it.
    - What to notice: the minimal flow of opening a socket, accepting a client connection, and echoing back whatever data arrives. There are only two files (`echoServer.go`, `main.go`), so it's a good way to see the whole picture at once.
4. **`chatServer`** — the first "real" game-server structure, with the concept of rooms and login. Watch the "Chat server code walkthrough" video alongside it. **Understanding this server well is a prerequisite for every step after it.**
    - Suggested reading order: `main.go` (startup) → `distributePacket.go` (how packets get dispatched) → `connectedSessions/session.go` (what state one session holds) → `roomPkg/room.go` (what state one room holds) → `roomPkg/room_PacketEnter.go` (how a logged-in user enters a room) → `roomPkg/room_PacketChat.go` (how a chat message reaches everyone else in the room).
    - What to notice: how the single scenario "log in → enter a room → chat" is stitched together across several files. Try tracing, by hand, the order function calls happen in when one packet comes in.
5. From here, the path splits depending on what you're interested in (the order among these doesn't matter):
    - **Interested in game logic → `baccaratServer`** (chatServer + card-game logic). Start with `roomPkg/Game.go`.
    - **Interested in multi-goroutine concurrency → `chatServer2`** (chatServer parallelized across multiple goroutines). Note that the current code doesn't build — read [that section](#chatserver2) first.
    - **Interested in serialization formats → `chatServer_msgpack`** (chatServer + msgpack serialization). Put it side by side with `chatServer` and just diff `protocol/packet.go` between the two — that's the fastest way to see what changed.
6. After finishing each step, [run that server yourself](#how-to-build-and-run-the-go-servers) and [connect to it with the client](#using-the-test-client) to watch the code you just read actually do something. It sticks a lot better than reading alone.


## Packet header
The packet header is 5 bytes total, and every packet starts with it.

| Field | Size | Meaning |
|---|---|---|
| Total packet size | 2 bytes | the size, in bytes, of this whole packet — the 5-byte header plus the body |
| Packet ID | 2 bytes | what kind of packet this is (a login request, a chat request, etc.) |
| Packet type | 1 byte | 0 = normal; other values are reserved for things like compression/encryption (not used by this repo's exercise code) |

Numbers are all written **little-endian** (the low-order byte comes first).

### A worked example

Say you're sending a login request (`PACKET_ID_LOGIN_REQ`, packet ID 701). `chatServer`'s login request body is a fixed size — "16 bytes of UserID + 16 bytes of password" (padded with zeros if shorter). So the whole packet is 5 (header) + 16 + 16 = **37 bytes**.

Laid out as bytes, it looks like this (hex, leftmost byte sent first):

```
25 00   BD 02   00   [16 bytes of UserID]   [16 bytes of password]
└─┬──┘ └─┬──┘  └┬┘
  │      │      └ packet type = 0
  │      └ packet ID = 701 (0x02BD, little-endian so it's BD, then 02)
  └ total size = 37 (0x0025, little-endian so it's 25, then 00)
```

The server just reads these bytes straight off the TCP socket and interprets them — this is raw binary, not human-readable text, so you can't test these servers with a text tool like telnet (that's exactly why you need the [test client](#using-the-test-client)).


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

Here's what each flag means:

| Flag | Meaning |
|---|---|
| `-c_IsTcp4Addr` | Whether to use an IPv4 address. Leave it `true` unless you have a reason not to. |
| `-c_BindAddress` | The IP and port the server binds to and listens on. `127.0.0.1` means "this machine, itself." |
| `-c_MaxSessionCount` | The maximum number of clients that can be connected at once. This many session slots are pre-allocated. |
| `-c_MaxPacketSize` | The maximum size, in bytes, of a single packet. |

> **Careful**: if you leave out `-c_MaxSessionCount`, it defaults to 0, and while the server does start, **nobody can connect to it** (internally, the pool of connectable session slots gets initialized to size zero). Always pass the flags as shown above.

The exact flags differ slightly per server, but the batch files in the `bin` directory already have a verified set of options for each server — just copy those.

| Server | Run command (run from inside that server's directory) |
|---|---|
| echoServer | `go run . -c_IsTcp4Addr=true -c_BindAddress=127.0.0.1:11021 -c_MaxSessionCount=200 -c_MaxPacketSize=1024` |
| chatServer | `go run . -c_IsTcp4Addr=true -c_BindAddress=127.0.0.1:11021 -c_MaxSessionCount=200 -c_MaxPacketSize=1024` |
| chatServer_msgpack | `go run . -c_IsTcp4Addr=true -c_BindAddress=127.0.0.1:11021 -c_MaxSessionCount=200 -c_MaxPacketSize=1024` |
| baccaratServer | `go run . -c_IsTcp4Addr=true -c_BindAddress=127.0.0.1:11021 -c_MaxSessionCount=200 -c_MaxPacketSize=1024` |

By default, all of them listen on `127.0.0.1:11021` (port 11021 on your own machine). If you want to run two or more servers at once, just give `-c_BindAddress` a different port number for each (e.g. 11021, 11022).

One more thing worth knowing: `chatServer`/`chatServer_msgpack`/`baccaratServer` don't take room-related settings (room count, max users per room, etc.) as command-line flags — those are **hardcoded directly in `main.go`** (`RoomMaxCount=1000, RoomStartNum=0, RoomMaxUserCount=4`). In other words, room numbers 0 through 999 are valid, and each room holds at most 4 people. If you want to test with different numbers, just edit those values in the relevant `main.go` and rebuild.

Instead of `go run .`, you can also build an executable with `go build`. The output binary is automatically named after the directory (building inside `echoServer` produces `echoServer.exe`). The `run_xxx.bat` files in the `bin` directory are meant to run exactly that kind of `.exe`, so if you move your built `.exe` into `bin`, you can use those batch files as-is.

### Method 2) Run from GoLand

Open a server directory (e.g. `echoServer`) as a project in GoLand, create a Run Configuration that runs `main.go`, and enter the flags from the table above (`-c_IsTcp4Addr=true -c_BindAddress=...` and the rest) as-is in **Program arguments**. This behaves the same as running it from a terminal. Use this method if you need debugging (breakpoints, etc.) — for example, set a breakpoint in `distributePacket.go`'s packet-dispatch code and try logging in from the client, and you can watch, line by line, how one packet gets handled inside the server.


## Using the test client

To see with your own eyes whether a server is working correctly, use the C# GUI client included in this repository. You can't test with a text-based tool like telnet — the packets aren't human-readable text, they're the binary format described in [Packet header](#packet-header).

1. Open `csharp_test_client/csharp_test_client.sln` in Visual Studio. (When testing `chatServer_msgpack`, open `csharp_test_client_msgpack/csharp_test_client_msgpack.sln` instead.)
2. Build and run it (F5). A WinForms window appears.
3. The IP field defaults to `0.0.0.0`, so either check the "localhost" checkbox or type `127.0.0.1` yourself. Set Port to match the server (the default is `11021`; if you started the server on a different port, use that one).
4. Click **접속하기** ("Connect") to connect to the server. If the log window shows a successful-connect message, you're good. If nothing happens or it disconnects immediately, see [Common issues](#common-issues).
5. What you can test next depends on the server:
    - **echoServer**: click "echo 보내기" ("Send echo") and check the log to see whatever you sent come right back. There's no login/room concept here, so that's the whole test.
    - **chatServer / baccaratServer (and chatServer2 too, once you've fixed it up and gotten it running)**: follow these steps.
        1. Fill in the UserID/PW fields and click **Login**. (UserID can't be empty; since there's no real database, any password passes without being checked. However, **logging in twice with the same UserID at the same time is blocked** — if you open multiple client windows, give each one a different UserID.)
        2. Enter a room number (0-999, see [the run-methods table above](#method-1-run-directly-from-a-terminal-simplest)) and click **방 입장** ("Enter room").
        3. Fill in the chat field and click **chat**.
        4. Open two client windows (each logged in with a different UserID) and put them both in the same room number — you'll see a chat message sent from one show up live in the other's log window. This is the most fun part of the exercise to watch happen.
    - **baccaratServer**: you can test the chat/room features above with this client, but there's no UI in it for actually playing the baccarat game itself (betting, etc.) — you'd need to add that yourself. See `baccaratServer/docs/*.puml` (PlantUML sequence diagrams — open them with something like the PlantUML extension for VS Code to see them rendered) for the game flow.
    - **chatServer_msgpack**: use `csharp_test_client_msgpack` and follow the same steps.


## Common issues

- **The server is up but the client can't connect / disconnects immediately.**
    - Make sure you didn't forget the `-c_MaxSessionCount` flag (see the [Careful] note above — if it defaults to 0, nobody can get in).
    - Make sure the server's and client's IP/port actually match (e.g. the server is on `127.0.0.1:11021` but the client is still trying `0.0.0.0:11021`).
- **You get an "address already in use" error.**
    - A server process from an earlier run is probably still bound to that port. Closing the terminal window doesn't always kill the process. Check Task Manager (Windows) for a leftover `echoServer.exe`/`chatServer.exe`-type process and end it, or just pick a different port with `-c_BindAddress` and run again.
- **The Login button fails.**
    - Check that the UserID field isn't empty.
    - Check whether another client is already logged in with that same UserID (duplicate logins are rejected) — use a different UserID per client window.
- **Entering a room fails.**
    - Check that the room number is in a valid range (0-999, for the chatServer family's default settings).
    - Each room only holds up to 4 people by default — a 5th attempt to join the same room will fail.
    - You can't enter a room before logging in.
- **`go run .` just keeps printing `go: downloading ...`.**
    - That's normal the first time you build — it's fetching dependencies over the network. Check your internet connection and give it a moment; once they're cached, subsequent runs start immediately.
    - Building `gohipernetFake` once on its own first (`cd gohipernetFake && go build ./...`) helps, since it prepares the local module the other servers reference.
    - This won't help with `chatServer2` — it can't be run this way because the code itself doesn't build yet. See [that section](#chatserver2).


## echoServer
- Directory: echoServer
- Use GoLand to write, build, and debug your Go program.
- Very small in scope — just two files, `echoServer.go` and `main.go`.
- What you'll learn: the most basic skeleton for opening a TCP server with the `gohipernetFake` library and handling client-connect/data-received events via callbacks.
- Running it: see [the table above](#method-1-run-directly-from-a-terminal-simplest). Test it with the "echo 보내기" ("Send echo") button in `csharp_test_client`.


## chatServer
- Directory: chatServer
- A chat server with the concept of rooms.
- Packet requests are all processed on a single goroutine (thread).
- About 3-4x the size of echoServer.
- What you'll learn: the typical shape of a "stateful" server — login, then room entry, then chat. Session management, room management, and dispatching by packet ID. See item 4 under [Learning order](#learning-order) for a suggested reading order.
- Running it: see [the table above](#method-1-run-directly-from-a-terminal-simplest). Test it with `csharp_test_client`: log in → enter a room → chat.

### Additional features implemented
- 1:1 whispers
- Room invitations



## baccaratServer
- Directory: baccaratServer
- An online version of the gambling card game baccarat.
    - Baccarat rules: https://namu.wiki/w/%EB%B0%94%EC%B9%B4%EB%9D%BC
- This is chatServer with baccarat game logic layered on top, so you need to understand chatServer first.
- What you'll learn: how to keep the same room "skeleton" and layer game-specific rules on top of it — card shuffling, betting, deciding the winner, and state transitions. `roomPkg/Game.go` owns the game rules themselves, and `roomPkg/room_PacketGame.go` wires up "when this betting-request packet arrives, call this function in Game.go."
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
- What you'll learn: how to split packet processing across multiple goroutines as your player count grows, and how to keep those goroutines from stepping on each other while touching the same shared data (sessions, room state) at the same time — synchronization. Put it side by side with chatServer's files and it's much easier to see exactly what changed to make it parallel.

> **This server currently does not build.** Several symbols — `NTELIB_LOG_INFO`, `NTELIB_LOG_ERROR`, `NetLib_IsRunningServer`, `NetLib_GetCurrnetUnixTime`, and more — are used throughout the code but aren't defined anywhere in this repository. So you can't run it directly using the [instructions above](#how-to-build-and-run-the-go-servers). It's still usable for reading through the code and learning "how was this split up across multiple goroutines," but to actually run it you'll first need to fill in those missing symbols yourself (e.g. write your own logging function and a function that reports whether the server is still running, and plug them in). Treat that as a small exercise of its own if you're up for it.

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
- What you'll learn: where `chatServer` assembles and disassembles packet bodies into bytes by hand, `chatServer_msgpack` hands that job off to msgpack, a general-purpose serialization library. Comparing `protocol/packet.go` between the two side by side gives you a feel for the trade-offs between a hand-rolled protocol and using an off-the-shelf serialization library.
- Running it: see [the table above](#method-1-run-directly-from-a-terminal-simplest). Test it with `csharp_test_client_msgpack`.


## References
- [YouTube: Learning Golang TCP Socket Server Programming from Open-Source Code](https://youtu.be/boDo8JoyHuo )
