# Getting Started

빠르게 noABI CLI를 설치하고 실행해보는 단계입니다.  
설치 직후 기본 체인과 지갑 설정으로 CLI를 실행하며, 간단한 명령어를 테스트합니다.

---

## 1. 설치

```bash
# GitHub에서 최신 소스 다운로드
https://github.com/DaeSob/noAbi/archive/refs/heads/main.zip

# 또는
https://github.com/DaeSob/noAbi/archive/refs/heads/main.tar.gz

# 압축 해제

```

## 2. CLI 실행
**2.1 기본 설정 파일로 실행**
```bash
./noabi
```

**2.2 특정 설정 파일 지정**
```bash
./noabi --config settings.json
```

실행하면 다음과 같은 화면이 표시됩니다:

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
```

**3. 기본 명령어**

설치 직후 바로 테스트할 수 있는 간단한 명령어입니다.

```bash
# 사용 가능한 명령어 확인
> help

# 활성화된 address
> whoami 

```

**4. Quick Tip**
- help 명령어로 CLI에서 지원하는 모든 명령어를 확인할 수 있습니다.
- 빠른 테스트를 위해 로컬 노드(http://127.0.0.1:8545)를 실행하고 기본 지갑을 준비하는 것을 추천합니다.

---
**문서 네비게이션**  
[← 이전: Installation](installation.md) | [다음: Quick Example →](quickExam.md)  
[↑ 목차로 돌아가기](../README.md)