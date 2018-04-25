#!/bin/bash

./manage.py migrate || exit 1
./manage.py runserver 0:8000 || exit 1
