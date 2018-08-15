package main

const composeTemplate = `
neo4j:
  image: neo4j:3.4.5
  ports:
   - "7474:7474"
   - "7687:7687"
  volumes:
   - {NEO4J}/data:/data
   - {NEO4J}/logs:/var/logs
   - {NEO4J}/conf:/conf
   - {NEO4J}/import:/import
   - {NEO4J}/plugins:/var/lib/neo4j/plugins
  environment:
   - NEO4J_dbms_memory_heap_max__size=4G
   - NEO4J_dbms_security_procedures_unrestricted=apoc.\\\*
`
