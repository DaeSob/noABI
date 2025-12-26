
const { textColor, printBlue, printWarning, printGray, printDefault } = require("../log/printScreen.js");



async function _commandHelp(_inputTokens) {

  printBlue("  set");
  printDefault("    Set global configuration values.");
  printDefault("    set (--var <json> | --gasLimit <num> | --output <format>)");
  // --var
  printDefault("    --var <json>               Set variable JSON");
  printGray("                                    exam) set --var {\"name\":\"value\"}");
  // --gasLimit
  printDefault("    --gasLimit <gasLimit>       Set default gas limit");
  printGray("                                    exam) set --gasLimit 15000000");
  // --output
  printDefault("    --output <minimal|full>     Change output format");
  printGray("                                    exam) set --output minimal");
  printGray("                                    exam) set --output full");
  console.log("");

  // show | ls 명령어 도움말 출력
  printBlue("  show | ls");
  printDefault("    Display stored collections, contracts, and configurations.");
  printDefault("    show [<target>] [--currentSet] [--rpcUrl] [--gasLimit] [--json]");
  printGray("                                     exam) show tokens");
  printGray("                                     exam) show tokens.ERC20");
  printGray("                                     exam) show tokens.ERC20 | symbol");
  // --currentSet
  printDefault("    --currentSet        Show current config set.");
  printGray("                                     exam) show --currentSet");
  // --rpcUrl
  printDefault("    --rpcUrl            Show current RPC URL.");
  printGray("                                     exam) show --rpcUrl");
  // --gasLimit
  printDefault("    --gasLimit          Show default gas limit.");
  printGray("                                     exam) show --gasLimit");
  // --json
  printDefault("    --json              Output in JSON format.");
  printGray("                                     exam) show tokens.ERC20 --json");
  printGray("                                     exam) show tokens.ERC20 | symbol --json");
  console.log("");

  printBlue("  whoami");
  printDefault("    Display wallet address.");
  console.log("");

  printBlue("  import");
  printDefault("    Import ABI files, functions, or stored data.");
  printDefault("    import (-abi [--add] --collection <name> --name <contract> (<abi_file | function_sig>)");
  printDefault("            | -memory <memory_file>)");
  // -abi (ABI 파일 가져오기)
  printDefault("    -abi                                 Import ABI file or add functions manually.");
  printDefault("        --collection <name>              (Required) Collection name.");
  printDefault("        --name <contract>                (Required) Contract name.");
  printDefault("        <abi_file>                       ABI file path (used when --add is not provided).");
  printGray("                                             exam)import -abi --collection sample --name ERC20 ./artifacts/abi/sample.abi");
  // -abi --add (함수 수동 추가)
  printDefault("    -abi --add                           Add a Solidity function manually.");
  printDefault("        --collection <name>              (Required) Collection name.");
  printDefault("        --name <contract>                (Required) Contract name.");
  printDefault("        <function_sig>                   Solidity function signature (used when --add is provided).");
  printGray("                                             exam)import -abi --add --collection sample --name ERC20 \"function balanceOf(address owner) public view returns(uint256 amount)\"");
  // -memory
  printDefault("    -memory <memory_file>                Load stored memory data.");
  printGray("                                             exam)import -memory ./artifacts/store/memory.txt");
  console.log("");

  printBlue("  export");
  printDefault("    Export ABI or memory file to local filesystem.");
  printDefault("    export (-memory <memory_file> | -abi --collection <collection> --name <contract> <abi_file>)");
  // 상세 설명
  printDefault("        -memory <memory_file>                    Export stored memory to file.");
  printGray("                                     exam)export -memory ./artifacts/store/memory.txt");
  printDefault("        -abi --collection <collection> --name <contract> <abi_file>");
  printDefault("                                                 Export ABI of a specific contract.");
  printGray("                                     exam)export -abi --collection tokens --name CJRC20 ./artifacts/abi/CJRC20.abi");
  console.log("");


  printBlue("  sh | run");
  printDefault("    Execute script file with an optional repeat loop.");
  // 명령어 시그니처
  printDefault("    sh <script_file> [--repeat <num>]");
  printDefault("    run <script_file> [--repeat <num>]");
  // 상세 설명
  printDefault("        <script_file>          Script file to execute.");
  printDefault("        [--repeat <num>]       Repeat execution. Use -1 for infinite loop.");
  printGray("                                     exam) sh echo.txt");
  printGray("                                     exam) sh echo.txt --repeat 5");
  printGray("                                     exam) sh echo.txt --repeat -1");
  console.log("");

  printBlue("  compose");
  printDefault("    Load and apply compose JSON for contract grouping.");
  printDefault("    compose <compose_file>");
  // 상세 설명
  printDefault("        <compose_file>         Path to compose JSON file.");
  printGray("                                     exam)compose ./compose-file/compose_tokens.json");
  console.log("");

  printBlue("  deploy");
  printDefault("    Deploy contracts using a compose file, with optional alias.");
  printDefault("    deploy <compose_file> [alias...]");
  // 상세 설명
  printDefault("        <compose_file>         Path to compose JSON file.");
  printDefault("        [alias...]             Optional alias for deployed contracts. Supports multiple values.");
  printGray("                                     exam)deploy ./compose-file/compose_tokens.json");
  printGray("                                     exam)deploy ./compose-file/compose.json CObjects CMine");
  printGray("                                     default) deploy ./compose-file/compose.json");
  console.log("");

  printBlue("  memory | mem");
  printDefault("    Show, search or clear stored memory.");
  printDefault("    memory [<key> | -clear]");
  // 상세 설명
  printDefault("        <key>                Search stored memory by key.");
  printDefault("        -clear               Clear all saved memory data.");
  printGray("                                     exam)memory");
  printGray("                                     exam)memory sp1");
  printGray("                                     exam)memory -clear");
  console.log("");

  printBlue("  use");
  printDefault("    Switch to another chain or wallet preset.");
  printDefault("    use (--chainSet <chainKey> | --walletSet <walletKey> | --chainSet <chainKey> --walletSet <walletKey>)");
  // 상세 설명
  printDefault("        --chainSet <chainKey>        Switch to the specified chain preset.");
  printDefault("        --walletSet <walletKey>      Switch to the specified wallet preset.");
  printDefault("        Note: At least one option is required.");
  printGray("                                     exam)use --chainSet chain-testnet");
  printGray("                                     exam)use --chainSet chain-testnet --walletSet wallet-testnet");
  console.log("");

  printBlue("  getBalance");
  printDefault("    Fetch native coin balance of an address.");
  printDefault("    getBalance <address>");
  printGray("                                     exam)getBalance 0x36c3f207d4ecc60bf6170b7c3c51f90ec6ce8145");
  printGray("                                     exam)getBalance ${USER}");
  console.log("");

  printBlue("  estimateGas");
  printDefault("    Estimate gas usage for a transaction call.");
  printDefault("    estimateGas [--from <address>] --to <address> [--gasPrice <num>] --callData <hex>");
  // 상세 설명
  printDefault("        --from <address>         (optional) Sender address.");
  printDefault("        --to <address>           (required) Target contract or address.");
  printDefault("        --gasPrice <num>         (optional) Gas price override.");
  printDefault("        --callData <hex>         (required) Encoded call data.");
  printGray("                                     exam)estimateGas --from 0x... --to 0x... --callData 0xa9059cbb...");
  printGray("                                     exam)estimateGas --to 0x... --callData 0xa9059cbb...");
  console.log("");

  printBlue("  toUtc");
  printDefault("    Convert timestamp or datetime to UTC.");
  printDefault("    toUtc [<datetime|unixtime>] [-d|-h|-m] [--add <sec>] [-toString]");
  // 상세 설명
  printDefault("        <datetime|unixtime>      (optional) Input time. e.g. 2023-05-08T16:43:00 or 1683507600");
  printDefault("        -d                        (optional) Convert to days.");
  printDefault("        -h                        (optional) Convert to hours.");
  printDefault("        -m                        (optional) Convert to minutes.");
  printDefault("        --add <sec>               (optional) Add seconds to the result.");
  printDefault("        -toString                 (optional) Output formatted UTC string.");
  printGray("                                     exam)toUtc 2023-05-08T16:43:00 -d --add 3600");
  printGray("                                     exam)toUtc 1683507600");
  printGray("                                     exam)toUtc");
  console.log("");

  printBlue("  abi");
  printDefault("    Generate function signatures or ABI from solidity function.");
  printDefault("    abi (-funcSig <function_name> | -create <solidity_function>)");
  // 상세 설명
  printDefault("        -funcSig <function_name>          Generate function selector/signature.");
  printDefault("        -create <solidity_function>       Parse solidity function and generate ABI.");
  printGray("                                     exam)abi -funcSig tokenURI(uint256)");
  printGray("                                     exam)abi -create \"function balanceOf(address owner) public view returns(uint256)\"");
  console.log("");

  printBlue("  keystore");
  printDefault("    Create new keystore file or import from private key.");
  printDefault("    keystore -create [--privateKey <hex>] [--password <pass>]");
  // 상세 설명
  printDefault("        -create                     (required) Create new keystore file.");
  printDefault("        --privateKey <hex>          (optional) Import from existing private key.");
  printDefault("        --password <pass>           (optional) Password for encrypting keystore.");
  printGray("                                     exam)keystore -create");
  printGray("                                     exam)keystore -create --privateKey 0xf208370... --password 1234");
  console.log("");

  printBlue("  curl");
  printDefault("    HTTP request utility with variable assignment.");
  printDefault("    curl <GET|POST|PUT|DELETE|PATCH> <url> [-H <header>] [-d <data>] [=> <variable>]");
  // 상세 설명
  printDefault("        <METHOD>                     (required) HTTP method.");
  printDefault("        <url>                        (required) Target request URL.");
  printDefault("        -H <header>                  (optional) Add request header.");
  printDefault("        -d <data>                    (optional) Send request body.");
  printDefault("        => <variable>                (optional) Save response into variable.");
  printGray("                                     exam)curl GET https://api.example.com/data => result");
  printGray("                                     exam)curl POST https://api.example.com/api -d '{\"key\":\"value\"}' => response");
  printGray("                                     exam)curl POST https://api.example.com/api -H \"Authorization: Bearer token\" -d '{\"data\":\"test\"}' => result");
  printGray("                                     exam)curl GET ${apiUrl} => response");
  console.log("");

  printBlue("  event");
  printDefault("    Query event logs from local files.");
  printDefault("    event -load <source> [--block <num> | --fromBlock <num> --toBlock <num>]");
  printDefault("         [--contractName <name>] [--contractAddress <addr>] [--eventName <name>]");
  printDefault("         [--txHash <hash>] [--logIndex <num>] [=> <variable>]");
  // 상세 설명
  printDefault("        -load                         (required) Flag to enable event loading.");
  printDefault("        <source>                      (required) Event source directory path.");
  printDefault("        --block <num>                 (optional) Single block number.");
  printDefault("        --fromBlock <num>             (optional) Start block number (use with --toBlock).");
  printDefault("        --toBlock <num>               (optional) End block number (use with --fromBlock).");
  printDefault("        --contractName <name>         (optional) Filter by contract name.");
  printDefault("        --contractAddress <addr>      (optional) Filter by contract address.");
  printDefault("        --eventName <name>            (optional) Filter by event name.");
  printDefault("        --txHash <hash>               (optional) Filter by transaction hash.");
  printDefault("        --logIndex <num>              (optional) Filter by log index.");
  printDefault("        => <variable>                 (optional) Save results into variable (always array).");
  printDefault("    Note: --block and (--fromBlock/--toBlock) cannot be used together.");
  printDefault("          All filter options use AND condition.");
  printDefault("          Option values support ${variable} syntax.");
  printGray("                                     exam)event -load C:/workspace/log/event --block 181816793 => events");
  printGray("                                     exam)event -load --block 181816793 C:/workspace/log/event => events");
  printGray("                                     exam)event -load C:/workspace/log/event --block ${blockNumber} => events");
  printGray("                                     exam)event -load C:/workspace/log/event --fromBlock 181816790 --toBlock 181816793 => events");
  printGray("                                     exam)event -load C:/workspace/log/event --block 181816793 --contractName DAO --eventName eventApproval => events");
  console.log("");

  printBlue("  snapshot");
  printDefault("    Create, revert, or list blockchain state snapshots (ganache-cli/anvil).");
  printDefault("    snapshot (--name <snapshot_name> | --revert <name_or_id> | -list)");
  // 상세 설명
  printDefault("        --name <snapshot_name>       Create a snapshot with the given name.");
  printGray("                                     exam)snapshot --name before-test");
  printDefault("        --revert <name_or_id>        Revert blockchain state to snapshot by name or ID.");
  printGray("                                     exam)snapshot --revert before-test");
  printGray("                                     exam)snapshot --revert 0x123");
  printDefault("        -list                        List all stored snapshots.");
  printGray("                                     exam)snapshot -list");
  printDefault("    Note: varStore is also saved/restored with snapshots.");
  console.log("");

  printBlue("  exit");
  printDefault("    Exit the CLI session.");
  console.log("");

}

module.exports._commandHelp = _commandHelp;
