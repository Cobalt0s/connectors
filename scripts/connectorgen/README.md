
# Connector generator

This script is CLI that allows to generate template to start implementing your connector.

At the beginning generate `Base`, then proceed to `Read`, then in any order `Write/Delete` or `Metadata`.

Every command requires 3 flags:
* `-o --output` Directory where to save all files. Since it overrides files specify temporary directory. Ex: `conn-tmpl-example`
* `-p --package` The name of golang package. Ex: `microsoftdynamicscrm`
* `-n --provider` Catalog name for this provider. Ex: `DynamicsCRM`

## Base

Start with base connector files. These will provide base struct, constructor method, params, etc.

```shell
./bin/cgen base -o microsoftdynamics-example -p msdcrm -n MicrosoftDynamicsCRM
```

# Methods

Every method needs ObjectName argument. 
Manual tests that perform real time requests to a server will request such object. The name is specified in singular form.
Ex: `contact, user, lead, event`. 

## Read 

Sample read method with mock and unit tests.
Test will read `Contacts` from Microsoft APIs.

```shell
./bin/cgen read contact -o microsoftdynamics-example -p msdcrm -n MicrosoftDynamicsCRM
```

## Write+Delete

Sample write and delete methods with mock and unit tests.
Template will provide test where `lead` will be created, updated and then removed.

* Write
```shell
./bin/cgen write lead -o microsoftdynamics-example -p msdcrm -n MicrosoftDynamicsCRM
```
* Delete
```shell
./bin/cgen delete lead -o microsoftdynamics-example -p msdcrm -n MicrosoftDynamicsCRM
```

## Metadata

Sample ListObjectMetadata method with mock and unit tests.
Template will have a manual test which will perform read request on `admin` and then ListObjectMetadata on `admin`.
It will then check properties between them match.

```shell
./bin/cgen metadata admin -o microsoftdynamics-example -p msdcrm -n MicrosoftDynamicsCRM
```
