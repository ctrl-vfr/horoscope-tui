# Horoscope TUI

Because science doesn't explain everything and computers are overrated, here's an app to calculate astrological charts in the terminal and answer your questions.

Calculates planetary positions with beautiful equations and generates a gorgeous zodiac wheel. Let an Oracle explain why Mercury retrograde broke prod again.

![Screenshot](./screenshot.png)

## Installation

```bash
go install github.com/ctrl-vfr/horoscope-tui@latest
```

## Prerequisites

### Compatible Terminal

A terminal supporting the Kitty graphics protocol:

- [Kitty](https://sw.kovidgoyal.net/kitty/)
- [Ghostty](https://ghostty.org/)
- [WezTerm](https://wezfurlong.org/wezterm/)

### resvg

```bash
brew install resvg
# or
cargo install resvg
```

### Environment Variables

```bash
export HOROSCOPE_CITY="Paris, France"
export OPENAI_API_KEY="sk-..."
```

## Localization

The application automatically detects your system locale (`LANG`, `LC_MESSAGES`, or `LC_ALL`) and displays the interface in the corresponding language.

Supported languages:
- English (default)
- French
- Spanish
- German

## Disclaimer

I am not responsible for the consequences of using this program or any decisions you make based on the Oracle's advice.
