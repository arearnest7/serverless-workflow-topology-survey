FROM python:3.10
WORKDIR /app
COPY . .
RUN pip install -r requirements.txt
RUN pip install flask flask_restful
CMD ["python","app.py"]
