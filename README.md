# noABI REPL
Zero-ABI Interactive Smart Contract Execution

> noABI는 여러 스마트 컨트랙트와 체인이 실제 환경에서 어떻게 함께 동작하는지를   
엔드투엔드로 시뮬레이션하고 검증하기 위한 인터랙티브 CLI 도구입니다.   
>
> noABI는 유닛 테스트를 위한 도구가 아닙니다.   
단일 컨트랙트의 함수 검증이 아니라, 수십 개의 컨트랙트가 서로 연결되어   
실제 운영 환경에서 어떻게 동작하는지를 확인하는 데 목적이 있습니다.   
>
> 즉, noABI는 스마트 컨트랙트 코드를 테스트하는 도구가 아니라,   
프로토콜이 실제로 어떻게 ‘움직이는지’를 검증하는 도구입니다.   

---

![License](https://img.shields.io/badge/license-MIT-green)
![Node Version](https://img.shields.io/badge/node-16.14.2-blue)
---

## Features

- ***Interactive Contract Execution (REPL)***   
CLI에서 컨트랙트 함수를 즉시 호출하고, 조회(Call)와 트랜잭션(Tx)을 인터랙티브하게 실행

- ***End-to-End Protocol Simulation***   
단일 컨트랙트가 아닌 다수의 컨트랙트 연동 흐름을 실제 운영 순서대로 시뮬레이션

- ***Zero-ABI Execution***   
ABI 파일 없이 solc output 기반으로 컨트랙트 배포 및 호출 수행

- ***Script + Interactive Hybrid Automation***   
DSL 스크립트로 배포·호출을 자동화하면서, 필요 시 인터랙티브 CLI로 즉시 이어서 작업 가능

- ***Stateful Memory System***   
컨트랙트 호출 결과를 memory에 저장하고, 이후 시나리오 및 다른 컨트랙트 호출에 재사용

- ***Snapshot & Repeatable Scenarios***   
체인 상태를 스냅샷으로 고정하여 동일한 E2E 시나리오를 반복 실행 가능

- ***Multi-Contract & Multi-Chain Flows***   
여러 컨트랙트 및 체인 간 결과를 연결하여 브릿지·오라클 등 복잡한 흐름 검증

- ***EIP-712 / EIP-191 Signing Support***   
지갑 없이 백엔드 서명 처리를 자동화하여 서명 기반 프로토콜 흐름을 UI 없이 테스트

- ***Automation-First Design***   
배포, 컨트랙트 호출, 외부 API 검증까지 하나의 DSL 스크립트로 자동화

- ***Workspace Consistency***   
환경 설정, 빌드 아티팩트, 배포 로그를 단일 워크스페이스에서 일관되게 관리

---

## Installation

### Requirements
- Node.js 16.14.2 (LTS 권장)
- Ganache CLI (개발용 로컬 블록체인)

```bash
npm install -g ganache-cli
npm install
```

**Optional**  
VSCode Solidity Plugin   
Ethereum Solidity Language for Visual Studio Code by    Juan Blanco   

---

**Quick Start**
```bash
# 기본 설정 파일로 CLI 실행
./noabi

# 특정 설정 파일 사용
./noabi --config settings.json
```

실행 시 CLI 화면에서 도움말(help)을 확인할 수 있습니다.

---

**Documentation**

### 시작하기
- [Installation](doc/installation.md) - 설치 및 환경 구성
- [Getting Started](doc/gettingStarted.md) - 빠른 시작 가이드
- [Quick Example](doc/quickExam.md) - 빠른 예제 가이드

### 설정 및 구성
- [Configuration](doc/configuration.md) - 설정 파일 작성 방법
- [Compose File](doc/composeFile.md) - Compose 파일 설정 가이드

### 사용법
- [Commands](doc/commands.md) - 핵심 명령어 상세 설명
- [Domain-Specific Language (DSL)](doc/dsl.md) - DSL 문법 및 규칙
- [DSL Sample](doc/dslSample.md) - DSL 예제 코드

### 고급 기능
- [Compose and Deploy Flow](doc/composeNdeploy.md) - 배포 전체 흐름
- [Optional Modules](doc/optionalModules.md) - 선택적 모듈


--- 
**Contributing**  
PR, Issue 환영합니다.  
모든 참여자에게 감사드립니다. 🙏  

---

**License**  
MIT License


```