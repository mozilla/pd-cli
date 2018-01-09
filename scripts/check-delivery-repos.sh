#!/bin/sh
cd $(dirname $0)

if [ -z $GH_ACCESS_TOKEN ]; then
    echo "GH_ACCESS_TOKEN env var not set"
    exit 1
fi

# check buildhub since it's in mozilla services
echo "Checking: github.com/mozilla-services/buildhub"
pd-cli repo check all -q "mozilla-services/buildhub"

for R in PollBot delivery-dashboard doorman; do
    echo
    echo "Checking: github.com/mozilla/$R"
    pd-cli repo check all -q "github.com/mozilla/$R"
done
echo
