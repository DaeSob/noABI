// 알려진 명령어 목록 (한 곳에서만 관리)
// 새로운 명령어를 추가할 때는 이 배열에만 추가하면 됩니다.
const knownCommands = ['show', 'ls', 'use', 'cd', 'import', 'export', 'set', 'deploy', 'compose', 
                       'memory', 'mem', 'help', 'version', 'whoami', 'sh', 'run', 'current', 
                       'toUtc', 'getBalance', 'estimateGas', 
                       'abi', 'keystore', 'curl', 'echo', 'snapshot', 'event'];

module.exports.knownCommands = knownCommands;

