# Webpacking

Webpacking is influenced by Rails Webpacker project, and provides similar
capabilities.

In dev mode, webpacking runs webpack in a background process and
generates script/style tags for its entry points.

In production mode, webpacking reads the manifest.json file and serves
it's pre-built assets from the same tags.
