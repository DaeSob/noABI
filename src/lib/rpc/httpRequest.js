const Axios = require("axios");

const _httpRequest = async (method, url, data) => {
  let axiosConfig = { method: method, url: url, data: data };
  let res = undefined;
  try {
    res = await Axios(axiosConfig);
    return res.data;
  } catch (e) {
    throw Error(e.message);
  } finally {
    axiosConfig = undefined;
    res = undefined;
  }
};

const _httpRequestEx = async (method, url, headers, data) => {
  let axiosConfig = {
    method: method,
    url: url,
    headers: headers,
    data: data,
  };
  let res = undefined;

  try {
    res = await Axios(axiosConfig);
    return res.data;
  } catch (e) {
    throw Error(e.message);
  } finally {
    axiosConfig = undefined;
    res = undefined;
  }
};

module.exports.HttpRequest = _httpRequest;
module.exports.HttpRequestEx = _httpRequestEx;
