
# dwrk
> A lightweight local/remote project management CLI written in Go.

## Features

-   Create, list, open, and clone local projects
    
-   GitHub integration for cloning and syncing
    
-   Automatic configuration file generation
    
-   Simple and fast CLI workflow for developers

## Requirements

-   Go 1.22+  
    Reference: Go install requires full module path  
    https://go.dev/ref/mod#go-install

## Installation
### Install using `go install`
```bash
go install github.com/okalexiiis/dwrk@latest
```
### Build from source
```bash
git clone https://github.com/okalexiiis/dwrk
cd dwrk
go install .
```

## Initial Configuration

After installation, set your GitHub username:
```bash
dwrk config set github_username your_username
```
The configuration file is created automatically at:
```bash
~/.config/dwrk/config.yaml
```
To inspect current settings:
```bash
dwrk config list
```
If the config file does not exist, dwrk will generate one when needed.

## Usage Examples
### List all projects
```bash
dwrk list
```

### Create a new project
```bash
dwrk new api-server
```

### Clone a GitHub project
```bash 
dwrk clone portfolio-site
```


#### To Do
- [ ] Refactorize the code
