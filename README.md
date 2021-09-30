[![Build Status](https://app.travis-ci.com/fall2021-csc510-group40/Project-1.svg?branch=main)](https://app.travis-ci.com/fall2021-csc510-group40/Project-1)
[![DOI](https://zenodo.org/badge/408212287.svg)](https://zenodo.org/badge/latestdoi/408212287)


<p align="center">
<img src="https://user-images.githubusercontent.com/43625082/135329921-51eeb5d9-b077-4a65-b130-bb4f7c327e53.png" alt="Citatron Logo" style="width:700px;"/>
</p>

# About the Citatron 5000 Project
The Citatron 5000 Project seeks an easier and quicker way to grab paper citations for your projects. The project currently includes an API which searches for the input paper name across multiple data sources including the ACM site, the CrossRef database, and the Citatron 5000 database and returns the cited paper in either plain text IEEE or Bibtex format. This API can be integrated with different extensions/services and is currently used with a telegram bot to demonstrate its current functionality.

# Citatron 5000 Bot Usage
The Citatron 5000 Bot is a telegram bot which allows users to quickly cite a paper by name in either plain text IEEE or Bibtex format. Simply send the name of the paper you are looking to cite and choose one of the two formatting options and the Citatron will return a list of 5 citations for papers which most closely match the input name.

# Citatron 5000 Bot Installation
1. Clone the Citatron 5000 Repository
2. From the root directory run: `docker-compose up --build`
3. **ADD CONFIG INFO HERE**

# Working with the Citatron 5000 API
The core holds all of the backend of the Citatron 5000 API. This includes the database, formatter, schema, server, source searchers, and util folders.

## Elements of Citatron Core
**Update READMEs for every core folder**  
* [db](/docs/db.md)  
* [formatter](/docs/formatter.md)  
* [schema](/docs/schema.md)  
* [server](/docs/server.md)  
* [source](/docs/source.md) 
* [util](/docs/util.md)  

# Tests
**Add info about tests**

# Documentation
**Add info about documentation if we want some**

# Funding

There is no funding.

# Deprecation

For any deprecation, we will pin an announcement to the README. 