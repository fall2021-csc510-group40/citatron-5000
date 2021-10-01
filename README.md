
[![GitHub license](https://img.shields.io/github/license/fall2021-csc510-group40/citatron-5000)](https://github.com/fall2021-csc510-group40/citatron-5000/blob/main/LICENSE)
[![Build Status](https://app.travis-ci.com/fall2021-csc510-group40/citatron-5000.svg?branch=main)](https://app.travis-ci.com/fall2021-csc510-group40/citatron-5000)
[![DOI](https://zenodo.org/badge/408212287.svg)](https://zenodo.org/badge/latestdoi/408212287)
![Golang](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)


<p align="center">
<img src="https://user-images.githubusercontent.com/43625082/135329921-51eeb5d9-b077-4a65-b130-bb4f7c327e53.png" alt="Citatron Logo" style="width:700px;"/>
</p>

# About the Citatron 5000 Project

The Citatron 5000 Project seeks an easier and quicker way to grab paper citations for your projects. It was created for people frustrated with inaccurate, slow, and ad-ridden citation services. The project currently includes an API which searches for the input paper name across multiple data sources including the ACM site, the CrossRef database, and the Citatron 5000 database and returns the cited paper in either plain text IEEE or Bibtex format. This API can be integrated with different extensions/services and is currently used with a telegram bot to demonstrate its functionality.
  
Click on the image below to learn more information about the Citatron Project:

[![Citatron 5000 Video](https://img.youtube.com/vi/Veipwehb6J4/0.jpg)](https://www.youtube.com/watch?v=Veipwehb6J4)  

Or follow [this link to the site](https://fall2021-csc510-group40.github.io/citatron-5000/).

# Project Structure

## Core

The `/core` holds all of the backend of the Citatron 5000 API. This includes the database, formatter, schema, server, source searchers, and util folders.

## Citration 5000 Telegram Bot

The `/telegram-bot` holds a simple Telegram bot as a front-end for the API.

### Bot Usage

The Citatron 5000 Bot is a telegram bot which allows users to quickly cite a paper by name in either plain text IEEE or Bibtex format. Simply send the name of the paper you are looking to cite and choose one of the two formatting options and the Citatron will return a list of 5 citations for papers which most closely match the input name.

### Bot Installation

1. Clone the Citatron 5000 Repository: `git clone git@github.com:fall2021-csc510-group40/citatron-5000.git``
2. Change the working directory to Citatron 5000: `cd citatron-5000`
2. Build and deploy: `docker-compose up --build`
3. **ADD CONFIG INFO HERE**

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