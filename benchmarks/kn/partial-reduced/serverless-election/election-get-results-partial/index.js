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
const redis = require('redis');
const moment = require('moment');

const client = redis.createClient({url: process.env.REDIS_URL, password: process.env.REDIS_PASSWORD});

if ("LOGGING_NAME" in process.env) {
        const loggingClient = redis.createClient({url: process.env.LOGGING_URL, password: process.env.LOGGING_PASSWORD});
}

const state_list = ['AK', 'AL', 'AR', 'AZ', 'CA', 'CO', 'CT', 'DC', 'DE', 'FL', 'GA', 'HI', 'IA', 'ID', 'IL', 'IN', 'KS', 'KY', 'LA', 'MA', 'MD', 'ME', 'MI', 'MN', 'MO', 'MS', 'MT', 'NC', 'ND', 'NE', 'NH', 'NJ', 'NM', 'NV', 'NY', 'OH', 'OK', 'OR', 'PA', 'RI', 'SC', 'SD', 'TN', 'TX', 'U'];

const handle = async (context, body) => {
	if ("LOGGING_NAME" in process.env) {
                await loggingClient.append(process.env.LOGGING_NAME, moment().format('MMMM Do YYYY, h:mm:ss a') + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "0" + "\n");
        }
	client.on('error', err => console.log('Redis Client Error', err));
        await client.connect();
        var results = [];
        for (var state in state_list) {
                state_results = await client.keys('election-results-' + state + '-*');
                var total_count = {"total": 0};
                for (var i = 0; i < state_results.length; i++) {
                        const cnt = await parseInt(client.get(state_results[i]));
                        total_count[state_results[i]] = cnt;
                        total_count["total"] += cnt;
                }

                var state_results_final = {};
                for (const [candidate, total] of total_count.entries()) {
                        if (candidate !== "total") {
                                state_results_final[candidate] = (total / total_count["total"] * 100.0).round();
                        }
                }
                results.push({
                        "state": state,
                        "disclaimer": "These vote counts are estimates. Visit https://github.com/tylerpearson/serverless-election-aws for more info.",
                        "total_count": total_count,
                        "results": state_results_final
                });
        }
        var response = {
                "isBase64Encoded": false,
                "statusCode": 200,
                "headers": {
                        "Access-Control-Allow-Origin": "*"
                },
                "body": results
        };
	if ("LOGGING_NAME" in process.env) {
                await loggingClient.append(process.env.LOGGING_NAME, moment().format('MMMM Do YYYY, h:mm:ss a') + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "1" + "\n");
        }
        return response;
}

// Export the function
module.exports = { handle };
