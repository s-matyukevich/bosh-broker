## Universal bosh service broker

Universal bosh service broker can provide any software, that can be deployed by BOSH, as a service for Coud Foundry.

### How it works

When new service instance is creaded (this happens when Cloud Foundry operator executes `cf create-service` command) this broker uses BOSH to deploy new service. It actually executes `bosh deploy` command and specifies the deployment manifest, that is generated for this service instance. This manifest is generated from a template. For each service, that bosh service broker should be able to deploy, manifest template should be created and put into `templates` folder.Each template is actually an usual BOSH manifest, that uses [Go template syntax](https://golang.org/pkg/text/template/) to define some parameters, that will be initialized only at service provision time. For example:

```
instance_groups:
- instances: {{.instances}} 
  name: redis_leader
  vm_type: {{.vm_type}}
  stemcell: {{.stemcell}}
  azs: [z1]
```

This is only a part of manifest template, but it should clearly show the main idea. This template defines 3 parameters (`instances`, `vm_type` and `stemcell`) Those parameters are initialized at service provision time with concrete values:

```
cf create-service bosh redis -c '{instances: 2, vm_type: "medium", stemcell: "trusty"}'
```

When service instance is bound to an application, special bash script is executed. This script should print to stdout a json object. This object should contain credentials that later are passed to the application in VCAP_SERVICES environment valuable. This script is also defined as a template and it is rendered using the same parameters, that were used for rendering manifest template. The following is an exaple of a bind script:

```
#!/bin/bash -e

host=$(bosh -u '{{.bosh_user}}' -p '{{.bosh_password}}'  -d  {{.deployment_name}} vms | grep 'redis_leader' | awk '{print $11}')

echo "{\"host\": \"$host\", \"password\": \"{{.password}}\", \"port\": 58301 }"
exit 0
```

In this script we first obtain value for the ip address of a redis_leader vm. This is done using a combination of `bosh vms`, `grep` and `awk` commands (probably later we will provide a more readable way of doing such kind of things).  Then we just print to stdout json object with credentials, that are required to connect to redis host. You can also see that in this script template 4 configuration parameters are used (`bosh_user`, `bosh_password`, `deployment_name` and `password`) The first 3 of them are standart and can be used in all templates. `password` parameter is specific to redis service plan.

### Broker configuration

Before using, bosh service broker must be configured with list of all service plan, that broker should deploy. The following is an example of a configuration file for the broker:

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
| broker_id | Broker Id (Must be unique per Cloud Foundry deployment). This Id is passed to Cloud Foundry in response to `/v2/catalog` method. See [Service broker API](https://docs.cloudfoundry.org/services/api.html#catalog-mgmt) for details. |
| bosh_targer | IP address or host mane of BOSH director |
| bosh_user | Username that broker uses to login to BOSH |
| bosh_password | Password that BOSH uses to login to BOSH |
| serivce_user | User that is used for [service broker auth](https://docs.cloudfoundry.org/services/api.html#authentication) |
| service_password | Password for service broker auth |
| plans | List of all service plans, that broker should be able to deploy |

#### Service Plan configuration

In configuration file, each service plan is defined as a separate property of a plans object. The name of this property is used as plan unique id. Each plan should contina the following properties.

| Property Name | Description |
| --- | --- |
| name | Name of the service plan |
| description | Description of the service plan |
| release | BOSH release url. This is actually a template, so things like {{.version}} can be used in this url |
| stemcell | BOSH stemcell. This is also a template |
| manifest_template | Manifest template file. This file should be located in a templates folder |
| bind_template | Bind script template |
| unbind_template | Unbind script template |
| params | Collection of all paramters, that can be passed to broker at provision time |

#### Parameter configuration

Each parameter, that can be passed to service broker at provision time, should be added to params collection. The following properties can be specified for each parameter.

| Property Name | Description |
| --- | --- |
| name | Parameter name. The same name should be used when parameter is specifed in `cf create-service` command. Also, with the same name this parameter will be exposed to all templates |
| default | Default value for a parameter. |
| random | Specifies, whether broker should generate random string for a parameter value, in case when this value is not specified by the user. |

### Installing service broker

The following steps should be made in order to deploy bosh service broker.

1. Clone service broker repository
  ```
  git clone https://github.com/s-matyukevich/bosh-broker
  cd bosh_broker
  ```

1. Modify `config.yml` file and add all necessary templates to `templates` folder.

1. Push service broker to cloud foundry as an application (For now this is the only avaliable option. Later sepparate BOSH release will be provided to deploy this broker to BOSH directly.)
  ```
  cf push bosh-broker
  ```

1. Add service broker to Cloud Foundry
  ```
  cf create-service-broker bosh <user> <password> <broker-url>
  ```

1. Create service instance
  ```
  cf create-service <instance-name> bosh <plan-name> -c <plan-parameters>
  ```

1. Bind service to an application
  ```
  cf bind-service <app-name> <instance-name>
  ```
 
