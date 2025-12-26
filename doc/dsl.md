## Domain-Specific Language (DSL)
noABI DSL은 스마트 컨트랙트 배포, ABI 관리, Contract Call 테스트를 스크립트 기반으로 자동화하기 위한 Domain-Specific Language입니다.
이 챕터는 DSL의 문법 구조, 명령어 사용법, Contract Call 규칙을 공식적으로 정의합니다.

---

### 1 noABI DSL 문법 규칙
스마트 컨트랙트 개발 환경을 안정적으로 구성하려면 아래 버전을 권장합니다:   
***1 Statement = 1 Line 구조***
- 일반 명령어는 반드시 한 줄로 작성
- 줄바꿈 불가, 멀티라인 금지
- 예외: curl, event와 Contract Call만 멀티라인 허용

**✅ 올바른 예시**
```bash
compose ./compose-file/noAbi/compose_tokens.json
deploy ./compose-file/noAbi/compose_tokens.json
import -abi --collection tokens --name CERC721 ./abi/CERC721.abi

```

**❌ 잘못된 예시**
```bash
import -abi
  --collection tokens 
  --name CERC721 
  ./abi/CERC721.abi
```
---
curl 명령 멀티라인 규칙
구조적 요청을 위해 멀티라인 허용
명령 종료 시 반드시 세미콜론(;) 포함

**✅ 올바른 예시**
```bash
curl GET
    http://localhost:8080/management/health/ping
    -H "Content-Type: application/json" => res;

curl http://localhost:8080/management/health/ping;
```

**❌ 잘못된 예시**
```bash
curl http://localhost:8080/management/health/ping
```

---
event 명령 멀티라인 규칙
여러 옵션을 명확하게 구분하기 위해 멀티라인 허용
명령 종료 시 반드시 세미콜론(;) 포함

**✅ 올바른 예시**
```bash
event -load C:/workspace/log/event
  --block 181816793
  --contractName DAO
  --eventName eventApproval
  => events;

event -load C:/workspace/log/event --block 181816793 => events;
```

**❌ 잘못된 예시**
```bash
event -load C:/workspace/log/event --block 181816793
```

---
Contract Call 멀티라인 규칙
JS-like method chain 구조 허용
마지막 줄 끝에는 반드시 세미콜론(;) 필요

**✅ 올바른 예시**
```bash
tokens.CERC20("${addr}")
  .balanceOf("${USER}")
  .then(res => {});

tokens.CERC20("${addr}").balanceOf("${USER}");
```

**❌ 잘못된 예시**
```bash
tokens.CERC20("${addr}")
  .transfer("${USER}", 1000)
```

---
### 2 Contract Instantiation 규칙
noABI에서 Contract는 호출 시점에 즉시 인스턴스화되며,    
두 가지 방식의 인스턴스 생성 규칙을 제공합니다.

**Address 직접 전달 방식 (기본)**
```bash
<collection>.<contract>(address).<function>(...);
```
- 가장 기본적인 방식
- 단순 호출 및 조회에 적합

**예시**
```bash
tokens.CERC20("0x1234...").transfer("${USER}", 1000);
```

---
**Object Parameter 방식 (확장)**
```bash
<collection>.<contract>({ to: address, value: amount }).<function>(...);
```
- 트랜잭션 옵션(to, value 등)을 명시적으로 전달
- value 전송이 필요한 함수 호출에 사용

**예시**
```bash
tokens.CERC20({
  to: "${address}",
  value: "1000000000000000000"
}).approve("${USER}", "1000");
```

---
### 3 Contract Call 결과 및 오류 처리
Contract Call은 결과 처리(.then) 와 오류 처리(.catch) 를
선택적으로 연결할 수 있습니다.

**3.1 기본 개념**   
- .then() / .catch() 는 Contract Call 뒤에서만 사용 가능
- 둘 다 선택 사항
- 사용하지 않아도 Contract Call 자체는 정상 실행됨

| 상황            | 동작                    |
| ------------- | --------------------- |
| `.then()` 없음  | 호출만 수행, 결과 미사용        |
| `.then()` 사용  | 반환값 후처리 가능            |
| `.catch()` 없음 | 오류 발생 시 스크립트 즉시 중단    |
| `.catch()` 사용 | 오류를 처리하고 이후 DSL 계속 실행 |


**3.2 단순 호출 (결과 처리 없음)**   
```bash
tokens.CERC20("0x1234...").balanceOf("0x9999...");
```
- 호출은 실행됨
- 반환값은 저장되지 않음
- 오류 발생 시 스크립트 중단

**3.3 결과 처리 (.then)**
```bash
tokens.CERC20("0x1234...").balanceOf("0x9999...")
  .then(res => {});
```
- res는 Contract Call의 반환값
- 반환값은 자동으로 memory에 저장
- 필요 시 varStore를 통해 명시적으로 저장

**3.4 .then() 블록 규칙**   
- JavaScript arrow function 형식만 허용
- 멀티라인 코드 작성 가능
- 마지막 줄은 반드시 `;`
- DSL 명령어 사용 불가
  - set, echo, compose, deploy 등 ❌

```bash
.then(res => {
  // JS-like logic only
});
```

**3.5 DSL Memory로 결과 전달 (varStore)**
.then() / .catch() 내부에는
DSL memory의 proxy 객체인 varStore가 제공됩니다.
- key-value 형태 저장
- 저장 후 DSL에서 ${key} 로 접근 가능

***값 저장 예시***
```bash
tokens.CERC20("0x1234...").balanceOf("0x9999...")
  .then(res => {
    varStore["savedBalance"] = res;
  });

echo "${savedBalance}";
```

| 구분     | 자동 저장          | varStore            |
| ------ | -------------- | ------------------- |
| 저장 시점  | `.then()` 호출 시점 | 코드 실행 중             |
| 용도     | 단순 결과 사용       | 가공 / 테스트 결과 / 상태 누적 |
| 필수 여부  | 자동             | 선택                  |


**3.6 오류 처리 (.catch)**
.catch() 는 revert, RPC 에러, 실행 실패를 처리합니다.    
.catch() 는 .then() 와 다르게 결과가 자동으로 memory에 저장 되지 않습니다.

```bash
tokens.CERC20("0x1234...").transfer("0x0", "1000")
  .catch(err => {
    varStore["error"] = err;
  });

echo "${error}";
```
- .catch() 가 없으면 오류 발생 시 스크립트 중단
- .catch() 가 있으면:
  - 오류를 memory에 저장
  - 이후 DSL 명령어 계속 실행

***3.7 .then() + .catch() 조합***   

```bash
tokens.CERC20("0x1234...").transfer("${USER}", "1000")
  .then(res => {
    varStore["txResult"] = "success";
  })
  .catch(err => {
    varStore["txResult"] = "failed";
  });
```

**3.8 Contract Call과 일반 DSL 명령어의 차이**
| 구분           | Contract Call                                | 일반 DSL 명령어     |
| ------------ | -------------------------------------------- | -------------- |
| 실행 방식        | RPC 호출 (eth_call / tx)                       | CLI 내부 실행      |
| 결과 처리 방식     | `.then(res => {})`                           | `=> key`       |
| 오류 처리        | `.catch(err => {})`                          | 즉시 중단          |
| memory 자동 저장 | ✅ (`.then()` 사용 시)                           | ❌              |
| memory 저장 방법 | `.then()` 변수명 자동 저장<br/>+ `varStore` (추가 저장) | `=> key`       |
| 목적           | 컨트랙트 실행 / 상태 조회 / 검증                         | 데이터 처리 / 흐름 제어 |

---
### 4 Literal Type 규칙
- 모든 type은 반드시 string으로 "..." 감싸야 함
- Contract Call, Compose JSON 동일 적용

**✅ 올바른 예시**
```bash
tokens.CERC20("0x1234...")
  .transfer("0x9999...", "1000");
```

**❌ 잘못된 예시**
```bash
tokens.CERC20("0x1234...").transfer("0x9999...", 1000);
```

| 타입           | 표현 규칙        | 예제                    |
| ------------ | ------------ | --------------------- |
| string       | 반드시 ""       | "1000"                |
| array        | string array | ["1","2"]             |
| tuple        | array로 표현    | ["1","hello"]         |
| nested array | OK           | [["1","2"],["3","4"]] |

---
### 5 Contract Call — Array & Tuple
- Array와 Tuple은 반드시 문자열 원소 기반
- Boolean(true/false) 제외 모든 값 string 사용
- 중첩 배열(nested array)도 허용

**✅ 올바른 예시**
```bash
// uint[]
tokens.Sample("0x1234...").setList(["1", "2", "3"]);

// address[]
tokens.Sample("0x1234...").setAccounts([
  "0x1111...",
  "0x2222...",
  "0x3333..."
]);

// tuple(uint,string)
tokens.Sample("0x1234...").submitData(["123", "hello"]);

// tuple[] (배열 of 튜플)
tokens.Sample("0x1234...").batchSubmit([
    ["1", "alpha"],
    ["2", "beta"],
    ["3", "gamma"]
]);

// nested array
tokens.Sample("0x1234...").matrix([
    ["1", "2", "3"],
    ["4", "5", "6"]
]);
```

**❌ 잘못된 예시**
```bash
tokens.Sample("0x1234...").setList(["1", 2, "3"]);
```

---
### 6 Memory & Variable 시스템
- Key-Value 기반 전역 저장소(memory) 제공
- 변수 정의: set --var { ... }
- 변수 호출: ${key} 또는 ${key.subKey} 형태
- JSON Object 형태만 허용
- 저장 가능한 타입: string, bool, array, object, nested JSON

**단일 값 저장 예시**
```bash
# 저장
set --var { "noabiAddress": "0x1234..." }

# 사용
tokens.Sample("${noabiAddress}").symbol();
```

**배열 저장 예시**
```bash
# 저장
set --var { "amounts": ["1000000", "2000000"] }
set --var { "obj": { "list": ["10", "20"] } }

# 사용
tokens.Sample("${noabiAddress}").transfer("${USER}", "${amounts.0}");
echo "${obj.list.1}"  # → "20"
```

**중첩 오브젝트 저장 예시**
```bash
# 저장
set --var {
  "a": {
    "a1": "val1"
  }
}

# 사용
echo "${a.a1}"
```

---
### 7 DSL 명령어 결과 저장 규칙

**Contract Call 결과 처리**
- JS 스타일 .then(res => { ... }) 사용
- 결과는 then 블록에서만 사용 가능
- memory로 자동 저장
- 멀티라인 가능, 마지막 줄 반드시 `;`

**예시**
```bash
tokens.CERC20("0x1234...").balanceOf("0x9999...")
  .then(res => {});

echo "${res}"
```

**Contract Call Result → External Memory**   
.then() 블록 내부에서는 DSL memory에 값을 저장할 수 있습니다.   
이를 위해 `varStore` 객체가 제공됩니다.   
- `varStore`는 DSL memory의 proxy 입니다
- key-value 형태로 저장됩니다
- 저장된 값은 이후 DSL에서 `${key}` 형태로 사용 가능합니다

**예시**
```bash
tokens.CERC20("0x1234...").balanceOf("0x9999...")
  .then(res => {
    varStore["saved"] = "value";
  });

echo "${saved}"
```

**일반 DSL 명령어 결과 처리**
- curl, toUtc, getBalance 함수만 저장 가능
- => resultKey 패턴으로 memory 저장 가능

**예시**
```bash
curl GET http://localhost:8080/management/health/ping => res;

echo "${res}"
```

---
### 8 EIP-712 / EIP-191 서명 처리 예제

**Signature 요청 & 저장**
```bash
curl GET http://localhost:8080/management/health/ping => res;

echo "${res}"
```
# EIP712 Signatrue 요청 과 사용 방법 예시
```bash
# Back-End로 EIP712 Signature 요청
curl POST http://127.0.0.1:8080/v1/sign/eip712
  -H "Content-Type: application/json"
  -d '{
    "requestId": "req-1234",
    "domain": {
      "chainId": "1",
      "name": "MyProtocol",
      "version": "1",
      "verifyingContract": "0x0000000000000000000000000000000000000000"
    },
    "data": [
      { "type": "address", "name": "sender", "value": "0x1111111111111111111111111111111111111111" },
      { "type": "uint256", "name": "nonce", "value": "123" },
      { "type": "address[]", "name": "whitelist", "value": [
        "0xaaaAAAaaaAAAaaaAAAaaaAAAaaaAAAaaaAAAaaaA",
        "0xbbbBBBbbbBBBbbbBBBbbbBBBbbbBBBbbbBBBbbbB"
      ]}
    ],
    "KeyPair": [
      { "account": "0x82cffb99d7983695f85311de792c0925bc09b32a", "phrase": "!rbdiwkfwhagkwk@2025" }
    ]
  }' => res;

# Signature 출력
echo "signature: ${res.result.data.signatures.0}"

# Contract Call에 signature 전달
defi.noabi("0x1234....").claim("1000", "${res.result.data.signatures.0}");

```

---
**문서 네비게이션**  
[← 이전: Compose and Deploy Flow](composeNdeploy.md) | [다음: DSL Sample →](dslSample.md)  
[↑ 목차로 돌아가기](../README.md)