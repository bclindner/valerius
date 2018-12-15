![valerius](https://i.imgur.com/EsgTukM.png)

Valerius (stylized valerius) is a modular Discord bot for personal use.

## Usage

Valerius uses JSON to configure all of its commands from generic shells.
The program looks for `valerius.json` in the current working directory, but you can set the config file with `valerius -conf <path_to_configfile>`.

It's best to learn by example here. This is an example config similar to the one the author uses, using all of the currently available command types:

```json
{
  "botToken": "YOUR_DISCORD_BOT_TOKEN_HERE",
  "commands": [
    {
      "name": "Hello",
      "type": "pingpong",
      "options": {
        "triggers": ["!hello"],
        "response": "Hello world!"
      }
    },
    {
      "name": "RegExample",
      "type": "regexpingpong",
      "options": {
        "trigger": "https?://[\\S]+",
        "response": "Yep, that looks like a link."
      }
    },
    {
      "name": "Bangers",
      "type": "randompingpong",
      "options": {
        "triggers": [
          "https://www.youtube.com/watch?v=kZf91MAwS7s",
          "https://www.youtube.com/watch?v=UGymAxj8DYI",
          "https://www.youtube.com/watch?v=6Wo-u8vQn4U",
          "https://www.youtube.com/watch?v=u3GTsFwJ5Uo",
          "https://www.youtube.com/watch?v=Emiu-xcLlJU",
          "https://www.youtube.com/watch?v=l7PD62YHRQk",
          "https://www.youtube.com/watch?v=B2jVbSI9H4o",
          "https://www.youtube.com/watch?v=pcamjcoRmrQ",
          "https://www.youtube.com/watch?v=B1lNhNHdoPI"
        ],
        "responses": [
          "https://i.kym-cdn.com/photos/images/original/001/331/773/518.gif",
          "https://media.giphy.com/media/PSKAppO2LH56w/giphy.gif",
          "https://i.kym-cdn.com/photos/images/newsfeed/000/427/549/b9d.gif",
          "https://i.kym-cdn.com/photos/images/newsfeed/000/032/802/ninja-dance.gif",
          "https://media1.tenor.com/images/e88f1c4b6d3ac98bde66db24fb73441d/tenor.gif?itemid=5586778",
          "https://media.giphy.com/media/7isbcNAx367qU/200.gif",
          "https://thumbs.gfycat.com/AgitatedGleefulEmperorshrimp-size_restricted.gif",
          "https://media0.giphy.com/media/CDzdJSkC4iyLC/giphy.gif"
        ],
        "responsePrefix": "🚨IT'S🚨A🚨BANGER🚨\n"
      }
    },
    {
      "name": "XKCD",
      "type": "xkcd"
      "options": {
        "prefix": "!xkcd",
      }
    },
    {
      "name": "IASIP",
      "type": "iasip",
      "options": {
        "trigger": "!iasip",
        "fontpath": "./textile.ttf",
        "quality": 100
      }
    }
  ]
}
```

In essence, a configuration will have a `botToken` property with (surprise) the Discord bot token to log in with, and a `commands` array which contains an array of command configuration objects, each with a `name` to call it by in logs. a `type` to base the command off of, and an `options` object to actually configure the command type with. Put the config file in the same directory as the executable and fire it up and you should have a working bot you can customize on-the-fly, without need for recompilation or source code editing.

You can also set a log file with the `-log` argument. Using this argument automatically puts the logger in JSON output mode for easy log parsing.
