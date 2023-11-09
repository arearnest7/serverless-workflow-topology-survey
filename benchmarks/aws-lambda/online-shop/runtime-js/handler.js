
const currency = require('./currency.js');
const payment = require('./payment.js');


exports.handler = async (event) => {
    const request = event.body;
    const event_body = JSON.parse(request);
    console.log(event_body)
    if (event_body["type"] === 'currency') {
        result = currency.currencyFunction(event_body["body"]);
        return result
    } else {
        result = payment.paymentFunction(event_body["body"]);
        return result
    }
};



//local-testing
// function main(body) {
//     const input = JSON.parse(body);
//     if (input["type"] == "currency") {
//         //console.log(input["body"])
//         result = currency.currencyFunction(input["body"]);
//         console.log("Result")
//         console.log(result)
//         return result
//     }
//     else if (input["type"] == "payment") {
//         //console.log(input["body"])
//         result = payment.paymentFunction(input["body"]);
//         console.log("Result")
//         console.log(result)
//         return result
//     }
// }

// if (require.main == module) {
//     main(process.argv[2])
// }



