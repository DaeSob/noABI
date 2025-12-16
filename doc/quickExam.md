# Quick Example

noABI CLI를 처음 사용하는 사용자를 위한 빠른 시작 예제입니다.  
설치 후 기본 예제를 실행하여 전체 워크플로우를 한 번에 경험할 수 있습니다.

---

## 1. 프로젝트 구조 확인

예제 파일은 `examples/` 디렉토리에 포함되어 있습니다.

```bash
noAbi/
├─ bin/
│  └─ examples/
│     └─ sample-1-tokens/
│        └─ contracts/
│           └─ build/
│              └─ package-token-solc-output.json   ← Compile 결과
├─ examples/
│  ├─ config/
│  │   └─ settings.json                             ← noABI config file
│  └─ sample-1-tokens/
│     ├─ compose-files/
│     │   └─ compose_tokens.json                    ← Compose 파일
│     ├─ contracts/                                 ← Normalize Path
│     │   ├─ build/
│     │   │  └─ package-token.sol                  ← Solidity Package Entry Point
│     │   └─ src/
│     │      ├─ erc20.sol
│     │      └─ erc721.sol
│     └─ dsl/
|        ├─ deploy.txt                              ← DSL 스크립트
│        └─ tokens.txt                              ← DSL 스크립트
├─ node_modules/
├─ conf.d/
│  └─ settings.json                                 ← 기본 설정 파일
└─ noabi.bat
```

---

## 2. 로컬 블록체인 실행

로컬 개발 환경을 위해 Ganache CLI 또는 Anvil을 실행합니다.

**Ganache CLI 실행 예시:**
```bash
ganache-cli -d -m xen -l 15000000 -h 0.0.0.0 -p 8545 --chainId 1338 -e 10000000 -g 15932629
```

**Anvil 실행 예시:**
```bash
anvil --host 0.0.0.0 --port 8545 --chain-id 1337 --gas-limit 15000000
```

> **참고**: `-d` 옵션은 고정된 시드(mnemonic)를 사용하여 동일한 계정 주소를 생성합니다.  
> 테스트 반복 실행 시 유용합니다.

---

## 3. noABI CLI 실행

**3.1 기본 설정 파일로 실행**

설정 파일을 `conf.d/settings.json`에 준비한 경우:
```bash
./noabi
```

**3.2 예제 설정 파일로 실행**

예제 설정 파일을 직접 지정하여 실행:
```bash
./noabi --config examples/config/settings.json
```

**실행 화면 예시:**
```bash
################################################################################
Welcome to noABI REPL
https://github.com/DaeSob/noAbi
################################################################################

Usage:
   noAbi [--config <file>]

Current Environment:
   Active RPC   : http://127.0.0.1:8545
   Active Wallet: 0x1234....

type 'help' to see available commands.

wallet-xen-test-0> 
```

---

## 4. Solidity 컨트랙트 컴파일

VSCode에서 Solidity 파일을 컴파일하여 solc-output.json을 생성합니다.

**컴파일 단계:**

1. **VSCode에서 프로젝트 열기**
   - noABI 전체 프로젝트를 VSCode로 엽니다.

2. **Entry Point 파일 열기**
   - `examples/sample-1-tokens/contracts/build/package-token.sol` 파일을 엽니다.

3. **컴파일 실행**
   - 파일을 마우스 우클릭
   - `Solidity: Compile Contract` 선택
   - 또는 `Ctrl+Shift+P` → `Solidity: Compile Contract` 입력

4. **컴파일 결과 확인**
   - 컴파일 시간: 몇 초 ~ 수십 초 (컨트랙트 복잡도에 따라 다름)
   - 결과 파일 생성 확인:
     ```bash
     ./bin/examples/sample-1-tokens/contracts/build/package-token-solc-output.json
     ```

> **참고**: VSCode Solidity Plugin이 설치되어 있어야 합니다.  
> 설치 방법은 [Installation](installation.md) 문서를 참고하세요.

---

## 5. DSL 스크립트 실행

컴파일이 완료되면 DSL 스크립트를 실행하여 컨트랙트를 배포하고 테스트할 수 있습니다.

**DSL 스크립트 실행:**
```bash
wallet-xen-test-0> sh ./examples/sample-1-tokens/dsl/deploy.txt
```

**실행 흐름:**
1. `compose` - Compose 파일을 기반으로 빌드 아티팩트 생성
2. `deploy` - 컨트랙트 배포 및 주소를 memory에 저장
3. Contract Call - 배포된 컨트랙트 함수 호출 및 검증
4. 결과 출력 - 테스트 결과 및 상태 확인

**예상 출력:**
- 배포된 컨트랙트 주소
- 함수 호출 결과 (totalSupply, name, symbol)
- 트랜잭션 해시 및 상태
- 에러 처리 결과 (있는 경우)

---

## 6. 다음 단계

이 예제를 통해 다음을 확인할 수 있습니다:
- ✅ Compose 파일 기반 컨트랙트 배포
- ✅ 배포된 컨트랙트 함수 호출 (view/pure 함수)
- ✅ 트랜잭션 전송 및 결과 확인
- ✅ `.then()` / `.catch()` 블록을 통한 결과 처리
- ✅ Memory 변수 시스템 활용

더 자세한 내용은 다음 문서를 참고하세요:
- [Configuration](configuration.md) - 설정 파일 상세 설명
- [Compose File](composeFile.md) - Compose 파일 작성 방법
- [DSL](dsl.md) - DSL 문법 및 규칙
- [Commands](commands.md) - 사용 가능한 명령어 목록

---

**문서 네비게이션**  
[← 이전: Getting Started](gettingStarted.md) | [다음: Configuration →](configuration.md)  
[↑ 목차로 돌아가기](../README.md)
