# Changelog

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

### [1.0.2](https://github.com/areknoster/public-distributed-commit-log/compare/v1.0.1...v1.0.2) (2022-01-11)


### Features

* **acceptance-sentinel:** implement sentinel for acceptance testing ([cab5e6d](https://github.com/areknoster/public-distributed-commit-log/commit/cab5e6daeae9f9472c30537639d3afc1b75449c2))
* **committer:** Add interval and max buffer commiter ([1c42066](https://github.com/areknoster/public-distributed-commit-log/commit/1c42066b21ce4081d05b6cc240afb8b8ef125a4d))
* **ratelimiting:** implement basic rate limiting ([#10](https://github.com/areknoster/public-distributed-commit-log/issues/10)) ([832a4f4](https://github.com/areknoster/public-distributed-commit-log/commit/832a4f4841034f8a49eb9a548ff82f6b534fa6a4))
* **signing:** add package for signed messages wrappers ([72cc34f](https://github.com/areknoster/public-distributed-commit-log/commit/72cc34fe22b776cbf48cd4b46f08b43eadecdec0))
* **storage:** add basic ipfs based storage implementation ([feaf8f9](https://github.com/areknoster/public-distributed-commit-log/commit/feaf8f9e278aa8af75bb40e7229d881ae4d11bbf))


### Refactoring

* **storage:** separate codec logic ([abcbd6a](https://github.com/areknoster/public-distributed-commit-log/commit/abcbd6acc9fb903472827cd9fcc5635f6b48ae3c))


### Test

* **infra:** add image registry ([09259bc](https://github.com/areknoster/public-distributed-commit-log/commit/09259bc5255928aaedc024fbe6f3e104304bb6e0))
* **infra:** add initial pdcl-test project setup ([f57c9d0](https://github.com/areknoster/public-distributed-commit-log/commit/f57c9d0e44315d6c8508cd9952f99b3abffddc54))
* **infra:** setup GCP project for acceptance testing ([4756adc](https://github.com/areknoster/public-distributed-commit-log/commit/4756adcab247209b60dadbfbc291b9829fc1bbcb))

### 1.0.1 (2021-12-11)


### Features

* **go.mod:** initialize go mod ([f083de2](https://github.com/areknoster/public-distributed-commit-log/commit/f083de2d6a63c4cc78f17a395b5f82bfe55a4e9a))
* Initial commit ([1b21ebc](https://github.com/areknoster/public-distributed-commit-log/commit/1b21ebc4eeadc981e858714d82cb06485f51cc65))
* **openpollution:** implement PDCL first-to-last consumer ([45bdfa7](https://github.com/areknoster/public-distributed-commit-log/commit/45bdfa7c5e4a5354b51c51c676e5c3600b59a3d4))
* **producer:** add producer and consumer interfaces definitions ([8e745fa](https://github.com/areknoster/public-distributed-commit-log/commit/8e745fa2eee0788afc04de39e60459ce8b99f50a))
* **producer:** implement random producer based on local filesystem ([468ca38](https://github.com/areknoster/public-distributed-commit-log/commit/468ca384ef9979748a9efe958905cf6287860c42))
* **sentinel:** add part of sentinel abstract definitions and naive implementations ([3ec8b82](https://github.com/areknoster/public-distributed-commit-log/commit/3ec8b8229fcef41476a8f6d89ae25384c851d456))
* **sentinel:** first iteration of openpollution local file storage based sentinel ([09f9a4b](https://github.com/areknoster/public-distributed-commit-log/commit/09f9a4b6dedc4d9beb7ca24b76083d268592beb7))
* **storage:** add local file storage implementation ([969072f](https://github.com/areknoster/public-distributed-commit-log/commit/969072f94ca531f3999aa41e02f854233c9963fd))
* **storage:** adjust interfaces, so that it's possible to unmarshall to user-defined type ([0edd2df](https://github.com/areknoster/public-distributed-commit-log/commit/0edd2df17fa07ee00fe0a42d500e80b38e7801dc))


### Test

* add integration test ([2280882](https://github.com/areknoster/public-distributed-commit-log/commit/22808820b9acbfd17520d9ecdc1cbacf481d89f0))
* **consumer:** add unit tests for general consuming logic ([98af315](https://github.com/areknoster/public-distributed-commit-log/commit/98af31529d9586add939bf5a2f4bce6d0b486a7d))
