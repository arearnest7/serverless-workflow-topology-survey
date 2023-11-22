import json
import random
def ListRecommendations(request):
    data = request["data"]
    product_entries = data.strip().split("------------------------")
    random.shuffle(product_entries)
    selected_entries = product_entries[:5]
    return selected_entries

    
def lambda_handler(event,context):
    event_data = event['body']
    event_body = json.loads(event_data)

    return ListRecommendations(event_body)
