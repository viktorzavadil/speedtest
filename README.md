# Speed Test
Speed Test tools diagnostikují síťovou komunikaci. Měří latency, propustnost, ...

## Architecture
Jedná se o dvě nezávislé služby `client` a `target`.

```mermaid
---
config:
  theme: neutral
---
graph TD
    C[Client] --> N(network)
    N --> T(Target)
```

### Client
Služba `client` je zodpovědná za pouštění diagnostických testů proti službě `target`.

### Target
Služba `target` vystavuje REST API. 

#### API


## Test sequence diagram

```mermaid
---
config:
  theme: neutral
---
sequenceDiagram
    participant C as Client
    participant T as Target

    Note over C, T: blablabla
    loop lattency probe, runs every 500ms
        C ->> T: ping request
        T ->> C: ping response
    end

    loop download per streamer
        C ->> T: download request
        T ->> C: download response: infinity stream
    end

    loop upload per streamer
        C ->> T: upload request: infinity stream
        T ->> C: upload response
    end

```

## Run tests
Runs target service:
```bash
make run-target
```

Runs test with default settings:
```bash
make run-client
```

Runs test with customized settings:
```bash
make run-client DURATION=1 STREAMS=3
```

## Todo & Improvements
1. make decision: which endpoints need auth 
1. compute min, max lattency
1. implement CPU usage on client side
1. compute avg, min, max CPU usage
