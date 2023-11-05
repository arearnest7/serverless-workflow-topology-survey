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

const fs = require('fs')
const redis = require('redis');
const http = require('http');

const client = redis.createClient({url: process.env.REDIS_URL});

const state_list = ['AK', 'AL', 'AR', 'AZ', 'CA', 'CO', 'CT', 'DC', 'DE', 'FL', 'GA', 'HI', 'IA', 'ID', 'IL', 'IN', 'KS', 'KY', 'LA', 'MA', 'MD', 'ME', 'MI', 'MN', 'MO', 'MS', 'MT', 'NC', 'ND', 'NE', 'NH', 'NJ', 'NM', 'NV', 'NY', 'OH', 'OK', 'OR', 'PA', 'RI', 'SC', 'SD', 'TN', 'TX', 'U'];

const handle = async (context, body) => {
	var results = [];
	for (var state in state_list) {
		client.keys('election-results-' + state + '-*', function (err, state_results) {
    			if (err) {
				return console.log(err);
			}
    			for(var i = 0, len = keys.length; i < len; i++) {
        			console.log(keys[i]);
  			}
		});
		var total_count = {"total": 0};
		for (var i = 0; i < state_results; i++) {
			const cnt = await client.get(state_results[i]);
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
	return response;
}

// Export the function
module.exports = { handle };
