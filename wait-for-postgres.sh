#!/bin/sh
/wait-for-it.sh db:5432 --timeout=30 --strict -- ./main
