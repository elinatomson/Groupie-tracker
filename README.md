# Groupie-tracker

This project has been made according to the task and its sub-tasks described [here](https://github.com/01-edu/public/tree/master/subjects/groupie-tracker).

Groupie Tracker is a web application designed to interact with a provided [API](https://groupietrackers.herokuapp.com/api), manipulate the data, and display it on a site. You can see bands information  such as name(s), image, the year they began their activity, the date of their first album, members, their last and/or upcoming concert locations and dates.

## Key Features:
* Search Functionality: Allows users to search for specific text inputs within the website, focusing on artists, members, or other attributes from the data system.
* Filtering: Provides users with various filters to refine their search results. The project includes the following four filters:
    - Filter by Creation Date: Refine results based on the date the artist or band was created.
    - Filter by First Album Date: Search based on the date of the artist's or band's first album release.
    - Filter by Number of Members: Filter the list of artists or bands based on their number of members.
    - Filter by Concert Locations: Narrow down the search based on the locations where the artist or band has performed.
* Artist Details View: Clicking on an artist's logo will display detailed information about the artist, providing users with more insights into the selected artist.

## About repository
* main.go file is for server handling.
* Folder <code>templates</code> contains the html files.
* Folder <code>static</code> contains the css file and images.
* Folder <code>functions</code> contains the go files for getting data, the handlers and helper subfunctions. 

## How to use

* Option one with Docker
    - You should have Docker installed. If you don't have, install [Docker](https://docs.docker.com/get-started/get-docker/)
    - To build the image and run the container use following commands:
        - for building the docker image: docker build -t dockerize .
        - for running the docker container: docker run -it -p 8080:8080 dockerize
    - To check the app, open http://localhost:8080 in a browser. 
    - To terminate the server click CTRL + "C".

* Option two directly from your terminal
    - You should have Go installed. If you don't have, install [Go](https://go.dev/doc/install)
    - Type in your terminal: go run main.go
    - Open http://localhost:8080
    - To stop the server, click Ctrl + C in your terminal

## Authors
- [@elinat](https://01.kood.tech/git/elinat)
