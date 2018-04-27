#!/bin/sh

while true; do
	nc -z "${POSTGRES_HOST:-localhost}" "${POSTGRES_PORT:-5432}"
	if [ "$?" = "0" ]; then
		break
	fi
	echo "Waiting for postgres..."
	sleep 1
done
./manage.py migrate && \
./manage.py collectstatic --no-input && \
exec ./start.py
