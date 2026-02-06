/* eslint-disable no-undef */

import fetch from "node-fetch";
import Joi from "joi";

const apiServerHostName = process.env.apiserverhostname || "localhost";
const basePath = `http://${apiServerHostName}:9041/api`;

const callPostJsonApi = async (endPoint, postData, authToken = null) => {
  const url = basePath + endPoint;
  // console.log("POST " + url);

  const body = JSON.stringify(postData);

  let headers = { "Content-Type": "application/json" };
  if (authToken) {
    headers["Authorization"] = `Bearer ${authToken}`;
  }

  const response = await fetch(url, {
    method: "post",
    body,
    headers,
  });

  return response;
};

const callHappyPostJsonApi = async (expectedStatus, endPoint, postData) => {
  const response = await callPostJsonApi(endPoint, postData);

  if (response.status !== expectedStatus) {
    console.log("UNSUCCESSFUL Endpoint: ", endPoint);
    console.log("UNSUCCESSFUL Request: ", postData);
    console.log("UNSUCCESSFUL Status: ", response.status);
    console.log("UNSUCCESSFUL Response: ", await response.text());
    throw new Error(`Expected ${expectedStatus} got ${response.status}`);
  }

  return await response.json();
};

const callHappyPostJsonApiWithAuth = async (
  expectedStatus,
  authToken,
  endPoint,
  postData
) => {
  const response = await callPostJsonApi(endPoint, postData, authToken);

  if (response.status !== expectedStatus) {
    console.log("UNSUCCESSFUL Endpoint: ", endPoint);
    console.log("UNSUCCESSFUL authToken: ", authToken);
    console.log("UNSUCCESSFUL Request: ", postData);
    console.log("UNSUCCESSFUL Status: ", response.status);
    console.log("UNSUCCESSFUL Response: ", await response.text());
    throw new Error(`Expected ${expectedStatus} got ${response.status}`);
  }

  return await response.json();
};

const validateSchema = async (data, schema) => {
  return await schema.validateAsync(data);
};

const validateObject = async (data, objectKeysMap) => {
  return await validateSchema(
    data,
    Joi.object().keys(objectKeysMap).required()
  );
};

const callRawPostApi = async (
  endPoint,
  authToken,
  body,
  additionalHeaders = {}
) => {
  const url = basePath + endPoint;
  let headers = { "Content-Type": "text/plain" };
  Object.assign(headers, additionalHeaders);
  if (authToken) {
    headers["Authorization"] = `Bearer ${authToken}`;
  }
  const response = await fetch(url, {
    method: "post",
    body,
    headers,
  });
  return response;
};

export {
  callHappyPostJsonApiWithAuth,
  callPostJsonApi,
  callHappyPostJsonApi,
  validateObject,
  validateSchema,
  callRawPostApi,
};
