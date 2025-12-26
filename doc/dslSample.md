## DSL Sample
아래는 noABI DSL을 사용한 Script-Based Testing 템플릿입니다.  
각 테스트 케이스별로 DSL CODE를 작성하고, memory, import/export, deploy 등의 흐름을 구성할 수 있습니다.

---
### Sample 1
#### Overview
이 샘플은 tokens 컬렉션에 포함된 CERC20 컨트랙트의 기본 동작과 실패 시나리오를 함께 검증하는 E2E 테스트입니다.

본 테스트는 다음 흐름을 검증합니다:

1. compose / deploy를 통한 컨트랙트 배포
2. 배포된 컨트랙트의 view 함수 호출
3. JavaScript 기반 .then() 로직을 이용한 결과 검증
4. revert 발생 시 .catch() 처리
5. DSL memory(varStore)를 활용한 상태 누적 및 출력

#### Key Features
- ***Compose & Deploy***   
    - Solidity 패키지 기반 컬렉션 초기화
    - 배포 결과(contract address) 자동 memory 저장

- ***Contract Read (eth_call)***   
    - totalSupply()
    - name()
    - symbol()
    - balanceOf(address)

- ***Contract Write (tx) + Revert 처리***   
    - 잘못된 주소(0x0)로 transfer 시 revert 발생
    - .catch()를 이용해 에러를 흡수하고 이후 로직 계속 수행

- ***JavaScript Assertion***
    - .then() 내부에서 Node.js assert 사용
    - 테스트 실패 시 에러 메시지를 memory에 저장

- ***Memory / Variable System***   
    - varStore를 통한 테스트 상태 누적
    - ${USER}, ${noAbitAddress} 등 전역 변수 활용

#### Test Flow Summary
```bash
compose
  ↓
deploy
  ↓
getBalance / toUtc
  ↓
totalSupply (view)
  ↓
name (view + assert)
  ↓
symbol (view)
  ↓
balanceOf (view → inc++)
  ↓
transfer (tx → revert → catch)
  ↓
balanceOf (view → inc++)
  ↓
echo (결과 출력)
```

#### DSL CODE
```bash
# ganache-cli -d -m xen -l 15000000 -h 0.0.0.0 -p 8545 --chainId 1338 -e 10000000 -g 15932629
# sh ./examples/sample-1-tokens/dsl/tokens.txt

set --output minimal

set --var { "inc": "0" }
set --var { "zeroAddress": "0x0000000000000000000000000000000000000000" }

###############################################################################################################
# 0. legacy tokens
# Collection Name: tokens
# deploy로 실행 되어 import 된 collection은 초기화 되어 abi가 import 됩니다
# 따라서 deploy 되기전 import 되어 있는 abi는 모두 사라지게 됩니다
compose ./examples/sample-1-tokens/compose-files/compose_tokens.json
deploy ./examples/sample-1-tokens/compose-files/compose_tokens.json

toUtc => now

getBalance ${USER} => ehtBalance

tokens.CERC20("${noAbitAddress}").totalSupply()
  .then(totalSupply => {});

tokens.CERC20("${noAbitAddress}").name()
  .then(name => {
    try {
      const assert = require("assert");
      assert.equal(name["0"], "noABI Token1", "Name should be 'noABI Token'");
      console.log("pass");
    } catch (e) {
      varStore["testRes"] = e.message;
      console.log(e.message);
    }
  });

tokens.CERC20("${noAbitAddress}").symbol()
  .then(symbol => {});

# 연산 예제
tokens.CERC20("${noAbitAddress}").balanceOf("${USER}")
  .then(balance => {
    varStore["inc"] = (Number(varStore["inc"]) + 1).toString();
  });

# revert 예제
tokens.CERC20("${noAbitAddress}").transfer("${zeroAddress}", "1000000")
  .then(balance => {})
  .catch(err => {
    varStore["catch"] = err;
  });

tokens.CERC20("${noAbitAddress}").balanceOf("${USER}")
  .then(balance => {
    varStore["inc"] = (Number(varStore["inc"]) + 1).toString();
  });

echo "----------------------------------------------------"
echo "now: ${now}"
echo "noAbi Address: ${noAbitAddress}"
echo "name: ${name}"
echo "symbol: ${symbol}"
echo "totol supply: ${totalSupply}"
echo "name test: ${testRes}"
echo "balance: ${balance}"
echo "inc: ${inc}"
echo "catch: ${catch}"
```

---

### Test Case 2
#### Overview
이 샘플은 EIP-712 서명 기반 프로토콜의 엔드투엔드 검증을 위한 테스트입니다.  
백엔드 API를 통해 EIP-712 서명을 생성하고, 스마트 컨트랙트에서 서명 검증을 수행하는 전체 흐름을 검증합니다.

본 테스트는 다음 흐름을 검증합니다:

1. compose / deploy를 통한 EIP-712 검증 컨트랙트 배포
2. curl 명령어를 통한 백엔드 API 호출 및 EIP-712 서명 요청
3. curl 응답에서 서명 추출 및 memory 저장
4. 배포된 컨트랙트의 recoverSigner 함수로 서명 검증
5. 외부 API와 스마트 컨트랙트 간 데이터 연동 검증

#### Key Features
- ***Compose & Deploy***   
    - EIP-712 검증 컨트랙트 배포
    - 배포 결과(contract address) 자동 memory 저장

- ***External API Integration (curl)***   
    - 백엔드 API를 통한 EIP-712 서명 생성
    - 멀티라인 JSON 데이터 전송
    - `${변수명}` 형식의 변수 치환 지원
    - 응답 데이터를 memory에 자동 저장

- ***EIP-712 Signature Verification***   
    - recoverSigner 함수를 통한 서명 검증
    - 서명자 주소 복구 및 검증
    - Permit 타입 서명 처리

- ***Memory / Variable System***   
    - curl 응답 데이터를 varStore에 저장
    - 중첩된 JSON 구조 접근 (예: `${res.result.data.signatures.0}`)
    - 배포된 컨트랙트 주소와 서명 데이터 연동

- ***End-to-End Protocol Testing***   
    - 백엔드 서명 생성 → 스마트 컨트랙트 검증의 전체 흐름
    - 실제 운영 환경과 유사한 시나리오 검증

#### Test Flow Summary
```bash
compose
  ↓
deploy (EIP-712 검증 컨트랙트)
  ↓
curl POST (백엔드 API로 EIP-712 서명 요청)
  ↓
서명 데이터 memory 저장 (${res.result.data.signatures.0})
  ↓
recoverSigner (컨트랙트에서 서명 검증)
  ↓
검증 결과 확인
```

#### DSL CODE
```bash
# ganache-cli -d -m xen -l 15000000 -h 0.0.0.0 -p 8545 --chainId 1338 -e 10000000 -g 15932629
# sh ./examples/sample-2-eip712/dsl/test.txt

set --output minimal

###############################################################################################################
# 0. eip712 test
# Collection Name: eip712
# deploy로 실행 되어 import 된 collection은 초기화 되어 abi가 import 됩니다
# 따라서 deploy 되기전 import 되어 있는 abi는 모두 사라지게 됩니다
compose ./examples/sample-2-eip712/compose-files/compose_eip712.json
deploy ./examples/sample-2-eip712/compose-files/compose_eip712.json

# 백엔드 API를 통한 EIP-712 서명 요청
curl POST http://127.0.0.1:8080/v1/sign/eip712
  -H "Content-Type: application/json"
  -d '{
        "requestId": "req-1234",
        "domain": {
          "chainId": "1338",
          "name": "NOABI",
          "version": "V0.6.0",
          "verifyingContract": "${testAddress}"
        },
        "primaryType": "Permit",
        "data": [
          {
            "type": "address",
            "name": "owner",
            "value": "0xb171fe0b0804651446a50344ae14e56596190bcf"
          },
          {
            "type": "address",
            "name": "spender",
            "value": "0x71ad669844d04f7eb028d6014473c4c099ad1a60"
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
      => res1;

# 서명 검증 (컨트랙트에서 서명자 주소 복구)
eip712.EIP712RecoverTest("${testAddress}").recoverSigner(
    "0xb171fe0b0804651446a50344ae14e56596190bcf",
    "0x71ad669844d04f7eb028d6014473c4c099ad1a60",
    "1000000",
    "${res1.result.data.signatures.0}"
  )
  .then(res2 => {
    // 서명 검증 결과 확인
    const expectedAddress = "0xb171fe0b0804651446a50344ae14e56596190bcf";
    const assert = require("assert");
    try {
      assert.equal(res2.recoveredAddress.toLowerCase(), expectedAddress.toLowerCase(), 
        "Recovered address should match signer");
      console.log("✓ Signature verification passed");
      varStore["verificationResult"] = "success";
    } catch (e) {
      varStore["verificationResult"] = "failed: " + e.message;
      console.log("✗ Signature verification failed: " + e.message);
    }
  })
  .catch(err => {
    varStore["verificationError"] = err.message;
    console.log("Error: " + err.message);
  });

echo "----------------------------------------------------"
echo "Test Address: ${testAddress}"
echo "Signature: ${res.result.data.signatures.0}"
echo "Verification Result: ${verificationResult}"
echo "Verification Error: ${verificationError}"
```

---
**문서 네비게이션**  
[← 이전: DSL](dsl.md) | [다음: Commands →](commands.md)  
[↑ 목차로 돌아가기](../README.md)