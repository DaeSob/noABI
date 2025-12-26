## Commands

CLI에서 사용 가능한 명령어와 옵션, 예시를 정리합니다.

**noABI DSL에서 사용 가능한 명령어**   
> memory, show 명령어는 인터랙티브 CLI 전용으로, DSL 스크립트에서는 사용할 수 없습니다.

| 명령어          | 설명                            |
| ------------ | ----------------------------- |
| `sh`         | DSL Script 실행            |
| `set`        | 변수 정의 및 출력 옵션 설정 등            |
| `use`        | 체인 또는 지갑 전환                   |
| `snapshot`   | 현재 상태 스냅샷 생성                  |
| `compose`    | Solidity 패키지 분석 및 배포 준비 파일 생성 |
| `deploy`     | Contract 배포                   |
| `import`     | ABI 또는 Contract import        |
| `export`     | ABI 또는 결과물 export             |
| `toUtc`      | 날짜/시간 변환                      |
| `getBalance` | 계정 잔액 조회                      |
| `curl`       | HTTP API 호출 (멀티라인 가능)         |
| `event`      | Event Log 조회 (멀티라인 가능)        |
| `echo`       | 변수 또는 메시지 출력                  |
| `sleep`      | 지정 시간만큼 대기                    |


<details>
<summary>sh 명령어</summary>   
스크립트 파일을 실행하는 명령입니다.
테스트용 DSL 스크립트, 배포 스크립트 등을 반복 실행하거나 자동화할 때 사용합니다.

**사용 형식**
```bash
sh <script_file> [--repeat <num>]
```
옵션 설명:
- --repeat <num>
  스크립트 반복 실행
  5 → 5회 반복
  1 → 기본값(한 번 실행)
  -1 → 무한 반복

**예시**
1. 스크립트 한 번 실행
    ```bash
    sh echo.txt
    ```
2. 스크립트 5회 반복 실행
    ```bash
    sh echo.txt --repeat 5
    ```
---
</details>

<details>
<summary>set 명령어</summary>   
전역 설정(global configuration)을 변경하거나 varStore 변수 값을 저장하는 명령입니다.    

테스트 스크립트, 배포 스크립트 등에서 환경을 통일할 때 사용됩니다.

**사용 형식**
```bash
set (--var <json> | --gasLimit <num> | --output <format>)
```
옵션 설명:
- --var {json}   
  varStore에 JSON 형태로 값을 저장합니다.
  저장된 값은 ${key} 형태로 DSL 전체에서 참조할 수 있습니다.

- --gasLimit value   
  기본 gas limit을 설정합니다.
  이후 모든 트랜잭션 호출에 기본적으로 적용됩니다.

- --output minimal|full   
  DSL 실행 시 출력 형식을 변경합니다.
  minimal : 깔끔한 요약 출력
  full : 상세한 full-trace 출력

**예시**
1. 변수 저장
  ```bash
  set --var { "name": "value" }
  ```   

2. 기본 gas limit 설정
  ```bash
  set --gasLimit 15000000
  ```   

3. 출력 포맷 minimal 변경
  ```bash
  set --output minimal
  ```   
---
</details>

<details>
<summary>memory 명령어</summary>
memory 명령은 스크립트 실행 중 생성된 메모리 공간에 저장된 값을 
조회하거나 초기화하는 기능을 제공합니다.   

데이터의 저장 자체는 set --var 또는 함수 호출의 res => {} 구문을 통해 수행되며,    
memory 명령은 저장된 데이터를 확인(read) 하거나 초기화(clear) 하는 역할만 담당합니다.    

**사용 형식**
```bash
memory [<key> | -clear]
```

**How Data is Stored**   
1. set --var 로 값 저장   
    ```bash
    set --var { "owner": "0x1111111111111111111111111111111111111111" }
    ```
    조회:   
    ```bash
    memory owner
    ```
2. 함수 호출 결과를 메모리에 저장
컨트랙트 함수 호출 시 .then(res=>{}) 사용:
    ```bash
    tokens.CERC20("${noAbitAddress}")
      .balanceOf("${USER}")
      .then(balance => {});
    ```
    조회:   
    ```bash
    memory balance
    ```

**When to Use**
- 배포 또는 함수 호출 후 값을 확인할 때
- 스크립트 흐름이 정상적으로 동작하는지 검증할 때
- 테스트 반복 실행을 위해 메모리를 초기화하고 싶을 때
- API 호출/컨트랙트 호출 결과를 이어지는 단계에서 활용할 때
---
</details>

<details>
<summary>show 명령어</summary>
show 명령어는 noABI CLI에서 저장된 컬렉션, 컨트랙트, 구성 정보를 조회하는 기능을 제공합니다.

**사용 형식**
```bash
show [<target>] [--currentSet] [--rpcUrl] [--gasLimit] [--json]
```

**예시**
```bash
show                                  # 모든 컬렉션 출력
show tokens                           # tokens 컬렉션 내 모든 컨트랙트 출력
show tokens.ERC20                     # ERC20 컨트랙트 정보 출력
show tokens.ERC20 | symbol            # ERC20 컨트랙트의 symbol 함수 결과 조회
```

**옵션**
| 옵션            | 설명                             |
| ------------- | ------------------------------ |
| `-currentSet` | 현재 CLI에서 선택된 체인/지갑 구성 조회       |
| `-rpcUrl`     | 현재 사용 중인 RPC URL 조회            |
| `-gasLimit`   | 기본 가스 한도(default gas limit) 조회 |
| `-json`       | 조회 결과를 JSON 형식으로 출력            |

**옵션 사용 예시**
```bash
show -currentSet                   # 현재 체인/지갑 설정 출력
show -rpcUrl                       # 현재 RPC URL 출력
show -gasLimit                     # 기본 가스 한도 출력
show tokens.ERC20 -json            # ERC20 컨트랙트 정보를 JSON 형식으로 출력
show tokens.ERC20 | symbol -json   # symbol 함수 결과를 JSON 형식으로 출력
```
---
</details>

<details>
<summary>import 명령어</summary>
import 명령어는 noABI CLI에서 ABI 파일, 개별 함수, 또는 저장된 메모리 데이터를 불러오는 기능을 제공합니다.

**사용 형식**
```bash
import (-abi [--add] --collection <name> --name <contract> (<abi_file | function_sig>)
```

**옵션 및 예시**

1. ABI 파일 불러오기   
컬렉션 단위로 ABI 파일을 가져와 컨트랙트 호출에 사용할 수 있습니다.
기존 import된 ABI는 deploy 시 초기화됩니다.

    ```bash
    import -abi --collection sample --name ERC20 ./artifacts/abi/sample.abi
    ```

2. Solidity 함수 수동 추가   
특정 함수만 수동으로 import할 때 사용합니다.   
함수 시그니처는 반드시 문자열로 감싸야 합니다.

    ```bash
    import -abi -add --collection sample --name ERC20 "function balanceOf(address owner) public view returns(uint256 amount)"
    ```

3. 저장된 메모리 불러오기
이전에 export 명령으로 저장한 store 데이터를 불러와 메모리 상태를 복원합니다.  
스크립트 기반 테스트나 반복 시나리오 실행 시 유용합니다.   

    ```bash
    import -memory ./artifacts/store/memory.txt
    ```
---
</details>

<details>
<summary>export 명령어</summary>
export 명령어는 현재 로드된 ABI 정보 또는 메모리 데이터를 파일로 저장하는 기능을 제공합니다.   
테스트 결과 저장, ABI 백업, 또는 반복 테스트 환경 구성에 활용할 수 있습니다.    

**사용 형식**
```bash
export (-memory <memory_file> | -abi --collection <collection> --name <contract> <abi_file>)
```

**옵션 및 예시**
1. Store(메모리) 내보내기   
현재 memory(storage)에 저장되어 있는 key-value 데이터를
지정한 파일로 저장합니다.

    ```bash
    export -memory ./artifacts/store/memory.txt
    ```
    - 스크립트 실행 후 메모리 상태 기록
    - 다음 테스트에서 import -memory로 복원 가능
    - 반복 테스트나 시나리오 재현에 유용

2. ABI 내보내기   
현재 로드되어 있는 특정 컬렉션/컨트랙트의 ABI를 파일로 저장합니다.   

    ```bash
    export -abi --collection tokens --name CJRC20 ./artifacts/abi/CJRC20.abi
    ```
    - deploy 또는 import 후 메모리에 적재된 ABI를 파일로 저장
    - 로컬 ABI 버전 관리 가능
    - 테스팅/검증/배포 파이프라인에서 활용 가능
---
</details>

<details>
<summary>compose 명령어</summary>
compose-file을 기반으로 패키지 단위 Build Artifact를 실제로 생성합니다.    
배포 정보, 초기 컬렉션 세팅, 반환값 매핑 등을 자동화할 수 있어 Script-Based Testing 및 대규모 배포 구성 시 핵심 역할을 합니다.

**사용 형식**
```bash
compose <compose_file>
```

**예시**
```bash
compose ./compose-file/compose_tokens.json
```
해당 명령을 실행하면 다음이 자동 처리됩니다:
 - compose JSON에 선언된 컬렉션(collection) 생성
 - 컨트랙트 이름 및 key-value 매핑 등록
  - deploy 후 반환되는 address 등을 store(memory)에 자동 저장
 - 이후 스크립트에서 ${변수명} 형태로 즉시 사용 가능

**Compose 파일의 역할**
 - 여러 컨트랙트를 논리적으로 묶어 관리
  - 배포(deploy) 흐름에서 사용할 주소/심볼/초기값 자동 매핑
  - 스크립트 실행 중 사용할 변수(store) 자동 생성
  - ABI import 전에 컬렉션 구조 확보
  > 즉, compose는 워크스페이스의 기본 구조를 정의하는 설정 엔진 역할을 합니다   

**Script-Based Testing에서의 위치**   
compose → deploy → import → 테스트 흐름에서   
compose는 가장 첫 단계에서 워크스페이스 형태를 정의하는 부분이며,   
모든 후속 테스트가 해당 구조 위에서 진행됩니다.

---
</details>

<details>
<summary>deploy 명령어</summary>
deploy 명령어는 compose JSON 파일에 정의된 컨트랙트들을 순차적으로 배포(deploy) 하고,     
배포 결과(address, txHash 등)를 memory에 자동 저장하는 핵심 기능입니다.     
compose가 “구조 정의”라면, deploy는 “실제 컨트랙트 배포 단계”에 해당합니다.     

**사용 형식**
```bash
deploy <compose_file> [alias...]
```
- compose_file   
  - 배포할 컨트랙트 구성(JSON)

- alias  
  - 선택적: compose_file 내에 정의된 특정 alias 그룹만 선택적으로 deploy할 때 사용합니다.  
    alias를 지정하면 해당 alias에 매핑된 컨트랙트만 배포됩니다.

**예시**
```bash
deploy ./compose-file/compose_tokens.json
```
compose 파일에 정의된 순서대로 모든 컨트랙트를 배포하며,
배포 결과는 compose 내 returns 규칙에 기반해 memory에 자동 매핑됩니다.

**Deploy 동작 요약**
- compose 파일 로딩
- `files.manifast` 파일 로딩
- collection 및 contract 구조 확인
- 각 contract를 순차적으로 배포
- 배포 결과를 memory에 자동 저장
  - 예: ${noAbitAddress}, ${symbol}, ${owner}, ${txHash}, etc
- 이후 DSL script에서 배포된 값을 바로 활용 가능
---
</details>

<details>
<summary>snapshot 명령어</summary>   
블록체인 실행 환경(Ganache / Anvil)의 상태 스냅샷을 생성, 복구, 조회하는 명령입니다.    

스냅샷에는 블록 상태뿐 아니라 varStore(내부 변수 저장소) 도 함께 저장·복원됩니다.   


**사용 형식**
```bash
snapshot (--name <snapshot_name> | --revert <name_or_id> | -list)
```
옵션 설명:   
- --name snapshot_name   
현재 블록체인 상태를 지정된 이름으로 스냅샷 저장
- --revert name_or_id   
스냅샷 이름 또는 ID로 상태 복원
- -list   
생성된 전체 스냅샷 목록 조회

> varStore는 snapshot 생성/복원 시 자동으로 함께 저장 및 복원됩니다.

**예시**
1. 스냅샷 생성   
```bash
snapshot --name before-test
```

2. 스냅샷 복원 (이름 기준)
```bash
snapshot --revert before-test
```

3. 스냅샷 복원 (ID 기준) 
```bash
snapshot --revert 0x123

```

4. 스냅샷 목록 확인
```bash
snapshot -list

```
---
</details>

<details>
<summary>use 명령어</summary>
체인 설정 또는 지갑 설정을 전환하는 명령입니다.    

사전에 정의된 chainSet / walletSet 프리셋을 기반으로 즉시 네트워크·지갑 환경을 변경할 수 있습니다.  

**사용 형식**
```bash
use (--chainSet <chainKey> | --walletSet <walletKey> | --chainSet <chainKey> --walletSet <walletKey>)
```

옵션 설명:
- --chainSet chainSetKey   
  chain-config에 등록된 체인 프리셋으로 전환합니다.
- --walletSet walletSetKey   
  wallet-config에 등록된 지갑 프리셋으로 전환합니다.
두 옵션은 단독 또는 함께 사용할 수 있습니다.
> ...SetKey는 config 파일에서 설정한 key입니다

**예시**
1. 지갑 프리셋만 변경
```bash
use --walletSet wallet-testnet
```
2. 체인 + 지갑 프리셋 둘 다 변경
```bash
use --chainSet chain-testnet --walletSet wallet-testnet
```
---
</details>

<details>
<summary>curl 명령어</summary>
HTTP 요청을 수행하고, 응답 데이터를 변수로 저장할 수 있는 유틸리티입니다.   

테스트 스크립트 내에서 API 호출(EIP-712, 백엔드 검증 등)에 활용됩니다.

**사용 형식**
```bash
curl <GET|POST|PUT|DELETE|PATCH> <url> [-H <header>] [-d <data>] [=> <variable>]
```
옵션 설명:
- GET | POST | PUT | DELETE | PATCH
  HTTP 요청 메서드 지정
- url
  요청을 보낼 대상 URL

- -H header
  요청에 헤더 추가
  여러 개 사용 가능

- -d data
  요청 Body 데이터(JSON 문자열 등)
- => variable
  응답 결과(JSON)를 memory 변수로 저장


**예시**
1. GET 요청 후 응답 저장   
    ```bash
    curl GET https://api.example.com/data => result
    ```

2. JSON Body 를 포함한 POST 요청
    ```bash
    curl POST https://api.example.com/api -d '{ "key": "value" }' => response
    ```

3. 헤더 + Body 포함 POST 요청
    ```bash
    curl POST https://api.example.com/api
        -H "Authorization: Bearer token"
        -d '{
          "data":"test"
        }' => result
    ```

---
</details>

<details>
<summary>event 명령어</summary>
로컬 파일에서 Event Log를 조회하고 필터링하는 명령입니다.   
이전 트랜잭션의 이벤트 결과를 다음 트랜잭션 입력 및 조회 용으로 사용하기 위해 활용됩니다.

**사용 형식**
```bash
event -load <source> [--block <num> | --fromBlock <num> --toBlock <num>]
       [--contractName <name>] [--contractAddress <addr>] [--eventName <name>]
       [--txHash <hash>] [--logIndex <num>] [=> <variable>]
```
옵션 설명:
- -load
  Event Log 조회 활성화 플래그 (필수)
- source
  Event Log 파일이 저장된 디렉토리 경로 (필수)
- --block num
  단일 블록 번호 지정
- --fromBlock num / --toBlock num
  블록 범위 지정 (둘 다 필요)
- --contractName name
  컨트랙트 이름으로 필터링
- --contractAddress addr
  컨트랙트 주소로 필터링
- --eventName name
  이벤트 이름으로 필터링
- --txHash hash
  트랜잭션 해시로 필터링
- --logIndex num
  로그 인덱스로 필터링
- => variable
  조회 결과를 memory 변수로 저장 (항상 배열)

**주의사항**
- `--block`과 `(--fromBlock/--toBlock)`은 동시에 사용할 수 없습니다.
- 모든 필터 옵션은 AND 조건으로 적용됩니다.
- 옵션 값에는 `${variable}` 형식의 변수를 사용할 수 있습니다.
- 결과는 항상 배열 형식으로 저장됩니다 (이벤트가 없으면 빈 배열 `[]`).

**예시**
1. 단일 블록의 이벤트 조회
    ```bash
    event -load C:/workspace/log/event --block 181816793 => events
    ```

2. 멀티라인으로 여러 필터 적용
    ```bash
    event -load C:/workspace/log/event
      --block 181816793
      --contractName DAO
      --eventName eventApproval
      => events;
    ```

3. 블록 범위로 조회
    ```bash
    event -load C:/workspace/log/event --fromBlock 181816790 --toBlock 181816793 => events
    ```

4. 변수 사용
    ```bash
    set --var {"blockNumber": "181816793", "eventPath": "C:/workspace/log/event"}
    event -load ${eventPath} --block ${blockNumber} => events
    ```

5. 조회 결과 활용
    ```bash
    event -load C:/workspace/log/event --block 181816793 => events
    echo "Found ${events.length} events"
    echo "First event docId: ${events[0].decoded.docId}"
    ```

**Event Summary 구조**
조회된 이벤트는 다음 구조로 저장됩니다:
```json
{
  "blockNumber": "181816793",
  "blockHash": "...",
  "transactionHash": "...",
  "transactionIndex": "2",
  "logIndex": "3",
  "contractName": "DAO",
  "contractAddress": "0xd0a5...",
  "eventName": "eventApproval",
  "decoded": {
    "approvor": "0xfc66...",
    "docId": "4",
    "invokable": true
  }
}
```

---
</details>

---
**문서 네비게이션**  
[← 이전: DSL Sample](dslSample.md) | [다음: Plugins →](plugins.md)  
[↑ 목차로 돌아가기](../README.md)