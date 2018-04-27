#!/bin/sh

./manage.py migrate && \
./manage.py collectstatic --no-input && \
exec ./start.py
