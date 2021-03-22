# Word Define

A cli dictionary client for the Oxford English Dictionary.

## Video

[https://youtu.be/RNW_4nCzG0A](https://youtu.be/RNW_4nCzG0A)

## Installation

1. Install Go and run `go get github.com/tom-on-the-internet/word-define`. Alternatively, download the latest release of word-define from GitHub.

2. Run the command once to create config. `word-define example`

3. Get a _free_ `AppID` and `AppKey` from [https://developer.oxforddictionaries.com/](https://developer.oxforddictionaries.com/). This allows you 1000 definition lookups a month. Quick math: about 30 definitions a day.

4. Add your `AppID` and `AppKey` to `$HOME/.config/word-define`.

5. Done!

## Usage

`word-define your-search-term`

Ex:

`word-define portent`

results in

```
[ PORTENT ]

(1)
DEFINITION: a sign or warning that a momentous or calamitous event is likely to happen
EXAMPLES: many birds are regarded as being portents of death
ETYMOLOGIES: late 16th century: from Latin portentum ‘omen, token’, from the verb portendere (see portend)

(2)
DEFINITION: an exceptional or wonderful person or thing
EXAMPLES: many birds are regarded as being portents of death
ETYMOLOGIES: late 16th century: from Latin portentum ‘omen, token’, from the verb portendere (see portend)
```

## Caching

Caching definitions is off by default.

To turn on caching, update the config `$HOME/.config/word-define` JSON and set `cache` to `true`.

Please note: It is your responsibility to make sure you are adhering to the terms and services of the Oxford English Dictionary regarding caching.
