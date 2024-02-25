# Powershell Core Mixin for Porter

This is a mixin for Porter that provides the Powershell Core.

## Install or Upgrade

Currently we only support the installation via `--url`. Please make sure to install the mixin as follow:

```
porter mixin install pwsh --url https://github.com/kichristensen/pwsh-mixin/releases/download
```

or for a specific version

```
porter mixin install pwsh --version VERSION --url https://github.com/kichristensen/pwsh-mixin/releases/download
```

## Mixin Configuration

### Client version

By default, the most recent version of Powershell Core is installed. You can specify a specific version with the
`clientVersion` setting.

```yaml
mixins:
- pwsh:
   clientVersion: 7.4.0
   
```

### Modules

When you declare the mixin, you can also configure additional modules to install.

#### Use vanilla Powershell Core

```yaml
mixins:
- pwsh
```

#### Install additional modules

By default, additional modules are installed in the most recent version. You can specify a specific version with the
`version` setting.

When installing without a specific version, it is possible to specify the name of the module as a string instead of an object.

NOTE: The version format is the same as the format in PSResourceGet,
[Searching by NuGet version ranges](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.psresourceget/about/about_psresourceget?view=powershellget-3.x#searching-by-nuget-version-ranges)

```yaml
mixins:
- pwsh:
   modules:
   - Microsoft.Graph
   - name: Microsoft.Graph.Beta
     version: 2.12.0
```

## Mixin Syntax

The mixin support both inline scripts and files.

```yaml
pwsh:
  description: Description of the script
  inlineScript: |- # The Powershell script to run
    POWERSHELL SCRIPT
  file: FILEPATH # The .ps1 file to execute
  arguments: # The arguments to pass on to the script
  - ARGUMENT
  suppress-output: false
  ignoreError: # Conditions when execution should continue even if the command fails
    all: true # Ignore all errors
    exitCodes: # Ignore failed commands that return the following exit codes
    - 1
    - 2
    output: # Ignore failed commands based on the contents of stderr
      contains: # Ignore when stderr contains a substring
      - "SUBSTRING IN STDERR"
      regex: # Ignore when stderr matches a regular expression
      - "GOLANG_REGULAR_EXPRESSION"      
    outputs: # Collect values from the script and make it available as an output
      - name: NAME
        jsonPath: JSONPATH # Scrape stdout with a json path expression
      - name: NAME
        regex: GOLANG_REGULAR_EXPRESSION
      - name: NAME
        path: FILEPATH # Save the contents to a file
```

### Suppress Output

See [Exec Mixin - Suppress Output](https://porter.sh/mixins/exec/#suppress-output).

### Ignore Error

The mixin supports the same functionality as the [Exec Mixin - Ignore Error](https://porter.sh/mixins/exec/#ignore-error).

### Outputs

The mixin supports the same output types as the [Exec Mixin - Outputs](https://porter.sh/mixins/exec/#outputs).

## Examples

### Install Microsoft.Graph.Module

```yaml
mixins:
- pwsh:
  modules:
  - name: Microsoft.Graph
```
### Install Microsoft.Graph.Module in specific version

```yaml
mixins:
- pwsh:
  modules:
  - name: Microsoft.Graph
    version: 2.12.0
```

# Run inline script

```yaml
pwsh:
  description: "Print Victory"
  inlineScript: |-
    Write-Host "VICTORY"
```
# Run inline script with arguments

```yaml
pwsh:
  description: "Print Victory"
  inlineScript: |-
    Write-Host "VICTORY to $($args[0])"
  arguments:
  - Porter
```
# Run file script

```yaml
pwsh:
  description: "Print Victory"
  file: ./script.ps1
```
# Run inline script with arguments

```yaml
pwsh:
  description: "Print Victory"
  file: ./script.ps1
  arguments:
  - Porter
```
