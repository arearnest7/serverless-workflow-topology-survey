#!/bin/bash
BUILD=${1:-false}
PUSH=${2:-false}
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/feature-generation/feature-extractor
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/feature-generation/feature-orchestrator
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/feature-generation/feature-reducer
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/feature-generation/feature-status
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/feature-generation/feature-wait

