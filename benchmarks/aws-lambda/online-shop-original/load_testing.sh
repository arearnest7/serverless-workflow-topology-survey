hey  -o csv -n 2500 -c 1 -d  '{"call": "cart","body": ""}' https://l7vqjvcteh375db5hw5hftd5oe0zlsyh.lambda-url.us-east-2.on.aws/ > cart.csv
hey  -o csv -n 2500 -c 1 -d '{"call": "list","body": ""}' https://l7vqjvcteh375db5hw5hftd5oe0zlsyh.lambda-url.us-east-2.on.aws/  > list.csv
hey  -o csv -n 2500 -c 1 -d '{"call": "ad","body": ""}' https://l7vqjvcteh375db5hw5hftd5oe0zlsyh.lambda-url.us-east-2.on.aws/ > ad.csv
hey   -o csv -n 2500 -c 1 -d '{"call": "checkout","body": "{\"UserId\":\"123\",\"UserCurrency\":\"EUR\",\"Address\":{\"StreetAddress\":\"123 Main St\",\"City\":\"Example City\",\"State\":\"CA\",\"Country\":\"US\",\"ZipCode\":12345},\"Email\":\"user@example.com\",\"CreditCard\":{\"CreditCardNumber\":\"4111111111111111\",\"CreditCardCvv\":123,\"CreditCardExpirationYear\":2025,\"CreditCardExpirationMonth\":12}}"}' https://l7vqjvcteh375db5hw5hftd5oe0zlsyh.lambda-url.us-east-2.on.aws/ > checkout.csv