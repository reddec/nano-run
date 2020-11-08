# Nano-Run

![sample-nano-run](https://user-images.githubusercontent.com/6597086/98463432-303e6900-21f6-11eb-9632-806b1c99813b.gif)

Lightweight async request runner. 

A simplified version of [trusted-cgi](https://github.com/reddec/trusted-cgi) designed
for async processing extreme amount of requests.

Main goals:

* Should have semi-constant resource consumption regardless of: 
  * number of requests,
  * size of requests,
  * kind of requests;
* Should be ready to run without configuration;
* Should be ready for deploying in clouds;
* Should support extending for another providers;
* Can be used as library and as a complete solution;
* **Performance (throughput/latency) has less priority** than resource usage.

Please note that the project is being developed in free time, as a non-profitable hobby project. 
All codes, bugs, opinions, and other related subjects should not be considered as the official position, official project,
or company-backed project to any of the companies for/with which I worked before or/and at present.   


![image](https://user-images.githubusercontent.com/6597086/95172239-9b58e200-07e9-11eb-8ca7-bf48d93a178b.png)

## Documentation

* [Installation](#installation)
* [Quick start](#quick-start)
* [Architecture overview](_docs/flow.md)
* [API](_docs/api.md)
* [API Authorization](_docs/authorization.md)
* [UI](_docs/ui.md)
* [UI Authorization](_docs/ui_authorization.md)
* [Unit configuration](_docs/unit.md)
* [Docker](_docs/docker.md)
* [Cron-like scheduler](_docs/cron.md)

## Stability

(After 1.0.0)

We are trying to follow semver:

* Patch releases provides fixes or light improvements without migration
* Minor releases provides new functionality and provides automatic migration if needed
* Major releases provides serious platform changes and may require manual or automatic migration

Within one major release, it guarantees forward compatibility: new versions can use data from previous versions, but not vice-versa.

## Reproducible binaries

The project tries to follow best practices providing reproducible binaries: it means, that
you can verify that complied binaries will be exactly the same (byte to byte) as if you will compile it by yourself
by following our public instructions.  

## License

The project (source code and provided official binaries) are licensed
under Apache-2.0 (see License file, or in [plain English](https://tldrlegal.com/license/apache-license-2.0-(apache-2.0))) and suitable 
for personal, commercial, government, and others usage without restrictions as long as it used with abiding
license agreement.

Do not forget to bring your changes back to the project. I will
 be happy to assist you with PR. 

The project uses external libraries that may be distributed
under other licenses.   

## Installation

### Debian/Ubuntu

(recommended)

Tested on 20.04 and 18.04, but should good for any x64 version.

Add the repository (only once)

```bash
sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 379CE192D401AB61
echo "deb https://dl.bintray.com/reddec/debian all main" | sudo tee /etc/apt/sources.list.d/reddec.list
```

Install or update nano-run

```bash
sudo apt update
sudo apt install nano-run
```

Automatically creates service `nano-run.service`.

### Binary

Download and unpack desired version in [releases](https://github.com/reddec/nano-run/releases).

### Docker

`docker pull reddec/nano-run`

### From source

Requires go 1.14+

`go get -v github.com/reddec/nano-run/cmd/...`

### Ansible for debian servers


```yaml
- name: Add an apt key by id from a keyserver
  become: yes
  apt_key:
    keyserver: keyserver.ubuntu.com
    id: 379CE192D401AB61
- name: Add repository
  become: yes
  apt_repository:
    repo: deb https://dl.bintray.com/reddec/debian all main
    state: present
    filename: reddec
- name: Install nano-run
  become: yes
  apt:
    name: nano-run
    update_cache: yes
    state: latest
```

## Quick start

**(optional) initialize configuration**

    nano-run server init

it will create required directories and files in a current working directory. 

**define a unit file**

Create minimal unit file (date.yaml) that will return current date (by unix command `date`) and put it
in directory `run/conf.d/`

_run/conf.d/date.yaml_
```yaml
command: date
```

**start nano-run**

    nano-run server run
    
    
now you can open ui over http://localhost:8989 or do API call: `curl -X POST http://localhost:8989/api/date/`
