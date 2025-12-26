# Plugins

noABI는 기본적인 Compose / Deploy / Contract Execution 기능 외에도  
확장 가능하도록 설계되어 있으며, 선택적으로 활성화할 수 있는 플러그인을 제공합니다.  
아래 플러그인들은 프로젝트 또는 배포 환경에 따라 개별적으로 독립 사용이 가능합니다.

---

## 1. EIP-712 / EIP-191 Signing Backend

noABI는 EIP-712(typed structured data) 및 EIP-191(personal_sign) 메시지 서명을  
로컬에서 자동 처리하기 위한 Signing Backend 플러그인을 제공합니다.

### 1.1 개요

EIP-712 / EIP-191 Signing Backend는 서명 기반 프로토콜 테스트를 위한 통합 솔루션입니다.  
별도의 백엔드 서버 구축 없이 noABI 패키지에 포함된 플러그인을 실행하여 서명 생성 및 검증을 수행할 수 있습니다.

**플러그인 위치**  
```
plugins/backend/
├─ src/
│  ├─ api/
│  │  └─ routers/
│  │     └─ v3/
│  │        └─ man/
│  │           ├─ routers.eip712DataSign.go     ← EIP-712 서명 엔드포인트
│  │           ├─ routers.eip191DataSign.go     ← EIP-191 서명 엔드포인트
│  │           └─ routers.recoverDataSigner.go  ← 서명자 복구 엔드포인트
│  ├─ common/
│  │  └─ blockchain/
│  │     ├─ signature/                          ← 서명 처리 로직
│  │     └─ keystore/                           ← 키스토어 관리
│  └─ build/
│     └─ signer/
│        └─ main.go                             ← 플러그인 진입점
└─ artifacts/
   └─ imports/                                  ← 설정 파일
```

### 1.2 목적

**개발 편의성 향상**
- 매번 지갑(Metamask 등)을 열어 서명 팝업을 확인해야 하는 번거로움 제거
- 프론트엔드나 dApp 흐름에 종속되지 않고 단독으로 서명 테스트 가능
- CI/CD나 스크립트 기반 테스트에서 자동 서명 지원

**통합 솔루션 제공**
- 별도 백엔드 서버 구축 불필요
- noABI 패키지에 포함된 플러그인으로 즉시 사용 가능
- 로컬 환경에서 완전히 독립적으로 동작

### 1.3 주요 기능

**서명 생성**
- ✅ EIP-712(typed structured data) 서명 생성
- ✅ EIP-191(personal_sign) 서명 생성
- ✅ 테스트 환경에서 지갑 없이 서명 가능
- ✅ 버튼 클릭 없이 자동 서명 반환

**API 엔드포인트**
- ✅ RESTful API로 서명 요청 처리
- ✅ curl 또는 외부 클라이언트(Postman, Golang, Node.js 등) 요청 가능
- ✅ 프론트엔드와 독립적으로 동작하는 서명 엔드포인트 제공

**통합 테스트**
- ✅ 스마트 컨트랙트 인터랙션 테스트 시 UI 없이 자동 진행 가능
- ✅ noABI DSL에서 curl 명령어로 직접 호출
- ✅ 서명 생성부터 컨트랙트 검증까지 하나의 스크립트로 자동화

### 1.4 사용 예시

**플러그인 실행**
```bash
# plugins/backend/src/build/signer/main.go 빌드 후 실행

```

**noABI DSL에서 사용**
```bash

set --var {"verifyingAddr": "0x1234...."}

# EIP-712 서명 요청
curl POST http://127.0.0.1:8080/v1/sign/eip712
  -H "Content-Type: application/json"
  -d '{
        "requestId": "req-1234",
        "domain": {
          "chainId": "1338",
          "name": "NOABI",
          "version": "V0.6.0",
          "verifyingContract": "${verifyingAddr}"
        },
        "primaryType": "Permit",
        "data": [
          {
            "type": "address",
            "name": "owner",
            "value": "0xb171fe0b0804651446a50344ae14e56596190bcf"
          },
          {
            "type": "uint256",
            "name": "value",
            "value": "1000000"
          }
        ],
        "KeyPair": [
          {
            "account": "0xb171fe0b0804651446a50344ae14e56596190bcf",
            "phrase": "old"
          }
        ]
      }'
      => res;

# 서명 검증 (컨트랙트에서)
eip712.EIP712RecoverTest("${verifyingAddr}").recoverSigner(
    "0xb171fe0b0804651446a50344ae14e56596190bcf",
    "1000000",
    "${res.result.data.signatures.0}"
  );
```

### 1.5 지원 기능

| 기능 | 설명 | 엔드포인트 |
|------|------|-----------|
| **EIP-712 서명** | Typed structured data 서명 생성 | `POST /v1/sign/eip712` |
| **EIP-191 서명** | Personal sign 메시지 서명 생성 | `POST /v1/sign/eip191` |
| **서명자 복구** | 서명에서 서명자 주소 복구 | `POST /v1/recover/signer` |

### 1.6 설정

플러그인 설정 파일은 `plugins/backend/src/build/signer/config/local/` 디렉토리에 위치합니다.

**주요 설정 파일**
- `preference.yaml` - 플러그인 기본 설정 (서버 포트, 로그 레벨 등)
- `log.yaml` - 로깅 설정

---

## 2. Event Log 수집기

noABI는 블록체인에서 발생하는 스마트 컨트랙트 이벤트를 자동으로 수집하고 저장하기 위한 Event Logger 플러그인을 제공합니다.

### 2.1 개요

Event Log 수집기는 블록체인 이벤트 로그를 자동으로 수집, 저장, 처리하는 통합 솔루션입니다.  
별도의 인덱서나 외부 서비스 없이 로컬에서 이벤트 로그를 수집하여 파일로 저장하고, 필요 시 데이터베이스에 저장할 수 있습니다.

**플러그인 위치**  
```
plugins/backend/
├─ src/
│  ├─ api/
│  │  └─ routers/
│  │     └─ v1/
│  │        └─ eventLogger/
│  │           └─ routers.handler.go     ← API 엔드포인트 핸들러
│  ├─ pkg/
│  │  └─ eventLogger/
│  │     ├─ eventCollector/              ← 이벤트 수집 로직
│  │     └─ eventActor/                  ← 이벤트 처리 로직
│  └─ build/
│     └─ eventLogger/
│        ├─ main.go                      ← 플러그인 진입점
│        └─ config/
│           └─ local/
│              ├─ preference.yaml        ← 플러그인 설정
│              └─ log.yaml                ← 로깅 설정
└─ artifacts/
   └─ imports/
      └─ local/
         └─ eventLogger.yaml              ← 이벤트 로거 설정
```

### 2.2 목적

**자동 이벤트 수집**
- 블록체인에서 스마트 컨트랙트 이벤트를 자동으로 수집
- 지정된 블록 범위에서 이벤트 로그를 주기적으로 조회
- 수집된 이벤트를 JSON 파일로 저장

**이벤트 로그 관리**
- 컨트랙트별, 체인별로 이벤트 로그 분리 저장
- 블록 번호 기반 이벤트 로그 파일 관리
- 이전 tx의 이벤트 결과를 다음 tx 입력으로 활용 가능

**통합 테스트 지원**
- noABI의 `event` 명령어와 연동하여 이벤트 로그 조회
- E2E 시나리오 테스트에서 이벤트 기반 상태 추적
- 배포 및 Contract Call 결과를 이벤트로 추적

### 2.3 주요 기능

**이벤트 수집 (Event Collector)**
- ✅ 블록체인 RPC를 통한 이벤트 로그 수집
- ✅ 지정된 컨트랙트 주소의 이벤트만 필터링
- ✅ 블록 범위 기반 이벤트 조회
- ✅ 주기적으로 최신 블록까지 이벤트 수집
- ✅ 수집된 이벤트를 JSON 파일로 저장

**이벤트 처리 (Event Actor)**
- ✅ 수집된 이벤트 로그 파일 읽기
- ✅ 이벤트 디코딩 및 구조화
- ✅ 데이터베이스 저장 (선택적)
- ✅ 주기적으로 이벤트 처리 작업 실행

**체인 관리 API**
- ✅ RESTful API로 체인 추가/수정/삭제
- ✅ 동적으로 체인 설정 변경 가능
- ✅ 체인별 설정 조회

### 2.4 사용 예시

**플러그인 실행**
```bash
# plugins/backend/src/build/eventLogger/main.go 빌드 후 실행
# 기본 포트: 8080
```

**noABI DSL에서 사용**
```bash
# event 명령어로 수집된 이벤트 로그 조회
event -load C:/data/eventLogger/event/local/sample
  --block 181816793
  --contractName ERC20
  --eventName Transfer
  => events;

# 조회된 이벤트 활용
echo "${events[0].decoded.from}";
echo "${events[0].decoded.to}";
echo "${events[0].decoded.value}";
```

**API를 통한 체인 관리**
```bash
# 모든 체인 정보 조회
curl GET http://127.0.0.1:8080/v1/event/chains => chains;

# 특정 체인 정보 조회
curl GET http://127.0.0.1:8080/v1/event/chains/local => chain;

# 새 체인 추가 (dev 모드만 가능)
curl POST http://127.0.0.1:8080/v1/event/chains/testnet
  -H "Content-Type: application/json"
  -d '{
  "mode": "dev",
  "chainId": "0x539",
  "rpc": "http://127.0.0.1:8545",
  "path": "c:/data/eventLogger/event",
  "commitBlockCount": 0,
  "period": 5,
  "logRange": 10000,
  "collections": {
    "ERC20": {
      "address": "0xa992e822040fbff8392ba3248bd7de55da77ecb9",
      "enable": true,
      "interface": "ERC20"
    }
  },
  "collectors": {
    "sample": {
      "startBlockNumber": 0,
      "collections": [
        "ERC20"
      ]
    }
  }
}' => result;

# 체인 제거
curl GET http://127.0.0.1:8080/v1/event/remove/testnet => result;
```

### 2.5 지원 기능

| 기능 | 설명 | 엔드포인트 |
|------|------|-----------|
| **체인 목록 조회** | 등록된 모든 체인 정보 조회 | `GET /v1/event/chains` |
| **체인 정보 조회** | 특정 체인 상세 정보 조회 | `GET /v1/event/chains/:alias` |
| **체인 추가/수정** | 새 체인 추가 또는 기존 체인 수정 | `POST /v1/event/chains/:alias` |
| **체인 제거** | 등록된 체인 제거 | `GET /v1/event/remove/:alias` |

### 2.6 설정

플러그인 설정 파일은 `plugins/backend/src/build/eventLogger/config/local/` 디렉토리에 위치합니다.

**주요 설정 파일**
- `preference.yaml` - 플러그인 기본 설정 (서버 포트, 로그 레벨 등)
- `log.yaml` - 로깅 설정

**이벤트 로거 설정 파일**
이벤트 로거 설정은 `plugins/backend/artifacts/imports/local/eventLogger.yaml`에 위치합니다.

**설정 예시**
```yaml
event-logger:
  version: "V2.0.0"
  chains:
    - "local"
  "local":
    mode: live                    # live: 운영 모드, dev: 개발 모드
    chainId: "0x539"              # 체인 ID
    rpc: http://127.0.0.1:8545   # RPC 엔드포인트
    db: http://127.0.0.1:15000/dc/v1/documents/queries/op  # DB 엔드포인트 (선택)
    path: "c:/data/eventLogger/event"  # 이벤트 로그 저장 경로
    commitBlockCount: 0          # 커밋 블록 수
    period: 5                     # 수집 주기 (초)
    logRange: 10000               # 한 번에 수집할 최대 블록 수
    collections:                  # 수집할 컨트랙트 목록
      ERC20: 
        address: "0xa992e822040fbff8392ba3248bd7de55da77ecb9"
        enable: true
        interface: "ERC20"
    collectors:                  # 수집기 설정
      sample:
        startBlockNumber: 0       # 시작 블록 번호
        collections:              # 수집할 컨트랙트 이름 목록
          - ERC20
    actors:                       # 액터 설정 (이벤트 처리)
      actor1:
        collectors:               # 사용할 수집기 목록
          - sample
```

**모드 설명**
- **live 모드**: 운영 모드, DB 저장 및 Actors 활성화, 설정 변경 불가
- **dev 모드**: 개발 모드, DB 저장 및 Actors 비활성화, 설정 변경 가능

**주요 설정 항목**
- `period`: 이벤트 수집 주기 (초 단위)
- `logRange`: 한 번에 수집할 최대 블록 수
- `startBlockNumber`: 수집 시작 블록 번호
- `collections`: 수집할 컨트랙트 주소 및 인터페이스 설정
- `collectors`: 수집기 그룹 설정
- `actors`: 이벤트 처리 액터 설정

### 2.7 이벤트 로그 파일 구조

수집된 이벤트 로그는 다음 경로 구조로 저장됩니다:
```
{path}/{chainAlias}/{collectorName}/{blockNumber}.json
```

**파일 내용 예시**
```json
{
  "blockHash": "0x9a251c85d49b0b695ee56278515b7de7659969cb72d38eb26256f10b517614f6",
  "blockNumber": "181816793",
  "eventLogs": [
    {
      "blockHash": "0x9a251c85d49b0b695ee56278515b7de7659969cb72d38eb26256f10b517614f6",
      "blockNumber": "181816793",
      "contractAddress": "0xd0a549cbe7d9605e119db3a2a4a3938cdcb6eac7",
      "contractName": "ERC20",
      "data": "0x000000000000000000000000fc662e967d2054973b3982bc3ae82fc40a168034...",
      "decodeLog": {
        "from": "0xfc662e967d2054973b3982bc3ae82fc40a168034",
        "to": "0x5678...",
        "value": "1000000"
      },
      "eventName": "Transfer",
      "logIndex": "3",
      "topics": ["0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"],
      "transactionHash": "0x9a67fadb5c6292e46c7714e0a15bb2a7e2e1609396e30e67bff30a9bb4889b9f",
      "transactionIndex": "2"
    }
  ]
}
```

### 2.8 noABI와의 연동

Event Logger 플러그인은 noABI의 `event` 명령어와 연동하여 사용할 수 있습니다.

**연동 흐름**
1. Event Logger가 블록체인에서 이벤트를 수집하여 파일로 저장
2. noABI DSL에서 `event` 명령어로 저장된 이벤트 로그 조회
3. 조회된 이벤트를 `varStore`에 저장하여 다음 tx에 활용

**사용 예시**
```bash
# 1. 이벤트 로그 수집 (Event Logger 플러그인 자동 실행)

# 2. noABI에서 이벤트 조회
event -load C:/data/eventLogger/event/local/sample
  --block ${currentBlock}
  --contractName ERC20
  --eventName Transfer
  => transferEvents;

# 3. 조회된 이벤트 활용
tokens.ERC20("${tokenAddress}").transfer(
  "${transferEvents[0].decoded.to}",
  "${transferEvents[0].decoded.value}"
);
```

---
**문서 네비게이션**  
[← 이전: Commands](commands.md) | [다음: README →](../README.md)  
[↑ 목차로 돌아가기](../README.md)



