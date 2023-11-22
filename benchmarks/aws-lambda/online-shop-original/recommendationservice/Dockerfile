# Use the AWS Lambda base image for Python 3.9
FROM public.ecr.aws/lambda/python:3.9
COPY . ${LAMBDA_TASK_ROOT}
COPY requirements.txt ./requirements.txt
RUN  pip3 install -r requirements.txt --target "${LAMBDA_TASK_ROOT}"
CMD [ "server.lambda_handler" ]

