#!/usr/bin/env bash

DEC_PASS=$(gpg -dq "$1")

echo "$DEC_PASS"
