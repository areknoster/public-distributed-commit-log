# Public Distributed Commit Log
![Logo](logo.svg)

Public Distributed Commit Log (PDCL) is an open source Go library that provides components for building public and distributed message exchange platform.

## Why?
There is a category of data, that reasonably should be public. Think of everything that you can go and observe yourself. 
Traffic on the roads, weather and space measurements or public companies notations. 
Currently, despite the fact, that you can access them for free, 
in order to fully utilize them in scale you need to oftentimes pay to some vendor.

We think, that there should be a solution to make that data accessible and independent. 

Let's consider event or commit log platforms that are currently available:
* Platforms like Kafka and RabbitMQ, whereas could theoretically be opened for public writes and reads do not easily
  replicate and the access to written data can be easily witheld by provider. Also, they don't provide tooling that would be required for such public message exchange.
* Blockchain solutions. PDCL at first glance might seem similar to some. Yet PDCL abandons distributed consensus mechanism 
in favour of central manager (Sentinel) who provides messages validation and controls their inflow.
Is it the right call? Distributed consensus is expensive by design. That's not a very good fit for place to store public messages.
Public means cheap and accessible. Our argument is, that the freedom lies in undeniable access to past data and freedom of fork.
That's very similar to how Open Source works. And that's PDCL aims to provide.
* There exist public and distributed p2p platforms like torrents and IPFS. But they are used for storing static files. They don't provide discovery tooling. 
In fact, basic PDCL setup is on top of IPFS. What PDCL enables is dynamic publishing and dynamic consumption of validated messages.

So this is a niche that PDCL tries to fill. 
A platform for append-only message log, that's scalable, distributed and public, but that's not expensive to operate.

## Documentation & architecture
TBD, but here is a nice video to give an intuition on how it works in general.

https://user-images.githubusercontent.com/38364298/151680034-d5f2c57b-2ee0-465f-bbb5-bbc6e7e58a96.mp4

The cloud in the middle is common, [content addressable storage](https://en.wikipedia.org/wiki/Content-addressable_storage) - in current main implementation it is based on [IPFS](https://ipfs.io/)

## Status and contributing
Depending on the general interest and personal time commitment of the authors and other contributors, 
the project can be further developed. Please, create an issue, PR or reach out to the authors 
for assistance with deployment, questions or ideas.

## Roadmap
* Investigate possibilities for more Message Storage implementations and extending current ones. 
In particular implementation based on Filecoin or other free and distributed storage could be looked upon.
* Add components for error handling and traversing logic for consumer. Current implementation doesn't handle missing messages very gracefully.
A mechanism to implement traverse stop could also be useful.
* Implement Broker, so that messages can be batched from multiple sources and writing and reading does not require running your own message storage solution.
* Cluster mode for both Sentinel and IPFS nodes for redundancy and better scalability
* Terraform modules for other cloud providers, docker compose and helm chart for wider deployment support
* Template repo with example topic definition
* Performance tuning (especially investigate IPFS node utilization under load which seems to be abnormally high)


## Authors
- [Arkadiusz Noster](https://github.com/areknoster)
- [Jakub Michalak](https://github.com/jmichalak9)

The project is a part of Engineering Thesis under the same title, 
made in Computer Science course in Warsaw University of Technology.

