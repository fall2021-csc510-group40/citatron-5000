FROM python:3.9.7

ENV PYTHONUNBUFFERED=1
COPY requirements.txt /
RUN pip install -r /requirements.txt

VOLUME /citation-app
WORKDIR /citation-app

CMD ["python", "/citation-app/main.py", "/citation-app/config.json"]

