#!/bin/bash
PATH=${1:-kn}
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/feature-extractor-partial
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/feature-orchestrator-wsr

