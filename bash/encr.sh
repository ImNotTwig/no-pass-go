#!/usr/bin/env bash

ENC_TEXT=$(gpg --encrypt --armor --always-trust --batch --yes --recipient "$1" --output - "$2")

echo "$ENC_TEXT"
