#!/bin/bash
BUILD=${1:-false}
PUSH=${2:-false}
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/sentiment-analysis/sentiment-cfail
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/sentiment-analysis/sentiment-db
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/sentiment-analysis/sentiment-main
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/sentiment-analysis/sentiment-product-or-service
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/sentiment-analysis/sentiment-product-result
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/sentiment-analysis/sentiment-product-sentiment
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/sentiment-analysis/sentiment-read-csv
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/sentiment-analysis/sentiment-service-result
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/sentiment-analysis/sentiment-service-sentiment
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/sentiment-analysis/sentiment-sfail
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/sentiment-analysis/sentiment-sns

