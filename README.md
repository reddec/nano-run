# Nano-Run

Lightweight async request runner. 

A simplified version of [trusted-cgi](https://github.com/reddec/trusted-cgi) designed
for async processing extreme amount of requests.

## Goals

* Minimal requirements for host;
* Should have semi-constant resource consumption regardless of: 
  * number of requests,
  * size of requests,
  * kind of requests;
* Should be ready to run without configuration;
* Should be ready for deploying in clouds;
* Should support extending for another providers;
* Can be used as library and as a complete solution;
* **Performance (throughput/latency) has less priority** than resource usage.