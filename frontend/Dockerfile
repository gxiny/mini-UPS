FROM python:3-alpine
RUN mkdir /static && chown nobody:nogroup /static
ENV PYTHONUNBUFFERED 1
WORKDIR /code
COPY requirements.txt .
RUN apk add --no-cache --virtual .build-deps \
	libc-dev \
	libev-dev \
	postgresql-dev \
	gcc \
	&& pip install --no-cache-dir -r requirements.txt \
	&& apk del .build-deps && apk add --no-cache libev libpq
COPY . .
EXPOSE 8000
CMD ["./start.sh"]
