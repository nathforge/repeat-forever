# repeat-forever

Run a command every X seconds. It's a bit like cron, but stupider and designed
for Docker images.

e.g:

```
$ repeat-forever --every=5s ls /
$ repeat-forever --every=5s --timeout=1s find /
```
