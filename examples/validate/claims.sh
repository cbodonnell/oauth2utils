#!/bin/bash
ACCESS_TOKEN=$(cat $HOME/.tfarm/token.json | jq -r ".access_token")
curl -v http://localhost:9091/claims -H "Authorization: Bearer ${ACCESS_TOKEN}" | jq
