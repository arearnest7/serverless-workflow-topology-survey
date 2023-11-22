hey -n 200 -c 10 -d '{"call": "cart","body": ""}' https://u6revnmt4q4wxtw3euqojxogua0rllel.lambda-url.us-east-2.on.aws/


hey -n 200 -c 10 -d '{"call": "list","body": ""}' https://u6revnmt4q4wxtw3euqojxogua0rllel.lambda-url.us-east-2.on.aws/



hey -n 200 -c 10 -d '{"call": "checkout","body": "{\"UserId\":\"123\",\"UserCurrency\":\"EUR\",\"Address\":{\"StreetAddress\":\"123 Main St\",\"City\":\"Example City\",\"State\":\"CA\",\"Country\":\"US\",\"ZipCode\":12345},\"Email\":\"user@example.com\",\"CreditCard\":{\"CreditCardNumber\":\"4111111111111111\",\"CreditCardCvv\":123,\"CreditCardExpirationYear\":2025,\"CreditCardExpirationMonth\":12}}"}' https://u6revnmt4q4wxtw3euqojxogua0rllel.lambda-url.us-east-2.on.aws/