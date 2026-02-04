# DB Migration Tool

i need a db migration tool for this project that's as flexible as possible:

## Requirements

- language agnostic
- free; not freemium (i'm cheap)
- migrations defined in sql (i'm not learning another language for this)
- support multiple databases (sql + nosql where possible)

## Alternatives Considered

| Option     | Remarks                                                                      |
| ---------- | ---------------------------------------------------------------------------- |
| Liquibase  | Too complex, freemium                                                        |
| Atlas      | LOVE the idea of declarative code, but views are paid features (dealbreaker) |
| Goose      | no NoSQL option                                                              |
| FlyWay     | overkill                                                                     |
| db-migrate | written in Go, but seems like best open-source option                        |

## Decision

db-migrate seems like the most flexible, free option that i know of.
revisit if i find more
