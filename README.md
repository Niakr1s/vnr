# Screenshot

![Screenshot](/screenshot.png?raw=true "Screenshot")

# Needed software

-   [Chromium](https://www.chromium.org/getting-involved/download-chromium). Vnr will use this in headless mode to fetch translations. _Note_: You should add chromium directory to PATH environment variable.
-   [Textractor](https://github.com/Artikash/Textractor). Configure it to extract sentences into system clipboard.
-   [Clipboard Inserter](https://chrome.google.com/webstore/detail/clipboard-inserter/deahejllghicakhplliloeheabddjajm) extension. Don't change default options, just don't forget to enable it every time you launch browser.

# Usage

-   Install needed software.
-   Run precompiled binary from [releases](https://github.com/Niakr1s/vnr/releases) or run with `go run src/main.go`.
-   In your browser open http://localhost:5322.

# Installing as windows service via [nssm](http://nssm.cc)

```
// vnr - name of a service, you can choose anything you want
nssm install vnr c:\vnr.exe
nssm start vnr
```

# Disabling google api warning in chromium

Sometimes you can view ugly notification in chromium about google api warning every time you open it. You can disable such warnings via next commands (in windows os).

```
setx GOOGLE_API_KEY "no"
setx GOOGLE_DEFAULT_CLIENT_ID "no"
setx GOOGLE_DEFAULT_CLIENT_SECRET "no"
```
