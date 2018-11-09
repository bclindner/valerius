![valerius](https://i.imgur.com/EsgTukM.png)

Valerius (stylized valerius) is a modular Discord bot for personal use.

## Features

### Banger Alerts

Valerius can detect if a banger is posted and respond with a dancing GIF.

You can enable this by specifying an array of `bangers` in the JSON config,
and optionally an array of `danceGifs` to respond with.

After this is enabled, simply type:

```
!banger
```

and Valerius will get a banger for you. If you post a banger Valerius
recognizes, it will respond with its own message and, optionally, a dance
GIF, if those are set.

### XKCD Search

Valerius can get the latest XKCD comic and its alt text:

```
!xkcd
```

You can also get a comic by number:

```
!xkcd <number of comic>
```
