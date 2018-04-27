#!/bin/bash

./manage.py migrate && \
./manage.py collectstatic --no-input && \
./start.py
