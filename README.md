### bright(ness)

rewrite of [this gist](https://gist.github.com/konradstrack/18fd96bd9d734f17f62f) from [this blog article](https://konradstrack.ninja/blog/changing-screen-brightness-in-accordance-with-human-perception/) (credit to @konradstrack) in go. Functionality, except the program won't fail if the display is turned off, brightness can increase all the way to 100%, and the args are slightly different.

```shell
$ git clone https://github.com/hrfee/brightness.git
$ cd brightness
$ sed -i "s/BACKLIGHT_PATH/PATH_TO_BACKLIGHT_FILE/g" bright.go
$ go build
$ ./bright [-U/-B/-D]
```

replace `PATH_TO_BACKLIGHT_FILE` with the path to the file which contains and controls your displays backlight. For me, this is `/sys/class/backlight/intel_backlight/brightness`.

where -U(up)/-B(big up) are small/large increase, and -D(down) is decrease.
