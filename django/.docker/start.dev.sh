#!/bin/bash

pipenv install
pipenv run python manage.py migrate
pipenv run python manage.py loaddata initial_data

tail -f /dev/null