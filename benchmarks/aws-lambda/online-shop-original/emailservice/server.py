
from concurrent import futures
from jinja2 import Environment, FileSystemLoader, select_autoescape, TemplateError
from google.api_core.exceptions import GoogleAPICallError
import json


env = Environment(
    loader=FileSystemLoader('./templates'),
    autoescape=select_autoescape(['html', 'xml'])
)
template = env.get_template('confirmation.html')   #confirmation.html -> not found 

def send_message(sender, envelope_from_authority, header_from_authority, envelope_from_address, simple_message):
	return {"rfc822_message_id": 1234}

def send_email(email_address, content):
    response = send_message(
        sender = [1234, "us-east", 3456],
        envelope_from_authority = '',
        header_from_authority = '',
        envelope_from_address = '',
        simple_message = {
            "from": {
            "address_spec": '',
            },
            "to": [{
                "address_spec": email_address
            }],
            "subject": "Your Confirmation Email",
            "html_body": content
        }
    )
    #print("Message sent: {}".format(response["rfc822_message_id"]))
    return content


def SendOrderConfirmation(request):
    email = request["email"]
    order = request["order"]
    #print(order)
    try:
        confirmation = template.render(order=order)
    except TemplateError as err:
        return {"error": "Template error"}

    try:
        return send_email(email, confirmation)
    except GoogleAPICallError as err:
        return {"error": "Google API error"}
        

def lambda_handler(event,context):
    #print(event)
    event_data = event['body']
    event_body = json.loads(event_data)
    return SendOrderConfirmation(event_body)
