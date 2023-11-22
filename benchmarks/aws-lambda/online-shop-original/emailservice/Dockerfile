FROM public.ecr.aws/lambda/python:3.9 as authPythonLambda


COPY tracing.py ${LAMBDA_TASK_ROOT}
COPY server.py ${LAMBDA_TASK_ROOT}

COPY requirements.txt ./requirements.txt

RUN  pip3 install -r requirements.txt --target "${LAMBDA_TASK_ROOT}"

# Set the CMD to handler
CMD [ "server.lambda_handler" ]
