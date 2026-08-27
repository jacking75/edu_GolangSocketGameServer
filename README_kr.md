# golang_socketGameServer_codelab
>>> 리팩토링과 문서 정리 필요

[English version](README.md)

- golang을 이용하여 실시간 통신 게임 서버 만들기 실습.  
- 각 서버의 원본 코드를 하나씩 따라서 코딩하면서 서버 만드는 방법을 배운다.
    - 코딩하면서 해당 코드의 구현 방법과 이유를 설명 듣는다.
  
**버그가 있을 수 있습니다**. 버그 잡아서 수정하는 것도 학습 중 일부라고 생각해 주세요^^;.  
  
  
## 목적 
- golang으로 소켓 통신용 서버를 만들 수 있는 기술고 경험을 쌓는 것이 목표이다.
- golang으로 소켓 통신용 서버를 만든 경험이 없는(있더라도 작은) 사람을 대상으로 한다.  
- golang의 socket API를 사용하지 않고, 센트럴서버팀에서 만든 goHiperNet(golang 네트워크 라이브러리)의 짝퉁(?)을 사용한다.
    - 이 라이브러리는 goHiperNet과 API만 같고, 내부 구현은 완전 다르다. 
	- 학습용으로 충분히 사용할 수 있다.
	- golang의 socket API를 사용하여 밑바닥부터 개발하는 방법을 배우고 싶다면 별도 요청을 바람.
- 실습은 단계 별로 진행하고, 각 단계 별로 소요 시간은 다르다.
    - 시간이 많이 필요한 경우라도 1번에 최대 3시간을 넘지 않는다.
	- 한번에 너무 많이 나가면 뒤에 복습이 어려워지기 때문이다.
   
   
## 준비
- 1인 1 노트북(Windows or OSX)
- 최신 버전의 golang SDK
- 최신 버전의 GoLand  
- 기본 golang의 문법 학습
    - 코딩은 한번도 해본적인 없는 경우도 괜찮음.  
- (선택) 클라이언트로 접속 테스트를 해보고 싶다면: Windows + Visual Studio 2022 이상(.NET 8.0 SDK 포함)
    - 클라이언트 없이 서버 코드만 읽고 빌드해 보는 것만으로도 실습은 진행할 수 있다.
     
  	 
## 이 저장소에는 무엇이 있나

이 저장소는 크게 "① golang으로 만든 서버들", "② 서버가 공통으로 쓰는 가짜 네트워크 라이브러리", "③ 서버에 접속해서 테스트해 볼 수 있는 C# 클라이언트"로 이루어져 있다.

| 디렉토리 | 무엇인가 | 언어 |
|---|---|---|
| `gohipernetFake` | 모든 서버가 공통으로 사용하는 네트워크 라이브러리(짝퉁 goHiperNet). 소켓을 직접 다루는 코드는 전부 여기 있다. | Go |
| `echoServer` | 가장 단순한 서버. 받은 데이터를 그대로 돌려준다(에코). 실습의 첫 단계. | Go |
| `chatServer` | 방(room) 개념이 있는 채팅 서버. 패킷 처리를 고루틴 1개에서만 한다. | Go |
| `chatServer2` | `chatServer`를 여러 고루틴으로 병렬 처리하도록 발전시킨 버전. **현재 저장소 상태로는 빌드가 되지 않는다**(아래 [chatServer2](#chatserver2) 절 참고). | Go |
| `chatServer_msgpack` | `chatServer`와 거의 같지만, 패킷 데이터를 JSON 비슷한 바이너리 포맷인 msgpack으로 주고받는다. | Go |
| `baccaratServer` | `chatServer` 위에 바카라 카드 게임 로직을 얹은 서버. | Go |
| `csharp_test_client` | `chatServer` / `chatServer2` / `baccaratServer`(방 입장, 채팅, 릴레이 기능만) 접속 테스트용 Windows GUI 클라이언트. | C# (.NET 8, WinForms) |
| `csharp_test_client_msgpack` | `chatServer_msgpack` 전용 접속 테스트 클라이언트. | C# (.NET 8, WinForms) |
| `bin` | 빌드된 서버 `.exe`를 실행할 때 쓰는 예시 배치 파일(`run_xxx.bat`)과 로그 설정 파일이 들어있다. | - |
| `thirdparty/SimpleMsgPack.Net` | C# 클라이언트가 msgpack을 다루기 위해 가져다 쓰는 오픈소스 라이브러리(일부 코드 수정됨). | C# |
| `lib_opensource` | 참고용으로 담아 둔 다른 오픈소스 Go 네트워크 라이브러리 압축 파일들. 실습에는 사용하지 않는다. | - |

각 Go 서버 디렉토리(`echoServer`, `chatServer`, `chatServer2`, `chatServer_msgpack`, `baccaratServer`)는 **자기 자신만으로 완전히 독립된 Go 모듈**이다(각자 `go.mod`가 있다). 그래서 어느 폴더를 열어도 그 폴더 하나만 보면서 코드를 따라갈 수 있고, 서버 사이에 겹치는 코드(세션 관리, 방 관리 등)는 실습 목적상 폴더마다 그대로 복사되어 있다(공통 라이브러리로 뽑아내지 않았다). `gohipernetFake`만 예외로, 나머지 5개 서버가 전부 이 라이브러리를 `go.mod`의 `replace` 지시자로 상대 경로(`../gohipernetFake`)를 통해 가져다 쓴다.


## 학습 순서

권장 순서는 다음과 같다. 각 단계는 이전 단계의 코드/개념을 이해하고 있다는 전제로 만들어져 있으니, 순서를 건너뛰지 않는 것이 좋다.

1. **[준비](#준비)** 를 마치고 **[패킷 헤더](#패킷-헤더)** 절을 한 번 읽어 이 저장소의 서버들이 어떤 형식으로 패킷을 주고받는지 감을 잡는다.
2. **[설명 영상](#설명-영상)** 의 "실습 목적과 방법 설명" 영상을 본다.
3. **`echoServer`** — 가장 간단한 서버부터 시작한다. "Echo Server 만들기" 영상을 함께 본다.
4. **`chatServer`** — 방/로그인 개념이 있는 첫 번째 "진짜" 게임 서버 구조. "채팅 서버 코드 설명" 영상을 함께 본다. **이 서버를 제대로 이해하는 것이 이후 모든 단계의 전제조건이다.**
5. 여기서부터는 목적에 따라 두 갈래로 나뉜다(순서는 상관없다).
    - **게임 로직에 관심이 있다면 → `baccaratServer`** (chatServer + 카드 게임 로직)
    - **멀티 고루틴 동시성 처리에 관심이 있다면 → `chatServer2`** (chatServer를 여러 고루틴으로 병렬화). 다만 현재 코드는 빌드가 안 되는 상태이니 [해당 절](#chatserver2)을 먼저 읽어보길 권한다.
    - **직렬화 포맷에 관심이 있다면 → `chatServer_msgpack`** (chatServer + msgpack 직렬화)
  
   
## 패킷 헤더
패킷 허더의 크기는 총 5바이트   
- 패킷의 총 크기(2바이트. 헤더와 보디 합친) + 패킷ID(2바이트) + 패킷Type(1바이트)  
  
    
## 설명 영상
- [실습 목적과 방법 설명](https://youtu.be/zR_zcY7SXio )
- [Echo Server 만들기](https://youtu.be/OSiwcsPAO2o )
- [채팅 서버 코드 설명](https://youtu.be/2rppKuW-wQg )
  
   
## Go 서버를 빌드하고 실행하는 방법

`echoServer`, `chatServer`, `chatServer_msgpack`, `baccaratServer` 4개는 모두 같은 방식으로 빌드/실행한다(`chatServer2`는 현재 빌드가 안 되므로 제외 — [아래 절](#chatserver2) 참고). 두 가지 방법이 있다.

### 방법 1) 터미널에서 바로 실행 (가장 간단함)

서버 디렉토리로 들어가서 `go run .` 뒤에 접속 주소/최대 접속자 수 등을 옵션(플래그)으로 넘겨서 실행한다. 예를 들어 `echoServer`는 다음과 같다.

```bash
cd echoServer
go run . -c_IsTcp4Addr=true -c_BindAddress=127.0.0.1:11021 -c_MaxSessionCount=200 -c_MaxPacketSize=1024
```

실행하면 아래처럼 로그가 출력되고, 그 상태로 접속을 기다린다(멈춘 게 아니라 정상이다. 종료하려면 `Ctrl+C`).

```
[ info ] tcpServerStart - Start
2026/08/27 23:53:28 Server Listen ...
```

> **주의**: `-c_MaxSessionCount`를 빼먹고 실행하면 기본값이 0이 되어, 서버는 켜지지만 **아무도 접속할 수 없다**(내부적으로 접속 가능한 세션 슬롯이 0개로 초기화되기 때문). 반드시 위 예시처럼 플래그를 함께 넘기자.

서버마다 넘겨야 하는 플래그가 조금씩 다른데, `bin` 디렉토리의 배치 파일에 서버별로 검증된 옵션이 이미 적혀 있으니 그대로 가져다 쓰면 된다.

| 서버 | 실행 명령 (해당 디렉토리 안에서 실행) |
|---|---|
| echoServer | `go run . -c_IsTcp4Addr=true -c_BindAddress=127.0.0.1:11021 -c_MaxSessionCount=200 -c_MaxPacketSize=1024` |
| chatServer | `go run . -c_IsTcp4Addr=true -c_BindAddress=127.0.0.1:11021 -c_MaxSessionCount=200 -c_MaxPacketSize=1024` |
| chatServer_msgpack | `go run . -c_IsTcp4Addr=true -c_BindAddress=127.0.0.1:11021 -c_MaxSessionCount=200 -c_MaxPacketSize=1024` |
| baccaratServer | `go run . -c_IsTcp4Addr=true -c_BindAddress=127.0.0.1:11021 -c_MaxSessionCount=200 -c_MaxPacketSize=1024` |

전부 기본적으로 `127.0.0.1:11021`(내 컴퓨터의 11021번 포트)에서 접속을 기다린다. 서버를 동시에 두 개 이상 띄워보고 싶다면 `-c_BindAddress`의 포트 번호를 서로 다르게 주면 된다.

`go run .` 대신 `go build`로 실행 파일을 만들어 둘 수도 있다. 실행 파일 이름은 디렉토리 이름을 따라 자동으로 정해진다(`echoServer` 디렉토리에서 빌드하면 `echoServer.exe`가 만들어짐). `bin` 디렉토리의 `run_xxx.bat` 파일은 바로 이렇게 만든 `.exe`를 실행하는 용도이므로, 빌드한 `.exe`를 `bin` 디렉토리로 옮기면 그 배치 파일을 그대로 쓸 수 있다.

### 방법 2) GoLand에서 실행

GoLand로 서버 디렉토리(예: `echoServer`)를 프로젝트로 연 뒤, `main.go`를 실행하는 Run Configuration을 만들고 **Program arguments**에 위 표의 플래그들(`-c_IsTcp4Addr=true -c_BindAddress=...` 이하)을 그대로 입력하면 터미널로 실행한 것과 동일하게 동작한다. 디버깅(중단점 걸기 등)이 필요하면 이 방법을 쓰자.


## 테스트 클라이언트 사용법

서버가 잘 만들어졌는지 눈으로 확인해 보려면 이 저장소에 포함된 C# GUI 클라이언트를 쓰면 된다. 텍스트 기반 도구(telnet 등)로는 테스트할 수 없다 — 패킷이 사람이 읽을 수 있는 텍스트가 아니라 [패킷 헤더](#패킷-헤더)에서 설명한 바이너리 포맷이기 때문이다.

1. `csharp_test_client/csharp_test_client.sln`을 Visual Studio로 연다. (`chatServer_msgpack`을 테스트할 때는 대신 `csharp_test_client_msgpack/csharp_test_client_msgpack.sln`을 연다.)
2. 빌드 후 실행한다(F5). WinForms 창이 뜬다.
3. IP 입력란은 기본값이 `0.0.0.0`이므로, "localhost" 체크박스를 켜거나 직접 `127.0.0.1`을 입력한다. Port는 기본값 `11021`을 서버와 동일하게 맞춘다(서버를 다른 포트로 띄웠다면 그 포트로 맞춘다).
4. **접속하기** 버튼을 눌러 서버에 연결한다.
5. 서버별로 이어서 테스트할 기능이 다르다.
    - **echoServer**: "echo 보내기" 버튼을 누르면 보낸 내용이 그대로 되돌아오는 것을 로그에서 확인할 수 있다.
    - **chatServer / baccaratServer(그리고 나중에 chatServer2를 직접 고쳐서 띄운다면 그것도)**: UserID/PW 입력란을 채우고 **Login** 버튼 → 방 번호를 입력하고 **방 입장** 버튼 → 채팅 입력란을 채우고 **chat** 버튼 순서로 눌러본다. 클라이언트 창을 두 개 띄워서 같은 방에 각각 입장시키면, 한쪽에서 보낸 채팅이 다른 쪽에도 전달되는 것을 볼 수 있다.
    - **baccaratServer**: 위 채팅/방 기능까지는 이 클라이언트로 테스트할 수 있지만, 배팅 등 바카라 게임 자체를 조작하는 UI는 클라이언트에 없다(직접 추가해야 한다). 게임 진행 흐름은 `baccaratServer/docs/*.puml`(PlantUML 시퀀스 다이어그램)을 참고하자.
    - **chatServer_msgpack**: 같은 순서로 `csharp_test_client_msgpack`을 사용한다.


## echoServer
- 디렉토리: echoServer
- GoLand를 사용하여 golang용 프로그램을 만들고, 빌드/디버깅을 한다.
- 아주 간단한 규모이다.
- 실행: [위 표](#방법-1-터미널에서-바로-실행-가장-간단함) 참고. 테스트는 `csharp_test_client`의 "echo 보내기" 버튼으로 한다.
  
  
## chatServer
- 디렉토리: chatServer
- 방 개념의 채팅 서버
- 패킷 요청 처리를 1개의 고루틴(스레드)에서만 한다.
- echoServer에 비해 규모는 3~4배 크다.
- 실행: [위 표](#방법-1-터미널에서-바로-실행-가장-간단함) 참고. 테스트는 `csharp_test_client`로 로그인 → 방 입장 → 채팅 순서로 한다.
  
### 추가 기능 구현
- 1:1 귓속말
- 방 초대
    
   
   
## baccaratServer 
- 디렉토리: baccaratServer
- 겜블 게임인 바카라 게임을 온라인화 한 것이다.
    - 바카라 룰: https://namu.wiki/w/%EB%B0%94%EC%B9%B4%EB%9D%BC
- chatServer에 바카라 게임 로직이 올라간 것으로 chatServer에 대한 이해가 꼭 필요하다.
- 실행: [위 표](#방법-1-터미널에서-바로-실행-가장-간단함) 참고. 게임 진행 흐름은 `baccaratServer/docs/*.puml`을 참고. 클라이언트에 배팅 UI는 없으므로 방 입장/채팅까지만 `csharp_test_client`로 눈으로 확인할 수 있다.
  
### 추가 기능 구현
- 게임 서버 Scale-Out 기능 구현
-  API Server(http)와 연동  
    - 유저를 특정 게임 서버에 할당하는 기능
    - 매칭 기능	

> 위 "추가 기능 구현" 항목은 실습 목표이며, 현재 커밋된 코드에는 아직 구현되어 있지 않다.
	 
	 
	 
## chatServer2
- 디렉토리: chatServer2
- 방 개념의 채팅 서버
- 패킷 요청 처리를 N개의 고루틴(스레드)에서 한다.
    - 패킷 처리를 멀티 고루틴에서 하므로 공유 객체 동기화를 조심해야 한다.
- chatServer의 코드와 겹치는 부분이 많으므로 chatServer에 대한 이해가 꼭 필요하다

> **현재 이 서버는 빌드되지 않는다.** `NTELIB_LOG_INFO`, `NTELIB_LOG_ERROR`, `NetLib_IsRunningServer`, `NetLib_GetCurrnetUnixTime` 등 여러 심볼이 코드 곳곳에서 쓰이지만 이 저장소 어디에도 정의되어 있지 않다. 그래서 위 [실행 방법](#go-서버를-빌드하고-실행하는-방법)으로 바로 실행할 수 없다. 코드를 읽으면서 "멀티 고루틴으로 어떻게 구조를 나눴는가"를 학습하는 용도로는 그대로 쓸 수 있지만, 직접 실행해 보려면 먼저 이 심볼들을 채워 넣어야 한다.
    
### 추가 기능 구현
- Redis 연동
- API Server(http)와 연동  
    - 로그인을 API Server에서 한다.  

> 위 "추가 기능 구현" 항목은 실습 목표이며, 현재 커밋된 코드에는 아직 구현되어 있지 않다.  
     
	 
## msgpack을 사용한 chatServer	
- 디렉토리: chatServer_msgpack
- 클라이언트 디렉토리: csharp_test_client_msgpack	
- 서버와 클라이언트가 네트워크로 주고 받는 패킷 데이터 포맷을 msgpack을 사용한다.
    - [Go](https://github.com/vmihailenco/msgpack )
	- [C#](https://github.com/ymofen/SimpleMsgPack.Net  ) 
	    - golang 라이브러리와 데이터 포맷이 일치하지 않는 부분이 있어서 코드를 수정하였음.
		- thirdparty/SimpleMsgPack.Net  디렉토리에 코드가 있다/
- 실행: [위 표](#방법-1-터미널에서-바로-실행-가장-간단함) 참고. 테스트는 `csharp_test_client_msgpack`으로 한다.
      
   
## 참고
- [유튜브: 오픈소스 코드로 배우는 Golang TCP Socket Server 프로그래밍 ](https://youtu.be/boDo8JoyHuo )

  
