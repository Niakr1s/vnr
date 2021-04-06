# Disabling google api warning

```
setx GOOGLE_API_KEY "no"
setx GOOGLE_DEFAULT_CLIENT_ID "no"
setx GOOGLE_DEFAULT_CLIENT_SECRET "no"
```

# Installing as windows service via [nssm](nssm.cc)

```
// vnr - name of a service, you can choose anything you want
nssm install vnr c:\vnr.exe

// c:\chrome - example path to chrome dir (not to exe!)
nssm set vnr AppEnvironmentExtra PATH=c:\chrome

nssm start vnr
```
