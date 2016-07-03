## Universal BOSH Service Broker

The Universal BOSH Service Broker can provide any software, which can be deployed by BOSH, as a service for Coud Foundry.

### How it works

When a new service instance is created (this happens when a Cloud Foundry operator runs the `cf create-service` command), this broker uses BOSH to deploy the new service. It actually executes the `bosh deploy` command and specifies the deployment manifest, which is generated for this service instance. The manifest is generated from a template. For each service that the BOSH Service Broker should be able to deploy, a manifest template needs to be created and put into the `templates` folder. Each template is actually a regular BOSH manifest, which uses the [Go template syntax](https://golang.org/pkg/text/template/) to define some parameters that will be initialized only at service provision time. For example:

```
instance_groups:
- instances: {{.instances}} 
  name: redis_leader
  vm_type: {{.vm_type}}
  stemcell: {{.stemcell}}
  azs: [z1]
```

This is only a part of the manifest template, but it should illustrate the main idea. The template defines 3 parameters (`instances`, `vm_type`, and `stemcell`). At service provision time, those parameters are initialized with concrete values:

```
cf create-service bosh redis -c '{instances: 2, vm_type: "medium", stemcell: "trusty"}'
```

When a service instance is bound to an application, a special bash script is executed. This script should print to stdout a JSON object. This object should contain credentials that later are passed to the application in the VCAP_SERVICES environment variable. The script is also defined as a template, rendered with the same parameters that were used for rendering the manifest template. The following is an exaple of a bind script:

```
#!/bin/bash -e

host=$(bosh -u '{{.bosh_user}}' -p '{{.bosh_password}}'  -d  {{.deployment_name}} vms | grep 'redis_leader' | awk '{print $11}')

echo "{\"host\": \"$host\", \"password\": \"{{.password}}\", \"port\": 58301 }"
exit 0
```

In this script, we first obtain a value for the IP address of a `redis_leader` VM. This is done using a combination of `bosh vms`, `grep`, and `awk` commands (later, we may provide a more readable way of doing such kind of things). Then, we just print to stdout a JSON object with credentials required to connect to the Redis host. You can also see that there are 4 configuration parameters in this script template: `bosh_user`, `bosh_password`, `deployment_name`, and `password`. The first 3 of them are standard and can be employed in all templates. The `password` parameter is specific to the Redis service plan.

### Broker configuration

Before it can be used, the BOSH Service Broker must be configured with a list of all service plans that it should be able to deploy. The following is an example of a configuration file for the broker:

```
broker_id: 'bosh'
bosh_target: '52.9.99.22'
bosh_user: 'admin'
bosh_password: 'admin'
service_user: 'user'
service_password: 'password'
plans:
  redis:
    name:  redis
    description: "Redis database"
    release: 'https://bosh.io/d/github.com/cloudfoundry-community/redis-boshrelease?v={{.version}}'
    stemcell: "https://bosh.io/d/stemcells/bosh-aws-xen-hvm-ubuntu-trusty-go_agent"
    manifest_template: 'redis.yml.tmpl'
    bind_template: 'redis-bind.sh.tmpl'
    params:
      - name: 'version'
        default: '12'
      - name: 'slave_instances'
        default: 0
      - name: 'password'
        random: true
```

The following properties must be defined in a configuration file:

| Property Name | Description |
| --- | --- |
| broker_id | The broker's ID (must be unique for each Cloud Foundry deployment). This ID is passed to Cloud Foundry in response to the `/v2/catalog` method. See [Service Broker API](https://docs.cloudfoundry.org/services/api.html#catalog-mgmt) for details. |
| bosh_targer | An IP address or a host name of BOSH Director |
| bosh_user | A username that the broker uses to login to BOSH |
| bosh_password | A password that BOSH uses to login to BOSH |
| serivce_user | A user that serves for [service broker auth](https://docs.cloudfoundry.org/services/api.html#authentication) |
| service_password | A password for service broker auth |
| plans | A list of all service plans that the broker should be able to deploy |

#### Service plan configuration

In the configuration file, each service plan is defined as a separate property of a plan's object. The name of this property is used as a plan-unique ID. Each plan should contain the following properties.

| Property Name | Description |
| --- | --- |
| name | The name of the service plan |
| description | A description of the service plan |
| release | BOSH release URL. This is actually a template, so things, like {{.version}}, can be used in this URL. |
| stemcell | BOSH stemcell. This is also a template. |
| manifest_template | The manifest template file. This file should be located in the templates folder. |
| bind_template | Bind script template |
| unbind_template | Unbind script template |
| params | A collection of all paramters that can be passed to the broker at provisioning time. |

#### Parameter configuration

Each parameter that can be passed to the service broker at provisioning time should be added to a params collection. The following properties can be specified for each parameter.

| Property Name | Description |
| --- | --- |
| name | Parameter name. The same name should be used when a parameter is specifed in the `cf create-service` command. The parameter will be referenced by this name when being exposed to all templates. |
| default | A default value for a parameter. |
| random | Specifies whether the broker should generate a random string for a parameter value, in case this value is not specified by the user. |

### Installing the service broker

The following are the steps for deploying the BOSH Service Broker.

1. Clone the service broker repository:
  ```
  git clone https://github.com/s-matyukevich/bosh-broker
  cd bosh_broker
  ```

1. Modify the `config.yml` file and add all the necessary templates to the `templates` folder.

1. Push the service broker to Cloud Foundry as an application.u (For now, this is the only avaliable option. Later, a separate BOSH release will be provided to deploy this broker to BOSH directly.)
  ```
  cf push bosh-broker -b https://github.com/cloudfoundry/go-buildpack
  ```

1. Add the service broker to Cloud Foundry:
  ```
  cf create-service-broker bosh <user> <password> <broker-url>
  cf enable-service-access bosh
  ```

1. Create a service instance:
  ```
  cf create-service  bosh <plan-name> <instance-name> -c <plan-parameters>
  ```

1. Bind the service to an application:
  ```
  cf bind-service <app-name> <instance-name>
  ```
 
