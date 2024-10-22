#!/usr/bin/env python
import os
from dataclasses import dataclass

import yaml


@dataclass
class BotConfig:
    """
    Base bot config class. Now by calling `BotConfig.get(config_path)` static
    method provided config file as well as environment variables will be parsed.
    """

    token: str  # Bot token, defined only in environment variable BOT_TOKEN
    api_url: str  # API url, must contain 'http://' or 'https://', defined only in config file

    @classmethod
    def get(cls, config_path: str) -> "BotConfig":
        if os.path.exists(config_path):
            yaml_config = BotConfig._from_yaml(config_path)
        else:
            raise FileNotFoundError(f"Config file not found: {config_path}")

        env_config = BotConfig._from_env()
        if "app" not in yaml_config or "api_url" not in yaml_config["app"]:
            raise ValueError("api_url is not defined in config file")

        api_url = yaml_config["app"]["api_url"]
        token = env_config["token"]

        config = {**yaml_config["bot"], **env_config}

        return cls(token, api_url)

    @staticmethod
    def _from_yaml(config_path: str) -> dict:
        with open(config_path, "r") as f:
            config = yaml.safe_load(f)
        return config

    @staticmethod
    def _from_env() -> dict:
        token = os.environ.get("BOT_TOKEN", None)
        if token is None:
            raise ValueError("BOT_TOKEN environment variable is not set")

        config = {"token": token}
        return config
