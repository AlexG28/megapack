FROM python:3.11-slim-buster 

WORKDIR /app

COPY megapack.py .
COPY requirements.txt .

RUN pip install -r requirements.txt

CMD ["python", "megapack.py"]
