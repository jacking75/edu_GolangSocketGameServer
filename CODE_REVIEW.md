# 코드 리뷰 및 개선 제안

> 이 문서는 `golang_socketGameServer_codelab` 저장소 전체(핵심 네트워크 라이브러리 `gohipernetFake` + `echoServer`, `chatServer`, `chatServer2`, `chatServer_msgpack`, `baccaratServer` 5개 서버)를 정밀 리뷰한 결과와 구체적인 개선 방법을 정리한 것이다. 성능 테스트는 범위에서 제외했다.
>
> 리뷰 방식: 각 디렉토리의 `.go` 파일을 전부 읽고 실제 실행 경로를 추적했으며, 아래 표시된 파일:라인 번호는 재확인을 거쳤다.

## 적용 현황

[6절 로드맵](#6-우선순위별-개선-로드맵)의 1~5단계 모두 이 저장소에 실제로 반영되었다. 예외는 다음과 같다.

- **1.4(공통 모듈 추출)는 의도적으로 건너뛰었다.** 저장소 소유자와 상의한 결과, 각 서버 폴더를 "한 폴더씩 따라 치며 배우는" 실습 구조로 유지하기로 했다. 코드 중복 자체는 여전히 남아 있다.
- **chatServer2는 이번 작업 범위 밖의 사전 존재 문제로 여전히 빌드되지 않는다.** `NTELIB_LOG_INFO`, `NTELIB_LOG_ERROR`, `NTELIB_LOG_DEBUG`, `NetLib_IsRunningServer`, `NetLib_GetCurrnetUnixTime`, `NetLibInitNetwork` 심볼이 15개 넘는 파일에서 참조되지만 저장소 어디에도 정의되어 있지 않다(gohipernetFake에도 없음). 이 저장소를 정적으로 읽는 리뷰만으로는 드러나지 않고 실제 `go build`를 해봐야 드러나는 문제라, 이번 코드 리뷰 자체에는 기록되지 않았던 별도의 이슈다. chatServer2에 대한 3단계 수정(4.1~4.4)은 코드 자체는 모두 반영했고 `gofmt`로 구문 유효성은 확인했지만, 위 문제 때문에 `go build`/`go vet`/`go test`로 실제 검증하지는 못했다. 이 심볼들을 무엇으로 채워야 할지(로깅 방식, 실행 상태 플래그 등)는 설계 판단이 필요해 별도로 다뤄야 한다.
- 3~5단계를 적용하는 과정에서 확인된 동일 패턴의 버그(예: `getRoomByNumber`의 `roomIndex < 0` 경계 검사 누락)는 로드맵에 명시적으로 배정되어 있지 않았더라도 해당 파일을 수정하는 김에 함께 고쳤다.
- 5단계의 "죽은 코드 정리"는 사용처가 전혀 없음을 grep으로 확인한 항목만 제거했다. `chatServer/logger.go`의 `init_Log`처럼 미완성이지만 실제로 쓸모 있어 보이는 코드(zap 기반 파일 로깅 전체 구성)는 임의로 지우거나 활성화하지 않고 그대로 두었다.

## 총평

이 저장소는 "실습 중 버그를 찾아 고치는 것도 학습"이라는 README의 취지에 맞게, 단순 스타일 문제가 아니라 **네트워크 계층에서 실제로 서버 전체를 죽일 수 있는 버그**, **인증/권한 우회**, **게임 로직 자체의 오류**가 다수 존재한다. 특히 가장 심각한 문제들은 모든 서버가 공유하는 `gohipernetFake` 라이브러리에 몰려 있어, 이 한 곳만 고쳐도 5개 서버 전부의 안정성이 크게 개선된다.

문서는 다음 순서로 구성했다.

1. [저장소 공통/구조적 이슈](#1-저장소-공통구조적-이슈)
2. [gohipernetFake — 모든 서버에 영향 (최우선)](#2-gohipernetfake--모든-서버에-영향-최우선)
3. [echoServer / chatServer](#3-echoserver--chatserver)
4. [chatServer2 (멀티 고루틴)](#4-chatserver2-멀티-고루틴)
5. [chatServer_msgpack / baccaratServer](#5-chatserver_msgpack--baccaratserver)
6. [우선순위별 개선 로드맵](#6-우선순위별-개선-로드맵)

---

## 1. 저장소 공통/구조적 이슈

### 1.1 테스트 코드가 전혀 없음
저장소 전체(9,000줄+)에 `*_test.go` 파일이 단 하나도 없다. 특히 다음 세 영역은 순수 함수라 단위 테스트를 붙이기 쉽고, 붙였다면 이번에 발견된 버그 다수를 즉시 잡아냈을 것이다.
- `gohipernetFake/packetEnDecoder.go`의 `RawPacketData` 인코딩/디코딩 (`ReadU16`, `WriteS8` 등)
- `gohipernetFake/TcpSession.go`의 `makePacket` (다양한 길이의 스트림을 흘려보내는 테이블 테스트)
- `baccaratServer/roomPkg/Game.go`의 `_baccaratCardIndexToScore`, `doBaccarat` (52장 전체를 대상으로 한 카드값 검증)

**개선 방법**: 최소한 위 세 함수부터 `go test`로 커버하라. 특히 `_baccaratCardIndexToScore`는 `for cardIndex := int8(0); cardIndex < 52; cardIndex++`로 전수 검사하는 테이블 테스트 하나만 있어도 2절에서 설명할 게임 로직 버그를 즉시 잡아낼 수 있었다.

### 1.2 CI가 없음
`.github/workflows` 등 CI 설정이 전혀 없다. `go build ./...`, `go vet ./...`조차 자동으로 돌지 않는다.

**개선 방법**: GitHub Actions로 각 서버 디렉토리에서 `go build ./... && go vet ./...`를 PR마다 실행하는 워크플로를 추가하라. `go vet`만 돌려도 일부 미사용 변수/의심스러운 형변환을 잡아낼 수 있다.

### 1.3 `go.mod`가 Go 1.13 기준으로 고정되어 있음
6개 모듈의 `go.mod`가 모두 `go 1.13`을 명시한다. 로컬에 설치된 Go는 1.24.1이라 최신 언어 기능(제네릭 등)이나 `golangci-lint`, `staticcheck` 같은 최신 정적분석 도구의 이점을 못 받는다.

**개선 방법**: `go mod tidy`와 함께 `go 1.21` 이상으로 올리고, 가능하면 `golangci-lint run`을 CI에 추가하라.

### 1.4 5개 서버 간 코드가 대량으로 복사·붙여넣기 되어 있고, 같은 버그가 여러 곳에 반복됨
`connectedSessions/session.go`의 `getRoomNumber()`가 대표적 사례다. 아래 두 줄은 chatServer2, chatServer_msgpack, baccaratServer에 **거의 동일하게** 존재한다.

```go
func (session *session) getRoomNumber() (int32, int32) {
    roomNum := atomic.LoadInt32(&session._RoomNum)
    roomNumOfEntering := atomic.LoadInt32(&session._RoomNum) // 버그: _RoomNumOfEntering이어야 함
    return roomNum, roomNumOfEntering
}
```
두 번째 반환값이 `_RoomNum`을 또 읽어서, "입장 중인 방 번호"가 항상 "현재 방 번호"와 같은 값으로 나온다. 이 복붙 패턴이 chatServer2에서는 **방 유저 슬롯이 영구적으로 새는 실제 버그**로 이어진다(4.1절 참고).

**개선 방법**: 이런 복붙 구조 자체가 문제다. `connectedSessions`, `roomPkg`의 공통 뼈대(세션 상태, 방 관리자 골격, 패킷 헤더 처리)를 별도 Go 모듈(예: `gamenetcommon`)로 뽑아내고, 각 서버는 이를 import해서 게임 고유 로직만 얹는 구조로 리팩토링하라. 이렇게 하면 위와 같은 버그를 한 곳에서만 고치면 모든 서버에 전파된다. README에도 이미 "리팩토링과 문서 정리 필요"라고 적혀 있으니, 학습 단계가 끝난 뒤 이 리팩토링 자체를 다음 실습 주제로 삼는 것도 좋다.

### 1.5 로그 레벨/설정값이 실제로는 반영되지 않는 "죽은 설정"이 다수
`gohipernetFake/configNetwork.go`의 `NetworkConfig.MaxPacketSize`, `MaxReceiveBufferSize` 필드는 사용자가 값을 넣어도 무시되고, 실제로는 `gohipernetFake/define.go`의 상수 `MAX_PACKET_SIZE(1024)`, `MAX_RECEIVE_BUFFER_SIZE(8012)`가 하드코딩되어 쓰인다(`TcpSession.go:22`, `TcpSession.go:69`). 설정 구조체가 실제 동작을 제어한다고 오해하기 쉽다.

**개선 방법**: `TcpSession`이 세션 생성 시 `NetworkConfig` 값을 필드로 받아 그 값을 쓰도록 배선하거나, 당장 안 쓸 거라면 필드 자체를 제거해 API를 정직하게 만들어라.

---

## 2. gohipernetFake — 모든 서버에 영향 (최우선)

`gohipernetFake`는 6개 서버 전부가 의존하는 소켓 계층이다. 여기서 발견한 문제는 **모든 서버에 그대로 상속**되므로 가장 먼저 고쳐야 한다.

### 2.1 [심각] 패킷 송신 함수가 초기화되지 않아 첫 패킷 전송 시 반드시 패닉
- **파일**: `gohipernetFake/goHiperNet.go:24-29`, `gohipernetFake/goHiperNet_Impl.go:54-61`
- **내용**: `NetLibISendToClient`, `NetLibISendToAllClient`, `NetLibIPostSendToClient`, `NetLibIPostSendToAllClient`는 `goHiperNet.go`에 `func(...) bool` 타입의 **nil 함수 변수**로만 선언되어 있다. 이 변수들을 실제 함수(`sendToClient` 등)로 연결하는 코드는 `_InitNetworkSendFunction()`(`goHiperNet_Impl.go:54`) 하나뿐인데, 저장소 전체를 검색해도 이 함수를 호출하는 곳이 없다. `NetLibStartNetwork` → `start_Network_Impl`(`goHiperNet_Impl.go:10-16`)도 이걸 호출하지 않는다.
- **왜 문제인가**: 5개 서버 모두 `NetLibISendToClient(...)`(예: `sendPacket()` 계열 함수)로 클라이언트에게 응답을 보낸다. 이 함수 변수가 nil인 채로 호출되면 **"call of nil func value" 런타임 패닉**이 발생한다. 즉 이 라이브러리를 그대로 쓰는 한, 클라이언트가 서버에 접속해 응답이 필요한 첫 요청(로그인 등)을 보내는 순간 서버가 죽어야 정상이다. (실습 중 우연히 이 경로를 다르게 배선했거나, 이 버그를 아직 아무도 밟지 않았을 가능성이 있다 — 어느 쪽이든 반드시 고쳐야 한다.)
- **개선 방법**: `start_Network_Impl`(`goHiperNet_Impl.go:10`) 맨 앞, 또는 `NetLibStartNetwork`(`goHiperNet.go:15`) 진입부에서 `_InitNetworkSendFunction()`을 호출하도록 한 줄만 추가하면 된다.

### 2.2 [심각] 패킷 길이 필드의 부호 반전 + 하한 검증 누락 → 원격 크래시/무한루프 DoS
- **파일**: `gohipernetFake/packetEnDecoder.go:11-14` (`PacketTotalSize`), `gohipernetFake/TcpSession.go:49-88` (`makePacket`)
- **내용**:
  ```go
  func PacketTotalSize(data []byte) int16 {
      totalsize := binary.LittleEndian.Uint16(data)
      return int16(totalsize) // uint16 → int16 강제 캐스팅
  }
  ```
  클라이언트가 패킷 헤더의 "전체 크기" 필드에 32768~65535 사이 값을 넣으면 `int16`으로 캐스팅되며 **음수로 뒤집힌다**. `makePacket`(`TcpSession.go:63-75`)은 이 값을 `requireDataSize`로 받아 `requireDataSize > readAbleByte`, `requireDataSize > MAX_PACKET_SIZE` 두 검증만 하는데, 값이 음수면 두 조건 모두 통과해버린다. 이어서 `recviveBuff[readPos:(readPos + requireDataSize)]`가 실행되며 **음수 상한 슬라이싱으로 panic**이 발생한다.
  또한 `requireDataSize`가 `0`처럼 매우 작은 값(헤더 크기 5보다 작음)이어도 하한 검증이 전혀 없다. 이 경우 `readPos += requireDataSize`, `readAbleByte -= requireDataSize`가 사실상 진행되지 않아 **같은 위치를 무한 반복하며 `OnReceive`를 계속 호출하는 CPU 100% 무한루프**에 빠진다.
- **왜 문제인가**: 조작된 5바이트짜리 TCP 페이로드 하나로 원격에서 서버 프로세스를 크래시시키거나(2.3절과 결합 시 전체 다운), 해당 연결의 고루틴을 영구히 CPU 스핀시킬 수 있는 전형적인 DoS 벡터다.
- **개선 방법**:
  - `PacketTotalSize`가 `uint16`을 그대로 반환하도록 시그니처를 바꾸거나, 최소한 `int32`로 승격해서 반환하라.
  - `makePacket` 루프 초입에 `if requireDataSize < PacketHeaderSize { /* 잘못된 패킷: 연결 종료 */ return startRecvPos, NET_ERROR_RECV_MAKE_PACKET_TOO_LARGE_PACKET_SIZE }` 같은 하한 검증을 추가하라.

### 2.3 [심각] 연결별 수신 고루틴에 panic recover가 없어, 한 연결의 오류가 서버 프로세스 전체를 종료시킴
- **파일**: `gohipernetFake/goHiperNet_Impl.go:47` (`go client.handleTcpRead(...)`), `gohipernetFake/TcpSession.go:16-47` (`handleTcpRead`)
- **내용**: `_start_TCPServer_block`과 `start_Network_Impl`에는 `defer PrintPanicStack()`이 있지만(각각 accept 루프와 시작 함수 자신을 보호), 정작 각 연결을 처리하는 `handleTcpRead` 고루틴 내부에는 recover 로직이 없다. Go에서 recover되지 않은 고루틴 패닉은 **프로세스 전체를 즉시 종료**시킨다.
- **왜 문제인가**: 2.2절 같은 악성 패킷 하나, 혹은 앞으로 추가될 패킷 처리 로직의 사소한 버그 하나가 **한 명의 클라이언트 때문에 접속 중인 모든 유저의 서버가 함께 죽는** 결과로 이어진다.
- **개선 방법**: `handleTcpRead` 맨 앞에 `defer session.closeProcess()`와 `defer PrintPanicStack()`을 추가해 해당 세션만 안전하게 종료되도록 하라(패닉이 나면 그 연결만 끊고 서버는 계속 동작).

### 2.4 [높음] `Accept()` 에러 무시 → 리스너 종료 후 CPU 100% 무한루프
- **파일**: `gohipernetFake/goHiperNet_Impl.go:38`
  ```go
  conn, _ := _mClientListener.Accept()
  ```
- **왜 문제인가**: `NetLibStopListen()`(`stopListen_impl`, `goHiperNet.go:19-21`) 호출 후에는 `Accept()`가 즉시 에러를 반환하는데, 이 에러를 무시하고 `conn == nil`인 채로 계속 `&TcpSession{TcpConn: nil, ...}`을 만들어 `addSession`, `go client.handleTcpRead(...)`까지 실행한다. `handleTcpRead`가 nil `TcpConn.Read()`를 호출하면 nil pointer dereference 패닉(2.3절 수정 전이라면 서버 전체 다운)이 나고, 그 전까지는 `for { conn, _ := Accept() ... }`가 즉시 에러를 반환하며 무한히 도는 CPU 스핀 루프가 된다.
- **개선 방법**: `if err != nil { break }` (또는 로그 후 `continue` — 일시적 에러와 리스너 종료를 구분하려면 `errors.Is`로 판별)를 추가하라.

### 2.5 [높음] `addSession()`의 실패(bool) 반환값을 무시 → 세션 풀 상한이 무력화되고, 세션 인덱스 충돌 발생
- **파일**: `gohipernetFake/goHiperNet_Impl.go:45`
  ```go
  _tcpSessionManager.addSession(client)
  go client.handleTcpRead(networkFunctor)
  ```
- **왜 문제인가**: `addSession`(`clientSessionManager.go:50-66`)은 세션 인덱스 풀이 소진되면(`MaxSessionCount` 초과) `false`를 반환하고 `client.Index`를 세팅하지 않는다(0으로 남음). 그런데 반환값을 확인하지 않아 **`MaxSessionCount` 제한이 사실상 무의미**해지고, `client.Index == 0`인 이 세션이 나중에 끊기면 `closeProcess()` → `removeSession(0, seq)`가 **실제로 인덱스 0을 쓰고 있는 다른 정상 세션의 슬롯을 훔쳐 반환**해버린다. 이후 서로 다른 두 연결이 같은 인덱스를 공유하는 상태로 이어진다.
- **개선 방법**: `if !ok := _tcpSessionManager.addSession(client); !ok { conn.Close(); continue }` 형태로 실패 시 연결을 닫고 고루틴을 만들지 않도록 고쳐라.

### 2.6 [높음] `closeProcess()`가 멱등적이지 않아 이중 호출 시 세션 인덱스가 이중 반환(double free)됨
- **파일**: `gohipernetFake/TcpSession.go:90-95`
- **내용**: 애플리케이션이 `NetLibForceDisconnectClient` → `forceDisconnectClient`(`clientSessionManager.go:115-125`)를 호출해 `closeProcess()`를 실행하는 경로와, 그 결과로 소켓이 닫히면서 블로킹 중이던 `handleTcpRead`의 `Read()`가 에러를 반환해 **자기 자신도 `closeProcess()`를 호출**하는 경로가 동시에 존재할 수 있다. `closeProcess()`에는 "이미 처리됨" 여부를 막는 가드가 없다.
- **왜 문제인가**: 두 경로가 겹치면 `removeSession`(`clientSessionManager.go:68-71`) → `_freeSessionIndex`가 같은 인덱스를 세션 풀에 **두 번 반환**한다. 이후 서로 다른 두 신규 연결이 같은 세션 인덱스를 할당받아 데이터가 뒤섞이고, `OnClose` 콜백도 애플리케이션 쪽에서 두 번 호출되어 상위 로직(방 퇴장 처리 등)이 이중으로 실행된다.
- **개선 방법**: `TcpSession`에 `closed int32`(atomic) 또는 `sync.Once` 필드를 추가하고, `closeProcess()` 시작부에서 `if !atomic.CompareAndSwapInt32(&session.closed, 0, 1) { return }`로 1회 실행을 보장하라.

### 2.7 [중간] 헤더 수신 판단이 "이번 Read 호출에서 받은 바이트 수"만 검사함
- **파일**: `gohipernetFake/TcpSession.go:32-36`
  ```go
  if recvBytes < PACKET_HEADER_SIZE {
      session.closeProcess()
      return
  }
  ```
- **왜 문제인가**: `recvBytes`는 이번 `Read()` 호출에서 받은 바이트 수일 뿐, 누적 버퍼(`startRecvPos`)를 포함하지 않는다. 정상적인 클라이언트라도 네트워크 상황에 따라 헤더 5바이트가 여러 번의 TCP 세그먼트로 나뉘어 도착하면(예: 3바이트 먼저 도착 후 나머지 2바이트), 두 번째 `Read()`가 2바이트만 반환하는 순간 `recvBytes(2) < PACKET_HEADER_SIZE(5)`가 참이 되어 **정상 연결이 스퓨리어스하게 끊긴다**.
- **개선 방법**: 연결 종료 판단은 `recvBytes == 0`(EOF) 또는 `err != nil`로만 하고, "아직 헤더/바디가 덜 왔다"는 판단은 이미 `makePacket`이 누적 `readAbleByte` 기준으로 올바르게 하고 있으므로(`TcpSession.go:59`) 이 라인의 조기 종료 검사는 제거하라.

### 2.8 [중간] `RawPacketData.ReadBytes`만 범위 검사가 빠져 있음
- **파일**: `gohipernetFake/packetEnDecoder.go:166-170`
  ```go
  func (p *RawPacketData) ReadBytes(readSize int) (refSlice []byte) {
      refSlice = p.data[p.pos : p.pos+readSize]
      p.pos += readSize
      return
  }
  ```
  같은 파일의 `ReadU16`/`ReadU32`/`ReadU64`/`ReadByte`/`ReadString`은 전부 `p.pos+n > len(p.data)`를 검사하지만 `ReadBytes`만 없다.
- **왜 문제인가**: 손상되었거나 악의적으로 조작된 페이로드 안의 길이 필드를 이 함수에 그대로 넘기면 `slice bounds out of range` 패닉으로 이어진다(2.3절 수정 전이면 서버 전체 다운).
- **개선 방법**: 다른 Read 함수들과 동일하게 `if p.pos+readSize > len(p.data) { return nil }` (또는 error 반환하도록 시그니처 통일) 가드를 추가하라.

### 2.9 [중간] `sendPacket`의 쓰기 에러를 무시
- **파일**: `gohipernetFake/clientSessionManager.go:81`
  ```go
  session.sendPacket(sendData)
  return true
  ```
  `TcpSession.sendPacket`(`TcpSession.go:98-101`)은 `error`를 반환하지만 호출부에서 버린다.
- **왜 문제인가**: 상대가 이미 연결을 끊었는데 서버가 계속 write를 시도하는 경우(broken pipe 등) 이를 감지하지 못해 죽은 세션이 정리되지 않고 방치된다.
- **개선 방법**: 에러 발생 시 `session.closeProcess()`를 트리거하도록 고쳐라.

### 2.10 [낮음] 그 외 정리하면 좋은 항목
- `gohipernetFake/goHiperNet_Impl.go:27` — `net.ResolveTCPAddr`의 에러를 무시(`tcpAddr, _ := ...`). `BindAddress` 설정 오타가 조용히 묻힌다. `err != nil`이면 `log.Fatal`로 즉시 실패시켜라.
- `gohipernetFake/clientSessionManager.go:92-103` — `_curSessionCount`와 `_IncConnectedessionCount`/`_DecConnectedessionCount`가 정의만 되고 `addSession`/`removeSession`에서 호출되지 않아 `_connectedessionCount()`는 항상 0을 반환하는 죽은 코드다. 실제 연결 수 카운트가 필요하면 성공 경로에서 호출을 배선하라.
- `gohipernetFake/utilDeque.go` — `Deque` 구조체가 `sync.RWMutex`를 이름 없이 임베딩해 `Lock()`/`Unlock()`이 외부에 그대로 노출된다. `mu sync.RWMutex` 같은 이름 있는 필드로 캡슐화해 실수로 외부에서 잠그고 안 푸는 사고를 막아라.
- `gohipernetFake/clientSessionManager.go:105-113` — `_findSession(sessionIndex, sessionUniqueId)`가 `sessionIndex` 파라미터를 전혀 쓰지 않고 `sessionUniqueId`만으로 조회한다. API가 두 값을 다 요구하는 것처럼 보이지만 실제로는 하나만 검증하므로, 세션이 재사용된 상황에서 "옛 sessionIndex + 새 sessionUniqueId" 같은 불일치를 걸러내지 못한다. 검증에 실제로 사용하거나 파라미터를 제거하라.

---

## 3. echoServer / chatServer

echoServer는 구조가 단순해 특별한 결함이 없다. 문제는 chatServer(및 이를 그대로 물려받은 chatServer2/chatServer_msgpack/baccaratServer)에 집중된다.

### 3.1 [높음] 패킷 헤더 길이 미검증 → 조작된 짧은 패킷으로 서버 전체 크래시
- **파일**: `chatServer/protocol/packet.go:62-78` (`PeekPacketID`, `PeekPacketBody`), 호출부 `chatServer/distributePacket.go:17-18`
- **내용**: `PeekPacketID`는 `rawData[2:]`를, `PeekPacketBody`는 `rawData[0:2]`를 **길이 검증 없이** 슬라이싱한다. 2.2절에서 설명했듯 `gohipernetFake`의 `makePacket`은 조작된 헤더에 대해 검증이 허술해 5바이트 미만짜리 데이터도 `OnReceive`로 넘길 수 있다.
- **왜 문제인가**: `handleTcpRead`(연결별 고루틴, 2.3절)에 recover가 없는 상태에서 이 두 함수가 패닉하면 **서버 프로세스 전체가 죽는다**. 2.2/2.3/2.8이 근본 원인이지만, 애플리케이션 계층에서도 자체 방어선을 하나 더 두는 것이 안전하다.
- **개선 방법**: `PeekPacketID`/`PeekPacketBody` 진입부에 `if len(rawData) < protocol.PacketHeaderSize { return 0, nil, false }`류의 방어 코드를 추가하고, 실패 시 `distributePacket`이 해당 세션을 끊도록 하라. 네트워크 계층 수정(2.2, 2.3)과 애플리케이션 계층 방어(이 항목)는 **둘 다** 해야 한다 — 계층 하나에만 의존하면 나중에 다른 계층이 바뀔 때 방어가 사라진다.

### 3.2 [높음] 로그인 없이 방에 입장 가능 (인증 우회)
- **파일**: `chatServer/connectedSessions/sessionManager.go:88-94` (`GetUserID`), `chatServer/roomPkg/room_PacketEnter.go:20-24`
- **내용**: `GetUserID`는 세션 인덱스가 유효한지(`_validSessionIndex`)만 확인하고 로그인 여부(`IsAuth()`)는 확인하지 않는다. `room_PacketEnter.go`는 `ok` 값만 보고 통과시킨다.
- **왜 문제인가**: 로그인하지 않은 클라이언트가 곧바로 `PACKET_ID_ROOM_ENTER_REQ`를 보내면 빈 `userID`(길이 0)로 `ok=true`가 되어 그대로 방에 등록된다. 로그인 절차를 완전히 우회한 "유령 유저"가 채팅 릴레이 등에 참여할 수 있다.
- **개선 방법**: `GetUserID` 내부 또는 `room_PacketEnter.go`의 호출부에서 `session.IsAuth()`를 함께 검사하고, 인증되지 않은 세션의 방 관련 요청은 즉시 에러 응답 후 무시(또는 연결 종료)하라.

### 3.3 [높음] 방 입장 실패 시 롤백 누락 → 유령 유저가 영구히 쌓임(자원 고갈)
- **파일**: `chatServer/roomPkg/room_PacketEnter.go:31-41`
- **내용**: `room.addUser(userInfo)` 성공 후 `connectedSessions.SetRoomNumber(...)`가 실패(이미 다른 방에 있는 유저 등)하면 에러 응답만 보내고 **`room._removeUser`로 되돌리지 않는다**.
- **왜 문제인가**: 이미 방 A에 있는 유저가 방 B에 입장을 시도하면, 방 B의 내부 맵(`_userSessionUniqueIdMap`)과 `_curUserCount`에는 등록되지만 세션 쪽 기록은 여전히 방 A를 가리킨다. 이 유저가 끊어져도 연결 종료 처리(`ProcessPacketSessionClosed`, `chatServer/distributePacket.go:106-122`)는 세션이 "알고 있는" 방 A에만 퇴장 통지를 보내므로, **방 B의 유령 엔트리는 서버 재시작 전까지 절대 정리되지 않는다**. 이 시나리오를 반복하면 방 정원이 실제 인원 없이도 소진돼 신규 유저가 못 들어오는 서비스 거부 상태가 된다.
- **개선 방법**: `SetRoomNumber` 실패 분기에 `room._removeUser(newUser)` 호출을 추가해 방 등록을 원자적으로 롤백하라. (이 패턴은 4.1절에서 chatServer2에도 동일하게 나타난다 — 공통 라이브러리로 뽑을 때 반드시 함께 고쳐야 한다.)

### 3.4 [중간] 세션 인덱스 재사용 경쟁 상태
- **파일**: `chatServer/chatServer.go:101-119` (`disConnectClient`), `gohipernetFake/TcpSession.go:90-95` (`closeProcess`)
- **내용**: `disConnectClient`는 로그인 유저 정리를 `PacketChan`(버퍼 256, 비동기)에 위임하지만, `gohipernetFake`의 `closeProcess`는 `OnClose` 직후 곧바로 세션 인덱스를 재사용 풀에 반환한다. 새 접속이 그 인덱스를 즉시 재할당받는데, `connectedSessions.AddSession`(`sessionManager.go:40-54`)은 반환값 미확인 상태에서 "옛 세션이 아직 `ConnectTimeSec > 0`"이라는 이유로 신규 등록을 조용히 거부할 수 있다. 반대로 지연된 세션 종료 패킷이 나중에 처리되면 `RemoveSession`이 `sessionUniqueId` 검증 없이 슬롯을 무조건 지워, 이미 새 유저가 로그인한 슬롯이 삭제될 수도 있다.
- **개선 방법**: `AddSession`/`RemoveSession`에서 `SetRoomNumber`처럼 `sessionUniqueId` 일치 여부를 검증하고, `AddSession`의 반환값을 호출부에서 반드시 확인하라. 근본적으로는 2.6절처럼 `gohipernetFake` 쪽에서 세션 종료를 멱등적/원자적으로 만드는 것이 우선이다.

### 3.5 [중간] 방 번호 경계 검사 오류
- **파일**: `chatServer/roomPkg/roomManager.go:41-49` (`getRoomByNumber`)
- **내용**: `roomIndex := roomNumber - roomMgr._roomStartNum` 계산 후 `roomNumber < 0`만 검사하고 `roomIndex < 0`은 검사하지 않는다. 현재 `main.go`는 `RoomStartNum`을 0으로 두어 미발현이지만, 설정값을 양수로 바꾸는 순간 배열 음수 인덱스 패닉으로 이어진다.
- **개선 방법**: `if roomIndex < 0 || roomIndex >= _maxRoomCount { return nil, false }`로 고쳐라.

### 3.6 [중간] 패킷 디코딩 실패 반환값 무시
- **파일**: `chatServer/roomPkg/roomManager.go:62-63`
- **내용**: `(&requestPacket).Decoding(packet.Data)`의 bool 반환값을 버린다. 손상된 ENTER_REQ 페이로드가 오면 `RoomNumber`가 0으로 조용히 처리되어 엉뚱한 방으로 라우팅을 시도한다.
- **개선 방법**: 디코딩 실패 시 에러 응답 후 처리를 중단하라.

### 3.7 [낮음] 죽은 코드 정리
`connectedSessions/session.go:121-125`의 `getRoomNumber`가 1.4절에서 지적한 복붙 버그의 원형이다(`_RoomNumOfEntering` 대신 `_RoomNum`을 두 번 반환). 이를 세팅하는 `setRoomEntering`도 어디서도 호출되지 않는다. 그 외 `roomPkg/room.go`의 `EnableEnterUser`, `disConnectedUser`, `secondTimeEvent`, `roomManager.go`의 `GetAllChannelUserCount`, `chatServer/logger.go`의 `init_Log`(어디서도 호출 안 됨)도 미사용 죽은 코드다. 학습용 저장소일수록 "쓰이지 않는 코드"가 헷갈림을 유발하니 정리하거나, 향후 사용 예정이면 TODO 주석으로 의도를 명시하라.

---

## 4. chatServer2 (멀티 고루틴)

chatServer2는 패킷 처리를 N개 고루틴으로 병렬화한 버전이라 **동시성 버그가 가장 집중되는 지점**이다. (참고: README에는 Redis 연동/API 서버 연동이 chatServer2의 추가 기능으로 소개되지만, 현재 `chatServer2` 디렉토리에는 이 코드가 아직 구현되어 있지 않다 — 문서와 코드 상태가 불일치하니 README 갱신이나 구현 계획 확정이 필요하다.)

### 4.1 [높음] `getRoomNumber()` 복붙 버그로 인한 방 유저 슬롯 영구 누수
- **파일**: `chatServer2/connectedSessions/session.go:129-133`
  ```go
  func (session *session) getRoomNumber() (int32, int32) {
      roomNum := atomic.LoadInt32(&session._RoomNum)
      roomNumOfEntering := atomic.LoadInt32(&session._RoomNum) // _RoomNumOfEntering이어야 함
      return roomNum, roomNumOfEntering
  }
  ```
- **재현 시나리오**: 클라이언트가 `ROOM_ENTER_REQ`를 보내면 `SetRoomEntering`으로 `_RoomNumOfEntering`이 먼저 세팅되고, 실제 입장 처리는 방 전용 고루틴 큐에서 비동기로 이뤄진다. 이 패킷이 큐에서 처리되기 **전에** 클라이언트가 접속을 끊으면 `clientSessionEvent.go:46-57`의 `disConnectClient`가 `GetRoomNumber()`로 두 값을 얻는데, 위 버그 때문에 항상 같은 값(둘 다 -1)이 반환돼 `roomNumOfEntering > -1 && roomNum != roomNumOfEntering` 조건이 절대 참이 될 수 없다. 즉 **"입장 처리 중이던 방"에 대한 정리 호출이 발생하지 않는다.**
  이후 세션이 `Clear()`되어 `_networkUniqueID`가 리셋된 뒤, 큐에 남아있던 `ROOM_ENTER_REQ`가 뒤늦게 처리되면(`roomPkg/room_PacketEnterLeave.go:32-42`) `room.addUser()`는 성공(유령 유저 등록, `_curUserCount` 증가)하지만 뒤이은 `connectedSessions.SetRoomNumber`는 uniqueId 불일치로 실패한다. **이 실패 분기가 앞서 추가한 유저를 롤백하지 않는다**(3.3절과 동일한 패턴).
- **왜 문제인가**: 접속 직후 즉시 끊기는 시나리오(모바일 네트워크 전환, 클라이언트 재시도 로직 등에서 흔함)가 반복될 때마다 방에 유령 유저가 하나씩 영구히 쌓인다. 흥미롭게도 코드 작성자 스스로 `room_PacketEnterLeave.go:11`에 "TODO 방 입장 중에 유저가 연결을 끊을 수 있으므로..."라는 주석을 남겨, 이 위험을 인지하고 있었지만 실제로는 막지 못한 상태다.
- **개선 방법**: (1) `getRoomNumber()`의 두 번째 반환값을 `_RoomNumOfEntering`으로 수정한다. (2) `_packetProcess_EnterUser`(방 입장 처리부)에서 `SetRoomNumber` 실패 시 반드시 `room._removeUser(newUser)`로 롤백한다. 두 수정이 함께 있어야 완전히 막힌다.

### 4.2 [높음] 로그인 처리의 TOCTOU 레이스로 중복 로그인 방지가 무력화됨
- **파일**: `chatServer2/connectedSessions/sessionManager.go:162-168`
  ```go
  if _, ok := _manager._UserIDsessionMap.Load(newUserID); ok { return false }
  _manager._sessionList[sessionIndex].SetUser(...)
  _manager._UserIDsessionMap.Store(newUserID, ...)
  ```
- **왜 문제인가**: 패킷 처리가 세션별로 여러 고루틴에서 병렬 실행되므로, 서로 다른 두 연결이 동시에 같은 `userID`로 로그인 요청을 보내면 둘 다 `Load`에서 "없음"을 확인한 뒤 각각 `Store`를 호출할 수 있다(check-then-act 레이스). 결과적으로 두 세션 모두 "로그인 성공" 응답을 받지만 맵에는 마지막 `Store`만 남아 상태가 꼬인다.
- **개선 방법**: `Load` + `Store`를 `sync.Map.LoadOrStore`로 원자화하라. `if _, loaded := _manager._UserIDsessionMap.LoadOrStore(newUserID, session); loaded { return false }` 형태로 바꾸면 두 단계가 원자적으로 처리된다.

### 4.3 [중간] `_userID`/`_userIDLength` 필드가 atomic/락 보호 없이 동시 접근됨 (data race)
- **파일**: `chatServer2/connectedSessions/session.go:14-15, 63-99`
- **내용**: `session` 구조체의 다른 필드(`_RoomNum`, `connectTimeSec` 등)는 전부 `atomic`으로 감쌌는데 `_userID`/`_userIDLength`만 평범한 필드다. 로그인 처리 고루틴이 `setUserID`로 쓰는 동시에, `checkState.go`의 주기적 상태 점검 고루틴이 같은 세션의 `IsAuth()`(내부적으로 `_userIDLength` 읽음)를 호출할 수 있어 `go test -race` 기준 명백한 data race다.
- **개선 방법**: `_userIDLength`를 `int32`로 바꿔 atomic으로 읽고 쓰거나, 이 두 필드를 `sync.RWMutex`로 묶어 보호하라.

### 4.4 [중간] 미로그인 세션에 대한 강제 종료 로직이 사실상 죽어있음
- **파일**: `chatServer2/connectedSessions/checkState.go:84-92, 111-125, 149-162`
- **내용**: `_disablePacketProcess`(149행)는 `DisConnectWaitStartTimeSec`만 세팅할 뿐, 실제 패킷 처리 차단 호출(`NetLibDisablePacketProcessClient`, 161행)은 주석 처리돼 있다. 로그인하지 않은 세션은 `IsLoginUser()`가 항상 false라 `_checkNotLogIn`(119행)만 반복 호출되는데, 이 함수는 `disConnectTime == 0`일 때 딱 한 번만 타임스탬프를 세팅하고 그 이후로는 아무 동작도 하지 않는다.
- **왜 문제인가**: 로그인하지 않은 채 하트비트(핑)만 계속 보내는 클라이언트는 세션 슬롯을 **영구히 점유한 채 강제 종료되지 않는다**. 악의적인 클라이언트가 로그인 없이 다수의 연결을 열어놓기만 해도 `MaxSessionCount`를 소진시킬 수 있다.
- **개선 방법**: 로그인 대기 시간 초과 세션도 `_checkDisConnectWait`와 동일하게 "일정 시간 후 강제 종료" 후속 체크를 받도록 상태 분기를 수정하거나, 로그인 타임아웃 시 즉시 `NetLibForceDisconnectClient`를 호출하라.

### 4.5 [중간] 방 패킷 채널이 무제한 블로킹 (백프레셔 부재)
- **파일**: `chatServer2/roomPkg/roomPacketDistributor.go:64-72, 80-89` (`PushPacket`, `PushInternalPacket`)
- **내용**: 채널 버퍼가 가득 차면 `chanPacket <- packet`에서 호출자 고루틴이 타임아웃 없이 무한정 블록된다.
- **왜 문제인가**: 특정 방(그 방을 처리하는 고루틴)이 느려지거나 트래픽이 몰리면, 그 방에 패킷을 넣으려는 네트워크 콜백 경로가 블록된다. 만약 이 경로가 워커풀을 다른 세션들과 공유하는 구조라면 무관한 세션들까지 연쇄적으로 지연될 위험이 있다.
- **개선 방법**: `select { case chanPacket <- packet: default: /* 드롭 + 로그 또는 에러 응답 */ }`처럼 논블로킹 전송으로 바꾸거나, `context.WithTimeout` 기반의 타임아웃 전송을 적용하라.

### 4.6 [낮음] 미사용 채널로 인한 잠재적 데드락 지뢰
- **파일**: `chatServer2/roomPkg/roomPacketPipe.go:19,29,80-101`
- **내용**: `PushMemberPacket`이 쓰는 `_chanMemebrPacket`을 소비하는 `select` 분기가 방 처리 고루틴(`roomProcess_goroutine_Impl`)에 없다. 현재는 `PushMemberPacket` 호출부 자체가 없어 죽은 코드지만, 나중에 이 API를 실제로 쓰면 채널이 가득 차자마자 호출자가 영구히 블록된다.
- **개선 방법**: `select`에 해당 채널을 소비하는 case를 추가하거나, 당장 쓰지 않을 API라면 제거해 함정을 없애라.

### 4.7 [낮음] 그 외
- `chatServer2/protocol/packetID.go` — `PACKET_ID_PING_REQ`와 `PACKET_ID_PING_RES`가 둘 다 `201`로 중복 정의되어 있다. 서로 다른 값으로 분리하라.
- `roomManager.go`의 `GetRoom`, `GetAllChannelUserCount`, `sessionManager.go`의 `EnableLogin`은 현재 미사용이다. 향후 API 서버 연동 시 이 함수들로 room의 내부 맵에 직접 접근하면, 그 맵은 "방 전용 고루틴을 통한 채널 처리"로만 동시성 안전이 보장되는 구조이므로 즉시 race가 발생한다. API 서버 연동을 실제로 구현할 때는 반드시 채널을 경유하는 요청/응답 방식을 유지하거나 별도 락을 추가하라.
- `session.go:156-166`의 `AddRequestPerSecondTime`은 주석으로 "동시 호출되면 버그"라고 스스로 명시했지만 호출부가 코드베이스에 전혀 없다(요청 속도 제한 기능이 미배선). 배선할 때 atomic 처리를 반드시 함께 넣어라.

---

## 5. chatServer_msgpack / baccaratServer

### 5.1 [매우 높음] 바카라 카드 점수 계산 로직 자체가 틀림 — 게임 결과가 정규 룰과 다르게 나옴
- **파일**: `baccaratServer/roomPkg/Game.go:154-164`
  ```go
  func _baccaratCardIndexToScore(cardIndex int8) int8 {
      score := cardIndex % CARD_ROW_COUNT // CARD_ROW_COUNT = 13
      if score == 0 {
          score = 1
      } else if 10 <= score {
          score = 0
      }
      return score
  }
  ```
- **내용**: 카드 인덱스 0~51을 13(`CARD_ROW_COUNT`)으로 나눈 나머지가 랭크(A,2,...,10,J,Q,K)다. A(랭크 0)→1점, J/Q/K(랭크 10~12)→0점은 맞다. 하지만 **랭크 1~8("2"~"9" 카드)은 실제 점수보다 1 작게 계산**된다(예: "9" 카드가 8점으로 계산됨). 랭크 9("10" 카드)는 `10 <= score` 조건을 만족하지 못해 **0점이 아니라 9점**으로 계산된다.
- **왜 문제인가**: 승패/타이 판정(`_End`, `Game.go:94-102`)이 이 잘못된 점수를 그대로 사용하므로, "2"~"10" 카드가 관여하는 거의 모든 핸드의 결과가 정규 바카라 룰과 다르게 나온다. 이건 성능이나 안정성 문제가 아니라 **게임 그 자체가 잘못 동작하는** 가장 심각한 버그다.
- **개선 방법**:
  ```go
  func _baccaratCardIndexToScore(cardIndex int8) int8 {
      rank := cardIndex % CARD_ROW_COUNT // 0=A, 1~8="2"~"9", 9="10", 10~12=J/Q/K
      switch {
      case rank == 0:
          return 1
      case rank >= 9:
          return 0
      default:
          return rank + 1
      }
  }
  ```
  이 함수 하나만 있는 순수 함수이므로, 1.1절에서 언급한 것처럼 `cardIndex 0~51` 전수 테이블 테스트를 붙이면(기대값: A=1, 2~9=2~9, 10/J/Q/K=0) 회귀를 영구히 방지할 수 있다.

### 5.2 [높음] 시간 단위 불일치(초 vs 밀리초)로 배팅/결과 대기시간이 1000배 늘어남
- **파일**: `baccaratServer/roomPkg/Game.go:32-34, 89`, `baccaratServer/roomPkg/room.go:67`, `baccaratServer/distributePacket.go:71`
- **내용**: `setBattingWaitTime`/`doBaccarat`는 이름과 주석 그대로 "밀리초" 단위 상수(`BATTING_WAIT_MILLISEC = 5000`, `NEXT_GAME_WAIT_MILLISEC = 10000`, `roomPkg/define.go:36-37`)를 더한다. 그런데 실제 호출부(`room.go:67`의 `setBattingWaitTime(time.Now().Unix())`, `Game.go:89`의 `time.Now().Unix() + NEXT_GAME_WAIT_MILLISEC`, `distributePacket.go:71`의 `CheckRoomState(curTime.Unix())`)는 전부 `time.Now().Unix()`(**초 단위**)를 넘긴다.
- **왜 문제인가**: 5초여야 할 배팅 대기 시간이 5000초(약 83분)가 되고, 10초여야 할 결과 표시 시간이 10000초(약 166분)가 된다. 게임이 사실상 멈춘 것처럼 보이는, 플레이 자체가 불가능한 수준의 기능 버그다.
- **개선 방법**: `time.Now().Unix()`를 전부 `time.Now().UnixMilli()`로 바꾸거나(Go 1.17+ 필요 — 1.3절의 go.mod 업그레이드와 함께 적용), 반대로 상수를 초 단위(`BATTING_WAIT_SEC = 5`, `NEXT_GAME_WAIT_SEC = 10`)로 바꾸고 변수/함수명에서 "MilliSec" 표기를 전부 제거하라. 어느 쪽이든 **단위를 코드 전체에서 하나로 통일**하는 것이 핵심이다.

### 5.3 [높음] msgpack UserID 길이 미검증 → 미인증 상태에서 서버 패닉 유발(DoS)
- **파일**: `chatServer_msgpack/connectedSessions/session.go:55-62` (`setUserID`, `getUserID`), `chatServer_msgpack/distributePacket.go:73-88` (`ProcessPacketLogin`)
- **내용**: 로그인 요청은 `msgpack.Unmarshal`로 임의 길이의 `UserID string`을 길이 검증 없이 받아 그대로 `SetLogin` → `setUserID`에 넘긴다.
  ```go
  func (session *session) setUserID(userID []byte) {
      session._userIDLength = int8(len(userID)) // 16바이트 초과 시 오버플로/음수 가능
      copy(session._userID[:], userID)          // 배열은 16바이트로 잘림
  }
  func (session *session) getUserID() []byte {
      return session._userID[0:session._userIDLength] // 길이 > 16 또는 음수면 panic
  }
  ```
- **왜 문제인가**: UserID를 16바이트보다 길게(17~127바이트, 또는 128바이트 이상이면 `int8` 오버플로로 음수) 보내면 이후 `getUserID()`에서 `slice bounds out of range` 패닉이 발생한다. **로그인 자체가 미인증 단계에서 처리되므로, 인증 없이 누구나 트리거 가능한 DoS다.** 참고로 `baccaratServer`는 고정 16바이트 배열이지만 `len(bodyData) != bodySize` 같은 길이 검증이 있어 이 문제가 없다.
- **개선 방법**: `setUserID` 진입부에서 `if len(userID) > 16 { /* 에러 응답 후 연결 종료, 또는 truncate */ }`로 방어하라. chatServer_msgpack에도 baccaratServer와 같은 수준의 길이 검증을 추가하면 된다.

### 5.4 [중간] 패킷 ID 상수 중복 정의
- **파일**: `baccaratServer/protocol/packetID.go:34,38`
- **내용**: `PACKET_ID_GAME_START_NTF = 753`과 `PACKET_ID_GAME_BATTING_NTF = 753`이 동일한 값이다.
- **왜 문제인가**: 클라이언트가 두 알림을 패킷 ID만으로 구분할 수 없어 게임 시작 통지와 배팅 통지 처리가 서로 오염된다.
- **개선 방법**: 둘 중 하나를 754 등 미사용 값으로 바꿔라. (재발 방지를 위해 1.2절의 CI에 "패킷 ID 중복 검사" 스크립트를 추가하는 것도 고려하라 — `packetID.go`의 상수를 파싱해 값 중복을 검사하는 간단한 Go 스크립트로 충분하다.)

### 5.5 [중간] 패킷 바디 크기 계산의 언더플로/음수 미검증
- **파일**: `chatServer_msgpack/protocol/packet.go:69-79`, `baccaratServer/protocol/packet.go:68-78`, 양쪽 `distributePacket.go`의 `packet.Data = make([]byte, packet.DataSize)`
- **내용**: `bodySize := totalSize - headerSize` 계산에서 `totalSize < headerSize`(변조/손상 패킷)일 때의 검증이 없다. `baccaratServer`는 `int16`이라 음수가 그대로 `make([]byte, 음수)`에 들어가 **panic**이 나고, `chatServer_msgpack`은 `uint16`이라 언더플로로 거대한 양수가 되어 데이터가 손상된다(패딩된 0바이트를 진짜 바디로 오인).
- **개선 방법**: `PeekPacketBody`(또는 동급 함수)에 `if totalSize < headerSize { return nil, false }` 방어 코드를 명시적으로 추가하라. 네트워크 계층(`gohipernetFake`, 2.2절)에서 걸러준다는 보장이 코드상 없으므로, 애플리케이션 계층에서도 독립적으로 검증해야 한다(3.1절과 같은 원칙).

### 5.6 [중간] 방별 카드 셔플 시드가 충돌할 수 있음
- **파일**: `baccaratServer/roomPkg/Game.go:15-17`
  ```go
  game._rand = rand.New(rand.NewSource(1))
  game._rand.Seed(time.Now().UTC().UnixNano()) // 다른 방의 시드 값과 같을 듯
  ```
  (주석에서 개발자 스스로 위험성을 인지하고 있다.)
- **왜 문제인가**: 서버 시작 시 방 개수만큼 연속 루프로 초기화되는데(`roomManager.go`), 나노초 타이머의 해상도가 거친 환경(특히 Windows)에서는 다수의 방이 동일하거나 근접한 나노초 시드를 공유해 카드 순서가 서로 상관관계를 가질 수 있다. 갬블 게임 학습 예제에서 "카드 순서 예측 가능성"은 다뤄볼 가치가 있는 주제다.
- **개선 방법**: `crypto/rand`로 시드용 난수를 생성하거나, 방 인덱스를 시드 계산에 섞어라(예: `seed := time.Now().UnixNano() ^ int64(roomIndex)*someLargePrime`).

### 5.7 [낮음] 그 외
- `connectedSessions/session.go:121-125`(양쪽 디렉토리 동일)의 `getRoomNumber` 복붙 버그(1.4/4.1절과 동일 패턴). 현재는 반환값이 호출부에서 버려지고 있어 실사용 영향은 없지만, "입장 중" 상태 추적 기능이 사실상 구현되어 있지 않다는 뜻이므로 4.1절 수정과 함께 정리하라.
- chatServer_msgpack의 `MAX_CHAT_MESSAGE_BYTE_LENGTH`/`ERROR_CODE_ROOM_CHAT_CHAT_MSG_LEN`은 정의만 있고 `room_PacketChat.go`에서 실제 검증 코드가 없다(baccaratServer에는 있음). 채팅 메시지 길이 제한을 실제로 적용하라.
- README가 설명하는 baccaratServer의 "Scale-Out 기능", "API Server 연동(유저 할당/매칭)"은 현재 코드베이스에 존재하지 않는다(관련 키워드로 검색해도 없음). 문서와 코드 상태가 어긋나 있으니, 이 기능을 실습 다음 단계로 남겨둘 계획이라면 README에 "(예정)"이라고 명시하거나, 이미 구현했다고 착각하지 않도록 상태를 갱신하라.
- chatServer_msgpack/baccaratServer의 `connectedSessions`, `roomPkg` 뼈대 코드가 chatServer/chatServer2와 거의 동일하게 복붙되어 있다. 1.4절의 공통 모듈화 제안이 이 두 서버에도 그대로 적용된다.

---

## 6. 우선순위별 개선 로드맵

학습용 저장소의 성격을 고려해, "지금 당장 실습이 막히는 것"부터 "구조적으로 다음 학습 주제로 삼을 만한 것" 순으로 정리했다.

### 1단계 — 즉시 수정 (모든 서버가 살아있기 위한 최소 조건)
`gohipernetFake`만 고치면 5개 서버 전부에 적용된다.
- 2.1 `_InitNetworkSendFunction()` 호출 배선
- 2.2 패킷 길이 부호 반전 + 하한 검증 추가
- 2.3 `handleTcpRead`에 세션 단위 panic recover 추가
- 2.4 `Accept()` 에러 처리
- 2.6 `closeProcess()` 멱등성 보장(`sync.Once`/CAS)

### 2단계 — 데이터 무결성/보안
- 3.1, 5.5 애플리케이션 계층 패킷 길이 방어 추가
- 3.2 로그인 인증 확인 없이 방 입장 가능한 취약점 수정
- 5.3 msgpack UserID 길이 검증 추가
- 2.8 `ReadBytes` 범위 검사 추가

### 3단계 — 자원 누수/동시성 정합성
- 3.3, 4.1 방 입장 실패 시 롤백 처리
- 4.2 로그인 TOCTOU 레이스 (`LoadOrStore`로 수정)
- 4.3 `_userID` 필드 data race 해소
- 4.4 미로그인 세션 강제 종료 로직 복구
- 2.5, 2.9 `addSession`/`sendPacket` 반환값 처리

### 4단계 — 게임 로직 정확성 (baccaratServer)
- 5.1 카드 점수 계산 로직 수정 (가장 눈에 띄는 "기능이 틀린" 버그)
- 5.2 시간 단위(초/밀리초) 통일
- 5.4 패킷 ID 중복 해소

### 5단계 — 구조 개선 (다음 학습 주제로 적합)
- 1.1 핵심 순수 함수(패킷 인/디코딩, 카드 점수 계산)에 단위 테스트 추가
- 1.2 GitHub Actions CI 추가 (`go build`, `go vet`)
- 1.4 `connectedSessions`/`roomPkg` 공통 뼈대를 별도 모듈로 추출
- 1.3 `go.mod`를 최신 Go 버전으로 업그레이드하고 `golangci-lint` 도입
- 3.7, 4.6, 5.7의 죽은 코드 정리 및 README-코드 상태 동기화

1~3단계는 "네트워크 계층이 신뢰할 수 있고, 인증/자원 관리가 올바른가"를 다루므로 실습생이 다음 단계(chatServer2, baccaratServer)로 넘어가기 전에 chatServer 단계에서 함께 고쳐보면 좋은 학습 소재다. 4단계는 baccaratServer 실습 중 "게임 로직 디버깅" 연습으로 적합하고, 5단계는 여러 서버를 다 만들어본 뒤 "지금까지 만든 걸 리팩토링해보기" 실습으로 자연스럽게 이어진다.
