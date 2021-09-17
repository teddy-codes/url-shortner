URL Shorter

# What is it?
This a quick program that I had written in go to help someone out on a twitch stream.

Takes a url from a post request and adds it to the database (postgres is the one of choice). Then, queries the database with a path and redirects using
the builtin redirect.

The store structure allows for testing, but at the moment I don't give a crap about doing that on something this small.