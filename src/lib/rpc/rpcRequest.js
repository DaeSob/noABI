const { HttpRequest } = require("./httpRequest.js");

const rpcRequest = async (url, method, param, id) => {
  const data = {
    jsonrpc: "2.0",
    method: method,
    params: param,
    id: id,
  };
  return await HttpRequest("POST", url, data);
};

module.exports.RpcRequest = rpcRequest;
