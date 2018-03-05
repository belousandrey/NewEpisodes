# new-episodes
Application that checks podcasts for new published episodes. In case of new episodes application notifies user by email.

Some of podcasts has their RSS and we can easily receive updates. But some of them don't, so we have to parse html to get links to files.

## Installation
    go get github.com/belousandrey/new-episodes

## Compilation
    make build

## Usage
Create your own config file, add/remove podcasts that you want to be notified. You can add podcasts that use supported engines.
 
Add record to cron to start app once an hour (or with your frequency).

    0 * * * * /path/to/binary/new-episodes --config /path/to/conf/myconfig.yaml

# Supported podcasts engines
At this point app can work with these podcast engines:
* [golangshow](golangshow.com) - single podcast site
* [changelog](https://changelog.com/) - podcast platform for developers 
* [rucast](https://radio-t.com/) - service used by Radio-T podcast
* [podfm](http://podfm.ru) - podcast platform
* [podster](https://podster.fm) - podcast platform
* [matchdaybiz](https://matchday.biz) - single podcast site (does not use RSS)

## Config format
Top level item `podcasts` is an array of structs with four fields.
- `last` date of last episode in `YYYY-MM-DD` format (ex. 2017-11-29)
- `link` link to RSS feed or page with episodes list (ex. https://site.com/podcast/feed.rss)
- `title`: human readable podcast name (ex. Awesome Podcast)
- `engine`: string with one of supported podcast engines name (ex. changelog)

Top level item `email` is a struct with two fields.
- `to` string with email address where to send email (ex. mike@gmail.com)
- `from` is a struct with five fiels
- - `smtp`: string with SMTP domain where to send email from (ex. smtp.gmail.com) 
- - `port`: port number where to send email from (ex. 587)
- - `username`: string with username where to send email from (ex. noreplysender)
- - `domain`: string with domain where `username` is registered (ex. gmail.com)
- - `password`: sring with password from account `username@domain`

# Third-party libraries
* [go-gomail/gomail](https://github.com/go-gomail/gomail)
* [mmcdole/gofeed](https://github.com/mmcdole/gofeed)
* [pkg/errors](https://github.com/pkg/errors)
* [PuerkitoBio/goquery](https://github.com/PuerkitoBio/goquery)
* [spf13/pflag](https://github.com/spf13/pflag)
* [spf13/viper](https://github.com/spf13/viper)

## TODO
* tests