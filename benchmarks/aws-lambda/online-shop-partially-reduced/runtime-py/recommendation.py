from concurrent import futures
from jinja2 import Environment, FileSystemLoader, select_autoescape, TemplateError
from google.api_core.exceptions import GoogleAPICallError
import json
import random
def ListRecommendations(request):
    data = request
    product_entries = data.strip().split("------------------------")
    random.shuffle(product_entries)
    selected_entries = product_entries[:5]
    return selected_entries

    
