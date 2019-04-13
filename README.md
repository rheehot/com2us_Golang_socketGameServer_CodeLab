# golang_socketGameServer_codelab
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
     
  	 
## echoServer
- 디렉토리: echoServer
- GoLand를 사용하여 golang용 프로그램을 만들고, 빌드/디버깅을 한다.
- 아주 간단한 규모이다.
  
  
## chatServer
- 디렉토리: chatServer
- 방 개념의 채팅 서버
- 패킷 요청 처리를 1개의 고루틴(스레드)에서만 한다.
- echoServer에 비해 규모는 3~4배 크다.
  
### 추가 기능 구현
- 1:1 귓속말
- 방 초대
    
  
## chatServer2
- 디렉토리: chatServer2
- 방 개념의 채팅 서버
- 패킷 요청 처리를 N개의 고루틴(스레드)에서 한다.
    - 패킷 처리를 멀티 고루틴에서 하므로 공유 객체 동기화를 조심해야 한다.
- chatServer의 코드와 겹치는 부분이 많으므로 chatServer에 대한 이해가 꼭 필요하다
    
### 추가 기능 구현
- Redis 연동
- API Server(http)와 연동  
    - 로그인을 API Server에서 한다.  
    
	
## baccaratServer 
- 디렉토리: baccaratServer
    - 2019.04.05 아직 구현하지는 않았다. 위 단계 실습 후 결정할 예정
- 겜블 게임인 바카라 게임을 온라인화 한 것이다.
    - 바카라 룰: https://namu.wiki/w/%EB%B0%94%EC%B9%B4%EB%9D%BC
- chatServer2에 바카라 게임 로직이 올라간 것으로 chatServer2에 대한 이해가 꼭 필요하다.
  
### 추가 기능 구현
- 게임 서버 Scale-Out 기능 구현
-  API Server(http)와 연동  
    - 유저를 특정 게임 서버에 할당하는 기능
    - 매칭 기능	
  
    
## TCP 내트워크 프로그래밍과 GoHiperNet 
- 이 [저장소](https://github.com/jacking75/golang_network_socket_study)에 있는 코드와 글을 기반으로 소켓 프로그래밍 방법을 배운다.
  