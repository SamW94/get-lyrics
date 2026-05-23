# get-lyrics

## Overview 🎵

Welcome to the repository for the get-lyrics project. This tool automates tagging your MP3s with lyrics. It's written in Go and uses the Genius API and a Levenshtein distance algorithm to find the best match for the given artist/track, download the lyrics, and embed them as an ID3v2 tag in the file. 

## Pre-Requisites 📋

1. A git client to clone the repository.
2. The Go toolchain installed to run the program.
3. A [Genius API Access Token](https://docs.genius.com/#/getting-started-h1)
4. Some MP3 files you want to tag.

## How To ✅

1. Clone the repository using git. 

```
git clone https://github.com/SamW94/get-lyrics.git
```

2. Create a `.env` file in the root of the repository with your genius client access token. 

```
cd get-lyrics
echo "GENIUS_CLIENT_ACCESS_TOKEN=<your-client-access-token> >> .env
```

3. Switch into the `src` directory and run the program. The argument provided to the program should be the directory containing the MP3 files (or subdirectories containing MP3 files) you want to tag.

```
cd src
go run . ~/Music/Converge
```

4. Wait for the program to complete. Check your MP3 to see embedded lyrics!