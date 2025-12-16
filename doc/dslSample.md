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
TODO - 테스트 목적, 시나리오, 검증 범위 설명

#### Key Features
- TODO - 이 테스트 케이스에서 확인할 주요 기능 나열
- TODO - memory 사용, ABI import/export, deploy/contract call 등

#### DSL CODE
<DSL CODE>
---

---
**문서 네비게이션**  
[← 이전: DSL](dsl.md) | [다음: Commands →](commands.md)  
[↑ 목차로 돌아가기](../README.md)