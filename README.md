# Citatron 5000

[![GitHub license](https://img.shields.io/github/license/fall2021-csc510-group40/citatron-5000)](https://github.com/fall2021-csc510-group40/citatron-5000/blob/main/LICENSE)
[![Build Status](https://app.travis-ci.com/fall2021-csc510-group40/citatron-5000.svg?branch=main)](https://app.travis-ci.com/fall2021-csc510-group40/citatron-5000)
[![DOI](https://zenodo.org/badge/408212287.svg)](https://zenodo.org/badge/latestdoi/408212287)
![Golang](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)


<p align="center">
<img src="https://user-images.githubusercontent.com/43625082/135329921-51eeb5d9-b077-4a65-b130-bb4f7c327e53.png" alt="Citatron Logo" style="width:700px;"/>
</p>

The Citatron 5000 Project seeks an easier and quicker way to grab paper citations for your projects. It was created for people frustrated with inaccurate, slow, and ad-ridden citation services. The project currently includes an API which searches for the input paper name across multiple data sources including the ACM site, the CrossRef database, and the Citatron 5000 database and returns the cited paper in either plain text IEEE or Bibtex format. This API can be integrated with different extensions/services and is currently used with a telegram bot to demonstrate its functionality.
  
Click on the image below to learn more information about the Citatron Project:

[![Citatron 5000 Video](https://img.youtube.com/vi/Veipwehb6J4/0.jpg)](https://www.youtube.com/watch?v=Veipwehb6J4)  

Or follow [this link to the site](https://fall2021-csc510-group40.github.io/citatron-5000/).

## Project Structure

### Core

The `/core` holds all of the backend of the Citatron 5000 API. This includes the database, formatter, schema, server, source searchers, and util folders.

### Citration 5000 Telegram Bot

The `/telegram-bot` holds a simple Telegram bot as a front-end for the API. The bot allows users to quickly cite a paper by name in either plain text IEEE or Bibtex format. 
Simply send the title of the paper you are looking to cite, select the specific paper that you were looking for using the provided buttons, then choose the format, and enjoy
the citation.

## Deployment

Currently, the core with the API and the bot are deployed together using the provided docker-compose config. To deploy them, you will need to have Docker and docker-compose on your machine:
```bash
sudo apt update
sudo apt install docker docker-compose
```

To deploy the service, you need to first clone the repository to a location of your choice:
```bash
$ git clone git@github.com:fall2021-csc510-group40/citatron-5000.git
$ cd citatron-5000
```

After that, you need to configure the components. There are several points of configuration that allow for easy deployment with default values
as well as more fine-tuning if you wish to deploy the components manually:

* `mongo-init.js`: here you can change the database setup, in particular, the user configs, including the user for the core service
* `core/config.json`: that config is currently responsible for the database connection, timeouts associated with the search as well as log level for the search part of the service
* `telegram-bot/config.json`: this file is used to configure the bot, allowing to change the token, update the available formatting options and change the log level
    - to get the token for the bot, talk to the @BotFather in Telegram
* `docker-compose.yml`:

#### Configuring the docker-compose.yml

Main points of configuration for the easy deployment option are in the database section, allowing to set the login and password for the mongo-express:
```yaml
mongo-express:
    # ...
    environment:
        # ...
        ME_CONFIG_BASICAUTH_USERNAME: <YOUR_USERNAME>
        ME_CONFIG_BASICAUTH_PASSWORD: <YOUR_PASSWORD>
```

By default, these options are absent, allowing anyone to access the database admin panel.

If you wish to expose your components to the outside world, you will definitely have to update default root password as well and potentially setup nginx to aid you
in your deployment.

---

Finally, after the configuration is done, you can start your setup:
```bash
$ docker-compose up --build -d
```

And to access its logs:
```bash
$ docker-compose logs -f
```

## Tests

Travis tests are automatically run each commit. To test the Telegram bot manually, see the demo video.

## Documentation

For extra documentation, including API and troubleshooting, click [here](https://github.com/fall2021-csc510-group40/citatron-5000/tree/main/docs).

# Funding

There is no funding.

# Deprecation

For any deprecation, we will pin an announcement to the README.

# Contributors

Citatron 5000 is currently a standalone project, and is maintained by 5 users. We would also like to give special thanks to Crossref for providing a convenient API for gathering sources.

# Support

If you would like to contact us about an issue, please submit the issue via the issues tab in the Github repo and our support team will look into it.

# Success Stories

As Citatron 5000 is introduced to the public, real-world stories of its use will be included in this section.

# Licensing

MIT License

Copyright (c) 2021 fall2021-csc510-group40

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
