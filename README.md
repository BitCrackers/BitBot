# BitBot [![Discord](https://img.shields.io/badge/Discord-Invite-7289DA.svg?logo=Discord&style=flat-square)](https://discord.gg/AUpXd3VUh8) [![Paypal](https://img.shields.io/badge/PayPal-Donate-Green.svg?logo=Paypal&style=flat-square)](https://www.paypal.com/donate/?hosted_button_id=TYMU92FD9D9UW)

<p align="center">
    BitBot is the open source utility and moderation bot for BitCrackers. Written in Go and made with ❤️.
</p>

## Getting Started

*Disclaimer: This assumes you are planning on running this bot on macOS or Linux.*

1. Install [Go](https://golang.org) and get it properly configured.
    1. [This](https://github.com/alco/gostart#faq10) is a great resource for learning how Go works.
2.

``` bash
 $ go get -u github.com/BitCrackers/BitBot
 $ go install github.com/BitCrackers/BitBot
```

3. `go install` will produce the BitBot binary in `$GOPATH/bin`. The executable will depend on a handful of environment variables being properly exported, the list of which you can find in the next section.

### Environment Variables

BitBot checks for a number of environment variables. The current list is:

* `$BITBOT_TOKEN`
* `$BITBOT_DEBUG`
* `$BITBOT_OWNERID`

### TOML Config

In addition to the required environment variables, BitBot has a TOML based config file (`config.toml`) that is generated on the first run of the bot. None of the values in the config file are required, but taking the time to go through them can expand the functionality of the bot, such as setting up moderation users, muted roles, and more!


## Release Candidate 1 Milestones

* [ ] TOML config for extra functionality.
* [ ] BITBOT_DEBUG environment variable actually doing something.
* [ ] Moderation commands.
    * [ ] Kicking
    * [ ] Banning
    * [ ] Muting
    * [ ] Temp-ban and Temp-mute
* [ ] sqlite3 database integration
* [ ] Moderation logging in channel
* [ ] Parsing messages for common phrases (listed in config.toml probably)
* [ ] Possible: Pull build artifacts for menu every update and send to channel

## Contributing

1. Fork it ([https://github.com/BitCrackers/BitBot/fork](https://github.com/BitCrackers/BitBot/fork))
2. Create your feature branch (`git checkout -b feature/fooBar`)
3. Commit your changes (`git commit -am 'Add some fooBar'`)
4. Push to the branch (`git push origin feature/fooBar`)
5. Create a new Pull Request