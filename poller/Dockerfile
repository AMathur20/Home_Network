FROM python:3.12-slim

WORKDIR /app

COPY . /app
RUN apt-get update && apt-get install -y curl \
    && rm -rf /var/lib/apt/lists/* \
    && pip install --no-cache-dir -r requirements.txt

CMD ["python", "main.py"]
