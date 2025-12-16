## Installation

로컬 환경에서 스마트 컨트랙트를 빠르게 개발·테스트하고 싶다면 Ganache CLI를 활용하여 개인 이더리움 네트워크를 구성할 수 있습니다.  
이를 통해 개발 속도와 테스트 효율이 크게 향상됩니다.

---

### 1 Requirements
스마트 컨트랙트 개발 환경을 안정적으로 구성하려면 아래 버전을 권장합니다:

- Node.js 16.14.2 (LTS 권장)  
  해당 버전은 Ganache CLI 및 여러 개발 도구와의 호환성이 가장 뛰어납니다.

---

### 2 로컬 블록체인 설치

로컬 개발 환경을 구성하기 위해 다음 중 하나를 선택하여 설치할 수 있습니다.

#### 2.1 Ganache CLI 설치

```bash
npm install -g ganache-cli
```

설치가 완료되면 아래 명령어로 정상 설치 여부를 확인할 수 있습니다:
```bash
ganache-cli --version
```

**Ganache 실행 예시:**
```bash
# 기본 설정으로 실행
ganache-cli

# 커스텀 설정으로 실행
ganache-cli -d -m xen -l 15000000 -h 0.0.0.0 -p 8545 --chainId 1338 -e 10000000 -g 15932629
```

#### 2.2 Anvil 설치 (권장)

Anvil은 Foundry 프로젝트의 일부로 제공되는 로컬 블록체인 노드입니다.  
Ganache CLI보다 빠르고 최신 이더리움 기능을 지원합니다.

**설치 방법:**

1. **Foundry 설치 (Anvil 포함)**
   ```bash
   # Linux / macOS
   curl -L https://foundry.paradigm.xyz | bash
   foundryup
   
   # 또는 직접 설치
   cargo install --git https://github.com/foundry-rs/foundry foundry-cli anvil --locked
   ```

2. **Windows 설치**
   ```bash
   # PowerShell에서 실행
   irm https://aka.ms/foundry/install.ps1 | iex
   ```

3. **설치 확인**
   ```bash
   anvil --version
   ```

**Anvil 실행 예시:**
```bash
# 기본 설정으로 실행
anvil

# 커스텀 설정으로 실행
anvil --host 0.0.0.0 --port 8545 --chain-id 1337 --gas-limit 15000000
```

**Anvil vs Ganache CLI 비교:**
| 기능 | Anvil | Ganache CLI |
|------|-------|-------------|
| 성능 | 빠름 | 보통 |
| 최신 EVM 지원 | ✅ | 제한적 |
| 스냅샷 지원 | ✅ | ✅ |
| 설치 방법 | Foundry 필요 | npm으로 간단 |
| 활발한 개발 | ✅ | 유지보수 모드 |

> **참고**: noABI CLI는 Ganache CLI와 Anvil 모두 지원합니다. 프로젝트 요구사항에 따라 선택하세요.

### 3 프로젝트 패키지 설치
프로젝트 내 의존성을 설치하려면 다음 명령어를 실행합니다:
```bash
npm install
```
위 명령은 package.json에 정의된 모든 패키지를 자동으로 설치합니다.

### 4 추가 설치 (VSCode Solidity Plugin)
VSCode에서 직접 스마트 컨트랙트를 컴파일하고, 생성된 noABI 결과물을 활용하여 배포 및 테스트를 진행하고 싶다면 Solidity 전용 확장 기능을 설치해야 합니다.

**VSCode Plugin 설치 방법:**
1. VSCode를 실행합니다.
2. 아래 단축키로 Extensions 패널을 엽니다:
    ```bash
    Shift + Ctrl + X
    ```
3. Extensions 검색창에 solidity 입력

4. 검색 결과 목록에서 아래 확장 프로그램을 찾습니다:   
    Ethereum Solidity Language for Visual Studio Code by Juan Blanco   
    선택 후 Install 버튼을 클릭하여 설치합니다.

설치가 완료되면, VSCode에서 .sol 파일을 열었을 때 자동으로 컴파일 기능과 언어 지원 기능을 사용할 수 있습니다.

---
**문서 네비게이션**  
[← 이전: README](../README.md) | [다음: Getting Started →](gettingStarted.md)  
[↑ 목차로 돌아가기](../README.md)