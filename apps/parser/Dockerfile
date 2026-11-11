FROM python:3.12-slim

WORKDIR /srv
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY app ./app

ENV PARSER_PORT=8090
EXPOSE 8090
CMD ["sh", "-c", "uvicorn app.main:app --host 0.0.0.0 --port ${PARSER_PORT}"]
