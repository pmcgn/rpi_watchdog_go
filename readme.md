# Raspberry PI Watchdog

Make use of the Raspberry Pi's internal hardware watchdog via a docker container. It also supports monitoring of one webinterface. If the given webinterface is down for a given time, the watchdog will perform a hardreset of the RPi.

# Prerequisites

To prepare the host for using the watchdog, you have to run 'sudo apt install watchdog' on the host! (In case of raspbian)

# Running the container

To run the container, the only mandatory parameter is:
> --device /dev/watchdog:/dev/watchdog

The shortest way to make use of the watchdog is:
> docker run -d --name watchdog --restart always --device /dev/watchdog:/dev/watchdog pmcgn/rpiwatchdog

In this case the watchdog hardware will be triggered with a delay of 120s. After this timer is elapsed, the container has to run forever. In case of a crash of the docker runtime or the hosts operating system, the hardware watchdog of the Pi will trigger a hard reset.

## Optional environment variables

As mentioned in the introduction, this container is able to monitor a webinterface (which should run on the same host). When a URL is provided, this functionality will be enabled and can be controlled by other variables. The application checks for HTTP 2xx status codes. Any other error codes will be treated as errors and cause a reset. TLS is supported, but certificate errors are NOT treated as errors. With this, you are able to use selfsigned certificates.

The following variables are being used:

| Environment Variable             | Description                                                                                                               |
|----------------------------------|---------------------------------------------------------------------------------------------------------------------------|
| WATCHDOG_START_DELAY             | Delay before the watchdog starts in seconds. Required to avoid infinite loops of resets. Default is 120s, minimum is 60s. |
| HTTP_HEALTH_CHECK_URL            | Once a URL is provided, a periodic HTTP call against this URL will be performed.                                          |
| HTTP_HEALTH_CHECK_DELAY          | Delay between HTTP calls. Parameter is optional, default is 10s.                                                          |
| HTTP_HEALTH_CHECK_ERRORTHRESHOLD | Provides the maximum number of consecutive HTTP errors before a reset is triggerd. Parameter is optional, default is 5.   |