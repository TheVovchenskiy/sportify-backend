#!/usr/bin/env python
import click

from bot import Bot
from config import BotConfig


@click.command
@click.argument("config-file", type=click.Path(exists=True))
def run_bot(config_file: str):
    """
    Run parser bot.

    CONFIG_FILE: Path to config file. It must exist.
    
    Environment variables:

    - BOT_TOKEN: Token of the bot.
    """
    config = BotConfig.get(config_file)
    bot = Bot(config)

    bot.run()


if __name__ == "__main__":
    run_bot()
