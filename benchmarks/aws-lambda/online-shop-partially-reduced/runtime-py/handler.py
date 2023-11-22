import recommendation
import json
import email_service

def lambda_handler(event,context):
    event_data = event['body']
    event_body = json.loads(event_data)
    print(event_body)
    print(event_body['type'])
    print(event_body['data'])
    if event_body['type'] == 'recommendation':
        print("recommendinggggg.....")
        recommendation_result = recommendation.ListRecommendations(event_body['data'])
        return recommendation_result
    elif event_body['type'] == 'email':
        print("email...........")
        return email_service.SendOrderConfirmation(event_body['data'])
    else:
        raise Exception('Invalid request type: {}'.format(event['type']))

