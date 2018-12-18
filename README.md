![valerius](https://i.imgur.com/EsgTukM.png)

Valerius (stylized valerius) is a modular Discord bot for personal use.

## Usage

Valerius uses JSON to describe all of its commands from generic shells.
The program looks for `valerius.json` in the current working directory, but you can set the config file with `valerius -conf <path_to_configfile>`.

In essence, a configuration will have a `botToken` property with (surprise) the Discord bot token to log in with, and a `commands` array which contains an array of command configuration objects, each with a `name` to call it by in logs. a `type` to base the command off of, and an `options` object to actually configure the command type with. Put the config file in the same directory as the executable and fire it up and you should have a working bot you can customize on-the-fly, without need for recompilation or source code editing.

You can also set a log file with the `-log` argument.
