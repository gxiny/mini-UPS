FROM python:3
RUN mkdir /static && chown nobody:nogroup /static
ENV PYTHONUNBUFFERED 1
RUN apt-get update && \
	apt-get install -y libev-dev && \
	rm -rf /var/lib/apt/lists/*
WORKDIR /code
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY . .
EXPOSE 8000
CMD ["./start.sh"]
