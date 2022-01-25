# Changelog

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

## [1.2.0](https://github.com/areknoster/public-distributed-commit-log/compare/v1.1.1...v1.2.0) (2022-01-25)


### Features

* **cmd/acceptance-*:** implement applications for manual operations on acceptance topic ([118817b](https://github.com/areknoster/public-distributed-commit-log/commit/118817b13b02686af3082bc6bc7cfdbdd83d0d7b))
* **producer:** implement concurrent producer ([da43802](https://github.com/areknoster/public-distributed-commit-log/commit/da43802ebf85672a81a58e1d4d2efb42b11df754))
* **sentinel:** configurable maxbuffer commmiter ([d536f73](https://github.com/areknoster/public-distributed-commit-log/commit/d536f73ca4ac803b90daef9d7e48845c53a9ce1d))


### Refactoring

* **terraform:** modularize deployment ([a52e2c9](https://github.com/areknoster/public-distributed-commit-log/commit/a52e2c9f43db785369039155dd7a99e1f98ac7fe))


### Test

* **acceptance:** basic acceptance performance benchmark ([df6b7c2](https://github.com/areknoster/public-distributed-commit-log/commit/df6b7c2b5fe51aa1cf717dd8924ffc0b579a3d19))

### [1.1.1](https://github.com/areknoster/public-distributed-commit-log/compare/v1.1.0...v1.1.1) (2022-01-17)


### Features

* **committer:** support publishing head in ipns ([fbe200f](https://github.com/areknoster/public-distributed-commit-log/commit/fbe200f7550c2b1587437cd191019e7737d8f810))
* **consumer:** handle message reading for commits and messages separately ([5305c18](https://github.com/areknoster/public-distributed-commit-log/commit/5305c184dbceec1091db74ff0e0391965a3a2fdb))
* expose grpc endpoint for returning ipns address ([c49245d](https://github.com/areknoster/public-distributed-commit-log/commit/c49245dc6b955f8f7379704a7f2693e33ef1ac0c))
* **ipns:** use ipns resolver in consumer ([f3b77a7](https://github.com/areknoster/public-distributed-commit-log/commit/f3b77a7b2b1ff3663bcc100f2e9c88d6bf299cc2))


### Bug Fixes

* **acceptance-consumer:** minor fixes so that acceptance consumer works ([bcc4582](https://github.com/areknoster/public-distributed-commit-log/commit/bcc4582d6592332cc416e53d1a4c7b79bf8561ee))
* **crypto:** adjust for buggy libp2p ed25519 key pointer passing ([6229faa](https://github.com/areknoster/public-distributed-commit-log/commit/6229faad73e149346669710a2f1bb724088bc0ab))


### Test

* **acceptance-sentinel:** adjust acceptance-sentinel for getting secret from GCP ([3653d61](https://github.com/areknoster/public-distributed-commit-log/commit/3653d61277af7188ccdde06e9245cc1adbfc4a9a))
* **infra:** add definitions for IPNS secret ([83dd780](https://github.com/areknoster/public-distributed-commit-log/commit/83dd780ba00565a11836a529cadab970195378e8))
* **infra:** add ipfs server deployment ([b0d2887](https://github.com/areknoster/public-distributed-commit-log/commit/b0d2887f6894a9d2699c92e529c3fa7621d28558))
* **infra:** add key for ipns to KMS ([aacd459](https://github.com/areknoster/public-distributed-commit-log/commit/aacd459c3241e67e3ea7ce98541253159f5add87))
* **infra:** add network rules to access ipfs ([83c0605](https://github.com/areknoster/public-distributed-commit-log/commit/83c06052b9d9a9b7e0e01977ec3c41726363ee68))
* **infra:** add sentinel cloud run deployment ([eb2e74e](https://github.com/areknoster/public-distributed-commit-log/commit/eb2e74eff6dcd2ee4404cb218724bc94d6705919))
* **infra:** adjust sentinel and ipfs node setup ([2df41cd](https://github.com/areknoster/public-distributed-commit-log/commit/2df41cd03059a1778cf2674c88fe87b94f49f0ce))


### Refactoring

* **consumer:** adjust to ipns head resolver ([8101872](https://github.com/areknoster/public-distributed-commit-log/commit/8101872cfb4cbc908ebc58cfd50ab1cbce13998f))
* **pdclcrypto:** move common key logic to pdclcrypto package ([9cdc088](https://github.com/areknoster/public-distributed-commit-log/commit/9cdc0883e1f0643ccf5836a94ab46ceea98051a5))
* **storage:** refactor content storage, message storage and codec ([01e4295](https://github.com/areknoster/public-distributed-commit-log/commit/01e4295c29d3ce919575b09d518bae9561237ccd))

## [1.1.0](https://github.com/areknoster/public-distributed-commit-log/compare/v1.0.1...v1.1.0) (2022-01-13)


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

## [1.1.0](https://github.com/areknoster/public-distributed-commit-log/compare/v1.0.1...v1.1.0) (2022-01-11)


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

>>>>>>> 96af461... chore(release): 1.1.0
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
