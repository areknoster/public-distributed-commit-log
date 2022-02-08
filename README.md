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

## Current status
PDCL is in very early stages of the development. In fact, besides our test deployments for [acceptance testing](test/acceptance/acceptance_benchmark.txt) and [Open Pollution - demonstration app](https://github.com/jmichalak9/open-pollution) we're not aware of any usage right now.

So we're on the stage of discovery. It's bad if you want to use PDCL as a core feature for your enterprise system. But it's good if you want to innovate. [Say hello](https://github.com/areknoster/public-distributed-commit-log/discussions/34) and let's have fun together 💪
## Architecture
For reference on the terms used below, see [Glossary](#glossary).

PDCL distinguished 3 stakeholders: Topic Owner, Publisher, and Subscriber.

**Topic Owner** starts and maintains topic. They would create public repository, where they define:
- Validation rules for accepted messages: e.g. schema, data values ranges or signature owners.
  For an example, [see logic we used for acceptance testing](cmd/acceptance-sentinel/internal/validator/signed_validator.go)
- [Encoding of messages](storage/pbcodec/)
- [Storage implementation](storage/) - an implementation of [content addressable storage](https://en.wikipedia.org/wiki/Content-addressable_storage). Currently only [IPFS](https://ipfs.io/) based implementation is supported for publicly available topics. Remaining implementations can be used for local testing.
- [Head management implementation](thead) - the way to propagate information about new commits. Currently the main implementation for public topics is based on [IPNS](https://docs.ipfs.io/concepts/ipns/) 
- All of the above requirements are enforced by [Sentinel](sentinel) - an application, that validates messages and maintains reference to the topic head. Sentinel can be easily deployed using provided [Terraform modules](terraform/modules) 

**Publisher** (or producer in current code, see #36) - Adds data to chosen topic. They need to have access to the data that's useful for some topic. They use definitions from Topic repository, to transform the data they have and save it for public access. Blocking and concurrent implementation for producing is in `producer` package.

**Subscriber** (or consumer in current code, see #36) - Reads data from chosen topic. They can expect, that the messages that they read adhere the definitions defined by topic. The implementations that can be used for fetching newest commits and messages and handling them can be found in [consumer](consumer) package.


For an overview of how the logic of message exchange can look like see video below.

https://user-images.githubusercontent.com/38364298/151680034-d5f2c57b-2ee0-465f-bbb5-bbc6e7e58a96.mp4
## FAQ
### Do I need to deploy PDCL to start working with it?
No. The core logic PDCL provides is about publishing and subscribing to messages.
For instance, [integration tests](test/itest) use in-memory hash map as Storage implementation to test our logic. Similarly, you can start up your solution locally and move to public later.
## Status and contributing
Depending on the general interest and personal time commitment of the authors and other contributors, 
the project can be further developed. Please, create an issue, PR or reach out to the authors 
for assistance with deployment, questions or ideas.

## Roadmap
See [Project](https://github.com/users/areknoster/projects/3) to keep track on current development effort. Also, feel free to add new issues and discuss your needs.

## Glossary
**topic** - a unit which aggregates messages with same encoding, schema and semantic meaning
**message** - 
**data message** - a structured unit of data saved to PDCL. The schema, encoding and valid values are defined per topic. Data messages are immutable.
**commit** - a message belonging to PDCL topic (or topics) that refers to data messages and previous commit. Having a reference to a commit allows subscriber to find all previous commits and thus all previous messages. 
**head** (or thead, topic head) - the latest commit. 

## Authors
- [Arkadiusz Noster](https://github.com/areknoster)
- [Jakub Michalak](https://github.com/jmichalak9)

The project is a part of Engineering Thesis under the same title, 
made in Computer Science course in Warsaw University of Technology.

