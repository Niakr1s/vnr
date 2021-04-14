# Screenshot

![Screenshot](/screenshot.png?raw=true "Screenshot")

# Needed software

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

// c:\chrome - example path to chrome dir (not to exe!)
nssm set vnr AppEnvironmentExtra PATH=c:\chrome

nssm start vnr
```
