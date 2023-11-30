/*
 * Copyright 2018 Google LLC.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
const data = require('./data/currency_conversion.json');

function _getCurrencyData (code) {
  return data[code];
}

function _carry (amount) {
  const fractionSize = Math.pow(10, 9);
  amount.nanos += (amount.units % 1) * fractionSize;
  amount.units = Math.floor(amount.units) + Math.floor(amount.nanos / fractionSize);
  amount.nanos = amount.nanos % fractionSize;
  return amount;
}

function convert (request) {

      // Convert: from_currency --> EUR
      const from_code = request.currency_code;
      const euros = _carry({
        units: request.units / _getCurrencyData(from_code),
        nanos: request.nanos / _getCurrencyData(from_code)
      });
      euros.nanos = Math.round(euros.nanos);
      //console.log(euros.nanos);
      // Convert: EUR --> to_currency
      const result = _carry({
        units: euros.units * data[request.to_code],
        nanos: euros.nanos * data[request.to_code]
      });

      result.units = Math.floor(result.units);
      result.nanos = Math.floor(result.nanos);
      result.currency_code = request.to_code;

      return result;

}

exports.handler = async function (event, context) {
  try {
    // Parse the input event, assuming it's a JSON object
    const requestBody = event ;
    // console.log(event);
    // console.log(event.currency_code);
    
    // Call the charge function with the parsed request
    const result = convert(requestBody);

    // Return a successful response
    const response = {
      statusCode: 200,
      body: JSON.stringify(result),
    };

    return response;
  } catch (error) {
    // Handle errors and return an appropriate error response
    const response = {
      statusCode: error.statusCode || 500, // You can customize the status code as needed
      body: JSON.stringify({ error: error.message }),
    };

    return response;
  }
};
