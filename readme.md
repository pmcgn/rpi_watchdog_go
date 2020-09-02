# Raspberry PI Watchdog

Make use of the Raspberry Pi's internal hardware watchdog via a docker container. It also supports monitoring of one webinterface. If the given webinterface is down for a given time, the watchdog will perform a hardreset of the RPi.

# Prerequisites

To prepare the host for using the watchdog, you have to run ```sudo apt install watchdog``` on the host! (In case of raspbian)

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

## Running the container with activated HTTP monitoring

In this case, the container needs to know which URL shall be monitored. An optional delay between two consecutive health checks, can be provided. This delay can be larger than the internal timelimit of 15s from the watchdog hardware. This application will count the failed responses and if the limit is reached, stop writing to the watchdog device. This will trigger the reset. Yes, this is the most brutal way to get a webinterface up and running again. You're welcome to send me pull requests.  To avoid reset loops after a reboot, there must be one successful HTTP call before a reset will be performed.

To run the container with active HTTP monitoring, run the following docker command (change your URL):
> docker run -d --name watchdog \\
> --device /dev/watchdog:/dev/watchdog \\
> -e WATCHDOG_START_DELAY=120 \\
> -e HTTP_HEALTH_CHECK_URL=http://<i></i>192.168.1.10 \\
> -e HTTP_HEALTH_CHECK_DELAY=30 \\
> -e HTTP_HEALTH_CHECK_ERRORTHRESHOLD=10 \\
> --restart always \\
> pmcgn/rpiwatchdog

# IMPORTANT: Stopping the watchdog
Please note, that the watchdog hardware can't be deactivated once it is active! In future releases, I plan to implement a mechanism to pause the HTTP healtch check. For now, it runs forever. **A stop of this container will cause a reboot after 15 seconds!!!** 
If you want to upgrade the monitored container, you have to change the restart policy of the watchdog container to 'no' and then do a reboot of the host. After the reboot, the watchdog is inactive. This is the exact procedure:

1. Change restart policy of watchdog container via: ```docker update --restart=no watchdog```
1. On Raspian, reboot with ```sudo shutdown -r now```
1. Upgrade/replace the monitored container
1. Change restart policy of watchdog container ```docker update --restart=no watchdog``` 
1. start watchdog container wit ```docker start watchdog```
