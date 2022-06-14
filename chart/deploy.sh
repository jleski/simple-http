#!/usr/bin/env bash

helm upgrade -i simple-http --reset-values -f values.yaml -n simple-http ./