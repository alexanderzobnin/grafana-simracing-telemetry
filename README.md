# Simracing telemetry

Simracing Telemetry data source plugin makes it possible to visualize telemetry
data from various simracing titles such as Assetto Corsa Competizione, iRacing 
and others in [Grafana](https://grafana.com/).

<img src="https://user-images.githubusercontent.com/4932851/166692176-6867ccf4-1726-438e-ba52-783696f412b1.png" />

## Features

- Real-time telemetry data visualization
- Highly customizable dashboards

Here's a [demo video](https://vimeo.com/571685229) showing plugin capabilities.

## Supported titles

Currently, plugin supports following games:

- Assetto Corsa Competizione
- Assetto Corsa (experimental)
- iRacing
- Dirt Rally 2.0
- Forza Horizon 5
- Forza Horizon 4 (experimental)

## Supported platforms

Most of simracing titles are intended to be run on Windows. In general, 
Grafana and plugin can be run on any supported platform, since some titles
send telemetry over the network. But Assetto Corsa Competizione and iRacing
provide its telemetry via [Memory-mapped files](https://docs.microsoft.com/en-us/dotnet/standard/io/memory-mapped-files),
so Grafana, plugin and game should be running on the same system (Windows).
Thus, Windows now is the only supported platform.

## Getting started

If you're not familiar with [Grafana](https://grafana.com/) software, [download](https://grafana.com/grafana/download?edition=oss&platform=windows)
and install it first (select OSS edition and Windows platform). Read the Windows [installation guide](https://grafana.com/docs/grafana/latest/installation/windows/).
Refer to [plugin istallation](https://grafana.com/docs/grafana/latest/plugins/installation/) guide and
read how to install plugin. There are few ways how to do it.

### Install pluign

#### Install pluign from package
Go to the [plugin page](https://grafana.com/grafana/plugins/grafana-simracing-telemetry-datasource/)
and select _Installation_ tab. Scroll down and click _Download_ link under the _Install the Data Source_ section.
Unpack it into your grafana plugins directory. Then, restart Grafana server.

#### Install via Plugin catalog
In order to be able to install / uninstall / update plugins using plugin catalog, 
you must enable it via the `plugin_admin_enabled` flag in the [configuration](https://grafana.com/docs/grafana/latest/administration/configuration/#plugin_admin_enabled) file. 
Before following the steps below, make sure you are logged in as a Grafana administrator.

To install a plugin:

1. In Grafana, [navigate to the Plugin catalog](https://grafana.com/docs/grafana/latest/plugins/catalog/#plugin-catalog-entry) to view installed plugins.
2. Browse and find a plugin.
3. Click on the plugin logo.
4. Click Install.

### Create data source

Navigate to the [Configuration -> Data sources](http://localhost:3000/datasources) in Grafana
side menu and click _Add data source_. Select simracing-telemetry-datasource from the list 
and press _Save & test_. You can also enable _Default_ toggle near data source name
to make data source default source when you create new queries.

### Import dashboards

At the data source config page, navigate to _Dashboards_ tab and import
dashboards for the game you want to use.

Then, go to the Grafana [home page](http://localhost:3000/) and select dashboard 
from the dropdown at the top left. Now, you can run the game and start exploring your telemetry.

## Game-specific configs

Some games do not send telemetry by default, so you need to perform some
extra steps to enable it.

### Dirt Rally 2.0

Here's instruction from [Race Department](https://www.racedepartment.com/downloads/dirt-rally-2-0-dashboard-telemetry-tool.26703/): 

1. Locate the config file by going to the following path in windows explorer (Note: this may be different for you, depending on your operating system). Replace “{Yourusername}" with your username.
   `C:\Users\{Yourusername}\Documents\My Games\DiRT Rally 2.0\hardwaresettings`
2. Locate a file called “hardware_settings_config"
3. Make a backup of this file (just in case).
4. Open the file in a text editor (like notepad).
5. Locate the following line:

   `<udp enabled="false" extradata="0" ip="127.0.0.1" port="20777" delay="1" />`
6. Update “enabled" to “true".
7. Update “extradata" to “3"
8. Your new line of code should look like this:

   `<udp enabled="true" extradata="3" ip="127.0.0.1" port="20777" delay="1" />`
9. Save the file and start the game.
