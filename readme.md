# LocalHoster

This proof of concept allows you to host applications which are accessible with conveniently short urls in the browser. It starts an http server at 127.0.0.1:80 and adds the configured hosts to the hosts file of the computer.

Currently a default file browser that exposes a directory at a certain host is the only supported application.

## Configuration
The configuration is read from the following path(s):

Windows: `%USERPROFILE%/localhoster.yaml`

Unix: `/etc/localhoster.yaml`

## Example configuration

This example will make the Font files accessible by navigating to `files/` in the browser.

``` yaml
hosts:
  fonts:
    type: files
    path: C:\Windows\Fonts # could be /home/your-name/Downloads on unix
```
