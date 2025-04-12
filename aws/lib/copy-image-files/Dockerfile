ARG IMAGE
FROM $IMAGE AS sourceimage

FROM public.ecr.aws/lambda/python:3.13

ARG SRC
COPY --from=sourceimage ${SRC} /staging

COPY requirements.txt ${LAMBDA_TASK_ROOT}

RUN pip install -r requirements.txt

COPY handler.py ${LAMBDA_TASK_ROOT}

CMD ["handler.handler"]
