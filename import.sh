#!/bin/sh

cat import.cql | cypher-shell -u neo4j -p password --format verbose
