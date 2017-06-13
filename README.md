# kagami

Small daemon that helps mirror git repository between various
providers such as github/bitbucket/etc...

__DOES NOT WORK YET!__

## Usage

Copy `config.sample.hcl` to `/etc/kagami.hcl`, change based on
your needs.

```
usage: kagami [<flags>] <command> [<args> ...]

Git mirroring agent

Flags:
      --help           Show context-sensitive help (also try --help-long and --help-man).
  -c, --config="/etc/kagami.hcl"  
                       Configuration file.
  -l, --loglevel=INFO  Log level.

Commands:
  help [<command>...]
    Show help.

  check
    Check if configuration is valid.

  serve
    Start the kagami server.
```
