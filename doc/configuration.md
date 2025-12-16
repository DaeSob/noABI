## Configuration

noABI CLI에서 프로젝트 및 개발 환경을 설정하는 방법을 정리합니다.
체인, 지갑, currentSet, 프로젝트, 워크스페이스, 옵션 등 각 항목별 필드 의미와 JSON 예시를 포함합니다.

---

### 1. Schema
settings.json 또는 config.json에 작성되는 전체 구조 예시입니다.
```json
{
  "chain-sample1": {
    "rpc": "https://arbitrum-sepolia.drpc.org",
    "custom": {
      "baseChain": "mainnet",
      "param": {
        "name": "arbitrum-sepolia",
        "networkId": 421614,
        "chainId": 421614
      },
      "hardFork": "petersburg"
    }
  },
  "chain-sample2": {
    "rpc": "http://127.0.0.1:8545",
    "blockCreationCycle": 1,
    "gasBufferPercent": 5
  },  
  "wallet-sample1": {
    "keyStore": "",
    "privKey": ""
  },  
  "currentSet": {
    "chain": "chain-sample2",
    "wallet": "wallet-sample1"
  },
  "workspace": {
    "root": "c:/workspace/noAbi",
    "output": {
      "build": "./artifacts/build"
    },
    "logs": {
      "deploy": "./logs/deploy"
    }
  },
  "options": {
    "multiline": true
  }
}
```

### 2. Chains 설정 방법
여러 블록체인 노드를 정의할 수 있으며, CLI 실행 시 currentSet.chain으로 선택된 체인이 사용됩니다.

| 필드                 | 설명                      | 선택  | 비고                            |
| ------------------ | ------------------------- | ------ | ----------------------------- |
| rpc                | JSON-RPC 엔드포인트 URL     | 필수   | -                            |
| blockCreationCycle | 블록 생성 주기               | 옵션   | 기본 1                      |
| gasBufferPercent   | 가스 계산 시 버퍼 (%)        | 옵션   | -                            |
| custom             | 특정 체인 옵션               | 옵션   | baseChain, param, hardFork 포함 |
| custom.baseChain   | 기준 체인 (예: mainnet)      | 옵션 | -                             |
| custom.param       | 체인 이름, networkId, chainId | 옵션 | -                             |
| custom.hardFork    | 사용할 hardfork              | 옵션 | -                             |

**예시**
```bash
{
  "chain-sample2": {
    "rpc": "http://127.0.0.1:8545",
    "blockCreationCycle": 1,
    "gasBufferPercent": 5
  },
  "chain-sample1": {
    "rpc": "https://arbitrum-sepolia.drpc.org",
    "custom": {
      "baseChain": "mainnet",
      "param": {
        "name": "arbitrum-sepolia",
        "networkId": 421614,
        "chainId": 421614
      },
      "hardFork": "petersburg"
    }
  }
}
```
> 새로운 체인을 추가하려면 chain-{이름} 키를 만들고 필요한 필드를 정의하세요.

### 3. Wallets 설정 방법
여러 지갑을 정의할 수 있으며, Private Key 또는 Keystore를 사용할 수 있습니다.

| 필드       | 설명                             | 비고                       |
| -------- | ------------------------------ | ------------------------ |
| privKey  | 개인키 직접 사용                      | 우선순위는 keystore보다 낮음      |
| keyStore | keystore 파일 경로 + passphrase 필요 | privKey 존재 시 keystore 우선 |

**예시**
```bash
{
  "wallet-sample2": {
    "privKey": "0x12345..."
  },
  "wallet-sample3": {
    "keyStore": "C:/workspace/daesob/cli/artifacts/keystores/12345",
    "description": "Keystore 사용 시 passphrase 입력 필요"
  }
}
```
> 새로운 지갑을 추가하려면 wallet-{이름} 키를 만들고, privKey 또는 keystore 정보를 입력하세요.

### 4. Current Set 설정
CLI 실행 시 기본으로 사용할 체인과 지갑을 지정합니다.

**예시**
```bash
{
  "currentSet": {
    "chain": "chain-sample2",
    "wallet": "wallet-sample3",
    "description": "CLI 기본 체인/지갑"
  }
}
```
> chain과 wallet 값은 위에서 정의한 키 중 하나를 지정하면 됩니다.

### 5. Workspace 설정
CLI가 작업할 공간을 정의하며, 파일 읽기/쓰기 시 상대 경로 기준이 됩니다.

| 필드                     | 설명                               |
| ---------------------- | -------------------------------- |
| workspace.root         | CLI 작업 기준 루트 경로. 지정 없으면 실행 위치 사용 |
| workspace.output.build | 빌드 결과물 저장 경로                     |
| workspace.logs.deploy  | 배포 과정 로그 저장 경로                   |

**예시**
```bash
{
  "workspace": {
    "root": "c:/workspace/noAbi",
    "output": {
      "build": "./artifacts/build"
    },
    "logs": {
      "deploy": "./logs/deploy"
    },
    "description": "CLI 산출물 관리 루트"
  }
}
```

### 6. Options 설정
CLI 실행 관련 옵션을 정의합니다.

| 필드        | 설명                                      |
| --------- | --------------------------------------- |
| multiline | true → contract function call에서 멀티라인 사용 |

**예시**
```bash
{
  "options": {
    "multiline": true
  }
}
```

---
**문서 네비게이션**  
[이전: Quick Example →](quickExample.md)   | [다음: Compose File →](composeFile.md)  
[↑ 목차로 돌아가기](../README.md)