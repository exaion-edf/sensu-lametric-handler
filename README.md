### Overview

This plugin provides a lametric handler to send notifications when alerts are raised

### Files
 * bin/sensu-lametric-handler

## Usage example

### Help

**sensu-lametric-handler**

```
The Sensu Go lametric notifications handler

Usage:
  sensu-lametric-handler [flags]
  sensu-lametric-handler [command]

Available Commands:
  help        Help about any command
  version     Print the version number of this plugin

Flags:
  -c, --critical-icon string    The critical state notification icon (default "a2715")
  -C, --critical-sound string   The critical state notification sound (default "negative1")
  -e, --entity-icon string      The entity notification icon (default "i31916")
  -h, --help                    help for sensu-lametric-handler
  -i, --ip string               The lametric ip to send notifications to, defaults to value of SENSU_LAMETRIC_IP env variable
  -k, --key string              The lametric auth key, defaults to value of SENSU_LAMETRIC_KEY env variable
  -o, --ok-icon string          The ok state notification icon (default "a25939")
  -O, --ok-sound string         The ok state notification sound (default "positive1")
  -w, --warning-icon string     The warning state notification icon (default "a7756")
  -W, --warning-sound string    The warning state notification sound (default "negative5")
```