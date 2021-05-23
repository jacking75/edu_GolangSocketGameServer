# Socket API
socket api 관련 문서인 [.net](https://golang.org/pkg/net)를 보면 API의 상세한 설명을 볼 수 있다.  
[Network programming with Go](https://www.joinc.co.kr/w/man/12/golang/networkProgramming)  
[socket API 설명 문서](https://docs.google.com/document/d/1ThEKfZBeVnViRpkJeNOALVPMqklIl5kEc7gFZUEI1KE/edit?usp=sharing  )    

<br>    

## 연결 하기(server) - listen, accept
- IP는 `0.0.0.0`에 지정된 port로 listen을 하고 접속을 받는다      
  
```  
port := ":8080"
ln, err := net.Listen("tcp", port)
if err != nil {
    log.Fatal(err)
}

conn, err := l.Accept()
```    

### listen
  
```
tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:12345")
if err != nil {
    log.Println("ResolveTCPAddr", err)
    return
}

l, err := net.ListenTCP("tcp", tcpAddr)
if err != nil {
    log.Println("ListenTCP", err)
    return
}
```
  


### accept      
accept는 `Accept` 와 `AcceptTCP` 2가지가 있다.    
```func (l *TCPListener) Accept() (Conn, error)```    
`Conn`의 [API 문서](https://golang.org/pkg/net/#Conn).   
`Conn`은 소켓 옵션 설정 등 세세한 동작을 할 수 없다.    
  
<br>  

```func (l *TCPListener) AcceptTCP() (*TCPConn, error)```  
`TCPConn`의 [API 문서](https://golang.org/pkg/net/#TCPConn)      

<br>  
<br>  

## 연결 하기(client) - connect
- `ResolveTCPAddr`로 연결 주소를 만든다
- `DialTCP`로 연결한다
   
```
tcp_addr, err := net.ResolveTCPAddr("tcp", "localhost:6666")
if err != nil {
    println("error tcp resolve failed", err.Error())
    os.Exit(1)
}
tcp_conn, err := net.DialTCP("tcp", nil, tcp_addr)  
```  
  
<br>  

## 데이터 보내기 
- `Write` 로 보낸다
- 송신 버퍼가 다 찬 경우 게속 시도를 해서 대기 상태에 빠질 수 있으므로 주의해애 한다
    
```
func SendEcho(conn *net.TCPConn, msg string) {
    _, err := conn.Write([]byte(msg))
    if err != nil {
        println("Error send request:", err.Error())
    } else {
        println("Request sent")
    }
}
```  
  
<br>    

## 데이터 받기
- `Read` 로 받는다
      
```
func GetEcho(conn *net.TCPConn) string {
    buf_recever := make([]byte, RECV_BUF_LEN)
    _, err := conn.Read(buf_recever)
    if err != nil {
        println("Error while receive response:", err.Error())
        return ""
    }
    return string(buf_recever)
}
```

<br>    
  
## 연결 끊기
`Close` 함수를 사용한다  
  
  
<br>  

## 일시적인 에러 무시하기  
`net.Error`의 `Temporary()` 함수를 사용하여 일시적인 에러는 다시 시도하도록 한다.   
    
```
func handleListener(l *net.TCPListener) {    
    for {
        conn, err := l.AcceptTCP()
        if err != nil {
            if ne, ok := err.(net.Error); ok {
                if ne.Temporary() {
                    log.Println("AcceptTCP", err)
                    continue
                }
            }
            
            log.Println("AcceptTCP", err)
            return
        }
    }
}
```

- `AcceptTCP()`, `Read()` 각각의 [내부 구현](libexec/src/net/tcpsock.go)을 보면, 발생할 수 있는 오류는 `syscall.EINVAL`, `*net.OpError` 중 하나이다.
- `*net.OpError`는 `func(e *OpError) Timeout() bool` 과 `func (e *OpError) Temporary() bool` 메소드를 가지고 있으며, 복구 가능한 오류를 쉽게 알 수 있다.
- `*net.OpError`는 `net.Error`를 충족하므로 `net.Error`로 캐스팅 하여 오류 유형을 판별한다.
  
<br>  

## 커널의 Socket Buffer 크기 설정하가  
`SetWriteBuffer` , `SetReadBuffer`    
```
func setBuffer(conn *net.TCPConn) {
	conn.SetWriteBuffer(4000)
	conn.SetReadBuffer(4000)
}
```  
  
  
## backlog 설정
- https://qiita.com/kawasin73/items/7a24077fa3f89ce240c3
    - 이 글을 보면 backlog 수를 리눅스의 설정 정보에 따라가는 것 같음 -_-
	- 현재 backlog 수  `$ cat /proc/sys/net/core/somaxconn`
	- backlog 수 늘리기. `$ sysctl -w net.core.somaxconn=1024`
	- (중국어 문서)[이 글](https://blog.csdn.net/Neuliudapeng/article/details/73106809)을 보면 맞는 듯
  
  