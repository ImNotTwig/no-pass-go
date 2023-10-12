#!/usr/bin/env bash

ENC_TEXT=$(echo "$2" | gpg --encrypt --armor --always-trust --batch --yes --recipient "$1" --output -)

echo "$ENC_TEXT"
