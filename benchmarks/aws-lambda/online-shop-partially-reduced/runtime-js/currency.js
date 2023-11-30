/**
 * Your HTTP handling function, invoked with each request. This is an example
 * function that echoes its input to the caller, and returns an error if
 * the incoming request is something other than an HTTP POST or GET.
 *
 * In can be invoked with 'func invoke'
 * It can be tested with 'npm test'
 *
 * @param {Context} context a context object.
 * @param {object} context.body the request body if any
 * @param {object} context.query the query string deserialized as an object, if any
 * @param {object} context.log logging object with methods for 'info', 'warn', 'error', etc.
 * @param {object} context.headers the HTTP request headers
 * @param {string} context.method the HTTP request method
 * @param {string} context.httpVersion the HTTP protocol version
 * See: https://github.com/knative/func/blob/main/docs/function-developers/nodejs.md#the-context-object
 */

/**
 * Helper function that gets currency data from a stored JSON file
 * Uses public data from European Central Bank
 */


const data = require('./data/currency_conversion.json');
function _getCurrencyData(code) {
    if (data[code]) {
        return data[code];
    } else {
        return null;
    }
}

function getSupportedCurrencies(data) {
    logger.info('Getting supported currencies...');
    ret = _getCurrencyData(data)
    return { currency_codes: Object.keys(data) };
}



function _carry(amount) {
    const fractionSize = Math.pow(10, 9);
    amount.nanos += (amount.units % 1) * fractionSize;
    amount.units = Math.floor(amount.units) + Math.floor(amount.nanos / fractionSize);
    amount.nanos = amount.nanos % fractionSize;
    return amount;
}

function convert(request) {
    const from_code = request.currency_code;
    const euros = _carry({
        units: request.units / _getCurrencyData(from_code),
        nanos: request.nanos / _getCurrencyData(from_code)
    });
    euros.nanos = Math.round(euros.nanos);
    const result = _carry({
        units: euros.units * data[request.to_code],
        nanos: euros.nanos * data[request.to_code]
    });

    result.units = Math.floor(result.units);
    result.nanos = Math.floor(result.nanos);
    result.currency_code = request.to_code;
    const jsonString = JSON.stringify(result, null, 2);

    return jsonString;

}

function currencyFunction(eventBody) {
    //console.log("currency_body")
    //console.log(eventBody)
    const input = eventBody
    if (input.requestType === 'convert') {
        //console.log(convert(input))
        return convert(input);
    } else if (input.requestType === 'supported') {
        //console.log(getSupportedCurrencies(input))
        return getSupportedCurrencies(input);
    }
}

module.exports = {
    currencyFunction,
};


