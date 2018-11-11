![valerius](https://i.imgur.com/EsgTukM.png)

Valerius (stylized valerius) is a modular Discord bot for personal use.

## Usage

Valerius uses JSON to configure all of its commands from generic shells.

It's best to learn by example here. This is an example config similar to the one the author uses:

```json
{
  "botToken": "YOURDISCORDBOTTOKENGOESHERE0123456789",
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
      "name": "Bangers",
      "type": "randompingpong",
      "options": {
        "triggers": [
          "https://www.youtube.com/watch?v=kZf91MAwS7s",
          "https://www.youtube.com/watch?v=Id-ituXzPvQ",
          "https://www.youtube.com/watch?v=Nq5LMGtBmis"
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
        "responsePrefix": "🚨OH🚨SHIT🚨IT'S🚨A🚨BANGER🚨\n"
      }
    },
    {
      "name": "XKCD",
      "type": "xkcd"
    }
  ]
}
```

In essence, a configuration will have a `botToken` property with (surprise) the Discord bot token to log in with, and a `commands` array which contains an array of command configuration objects, each with a `name` to call it by in logs. a `type` to base the command off of, and an `options` object to actually configure the
command type with.
