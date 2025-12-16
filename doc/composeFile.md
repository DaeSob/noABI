# Compose File
noABI CLI에서 Solidity 패키지를 빌드하고 배포하기 위한 Compose 파일 설정 문서입니다.
빌드, 배포, 결과물 관리에 필요한 모든 메타데이터를 정의합니다.

### 1. Schema
Compose 파일 기본 구조 예시:
```json
{
  "collection": "",
  "build": {
    "packages": [
      {
        "source": {
          "normalize": {
            "path": ""
          },          
          "entry": ""
        },
        "artifacts": {
          "solcOutput": ""
        },
        "contracts": [
          {
            "name": "",
            "alias": "",
            "ignore": false,
            "inputs": {
              "desc": [],
              "types": [],
              "value": []
            },
            "returns": {
              "varNames": [],
              "txReceipts": []
            }
          }
        ]
      }
    ]
  },
  "files": {
    "manifest": ""
  }
}

```

| 필드             | 설명                                |
| -------------- | --------------------------------- |
| collection     | 이 빌드/패키지 세트가 속한 컬렉션 이름            |
| build.packages | VSCode에서 빌드한 결과물 정보 및 배포할 컨트랙트 정의 |
| files.manifest | `compose` 명령 결과물. 배포 시 사용됨        |


### 2. Normalize Path 와 Solidity Package Entry Point 설정
VSCode Solidity 플러그인을 이용해 빌드할 때 Entry Point 파일을 지정합니다.   
이 파일을 기준으로 dependency graph가 생성되고, solc 컴파일이 시작됩니다.    

| 필드                            | 설명                           |
| ----------------------------- | ---------------------------- |
| build.packages[].source.normalize.path | solc-output.json에 기록된 source 경로를 정규화하기 위한 기준 경로       |
| build.packages[].source.entry | 빌드 시작점(Entry Point) 경로       |


***왜 Source Normalization이 필요한가***   
solc 컴파일 결과(solc-output.json)에는 다음과 같은 특징이 있습니다.
- sources 필드는 절대 경로 또는 컴파일 시점의 전체 경로를 key로 사용
- OS, 개발 환경, 워크스페이스 위치에 따라 경로가 달라짐
- 동일한 Solidity 코드라도 경로가 달라지면 다른 Contract로 인식될 수 있음

noABI는 다음을 목표로 합니다:
- 환경이 달라져도 동일한 Contract를 안정적으로 식별
- solc-output.json과 Compose 정의 간의 명확한 매칭
- CI / 로컬 / 운영 환경 간 경로 차이 제거   

이를 위해 Source Normalization 개념을 사용합니다.    

**프로젝트 디렉토리 예시**
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

**Entry Point 파일 예시 (package-token.sol)**
```solidity
// SPDX-License-Identifier: MIT
// XeN Contracts v0.6.0
pragma solidity ^0.8.13;

import '../src/erc20.sol';
import '../src/erc721.sol';
```

**설정 예시**
```json
"source": {
  "normalize": {
    "path": "c:/workspace/noAbi/examples/sample-1-tokens/contracts"
  },  
  "entry": "c:/workspace/noAbi/examples/sample-1-tokens/contracts/build/package-token.sol"
}
```

### 3. Solc Build Output 파일 설정
solc가 컴파일 후 생성한 Output 파일은 스마트 컨트랙트 배포 시 Bytecode와 ABI의 근거 자료입니다.

| 필드                                    | 설명                |
| ------------------------------------- | ----------------- |
| build.packages[].artifacts.solcOutput | solc 컴파일 결과 파일 경로 |


**예시**
```json
"artifacts": {
  "solcOutput": "c:/workspace/noAbi/bin/examples/sample-1-tokens/contracts/build/package-token-solc-output.json"
}
```

**Solc Build Output 파일 예시**
- package-token-solc-output.json
```json
{
    "contracts": {
        "c:/workspace/noAbi/examples/sample-1-tokens/contracts/src/erc20.sol": {
            "CERC20": {
                "abi": [...],
                ...
                "evm": {
                  "bytecode": {...}
                }
	          }       
        }
	      ...
    }
}	
```
- VSCode의 plugin으로 compile한 결과 파일 입니다

**Entry Point와 solc-output JSON 경로 정규화**   
Compose 파일에서 설정한 build.packages[].contracts[].name은 실제 배포 시 solc-output JSON에서 해당 Contract의 bytecode와 ABI를 찾아 사용해야 합니다.

하지만 solc-output JSON에는 컴파일 시점의 절대 경로가 key로 저장되어 있고, Entry Point에서는 상대 경로(import '../src/erc20.sol';) 형태로 Contract를 참조하기 때문에, 단순 비교만으로는 매칭이 어렵습니다.

이 문제를 해결하기 위해 normalize path를 기준으로 경로를 정규화합니다.
즉, Entry Point에서 import된 상대 경로와 solc-output JSON의 절대 경로를 normalize path를 기반으로 **일관된 키(key)**로 변환하여 매핑합니다.

그 결과 Compose 단계에서는, Contract 이름(build.packages[].contracts[].name)과 solc-output JSON의 Contract 정보를 정확히 연결할 수 있으며, 이를 통해 필요한 bytecode와 ABI를 찾아 안전하게 배포(deploy)할 수 있습니다.

### 4. Contract Metadata 설정
contracts 항목은 배포할 컨트랙트 메타데이터를 정의합니다.

| 필드                 | 설명                                                                    |
| ------------------ | --------------------------------------------------------------------- |
| name               | 배포할 컨트랙트 이름. solcOutput에서 bytecode/ABI 매칭                             |
| alias              | 동일한 contract name을 2번 이상 deploy시 구분하기 위한 별칭                        |
| ignore             | true일 경우 배포에서 제외 (기본: false)                                          |
| inputs             | 컨트랙트 생성자(Constructor) 파라미터                                            |
| inputs.desc        | 각 파라미터 설명                                                             |
| inputs.types       | Solidity 타입 목록                                                        |
| inputs.value       | 실제 배포 시 전달 값                                                          |
| returns            | 배포 후 eth_getTransactionReceipt에서 원하는 필드를 memory에 자동 저장                |
| returns.varNames   | 저장할 변수 이름                                                             |
| returns.txReceipts | 트랜잭션 Receipt에서 가져올 필드명 (예: contractAddress, gasUsed, transactionHash) |


**예시:inputs**
```json
"inputs": {
  "desc": [
    "contract name",
    "contract symbol",
    "decimal",
    "auto"
  ],
  "types": [
    "string",
    "string",
    "uint256",
    "bool"
  ],
  "value": [
    "noABI Token",
    "NOAT",
    "18",
    "false"
  ]
}
```
> inputs은 construct의 파라미터와 매핑 됩니다
> value는 모두 string type 이어야 합니다

**예시:Returns**
```json
"returns": {
  "varNames": [
    "noAbitAddress"
  ],
  "txReceipts": [
    "contractAddress"
  ]
}
```
> returns를 활용하면 이후 스크립트나 Compose에서 ${noAbitAddress} 형태로 값을 재사용할 수 있습니다.

---
**문서 네비게이션**  
[← 이전: Configuration](configuration.md) | [다음: Compose and Deploy Flow →](composeNdeploy.md)  
[↑ 목차로 돌아가기](../README.md)