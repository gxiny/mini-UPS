#!/bin/bash

./manage.py migrate && \
./manage.py collectstatic --no-input && \
./manage.py runserver 0:8000
