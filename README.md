# Testing Implicit scans - copy of Reports Service Repo for findings

Cloudbees SaaS platform Reports-service manage data for VSM metrics/Widget.

## Purpose

- It exposes an single endpoint for UI to get Widget specific data from OpenSearch using Widget identfier
- Contains a framework to map a Widget Identifier with associated set of queries to support multiple metrics
- The Mapping is currently in Repo. Will be moved to Cassandra DB to support Custom Widgets later on
- The Framework will also maintain UI layout definition
- Opensearch Client intergated to execute search queries. 
