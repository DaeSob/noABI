## Compose & Deploy
본 섹션은 VSCode 환경에서 Solidity 코드를 작성하고, noABI CLI를 이용해 Compose 및 Deploy까지 진행되는 전체 흐름을 설명합니다.

---
### 1 개발 환경 구성
- VSCode 확장 설치: Solidity (by Juan Blanco)
- Solidity 파일 저장 시 자동으로 solc 컴파일 수행
- solc JSON Output(solc-output.json) 생성

---
### 2 VSCode Solidity Compile 흐름

**1. 프로젝트 디렉토리 예시**
```bash
noAbi/
├─ bin/
│  └─ examples/
|     └─ sample-1-tokens/
|       └─ contracts/
|         └─ build/
|           └─ package-token-solc-output.json   ← Compile 결과
├─ examples/
│  ├─ sample-1-tokens/
|  │   └─ contracts/                              ← Normalize Path
|  │      ├─ build/
│  │      |  └─ package-token.sol                ← Solidity Package Entry Point
│  │      └─ src
|  |         ├─ erc20.sol
|  |         └─ erc721.sol
|  └─ ...
├─ node_modules/
└─ ...
```

| 디렉토리                | 설명                        |
| ------------------- | ------------------------- |
| `sample-1-tokens/contracts/build/*`  | 패키지 단위 Entry Point 파일 위치  |
| `src/*`                             | 실제 Contract 구현 코드         |
| `node_modules`                      | 외부 Solidity 패키지 및 타입 정의   |
| `bin`                               | Compile 결과 JSON Output 저장 |


**2. Solidity Package Entry Point**
```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

import '../src/erc20.sol';
import '../src/erc721.sol';
```
- Entry Point: 패키지 내 모든 Contract를 단일 컴파일 대상으로 통합
- VSCode에서 저장 후 Solidity: Compile Contract 실행

**3. Compile 결과**   
package-token-solc-output.json 생성
- 포함 정보:
  - abi
  - bytecode / deployedBytecode
  - metadata, devdoc, userdoc
  - source path 리스트
  - dependency 구조

---
### 3 noABI Compose 파일 생성
- Compose 파일: 패키지 단위 배포 설정
- Entry Point와 solc-output 지정
```json
{
  "collection": "tokens",
  "build": {
    "packages": [
      {
        "source": {
          "normalize": {
            "path": "c:/workspace/noAbi/examples/sample-1-tokens/contracts"
          },  
          "entry": "c:/workspace/noAbi/examples/sample-1-tokens/contracts/build/package-token.sol"
        },
        "artifacts": {
           "solcOutput": "c:/workspace/noAbi/bin/examples/sample-1-tokens/contracts/build/package-token-solc-output.json"
        },
        "contracts": [
          {
            "name": "CERC20",
            "inputs": {
              "desc": ["contract name","contract symbol","decimal"],
              "types": ["string","string","uint256"],
              "value": ["noABI Token","NOAT","18"]
            },
            "returns": {
              "varNames": ["noabiAddress"],
              "txReceipts": ["contractAddress"]
            }
          }
        ]
      }
    ]
  },
  "files": {}
}
```
- Compose 파일 생성 후 Build Artifact 및 manifest.obj가 자동 생성
- files.manifest에 Compose 결과 등록

---
### 4 noABI Compose
- Solidity Package Entry Point와 solc-output JSON을 기반으로 패키지 단위 Contract 정보를 분석하고, 배포 가능한 Build Artifact를 생성.
- 핵심 입력:   
  - Solidity Package Entry Point (.sol)    
  - solc-output JSON (package-token-solc-output.json)
- 처리 과정:
  - 패키지 내 모든 Contract 분석
  - Constructor 인자 매핑
  - 내부 Dependency 정리
- 출력:   
  - manifest.obj 생성: 배포에 필요한 모든 정보 포함
  - Compose에서 지정한 files.manifest 경로에 Build Artifact 등록

**예시**
```bash
compose ./examples/sample-1-tokens/compose-files/compose_tokens.json
```

> 요약: Compose는 단순히 JSON 파일 생성이 아니라 배포 준비 단계의 핵심 파이프라인이며, Deploy 이전에 반드시 수행해야 하는 단계입니다.

### 5 noABI Deploy
- Deploy 입력: manifest.obj
- Compose가 계산한 배포 순서 및 constructor-map 사용
- 트랜잭션 전송 및 receipt 기록
- 배포 결과: contract address, 지정된 returns 항목에 매핑 -> memory에 저장


**예시**
```bash
deploy ./examples/sample-1-tokens/compose-files/compose_tokens.json
```

---
### 6 전체 흐름도
```bash
1. 개발자 코드 작성
   - src/* 아래 Contract 구현

2. Solidity Entry Point 생성
   - contracts/build/package-*.sol
   - 패키지 내 모든 Contract import

3. VSCode → Solidity Compile
   - Entry Point 저장 → solc 실행
   - solc-output.json 생성 (abi, bytecode, metadata)

4. noABI Compose 파일 구성
   - compose_noabi.json 작성
   - Entry Point + solc-output 지정
   - Contract constructor 정의

5. noABI Compose 실행
   - 입력: Entry Point, solc-output
   - 처리: Contract 분석, constructor 매핑, 배포 순서 계산
   - 출력: manifest.obj 생성

6. manifest.obj 생성 완료
   - Contract 의존성 그래프
   - constructor-map
   - bytecode / deploy bytecode
   - ABI 불필요

7. noABI Deploy 실행
   - 입력: manifest.obj
   - 처리: 배포 순서대로 트랜잭션 송신
   - 출력: contract address, returns 매핑

8. 배포 결과 저장 및 후처리
   - ContractAddress 반환
   - 다른 패키지에서 참조 가능
```

---
**문서 네비게이션**  
[← 이전: Compose File](composeFile.md) | [다음: DSL →](dsl.md)  
[↑ 목차로 돌아가기](../README.md)